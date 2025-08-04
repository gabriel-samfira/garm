package handlers

import (
	"encoding/json"
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
	
	// Check if filtering by forge type is requested
	forgeType := r.URL.Query().Get("forge_type")
	
	var credentials []params.ForgeCredentials
	
	switch params.EndpointType(forgeType) {
	case params.GithubEndpointType:
		// Only get GitHub credentials
		githubCredentials, err := h.runner.runner.ListCredentials(ctx)
		if err != nil {
			slog.ErrorContext(ctx, "Failed to list GitHub credentials", "error", err)
			http.Error(w, "Failed to list GitHub credentials", http.StatusInternalServerError)
			return
		}
		credentials = githubCredentials
	case params.GiteaEndpointType:
		// Only get Gitea credentials
		giteaCredentials, err := h.runner.runner.ListGiteaCredentials(ctx)
		if err != nil {
			slog.ErrorContext(ctx, "Failed to list Gitea credentials", "error", err)
			http.Error(w, "Failed to list Gitea credentials", http.StatusInternalServerError)
			return
		}
		credentials = giteaCredentials
	default:
		// Get both GitHub and Gitea credentials (default behavior)
		githubCredentials, err := h.runner.runner.ListCredentials(ctx)
		if err != nil {
			slog.ErrorContext(ctx, "Failed to list GitHub credentials", "error", err)
			http.Error(w, "Failed to list GitHub credentials", http.StatusInternalServerError)
			return
		}
		
		giteaCredentials, err := h.runner.runner.ListGiteaCredentials(ctx)
		if err != nil {
			slog.ErrorContext(ctx, "Failed to list Gitea credentials", "error", err)
			http.Error(w, "Failed to list Gitea credentials", http.StatusInternalServerError)
			return
		}
		
		// Combine all credentials
		credentials = append(githubCredentials, giteaCredentials...)
	}

	if r.Header.Get("Accept") == "application/json" {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(credentials)
		return
	}

	// Render as HTML table rows for HTMX
	w.Header().Set("Content-Type", "text/html")
	
	if len(credentials) == 0 {
		emptyData := EmptyTableRowData{
			ColSpan: 5,
			Icon: template.HTML(`<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 7a2 2 0 012 2m4 0a6 6 0 01-7.743 5.743L11 17H9v2H7v2H4a1 1 0 01-1-1v-2.586a1 1 0 01.293-.707l5.964-5.964A6 6 0 1121 9z"></path>`),
			Message:    "No credentials configured",
			SubMessage: "Add a credential to get started",
		}
		if err := h.templates.ExecuteTemplate(w, "empty-table-row.html", emptyData); err != nil {
			slog.Error("Failed to execute empty table row template", "error", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	for _, credential := range credentials {
		authTypeBadge := "bg-gray-100 dark:bg-gray-700 text-gray-800 dark:text-gray-200"
		authTypeText := string(credential.AuthType)
		
		switch credential.AuthType {
		case params.ForgeAuthTypePAT:
			authTypeBadge = "bg-green-100 dark:bg-green-900 text-green-800 dark:text-green-200"
			authTypeText = "Personal Access Token"
		case params.ForgeAuthTypeApp:
			authTypeBadge = "bg-blue-100 dark:bg-blue-900 text-blue-800 dark:text-blue-200"
			authTypeText = "GitHub App"
		}
		
		// Determine forge icon based on endpoint type
		var forgeIcon string
		switch credential.ForgeType {
		case params.GithubEndpointType:
			forgeIcon = `<svg class="w-4 h-4 inline mr-2 dark:hidden" width="98" height="96" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 98 96"><path fill-rule="evenodd" clip-rule="evenodd" d="M48.854 0C21.839 0 0 22 0 49.217c0 21.756 13.993 40.172 33.405 46.69 2.427.49 3.316-1.059 3.316-2.362 0-1.141-.08-5.052-.08-9.127-13.59 2.934-16.42-5.867-16.42-5.867-2.184-5.704-5.42-7.17-5.42-7.17-4.448-3.015.324-3.015.324-3.015 4.934.326 7.523 5.052 7.523 5.052 4.367 7.496 11.404 5.378 14.235 4.074.404-3.178 1.699-5.378 3.074-6.6-10.839-1.141-22.243-5.378-22.243-24.283 0-5.378 1.94-9.778 5.014-13.2-.485-1.222-2.184-6.275.486-13.038 0 0 4.125-1.304 13.426 5.052a46.97 46.97 0 0 1 12.214-1.63c4.125 0 8.33.571 12.213 1.63 9.302-6.356 13.427-5.052 13.427-5.052 2.67 6.763.97 11.816.485 13.038 3.155 3.422 5.015 7.822 5.015 13.2 0 18.905-11.404 23.06-22.324 24.283 1.78 1.548 3.316 4.481 3.316 9.126 0 6.6-.08 11.897-.08 13.526 0 1.304.89 2.853 3.316 2.364 19.412-6.52 33.405-24.935 33.405-46.691C97.707 22 75.788 0 48.854 0z" fill="#24292f"/></svg><svg class="w-4 h-4 inline mr-2 hidden dark:block" width="98" height="96" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 98 96"><path fill-rule="evenodd" clip-rule="evenodd" d="M48.854 0C21.839 0 0 22 0 49.217c0 21.756 13.993 40.172 33.405 46.69 2.427.49 3.316-1.059 3.316-2.362 0-1.141-.08-5.052-.08-9.127-13.59 2.934-16.42-5.867-16.42-5.867-2.184-5.704-5.42-7.17-5.42-7.17-4.448-3.015.324-3.015.324-3.015 4.934.326 7.523 5.052 7.523 5.052 4.367 7.496 11.404 5.378 14.235 4.074.404-3.178 1.699-5.378 3.074-6.6-10.839-1.141-22.243-5.378-22.243-24.283 0-5.378 1.94-9.778 5.014-13.2-.485-1.222-2.184-6.275.486-13.038 0 0 4.125-1.304 13.426 5.052a46.97 46.97 0 0 1 12.214-1.63c4.125 0 8.33.571 12.213 1.63 9.302-6.356 13.427-5.052 13.427-5.052 2.67 6.763.97 11.816.485 13.038 3.155 3.422 5.015 7.822 5.015 13.2 0 18.905-11.404 23.06-22.324 24.283 1.78 1.548 3.316 4.481 3.316 9.126 0 6.6-.08 11.897-.08 13.526 0 1.304.89 2.853 3.316 2.364 19.412-6.52 33.405-24.935 33.405-46.691C97.707 22 75.788 0 48.854 0z" fill="#fff"/></svg>`
		case params.GiteaEndpointType:
			forgeIcon = `<svg class="w-4 h-4 inline mr-2" xmlns="http://www.w3.org/2000/svg" xml:space="preserve" viewBox="0 0 640 640"><path d="m395.9 484.2-126.9-61c-12.5-6-17.9-21.2-11.8-33.8l61-126.9c6-12.5 21.2-17.9 33.8-11.8 17.2 8.3 27.1 13 27.1 13l-.1-109.2 16.7-.1.1 117.1s57.4 24.2 83.1 40.1c3.7 2.3 10.2 6.8 12.9 14.4 2.1 6.1 2 13.1-1 19.3l-61 126.9c-6.2 12.7-21.4 18.1-33.9 12" style="fill:#fff"/><path d="M622.7 149.8c-4.1-4.1-9.6-4-9.6-4s-117.2 6.6-177.9 8c-13.3.3-26.5.6-39.6.7v117.2c-5.5-2.6-11.1-5.3-16.6-7.9 0-36.4-.1-109.2-.1-109.2-29 .4-89.2-2.2-89.2-2.2s-141.4-7.1-156.8-8.5c-9.8-.6-22.5-2.1-39 1.5-8.7 1.8-33.5 7.4-53.8 26.9C-4.9 212.4 6.6 276.2 8 285.8c1.7 11.7 6.9 44.2 31.7 72.5 45.8 56.1 144.4 54.8 144.4 54.8s12.1 28.9 30.6 55.5c25 33.1 50.7 58.9 75.7 62 63 0 188.9-.1 188.9-.1s12 .1 28.3-10.3c14-8.5 26.5-23.4 26.5-23.4S547 483 565 451.5c5.5-9.7 10.1-19.1 14.1-28 0 0 55.2-117.1 55.2-231.1-1.1-34.5-9.6-40.6-11.6-42.6M125.6 353.9c-25.9-8.5-36.9-18.7-36.9-18.7S69.6 321.8 60 295.4c-16.5-44.2-1.4-71.2-1.4-71.2s8.4-22.5 38.5-30c13.8-3.7 31-3.1 31-3.1s7.1 59.4 15.7 94.2c7.2 29.2 24.8 77.7 24.8 77.7s-26.1-3.1-43-9.1m300.3 107.6s-6.1 14.5-19.6 15.4c-5.8.4-10.3-1.2-10.3-1.2s-.3-.1-5.3-2.1l-112.9-55s-10.9-5.7-12.8-15.6c-2.2-8.1 2.7-18.1 2.7-18.1L322 273s4.8-9.7 12.2-13c.6-.3 2.3-1 4.5-1.5 8.1-2.1 18 2.8 18 2.8L467.4 315s12.6 5.7 15.3 16.2c1.9 7.4-.5 14-1.8 17.2-6.3 15.4-55 113.1-55 113.1" style="fill:#609926"/><path d="M326.8 380.1c-8.2.1-15.4 5.8-17.3 13.8s2 16.3 9.1 20c7.7 4 17.5 1.8 22.7-5.4 5.1-7.1 4.3-16.9-1.8-23.1l24-49.1c1.5.1 3.7.2 6.2-.5 4.1-.9 7.1-3.6 7.1-3.6 4.2 1.8 8.6 3.8 13.2 6.1 4.8 2.4 9.3 4.9 13.4 7.3.9.5 1.8 1.1 2.8 1.9 1.6 1.3 3.4 3.1 4.7 5.5 1.9 5.5-1.9 14.9-1.9 14.9-2.3 7.6-18.4 40.6-18.4 40.6-8.1-.2-15.3 5-17.7 12.5-2.6 8.1 1.1 17.3 8.9 21.3s17.4 1.7 22.5-5.3c5-6.8 4.6-16.3-1.1-22.6 1.9-3.7 3.7-7.4 5.6-11.3 5-10.4 13.5-30.4 13.5-30.4.9-1.7 5.7-10.3 2.7-21.3-2.5-11.4-12.6-16.7-12.6-16.7-12.2-7.9-29.2-15.2-29.2-15.2s0-4.1-1.1-7.1c-1.1-3.1-2.8-5.1-3.9-6.3 4.7-9.7 9.4-19.3 14.1-29-4.1-2-8.1-4-12.2-6.1-4.8 9.8-9.7 19.7-14.5 29.5-6.7-.1-12.9 3.5-16.1 9.4-3.4 6.3-2.7 14.1 1.9 19.8z" style="fill:#609926"/></svg>`
		}

		rowData := CredentialRowData{
			ID:            credential.ID,
			Name:          credential.Name,
			Description:   credential.Description,
			ForgeIcon:     template.HTML(forgeIcon),
			EndpointName:  credential.Endpoint.Name,
			AuthTypeBadge: authTypeBadge,
			AuthTypeText:  authTypeText,
		}
		
		if err := h.templates.ExecuteTemplate(w, "credential-table-row.html", rowData); err != nil {
			slog.Error("Failed to execute credential table row template", "error", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// NewCredentialFormHandler serves the forge type selector or credential form
func (h *CredentialsHandler) NewCredentialFormHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	
	// Check if forge_type is specified in query params
	forgeType := r.URL.Query().Get("forge_type")
	
	if forgeType == "" {
		// Show forge type selector
		if err := h.templates.ExecuteTemplate(w, "credential-forge-selector.html", nil); err != nil {
			slog.ErrorContext(ctx, "Failed to execute forge selector template", "error", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}
	
	// Show credential form for the selected forge type
	var endpoints []params.ForgeEndpoint
	var isGitea bool
	
	if forgeType == string(params.GiteaEndpointType) {
		giteaEndpoints, err := h.runner.runner.ListGiteaEndpoints(ctx)
		if err != nil {
			slog.ErrorContext(ctx, "Failed to list Gitea endpoints", "error", err)
			endpoints = []params.ForgeEndpoint{}
		} else {
			endpoints = giteaEndpoints
		}
		isGitea = true
	} else {
		// Default to GitHub
		githubEndpoints, err := h.runner.runner.ListGithubEndpoints(ctx)
		if err != nil {
			slog.ErrorContext(ctx, "Failed to list GitHub endpoints", "error", err)
			endpoints = []params.ForgeEndpoint{}
		} else {
			endpoints = githubEndpoints
		}
		isGitea = false
	}

	data := struct {
		Credential *params.ForgeCredentials
		Endpoints  []params.ForgeEndpoint
		IsGitea    bool
		ForgeType  string
	}{
		Credential: nil, // No credential for new form
		Endpoints:  endpoints,
		IsGitea:    isGitea,
		ForgeType:  forgeType,
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

	// Parse the JSON to get the endpoint name first to determine forge type
	var rawData map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&rawData); err != nil {
		slog.ErrorContext(ctx, "Failed to decode create credential request", "error", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	endpointName, ok := rawData["endpoint"].(string)
	if !ok || endpointName == "" {
		slog.ErrorContext(ctx, "Missing or invalid endpoint in request")
		http.Error(w, "Endpoint is required", http.StatusBadRequest)
		return
	}

	// Determine if this endpoint is GitHub or Gitea by trying to find it in both lists
	githubEndpoints, _ := h.runner.runner.ListGithubEndpoints(ctx)
	giteaEndpoints, _ := h.runner.runner.ListGiteaEndpoints(ctx)
	
	var isGitea bool
	for _, ep := range githubEndpoints {
		if ep.Name == endpointName {
			isGitea = false
			break
		}
	}
	if !isGitea {
		// Check if it's in Gitea endpoints
		for _, ep := range giteaEndpoints {
			if ep.Name == endpointName {
				isGitea = true
				break
			}
		}
	}

	// Convert rawData back to JSON for decoding into the appropriate struct
	rawJSON, err := json.Marshal(rawData)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to marshal request data", "error", err)
		http.Error(w, "Invalid request data", http.StatusInternalServerError)
		return
	}

	if isGitea {
		var createParams params.CreateGiteaCredentialsParams
		if err := json.Unmarshal(rawJSON, &createParams); err != nil {
			slog.ErrorContext(ctx, "Failed to decode create Gitea credential request", "error", err)
			http.Error(w, "Invalid JSON for Gitea credentials", http.StatusBadRequest)
			return
		}

		credential, err := h.runner.runner.CreateGiteaCredentials(ctx, createParams)
		if err != nil {
			slog.ErrorContext(ctx, "Failed to create Gitea credential", "error", err)
			w.Header().Set("HX-Trigger", "credentialCreateError")
			http.Error(w, "Failed to create credential", http.StatusInternalServerError)
			return
		}

		w.Header().Set("HX-Trigger", "credentialCreated")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(credential)
	} else {
		var createParams params.CreateGithubCredentialsParams
		if err := json.Unmarshal(rawJSON, &createParams); err != nil {
			slog.ErrorContext(ctx, "Failed to decode create GitHub credential request", "error", err)
			http.Error(w, "Invalid JSON for GitHub credentials", http.StatusBadRequest)
			return
		}

		credential, err := h.runner.runner.CreateGithubCredentials(ctx, createParams)
		if err != nil {
			slog.ErrorContext(ctx, "Failed to create GitHub credential", "error", err)
			w.Header().Set("HX-Trigger", "credentialCreateError")
			http.Error(w, "Failed to create credential", http.StatusInternalServerError)
			return
		}

		w.Header().Set("HX-Trigger", "credentialCreated")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(credential)
	}
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

	// Try to get the credential from both GitHub and Gitea
	var credential params.ForgeCredentials
	var isGitea bool
	
	// First try GitHub
	githubCred, err := h.runner.runner.GetGithubCredentials(ctx, uint(id))
	if err == nil {
		credential = githubCred
		isGitea = false
	} else {
		// If GitHub fails, try Gitea
		giteaCred, err := h.runner.runner.GetGiteaCredentials(ctx, uint(id))
		if err != nil {
			slog.ErrorContext(ctx, "Failed to get credential from both GitHub and Gitea", "error", err)
			http.Error(w, "Credential not found", http.StatusNotFound)
			return
		}
		credential = giteaCred
		isGitea = true
	}

	// Get available endpoints based on the credential type
	var endpoints []params.ForgeEndpoint
	if isGitea {
		giteaEndpoints, err := h.runner.runner.ListGiteaEndpoints(ctx)
		if err != nil {
			slog.ErrorContext(ctx, "Failed to list Gitea endpoints", "error", err)
			endpoints = []params.ForgeEndpoint{}
		} else {
			endpoints = giteaEndpoints
		}
	} else {
		githubEndpoints, err := h.runner.runner.ListGithubEndpoints(ctx)
		if err != nil {
			slog.ErrorContext(ctx, "Failed to list GitHub endpoints", "error", err)
			endpoints = []params.ForgeEndpoint{}
		} else {
			endpoints = githubEndpoints
		}
	}

	data := struct {
		Credential params.ForgeCredentials
		Endpoints  []params.ForgeEndpoint
		IsGitea    bool
		ForgeType  string
	}{
		Credential: credential,
		Endpoints:  endpoints,
		IsGitea:    isGitea,
		ForgeType:  string(credential.ForgeType),
	}

	if err := h.templates.ExecuteTemplate(w, "credential-form.html", data); err != nil {
		slog.ErrorContext(ctx, "Failed to execute credential form template", "error", err)
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

	// Determine if this is a GitHub or Gitea credential
	var isGitea bool
	_, err = h.runner.runner.GetGithubCredentials(ctx, uint(id))
	if err != nil {
		// Try Gitea
		_, err = h.runner.runner.GetGiteaCredentials(ctx, uint(id))
		if err != nil {
			slog.ErrorContext(ctx, "Credential not found in both GitHub and Gitea", "error", err)
			http.Error(w, "Credential not found", http.StatusNotFound)
			return
		}
		isGitea = true
	}

	if isGitea {
		var updateParams params.UpdateGiteaCredentialsParams
		if err := json.NewDecoder(r.Body).Decode(&updateParams); err != nil {
			slog.ErrorContext(ctx, "Failed to decode update Gitea credential request", "error", err)
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		credential, err := h.runner.runner.UpdateGiteaCredentials(ctx, uint(id), updateParams)
		if err != nil {
			slog.ErrorContext(ctx, "Failed to update Gitea credential", "error", err)
			w.Header().Set("HX-Trigger", "credentialUpdateError")
			http.Error(w, "Failed to update credential", http.StatusInternalServerError)
			return
		}

		w.Header().Set("HX-Trigger", "credentialUpdated")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(credential)
	} else {
		var updateParams params.UpdateGithubCredentialsParams
		if err := json.NewDecoder(r.Body).Decode(&updateParams); err != nil {
			slog.ErrorContext(ctx, "Failed to decode update GitHub credential request", "error", err)
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		credential, err := h.runner.runner.UpdateGithubCredentials(ctx, uint(id), updateParams)
		if err != nil {
			slog.ErrorContext(ctx, "Failed to update GitHub credential", "error", err)
			w.Header().Set("HX-Trigger", "credentialUpdateError")
			http.Error(w, "Failed to update credential", http.StatusInternalServerError)
			return
		}

		w.Header().Set("HX-Trigger", "credentialUpdated")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(credential)
	}
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