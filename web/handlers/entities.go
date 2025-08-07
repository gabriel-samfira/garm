package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/cloudbase/garm/params"
	"github.com/gorilla/mux"
)


// EntityConfig holds configuration for different entity types
type EntityConfig struct {
	Title           string
	PageTitle       string
	ActivePage      string
	CreateURL       string
	CreateText      string
	SearchPlaceholder string
	EntityType      string
}

// GetEntityConfigs returns configurations for all supported entity types
func GetEntityConfigs() map[string]EntityConfig {
	return map[string]EntityConfig{
		"repositories": {
			Title:             "Repositories",
			PageTitle:         "Repositories", 
			ActivePage:        "repositories",
			CreateURL:         "/web/repositories/new",
			CreateText:        "Add Repository",
			SearchPlaceholder: "Search repositories by name or owner...",
			EntityType:        "repositories",
		},
		"organizations": {
			Title:             "Organizations",
			PageTitle:         "Organizations",
			ActivePage:        "organizations", 
			CreateURL:         "/web/organizations/new",
			CreateText:        "Add Organization",
			SearchPlaceholder: "Search organizations by name...",
			EntityType:        "organizations",
		},
		"enterprises": {
			Title:             "Enterprises",
			PageTitle:         "Enterprises",
			ActivePage:        "enterprises",
			CreateURL:         "/web/enterprises/new", 
			CreateText:        "Add Enterprise",
			SearchPlaceholder: "Search enterprises by name...",
			EntityType:        "enterprises",
		},
	}
}

// EntityPageData extends PageData with entity-specific information
type EntityPageData struct {
	PageData
	EntityType        string
	SearchPlaceholder string
}

// GenericEntityListHandler handles list views for any entity type
func (h *WebHandler) GenericEntityListHandler(entityType string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		configs := GetEntityConfigs()
		config, exists := configs[entityType]
		if !exists {
			http.Error(w, "Unsupported entity type", http.StatusBadRequest)
			return
		}

		data := EntityPageData{
			PageData: PageData{
				Title:      config.Title,
				PageTitle:  config.PageTitle,
				ActivePage: config.ActivePage,
				CreateButton: &CreateButton{
					URL:  config.CreateURL,
					Text: config.CreateText,
				},
			},
			EntityType:        config.EntityType,
			SearchPlaceholder: config.SearchPlaceholder,
		}

		if err := h.templates.ExecuteTemplate(w, "base.html", data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// EntityGetter interface for fetching entities by ID
type EntityGetter interface {
	GetRepositoryByID(ctx context.Context, id string) (params.Repository, error)
	GetOrganizationByID(ctx context.Context, id string) (params.Organization, error) 
	GetEnterpriseByID(ctx context.Context, id string) (params.Enterprise, error)
}

// GenericEntityDetailHandler handles detail views for any entity type
func (h *WebHandler) GenericEntityDetailHandler(entityType string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		vars := mux.Vars(r)
		entityID := vars["id"]

		configs := GetEntityConfigs()
		config, exists := configs[entityType]
		if !exists {
			http.Error(w, "Unsupported entity type", http.StatusBadRequest)
			return
		}

		// Fetch entity based on type
		var entity interface{}
		var entityName string

		switch entityType {
		case "repositories":
			repo, repoErr := h.runner.GetRepositoryByID(ctx, entityID)
			if repoErr != nil {
				// If this is an HTMX request (auto-refresh), redirect to repositories page
				if r.Header.Get("HX-Request") == "true" {
					w.Header().Set("HX-Redirect", "/web/repositories")
					w.WriteHeader(http.StatusNotFound)
					return
				}
				// For regular requests, redirect normally
				http.Redirect(w, r, "/web/repositories", http.StatusFound)
				return
			}
			entity = repo
			entityName = repo.Name
		case "organizations":
			org, orgErr := h.runner.GetOrganizationByID(ctx, entityID)
			if orgErr != nil {
				// If this is an HTMX request (auto-refresh), redirect to organizations page
				if r.Header.Get("HX-Request") == "true" {
					w.Header().Set("HX-Redirect", "/web/organizations")
					w.WriteHeader(http.StatusNotFound)
					return
				}
				// For regular requests, redirect normally
				http.Redirect(w, r, "/web/organizations", http.StatusFound)
				return
			}
			entity = org
			entityName = org.Name
		case "enterprises":
			enterprise, entErr := h.runner.GetEnterpriseByID(ctx, entityID)
			if entErr != nil {
				// If this is an HTMX request (auto-refresh), redirect to enterprises page
				if r.Header.Get("HX-Request") == "true" {
					w.Header().Set("HX-Redirect", "/web/enterprises")
					w.WriteHeader(http.StatusNotFound)
					return
				}
				// For regular requests, redirect normally
				http.Redirect(w, r, "/web/enterprises", http.StatusFound)
				return
			}
			entity = enterprise
			entityName = enterprise.Name
		default:
			http.Error(w, "Unsupported entity type", http.StatusBadRequest)
			return
		}

		// Create title by capitalizing entity type and removing 's'
		entityTypeSingular := strings.TrimSuffix(entityType, "s")
		if len(entityTypeSingular) > 0 {
			entityTypeSingular = strings.ToUpper(entityTypeSingular[:1]) + entityTypeSingular[1:]
		}

		data := EntityPageData{
			PageData: PageData{
				Title:        fmt.Sprintf("%s - %s", entityName, entityTypeSingular),
				PageTitle:    entityName,
				ActivePage:   config.ActivePage,
				IsDetailView: true,
				Entity:       entity,
			},
			EntityType:        config.EntityType,
			SearchPlaceholder: config.SearchPlaceholder,
		}

		// For HTMX requests (auto-refresh), return just the content template
		if r.Header.Get("HX-Request") == "true" {
			if err := h.templates.ExecuteTemplate(w, "entity-detail", data); err != nil {
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
}

// Repository handlers using the unified approach
func (h *WebHandler) RepositoriesHandler(w http.ResponseWriter, r *http.Request) {
	h.GenericEntityListHandler("repositories")(w, r)
}

func (h *WebHandler) RepositoryDetailHandler(w http.ResponseWriter, r *http.Request) {
	h.GenericEntityDetailHandler("repositories")(w, r)
}

// Organization handlers using the unified approach  
func (h *WebHandler) OrganizationsHandler(w http.ResponseWriter, r *http.Request) {
	h.GenericEntityListHandler("organizations")(w, r)
}

func (h *WebHandler) OrganizationDetailHandler(w http.ResponseWriter, r *http.Request) {
	h.GenericEntityDetailHandler("organizations")(w, r)
}

// Enterprise handlers using the unified approach
func (h *WebHandler) EnterprisesHandler(w http.ResponseWriter, r *http.Request) {
	h.GenericEntityListHandler("enterprises")(w, r)
}

func (h *WebHandler) EnterpriseDetailHandler(w http.ResponseWriter, r *http.Request) {
	h.GenericEntityDetailHandler("enterprises")(w, r)
}

// EntityFormData holds data for the unified entity form
type EntityFormData struct {
	Entity      interface{}
	Credentials []params.ForgeCredentials
	ForgeType   string
	EntityType  string
}

// GenericEntityNewFormHandler handles new entity form for any entity type
func (h *WebHandler) GenericEntityNewFormHandler(entityType string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		
		// For enterprises, we don't need forge type selection since they're only available for GitHub
		if entityType == "enterprises" {
			credentials, err := h.runner.ListCredentials(ctx)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			data := EntityFormData{
				Credentials: credentials,
				EntityType:  entityType,
			}
			if err := h.templates.ExecuteTemplate(w, "entity-form", data); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			return
		}

		// For repositories and organizations, we need forge type selection
		forgeType := r.URL.Query().Get("forge_type")
		if forgeType == "" {
			// Show forge selector template based on entity type
			templateName := entityType + "-forge-selector.html"
			switch entityType {
			case "repositories":
				templateName = "repository-forge-selector.html"
			case "organizations":
				templateName = "organization-forge-selector.html"
			}
			
			if err := h.templates.ExecuteTemplate(w, templateName, nil); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			return
		}

		var credentials []params.ForgeCredentials
		var err error

		switch params.EndpointType(forgeType) {
		case params.GithubEndpointType:
			credentials, err = h.runner.ListCredentials(ctx)
		case params.GiteaEndpointType:
			credentials, err = h.runner.ListGiteaCredentials(ctx)
		default:
			http.Error(w, "Invalid forge type", http.StatusBadRequest)
			return
		}

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		data := EntityFormData{
			Credentials: credentials,
			ForgeType:   forgeType,
			EntityType:  entityType,
		}
		if err := h.templates.ExecuteTemplate(w, "entity-form", data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// GenericEntityEditFormHandler handles edit entity form for any entity type
func (h *WebHandler) GenericEntityEditFormHandler(entityType string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		vars := mux.Vars(r)
		entityID := vars["id"]

		// Fetch entity based on type
		var entity interface{}
		var forgeType string
		var err error

		switch entityType {
		case "repositories":
			repo, repoErr := h.runner.GetRepositoryByID(ctx, entityID)
			if repoErr != nil {
				http.Error(w, repoErr.Error(), http.StatusNotFound)
				return
			}
			entity = repo
			forgeType = string(repo.Endpoint.EndpointType)
		case "organizations":
			org, orgErr := h.runner.GetOrganizationByID(ctx, entityID)
			if orgErr != nil {
				http.Error(w, orgErr.Error(), http.StatusNotFound)
				return
			}
			entity = org
			forgeType = string(org.Endpoint.EndpointType)
		case "enterprises":
			enterprise, entErr := h.runner.GetEnterpriseByID(ctx, entityID)
			if entErr != nil {
				http.Error(w, entErr.Error(), http.StatusNotFound)
				return
			}
			entity = enterprise
			forgeType = string(params.GithubEndpointType) // Enterprises are always GitHub
		default:
			http.Error(w, "Unsupported entity type", http.StatusBadRequest)
			return
		}

		// Get credentials based on forge type
		var credentials []params.ForgeCredentials
		switch params.EndpointType(forgeType) {
		case params.GithubEndpointType:
			credentials, err = h.runner.ListCredentials(ctx)
		case params.GiteaEndpointType:
			credentials, err = h.runner.ListGiteaCredentials(ctx)
		default:
			http.Error(w, "Invalid forge type", http.StatusInternalServerError)
			return
		}

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		data := EntityFormData{
			Entity:      entity,
			Credentials: credentials,
			ForgeType:   forgeType,
			EntityType:  entityType,
		}
		if err := h.templates.ExecuteTemplate(w, "entity-form", data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// Repository form handlers using the unified approach
func (h *WebHandler) RepositoryNewFormHandler(w http.ResponseWriter, r *http.Request) {
	h.GenericEntityNewFormHandler("repositories")(w, r)
}

func (h *WebHandler) RepositoryEditFormHandler(w http.ResponseWriter, r *http.Request) {
	h.GenericEntityEditFormHandler("repositories")(w, r)
}

// Organization form handlers using the unified approach
func (h *WebHandler) OrganizationNewFormHandler(w http.ResponseWriter, r *http.Request) {
	h.GenericEntityNewFormHandler("organizations")(w, r)
}

func (h *WebHandler) OrganizationEditFormHandler(w http.ResponseWriter, r *http.Request) {
	h.GenericEntityEditFormHandler("organizations")(w, r)
}

// Enterprise form handlers using the unified approach
func (h *WebHandler) EnterpriseNewFormHandler(w http.ResponseWriter, r *http.Request) {
	h.GenericEntityNewFormHandler("enterprises")(w, r)
}

func (h *WebHandler) EnterpriseEditFormHandler(w http.ResponseWriter, r *http.Request) {
	h.GenericEntityEditFormHandler("enterprises")(w, r)
}