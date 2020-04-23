package store

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/mattermost/mattermost-marketplace/internal/model"
	"github.com/mattermost/mattermost-marketplace/internal/testlib"
	mattermostModel "github.com/mattermost/mattermost-server/v5/model"
)

func TestProxy(t *testing.T) {
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

	t.Run("empty stream", func(t *testing.T) {
		logger := testlib.MakeLogger(t)
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"invalid":`))
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
			w.Write([]byte(`[{"homepage_url":"https://github.com/mattermost/mattermost-plugin-demo","icon_data":"icon-data.svg","download_url":"https://github.com/mattermost/mattermost-plugin-demo/releases/download/v0.1.0/com.mattermost.demo-plugin-0.1.0.tar.gz","signature":"signature1", "release_notes_url":"https://github.com/mattermost/mattermost-plugin-demo/releases/v0.1.0","manifest":{}}]`))
		}))
		t.Cleanup(ts.Close)

		proxyStore, err := NewProxy(ts.URL, logger)
		require.NoError(t, err)

		plugins, err := proxyStore.GetPlugins(&model.PluginFilter{
			PerPage: model.AllPerPage,
		})
		require.NoError(t, err)
		require.Equal(t, []*model.Plugin{&model.Plugin{
			HomepageURL:     "https://github.com/mattermost/mattermost-plugin-demo",
			IconData:        "icon-data.svg",
			DownloadURL:     "https://github.com/mattermost/mattermost-plugin-demo/releases/download/v0.1.0/com.mattermost.demo-plugin-0.1.0.tar.gz",
			Signature:       "signature1",
			ReleaseNotesURL: "https://github.com/mattermost/mattermost-plugin-demo/releases/v0.1.0",
			Manifest:        &mattermostModel.Manifest{},
		}}, plugins)
	})
}
