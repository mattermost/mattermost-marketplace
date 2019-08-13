package store

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/mattermost/mattermost-marketplace/internal/model"
	"github.com/mattermost/mattermost-marketplace/internal/testlib"
	mattermostModel "github.com/mattermost/mattermost-server/model"
	"github.com/stretchr/testify/require"
)

func TestPlugins(t *testing.T) {
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

	data, err := json.Marshal([]*model.Plugin{
		plugin1,
		plugin2,
	})
	require.NoError(t, err)

	logger := testlib.MakeLogger(t)
	sqlStore, err := New(bytes.NewReader(data), logger)
	require.NoError(t, err)

	t.Run("page 0, per page 0", func(t *testing.T) {
		actualPlugins, err := sqlStore.GetPlugins(&model.PluginFilter{
			Page:    0,
			PerPage: 0,
			Filter:  "",
		})
		require.NoError(t, err)
		require.Empty(t, actualPlugins)
	})

	t.Run("page 0, per page 1", func(t *testing.T) {
		actualPlugins, err := sqlStore.GetPlugins(&model.PluginFilter{
			Page:    0,
			PerPage: 1,
			Filter:  "",
		})
		require.NoError(t, err)
		require.Equal(t, []*model.Plugin{plugin1}, actualPlugins)
	})

	t.Run("page 0, per page 10", func(t *testing.T) {
		actualPlugins, err := sqlStore.GetPlugins(&model.PluginFilter{
			Page:    0,
			PerPage: 10,
			Filter:  "",
		})
		require.NoError(t, err)
		require.Equal(t, []*model.Plugin{plugin1, plugin2}, actualPlugins)
	})

	t.Run("page 0, per page 1", func(t *testing.T) {
		actualPlugins, err := sqlStore.GetPlugins(&model.PluginFilter{
			Page:    0,
			PerPage: 1,
			Filter:  "",
		})
		require.NoError(t, err)
		require.Equal(t, []*model.Plugin{plugin1}, actualPlugins)
	})

	t.Run("page 0, per page 10", func(t *testing.T) {
		actualPlugins, err := sqlStore.GetPlugins(&model.PluginFilter{
			Page:    0,
			PerPage: 10,
			Filter:  "",
		})
		require.NoError(t, err)
		require.Equal(t, []*model.Plugin{plugin1, plugin2}, actualPlugins)
	})

	t.Run("default paging", func(t *testing.T) {
		actualPlugins, err := sqlStore.GetPlugins(&model.PluginFilter{PerPage: model.AllPerPage,
			Filter: "",
		})
		require.NoError(t, err)
		require.Equal(t, []*model.Plugin{plugin1, plugin2}, actualPlugins)
	})
}
