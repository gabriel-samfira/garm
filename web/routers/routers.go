package routers

import (
	"log/slog"
	"net/http"

	"github.com/cloudbase/garm/auth"
	"github.com/cloudbase/garm/config"
	dbCommon "github.com/cloudbase/garm/database/common"
	"github.com/cloudbase/garm/web/assets"
	"github.com/cloudbase/garm/web/handlers"
	"github.com/cloudbase/garm/web/middleware"
	"github.com/gorilla/mux"
)

func AddWebRoutes(router *mux.Router, webHandler *handlers.WebHandler, store dbCommon.Store, jwtConfig config.JWTAuth, authenticator *auth.Authenticator) {
	slog.Info("Adding web routes")

	// Asset serving (SVG logos, etc.)
	router.PathPrefix("/assets/").HandlerFunc(assets.ServeSVG).Methods("GET")

	// Authentication routes (no middleware required)
	router.HandleFunc("/web/login", func(w http.ResponseWriter, r *http.Request) {
		authHandler := handlers.NewSimpleAuthHandler()
		authHandler.SetAuthenticator(authenticator)

		if webHandler == nil {
			slog.Error("webHandler is nil")
			http.Error(w, "Web handler not initialized", http.StatusInternalServerError)
			return
		}

		if err := authHandler.LoadAuthTemplates(webHandler); err != nil {
			slog.Error("Failed to load auth templates", "error", err)
			http.Error(w, "Template loading failed", http.StatusInternalServerError)
			return
		}
		authHandler.LoginPageHandler(w, r)
	}).Methods("GET")

	router.HandleFunc("/web/init", func(w http.ResponseWriter, r *http.Request) {
		authHandler := handlers.NewSimpleAuthHandler()
		authHandler.SetAuthenticator(authenticator)
		if err := authHandler.LoadAuthTemplates(webHandler); err != nil {
			http.Error(w, "Template loading failed", http.StatusInternalServerError)
			return
		}
		authHandler.InitPageHandler(w, r)
	}).Methods("GET")

	// Create web-specific middleware
	webInitMiddleware := middleware.NewWebInitMiddleware(store)
	webAuthMiddleware := middleware.NewWebAuthMiddleware(store, jwtConfig)

	// Web UI routes (serve HTML pages)
	webRouter := router.PathPrefix("/web").Subrouter()
	webRouter.Use(webInitMiddleware.Middleware)
	webRouter.Use(webAuthMiddleware.Middleware)

	// Debug endpoint
	webRouter.HandleFunc("/debug", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("GARM Web Interface Debug - Routes are working!"))
	}).Methods("GET")

	// Dashboard
	webRouter.HandleFunc("/", webHandler.DashboardHandler).Methods("GET")
	webRouter.HandleFunc("", webHandler.DashboardHandler).Methods("GET")

	// Repositories
	webRouter.HandleFunc("/repositories", webHandler.RepositoriesHandler).Methods("GET")
	webRouter.HandleFunc("/repositories/", webHandler.RepositoriesHandler).Methods("GET")
	webRouter.HandleFunc("/repositories/new", webHandler.NewRepositoryHandler).Methods("GET")
	webRouter.HandleFunc("/repositories/{id}", webHandler.RepositoryDetailHandler).Methods("GET")
	webRouter.HandleFunc("/repositories/{id}/edit", webHandler.EditRepositoryHandler).Methods("GET")

	// Organizations
	webRouter.HandleFunc("/organizations", webHandler.OrganizationsHandler).Methods("GET")
	webRouter.HandleFunc("/organizations/", webHandler.OrganizationsHandler).Methods("GET")
	webRouter.HandleFunc("/organizations/new", webHandler.NewOrganizationHandler).Methods("GET")
	webRouter.HandleFunc("/organizations/{id}", webHandler.OrganizationDetailHandler).Methods("GET")
	webRouter.HandleFunc("/organizations/{id}/edit", webHandler.EditOrganizationHandler).Methods("GET")

	// Enterprises
	webRouter.HandleFunc("/enterprises", webHandler.EnterprisesHandler).Methods("GET")
	webRouter.HandleFunc("/enterprises/", webHandler.EnterprisesHandler).Methods("GET")
	webRouter.HandleFunc("/enterprises/new", webHandler.NewEnterpriseHandler).Methods("GET")
	webRouter.HandleFunc("/enterprises/{id}", webHandler.EnterpriseDetailHandler).Methods("GET")
	webRouter.HandleFunc("/enterprises/{id}/edit", webHandler.EditEnterpriseHandler).Methods("GET")
	webRouter.HandleFunc("/enterprises/{id}/details", webHandler.EnterpriseDetailsHandler).Methods("GET")

	// Scale Sets
	webRouter.HandleFunc("/scalesets", webHandler.ScaleSetsHandler).Methods("GET")
	webRouter.HandleFunc("/scalesets/", webHandler.ScaleSetsHandler).Methods("GET")
	webRouter.HandleFunc("/scalesets/new", webHandler.NewScaleSetHandler).Methods("GET")
	webRouter.HandleFunc("/scalesets/{id}/details", webHandler.ScaleSetDetailsHandler).Methods("GET")
	webRouter.HandleFunc("/scalesets/{id}/edit", webHandler.ScaleSetEditHandler).Methods("GET")

	// Pools
	webRouter.HandleFunc("/pools", webHandler.PoolsHandler).Methods("GET")
	webRouter.HandleFunc("/pools/", webHandler.PoolsHandler).Methods("GET")
	webRouter.HandleFunc("/pools/new", webHandler.NewPoolHandler).Methods("GET")
	webRouter.HandleFunc("/pools/{id}/details", webHandler.PoolDetailsHandler).Methods("GET")
	webRouter.HandleFunc("/pools/{id}/edit", webHandler.PoolEditHandler).Methods("GET")

	// Instances
	webRouter.HandleFunc("/instances", webHandler.InstancesHandler).Methods("GET")
	webRouter.HandleFunc("/instances/", webHandler.InstancesHandler).Methods("GET")
	webRouter.HandleFunc("/instances/{name}/details", webHandler.InstanceDetailsHandler).Methods("GET")

	// Credentials
	credentialsHandler := handlers.NewCredentialsHandler(webHandler)
	webRouter.HandleFunc("/credentials", credentialsHandler.CredentialsPageHandler).Methods("GET")
	webRouter.HandleFunc("/credentials/", credentialsHandler.CredentialsPageHandler).Methods("GET")
	webRouter.HandleFunc("/credentials/new", credentialsHandler.NewCredentialFormHandler).Methods("GET")
	webRouter.HandleFunc("/credentials/{id}/edit", credentialsHandler.EditCredentialFormHandler).Methods("GET")

	// Endpoints
	endpointsHandler := handlers.NewEndpointsHandler(webHandler)
	webRouter.HandleFunc("/endpoints", endpointsHandler.EndpointsPageHandler).Methods("GET")
	webRouter.HandleFunc("/endpoints/", endpointsHandler.EndpointsPageHandler).Methods("GET")
	webRouter.HandleFunc("/endpoints/new", endpointsHandler.NewEndpointFormHandler).Methods("GET")
	webRouter.HandleFunc("/endpoints/{name}/edit", endpointsHandler.EditEndpointFormHandler).Methods("GET")

	// Authentication API routes (no auth middleware)
	authAPIRouter := router.PathPrefix("/web/api/auth").Subrouter()
	authAPIRouter.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		authHandler := handlers.NewSimpleAuthHandler()
		authHandler.SetAuthenticator(authenticator)
		if err := authHandler.LoadAuthTemplates(webHandler); err != nil {
			http.Error(w, "Template loading failed", http.StatusInternalServerError)
			return
		}
		authHandler.LoginHandler(w, r)
	}).Methods("POST")

	authAPIRouter.HandleFunc("/init", func(w http.ResponseWriter, r *http.Request) {
		authHandler := handlers.NewSimpleAuthHandler()
		authHandler.SetAuthenticator(authenticator)
		if err := authHandler.LoadAuthTemplates(webHandler); err != nil {
			http.Error(w, "Template loading failed", http.StatusInternalServerError)
			return
		}
		authHandler.InitHandler(w, r)
	}).Methods("POST")

	authAPIRouter.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		authHandler := handlers.NewSimpleAuthHandler()
		authHandler.LogoutHandler(w, r)
	}).Methods("POST")

	// HTMX API routes (return HTML fragments)
	apiRouter := webRouter.PathPrefix("/api").Subrouter()

	// Dashboard API
	apiRouter.HandleFunc("/recent-activity", webHandler.RecentActivityHandler).Methods("GET")

	// Repository APIs
	apiRouter.HandleFunc("/repositories", webHandler.RepositoriesAPIHandler).Methods("GET")
	apiRouter.HandleFunc("/repositories", webHandler.CreateRepositoryHandler).Methods("POST")
	apiRouter.HandleFunc("/repositories/search", webHandler.RepositoriesAPIHandler).Methods("POST")
	apiRouter.HandleFunc("/repositories/{id}", webHandler.HandleMethodOverride(webHandler.UpdateRepositoryHandler)).Methods("POST", "PUT")
	apiRouter.HandleFunc("/repositories/{id}", webHandler.DeleteRepositoryHandler).Methods("DELETE")
	apiRouter.HandleFunc("/repositories/{id}/pools", webHandler.RepositoryPoolsAPIHandler).Methods("GET")
	apiRouter.HandleFunc("/repositories/{id}/instances", webHandler.RepositoryInstancesAPIHandler).Methods("GET")
	apiRouter.HandleFunc("/repositories/{id}/events", webHandler.RepositoryEventsAPIHandler).Methods("GET")

	// Organization APIs
	apiRouter.HandleFunc("/organizations", webHandler.OrganizationsAPIHandler).Methods("GET")
	apiRouter.HandleFunc("/organizations", webHandler.CreateOrganizationHandler).Methods("POST")
	apiRouter.HandleFunc("/organizations/search", webHandler.OrganizationsAPIHandler).Methods("POST")
	apiRouter.HandleFunc("/organizations/{id}", webHandler.HandleMethodOverride(webHandler.UpdateOrganizationHandler)).Methods("POST", "PUT")
	apiRouter.HandleFunc("/organizations/{id}", webHandler.DeleteOrganizationHandler).Methods("DELETE")
	apiRouter.HandleFunc("/organizations/{id}/pools", webHandler.OrganizationPoolsAPIHandler).Methods("GET")
	apiRouter.HandleFunc("/organizations/{id}/instances", webHandler.OrganizationInstancesAPIHandler).Methods("GET")
	apiRouter.HandleFunc("/organizations/{id}/events", webHandler.OrganizationEventsAPIHandler).Methods("GET")

	// Enterprise APIs
	apiRouter.HandleFunc("/enterprises", webHandler.EnterprisesAPIHandler).Methods("GET")
	apiRouter.HandleFunc("/enterprises", webHandler.HandleMethodOverride(webHandler.CreateEnterpriseHandler)).Methods("POST")
	apiRouter.HandleFunc("/enterprises/{id}", webHandler.HandleMethodOverride(webHandler.UpdateEnterpriseHandler)).Methods("POST", "PUT")
	apiRouter.HandleFunc("/enterprises/{id}", webHandler.DeleteEnterpriseHandler).Methods("DELETE")
	apiRouter.HandleFunc("/enterprises/{id}/pools", webHandler.EnterprisePoolsAPIHandler).Methods("GET")
	apiRouter.HandleFunc("/enterprises/{id}/instances", webHandler.EnterpriseInstancesAPIHandler).Methods("GET")
	apiRouter.HandleFunc("/enterprises/{id}/events", webHandler.EnterpriseEventsAPIHandler).Methods("GET")

	// Scale Set APIs
	apiRouter.HandleFunc("/scalesets", webHandler.ScaleSetsAPIHandler).Methods("GET")
	apiRouter.HandleFunc("/scalesets", webHandler.CreateScaleSetHandler).Methods("POST")
	apiRouter.HandleFunc("/scalesets/{id}", webHandler.HandleMethodOverride(webHandler.UpdateScaleSetHandler)).Methods("POST", "PUT")
	apiRouter.HandleFunc("/scalesets/{id}", webHandler.DeleteScaleSetHandler).Methods("DELETE")

	// Pool APIs
	apiRouter.HandleFunc("/pools", webHandler.PoolsAPIHandler).Methods("GET")
	apiRouter.HandleFunc("/pools", webHandler.CreatePoolHandler).Methods("POST")
	apiRouter.HandleFunc("/pools/search", webHandler.PoolsAPIHandler).Methods("POST")
	apiRouter.HandleFunc("/pools/{id}", webHandler.HandleMethodOverride(webHandler.UpdatePoolHandler)).Methods("POST", "PUT")
	apiRouter.HandleFunc("/pools/{id}", webHandler.DeletePoolHandler).Methods("DELETE")

	// Instance APIs
	apiRouter.HandleFunc("/instances", webHandler.InstancesAPIHandler).Methods("GET")
	apiRouter.HandleFunc("/instances/{name}", webHandler.DeleteInstanceHandler).Methods("DELETE")

	// Credentials APIs
	apiRouter.HandleFunc("/credentials", credentialsHandler.ListCredentialsHandler).Methods("GET")
	apiRouter.HandleFunc("/credentials", credentialsHandler.CreateCredentialHandler).Methods("POST")
	apiRouter.HandleFunc("/credentials/{id}", credentialsHandler.UpdateCredentialHandler).Methods("PUT")
	apiRouter.HandleFunc("/credentials/{id}", credentialsHandler.DeleteCredentialHandler).Methods("DELETE")

	// Endpoints APIs
	apiRouter.HandleFunc("/endpoints", endpointsHandler.ListEndpointsHandler).Methods("GET")
	apiRouter.HandleFunc("/endpoints", endpointsHandler.CreateEndpointHandler).Methods("POST")
	apiRouter.HandleFunc("/endpoints/{name}", endpointsHandler.UpdateEndpointHandler).Methods("PUT")
	apiRouter.HandleFunc("/endpoints/{name}", endpointsHandler.DeleteEndpointHandler).Methods("DELETE")

	// Static file serving (using embedded assets)
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(assets.GetStaticFS()))).Methods("GET")

	slog.Info("Web routes added successfully")
}
