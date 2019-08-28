package client_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/mattermost/mattermost-marketplace/client"
	"github.com/mattermost/mattermost-marketplace/internal/model"
	mattermostModel "github.com/mattermost/mattermost-server/model"
	"github.com/stretchr/testify/require"
)

func TestPlugins(t *testing.T) {
	t.Run("no plugins", func(t *testing.T) {
		c, _, tearDown := SetupApi(t, nil)
		defer tearDown()

		plugins, err := c.GetPlugins(client.GetPluginsRequest{
			Page:    0,
			PerPage: 10,
		})
		require.NoError(t, err)
		require.Empty(t, plugins)
	})

	t.Run("parameter edge cases", func(t *testing.T) {
		t.Run("invalid page", func(t *testing.T) {
			_, url, tearDown := SetupApi(t, nil)
			defer tearDown()

			resp, err := http.Get(fmt.Sprintf("%s/api/v1/plugins?page=invalid&per_page=100", url))
			require.NoError(t, err)
			require.Equal(t, http.StatusBadRequest, resp.StatusCode)
		})

		t.Run("invalid perPage", func(t *testing.T) {
			_, url, tearDown := SetupApi(t, nil)
			defer tearDown()

			resp, err := http.Get(fmt.Sprintf("%s/api/v1/plugins?page=0&per_page=invalid", url))
			require.NoError(t, err)
			require.Equal(t, http.StatusBadRequest, resp.StatusCode)
		})

		t.Run("no paging parameters", func(t *testing.T) {
			_, url, tearDown := SetupApi(t, nil)
			defer tearDown()

			resp, err := http.Get(fmt.Sprintf("%s/api/v1/plugins", url))
			require.NoError(t, err)
			require.Equal(t, http.StatusOK, resp.StatusCode)
		})

		t.Run("missing page", func(t *testing.T) {
			_, url, tearDown := SetupApi(t, nil)
			defer tearDown()

			resp, err := http.Get(fmt.Sprintf("%s/api/v1/plugins?per_page=100", url))
			require.NoError(t, err)
			require.Equal(t, http.StatusOK, resp.StatusCode)
		})

		t.Run("missing perPage", func(t *testing.T) {
			_, url, tearDown := SetupApi(t, nil)
			defer tearDown()

			resp, err := http.Get(fmt.Sprintf("%s/api/v1/plugins?page=1", url))
			require.NoError(t, err)
			require.Equal(t, http.StatusOK, resp.StatusCode)
		})
	})

	t.Run("plugins", func(t *testing.T) {
		plugin1 := &model.Plugin{
			HomepageURL:       "https://github.com/mattermost/mattermost-plugin-demo",
			DownloadURL:       "https://github.com/mattermost/mattermost-plugin-demo/releases/download/v0.1.0/com.mattermost.demo-plugin-0.1.0.tar.gz",
			DownloadSignature: []byte("signature"),
			Manifest:          &mattermostModel.Manifest{},
		}
		plugin2 := &model.Plugin{
			HomepageURL:       "https://github.com/mattermost/mattermost-plugin-starter-template",
			DownloadURL:       "https://github.com/mattermost/mattermost-plugin-starter-template/releases/download/v0.1.0/com.mattermost.plugin-starter-template-0.1.0.tar.gz",
			DownloadSignature: []byte("signature2"),
			Manifest:          &mattermostModel.Manifest{},
		}
		plugin3 := &model.Plugin{
			HomepageURL:       "https://github.com/matterpoll/matterpoll",
			DownloadURL:       "https://github.com/matterpoll/matterpoll/releases/download/v1.1.0/com.github.matterpoll.matterpoll-1.1.0.tar.gz",
			DownloadSignature: []byte("signature3"),
			Manifest:          &mattermostModel.Manifest{},
		}
		plugins := []*model.Plugin{plugin1, plugin2, plugin3}

		t.Run("get plugins, page 0, perPage 2", func(t *testing.T) {
			c, _, tearDown := SetupApi(t, plugins)
			defer tearDown()

			plugins, err := c.GetPlugins(client.GetPluginsRequest{
				Page:    0,
				PerPage: 2,
			})
			require.NoError(t, err)
			require.Equal(t, []*model.Plugin{plugin1, plugin2}, plugins)
		})

		t.Run("get plugins, page 1, perPage 2", func(t *testing.T) {
			c, _, tearDown := SetupApi(t, plugins)
			defer tearDown()

			plugins, err := c.GetPlugins(client.GetPluginsRequest{
				Page:    1,
				PerPage: 2,
			})
			require.NoError(t, err)
			require.Equal(t, []*model.Plugin{plugin3}, plugins)
		})
	})
}
