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
		plugin1_V1Min515 := &model.Plugin{
			HomepageURL: "https://github.com/mattermost/mattermost-plugin-demo",
			IconData:    "icon-data.svg",
			DownloadURL: "https://github.com/mattermost/mattermost-plugin-demo/releases/download/v0.1.0/com.mattermost.demo-plugin-0.1.0.tar.gz",
			Manifest: &mattermostModel.Manifest{
				Id:               "mattermost-plugin-demo",
				Name:             "mattermost-plugin-demo",
				Version:          "0.1.0",
				MinServerVersion: "5.15.0",
			},
			Signature: "signature1",
		}
		plugin1_V2Min515 := &model.Plugin{
			HomepageURL: "https://github.com/mattermost/mattermost-plugin-demo",
			IconData:    "icon-data.svg",
			DownloadURL: "https://github.com/mattermost/mattermost-plugin-demo/releases/download/v0.2.0/com.mattermost.demo-plugin-0.2.0.tar.gz",
			Manifest: &mattermostModel.Manifest{
				Id:               "mattermost-plugin-demo",
				Name:             "mattermost-plugin-demo",
				Version:          "0.2.0",
				MinServerVersion: "5.15.0",
			},
			Signature: "signature1",
		}
		plugin1_V3Min515 := &model.Plugin{
			HomepageURL: "https://github.com/mattermost/mattermost-plugin-demo",
			IconData:    "icon-data.svg",
			DownloadURL: "https://github.com/mattermost/mattermost-plugin-demo/releases/download/v0.3.0/com.mattermost.demo-plugin-0.3.0.tar.gz",
			Manifest: &mattermostModel.Manifest{
				Id:               "mattermost-plugin-demo",
				Name:             "mattermost-plugin-demo",
				Version:          "0.3.0",
				MinServerVersion: "5.15.0",
			},
			Signature: "signature1",
		}
		plugin2_V1Min516 := &model.Plugin{
			HomepageURL: "https://github.com/mattermost/mattermost-plugin-starter-template",
			IconData:    "icon-data2.svg",
			DownloadURL: "https://github.com/mattermost/mattermost-plugin-starter-template/releases/download/v0.1.0/com.mattermost.plugin-starter-template-0.1.0.tar.gz",
			Manifest: &mattermostModel.Manifest{
				Id:               "mattermost-plugin-starter-template",
				Name:             "mattermost-plugin-starter-template",
				Version:          "0.1.0",
				MinServerVersion: "5.16.0",
			},
			Signature: "signature2",
		}
		plugin3_V1NoMin := &model.Plugin{
			HomepageURL: "https://github.com/matterpoll/matterpoll",
			IconData:    "icon-data3.svg",
			DownloadURL: "https://github.com/matterpoll/matterpoll/releases/download/v1.1.0/com.github.matterpoll.matterpoll-1.1.0.tar.gz",
			Manifest: &mattermostModel.Manifest{
				Id:      "matterpoll",
				Name:    "matterpoll",
				Version: "1.1.0",
			},
			Signature: "signature3",
		}

		plugin3_V2Min516 := &model.Plugin{
			HomepageURL: "https://github.com/matterpoll/matterpoll",
			IconData:    "icon-data3.svg",
			DownloadURL: "https://github.com/matterpoll/matterpoll/releases/download/v1.2.0/com.github.matterpoll.matterpoll-1.2.0.tar.gz",
			Manifest: &mattermostModel.Manifest{
				Id:               "matterpoll",
				Name:             "matterpoll",
				Version:          "1.2.0",
				MinServerVersion: "5.16.0",
			},
			Signature: "signature3",
		}

		plugin3_V3Min517 := &model.Plugin{
			HomepageURL: "https://github.com/matterpoll/matterpoll",
			IconData:    "icon-data3.svg",
			DownloadURL: "https://github.com/matterpoll/matterpoll/releases/download/v1.3.0/com.github.matterpoll.matterpoll-1.3.0.tar.gz",
			Manifest: &mattermostModel.Manifest{
				Id:               "matterpoll",
				Name:             "matterpoll",
				Version:          "1.3.0",
				MinServerVersion: "5.17.0",
			},
			Signature: "signature3",
		}

		plugin4_V1NoMin := &model.Plugin{
			HomepageURL: "fake_plugin",
			IconData:    "icon-data3.svg",
			DownloadURL: "fake_plugin.tar.gz",
			Manifest: &mattermostModel.Manifest{
				Id:      "fake_plugin",
				Name:    "Zfake_plugin",
				Version: "1.2.4",
			},
			Signature: "signature3",
		}

		allPlugins := []*model.Plugin{
			plugin1_V1Min515,
			plugin1_V2Min515,
			plugin1_V3Min515,
			plugin2_V1Min516,
			plugin3_V1NoMin,
			plugin3_V2Min516,
			plugin3_V3Min517,
		}

		t.Run("get plugins, page 0, perPage 2", func(t *testing.T) {
			client, tearDown := setupApi(t, allPlugins)
			defer tearDown()

			plugins, err := client.GetPlugins(&api.GetPluginsRequest{
				Page:    0,
				PerPage: 2,
			})
			require.NoError(t, err)
			require.Equal(t, []*model.Plugin{plugin1_V3Min515, plugin2_V1Min516}, plugins)
		})

		t.Run("get plugins, page 1, perPage 2", func(t *testing.T) {
			client, tearDown := setupApi(t, allPlugins)
			defer tearDown()

			plugins, err := client.GetPlugins(&api.GetPluginsRequest{
				Page:    1,
				PerPage: 2,
			})
			require.NoError(t, err)
			require.Equal(t, []*model.Plugin{plugin3_V3Min517}, plugins)
		})

		t.Run("server version that satisfies all plugins", func(t *testing.T) {
			client, tearDown := setupApi(t, allPlugins)
			defer tearDown()

			plugins, err := client.GetPlugins(&api.GetPluginsRequest{
				PerPage:       3,
				ServerVersion: "5.18.0",
			})
			require.NoError(t, err)
			require.Equal(t, []*model.Plugin{plugin1_V3Min515, plugin2_V1Min516, plugin3_V3Min517}, plugins)
		})

		t.Run("server version that satisfies plugin1_V3Min515 and plugin3_V1NoMin", func(t *testing.T) {
			client, tearDown := setupApi(t, allPlugins)
			defer tearDown()

			plugins, err := client.GetPlugins(&api.GetPluginsRequest{
				PerPage:       3,
				ServerVersion: "5.15.0",
			})
			require.NoError(t, err)
			require.Equal(t, []*model.Plugin{plugin1_V3Min515, plugin3_V1NoMin}, plugins)
		})

		t.Run("server version that satisfies no plugin", func(t *testing.T) {
			client, tearDown := setupApi(t, []*model.Plugin{plugin1_V1Min515, plugin1_V2Min515, plugin1_V3Min515, plugin2_V1Min516, plugin3_V2Min516, plugin3_V3Min517})

			defer tearDown()

			plugins, err := client.GetPlugins(&api.GetPluginsRequest{
				PerPage:       3,
				ServerVersion: "5.14.0",
			})
			require.NoError(t, err)
			require.Equal(t, []*model.Plugin{}, plugins)
		})

		t.Run("no server version satisfies all plugins", func(t *testing.T) {
			client, tearDown := setupApi(t, allPlugins)
			defer tearDown()

			plugins, err := client.GetPlugins(&api.GetPluginsRequest{
				PerPage: 3,
			})
			require.NoError(t, err)
			require.Equal(t, []*model.Plugin{plugin1_V3Min515, plugin2_V1Min516, plugin3_V3Min517}, plugins)
		})

		t.Run("no server version that satisfies plugin3_V3Min517", func(t *testing.T) {
			client, tearDown := setupApi(t, allPlugins)
			defer tearDown()

			plugins, err := client.GetPlugins(&api.GetPluginsRequest{
				Filter:  "matterpoll",
				PerPage: 3,
			})
			require.NoError(t, err)
			require.Equal(t, []*model.Plugin{plugin3_V3Min517}, plugins)
		})

		t.Run("server version 1.16 that satisfies plugin3_V2Min516", func(t *testing.T) {
			client, tearDown := setupApi(t, []*model.Plugin{plugin1_V1Min515, plugin1_V2Min515, plugin1_V3Min515, plugin2_V1Min516, plugin3_V2Min516, plugin3_V3Min517, plugin4_V1NoMin})
			defer tearDown()

			plugins, err := client.GetPlugins(&api.GetPluginsRequest{
				Filter:        "matterpoll",
				ServerVersion: "5.16.0",
				PerPage:       3,
			})
			require.NoError(t, err)
			require.Equal(t, []*model.Plugin{plugin3_V2Min516}, plugins)
		})

		t.Run("server version 1.17 that satisfies plugin3_V3Min517", func(t *testing.T) {
			client, tearDown := setupApi(t, allPlugins)
			defer tearDown()

			plugins, err := client.GetPlugins(&api.GetPluginsRequest{
				Filter:        "matterpoll",
				ServerVersion: "5.17.0",
				PerPage:       3,
			})
			require.NoError(t, err)
			require.Equal(t, []*model.Plugin{plugin3_V3Min517}, plugins)
		})

		t.Run("no server version gets all the latest plugins", func(t *testing.T) {
			client, tearDown := setupApi(t, append(allPlugins, plugin4_V1NoMin))
			defer tearDown()

			plugins, err := client.GetPlugins(&api.GetPluginsRequest{
				PerPage: -1,
			})
			require.NoError(t, err)
			require.Equal(t, []*model.Plugin{plugin1_V3Min515, plugin2_V1Min516, plugin3_V3Min517, plugin4_V1NoMin}, plugins)
		})

		t.Run("invalid server_version format", func(t *testing.T) {
			client, tearDown := setupApi(t, allPlugins)
			defer tearDown()

			plugins, err := client.GetPlugins(&api.GetPluginsRequest{
				PerPage:       -1,
				ServerVersion: "1",
			})
			require.Error(t, err)
			require.Nil(t, plugins)

			plugins, err = client.GetPlugins(&api.GetPluginsRequest{
				PerPage:       -1,
				ServerVersion: "a",
			})
			require.Error(t, err)
			require.Nil(t, plugins)
		})
	})
}
