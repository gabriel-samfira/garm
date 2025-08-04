package handlers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/cloudbase/garm/auth"
	"github.com/cloudbase/garm/params"
)

type AuthHandler struct {
	runner        *WebHandler // Reuse the WebHandler struct that has the runner
	templates     *template.Template
	authenticator *auth.Authenticator
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token       string    `json:"token"`
	ExpiresAt   time.Time `json:"expires_at"`
	User        User      `json:"user"`
	RedirectURL string    `json:"redirect_url,omitempty"`
}

type User struct {
	Username string `json:"username"`
	IsAdmin  bool   `json:"is_admin"`
}

type InitRequest struct {
	Username        string `json:"username"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}

func NewAuthHandler(webHandler *WebHandler, authenticator *auth.Authenticator) (*AuthHandler, error) {
	h := &AuthHandler{
		runner:        webHandler,
		authenticator: authenticator,
	}

	// Load auth-specific templates
	if err := h.loadAuthTemplates(); err != nil {
		return nil, fmt.Errorf("loading auth templates: %w", err)
	}

	return h, nil
}

// NewSimpleAuthHandler creates a simple auth handler for testing
func NewSimpleAuthHandler() *AuthHandler {
	return &AuthHandler{
		// No authenticator for simple handler - will need to be set later
	}
}

// SetAuthenticator sets the authenticator for the auth handler
func (h *AuthHandler) SetAuthenticator(authenticator *auth.Authenticator) {
	h.authenticator = authenticator
}

func (h *AuthHandler) loadAuthTemplates() error {
	// For now, reuse the main template loader
	h.templates = h.runner.templates
	return nil
}

// LoadAuthTemplates is a public method to load templates from an existing WebHandler
func (h *AuthHandler) LoadAuthTemplates(webHandler *WebHandler) error {
	h.runner = webHandler
	h.templates = webHandler.templates
	return nil
}

// LoginPageHandler serves the login page
func (h *AuthHandler) LoginPageHandler(w http.ResponseWriter, r *http.Request) {
	slog.Info("LoginPageHandler called")

	// Check if templates are loaded
	if h.templates == nil {
		slog.Error("Templates not loaded in auth handler")
		http.Error(w, "Templates not loaded", http.StatusInternalServerError)
		return
	}

	data := PageData{
		Title:      "Login",
		PageTitle:  "Login to GARM",
		ActivePage: "login",
	}

	slog.Info("Executing login-base template with login content")
	if err := h.templates.ExecuteTemplate(w, "login-base.html", data); err != nil {
		slog.Error("Failed to execute login-base template", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	slog.Info("Login-base template executed successfully")
}

// InitPageHandler serves the initialization page for first-time setup
func (h *AuthHandler) InitPageHandler(w http.ResponseWriter, r *http.Request) {
	// Check if system is already initialized
	ctx := r.Context()
	_, err := h.runner.runner.GetControllerInfo(ctx)
	if err == nil {
		// System is initialized, redirect to login
		http.Redirect(w, r, "/web/login", http.StatusFound)
		return
	}

	data := PageData{
		Title:      "Initialize GARM",
		PageTitle:  "Initialize GARM",
		ActivePage: "init",
	}

	if err := h.templates.ExecuteTemplate(w, "init-base.html", data); err != nil {
		slog.Error("Failed to execute init-base template", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// LoginHandler handles login API requests
func (h *AuthHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var loginReq LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil {
		slog.ErrorContext(ctx, "Failed to decode login request", "error", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Proxy to the REST API login endpoint to ensure compatibility
	apiLoginReq := params.PasswordLoginParams{
		Username: loginReq.Username,
		Password: loginReq.Password,
	}

	apiLoginJSON, err := json.Marshal(apiLoginReq)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to marshal API login request", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Make internal request to the REST API login endpoint
	resp, err := http.Post("http://localhost:9997/api/v1/auth/login", "application/json", strings.NewReader(string(apiLoginJSON)))
	if err != nil {
		slog.ErrorContext(ctx, "Failed to call REST API login", "error", err)
		http.Error(w, "Authentication service unavailable", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// Forward the error from the REST API
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(resp.StatusCode)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid credentials"})
		return
	}

	// Parse the REST API response
	var apiResp params.JWTResponse
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&apiResp); err != nil {
		slog.ErrorContext(ctx, "Failed to parse REST API token response", "error", err)
		http.Error(w, "Authentication service error", http.StatusInternalServerError)
		return
	}

	tokenString := apiResp.Token
	expiresAt := time.Now().Add(24 * time.Hour) // Use reasonable default

	// Set JWT token as HTTP-only cookie for web interface authentication
	http.SetCookie(w, &http.Cookie{
		Name:     "garm_token",
		Value:    tokenString,
		Path:     "/",
		HttpOnly: true,
		Secure:   r.TLS != nil,
		SameSite: http.SameSiteLaxMode,
		Expires:  expiresAt,
	})

	// Check if there's a return URL to redirect to
	returnURL := "/web/"
	if cookie, err := r.Cookie("garm_return_url"); err == nil && cookie.Value != "" {
		returnURL = cookie.Value
		// Clear the return URL cookie
		http.SetCookie(w, &http.Cookie{
			Name:     "garm_return_url",
			Value:    "",
			Path:     "/",
			HttpOnly: true,
			Secure:   r.TLS != nil,
			SameSite: http.SameSiteLaxMode,
			MaxAge:   -1,
		})
	}

	response := LoginResponse{
		Token:     tokenString,
		ExpiresAt: expiresAt,
		User: User{
			Username: loginReq.Username, // Use the original username
			IsAdmin:  true,              // Default to admin for now
		},
		RedirectURL: returnURL,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// InitHandler handles system initialization
func (h *AuthHandler) InitHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var initReq InitRequest
	if err := json.NewDecoder(r.Body).Decode(&initReq); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Validate input
	if initReq.Username == "" || initReq.Password == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Username and password are required"})
		return
	}

	if initReq.Password != initReq.ConfirmPassword {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Passwords do not match"})
		return
	}

	if len(initReq.Password) < 8 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Password must be at least 8 characters long"})
		return
	}

	// Initialize the system using GARM's NewUserParams
	newUserParams := params.NewUserParams{
		Username: initReq.Username,
		Password: initReq.Password,
		Email:    fmt.Sprintf("%s@garm.local", initReq.Username), // Default email
		FullName: initReq.Username,
		IsAdmin:  true,
		Enabled:  true,
	}

	ctx := r.Context()
	slog.Info("System initialization request", "username", newUserParams.Username)

	// Call the actual InitController function
	if h.authenticator == nil {
		slog.Error("Authenticator not available for initialization")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Authentication service not available"})
		return
	}

	user, err := h.authenticator.InitController(ctx, newUserParams)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to initialize controller", "error", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to initialize system: " + err.Error()})
		return
	}

	slog.Info("System initialized successfully", "username", user.Username, "user_id", user.ID)

	// After successful initialization, log the user in by creating a JWT token
	// We'll use the same approach as the login handler
	apiLoginReq := params.PasswordLoginParams{
		Username: initReq.Username,
		Password: initReq.Password,
	}

	apiLoginJSON, err := json.Marshal(apiLoginReq)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to marshal API login request after init", "error", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to complete initialization"})
		return
	}

	// Make internal request to the REST API login endpoint
	resp, err := http.Post("http://localhost:9997/api/v1/auth/login", "application/json", strings.NewReader(string(apiLoginJSON)))
	if err != nil {
		slog.ErrorContext(ctx, "Failed to call REST API login after init", "error", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "System initialized but login failed"})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		slog.ErrorContext(ctx, "REST API login failed after init", "status", resp.StatusCode)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "System initialized but login failed"})
		return
	}

	// Parse the REST API response
	var apiResp params.JWTResponse
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&apiResp); err != nil {
		slog.ErrorContext(ctx, "Failed to parse REST API token response after init", "error", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "System initialized but login failed"})
		return
	}

	tokenString := apiResp.Token
	expiresAt := time.Now().Add(24 * time.Hour)

	// Set JWT token as HTTP-only cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "garm_token",
		Value:    tokenString,
		Path:     "/",
		HttpOnly: true,
		Secure:   r.TLS != nil,
		SameSite: http.SameSiteLaxMode,
		Expires:  expiresAt,
	})

	response := LoginResponse{
		Token:     tokenString,
		ExpiresAt: expiresAt,
		User: User{
			Username: user.Username,
			IsAdmin:  user.IsAdmin,
		},
		RedirectURL: "/web/",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// LogoutHandler handles logout requests
func (h *AuthHandler) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// Clear the JWT token cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "garm_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   r.TLS != nil,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1, // Delete the cookie
	})

	// Also clear the return URL cookie if it exists
	http.SetCookie(w, &http.Cookie{
		Name:     "garm_return_url",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   r.TLS != nil,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1, // Delete the cookie
	})

	// For JWT, we just need to tell the client to remove the token
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Logged out successfully"})
}
