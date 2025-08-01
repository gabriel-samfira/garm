package middleware

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"

	runnerErrors "github.com/cloudbase/garm-provider-common/errors"
	"github.com/cloudbase/garm/auth"
	"github.com/cloudbase/garm/config"
	dbCommon "github.com/cloudbase/garm/database/common"
)

// WebAuthMiddleware is a web-specific authentication middleware that redirects
// to login page instead of returning JSON errors
type WebAuthMiddleware struct {
	store dbCommon.Store
	cfg   config.JWTAuth
}

// NewWebAuthMiddleware creates a new web authentication middleware
func NewWebAuthMiddleware(store dbCommon.Store, cfg config.JWTAuth) *WebAuthMiddleware {
	return &WebAuthMiddleware{
		store: store,
		cfg:   cfg,
	}
}

// Middleware implements the auth.Middleware interface for web routes
func (w *WebAuthMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		
		// Try to get JWT token from Authorization header first
		token := w.getTokenFromHeader(r)
		
		// If no header token, try to get from cookie
		if token == "" {
			token = w.getTokenFromCookie(r)
		}
		
		if token == "" {
			w.redirectToLogin(rw, r)
			return
		}
		
		claims, err := w.validateToken(token)
		if err != nil {
			slog.With(slog.Any("error", err)).DebugContext(ctx, "token validation failed")
			w.redirectToLogin(rw, r)
			return
		}
		
		ctx, err = w.claimsToContext(ctx, claims)
		if err != nil {
			slog.With(slog.Any("error", err)).DebugContext(ctx, "failed to add claims to context")
			w.redirectToLogin(rw, r)
			return
		}
		
		if !auth.IsEnabled(ctx) {
			slog.DebugContext(ctx, "user account is disabled")
			w.redirectToLogin(rw, r)
			return
		}
		
		next.ServeHTTP(rw, r.WithContext(ctx))
	})
}

// getTokenFromHeader extracts JWT token from Authorization header
func (w *WebAuthMiddleware) getTokenFromHeader(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return ""
	}
	
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return ""
	}
	
	return parts[1]
}

// getTokenFromCookie extracts JWT token from cookie
func (w *WebAuthMiddleware) getTokenFromCookie(r *http.Request) string {
	cookie, err := r.Cookie("garm_token")
	if err != nil {
		return ""
	}
	return cookie.Value
}

// validateToken validates and parses the JWT token
func (w *WebAuthMiddleware) validateToken(tokenString string) (*auth.JWTClaims, error) {
	claims := &auth.JWTClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid signing method")
		}
		return []byte(w.cfg.Secret), nil
	})
	
	if err != nil {
		return nil, err
	}
	
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	
	return claims, nil
}

// claimsToContext adds user information from JWT claims to the request context
func (w *WebAuthMiddleware) claimsToContext(ctx context.Context, claims *auth.JWTClaims) (context.Context, error) {
	if claims == nil {
		return ctx, runnerErrors.ErrUnauthorized
	}

	if claims.UserID == "" {
		return nil, runnerErrors.ErrUnauthorized
	}

	userInfo, err := w.store.GetUserByID(ctx, claims.UserID)
	if err != nil {
		return ctx, runnerErrors.ErrUnauthorized
	}

	var expiresAt *time.Time
	if claims.ExpiresAt != nil {
		expires := claims.ExpiresAt.Time.UTC()
		expiresAt = &expires
	}

	if userInfo.Generation != claims.Generation {
		// Password was reset since token was issued. Invalidate.
		return ctx, runnerErrors.ErrUnauthorized
	}

	ctx = auth.PopulateContext(ctx, userInfo, expiresAt)
	return ctx, nil
}

// redirectToLogin redirects the user to the login page
func (w *WebAuthMiddleware) redirectToLogin(rw http.ResponseWriter, r *http.Request) {
	// Store the original URL in a cookie so we can redirect back after login
	returnURL := r.URL.Path
	if r.URL.RawQuery != "" {
		returnURL += "?" + r.URL.RawQuery
	}
	
	// Set return URL cookie
	http.SetCookie(rw, &http.Cookie{
		Name:     "garm_return_url",
		Value:    returnURL,
		Path:     "/",
		HttpOnly: true,
		Secure:   r.TLS != nil,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   300, // 5 minutes
	})
	
	// Redirect to login page
	http.Redirect(rw, r, "/web/login", http.StatusSeeOther)
}

// WebInitMiddleware is a web-specific initialization middleware that redirects
// to init page instead of returning JSON errors
type WebInitMiddleware struct {
	store dbCommon.Store
}

// NewWebInitMiddleware creates a new web initialization middleware
func NewWebInitMiddleware(store dbCommon.Store) *WebInitMiddleware {
	return &WebInitMiddleware{
		store: store,
	}
}

// Middleware implements the auth.Middleware interface for web routes
func (w *WebInitMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ctrlInfo, err := w.store.ControllerInfo()
		if err != nil || ctrlInfo.ControllerID.String() == "" {
			// Redirect to initialization page
			http.Redirect(rw, r, "/web/init", http.StatusSeeOther)
			return
		}

		next.ServeHTTP(rw, r.WithContext(ctx))
	})
}