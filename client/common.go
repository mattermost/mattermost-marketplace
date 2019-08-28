package client

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

// Client is the programmatic interface to the marketplace server API.
type Client struct {
	address    string
	httpClient *http.Client
}

// NewClient creates a client to the marketplace server at the given address.
func NewClient(address string) *Client {
	return &Client{
		address:    address,
		httpClient: &http.Client{},
	}
}

// closeBody ensures the Body of an http.Response is properly closed.
func closeBody(r *http.Response) {
	if r.Body != nil {
		_, _ = ioutil.ReadAll(r.Body)
		_ = r.Body.Close()
	}
}

func (c *Client) buildURL(urlPath string, args ...interface{}) string {
	return fmt.Sprintf("%s%s", c.address, fmt.Sprintf(urlPath, args...))
}

func (c *Client) doGet(u string) (*http.Response, error) {
	return c.httpClient.Get(u)
}
