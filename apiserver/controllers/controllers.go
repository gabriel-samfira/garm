// Copyright 2022 Cloudbase Solutions SRL
//
//    Licensed under the Apache License, Version 2.0 (the "License"); you may
//    not use this file except in compliance with the License. You may obtain
//    a copy of the License at
//
//         http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
//    WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
//    License for the specific language governing permissions and limitations
//    under the License.

package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"garm/apiserver/params"
	"garm/auth"
	gErrors "garm/errors"
	runnerParams "garm/params"
	"garm/runner"
	wsWriter "garm/websocket"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
)

func NewAPIController(r *runner.Runner, auth *auth.Authenticator, hub *wsWriter.Hub) (*APIController, error) {
	id, err := r.GetControllerID()
	if err != nil {
		return nil, errors.Wrap(err, "getting controller ID")
	}

	return &APIController{
		r:    r,
		auth: auth,
		hub:  hub,
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 16384,
		},
		id: id.String(),
	}, nil
}

type APIController struct {
	r        *runner.Runner
	auth     *auth.Authenticator
	hub      *wsWriter.Hub
	upgrader websocket.Upgrader
	// holds this controller's id
	id string
}

func handleError(w http.ResponseWriter, err error) {
	w.Header().Add("Content-Type", "application/json")
	origErr := errors.Cause(err)
	apiErr := params.APIErrorResponse{
		Details: origErr.Error(),
	}

	switch origErr.(type) {
	case *gErrors.NotFoundError:
		w.WriteHeader(http.StatusNotFound)
		apiErr.Error = "Not Found"
	case *gErrors.UnauthorizedError:
		w.WriteHeader(http.StatusUnauthorized)
		apiErr.Error = "Not Authorized"
		// Don't include details on 401 errors.
		apiErr.Details = ""
	case *gErrors.BadRequestError:
		w.WriteHeader(http.StatusBadRequest)
		apiErr.Error = "Bad Request"
	case *gErrors.DuplicateUserError, *gErrors.ConflictError:
		w.WriteHeader(http.StatusConflict)
		apiErr.Error = "Conflict"
	default:
		w.WriteHeader(http.StatusInternalServerError)
		apiErr.Error = "Server error"
		// Don't include details on server error.
		apiErr.Details = ""
	}

	json.NewEncoder(w).Encode(apiErr)
}

// GetControllerInfo returns means to identify this very garm instance.
// This is very useful for debugging and monitoring purposes.
func (a *APIController) GetControllerInfo() (hostname, controllerId string) {

	// the hostname is neither fixed nor in our control
	// so we get it every time to avoid confusion
	var err error
	hostname, err = os.Hostname()
	if err != nil {
		log.Printf("error getting hostname: %q", err)
		return "", ""
	}

	return hostname, a.id
}

func (a *APIController) authenticateHook(body []byte, headers http.Header) error {
	// signature := headers.Get("X-Hub-Signature-256")
	hookType := headers.Get("X-Github-Hook-Installation-Target-Type")
	var workflowJob runnerParams.WorkflowJob
	if err := json.Unmarshal(body, &workflowJob); err != nil {
		return gErrors.NewBadRequestError("invalid post body: %s", err)
	}

	switch hookType {
	case "repository":
	case "organization":
	default:
		return gErrors.NewBadRequestError("invalid hook type: %s", hookType)
	}
	return nil
}

// metric to count total webhooks received
// at this point the webhook is not yet authenticated and
// we don't know if it's meant for us or not
var webhooksReceived = prometheus.NewCounterVec(prometheus.CounterOpts{
	Name: "garm_webhooks_received",
	Help: "The total number of webhooks received",
}, []string{"valid", "reason", "hostname", "controller_id"})

func init() {
	err := prometheus.Register(webhooksReceived)
	if err != nil {
		log.Printf("error registering prometheus metric: %q", err)
	}
}

func (a *APIController) handleWorkflowJobEvent(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		handleError(w, gErrors.NewBadRequestError("invalid post body: %s", err))
		return
	}

	signature := r.Header.Get("X-Hub-Signature-256")
	hookType := r.Header.Get("X-Github-Hook-Installation-Target-Type")

	hostname, controllerId := a.GetControllerInfo()

	if err := a.r.DispatchWorkflowJob(hookType, signature, body); err != nil {
		if errors.Is(err, gErrors.ErrNotFound) {
			webhooksReceived.WithLabelValues("false", "owner_unknown", hostname, controllerId).Inc()
			log.Printf("got not found error from DispatchWorkflowJob. webhook not meant for us?: %q", err)
			return
		} else if strings.Contains(err.Error(), "signature") { // TODO: check error type
			webhooksReceived.WithLabelValues("false", "signature_invalid", hostname, controllerId).Inc()
		} else {
			webhooksReceived.WithLabelValues("false", "unknown", hostname, controllerId).Inc()
		}

		handleError(w, err)
		return
	}
	webhooksReceived.WithLabelValues("true", "", hostname, controllerId).Inc()
}

func (a *APIController) CatchAll(w http.ResponseWriter, r *http.Request) {
	headers := r.Header.Clone()

	event := runnerParams.Event(headers.Get("X-Github-Event"))
	switch event {
	case runnerParams.WorkflowJobEvent:
		a.handleWorkflowJobEvent(w, r)
	default:
		log.Printf("ignoring unknown event %s", event)
		return
	}
}

func (a *APIController) WSHandler(writer http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	if !auth.IsAdmin(ctx) {
		writer.WriteHeader(http.StatusForbidden)
		writer.Write([]byte("you need admin level access to view logs"))
		return
	}

	if a.hub == nil {
		handleError(writer, gErrors.NewBadRequestError("log streamer is disabled"))
		return
	}

	conn, err := a.upgrader.Upgrade(writer, req, nil)
	if err != nil {
		log.Printf("error upgrading to websockets: %v", err)
		return
	}

	// TODO (gsamfira): Handle ExpiresAt. Right now, if a client uses
	// a valid token to authenticate, and keeps the websocket connection
	// open, it will allow that client to stream logs via websockets
	// until the connection is broken. We need to forcefully disconnect
	// the client once the token expires.
	client, err := wsWriter.NewClient(conn, a.hub)
	if err != nil {
		log.Printf("failed to create new client: %v", err)
		return
	}
	if err := a.hub.Register(client); err != nil {
		log.Printf("failed to register new client: %v", err)
		return
	}
	client.Go()
}

// NotFoundHandler is returned when an invalid URL is acccessed
func (a *APIController) NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	apiErr := params.APIErrorResponse{
		Details: "Resource not found",
		Error:   "Not found",
	}
	w.WriteHeader(http.StatusNotFound)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(apiErr)
}

func (a *APIController) MetricsTokenHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	token, err := a.auth.GetJWTMetricsToken(ctx)
	if err != nil {
		handleError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"token": "%s"}`, token)
}

// LoginHandler returns a jwt token
func (a *APIController) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var loginInfo runnerParams.PasswordLoginParams
	if err := json.NewDecoder(r.Body).Decode(&loginInfo); err != nil {
		handleError(w, gErrors.ErrBadRequest)
		return
	}

	if err := loginInfo.Validate(); err != nil {
		handleError(w, err)
		return
	}

	ctx := r.Context()
	ctx, err := a.auth.AuthenticateUser(ctx, loginInfo)
	if err != nil {
		handleError(w, err)
		return
	}

	tokenString, err := a.auth.GetJWTToken(ctx)
	if err != nil {
		handleError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(runnerParams.JWTResponse{Token: tokenString})
}

func (a *APIController) FirstRunHandler(w http.ResponseWriter, r *http.Request) {
	if a.auth.IsInitialized() {
		err := gErrors.NewConflictError("already initialized")
		handleError(w, err)
		return
	}

	ctx := r.Context()

	var newUserParams runnerParams.NewUserParams
	if err := json.NewDecoder(r.Body).Decode(&newUserParams); err != nil {
		handleError(w, gErrors.ErrBadRequest)
		return
	}

	newUser, err := a.auth.InitController(ctx, newUserParams)
	if err != nil {
		handleError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newUser)
}

func (a *APIController) ListCredentials(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	creds, err := a.r.ListCredentials(ctx)
	if err != nil {
		handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(creds)
}

func (a *APIController) ListProviders(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	providers, err := a.r.ListProviders(ctx)
	if err != nil {
		handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(providers)
}
