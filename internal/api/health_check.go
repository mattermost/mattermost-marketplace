package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

var (
	buildTag       = ""
	buildHash      = ""
	buildHashShort = ""
)

type healthCheckResponse struct {
	Status      string                       `json:"status"`
	Version     string                       `json:"version"`
	ReleaseID   string                       `json:"releaseID"`
	Details     map[string]map[string]string `json:"details"`
	Description string                       `json:"description"`
}

func initHealthCheck(apiRouter *mux.Router, context *Context) {
	addContext := func(handler contextHandlerFunc) *contextHandler {
		return newContextHandler(context, handler)
	}

	pluginsRouter := apiRouter.PathPrefix("/health").Subrouter()
	pluginsRouter.Handle("", addContext(handleHealthCheck)).Methods("GET")
}

func handleHealthCheck(c *Context, w http.ResponseWriter, r *http.Request) {
	buildInfo := make(map[string]string)
	buildInfo["buildHash"] = buildHash
	buildInfo["buildHashShort"] = buildHashShort

	details := make(map[string]map[string]string)
	details["buildInfo"] = buildInfo

	response := healthCheckResponse{
		Status:      "pass",
		Version:     "1",
		ReleaseID:   buildTag,
		Details:     details,
		Description: "The stateless HTTP service backing the Mattermost marketplace",
	}
	w.Header().Set("Content-Type", "application/health+json")
	outputJSON(c, w, response)
}
