package handlers

import (
	"fmt"
	"log/slog"
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
		emptyData := SimpleEmptyTableRowData{
			ColSpan:    6,
			ItemType:   "instances",
			SearchTerm: "",
		}
		if err := h.templates.ExecuteTemplate(w, "simple-empty-table-row.html", emptyData); err != nil {
			slog.Error("Failed to execute empty table row template", "error", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
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

		rowData := InstanceRowData{
			Name:              instance.Name,
			ID:                instance.ID,
			PoolID:            instance.PoolID,
			CreatedTime:       createdTime,
			Status:            string(instance.Status),
			StatusClass:       statusClass,
			RunnerStatus:      string(instance.RunnerStatus),
			RunnerStatusClass: runnerStatusClass,
		}
		
		if err := h.templates.ExecuteTemplate(w, "instance-table-row.html", rowData); err != nil {
			slog.Error("Failed to execute instance table row template", "error", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
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
		// If this is an HTMX request (auto-refresh), redirect to instances page
		if r.Header.Get("HX-Request") == "true" {
			w.Header().Set("HX-Redirect", "/web/instances")
			w.WriteHeader(http.StatusNotFound)
			return
		}
		// For regular requests, redirect normally
		http.Redirect(w, r, "/web/instances", http.StatusFound)
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

	// For HTMX requests (auto-refresh), return just the content template
	if r.Header.Get("HX-Request") == "true" {
		if err := h.templates.ExecuteTemplate(w, "instance-detail-content", data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		// For regular requests, return the full page
		if err := h.templates.ExecuteTemplate(w, "base.html", data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
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

	// For HTMX requests, redirect to instances page
	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Redirect", "/web/instances")
		w.WriteHeader(http.StatusOK)
		return
	}
	
	// For regular requests, redirect normally
	http.Redirect(w, r, "/web/instances", http.StatusFound)
}

func (h *WebHandler) getInstanceStatusClass(status commonParams.InstanceStatus) string {
	return GetInstanceStatusClass(status)
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