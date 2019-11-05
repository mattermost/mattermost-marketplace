package api_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/mattermost/mattermost-marketplace/internal/api"
	"github.com/mattermost/mattermost-marketplace/internal/model"
	"github.com/mattermost/mattermost-marketplace/internal/store"
	"github.com/mattermost/mattermost-marketplace/internal/testlib"
	mattermostModel "github.com/mattermost/mattermost-server/model"
	"github.com/stretchr/testify/require"
)

func setupApi(t *testing.T, plugins []*model.Plugin) (*api.Client, func()) {
	logger := testlib.MakeLogger(t)

	data, err := json.Marshal(plugins)
	require.NoError(t, err)
	store, err := store.New(bytes.NewReader(data), logger)
	require.NoError(t, err)

	router := mux.NewRouter()
	api.Register(router, &api.Context{
		Store:  store,
		Logger: logger,
	})
	ts := httptest.NewServer(router)

	return api.NewClient(ts.URL), func() {
		ts.Close()
	}
}

func TestPlugins(t *testing.T) {
	t.Run("no plugins", func(t *testing.T) {
		client, tearDown := setupApi(t, nil)
		defer tearDown()

		plugins, err := client.GetPlugins(&api.GetPluginsRequest{
			Page:    0,
			PerPage: 10,
		})
		require.NoError(t, err)
		require.Empty(t, plugins)
	})

	t.Run("parameter edge cases", func(t *testing.T) {
		t.Run("invalid page", func(t *testing.T) {
			client, tearDown := setupApi(t, nil)
			defer tearDown()

			resp, err := http.Get(fmt.Sprintf("%s/api/v1/plugins?page=invalid&per_page=100", client.Address))
			require.NoError(t, err)
			require.Equal(t, http.StatusBadRequest, resp.StatusCode)
		})

		t.Run("invalid perPage", func(t *testing.T) {
			client, tearDown := setupApi(t, nil)
			defer tearDown()

			resp, err := http.Get(fmt.Sprintf("%s/api/v1/plugins?page=0&per_page=invalid", client.Address))
			require.NoError(t, err)
			require.Equal(t, http.StatusBadRequest, resp.StatusCode)
		})

		t.Run("no paging parameters", func(t *testing.T) {
			client, tearDown := setupApi(t, nil)
			defer tearDown()

			resp, err := http.Get(fmt.Sprintf("%s/api/v1/plugins", client.Address))
			require.NoError(t, err)
			require.Equal(t, http.StatusOK, resp.StatusCode)
		})

		t.Run("missing page", func(t *testing.T) {
			client, tearDown := setupApi(t, nil)
			defer tearDown()

			resp, err := http.Get(fmt.Sprintf("%s/api/v1/plugins?per_page=100", client.Address))
			require.NoError(t, err)
			require.Equal(t, http.StatusOK, resp.StatusCode)
		})

		t.Run("missing perPage", func(t *testing.T) {
			client, tearDown := setupApi(t, nil)
			defer tearDown()

			resp, err := http.Get(fmt.Sprintf("%s/api/v1/plugins?page=1", client.Address))
			require.NoError(t, err)
			require.Equal(t, http.StatusOK, resp.StatusCode)
		})
	})

	t.Run("plugins", func(t *testing.T) {
		plugin1_1 := &model.Plugin{
			HomepageURL: "https://github.com/mattermost/mattermost-plugin-demo",
			IconData:    "icon-data.svg",
			DownloadURL: "https://github.com/mattermost/mattermost-plugin-demo/releases/download/v0.1.0/com.mattermost.demo-plugin-0.1.0.tar.gz",
			Manifest:    &mattermostModel.Manifest{Id: "mattermost-plugin-demo", Name: "mattermost-plugin-demo", Version: "1.2.3", MinServerVersion: "1.15.0"},
			Signatures:  []*model.PluginSignature{{Signature: "signature1", PublicKeyHash: "hash1"}},
		}
		plugin1_2 := &model.Plugin{
			HomepageURL: "https://github.com/mattermost/mattermost-plugin-demo",
			IconData:    "icon-data.svg",
			DownloadURL: "https://github.com/mattermost/mattermost-plugin-demo/releases/download/v0.1.0/com.mattermost.demo-plugin-0.1.0.tar.gz",
			Manifest:    &mattermostModel.Manifest{Id: "mattermost-plugin-demo", Name: "mattermost-plugin-demo", Version: "1.2.4", MinServerVersion: "1.15.0"},
			Signatures:  []*model.PluginSignature{{Signature: "signature1", PublicKeyHash: "hash1"}},
		}
		plugin1_3 := &model.Plugin{
			HomepageURL: "https://github.com/mattermost/mattermost-plugin-demo",
			IconData:    "icon-data.svg",
			DownloadURL: "https://github.com/mattermost/mattermost-plugin-demo/releases/download/v0.1.0/com.mattermost.demo-plugin-0.1.0.tar.gz",
			Manifest:    &mattermostModel.Manifest{Id: "mattermost-plugin-demo", Name: "mattermost-plugin-demo", Version: "1.2.5", MinServerVersion: "1.15.0"},
			Signatures:  []*model.PluginSignature{{Signature: "signature1", PublicKeyHash: "hash1"}},
		}
		plugin2_1 := &model.Plugin{
			HomepageURL: "https://github.com/mattermost/mattermost-plugin-starter-template",
			IconData:    "icon-data2.svg",
			DownloadURL: "https://github.com/mattermost/mattermost-plugin-starter-template/releases/download/v0.1.0/com.mattermost.plugin-starter-template-0.1.0.tar.gz",
			Manifest:    &mattermostModel.Manifest{Id: "mattermost-plugin-starter-template", Name: "mattermost-plugin-starter-template", Version: "1.2.3", MinServerVersion: "1.16.0"},
			Signatures:  []*model.PluginSignature{{Signature: "signature2", PublicKeyHash: "hash2"}},
		}
		plugin3_1 := &model.Plugin{
			HomepageURL: "https://github.com/matterpoll/matterpoll",
			IconData:    "icon-data3.svg",
			DownloadURL: "https://github.com/matterpoll/matterpoll/releases/download/v1.1.0/com.github.matterpoll.matterpoll-1.1.0.tar.gz",
			Manifest:    &mattermostModel.Manifest{Id: "matterpoll", Name: "matterpoll", Version: "1.2.3", MinServerVersion: "1.16.0"},
			Signatures:  []*model.PluginSignature{{Signature: "signature3", PublicKeyHash: "hash3"}},
		}

		plugin3_2 := &model.Plugin{
			HomepageURL: "https://github.com/matterpoll/matterpoll",
			IconData:    "icon-data3.svg",
			DownloadURL: "https://github.com/matterpoll/matterpoll/releases/download/v1.1.0/com.github.matterpoll.matterpoll-1.1.0.tar.gz",
			Manifest:    &mattermostModel.Manifest{Id: "matterpoll", Name: "matterpoll", Version: "1.2.4", MinServerVersion: "1.17.0"},
			Signatures:  []*model.PluginSignature{{Signature: "signature3", PublicKeyHash: "hash3"}},
		}

		plugin4_1 := &model.Plugin{
			HomepageURL: "fake_plugin",
			IconData:    "icon-data3.svg",
			DownloadURL: "fake_plugin.tar.gz",
			Manifest:    &mattermostModel.Manifest{Id: "fake_plugin", Name: "Zfake_plugin", Version: "1.2.4"},
			Signatures:  []*model.PluginSignature{{Signature: "signature3", PublicKeyHash: "hash3"}},
		}

		plugins := []*model.Plugin{plugin1_1, plugin1_2, plugin1_3, plugin2_1, plugin3_1, plugin3_2}

		t.Run("get plugins, page 0, perPage 2", func(t *testing.T) {
			client, tearDown := setupApi(t, plugins)
			defer tearDown()

			plugins, err := client.GetPlugins(&api.GetPluginsRequest{
				Page:    0,
				PerPage: 2,
			})
			require.NoError(t, err)
			require.Equal(t, []*model.Plugin{plugin1_3, plugin2_1}, plugins)
		})

		t.Run("get plugins, page 1, perPage 2", func(t *testing.T) {
			client, tearDown := setupApi(t, plugins)
			defer tearDown()

			plugins, err := client.GetPlugins(&api.GetPluginsRequest{
				Page:    1,
				PerPage: 2,
			})
			require.NoError(t, err)
			require.Equal(t, []*model.Plugin{plugin3_2}, plugins)
		})

		t.Run("server version that satisfies all plugins", func(t *testing.T) {
			client, tearDown := setupApi(t, plugins)
			defer tearDown()

			plugins, err := client.GetPlugins(&api.GetPluginsRequest{
				PerPage:       3,
				ServerVersion: "1.18.0",
			})
			require.NoError(t, err)
			require.Equal(t, []*model.Plugin{plugin1_3, plugin2_1, plugin3_2}, plugins)
		})

		t.Run("server version that satisfies 1 plugin", func(t *testing.T) {
			client, tearDown := setupApi(t, plugins)
			defer tearDown()

			plugins, err := client.GetPlugins(&api.GetPluginsRequest{
				PerPage:       3,
				ServerVersion: "1.15.0",
			})
			require.NoError(t, err)
			require.Equal(t, []*model.Plugin{plugin1_3}, plugins)
		})

		t.Run("server version that satisfies no plugin", func(t *testing.T) {
			client, tearDown := setupApi(t, plugins)
			defer tearDown()

			plugins, err := client.GetPlugins(&api.GetPluginsRequest{
				PerPage:       3,
				ServerVersion: "1.14.0",
			})
			require.NoError(t, err)
			require.Equal(t, []*model.Plugin{}, plugins)
		})

		t.Run("no server version satisfies all plugins", func(t *testing.T) {
			client, tearDown := setupApi(t, plugins)
			defer tearDown()

			plugins, err := client.GetPlugins(&api.GetPluginsRequest{
				PerPage: 3,
			})
			require.NoError(t, err)
			require.Equal(t, []*model.Plugin{plugin1_3, plugin2_1, plugin3_2}, plugins)
		})

		t.Run("no server version that satisfies 1 plugin", func(t *testing.T) {
			client, tearDown := setupApi(t, plugins)
			defer tearDown()

			plugins, err := client.GetPlugins(&api.GetPluginsRequest{
				Filter:  "matterpoll",
				PerPage: 3,
			})
			require.NoError(t, err)
			require.Equal(t, []*model.Plugin{plugin3_2}, plugins)
		})

		t.Run("server version 1.16 that satisfies 1 plugin", func(t *testing.T) {
			client, tearDown := setupApi(t, plugins)
			defer tearDown()

			plugins, err := client.GetPlugins(&api.GetPluginsRequest{
				Filter:        "matterpoll",
				ServerVersion: "1.16.0",
				PerPage:       3,
			})
			require.NoError(t, err)
			require.Equal(t, []*model.Plugin{plugin3_1}, plugins)
		})

		t.Run("server version 1.17 that satisfies 1 plugin", func(t *testing.T) {
			client, tearDown := setupApi(t, plugins)
			defer tearDown()

			plugins, err := client.GetPlugins(&api.GetPluginsRequest{
				Filter:        "matterpoll",
				ServerVersion: "1.17.0",
				PerPage:       3,
			})
			require.NoError(t, err)
			require.Equal(t, []*model.Plugin{plugin3_2}, plugins)
		})

		t.Run("server version that satisfies 1 plugin with no min_server_version", func(t *testing.T) {
			client, tearDown := setupApi(t, []*model.Plugin{plugin1_1, plugin1_2, plugin1_3, plugin2_1, plugin3_1, plugin3_2, plugin4_1})
			defer tearDown()

			plugins, err := client.GetPlugins(&api.GetPluginsRequest{
				PerPage: -1,
			})
			require.NoError(t, err)
			require.Equal(t, []*model.Plugin{plugin1_3, plugin2_1, plugin3_2, plugin4_1}, plugins)
		})
	})
}
