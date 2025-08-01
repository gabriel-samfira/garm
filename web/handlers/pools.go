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

func (h *WebHandler) PoolsHandler(w http.ResponseWriter, r *http.Request) {
	data := PageData{
		Title:      "Pools",
		PageTitle:  "Runner Pools",
		ActivePage: "pools",
		CreateButton: &CreateButton{
			URL:  "/web/pools/new",
			Text: "Add Pool",
		},
	}

	if err := h.templates.ExecuteTemplate(w, "base.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *WebHandler) PoolsAPIHandler(w http.ResponseWriter, r *http.Request) {
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

	pools, err := h.runner.ListAllPools(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Filter pools by search term (search by entity name)
	filteredPools := pools
	if search != "" {
		filteredPools = []params.Pool{}
		for _, pool := range pools {
			entityName := ""
			if pool.RepoID != "" {
				entityName = pool.RepoName
			} else if pool.OrgID != "" {
				entityName = pool.OrgName
			} else if pool.EnterpriseID != "" {
				entityName = pool.EnterpriseName
			}
			
			if strings.Contains(strings.ToLower(entityName), strings.ToLower(search)) {
				filteredPools = append(filteredPools, pool)
			}
		}
	}

	// Apply pagination
	total := len(filteredPools)
	start := (page - 1) * perPage
	end := start + perPage
	if start >= total {
		start = 0
		end = 0
	} else if end > total {
		end = total
	}

	paginatedPools := []params.Pool{}
	if start < total {
		paginatedPools = filteredPools[start:end]
	}

	w.Header().Set("Content-Type", "text/html")
	
	if len(paginatedPools) == 0 {
		if search != "" {
			fmt.Fprintf(w, `<tr><td colspan="6" class="px-3 py-4 text-center text-gray-500 dark:text-gray-400">No pools found matching "%s"</td></tr>`, search)
		} else {
			fmt.Fprintf(w, `<tr><td colspan="6" class="px-3 py-4 text-center text-gray-500 dark:text-gray-400">No pools found</td></tr>`)
		}
		return
	}

	// Add a simple pagination hint by checking if we have exactly perPage results
	// If we do, there might be more pages available
	hasMorePages := len(paginatedPools) == perPage && end < total

	for _, pool := range paginatedPools {
		// Status
		status := "disabled"
		statusClass := "bg-red-100 dark:bg-red-900 text-red-800 dark:text-red-200"
		if pool.Enabled {
			status = "enabled"
			statusClass = "bg-green-100 dark:bg-green-900 text-green-800 dark:text-green-200"
		}

		// Entity info and clickable entity
		entityName := "Unknown"
		entityType := ""
		if pool.RepoID != "" {
			entityName = pool.RepoName
			entityType = "repository"
		} else if pool.OrgID != "" {
			entityName = pool.OrgName
			entityType = "organization" 
		} else if pool.EnterpriseID != "" {
			entityName = pool.EnterpriseName
			entityType = "enterprise"
		}

		// Forge icon - GitHub or Gitea - inline with endpoint name
		forgeIcon := ""
		if pool.Endpoint.EndpointType == params.GithubEndpointType {
			forgeIcon = `<svg class="w-4 h-4 inline mr-2 dark:hidden" width="98" height="96" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 98 96"><path fill-rule="evenodd" clip-rule="evenodd" d="M48.854 0C21.839 0 0 22 0 49.217c0 21.756 13.993 40.172 33.405 46.69 2.427.49 3.316-1.059 3.316-2.362 0-1.141-.08-5.052-.08-9.127-13.59 2.934-16.42-5.867-16.42-5.867-2.184-5.704-5.42-7.17-5.42-7.17-4.448-3.015.324-3.015.324-3.015 4.934.326 7.523 5.052 7.523 5.052 4.367 7.496 11.404 5.378 14.235 4.074.404-3.178 1.699-5.378 3.074-6.6-10.839-1.141-22.243-5.378-22.243-24.283 0-5.378 1.94-9.778 5.014-13.2-.485-1.222-2.184-6.275.486-13.038 0 0 4.125-1.304 13.426 5.052a46.97 46.97 0 0 1 12.214-1.63c4.125 0 8.33.571 12.213 1.63 9.302-6.356 13.427-5.052 13.427-5.052 2.67 6.763.97 11.816.485 13.038 3.155 3.422 5.015 7.822 5.015 13.2 0 18.905-11.404 23.06-22.324 24.283 1.78 1.548 3.316 4.481 3.316 9.126 0 6.6-.08 11.897-.08 13.526 0 1.304.89 2.853 3.316 2.364 19.412-6.52 33.405-24.935 33.405-46.691C97.707 22 75.788 0 48.854 0z" fill="#24292f"/></svg><svg class="w-4 h-4 inline mr-2 hidden dark:block" width="98" height="96" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 98 96"><path fill-rule="evenodd" clip-rule="evenodd" d="M48.854 0C21.839 0 0 22 0 49.217c0 21.756 13.993 40.172 33.405 46.69 2.427.49 3.316-1.059 3.316-2.362 0-1.141-.08-5.052-.08-9.127-13.59 2.934-16.42-5.867-16.42-5.867-2.184-5.704-5.42-7.17-5.42-7.17-4.448-3.015.324-3.015.324-3.015 4.934.326 7.523 5.052 7.523 5.052 4.367 7.496 11.404 5.378 14.235 4.074.404-3.178 1.699-5.378 3.074-6.6-10.839-1.141-22.243-5.378-22.243-24.283 0-5.378 1.94-9.778 5.014-13.2-.485-1.222-2.184-6.275.486-13.038 0 0 4.125-1.304 13.426 5.052a46.97 46.97 0 0 1 12.214-1.63c4.125 0 8.33.571 12.213 1.63 9.302-6.356 13.427-5.052 13.427-5.052 2.67 6.763.97 11.816.485 13.038 3.155 3.422 5.015 7.822 5.015 13.2 0 18.905-11.404 23.06-22.324 24.283 1.78 1.548 3.316 4.481 3.316 9.126 0 6.6-.08 11.897-.08 13.526 0 1.304.89 2.853 3.316 2.364 19.412-6.52 33.405-24.935 33.405-46.691C97.707 22 75.788 0 48.854 0z" fill="#fff"/></svg>`
		} else if pool.Endpoint.EndpointType == params.GiteaEndpointType {
			forgeIcon = `<svg class="w-4 h-4 inline mr-2" xmlns="http://www.w3.org/2000/svg" xml:space="preserve" viewBox="0 0 640 640"><path d="m395.9 484.2-126.9-61c-12.5-6-17.9-21.2-11.8-33.8l61-126.9c6-12.5 21.2-17.9 33.8-11.8 17.2 8.3 27.1 13 27.1 13l-.1-109.2 16.7-.1.1 117.1s57.4 24.2 83.1 40.1c3.7 2.3 10.2 6.8 12.9 14.4 2.1 6.1 2 13.1-1 19.3l-61 126.9c-6.2 12.7-21.4 18.1-33.9 12" style="fill:#fff"/><path d="M622.7 149.8c-4.1-4.1-9.6-4-9.6-4s-117.2 6.6-177.9 8c-13.3.3-26.5.6-39.6.7v117.2c-5.5-2.6-11.1-5.3-16.6-7.9 0-36.4-.1-109.2-.1-109.2-29 .4-89.2-2.2-89.2-2.2s-141.4-7.1-156.8-8.5c-9.8-.6-22.5-2.1-39 1.5-8.7 1.8-33.5 7.4-53.8 26.9C-4.9 212.4 6.6 276.2 8 285.8c1.7 11.7 6.9 44.2 31.7 72.5 45.8 56.1 144.4 54.8 144.4 54.8s12.1 28.9 30.6 55.5c25 33.1 50.7 58.9 75.7 62 63 0 188.9-.1 188.9-.1s12 .1 28.3-10.3c14-8.5 26.5-23.4 26.5-23.4S547 483 565 451.5c5.5-9.7 10.1-19.1 14.1-28 0 0 55.2-117.1 55.2-231.1-1.1-34.5-9.6-40.6-11.6-42.6M125.6 353.9c-25.9-8.5-36.9-18.7-36.9-18.7S69.6 321.8 60 295.4c-16.5-44.2-1.4-71.2-1.4-71.2s8.4-22.5 38.5-30c13.8-3.7 31-3.1 31-3.1s7.1 59.4 15.7 94.2c7.2 29.2 24.8 77.7 24.8 77.7s-26.1-3.1-43-9.1m300.3 107.6s-6.1 14.5-19.6 15.4c-5.8.4-10.3-1.2-10.3-1.2s-.3-.1-5.3-2.1l-112.9-55s-10.9-5.7-12.8-15.6c-2.2-8.1 2.7-18.1 2.7-18.1L322 273s4.8-9.7 12.2-13c.6-.3 2.3-1 4.5-1.5 8.1-2.1 18 2.8 18 2.8L467.4 315s12.6 5.7 15.3 16.2c1.9 7.4-.5 14-1.8 17.2-6.3 15.4-55 113.1-55 113.1" style="fill:#609926"/><path d="M326.8 380.1c-8.2.1-15.4 5.8-17.3 13.8s2 16.3 9.1 20c7.7 4 17.5 1.8 22.7-5.4 5.1-7.1 4.3-16.9-1.8-23.1l24-49.1c1.5.1 3.7.2 6.2-.5 4.1-.9 7.1-3.6 7.1-3.6 4.2 1.8 8.6 3.8 13.2 6.1 4.8 2.4 9.3 4.9 13.4 7.3.9.5 1.8 1.1 2.8 1.9 1.6 1.3 3.4 3.1 4.7 5.5 1.9 5.5-1.9 14.9-1.9 14.9-2.3 7.6-18.4 40.6-18.4 40.6-8.1-.2-15.3 5-17.7 12.5-2.6 8.1 1.1 17.3 8.9 21.3s17.4 1.7 22.5-5.3c5-6.8 4.6-16.3-1.1-22.6 1.9-3.7 3.7-7.4 5.6-11.3 5-10.4 13.5-30.4 13.5-30.4.9-1.7 5.7-10.3 2.7-21.3-2.5-11.4-12.6-16.7-12.6-16.7-12.2-7.9-29.2-15.2-29.2-15.2s0-4.1-1.1-7.1c-1.1-3.1-2.8-5.1-3.9-6.3 4.7-9.7 9.4-19.3 14.1-29-4.1-2-8.1-4-12.2-6.1-4.8 9.8-9.7 19.7-14.5 29.5-6.7-.1-12.9 3.5-16.1 9.4-3.4 6.3-2.7 14.1 1.9 19.8z" style="fill:#609926"/></svg>`
		}

		// Truncate image name if too long, with full name in title attribute
		displayImage := pool.Image
		imageTitle := pool.Image
		if len(pool.Image) > 20 {
			displayImage = pool.Image[:17] + "..."
		}

		// Truncate endpoint name
		endpointName := pool.Endpoint.Name
		if len(endpointName) > 15 {
			endpointName = endpointName[:12] + "..."
		}

		fmt.Fprintf(w, `
			<tr class="hover:bg-gray-50 dark:hover:bg-gray-700">
				<td class="px-3 py-4 whitespace-nowrap text-sm font-mono" title="%s">
					<a href="/web/pools/%s/details" 
							class="text-purple-600 dark:text-purple-400 hover:text-purple-900 dark:hover:text-purple-300 hover:underline">%s</a>
				</td>
				<td class="px-3 py-4 whitespace-nowrap text-sm text-gray-900 dark:text-white" title="%s">%s</td>
				<td class="px-3 py-4 whitespace-nowrap text-sm">
					<span onclick="alert('Navigate to %s details - not implemented yet')" class="text-purple-600 dark:text-purple-400 hover:text-purple-900 dark:hover:text-purple-300 cursor-pointer">%s</span>
					<div class="text-xs text-gray-500 dark:text-gray-400">%s</div>
				</td>
				<td class="px-3 py-4 whitespace-nowrap text-sm text-gray-900 dark:text-white" title="%s">
					%s%s
				</td>
				<td class="px-3 py-4 whitespace-nowrap">
					<span class="inline-flex px-2 py-1 text-xs font-semibold rounded-full %s">%s</span>
				</td>
				<td class="px-3 py-4 whitespace-nowrap text-right text-sm font-medium">
					<button hx-get="/web/pools/%s/edit" 
							hx-target="#modal-container" 
							class="text-orange-600 dark:text-orange-400 hover:text-orange-900 dark:hover:text-orange-300 mr-2">Edit</button>
					<button onclick="showDeleteConfirm('%s', '%s')"
							class="text-red-600 dark:text-red-400 hover:text-red-900 dark:hover:text-red-300">Delete</button>
				</td>
			</tr>`,
			pool.ID, pool.ID, pool.ID[:8]+"...", imageTitle, displayImage, entityType, entityName, entityType, 
			pool.Endpoint.Name, forgeIcon, endpointName, statusClass, status, pool.ID, pool.ID, pool.ID)
	}
	
	// Add a data attribute to help with pagination detection
	if hasMorePages {
		fmt.Fprintf(w, `<tr style="display:none" data-has-more="true"></tr>`)
	}
}

func (h *WebHandler) PoolDetailsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	poolID := vars["id"]

	pool, err := h.runner.GetPoolByID(ctx, poolID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Format extra specs JSON for display
	extraSpecsDecoded := ""
	if len(pool.ExtraSpecs) > 0 {
		// Pretty format the JSON
		var jsonObj interface{}
		if err := json.Unmarshal(pool.ExtraSpecs, &jsonObj); err == nil {
			if formatted, err := json.MarshalIndent(jsonObj, "", "  "); err == nil {
				extraSpecsDecoded = string(formatted)
			}
		} else {
			// If it fails to parse as JSON, just use as string
			extraSpecsDecoded = string(pool.ExtraSpecs)
		}
	}

	// Create custom data structure with decoded extra specs
	data := struct {
		PageData
		ExtraSpecsDecoded string
	}{
		PageData: PageData{
			Title:        "Pool " + poolID[:8] + "... - Pool Details",
			PageTitle:    "Pool Details",
			ActivePage:   "pools",
			IsDetailView: true,
			Entity:       pool,
		},
		ExtraSpecsDecoded: extraSpecsDecoded,
	}

	if err := h.templates.ExecuteTemplate(w, "base.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *WebHandler) DeletePoolHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	poolID := vars["id"]

	err := h.runner.DeletePoolByID(ctx, poolID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.PoolsAPIHandler(w, r)
}

func (h *WebHandler) PoolEditHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	poolID := vars["id"]

	pool, err := h.runner.GetPoolByID(ctx, poolID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Format extra specs JSON for editing
	extraSpecsDecoded := ""
	if len(pool.ExtraSpecs) > 0 {
		// Pretty format the JSON for editing
		var jsonObj interface{}
		if err := json.Unmarshal(pool.ExtraSpecs, &jsonObj); err == nil {
			if formatted, err := json.MarshalIndent(jsonObj, "", "  "); err == nil {
				extraSpecsDecoded = string(formatted)
			}
		} else {
			// If it fails to parse as JSON, just use as string
			extraSpecsDecoded = string(pool.ExtraSpecs)
		}
	}

	data := struct {
		Pool              interface{}
		ExtraSpecsDecoded string
	}{
		Pool:              pool,
		ExtraSpecsDecoded: extraSpecsDecoded,
	}

	if err := h.templates.ExecuteTemplate(w, "pool-edit.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *WebHandler) UpdatePoolHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	poolID := vars["id"]

	// Parse form data
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var updateParams params.UpdatePoolParams

	// Parse simple fields
	if image := r.FormValue("image"); image != "" {
		updateParams.Image = image
	}
	if flavor := r.FormValue("flavor"); flavor != "" {
		updateParams.Flavor = flavor
	}
	if runnerPrefix := r.FormValue("runner_prefix"); runnerPrefix != "" {
		updateParams.Prefix = runnerPrefix
	}
	if osType := r.FormValue("os_type"); osType != "" {
		updateParams.OSType = commonParams.OSType(osType)
	}
	if osArch := r.FormValue("os_arch"); osArch != "" {
		updateParams.OSArch = commonParams.OSArch(osArch)
	}
	if runnerGroup := r.FormValue("github-runner-group"); runnerGroup != "" {
		updateParams.GitHubRunnerGroup = &runnerGroup
	}

	// Parse numeric fields
	if minIdleStr := r.FormValue("min_idle_runners"); minIdleStr != "" {
		if minIdle, err := strconv.ParseUint(minIdleStr, 10, 32); err == nil {
			minIdleUint := uint(minIdle)
			updateParams.MinIdleRunners = &minIdleUint
		}
	}
	if maxRunnersStr := r.FormValue("max_runners"); maxRunnersStr != "" {
		if maxRunners, err := strconv.ParseUint(maxRunnersStr, 10, 32); err == nil {
			maxRunnersUint := uint(maxRunners)
			updateParams.MaxRunners = &maxRunnersUint
		}
	}
	if timeoutStr := r.FormValue("runner_bootstrap_timeout"); timeoutStr != "" {
		if timeout, err := strconv.ParseUint(timeoutStr, 10, 32); err == nil {
			timeoutUint := uint(timeout)
			updateParams.RunnerBootstrapTimeout = &timeoutUint
		}
	}
	if priorityStr := r.FormValue("priority"); priorityStr != "" {
		if priority, err := strconv.ParseUint(priorityStr, 10, 32); err == nil {
			priorityUint := uint(priority)
			updateParams.Priority = &priorityUint
		}
	}

	// Parse boolean fields
	enabled := r.FormValue("enabled") == "on"
	updateParams.Enabled = &enabled

	// Parse tags
	if tagsStr := r.FormValue("tags"); tagsStr != "" {
		tags := strings.Split(tagsStr, ",")
		for i, tag := range tags {
			tags[i] = strings.TrimSpace(tag)
		}
		updateParams.Tags = tags
	}

	// Parse extra specs JSON
	if extraSpecsStr := r.FormValue("extra_specs"); extraSpecsStr != "" {
		var extraSpecs json.RawMessage
		if err := json.Unmarshal([]byte(extraSpecsStr), &extraSpecs); err == nil {
			updateParams.ExtraSpecs = extraSpecs
		}
	}

	// Update the pool
	_, err := h.runner.UpdatePoolByID(ctx, poolID, updateParams)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return updated pools list
	h.PoolsAPIHandler(w, r)
}

func (h *WebHandler) NewPoolHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get providers
	providers, err := h.runner.ListProviders(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get repositories
	repos, err := h.runner.ListRepositories(ctx, params.RepositoryFilter{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get organizations  
	orgs, err := h.runner.ListOrganizations(ctx, params.OrganizationFilter{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get enterprises
	enterprises, err := h.runner.ListEnterprises(ctx, params.EnterpriseFilter{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := struct {
		Pool         *params.Pool
		Providers    []params.Provider
		Repositories []params.Repository
		Organizations []params.Organization
		Enterprises  []params.Enterprise
	}{
		Providers:     providers,
		Repositories:  repos,
		Organizations: orgs,
		Enterprises:   enterprises,
	}

	if err := h.templates.ExecuteTemplate(w, "pool-form.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *WebHandler) CreatePoolHandler(w http.ResponseWriter, r *http.Request) {
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
	var maxRunners, minIdleRunners, bootstrapTimeout, priority uint
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
	if p := r.FormValue("priority"); p != "" {
		if parsed, err := strconv.ParseUint(p, 10, 32); err == nil {
			priority = uint(parsed)
		}
	} else {
		priority = 100 // Default priority
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

	// Parse extra specs JSON
	var extraSpecs json.RawMessage
	if extraSpecsStr := r.FormValue("extra_specs"); extraSpecsStr != "" {
		if err := json.Unmarshal([]byte(extraSpecsStr), &extraSpecs); err != nil {
			http.Error(w, "Invalid extra specs JSON: "+err.Error(), http.StatusBadRequest)
			return
		}
	}

	createParams := params.CreatePoolParams{
		RunnerPrefix: params.RunnerPrefix{
			Prefix: r.FormValue("runner_prefix"),
		},
		ProviderName:           r.FormValue("provider_name"),
		Image:                  r.FormValue("image"),
		Flavor:                 r.FormValue("flavor"),
		MaxRunners:             maxRunners,
		MinIdleRunners:         minIdleRunners,
		RunnerBootstrapTimeout: bootstrapTimeout,
		Tags:                   tags,
		OSType:                 commonParams.OSType(r.FormValue("os_type")),
		OSArch:                 commonParams.OSArch(r.FormValue("os_arch")),
		GitHubRunnerGroup:      r.FormValue("github_runner_group"),
		Priority:               priority,
		Enabled:                r.FormValue("enabled") == "on",
		ExtraSpecs:             extraSpecs,
	}

	// Call appropriate creation function based on entity level
	var err error

	switch entityLevel {
	case "repo":
		_, err = h.runner.CreateRepoPool(ctx, entityID, createParams)
	case "org":
		_, err = h.runner.CreateOrgPool(ctx, entityID, createParams)
	case "enterprise":
		_, err = h.runner.CreateEnterprisePool(ctx, entityID, createParams)
	default:
		http.Error(w, "Invalid entity level", http.StatusBadRequest)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Close modal and refresh table
	w.Header().Set("HX-Trigger", "closeModal")
	h.PoolsAPIHandler(w, r)
}