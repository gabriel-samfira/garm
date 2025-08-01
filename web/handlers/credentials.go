package handlers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/cloudbase/garm/params"
)

type CredentialsHandler struct {
	runner    *WebHandler
	templates *template.Template
}

func NewCredentialsHandler(webHandler *WebHandler) *CredentialsHandler {
	return &CredentialsHandler{
		runner:    webHandler,
		templates: webHandler.templates,
	}
}

// CredentialsPageHandler serves the credentials management page
func (h *CredentialsHandler) CredentialsPageHandler(w http.ResponseWriter, r *http.Request) {
	data := PageData{
		Title:      "Credentials",
		PageTitle:  "Credential Management",
		ActivePage: "credentials",
		CreateButton: &CreateButton{
			URL:  "/web/credentials/new",
			Text: "Add Credential",
		},
	}

	if err := h.templates.ExecuteTemplate(w, "base.html", data); err != nil {
		slog.Error("Failed to execute credentials template", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// ListCredentialsHandler returns credentials as JSON for HTMX
func (h *CredentialsHandler) ListCredentialsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	
	credentials, err := h.runner.runner.ListCredentials(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to list credentials", "error", err)
		http.Error(w, "Failed to list credentials", http.StatusInternalServerError)
		return
	}

	if r.Header.Get("Accept") == "application/json" {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(credentials)
		return
	}

	// Render as HTML table rows for HTMX
	w.Header().Set("Content-Type", "text/html")
	
	if len(credentials) == 0 {
		w.Write([]byte(`<tr>
			<td colspan="5" class="px-6 py-4 text-center text-gray-500 dark:text-gray-400">
				<div class="flex flex-col items-center py-8">
					<svg class="w-12 h-12 text-gray-400 dark:text-gray-500 mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 7a2 2 0 012 2m4 0a6 6 0 01-7.743 5.743L11 17H9v2H7v2H4a1 1 0 01-1-1v-2.586a1 1 0 01.293-.707l5.964-5.964A6 6 0 1121 9z"></path>
					</svg>
					<p class="text-sm">No credentials configured</p>
					<p class="text-xs text-gray-400 dark:text-gray-500 mt-1">Add a credential to get started</p>
				</div>
			</td>
		</tr>`))
		return
	}

	for _, credential := range credentials {
		authTypeBadge := "bg-gray-100 dark:bg-gray-700 text-gray-800 dark:text-gray-200"
		authTypeText := string(credential.AuthType)
		
		if credential.AuthType == "pat" {
			authTypeBadge = "bg-green-100 dark:bg-green-900 text-green-800 dark:text-green-200"
			authTypeText = "Personal Access Token"
		} else if credential.AuthType == "app" {
			authTypeBadge = "bg-blue-100 dark:bg-blue-900 text-blue-800 dark:text-blue-200"
			authTypeText = "GitHub App"
		}

		html := fmt.Sprintf(`
		<tr class="hover:bg-gray-50 dark:hover:bg-gray-700">
			<td class="px-6 py-4 whitespace-nowrap">
				<div class="text-sm font-medium text-gray-900 dark:text-white">%s</div>
			</td>
			<td class="px-6 py-4 whitespace-nowrap">
				<div class="text-sm text-gray-900 dark:text-white">%s</div>
			</td>
			<td class="px-6 py-4 whitespace-nowrap">
				<div class="text-sm text-gray-900 dark:text-white">%s</div>
			</td>
			<td class="px-6 py-4 whitespace-nowrap">
				<span class="inline-flex px-2 py-1 text-xs font-semibold rounded-full %s">
					%s
				</span>
			</td>
			<td class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
				<div class="flex space-x-2">
					<button hx-get="/web/credentials/%d/edit" 
							hx-target="#modal-container"
							class="text-blue-600 dark:text-blue-400 hover:text-blue-900 dark:hover:text-blue-300 transition-colors">
						Edit
					</button>
					<button hx-delete="/web/api/credentials/%d" 
							hx-confirm="Are you sure you want to delete this credential?"
							hx-target="closest tr"
							hx-swap="outerHTML"
							class="text-red-600 dark:text-red-400 hover:text-red-900 dark:hover:text-red-300 transition-colors">
						Delete
					</button>
				</div>
			</td>
		</tr>`, credential.Name, credential.Description, credential.Endpoint.Name, authTypeBadge, authTypeText, credential.ID, credential.ID)
		w.Write([]byte(html))
	}
}

// NewCredentialFormHandler serves the new credential form
func (h *CredentialsHandler) NewCredentialFormHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	
	// Get available endpoints for the dropdown
	endpoints, err := h.runner.runner.ListGithubEndpoints(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to list endpoints", "error", err)
		endpoints = []params.ForgeEndpoint{} // Continue with empty list
	}

	data := struct {
		Endpoints []params.ForgeEndpoint
	}{
		Endpoints: endpoints,
	}

	if err := h.templates.ExecuteTemplate(w, "credential-form.html", data); err != nil {
		slog.ErrorContext(ctx, "Failed to execute credential form template", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// CreateCredentialHandler handles credential creation
func (h *CredentialsHandler) CreateCredentialHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var createParams params.CreateGithubCredentialsParams
	if err := json.NewDecoder(r.Body).Decode(&createParams); err != nil {
		slog.ErrorContext(ctx, "Failed to decode create credential request", "error", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	credential, err := h.runner.runner.CreateGithubCredentials(ctx, createParams)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to create credential", "error", err)
		http.Error(w, "Failed to create credential", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(credential)
}

// EditCredentialFormHandler serves the edit credential form
func (h *CredentialsHandler) EditCredentialFormHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	idParam := vars["id"]
	
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		slog.ErrorContext(ctx, "Invalid credential ID", "error", err)
		http.Error(w, "Invalid credential ID", http.StatusBadRequest)
		return
	}

	credential, err := h.runner.runner.GetGithubCredentials(ctx, uint(id))
	if err != nil {
		slog.ErrorContext(ctx, "Failed to get credential", "error", err)
		http.Error(w, "Failed to get credential", http.StatusInternalServerError)
		return
	}

	// Get available endpoints for the dropdown
	endpoints, err := h.runner.runner.ListGithubEndpoints(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to list endpoints", "error", err)
		endpoints = []params.ForgeEndpoint{} // Continue with empty list
	}

	data := struct {
		Credential params.ForgeCredentials
		Endpoints  []params.ForgeEndpoint
	}{
		Credential: credential,
		Endpoints:  endpoints,
	}

	if err := h.templates.ExecuteTemplate(w, "credential-edit-form.html", data); err != nil {
		slog.ErrorContext(ctx, "Failed to execute credential edit form template", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// UpdateCredentialHandler handles credential updates
func (h *CredentialsHandler) UpdateCredentialHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	idParam := vars["id"]
	
	if r.Method != "PUT" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		slog.ErrorContext(ctx, "Invalid credential ID", "error", err)
		http.Error(w, "Invalid credential ID", http.StatusBadRequest)
		return
	}

	var updateParams params.UpdateGithubCredentialsParams
	if err := json.NewDecoder(r.Body).Decode(&updateParams); err != nil {
		slog.ErrorContext(ctx, "Failed to decode update credential request", "error", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	credential, err := h.runner.runner.UpdateGithubCredentials(ctx, uint(id), updateParams)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to update credential", "error", err)
		http.Error(w, "Failed to update credential", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(credential)
}

// DeleteCredentialHandler handles credential deletion
func (h *CredentialsHandler) DeleteCredentialHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	idParam := vars["id"]
	
	if r.Method != "DELETE" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		slog.ErrorContext(ctx, "Invalid credential ID", "error", err)
		http.Error(w, "Invalid credential ID", http.StatusBadRequest)
		return
	}

	if err := h.runner.runner.DeleteGithubCredentials(ctx, uint(id)); err != nil {
		slog.ErrorContext(ctx, "Failed to delete credential", "error", err)
		http.Error(w, "Failed to delete credential", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}