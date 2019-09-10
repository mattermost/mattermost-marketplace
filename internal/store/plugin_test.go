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
	demoPlugin := &model.Plugin{
		HomepageURL:       "https://github.com/mattermost/mattermost-plugin-demo",
		IconURL:           "http://example.com/icon.svg",
		DownloadURL:       "https://github.com/mattermost/mattermost-plugin-demo/releases/download/v0.1.0/com.mattermost.demo-plugin-0.1.0.tar.gz",
		DownloadSignature: []byte("signature"),
		Manifest: &mattermostModel.Manifest{
			Id:          "com.mattermost.demo-plugin",
			Name:        "Demo Plugin",
			Description: "This plugin demonstrates the capabilities of a Mattermost plugin.",
		},
	}

	starterPlugin := &model.Plugin{
		HomepageURL:       "https://github.com/mattermost/mattermost-plugin-starter-template",
		IconURL:           "http://example.com/icon2.svg",
		DownloadURL:       "https://github.com/mattermost/mattermost-plugin-starter-template/releases/download/v0.1.0/com.mattermost.plugin-starter-template-0.1.0.tar.gz",
		DownloadSignature: []byte("signature2"),
		Manifest: &mattermostModel.Manifest{
			Id:          "com.mattermost.plugin-starter-template",
			Name:        "Plugin Starter Template",
			Description: "This plugin serves as a starting point for writing a Mattermost plugin.",
		},
	}

	data, err := json.Marshal([]*model.Plugin{
		demoPlugin,
		starterPlugin,
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
		require.Equal(t, []*model.Plugin{demoPlugin}, actualPlugins)
	})

	t.Run("page 0, per page 10", func(t *testing.T) {
		actualPlugins, err := sqlStore.GetPlugins(&model.PluginFilter{
			Page:    0,
			PerPage: 10,
			Filter:  "",
		})
		require.NoError(t, err)
		require.Equal(t, []*model.Plugin{demoPlugin, starterPlugin}, actualPlugins)
	})

	t.Run("page 0, per page 1", func(t *testing.T) {
		actualPlugins, err := sqlStore.GetPlugins(&model.PluginFilter{
			Page:    0,
			PerPage: 1,
			Filter:  "",
		})
		require.NoError(t, err)
		require.Equal(t, []*model.Plugin{demoPlugin}, actualPlugins)
	})

	t.Run("page 0, per page 10", func(t *testing.T) {
		actualPlugins, err := sqlStore.GetPlugins(&model.PluginFilter{
			Page:    0,
			PerPage: 10,
			Filter:  "",
		})
		require.NoError(t, err)
		require.Equal(t, []*model.Plugin{demoPlugin, starterPlugin}, actualPlugins)
	})

	t.Run("default paging", func(t *testing.T) {
		actualPlugins, err := sqlStore.GetPlugins(&model.PluginFilter{PerPage: model.AllPerPage,
			Filter: "",
		})
		require.NoError(t, err)
		require.Equal(t, []*model.Plugin{demoPlugin, starterPlugin}, actualPlugins)
	})

	t.Run("filter spaces", func(t *testing.T) {
		actualPlugins, err := sqlStore.GetPlugins(&model.PluginFilter{PerPage: model.AllPerPage,
			Filter: "  ",
		})
		require.NoError(t, err)
		require.Equal(t, []*model.Plugin{demoPlugin, starterPlugin}, actualPlugins)
	})

	t.Run("id match, exact", func(t *testing.T) {
		actualPlugins, err := sqlStore.GetPlugins(&model.PluginFilter{PerPage: model.AllPerPage,
			Filter: "com.mattermost.demo-plugin",
		})
		require.NoError(t, err)
		require.Equal(t, []*model.Plugin{demoPlugin}, actualPlugins)
	})

	t.Run("id match, case-insensitive", func(t *testing.T) {
		actualPlugins, err := sqlStore.GetPlugins(&model.PluginFilter{PerPage: model.AllPerPage,
			Filter: "com.mattermost.demo-PLUGIN",
		})
		require.NoError(t, err)
		require.Equal(t, []*model.Plugin{demoPlugin}, actualPlugins)
	})

	t.Run("name match, exact", func(t *testing.T) {
		actualPlugins, err := sqlStore.GetPlugins(&model.PluginFilter{PerPage: model.AllPerPage,
			Filter: "Plugin Starter Template",
		})
		require.NoError(t, err)
		require.Equal(t, []*model.Plugin{starterPlugin}, actualPlugins)
	})

	t.Run("name match, partial", func(t *testing.T) {
		actualPlugins, err := sqlStore.GetPlugins(&model.PluginFilter{PerPage: model.AllPerPage,
			Filter: "Starter",
		})
		require.NoError(t, err)
		require.Equal(t, []*model.Plugin{starterPlugin}, actualPlugins)
	})

	t.Run("name match, case-insensitive", func(t *testing.T) {
		actualPlugins, err := sqlStore.GetPlugins(&model.PluginFilter{PerPage: model.AllPerPage,
			Filter: "TEMPLATE",
		})
		require.NoError(t, err)
		require.Equal(t, []*model.Plugin{starterPlugin}, actualPlugins)
	})

	t.Run("description match, partial", func(t *testing.T) {
		actualPlugins, err := sqlStore.GetPlugins(&model.PluginFilter{PerPage: model.AllPerPage,
			Filter: "capabilities",
		})
		require.NoError(t, err)
		require.Equal(t, []*model.Plugin{demoPlugin}, actualPlugins)
	})

	t.Run("description match, case-insensitive, multiple matches", func(t *testing.T) {
		actualPlugins, err := sqlStore.GetPlugins(&model.PluginFilter{PerPage: model.AllPerPage,
			Filter: "MATTERMOST",
		})
		require.NoError(t, err)
		require.Equal(t, []*model.Plugin{demoPlugin, starterPlugin}, actualPlugins)
	})
}
