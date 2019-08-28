package client

import (
	"net/http"
	"net/url"

	"github.com/mattermost/mattermost-marketplace/internal/model"
	"github.com/pkg/errors"
)

// GetPlugins fetches the list of plugins from the configured server.
func (c *Client) GetPlugins(request GetPluginsRequest) ([]*model.Plugin, error) {
	u, err := url.Parse(c.buildURL("/api/v1/plugins"))
	if err != nil {
		return nil, err
	}

	request.ApplyToURL(u)

	resp, err := c.doGet(u.String())
	if err != nil {
		return nil, err
	}
	defer closeBody(resp)

	switch resp.StatusCode {
	case http.StatusOK:
		return model.PluginsFromReader(resp.Body)
	default:
		return nil, errors.Errorf("failed with status code %d", resp.StatusCode)
	}
}
