package api

import (
	"net/http"

	"github.com/mattermost/mattermost-marketplace/internal/model"

	"github.com/gorilla/mux"
)

// initLabels registers label endpoints on the given router.
func initLabels(apiRouter *mux.Router, context *Context) {
	addContext := func(handler contextHandlerFunc) *contextHandler {
		return newContextHandler(context, handler)
	}

	pluginsRouter := apiRouter.PathPrefix("/labels").Subrouter()
	pluginsRouter.Handle("", addContext(handleGetLabels)).Methods(http.MethodGet)
}

// handleGetPlugins responds to GET /api/v1/labels, returning a list of all defined labels.
func handleGetLabels(c *Context, w http.ResponseWriter, r *http.Request) {
	response := model.AllLabels

	w.Header().Set("Content-Type", "application/json")
	outputJSON(c, w, response)
}
