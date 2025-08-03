package handlers

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/cloudbase/garm/params"
	commonParams "github.com/cloudbase/garm-provider-common/params"
)

func (h *WebHandler) InstancesHandler(w http.ResponseWriter, r *http.Request) {
	data := PageData{
		Title:      "Instances",  
		PageTitle:  "Runner Instances",
		ActivePage: "instances",
	}

	if err := h.templates.ExecuteTemplate(w, "base.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *WebHandler) InstancesAPIHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	statusFilter := r.URL.Query().Get("status")

	instances, err := h.runner.ListAllInstances(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	
	if len(instances) == 0 {
		fmt.Fprintf(w, `<tr><td colspan="6" class="px-6 py-4 text-center text-gray-500 dark:text-gray-400">No instances found</td></tr>`)
		return
	}

	for _, instance := range instances {
		// Filter by runner status if specified
		if statusFilter != "" && string(instance.RunnerStatus) != statusFilter {
			continue
		}

		// Format status
		statusClass := h.getInstanceStatusClass(instance.Status)
		runnerStatusClass := h.getRunnerStatusClass(instance.RunnerStatus)

		// Format creation time
		createdTime := "Unknown"
		if !instance.CreatedAt.IsZero() {
			createdTime = instance.CreatedAt.Format("Jan 2, 15:04")
		}

		fmt.Fprintf(w, `
			<tr class="hover:bg-gray-50 dark:hover:bg-gray-700">
				<td class="px-4 py-4 whitespace-nowrap">
					<div class="flex items-center">
						<div class="min-w-0 flex-1">
							<a href="/web/instances/%s/detail" class="text-sm font-medium text-blue-600 dark:text-blue-400 hover:text-blue-800 dark:hover:text-blue-300 hover:underline">%s</a>
							<div class="text-sm text-gray-500 dark:text-gray-400 truncate">%s</div>
						</div>
					</div>
				</td>
				<td class="px-3 py-4 whitespace-nowrap text-sm text-gray-900 dark:text-white min-w-0">
					<div class="truncate">%s</div>
				</td>
				<td class="px-3 py-4 whitespace-nowrap w-24">
					<span class="inline-flex px-2 py-1 text-xs font-semibold rounded-full %s">%s</span>
				</td>
				<td class="px-3 py-4 whitespace-nowrap w-32">
					<span class="inline-flex px-2 py-1 text-xs font-semibold rounded-full %s">%s</span>
				</td>
				<td class="px-3 py-4 whitespace-nowrap w-28 text-sm text-gray-500 dark:text-gray-400">%s</td>
				<td class="px-3 py-4 whitespace-nowrap w-20 text-right text-sm font-medium">
					<button hx-delete="/web/api/instances/%s" 
							hx-target="#instances-table" 
							hx-confirm="Are you sure you want to delete this instance?"
							class="text-red-600 hover:text-red-900 dark:text-red-400 dark:hover:text-red-300">Delete</button>
				</td>
			</tr>`,
			instance.Name, instance.Name, instance.ID, instance.PoolID,
			statusClass, string(instance.Status), runnerStatusClass, string(instance.RunnerStatus),
			createdTime, instance.Name)
	}
}

func (h *WebHandler) InstanceDetailsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	instanceName := vars["name"]

	instance, err := h.runner.GetInstance(ctx, instanceName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	data := struct {
		Instance interface{}
	}{
		Instance: instance,
	}

	if err := h.templates.ExecuteTemplate(w, "instance-details.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *WebHandler) InstanceDetailHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	instanceName := vars["name"]

	instance, err := h.runner.GetInstance(ctx, instanceName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Determine return path based on query params or default to instances list
	fromPath := r.URL.Query().Get("from")
	fromEntityID := r.URL.Query().Get("entity_id")
	
	var returnPath string
	var returnLabel string
	
	if fromPath != "" && fromEntityID != "" {
		switch fromPath {
		case "repository":
			returnPath = fmt.Sprintf("/web/repositories/%s", fromEntityID)
			returnLabel = "Repository"
		case "organization":
			returnPath = fmt.Sprintf("/web/organizations/%s", fromEntityID)
			returnLabel = "Organization"
		case "enterprise":
			returnPath = fmt.Sprintf("/web/enterprises/%s", fromEntityID)
			returnLabel = "Enterprise"
		case "pool":
			returnPath = fmt.Sprintf("/web/pools/%s/details", fromEntityID)
			returnLabel = "Pool"
		default:
			returnPath = "/web/instances"
			returnLabel = "Instances"
		}
	} else {
		returnPath = "/web/instances"
		returnLabel = "Instances"
	}

	data := PageData{
		Title:        "Instance " + instanceName + " - Instance Details",
		PageTitle:    "Instance Details: " + instanceName,
		ActivePage:   "instances",
		IsDetailView: true,
		Entity:       instance,
		ReturnPath:   returnPath,
		ReturnLabel:  returnLabel,
	}

	if err := h.templates.ExecuteTemplate(w, "base.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *WebHandler) DeleteInstanceHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	instanceName := vars["name"]

	err := h.runner.DeleteRunner(ctx, instanceName, false, false)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.InstancesAPIHandler(w, r)
}

func (h *WebHandler) getInstanceStatusClass(status commonParams.InstanceStatus) string {
	switch status {
	case "running":
		return "bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200"
	case "pending":
		return "bg-yellow-100 text-yellow-800 dark:bg-yellow-900 dark:text-yellow-200"
	case "stopped", "terminated":
		return "bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-200"
	case "stopping":
		return "bg-orange-100 text-orange-800 dark:bg-orange-900 dark:text-orange-200"
	default:
		return "bg-gray-100 text-gray-800 dark:bg-gray-800 dark:text-gray-200"
	}
}

func (h *WebHandler) getRunnerStatusClass(status params.RunnerStatus) string {
	switch status {
	case params.RunnerIdle:
		return "bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-200"
	case params.RunnerActive:
		return "bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200"
	case params.RunnerPending:
		return "bg-yellow-100 text-yellow-800 dark:bg-yellow-900 dark:text-yellow-200"
	case params.RunnerInstalling:
		return "bg-purple-100 text-purple-800 dark:bg-purple-900 dark:text-purple-200"
	case params.RunnerFailed:
		return "bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-200"
	case params.RunnerOffline:
		return "bg-gray-100 text-gray-800 dark:bg-gray-800 dark:text-gray-200"
	case params.RunnerOnline:
		return "bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200"
	case params.RunnerTerminated:
		return "bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-200"
	default:
		return "bg-gray-100 text-gray-800 dark:bg-gray-800 dark:text-gray-200"
	}
}