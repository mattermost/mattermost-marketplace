package main

import (
	"context"
	"strings"
	"testing"

	mocks "github.com/mattermost/mattermost-marketplace/mocks/github/v3"
)

func TestGenerator(t *testing.T) {
	client := mocks.MockGitHubClient(t)
	ctx := context.Background()

	t.Run("Check the \"mocked\" GitHub client", func(t *testing.T) {
		t.Run("Assurance that the client is offline", func(t *testing.T) {
			onlineAPIServer := strings.ToLower("api.github.com")
			localhostDigitsServer := strings.ToLower("127.0.0.1")
			localhostNameServer := strings.ToLower("localhost")
			actual := strings.Split(strings.ToLower(client.BaseURL.Host), ":")[0]

			if strings.Compare(onlineAPIServer, actual) == 0 {
				t.Errorf("Client is still expecting to connect to main GitHub API")
			}
			if (len(actual) > 0) &&
				(strings.Compare(localhostDigitsServer, actual) != 0) &&
				(strings.Compare(localhostNameServer, actual) != 0) {
				// client does appear to be separate from the main GitHub API,
				// so add that to the log in case something breaks
				t.Log("GitHub Client interacting with API server at: " + client.BaseURL.String())
			}
		})
		t.Run("Get a Repository without erroring", func(t *testing.T) {
			_, _, err := client.Repositories.Get(ctx, "mattermost", "mattermost-plugin-github")
			if err != nil {
				t.Errorf("client.Repositories.Get() errored: " + err.Error())
			}
		})
		t.Run("List a Repository's Releases without erroring", func(t *testing.T) {
			_, _, err := client.Repositories.ListReleases(ctx, "mattermost", "mattermost-plugin-github", nil)
			if err != nil {
				t.Errorf("client.Repositories.ListReleases() errored: " + err.Error())
			}
		})
	})
}
