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
		plugin1 := &model.Plugin{
			HomepageURL: "https://github.com/mattermost/mattermost-plugin-demo",
			IconData:    "icon-data.svg",
			DownloadURL: "https://github.com/mattermost/mattermost-plugin-demo/releases/download/v0.1.0/com.mattermost.demo-plugin-0.1.0.tar.gz",
			Manifest:    &mattermostModel.Manifest{},
			Signatures:  []*model.PluginSignature{{Signature: "signature1", PublicKeyHash: "hash1"}},
		}
		plugin2 := &model.Plugin{
			HomepageURL: "https://github.com/mattermost/mattermost-plugin-starter-template",
			IconData:    "icon-data2.svg",
			DownloadURL: "https://github.com/mattermost/mattermost-plugin-starter-template/releases/download/v0.1.0/com.mattermost.plugin-starter-template-0.1.0.tar.gz",
			Manifest:    &mattermostModel.Manifest{},
			Signatures:  []*model.PluginSignature{{Signature: "signature2", PublicKeyHash: "hash2"}},
		}
		plugin3 := &model.Plugin{
			HomepageURL: "https://github.com/matterpoll/matterpoll",
			IconData:    "icon-data3.svg",
			DownloadURL: "https://github.com/matterpoll/matterpoll/releases/download/v1.1.0/com.github.matterpoll.matterpoll-1.1.0.tar.gz",
			Manifest:    &mattermostModel.Manifest{},
			Signatures:  []*model.PluginSignature{{Signature: "signature3", PublicKeyHash: "hash3"}},
		}
		plugins := []*model.Plugin{plugin1, plugin2, plugin3}

		t.Run("get plugins, page 0, perPage 2", func(t *testing.T) {
			client, tearDown := setupApi(t, plugins)
			defer tearDown()

			plugins, err := client.GetPlugins(&api.GetPluginsRequest{
				Page:    0,
				PerPage: 2,
			})
			require.NoError(t, err)
			require.Equal(t, []*model.Plugin{plugin1, plugin2}, plugins)
		})

		t.Run("get plugins, page 1, perPage 2", func(t *testing.T) {
			client, tearDown := setupApi(t, plugins)
			defer tearDown()

			plugins, err := client.GetPlugins(&api.GetPluginsRequest{
				Page:    1,
				PerPage: 2,
			})
			require.NoError(t, err)
			require.Equal(t, []*model.Plugin{plugin3}, plugins)
		})
	})
}
