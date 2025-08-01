package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/cloudbase/garm/params"
	commonParams "github.com/cloudbase/garm-provider-common/params"
)

func (h *WebHandler) ScaleSetsHandler(w http.ResponseWriter, r *http.Request) {
	data := PageData{
		Title:      "Scale Sets",
		PageTitle:  "Scale Sets",
		ActivePage: "scalesets",
		CreateButton: &CreateButton{
			URL:  "/web/scalesets/new",
			Text: "Add Scale Set",
		},
	}

	if err := h.templates.ExecuteTemplate(w, "base.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *WebHandler) ScaleSetsAPIHandler(w http.ResponseWriter, r *http.Request) {
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

	scalesets, err := h.runner.ListAllScaleSets(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Filter scale sets by search term (search by name or entity name)
	filteredScaleSets := scalesets
	if search != "" {
		filteredScaleSets = []params.ScaleSet{}
		for _, scaleset := range scalesets {
			entityName := ""
			if scaleset.RepoName != "" {
				entityName = scaleset.RepoName
			} else if scaleset.OrgName != "" {
				entityName = scaleset.OrgName
			} else if scaleset.EnterpriseName != "" {
				entityName = scaleset.EnterpriseName
			}
			
			if strings.Contains(strings.ToLower(scaleset.Name), strings.ToLower(search)) || 
			   strings.Contains(strings.ToLower(entityName), strings.ToLower(search)) {
				filteredScaleSets = append(filteredScaleSets, scaleset)
			}
		}
	}

	// Apply pagination
	total := len(filteredScaleSets)
	start := (page - 1) * perPage
	end := start + perPage
	if start >= total {
		start = 0
		end = 0
	} else if end > total {
		end = total
	}

	paginatedScaleSets := []params.ScaleSet{}
	if start < total {
		paginatedScaleSets = filteredScaleSets[start:end]
	}

	w.Header().Set("Content-Type", "text/html")
	
	if len(paginatedScaleSets) == 0 {
		if search != "" {
			fmt.Fprintf(w, `<tr><td colspan="7" class="px-3 py-4 text-center text-gray-500 dark:text-gray-400">No scale sets found matching "%s"</td></tr>`, search)
		} else {
			fmt.Fprintf(w, `<tr><td colspan="7" class="px-3 py-4 text-center text-gray-500 dark:text-gray-400">No scale sets found</td></tr>`)
		}
		return
	}

	// Add a simple pagination hint by checking if we have exactly perPage results
	// If we do, there might be more pages available
	hasMorePages := len(paginatedScaleSets) == perPage && end < total

	for _, scaleset := range paginatedScaleSets {
		// Status based on enabled state
		status := "disabled"
		statusClass := "bg-red-100 dark:bg-red-900 text-red-800 dark:text-red-200"
		if scaleset.Enabled {
			status = "enabled"
			statusClass = "bg-green-100 dark:bg-green-900 text-green-800 dark:text-green-200"
		}

		// Determine entity type and name
		entityName := ""
		entityType := ""
		if scaleset.RepoName != "" {
			entityName = scaleset.RepoName
			entityType = "repository"
		} else if scaleset.OrgName != "" {
			entityName = scaleset.OrgName
			entityType = "organization"
		} else if scaleset.EnterpriseName != "" {
			entityName = scaleset.EnterpriseName
			entityType = "enterprise"
		}


		// Truncate image name if too long
		displayImage := scaleset.Image
		imageTitle := scaleset.Image
		if len(displayImage) > 30 {
			displayImage = displayImage[:27] + "..."
		}

		// Truncate provider name if too long
		providerName := scaleset.ProviderName
		if len(providerName) > 15 {
			providerName = providerName[:12] + "..."
		}

		// Count instances
		instanceCount := len(scaleset.Instances)

		fmt.Fprintf(w, `
			<tr class="hover:bg-gray-50 dark:hover:bg-gray-700">
				<td class="px-3 py-4 whitespace-nowrap text-sm text-gray-900 dark:text-white">
					<a href="/web/scalesets/%d/details" 
						class="text-purple-600 dark:text-purple-400 hover:text-purple-900 dark:hover:text-purple-300 hover:underline">%s</a>
				</td>
				<td class="px-3 py-4 whitespace-nowrap text-sm text-gray-900 dark:text-white" title="%s">%s</td>
				<td class="px-3 py-4 whitespace-nowrap text-sm">
					<span onclick="alert('Navigate to %s details - not implemented yet')" class="text-purple-600 dark:text-purple-400 hover:text-purple-900 dark:hover:text-purple-300 cursor-pointer">%s</span>
					<div class="text-xs text-gray-500 dark:text-gray-400">%s</div>
				</td>
				<td class="px-3 py-4 whitespace-nowrap text-sm text-gray-500 dark:text-gray-400" title="%s">%s</td>
				<td class="px-3 py-4 whitespace-nowrap text-sm text-gray-900 dark:text-white">%d</td>
				<td class="px-3 py-4 whitespace-nowrap">
					<span class="inline-flex px-2 py-1 text-xs font-semibold rounded-full %s">%s</span>
				</td>
				<td class="px-3 py-4 whitespace-nowrap text-right text-sm font-medium">
					<button hx-get="/web/scalesets/%d/edit" 
							hx-target="#modal-container" 
							class="text-orange-600 dark:text-orange-400 hover:text-orange-900 dark:hover:text-orange-300 mr-2">Edit</button>
					<button onclick="showDeleteScaleSetConfirm('%d', '%s')"
							class="text-red-600 dark:text-red-400 hover:text-red-900 dark:hover:text-red-300">Delete</button>
				</td>
			</tr>`,
			scaleset.ID, scaleset.Name, imageTitle, displayImage, 
			entityType, entityName, entityType, 
			scaleset.ProviderName, providerName, instanceCount,
			statusClass, status, 
			scaleset.ID, scaleset.ID, scaleset.Name)
	}
	
	// Add a data attribute to help with pagination detection
	if hasMorePages {
		fmt.Fprintf(w, `<tr style="display:none" data-has-more="true"></tr>`)
	}
}

func (h *WebHandler) ScaleSetDetailsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	scalesetIDStr := vars["id"]

	// Convert string to uint
	var scalesetID uint
	if _, err := fmt.Sscanf(scalesetIDStr, "%d", &scalesetID); err != nil {
		http.Error(w, "Invalid scale set ID", http.StatusBadRequest)
		return
	}

	scaleset, err := h.runner.GetScaleSetByID(ctx, scalesetID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Format extra specs JSON for display
	extraSpecsDecoded := ""
	if len(scaleset.ExtraSpecs) > 0 {
		// Pretty format the JSON
		var jsonObj interface{}
		if err := json.Unmarshal(scaleset.ExtraSpecs, &jsonObj); err == nil {
			if formatted, err := json.MarshalIndent(jsonObj, "", "  "); err == nil {
				extraSpecsDecoded = string(formatted)
			}
		} else {
			// If it fails to parse as JSON, just use as string
			extraSpecsDecoded = string(scaleset.ExtraSpecs)
		}
	}

	// Create custom data structure with decoded extra specs
	data := struct {
		PageData
		ExtraSpecsDecoded string
	}{
		PageData: PageData{
			Title:        scaleset.Name + " - Scale Set Details",
			PageTitle:    "Scale Set Details",
			ActivePage:   "scalesets",
			IsDetailView: true,
			Entity:       scaleset,
		},
		ExtraSpecsDecoded: extraSpecsDecoded,
	}

	if err := h.templates.ExecuteTemplate(w, "base.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *WebHandler) ScaleSetEditHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	scalesetIDStr := vars["id"]

	// Convert string to uint
	var scalesetID uint
	if _, err := fmt.Sscanf(scalesetIDStr, "%d", &scalesetID); err != nil {
		http.Error(w, "Invalid scale set ID", http.StatusBadRequest)
		return
	}

	scaleset, err := h.runner.GetScaleSetByID(ctx, scalesetID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Format extra specs JSON for form editing
	extraSpecsDecoded := ""
	if len(scaleset.ExtraSpecs) > 0 {
		// Pretty format the JSON for editing
		var jsonObj interface{}
		if err := json.Unmarshal(scaleset.ExtraSpecs, &jsonObj); err == nil {
			if formatted, err := json.MarshalIndent(jsonObj, "", "  "); err == nil {
				extraSpecsDecoded = string(formatted)
			}
		} else {
			// If it fails to parse as JSON, just use as string
			extraSpecsDecoded = string(scaleset.ExtraSpecs)
		}
	}

	data := struct {
		ScaleSet          interface{}
		ExtraSpecsDecoded string
	}{
		ScaleSet:          scaleset,
		ExtraSpecsDecoded: extraSpecsDecoded,
	}

	if err := h.templates.ExecuteTemplate(w, "scaleset-edit.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *WebHandler) UpdateScaleSetHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	scalesetIDStr := vars["id"]

	// Convert string to uint
	var scalesetID uint
	if _, err := fmt.Sscanf(scalesetIDStr, "%d", &scalesetID); err != nil {
		http.Error(w, "Invalid scale set ID", http.StatusBadRequest)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	var updateParams params.UpdateScaleSetParams

	// Parse form fields
	if name := r.FormValue("name"); name != "" {
		updateParams.Name = name
	}

	if runnerGroup := r.FormValue("github_runner_group"); runnerGroup != "" {
		updateParams.GitHubRunnerGroup = &runnerGroup
	}

	if maxRunners := r.FormValue("max_runners"); maxRunners != "" {
		var maxRunnersUint uint
		if _, err := fmt.Sscanf(maxRunners, "%d", &maxRunnersUint); err == nil {
			updateParams.MaxRunners = &maxRunnersUint
		}
	}

	if minIdleRunners := r.FormValue("min_idle_runners"); minIdleRunners != "" {
		var minIdleRunnersUint uint
		if _, err := fmt.Sscanf(minIdleRunners, "%d", &minIdleRunnersUint); err == nil {
			updateParams.MinIdleRunners = &minIdleRunnersUint
		}
	}

	if bootstrapTimeout := r.FormValue("runner_bootstrap_timeout"); bootstrapTimeout != "" {
		var bootstrapTimeoutUint uint
		if _, err := fmt.Sscanf(bootstrapTimeout, "%d", &bootstrapTimeoutUint); err == nil {
			updateParams.RunnerBootstrapTimeout = &bootstrapTimeoutUint
		}
	}

	// DisableUpdate is not available in UpdateScaleSetParams, it's read-only

	enabled := r.FormValue("enabled") == "on"
	updateParams.Enabled = &enabled

	_, err := h.runner.UpdateScaleSetByID(ctx, scalesetID, updateParams)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.ScaleSetsAPIHandler(w, r)
}

func (h *WebHandler) DeleteScaleSetHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	scalesetIDStr := vars["id"]

	// Convert string to uint
	var scalesetID uint
	if _, err := fmt.Sscanf(scalesetIDStr, "%d", &scalesetID); err != nil {
		http.Error(w, "Invalid scale set ID", http.StatusBadRequest)
		return
	}

	err := h.runner.DeleteScaleSetByID(ctx, scalesetID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.ScaleSetsAPIHandler(w, r)
}

func (h *WebHandler) NewScaleSetHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get providers
	providers, err := h.runner.ListProviders(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get repositories (scale sets are GitHub only)
	repos, err := h.runner.ListRepositories(ctx, params.RepositoryFilter{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Filter for GitHub repositories only
	githubRepos := []params.Repository{}
	for _, repo := range repos {
		if repo.Endpoint.EndpointType == params.GithubEndpointType {
			githubRepos = append(githubRepos, repo)
		}
	}

	// Get organizations (GitHub only)
	orgs, err := h.runner.ListOrganizations(ctx, params.OrganizationFilter{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Filter for GitHub organizations only
	githubOrgs := []params.Organization{}
	for _, org := range orgs {
		if org.Endpoint.EndpointType == params.GithubEndpointType {
			githubOrgs = append(githubOrgs, org)
		}
	}

	// Get enterprises (GitHub only)
	enterprises, err := h.runner.ListEnterprises(ctx, params.EnterpriseFilter{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// All enterprises are GitHub by definition
	data := struct {
		ScaleSet      *params.ScaleSet
		Providers     []params.Provider
		Repositories  []params.Repository
		Organizations []params.Organization
		Enterprises   []params.Enterprise
	}{
		Providers:     providers,
		Repositories:  githubRepos,
		Organizations: githubOrgs,
		Enterprises:   enterprises,
	}

	if err := h.templates.ExecuteTemplate(w, "scaleset-form.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *WebHandler) CreateScaleSetHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	// Parse entity level and ID
	entityLevel := r.FormValue("entity_level")
	entityID := r.FormValue("entity_id")

	if entityLevel == "" || entityID == "" {
		http.Error(w, "Entity level and ID are required", http.StatusBadRequest)
		return
	}

	// Parse numeric fields
	var maxRunners, minIdleRunners, bootstrapTimeout uint
	if mr := r.FormValue("max_runners"); mr != "" {
		if parsed, err := strconv.ParseUint(mr, 10, 32); err == nil {
			maxRunners = uint(parsed)
		}
	}
	if mir := r.FormValue("min_idle_runners"); mir != "" {
		if parsed, err := strconv.ParseUint(mir, 10, 32); err == nil {
			minIdleRunners = uint(parsed)
		}
	}
	if bt := r.FormValue("runner_bootstrap_timeout"); bt != "" {
		if parsed, err := strconv.ParseUint(bt, 10, 32); err == nil {
			bootstrapTimeout = uint(parsed)
		}
	}

	// Parse tags
	var tags []string
	if tagsStr := r.FormValue("tags"); tagsStr != "" {
		for _, tag := range strings.Split(tagsStr, ",") {
			if trimmed := strings.TrimSpace(tag); trimmed != "" {
				tags = append(tags, trimmed)
			}
		}
	}

	createParams := params.CreateScaleSetParams{
		RunnerPrefix: params.RunnerPrefix{
			Prefix: r.FormValue("runner_prefix"),
		},
		Name:                     r.FormValue("name"),
		ProviderName:             r.FormValue("provider_name"),
		Image:                    r.FormValue("image"),
		Flavor:                   r.FormValue("flavor"),
		MaxRunners:               maxRunners,
		MinIdleRunners:           minIdleRunners,
		RunnerBootstrapTimeout:   bootstrapTimeout,
		Tags:                     tags,
		OSType:                   commonParams.OSType(r.FormValue("os_type")),
		OSArch:                   commonParams.OSArch(r.FormValue("os_arch")),
		GitHubRunnerGroup:        r.FormValue("github_runner_group"),
		Enabled:                  r.FormValue("enabled") == "on",
	}

	// Call appropriate creation function based on entity level
	var entityType params.ForgeEntityType
	switch entityLevel {
	case "repo":
		entityType = params.ForgeEntityTypeRepository
	case "org":
		entityType = params.ForgeEntityTypeOrganization
	case "enterprise":
		entityType = params.ForgeEntityTypeEnterprise
	default:
		http.Error(w, "Invalid entity level", http.StatusBadRequest)
		return
	}

	_, err := h.runner.CreateEntityScaleSet(ctx, entityType, entityID, createParams)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Close modal and refresh table
	w.Header().Set("HX-Trigger", "closeModal")
	h.ScaleSetsAPIHandler(w, r)
}