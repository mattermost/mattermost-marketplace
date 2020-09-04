package api

import (
	"net/http"
	"net/url"

	"github.com/gorilla/mux"

	"github.com/mattermost/mattermost-marketplace/internal/model"
)

// initPlugins registers plugin endpoints on the given router.
func initPlugins(apiRouter *mux.Router, context *Context) {
	addContext := func(handler contextHandlerFunc) *contextHandler {
		return newContextHandler(context, handler)
	}

	pluginsRouter := apiRouter.PathPrefix("/plugins").Subrouter()
	pluginsRouter.Handle("", addContext(handleGetPlugins)).Methods(http.MethodGet)
}

func parsePluginFilter(u *url.URL) (*model.PluginFilter, error) {
	page, err := parseInt(u, "page", 0)
	if err != nil {
		return nil, err
	}

	perPage, err := parseInt(u, "per_page", 100)
	if err != nil {
		return nil, err
	}

	filter := u.Query().Get("filter")
	serverVersion := u.Query().Get("server_version")
	arch := u.Query().Get("arch")

	enterprisePlugins, err := parseBool(u, "enterprise_plugins", false)
	if err != nil {
		return nil, err
	}

	return &model.PluginFilter{
		Page:              page,
		PerPage:           perPage,
		Filter:            filter,
		ServerVersion:     serverVersion,
		EnterprisePlugins: enterprisePlugins,
		Architecture:      arch,
	}, nil
}

// handleGetPlugins responds to GET /api/v1/plugins, returning the specified page of plugins.
func handleGetPlugins(c *Context, w http.ResponseWriter, r *http.Request) {
	filter, err := parsePluginFilter(r.URL)
	if err != nil {
		c.Logger.WithError(err).Error("failed to parse paging parameters")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	plugins, err := c.Store.GetPlugins(filter)
	if err != nil {
		c.Logger.WithError(err).Error("failed to query plugins")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if plugins == nil {
		plugins = []*model.Plugin{}
	}

	w.Header().Set("Content-Type", "application/json")
	outputJSON(c, w, plugins)
}
