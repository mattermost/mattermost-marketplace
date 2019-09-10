package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

var buildTag = ""
var buildHash = ""

type healthCheckResponse struct {
	Status      string
	Version     string
	ReleaseID   string
	Notes       []string
	Description string
}

func initHealthCheck(apiRouter *mux.Router, context *Context) {
	addContext := func(handler contextHandlerFunc) *contextHandler {
		return newContextHandler(context, handler)
	}

	pluginsRouter := apiRouter.PathPrefix("/health").Subrouter()
	pluginsRouter.Handle("", addContext(handleHealthCheck)).Methods("GET")
}

func handleHealthCheck(c *Context, w http.ResponseWriter, r *http.Request) {
	response := healthCheckResponse{
		Status:      "pass",
		Version:     "1",
		ReleaseID:   buildTag,
		Notes:       []string{buildHash},
		Description: "The stateless HTTP service backing the Mattermost marketplace",
	}
	w.Header().Set("Content-Type", "application/health+json")
	outputJSON(c, w, response)
}
