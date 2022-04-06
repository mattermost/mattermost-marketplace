package store

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	mattermostModel "github.com/mattermost/mattermost-server/v5/model"

	"github.com/mattermost/mattermost-marketplace/internal/api"
	"github.com/mattermost/mattermost-marketplace/internal/model"
	"github.com/mattermost/mattermost-marketplace/internal/testlib"
)

func TestProxyGetPlugins(t *testing.T) {
	t.Run("empty stream", func(t *testing.T) {
		logger := testlib.MakeLogger(t)
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))
		t.Cleanup(ts.Close)

		proxyStore, err := NewProxy(ts.URL, logger)
		require.NoError(t, err)

		plugins, err := proxyStore.GetPlugins(&model.PluginFilter{
			PerPage: model.AllPerPage,
		})
		require.NoError(t, err)
		require.Empty(t, plugins)
	})

	t.Run("empty stream with error", func(t *testing.T) {
		logger := testlib.MakeLogger(t)
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, err := w.Write([]byte(`{"invalid":`))
			require.NoError(t, err)
		}))
		t.Cleanup(ts.Close)

		proxyStore, err := NewProxy(ts.URL, logger)
		require.NoError(t, err)

		plugins, err := proxyStore.GetPlugins(&model.PluginFilter{
			PerPage: model.AllPerPage,
		})
		require.Error(t, err)
		require.Empty(t, plugins)
	})

	t.Run("valid stream", func(t *testing.T) {
		logger := testlib.MakeLogger(t)
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, err := w.Write([]byte(`[{"homepage_url":"https://github.com/mattermost/mattermost-plugin-demo","icon_data":"icon-data.svg","download_url":"https://github.com/mattermost/mattermost-plugin-demo/releases/download/v0.1.0/com.mattermost.demo-plugin-0.1.0.tar.gz","signature":"signature1", "release_notes_url":"https://github.com/mattermost/mattermost-plugin-demo/releases/v0.1.0","manifest":{}}]`))
			require.NoError(t, err)
		}))
		t.Cleanup(ts.Close)

		proxyStore, err := NewProxy(ts.URL, logger)
		require.NoError(t, err)

		plugins, err := proxyStore.GetPlugins(&model.PluginFilter{
			PerPage: model.AllPerPage,
		})
		require.NoError(t, err)
		require.Equal(t, []*model.Plugin{{
			HomepageURL:     "https://github.com/mattermost/mattermost-plugin-demo",
			IconData:        "icon-data.svg",
			DownloadURL:     "https://github.com/mattermost/mattermost-plugin-demo/releases/download/v0.1.0/com.mattermost.demo-plugin-0.1.0.tar.gz",
			Signature:       "signature1",
			ReleaseNotesURL: "https://github.com/mattermost/mattermost-plugin-demo/releases/v0.1.0",
			Manifest:        &mattermostModel.Manifest{},
		}}, plugins)
	})

	t.Run("valid stream with filters set", func(t *testing.T) {
		logger := testlib.MakeLogger(t)
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			filter, err := api.ParsePluginFilter(r.URL)
			require.NoError(t, err)
			assert.Equal(t, 1, filter.Page)
			assert.Equal(t, model.AllPerPage, filter.PerPage)
			assert.Equal(t, "some filter", filter.Filter)
			assert.Equal(t, "6.0.0", filter.ServerVersion)
			assert.Equal(t, true, filter.EnterprisePlugins)
			assert.Equal(t, true, filter.Cloud)
			assert.Equal(t, "linux-amd64", filter.Platform)
			assert.Equal(t, true, filter.ReturnAllVersions)
			assert.Equal(t, "demo", filter.PluginID)

			w.WriteHeader(http.StatusOK)
			_, err = w.Write([]byte(`[{"homepage_url":"https://github.com/mattermost/mattermost-plugin-demo","icon_data":"icon-data.svg","download_url":"https://github.com/mattermost/mattermost-plugin-demo/releases/download/v0.1.0/com.mattermost.demo-plugin-0.1.0.tar.gz","signature":"signature1", "release_notes_url":"https://github.com/mattermost/mattermost-plugin-demo/releases/v0.1.0","manifest":{}}]`))
			require.NoError(t, err)
		}))
		t.Cleanup(ts.Close)

		proxyStore, err := NewProxy(ts.URL, logger)
		require.NoError(t, err)

		plugins, err := proxyStore.GetPlugins(&model.PluginFilter{
			Page:              1,
			PerPage:           model.AllPerPage,
			Filter:            "some filter",
			ServerVersion:     "6.0.0",
			EnterprisePlugins: true,
			Cloud:             true,
			Platform:          "linux-amd64",
			ReturnAllVersions: true,
			PluginID:          "demo",
		})
		require.NoError(t, err)
		require.Equal(t, []*model.Plugin{{
			HomepageURL:     "https://github.com/mattermost/mattermost-plugin-demo",
			IconData:        "icon-data.svg",
			DownloadURL:     "https://github.com/mattermost/mattermost-plugin-demo/releases/download/v0.1.0/com.mattermost.demo-plugin-0.1.0.tar.gz",
			Signature:       "signature1",
			ReleaseNotesURL: "https://github.com/mattermost/mattermost-plugin-demo/releases/v0.1.0",
			Manifest:        &mattermostModel.Manifest{},
		}}, plugins)
	})
}
