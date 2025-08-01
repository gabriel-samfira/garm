package handlers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/cloudbase/garm/params"
)

type EndpointsHandler struct {
	runner    *WebHandler
	templates *template.Template
}

func NewEndpointsHandler(webHandler *WebHandler) *EndpointsHandler {
	return &EndpointsHandler{
		runner:    webHandler,
		templates: webHandler.templates,
	}
}

// EndpointsPageHandler serves the endpoints management page
func (h *EndpointsHandler) EndpointsPageHandler(w http.ResponseWriter, r *http.Request) {
	data := PageData{
		Title:      "Endpoints",
		PageTitle:  "Endpoint Management",
		ActivePage: "endpoints",
		CreateButton: &CreateButton{
			URL:  "/web/endpoints/new",
			Text: "Add Endpoint",
		},
	}

	if err := h.templates.ExecuteTemplate(w, "base.html", data); err != nil {
		slog.Error("Failed to execute endpoints template", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// ListEndpointsHandler returns endpoints as JSON for HTMX
func (h *EndpointsHandler) ListEndpointsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	
	// Get both GitHub and Gitea endpoints
	githubEndpoints, err := h.runner.runner.ListGithubEndpoints(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to list GitHub endpoints", "error", err)
		http.Error(w, "Failed to list endpoints", http.StatusInternalServerError)
		return
	}
	
	giteaEndpoints, err := h.runner.runner.ListGiteaEndpoints(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to list Gitea endpoints", "error", err)
		http.Error(w, "Failed to list endpoints", http.StatusInternalServerError)
		return
	}
	
	// Combine all endpoints
	allEndpoints := make([]params.ForgeEndpoint, 0, len(githubEndpoints)+len(giteaEndpoints))
	allEndpoints = append(allEndpoints, githubEndpoints...)
	allEndpoints = append(allEndpoints, giteaEndpoints...)
	endpoints := allEndpoints

	if r.Header.Get("Accept") == "application/json" {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(endpoints)
		return
	}

	// Render as HTML table rows for HTMX
	w.Header().Set("Content-Type", "text/html")
	
	if len(endpoints) == 0 {
		w.Write([]byte(`<tr>
			<td colspan="6" class="px-6 py-4 text-center text-gray-500 dark:text-gray-400">
				<div class="flex flex-col items-center py-8">
					<svg class="w-12 h-12 text-gray-400 dark:text-gray-500 mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 12a9 9 0 01-9 9m9-9a9 9 0 00-9-9m9 9H3m9 9v-9m0-9v9"></path>
					</svg>
					<p class="text-sm">No endpoints configured</p>
					<p class="text-xs text-gray-400 dark:text-gray-500 mt-1">Add an endpoint to get started</p>
				</div>
			</td>
		</tr>`))
		return
	}

	for _, endpoint := range endpoints {
		// Determine forge type and icon
		forgeIcon := ""
		forgeType := ""
		if endpoint.UploadBaseURL != "" {
			// GitHub endpoint (has UploadBaseURL)
			forgeType = "GitHub"
			forgeIcon = `<svg class="w-5 h-5 dark:hidden" width="98" height="96" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 98 96">
				<path fill-rule="evenodd" clip-rule="evenodd" d="M48.854 0C21.839 0 0 22 0 49.217c0 21.756 13.993 40.172 33.405 46.69 2.427.49 3.316-1.059 3.316-2.362 0-1.141-.08-5.052-.08-9.127-13.59 2.934-16.42-5.867-16.42-5.867-2.184-5.704-5.42-7.17-5.42-7.17-4.448-3.015.324-3.015.324-3.015 4.934.326 7.523 5.052 7.523 5.052 4.367 7.496 11.404 5.378 14.235 4.074.404-3.178 1.699-5.378 3.074-6.6-10.839-1.141-22.243-5.378-22.243-24.283 0-5.378 1.94-9.778 5.014-13.2-.485-1.222-2.184-6.275.486-13.038 0 0 4.125-1.304 13.426 5.052a46.97 46.97 0 0 1 12.214-1.63c4.125 0 8.33.571 12.213 1.63 9.302-6.356 13.427-5.052 13.427-5.052 2.67 6.763.97 11.816.485 13.038 3.155 3.422 5.015 7.822 5.015 13.2 0 18.905-11.404 23.06-22.324 24.283 1.78 1.548 3.316 4.481 3.316 9.126 0 6.6-.08 11.897-.08 13.526 0 1.304.89 2.853 3.316 2.364 19.412-6.52 33.405-24.935 33.405-46.691C97.707 22 75.788 0 48.854 0z" fill="#24292f"/>
			</svg>
			<svg class="w-5 h-5 hidden dark:block" width="98" height="96" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 98 96">
				<path fill-rule="evenodd" clip-rule="evenodd" d="M48.854 0C21.839 0 0 22 0 49.217c0 21.756 13.993 40.172 33.405 46.69 2.427.49 3.316-1.059 3.316-2.362 0-1.141-.08-5.052-.08-9.127-13.59 2.934-16.42-5.867-16.42-5.867-2.184-5.704-5.42-7.17-5.42-7.17-4.448-3.015.324-3.015.324-3.015 4.934.326 7.523 5.052 7.523 5.052 4.367 7.496 11.404 5.378 14.235 4.074.404-3.178 1.699-5.378 3.074-6.6-10.839-1.141-22.243-5.378-22.243-24.283 0-5.378 1.94-9.778 5.014-13.2-.485-1.222-2.184-6.275.486-13.038 0 0 4.125-1.304 13.426 5.052a46.97 46.97 0 0 1 12.214-1.63c4.125 0 8.33.571 12.213 1.63 9.302-6.356 13.427-5.052 13.427-5.052 2.67 6.763.97 11.816.485 13.038 3.155 3.422 5.015 7.822 5.015 13.2 0 18.905-11.404 23.06-22.324 24.283 1.78 1.548 3.316 4.481 3.316 9.126 0 6.6-.08 11.897-.08 13.526 0 1.304.89 2.853 3.316 2.364 19.412-6.52 33.405-24.935 33.405-46.691C97.707 22 75.788 0 48.854 0z" fill="#fff"/>
			</svg>`
		} else {
			// Gitea endpoint
			forgeType = "Gitea"
			forgeIcon = `<svg class="w-5 h-5" xmlns="http://www.w3.org/2000/svg" xml:space="preserve" viewBox="0 0 640 640">
				<path d="m395.9 484.2-126.9-61c-12.5-6-17.9-21.2-11.8-33.8l61-126.9c6-12.5 21.2-17.9 33.8-11.8 17.2 8.3 27.1 13 27.1 13l-.1-109.2 16.7-.1.1 117.1s57.4 24.2 83.1 40.1c3.7 2.3 10.2 6.8 12.9 14.4 2.1 6.1 2 13.1-1 19.3l-61 126.9c-6.2 12.7-21.4 18.1-33.9 12" style="fill:#fff"/>
				<path d="M622.7 149.8c-4.1-4.1-9.6-4-9.6-4s-117.2 6.6-177.9 8c-13.3.3-26.5.6-39.6.7v117.2c-5.5-2.6-11.1-5.3-16.6-7.9 0-36.4-.1-109.2-.1-109.2-29 .4-89.2-2.2-89.2-2.2s-141.4-7.1-156.8-8.5c-9.8-.6-22.5-2.1-39 1.5-8.7 1.8-33.5 7.4-53.8 26.9C-4.9 212.4 6.6 276.2 8 285.8c1.7 11.7 6.9 44.2 31.7 72.5 45.8 56.1 144.4 54.8 144.4 54.8s12.1 28.9 30.6 55.5c25 33.1 50.7 58.9 75.7 62 63 0 188.9-.1 188.9-.1s12 .1 28.3-10.3c14-8.5 26.5-23.4 26.5-23.4S547 483 565 451.5c5.5-9.7 10.1-19.1 14.1-28 0 0 55.2-117.1 55.2-231.1-1.1-34.5-9.6-40.6-11.6-42.6M125.6 353.9c-25.9-8.5-36.9-18.7-36.9-18.7S69.6 321.8 60 295.4c-16.5-44.2-1.4-71.2-1.4-71.2s8.4-22.5 38.5-30c13.8-3.7 31-3.1 31-3.1s7.1 59.4 15.7 94.2c7.2 29.2 24.8 77.7 24.8 77.7s-26.1-3.1-43-9.1m300.3 107.6s-6.1 14.5-19.6 15.4c-5.8.4-10.3-1.2-10.3-1.2s-.3-.1-5.3-2.1l-112.9-55s-10.9-5.7-12.8-15.6c-2.2-8.1 2.7-18.1 2.7-18.1L322 273s4.8-9.7 12.2-13c.6-.3 2.3-1 4.5-1.5 8.1-2.1 18 2.8 18 2.8L467.4 315s12.6 5.7 15.3 16.2c1.9 7.4-.5 14-1.8 17.2-6.3 15.4-55 113.1-55 113.1" style="fill:#609926"/>
				<path d="M326.8 380.1c-8.2.1-15.4 5.8-17.3 13.8s2 16.3 9.1 20c7.7 4 17.5 1.8 22.7-5.4 5.1-7.1 4.3-16.9-1.8-23.1l24-49.1c1.5.1 3.7.2 6.2-.5 4.1-.9 7.1-3.6 7.1-3.6 4.2 1.8 8.6 3.8 13.2 6.1 4.8 2.4 9.3 4.9 13.4 7.3.9.5 1.8 1.1 2.8 1.9 1.6 1.3 3.4 3.1 4.7 5.5 1.9 5.5-1.9 14.9-1.9 14.9-2.3 7.6-18.4 40.6-18.4 40.6-8.1-.2-15.3 5-17.7 12.5-2.6 8.1 1.1 17.3 8.9 21.3s17.4 1.7 22.5-5.3c5-6.8 4.6-16.3-1.1-22.6 1.9-3.7 3.7-7.4 5.6-11.3 5-10.4 13.5-30.4 13.5-30.4.9-1.7 5.7-10.3 2.7-21.3-2.5-11.4-12.6-16.7-12.6-16.7-12.2-7.9-29.2-15.2-29.2-15.2s0-4.1-1.1-7.1c-1.1-3.1-2.8-5.1-3.9-6.3 4.7-9.7 9.4-19.3 14.1-29-4.1-2-8.1-4-12.2-6.1-4.8 9.8-9.7 19.7-14.5 29.5-6.7-.1-12.9 3.5-16.1 9.4-3.4 6.3-2.7 14.1 1.9 19.8z" style="fill:#609926"/>
			</svg>`
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
				<div class="flex items-center space-x-2">
					%s
					<span class="text-sm text-gray-900 dark:text-white">%s</span>
				</div>
			</td>
			<td class="px-6 py-4 whitespace-nowrap">
				<div class="text-sm text-gray-900 dark:text-white">%s</div>
			</td>
			<td class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
				<div class="flex space-x-2">
					<button hx-get="/web/endpoints/%s/edit" 
							hx-target="#modal-container"
							class="text-blue-600 dark:text-blue-400 hover:text-blue-900 dark:hover:text-blue-300 transition-colors">
						Edit
					</button>
					<button hx-delete="/web/api/endpoints/%s" 
							hx-confirm="Are you sure you want to delete this endpoint?"
							hx-target="closest tr"
							hx-swap="outerHTML"
							class="text-red-600 dark:text-red-400 hover:text-red-900 dark:hover:text-red-300 transition-colors">
						Delete
					</button>
				</div>
			</td>
		</tr>`, endpoint.Name, endpoint.Description, endpoint.APIBaseURL, forgeIcon, forgeType, endpoint.BaseURL, endpoint.Name, endpoint.Name)
		w.Write([]byte(html))
	}
}

// NewEndpointFormHandler serves the new endpoint form
func (h *EndpointsHandler) NewEndpointFormHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	data := struct {
		// Add any data needed for the form
	}{}

	if err := h.templates.ExecuteTemplate(w, "endpoint-form.html", data); err != nil {
		slog.ErrorContext(ctx, "Failed to execute endpoint form template", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// CreateEndpointHandler handles endpoint creation
func (h *EndpointsHandler) CreateEndpointHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		slog.ErrorContext(ctx, "Failed to parse form", "error", err)
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	endpointType := r.FormValue("endpoint_type")
	name := r.FormValue("name")
	description := r.FormValue("description")
	baseURL := r.FormValue("base_url")
	apiBaseURL := r.FormValue("api_base_url")
	uploadBaseURL := r.FormValue("upload_base_url")
	caCertBase64 := r.FormValue("ca_cert_bundle")

	// Decode CA certificate bundle if provided
	var caCertBundle []byte
	if caCertBase64 != "" {
		decoded, err := base64.StdEncoding.DecodeString(caCertBase64)
		if err != nil {
			slog.ErrorContext(ctx, "Failed to decode CA certificate", "error", err)
			http.Error(w, "Invalid CA certificate encoding", http.StatusBadRequest)
			return
		}
		caCertBundle = decoded
	}

	var err error

	if endpointType == "github" {
		createParams := params.CreateGithubEndpointParams{
			Name:          name,
			Description:   description,
			APIBaseURL:    apiBaseURL,
			UploadBaseURL: uploadBaseURL,
			BaseURL:       baseURL,
			CACertBundle:  caCertBundle,
		}
		_, err = h.runner.runner.CreateGithubEndpoint(ctx, createParams)
	} else if endpointType == "gitea" {
		// For Gitea, use BaseURL as APIBaseURL if APIBaseURL is empty
		if apiBaseURL == "" {
			apiBaseURL = baseURL
		}
		createParams := params.CreateGiteaEndpointParams{
			Name:         name,
			Description:  description,
			APIBaseURL:   apiBaseURL,
			BaseURL:      baseURL,
			CACertBundle: caCertBundle,
		}
		_, err = h.runner.runner.CreateGiteaEndpoint(ctx, createParams)
	} else {
		http.Error(w, "Invalid endpoint type", http.StatusBadRequest)
		return
	}

	if err != nil {
		slog.ErrorContext(ctx, "Failed to create endpoint", "error", err, "type", endpointType)
		http.Error(w, "Failed to create endpoint: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Return success message and redirect
	w.Header().Set("Content-Type", "text/html")
	w.Header().Set("HX-Trigger", "endpointCreated")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`<div class="bg-green-100 dark:bg-green-900 border border-green-400 dark:border-green-600 text-green-700 dark:text-green-300 px-4 py-3 rounded">
		Endpoint created successfully!
	</div>`))
}

// EditEndpointFormHandler serves the edit endpoint form
func (h *EndpointsHandler) EditEndpointFormHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	name := vars["name"]
	
	endpoint, err := h.runner.runner.GetGithubEndpoint(ctx, name)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to get endpoint", "error", err)
		http.Error(w, "Failed to get endpoint", http.StatusInternalServerError)
		return
	}

	data := struct {
		Endpoint params.ForgeEndpoint
	}{
		Endpoint: endpoint,
	}

	if err := h.templates.ExecuteTemplate(w, "endpoint-edit-form.html", data); err != nil {
		slog.ErrorContext(ctx, "Failed to execute endpoint edit form template", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// UpdateEndpointHandler handles endpoint updates
func (h *EndpointsHandler) UpdateEndpointHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	name := vars["name"]
	
	if r.Method != "PUT" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var updateParams params.UpdateGithubEndpointParams
	if err := json.NewDecoder(r.Body).Decode(&updateParams); err != nil {
		slog.ErrorContext(ctx, "Failed to decode update endpoint request", "error", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	endpoint, err := h.runner.runner.UpdateGithubEndpoint(ctx, name, updateParams)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to update endpoint", "error", err)
		http.Error(w, "Failed to update endpoint", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(endpoint)
}

// DeleteEndpointHandler handles endpoint deletion
func (h *EndpointsHandler) DeleteEndpointHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	name := vars["name"]
	
	if r.Method != "DELETE" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := h.runner.runner.DeleteGithubEndpoint(ctx, name); err != nil {
		slog.ErrorContext(ctx, "Failed to delete endpoint", "error", err)
		http.Error(w, "Failed to delete endpoint", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}