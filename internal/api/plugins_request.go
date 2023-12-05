package api

import (
	"net/url"
	"strconv"
)

// GetPluginsRequest describes the parameters to request a list of plugins.
type GetPluginsRequest struct {
	Page              int
	PerPage           int
	Filter            string
	ServerVersion     string
	EnterprisePlugins bool
	Cloud             bool
	Platform          string
	ReturnAllVersions bool
	PluginID          string
}

// ApplyToURL modifies the given url to include query string parameters for the request.
func (request *GetPluginsRequest) ApplyToURL(u *url.URL) {
	q := u.Query()
	q.Add("page", strconv.Itoa(request.Page))
	q.Add("per_page", strconv.Itoa(request.PerPage))
	q.Add("filter", request.Filter)
	q.Add("server_version", request.ServerVersion)
	q.Add("enterprise_plugins", strconv.FormatBool(request.EnterprisePlugins))
	q.Add("cloud", strconv.FormatBool(request.Cloud))
	q.Add("platform", request.Platform)
	q.Add("return_all_versions", strconv.FormatBool(request.ReturnAllVersions))
	q.Add("plugin_id", request.PluginID)
	u.RawQuery = q.Encode()
}
