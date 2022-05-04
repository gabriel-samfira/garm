package routers

import (
	"io"
	"net/http"
	"os"

	gorillaHandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	"garm/apiserver/controllers"
	"garm/auth"
)

func NewAPIRouter(han *controllers.APIController, logWriter io.Writer, authMiddleware, initMiddleware, instanceMiddleware auth.Middleware) *mux.Router {
	router := mux.NewRouter()
	log := gorillaHandlers.CombinedLoggingHandler

	// Handles github webhooks
	webhookRouter := router.PathPrefix("/webhooks").Subrouter()
	webhookRouter.PathPrefix("/").Handler(log(logWriter, http.HandlerFunc(han.CatchAll)))
	webhookRouter.PathPrefix("").Handler(log(logWriter, http.HandlerFunc(han.CatchAll)))

	// Handles API calls
	apiSubRouter := router.PathPrefix("/api/v1").Subrouter()

	// FirstRunHandler
	firstRunRouter := apiSubRouter.PathPrefix("/first-run").Subrouter()
	firstRunRouter.Handle("/", log(os.Stdout, http.HandlerFunc(han.FirstRunHandler))).Methods("POST", "OPTIONS")

	// Instance callback
	callbackRouter := apiSubRouter.PathPrefix("/callbacks").Subrouter()
	callbackRouter.Handle("/status/", log(os.Stdout, http.HandlerFunc(han.InstanceStatusMessageHandler))).Methods("POST", "OPTIONS")
	callbackRouter.Handle("/status", log(os.Stdout, http.HandlerFunc(han.InstanceStatusMessageHandler))).Methods("POST", "OPTIONS")
	callbackRouter.Use(instanceMiddleware.Middleware)
	// Login
	authRouter := apiSubRouter.PathPrefix("/auth").Subrouter()
	authRouter.Handle("/{login:login\\/?}", log(os.Stdout, http.HandlerFunc(han.LoginHandler))).Methods("POST", "OPTIONS")
	authRouter.Use(initMiddleware.Middleware)

	apiRouter := apiSubRouter.PathPrefix("").Subrouter()
	apiRouter.Use(initMiddleware.Middleware)
	apiRouter.Use(authMiddleware.Middleware)

	// Runners (instances)
	// List pool instances
	apiRouter.Handle("/pools/instances/{poolID}/", log(os.Stdout, http.HandlerFunc(han.ListPoolInstancesHandler))).Methods("GET", "OPTIONS")
	apiRouter.Handle("/pools/instances/{poolID}", log(os.Stdout, http.HandlerFunc(han.ListPoolInstancesHandler))).Methods("GET", "OPTIONS")
	// Get instance
	apiRouter.Handle("/instances/{instanceName}/", log(os.Stdout, http.HandlerFunc(han.GetInstanceHandler))).Methods("GET", "OPTIONS")
	apiRouter.Handle("/instances/{instanceName}", log(os.Stdout, http.HandlerFunc(han.GetInstanceHandler))).Methods("GET", "OPTIONS")
	// Delete instance
	// apiRouter.Handle("/instances/{instanceName}/", log(os.Stdout, http.HandlerFunc(han.CatchAll))).Methods("DELETE", "OPTIONS")
	// apiRouter.Handle("/instances/{instanceName}", log(os.Stdout, http.HandlerFunc(han.CatchAll))).Methods("DELETE", "OPTIONS")

	/////////////////////
	// Repos and pools //
	/////////////////////
	// Get pool
	apiRouter.Handle("/repositories/{repoID}/pools/{poolID}/", log(os.Stdout, http.HandlerFunc(han.GetRepoPoolHandler))).Methods("GET", "OPTIONS")
	apiRouter.Handle("/repositories/{repoID}/pools/{poolID}", log(os.Stdout, http.HandlerFunc(han.GetRepoPoolHandler))).Methods("GET", "OPTIONS")
	// Delete pool
	apiRouter.Handle("/repositories/{repoID}/pools/{poolID}/", log(os.Stdout, http.HandlerFunc(han.DeleteRepoPoolHandler))).Methods("DELETE", "OPTIONS")
	apiRouter.Handle("/repositories/{repoID}/pools/{poolID}", log(os.Stdout, http.HandlerFunc(han.DeleteRepoPoolHandler))).Methods("DELETE", "OPTIONS")
	// Update pool
	apiRouter.Handle("/repositories/{repoID}/pools/{poolID}/", log(os.Stdout, http.HandlerFunc(han.UpdateRepoPoolHandler))).Methods("PUT", "OPTIONS")
	apiRouter.Handle("/repositories/{repoID}/pools/{poolID}", log(os.Stdout, http.HandlerFunc(han.UpdateRepoPoolHandler))).Methods("PUT", "OPTIONS")
	// List pools
	apiRouter.Handle("/repositories/{repoID}/pools/", log(os.Stdout, http.HandlerFunc(han.ListRepoPoolsHandler))).Methods("GET", "OPTIONS")
	apiRouter.Handle("/repositories/{repoID}/pools", log(os.Stdout, http.HandlerFunc(han.ListRepoPoolsHandler))).Methods("GET", "OPTIONS")
	// Create pool
	apiRouter.Handle("/repositories/{repoID}/pools/", log(os.Stdout, http.HandlerFunc(han.CreateRepoPoolHandler))).Methods("POST", "OPTIONS")
	apiRouter.Handle("/repositories/{repoID}/pools", log(os.Stdout, http.HandlerFunc(han.CreateRepoPoolHandler))).Methods("POST", "OPTIONS")

	// Repo instances list
	apiRouter.Handle("/repositories/{repoID}/instances/", log(os.Stdout, http.HandlerFunc(han.ListRepoInstancesHandler))).Methods("GET", "OPTIONS")
	apiRouter.Handle("/repositories/{repoID}/instances", log(os.Stdout, http.HandlerFunc(han.ListRepoInstancesHandler))).Methods("GET", "OPTIONS")

	// Get repo
	apiRouter.Handle("/repositories/{repoID}/", log(os.Stdout, http.HandlerFunc(han.GetRepoByIDHandler))).Methods("GET", "OPTIONS")
	apiRouter.Handle("/repositories/{repoID}", log(os.Stdout, http.HandlerFunc(han.GetRepoByIDHandler))).Methods("GET", "OPTIONS")
	// Update repo
	apiRouter.Handle("/repositories/{repoID}/", log(os.Stdout, http.HandlerFunc(han.UpdateRepoHandler))).Methods("PUT", "OPTIONS")
	apiRouter.Handle("/repositories/{repoID}", log(os.Stdout, http.HandlerFunc(han.UpdateRepoHandler))).Methods("PUT", "OPTIONS")
	// Delete repo
	apiRouter.Handle("/repositories/{repoID}/", log(os.Stdout, http.HandlerFunc(han.DeleteRepoHandler))).Methods("DELETE", "OPTIONS")
	apiRouter.Handle("/repositories/{repoID}", log(os.Stdout, http.HandlerFunc(han.DeleteRepoHandler))).Methods("DELETE", "OPTIONS")
	// List repos
	apiRouter.Handle("/repositories/", log(os.Stdout, http.HandlerFunc(han.ListReposHandler))).Methods("GET", "OPTIONS")
	apiRouter.Handle("/repositories", log(os.Stdout, http.HandlerFunc(han.ListReposHandler))).Methods("GET", "OPTIONS")
	// Create repo
	apiRouter.Handle("/repositories/", log(os.Stdout, http.HandlerFunc(han.CreateRepoHandler))).Methods("POST", "OPTIONS")
	apiRouter.Handle("/repositories", log(os.Stdout, http.HandlerFunc(han.CreateRepoHandler))).Methods("POST", "OPTIONS")

	/////////////////////////////
	// Organizations and pools //
	/////////////////////////////
	// Get pool
	apiRouter.Handle("/organizations/{repoID}/pools/{poolID:poolID\\/?}", log(os.Stdout, http.HandlerFunc(han.CatchAll))).Methods("GET", "OPTIONS")
	// Delete pool
	apiRouter.Handle("/organizations/{repoID}/pools/{poolID:poolID\\/?}", log(os.Stdout, http.HandlerFunc(han.CatchAll))).Methods("DELETE", "OPTIONS")
	// List pools
	apiRouter.Handle("/organizations/{repoID}/pools/", log(os.Stdout, http.HandlerFunc(han.CatchAll))).Methods("GET", "OPTIONS")
	apiRouter.Handle("/organizations/{repoID}/pools", log(os.Stdout, http.HandlerFunc(han.CatchAll))).Methods("GET", "OPTIONS")
	// Create pool
	apiRouter.Handle("/organizations/{repoID}/pools/", log(os.Stdout, http.HandlerFunc(han.CatchAll))).Methods("POST", "OPTIONS")
	apiRouter.Handle("/organizations/{repoID}/pools", log(os.Stdout, http.HandlerFunc(han.CatchAll))).Methods("POST", "OPTIONS")

	// Get repo
	apiRouter.Handle("/organizations/{repoID:repoID\\/?}", log(os.Stdout, http.HandlerFunc(han.CatchAll))).Methods("GET", "OPTIONS")
	// Delete repo
	apiRouter.Handle("/organizations/{repoID:repoID\\/?}", log(os.Stdout, http.HandlerFunc(han.CatchAll))).Methods("DELETE", "OPTIONS")
	// List repos
	apiRouter.Handle("/organizations/", log(os.Stdout, http.HandlerFunc(han.CatchAll))).Methods("GET", "OPTIONS")
	apiRouter.Handle("/organizations", log(os.Stdout, http.HandlerFunc(han.CatchAll))).Methods("GET", "OPTIONS")
	// Create repo
	apiRouter.Handle("/organizations/", log(os.Stdout, http.HandlerFunc(han.CatchAll))).Methods("POST", "OPTIONS")
	apiRouter.Handle("/organizations", log(os.Stdout, http.HandlerFunc(han.CatchAll))).Methods("POST", "OPTIONS")

	// Credentials and providers
	apiRouter.Handle("/credentials/", log(os.Stdout, http.HandlerFunc(han.ListCredentials))).Methods("GET", "OPTIONS")
	apiRouter.Handle("/credentials", log(os.Stdout, http.HandlerFunc(han.ListCredentials))).Methods("GET", "OPTIONS")
	apiRouter.Handle("/providers/", log(os.Stdout, http.HandlerFunc(han.ListProviders))).Methods("GET", "OPTIONS")
	apiRouter.Handle("/providers", log(os.Stdout, http.HandlerFunc(han.ListProviders))).Methods("GET", "OPTIONS")

	return router
}