package handlers

import (
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/cloudbase/garm/params"
	"github.com/cloudbase/garm/runner"
	"github.com/cloudbase/garm/web/templates"
	"github.com/gorilla/mux"
)

type WebHandler struct {
	runner    *runner.Runner
	templates *template.Template
}

type PageData struct {
	Title        string
	PageTitle    string
	ActivePage   string
	Flash        *FlashMessage
	CreateButton *CreateButton
	Stats        *DashboardStats
	IsDetailView bool
	Entity       interface{}
	ReturnPath   string
	ReturnLabel  string
}

type FlashMessage struct {
	Type    string
	Message string
}

type CreateButton struct {
	URL  string
	Text string
}

type DashboardStats struct {
	Repositories  int
	Organizations int
	Pools         int
	Instances     int
}

func NewWebHandler(r *runner.Runner) (*WebHandler, error) {
	slog.Info("Creating new web handler")
	h := &WebHandler{
		runner: r,
	}

	if err := h.loadTemplates(); err != nil {
		slog.Error("Failed to load templates", "error", err)
		return nil, fmt.Errorf("loading templates: %w", err)
	}

	slog.Info("Web handler created successfully")
	return h, nil
}

// GetTemplates returns the loaded templates
func (h *WebHandler) GetTemplates() *template.Template {
	return h.templates
}

func (h *WebHandler) loadTemplates() error {
	slog.Info("Loading embedded templates")
	tmpl, err := templates.GetTemplates()
	if err != nil {
		return fmt.Errorf("failed to load embedded templates: %w", err)
	}

	if tmpl == nil || len(tmpl.Templates()) == 0 {
		return fmt.Errorf("no templates found in embedded assets")
	}

	slog.Info("Successfully loaded embedded templates", "count", len(tmpl.Templates()))
	h.templates = tmpl
	return nil
}

func (h *WebHandler) DashboardHandler(w http.ResponseWriter, r *http.Request) {
	slog.Info("Dashboard handler called", "path", r.URL.Path)

	// First check if templates are loaded
	if h.templates == nil {
		slog.Error("Templates not loaded")
		http.Error(w, "Templates not loaded", http.StatusInternalServerError)
		return
	}

	ctx := r.Context()

	// Get stats
	repos, err := h.runner.ListRepositories(ctx, params.RepositoryFilter{})
	if err != nil {
		slog.Error("Failed to list repositories", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	orgs, err := h.runner.ListOrganizations(ctx, params.OrganizationFilter{})
	if err != nil {
		slog.Error("Failed to list organizations", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	pools, err := h.runner.ListAllPools(ctx)
	if err != nil {
		slog.Error("Failed to list pools", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	instances, err := h.runner.ListAllInstances(ctx)
	if err != nil {
		slog.Error("Failed to list instances", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := PageData{
		Title:      "Dashboard",
		PageTitle:  "Dashboard",
		ActivePage: "dashboard",
		Stats: &DashboardStats{
			Repositories:  len(repos),
			Organizations: len(orgs),
			Pools:         len(pools),
			Instances:     len(instances),
		},
	}

	slog.Info("Executing template", "template", "base.html")
	if err := h.templates.ExecuteTemplate(w, "base.html", data); err != nil {
		slog.Error("Failed to execute template", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	slog.Info("Dashboard rendered successfully")
}

func (h *WebHandler) RepositoriesHandler(w http.ResponseWriter, r *http.Request) {
	data := PageData{
		Title:      "Repositories",
		PageTitle:  "Repositories",
		ActivePage: "repositories",
		CreateButton: &CreateButton{
			URL:  "/web/repositories/new",
			Text: "Add Repository",
		},
	}

	if err := h.templates.ExecuteTemplate(w, "base.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *WebHandler) RepositoryDetailHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	repoID := vars["id"]

	repo, err := h.runner.GetRepositoryByID(ctx, repoID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	data := PageData{
		Title:        repo.Name + " - Repository",
		PageTitle:    repo.Name,
		ActivePage:   "repositories",
		IsDetailView: true,
		Entity:       repo,
	}

	if err := h.templates.ExecuteTemplate(w, "base.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *WebHandler) RepositoriesAPIHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Parse query parameters for search and pagination
	search := r.FormValue("search")
	page := 1
	perPage := 25

	if pageStr := r.FormValue("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if perPageStr := r.FormValue("per_page"); perPageStr != "" {
		if pp, err := strconv.Atoi(perPageStr); err == nil && pp > 0 && pp <= 100 {
			perPage = pp
		}
	}

	repos, err := h.runner.ListRepositories(ctx, params.RepositoryFilter{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Filter repositories by search term
	filteredRepos := repos
	if search != "" {
		filteredRepos = []params.Repository{}
		for _, repo := range repos {
			if strings.Contains(strings.ToLower(repo.Name), strings.ToLower(search)) ||
				strings.Contains(strings.ToLower(repo.Owner), strings.ToLower(search)) {
				filteredRepos = append(filteredRepos, repo)
			}
		}
	}

	// Apply pagination
	total := len(filteredRepos)
	start := (page - 1) * perPage
	end := start + perPage
	if start >= total {
		start = 0
		end = 0
	} else if end > total {
		end = total
	}

	paginatedRepos := []params.Repository{}
	if start < total {
		paginatedRepos = filteredRepos[start:end]
	}

	w.Header().Set("Content-Type", "text/html")

	if len(paginatedRepos) == 0 {
		if search != "" {
			fmt.Fprintf(w, `<tr><td colspan="6" class="px-6 py-4 text-center text-gray-500 dark:text-gray-400">No repositories found matching "%s"</td></tr>`, search)
		} else {
			fmt.Fprintf(w, `<tr><td colspan="6" class="px-6 py-4 text-center text-gray-500 dark:text-gray-400">No repositories found</td></tr>`)
		}
		return
	}

	// Add a simple pagination hint by checking if we have exactly perPage results
	// If we do, there might be more pages available
	hasMorePages := len(paginatedRepos) == perPage && end < total

	for _, repo := range paginatedRepos {
		status := "stopped"
		statusClass := "bg-red-100 dark:bg-red-900 text-red-800 dark:text-red-200"
		if repo.PoolManagerStatus.IsRunning {
			status = "running"
			statusClass = "bg-green-100 dark:bg-green-900 text-green-800 dark:text-green-200"
		}

		// Get forge type and icon
		forgeType := "GitHub"
		forgeIcon := `<div class="inline-flex w-4 h-4">
			<svg class="w-4 h-4 dark:hidden" width="98" height="96" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 98 96"><path fill-rule="evenodd" clip-rule="evenodd" d="M48.854 0C21.839 0 0 22 0 49.217c0 21.756 13.993 40.172 33.405 46.69 2.427.49 3.316-1.059 3.316-2.362 0-1.141-.08-5.052-.08-9.127-13.59 2.934-16.42-5.867-16.42-5.867-2.184-5.704-5.42-7.17-5.42-7.17-4.448-3.015.324-3.015.324-3.015 4.934.326 7.523 5.052 7.523 5.052 4.367 7.496 11.404 5.378 14.235 4.074.404-3.178 1.699-5.378 3.074-6.6-10.839-1.141-22.243-5.378-22.243-24.283 0-5.378 1.94-9.778 5.014-13.2-.485-1.222-2.184-6.275.486-13.038 0 0 4.125-1.304 13.426 5.052a46.97 46.97 0 0 1 12.214-1.63c4.125 0 8.33.571 12.213 1.63 9.302-6.356 13.427-5.052 13.427-5.052 2.67 6.763.97 11.816.485 13.038 3.155 3.422 5.015 7.822 5.015 13.2 0 18.905-11.404 23.06-22.324 24.283 1.78 1.548 3.316 4.481 3.316 9.126 0 6.6-.08 11.897-.08 13.526 0 1.304.89 2.853 3.316 2.364 19.412-6.52 33.405-24.935 33.405-46.691C97.707 22 75.788 0 48.854 0z" fill="#24292f"/></svg>
			<svg class="w-4 h-4 hidden dark:block" width="98" height="96" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 98 96"><path fill-rule="evenodd" clip-rule="evenodd" d="M48.854 0C21.839 0 0 22 0 49.217c0 21.756 13.993 40.172 33.405 46.69 2.427.49 3.316-1.059 3.316-2.362 0-1.141-.08-5.052-.08-9.127-13.59 2.934-16.42-5.867-16.42-5.867-2.184-5.704-5.42-7.17-5.42-7.17-4.448-3.015.324-3.015.324-3.015 4.934.326 7.523 5.052 7.523 5.052 4.367 7.496 11.404 5.378 14.235 4.074.404-3.178 1.699-5.378 3.074-6.6-10.839-1.141-22.243-5.378-22.243-24.283 0-5.378 1.94-9.778 5.014-13.2-.485-1.222-2.184-6.275.486-13.038 0 0 4.125-1.304 13.426 5.052a46.97 46.97 0 0 1 12.214-1.63c4.125 0 8.33.571 12.213 1.63 9.302-6.356 13.427-5.052 13.427-5.052 2.67 6.763.97 11.816.485 13.038 3.155 3.422 5.015 7.822 5.015 13.2 0 18.905-11.404 23.06-22.324 24.283 1.78 1.548 3.316 4.481 3.316 9.126 0 6.6-.08 11.897-.08 13.526 0 1.304.89 2.853 3.316 2.364 19.412-6.52 33.405-24.935 33.405-46.691C97.707 22 75.788 0 48.854 0z" fill="#fff"/></svg>
		</div>`
		
		if repo.Endpoint.EndpointType == "gitea" {
			forgeType = "Gitea"
			forgeIcon = `<svg class="w-4 h-4" xmlns="http://www.w3.org/2000/svg" xml:space="preserve" viewBox="0 0 640 640"><path d="m395.9 484.2-126.9-61c-12.5-6-17.9-21.2-11.8-33.8l61-126.9c6-12.5 21.2-17.9 33.8-11.8 17.2 8.3 27.1 13 27.1 13l-.1-109.2 16.7-.1.1 117.1s57.4 24.2 83.1 40.1c3.7 2.3 10.2 6.8 12.9 14.4 2.1 6.1 2 13.1-1 19.3l-61 126.9c-6.2 12.7-21.4 18.1-33.9 12" style="fill:#fff"/><path d="M622.7 149.8c-4.1-4.1-9.6-4-9.6-4s-117.2 6.6-177.9 8c-13.3.3-26.5.6-39.6.7v117.2c-5.5-2.6-11.1-5.3-16.6-7.9 0-36.4-.1-109.2-.1-109.2-29 .4-89.2-2.2-89.2-2.2s-141.4-7.1-156.8-8.5c-9.8-.6-22.5-2.1-39 1.5-8.7 1.8-33.5 7.4-53.8 26.9C-4.9 212.4 6.6 276.2 8 285.8c1.7 11.7 6.9 44.2 31.7 72.5 45.8 56.1 144.4 54.8 144.4 54.8s12.1 28.9 30.6 55.5c25 33.1 50.7 58.9 75.7 62 63 0 188.9-.1 188.9-.1s12 .1 28.3-10.3c14-8.5 26.5-23.4 26.5-23.4S547 483 565 451.5c5.5-9.7 10.1-19.1 14.1-28 0 0 55.2-117.1 55.2-231.1-1.1-34.5-9.6-40.6-11.6-42.6M125.6 353.9c-25.9-8.5-36.9-18.7-36.9-18.7S69.6 321.8 60 295.4c-16.5-44.2-1.4-71.2-1.4-71.2s8.4-22.5 38.5-30c13.8-3.7 31-3.1 31-3.1s7.1 59.4 15.7 94.2c7.2 29.2 24.8 77.7 24.8 77.7s-26.1-3.1-43-9.1m300.3 107.6s-6.1 14.5-19.6 15.4c-5.8.4-10.3-1.2-10.3-1.2s-.3-.1-5.3-2.1l-112.9-55s-10.9-5.7-12.8-15.6c-2.2-8.1 2.7-18.1 2.7-18.1L322 273s4.8-9.7 12.2-13c.6-.3 2.3-1 4.5-1.5 8.1-2.1 18 2.8 18 2.8L467.4 315s12.6 5.7 15.3 16.2c1.9 7.4-.5 14-1.8 17.2-6.3 15.4-55 113.1-55 113.1" style="fill:#609926"/><path d="M326.8 380.1c-8.2.1-15.4 5.8-17.3 13.8s2 16.3 9.1 20c7.7 4 17.5 1.8 22.7-5.4 5.1-7.1 4.3-16.9-1.8-23.1l24-49.1c1.5.1 3.7.2 6.2-.5 4.1-.9 7.1-3.6 7.1-3.6 4.2 1.8 8.6 3.8 13.2 6.1 4.8 2.4 9.3 4.9 13.4 7.3.9.5 1.8 1.1 2.8 1.9 1.6 1.3 3.4 3.1 4.7 5.5 1.9 5.5-1.9 14.9-1.9 14.9-2.3 7.6-18.4 40.6-18.4 40.6-8.1-.2-15.3 5-17.7 12.5-2.6 8.1 1.1 17.3 8.9 21.3s17.4 1.7 22.5-5.3c5-6.8 4.6-16.3-1.1-22.6 1.9-3.7 3.7-7.4 5.6-11.3 5-10.4 13.5-30.4 13.5-30.4.9-1.7 5.7-10.3 2.7-21.3-2.5-11.4-12.6-16.7-12.6-16.7-12.2-7.9-29.2-15.2-29.2-15.2s0-4.1-1.1-7.1c-1.1-3.1-2.8-5.1-3.9-6.3 4.7-9.7 9.4-19.3 14.1-29-4.1-2-8.1-4-12.2-6.1-4.8 9.8-9.7 19.7-14.5 29.5-6.7-.1-12.9 3.5-16.1 9.4-3.4 6.3-2.7 14.1 1.9 19.8z" style="fill:#609926"/></svg>`
		}

		fmt.Fprintf(w, `
			<tr class="hover:bg-gray-50 dark:hover:bg-gray-700">
				<td class="px-6 py-4 whitespace-nowrap">
					<div class="flex items-center">
						<div>
							<a href="/web/repositories/%s" class="text-sm font-medium text-blue-600 dark:text-blue-400 hover:text-blue-800 dark:hover:text-blue-300 hover:underline">%s</a>
						</div>
					</div>
				</td>
				<td class="px-3 py-4 whitespace-nowrap text-sm text-gray-500 dark:text-gray-400">%s</td>
				<td class="px-3 py-4 whitespace-nowrap">
					<div class="flex items-center space-x-2">
						<div class="flex items-center text-gray-600 dark:text-gray-400" title="%s">
							%s
						</div>
						<span class="text-sm text-gray-500 dark:text-gray-400">%s</span>
					</div>
				</td>
				<td class="px-3 py-4 whitespace-nowrap text-sm text-gray-500 dark:text-gray-400">%s</td>
				<td class="px-3 py-4 whitespace-nowrap">
					<span class="inline-flex px-2 py-1 text-xs font-semibold rounded-full %s">%s</span>
				</td>
				<td class="px-3 py-4 whitespace-nowrap text-right text-sm font-medium">
					<div class="flex justify-end space-x-2">
						<button hx-get="/web/repositories/%s/edit" 
								hx-target="#modal-container" 
								class="text-green-600 dark:text-green-400 hover:text-green-900 dark:hover:text-green-300">Edit</button>
						<button hx-delete="/web/api/repositories/%s" 
								hx-target="#repositories-table" 
								hx-confirm="Are you sure you want to delete this repository?"
								class="text-red-600 dark:text-red-400 hover:text-red-900 dark:hover:text-red-300">Delete</button>
					</div>
				</td>
			</tr>`,
			repo.ID, repo.Name, repo.Owner, forgeType, forgeIcon, repo.Endpoint.Name, repo.CredentialsName, 
			statusClass, status, repo.ID, repo.ID)
	}
	
	// Add a data attribute to help with pagination detection
	if hasMorePages {
		fmt.Fprintf(w, `<tr style="display:none" data-has-more="true"></tr>`)
	}
}

func (h *WebHandler) NewRepositoryHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get both GitHub and Gitea credentials for the form
	githubCreds, err := h.runner.ListCredentials(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	giteaCreds, err := h.runner.ListGiteaCredentials(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	// Combine all credentials
	allCreds := append(githubCreds, giteaCreds...)

	data := struct {
		Repository  *params.Repository
		Credentials []params.ForgeCredentials
	}{
		Credentials: allCreds,
	}

	if err := h.templates.ExecuteTemplate(w, "repository-form.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *WebHandler) EditRepositoryHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	repoID := vars["id"]

	repo, err := h.runner.GetRepositoryByID(ctx, repoID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Get credentials for the form
	creds, err := h.runner.ListCredentials(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := struct {
		Repository  *params.Repository
		Credentials []params.ForgeCredentials
	}{
		Repository:  &repo,
		Credentials: creds,
	}

	if err := h.templates.ExecuteTemplate(w, "repository-form.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *WebHandler) CreateRepositoryHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	// Check if webhook secret should be auto-generated
	webhookSecret := r.FormValue("webhook_secret")
	autoGenerateSecret := r.FormValue("auto_generate_secret") == "on"
	if autoGenerateSecret && webhookSecret == "" {
		// Leave webhook secret empty - it will be auto-generated by the backend
		webhookSecret = ""
	}

	createParams := params.CreateRepoParams{
		Name:            r.FormValue("name"),
		Owner:           r.FormValue("owner"),
		CredentialsName: r.FormValue("credentials_name"),
		WebhookSecret:   webhookSecret,
	}

	repo, err := h.runner.CreateRepository(ctx, createParams)
	if err != nil {
		w.Header().Set("HX-Trigger", "repositoryCreateError")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Install webhook if requested
	installWebhook := r.FormValue("install_webhook") == "on"
	if installWebhook {
		webhookParams := params.InstallWebhookParams{
			WebhookEndpointType: params.WebhookEndpointDirect,
			InsecureSSL:         false, // Default to secure SSL
		}

		_, err := h.runner.InstallRepoWebhook(ctx, repo.ID, webhookParams)
		if err != nil {
			// Log the error but don't fail the repository creation
			slog.ErrorContext(ctx, "Failed to install webhook for repository", "error", err, "repo_id", repo.ID)
			// We could optionally return this as a warning to the user
		}
	}

	// Close modal and trigger refresh
	w.Header().Set("HX-Trigger", "repositoryCreated")
	w.WriteHeader(http.StatusOK)
}

// RepositoryEventsAPIHandler returns events for a specific repository
func (h *WebHandler) RepositoryEventsAPIHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	repoID := vars["id"]

	repo, err := h.runner.GetRepositoryByID(ctx, repoID)
	if err != nil {
		http.Error(w, "Repository not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	
	if len(repo.Events) == 0 {
		w.Write([]byte(`
			<div class="px-4 py-8 text-center">
				<svg class="w-12 h-12 text-gray-400 dark:text-gray-500 mx-auto mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"></path>
				</svg>
				<p class="text-sm text-gray-500 dark:text-gray-400">No events found for this repository.</p>
				<p class="text-xs text-gray-400 dark:text-gray-500 mt-1">Events will appear here as they occur.</p>
			</div>
		`))
		return
	}

	// Render events in reverse chronological order (newest first)
	events := make([]params.EntityEvent, len(repo.Events))
	copy(events, repo.Events)
	
	// Sort events by timestamp (newest first)
	for i := 0; i < len(events)/2; i++ {
		events[i], events[len(events)-1-i] = events[len(events)-1-i], events[i]
	}

	w.Write([]byte(`<div class="divide-y divide-gray-200 dark:divide-gray-700 max-h-96 overflow-y-auto">`))
	
	for _, event := range events {
		levelClass := "bg-gray-100 text-gray-800 dark:bg-gray-900 dark:text-gray-200"
		eventLevel := string(event.EventLevel)
		switch strings.ToLower(eventLevel) {
		case "info":
			levelClass = "bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-200"
		case "warning":
			levelClass = "bg-yellow-100 text-yellow-800 dark:bg-yellow-900 dark:text-yellow-200"
		case "error":
			levelClass = "bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-200"
		}

		// Capitalize first letter of event level
		displayLevel := eventLevel
		if len(displayLevel) > 0 {
			displayLevel = strings.ToUpper(displayLevel[:1]) + strings.ToLower(displayLevel[1:])
		}

		html := fmt.Sprintf(`
			<div class="px-4 py-3 hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors">
				<div class="flex items-start justify-between">
					<div class="flex-1">
						<div class="flex items-center space-x-2">
							<span class="inline-flex px-2 py-1 text-xs font-semibold rounded-full %s">
								%s
							</span>
							<span class="text-sm text-gray-500 dark:text-gray-400">%s</span>
						</div>
						<p class="mt-1 text-sm text-gray-900 dark:text-white">%s</p>
					</div>
				</div>
			</div>
		`, levelClass, displayLevel, event.CreatedAt.Format("Jan 2, 15:04:05"), event.Message)
		
		w.Write([]byte(html))
	}
	
	w.Write([]byte(`</div>`))
}

func (h *WebHandler) UpdateRepositoryHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	repoID := vars["id"]

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	credName := r.FormValue("credentials_name")
	webhookSecret := r.FormValue("webhook_secret")
	poolBalancerType := r.FormValue("pool_balancer_type")

	updateParams := params.UpdateEntityParams{
		CredentialsName:  credName,
		WebhookSecret:   webhookSecret,
		PoolBalancerType: params.PoolBalancerType(poolBalancerType),
	}

	_, err := h.runner.UpdateRepository(ctx, repoID, updateParams)
	if err != nil {
		w.Header().Set("HX-Trigger", "repositoryUpdateError")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Close modal and trigger refresh
	w.Header().Set("HX-Trigger", "repositoryUpdated")
	w.WriteHeader(http.StatusOK)
}

func (h *WebHandler) DeleteRepositoryHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	repoID := vars["id"]

	err := h.runner.DeleteRepository(ctx, repoID, false)
	if err != nil {
		w.Header().Set("HX-Trigger", "repositoryDeleteError")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("HX-Trigger", "repositoryDeleted")
	w.WriteHeader(http.StatusOK)
}

func (h *WebHandler) RepositoryPoolsAPIHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	repoID := vars["id"]

	pools, err := h.runner.ListRepoPools(ctx, repoID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	
	if len(pools) == 0 {
		fmt.Fprintf(w, `<div class="px-4 py-8 text-center">
			<p class="text-sm text-gray-500 dark:text-gray-400">No pools found for this repository.</p>
		</div>`)
		return
	}

	fmt.Fprintf(w, `<div class="w-full">
		<table class="w-full table-fixed divide-y divide-gray-200 dark:divide-gray-700">
			<thead class="bg-gray-50 dark:bg-gray-700">
				<tr>
					<th class="w-1/3 px-3 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">Pool ID</th>
					<th class="w-1/3 px-3 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">Image</th>
					<th class="w-20 px-3 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">Provider</th>
					<th class="w-16 px-3 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">Status</th>
				</tr>
			</thead>
			<tbody class="divide-y divide-gray-200 dark:divide-gray-700">`)

	for _, pool := range pools {
		status := "disabled"
		statusClass := "bg-red-100 dark:bg-red-900 text-red-800 dark:text-red-200"
		if pool.Enabled {
			status = "enabled"
			statusClass = "bg-green-100 dark:bg-green-900 text-green-800 dark:text-green-200"
		}

		fmt.Fprintf(w, `
			<tr class="hover:bg-gray-50 dark:hover:bg-gray-700">
				<td class="px-3 py-4 text-sm font-mono" title="%s">
					<div class="truncate">
						<a href="/web/pools/%s/details?from=repository&entity_id=%s" 
								class="text-purple-600 dark:text-purple-400 hover:text-purple-900 dark:hover:text-purple-300 hover:underline">%s</a>
					</div>
				</td>
				<td class="px-3 py-4 text-sm text-gray-900 dark:text-white" title="%s">
					<div class="truncate">%s</div>
				</td>
				<td class="px-3 py-4 text-sm text-gray-900 dark:text-white">
					<div class="truncate">%s</div>
				</td>
				<td class="px-3 py-4">
					<span class="inline-flex px-2 py-1 text-xs font-semibold rounded-full %s">%s</span>
				</td>
			</tr>`,
			pool.ID, pool.ID, repoID, pool.ID, pool.Image, pool.Image, pool.ProviderName, statusClass, status)
	}

	fmt.Fprintf(w, `
			</tbody>
		</table>
	</div>`)
}

func (h *WebHandler) RepositoryInstancesAPIHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	repoID := vars["id"]

	instances, err := h.runner.ListRepoInstances(ctx, repoID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	
	if len(instances) == 0 {
		fmt.Fprintf(w, `<div class="px-4 py-8 text-center">
			<p class="text-sm text-gray-500 dark:text-gray-400">No instances found for this repository.</p>
		</div>`)
		return
	}

	fmt.Fprintf(w, `<div class="overflow-x-auto">
		<table class="min-w-full divide-y divide-gray-200 dark:divide-gray-700">
			<thead class="bg-gray-50 dark:bg-gray-700">
				<tr>
					<th class="px-3 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">Name</th>
					<th class="px-3 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">Status</th>
					<th class="px-3 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">Runner Status</th>
					<th class="px-3 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">Created</th>
					<th class="px-3 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">Actions</th>
				</tr>
			</thead>
			<tbody class="divide-y divide-gray-200 dark:divide-gray-700">`)

	for _, instance := range instances {
		status := string(instance.Status)
		statusClass := "bg-gray-100 dark:bg-gray-700 text-gray-800 dark:text-gray-200"
		
		switch status {
		case "running":
			statusClass = "bg-green-100 dark:bg-green-900 text-green-800 dark:text-green-200"
		case "stopped", "error":
			statusClass = "bg-red-100 dark:bg-red-900 text-red-800 dark:text-red-200"
		case "pending":
			statusClass = "bg-yellow-100 dark:bg-yellow-900 text-yellow-800 dark:text-yellow-200"
		}

		fmt.Fprintf(w, `
			<tr class="hover:bg-gray-50 dark:hover:bg-gray-700">
				<td class="px-3 py-4 whitespace-nowrap text-sm font-medium text-gray-900 dark:text-white">%s</td>
				<td class="px-3 py-4 whitespace-nowrap">
					<span class="inline-flex px-2 py-1 text-xs font-semibold rounded-full %s">%s</span>
				</td>
				<td class="px-3 py-4 whitespace-nowrap">
					<span class="inline-flex px-2 py-1 text-xs font-semibold rounded-full bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-200">%s</span>
				</td>
				<td class="px-3 py-4 whitespace-nowrap text-sm text-gray-500 dark:text-gray-400">%s</td>
				<td class="px-3 py-4 whitespace-nowrap text-right text-sm font-medium">
					<a href="/web/instances/%s/detail?from=repository&entity_id=%s" 
					   class="text-purple-600 dark:text-purple-400 hover:text-purple-900 dark:hover:text-purple-300 hover:underline">
						View
					</a>
				</td>
			</tr>`,
			instance.Name, statusClass, status, string(instance.RunnerStatus), instance.CreatedAt.Format("Jan 2, 15:04"), instance.Name, repoID)
	}

	fmt.Fprintf(w, `
			</tbody>
		</table>
	</div>`)
}

func (h *WebHandler) RecentActivityHandler(w http.ResponseWriter, r *http.Request) {
	// This would typically fetch from a database of recent events
	// For now, we'll return a placeholder
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, `
		<div class="space-y-3">
			<div class="flex items-center space-x-3 p-3 bg-gray-50 rounded-lg">
				<div class="w-2 h-2 bg-green-500 rounded-full"></div>
				<span class="text-sm text-gray-600">Repository example/repo was updated</span>
				<span class="text-xs text-gray-400">2 minutes ago</span>
			</div>
			<div class="flex items-center space-x-3 p-3 bg-gray-50 rounded-lg">
				<div class="w-2 h-2 bg-blue-500 rounded-full"></div>
				<span class="text-sm text-gray-600">New runner instance started</span>
				<span class="text-xs text-gray-400">5 minutes ago</span>
			</div>
			<div class="flex items-center space-x-3 p-3 bg-gray-50 rounded-lg">
				<div class="w-2 h-2 bg-yellow-500 rounded-full"></div>
				<span class="text-sm text-gray-600">Pool "ubuntu-20.04" scaling up</span>
				<span class="text-xs text-gray-400">10 minutes ago</span>
			</div>
		</div>
	`)
}

// HandleMethodOverride handles method override for forms
func (h *WebHandler) HandleMethodOverride(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			method := r.FormValue("_method")
			if method == "PUT" || method == "DELETE" {
				r.Method = method
			}
		}
		next(w, r)
	}
}
