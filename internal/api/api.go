package api

import "github.com/gorilla/mux"

// Register registers the API endpoints on the given router.
func Register(rootRouter *mux.Router, context *Context) {
	apiRouter := rootRouter.PathPrefix("/api/v1").Subrouter()

	initPlugins(apiRouter, context)
	initLabels(apiRouter, context)
	initHealthCheck(apiRouter, context)
}
