package v3

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/google/go-github/v28/github"
)

// Provides a GitHub API Client instance that interacts with a mocked API server
func MockGitHubClient(t *testing.T) *github.Client {
	mockAPI := httptest.NewServer(http.HandlerFunc(makeHandleAPIEndpoints(t, mockMMApiState)))
	client := github.NewClient(mockAPI.Client())
	client.BaseURL, _ = url.Parse(mockAPI.URL + "/")
	client.UploadURL, _ = url.Parse(mockAPI.URL + "/")
	return client
}
