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

	pluginRouter := apiRouter.PathPrefix("/plugin/{id}").Subrouter()
	pluginRouter.Handle("", addContext(handleGetPlugin)).Methods(http.MethodGet)
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
	platform := u.Query().Get("platform")

	enterprisePlugins, err := parseBool(u, "enterprise_plugins", false)
	if err != nil {
		return nil, err
	}

	cloud, err := parseBool(u, "cloud", false)
	if err != nil {
		return nil, err
	}

	return &model.PluginFilter{
		Page:              page,
		PerPage:           perPage,
		Filter:            filter,
		ServerVersion:     serverVersion,
		EnterprisePlugins: enterprisePlugins,
		Cloud:             cloud,
		Platform:          platform,
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

// handleGetPlugins responds to GET /api/v1/plugin/{id}, returning the specified page of plugin versions.
func handleGetPlugin(c *Context, w http.ResponseWriter, r *http.Request) {
	filter, err := parsePluginFilter(r.URL)
	if err != nil {
		c.Logger.WithError(err).Error("failed to parse paging parameters")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	pluginid, ok := mux.Vars(r)["id"]
	if !ok {
		c.Logger.WithError(err).Error("failed to get pluginid from url")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Remove params that don't apply to us
	filter.Filter = ""

	plugins, err := c.Store.GetPlugin(filter, pluginid)
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
