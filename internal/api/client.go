package api

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/pkg/errors"

	"github.com/mattermost/mattermost-marketplace/internal/model"
)

// Client is the programmatic interface to the Plugin Marketplace API.
type Client struct {
	Address    string
	httpClient *http.Client
}

// NewClient creates a client to the Plugin Marketplace at the given address.
func NewClient(address string) *Client {
	return &Client{
		Address:    address,
		httpClient: &http.Client{},
	}
}

// closeBody ensures the Body of an http.Response is properly closed.
func closeBody(r *http.Response) {
	if r.Body != nil {
		_, _ = io.ReadAll(r.Body)
		_ = r.Body.Close()
	}
}

func (c *Client) buildURL(urlPath string, args ...interface{}) string {
	return fmt.Sprintf("%s/%s", strings.TrimRight(c.Address, "/"), strings.TrimLeft(fmt.Sprintf(urlPath, args...), "/"))
}

func (c *Client) doGet(u string) (*http.Response, error) {
	return c.httpClient.Get(u)
}

// GetPlugins fetches the list of plugins from the configured server.
func (c *Client) GetPlugins(request *GetPluginsRequest) ([]*model.Plugin, error) {
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

// GetPlugin fetches all the specified version of a single plugin
func (c *Client) GetPlugin(request *GetPluginsRequest, pluginid string) ([]*model.Plugin, error) {
	u, err := url.Parse(c.buildURL("/api/v1/plugin/" + pluginid))
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
