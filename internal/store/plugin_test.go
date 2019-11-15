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
	demoPluginV1Min514 := &model.Plugin{
		HomepageURL: "https://github.com/mattermost/mattermost-plugin-demo",
		IconData:    "icon-data.svg",
		DownloadURL: "https://github.com/mattermost/mattermost-plugin-demo/releases/download/v0.1.0/com.mattermost.demo-plugin-0.1.0.tar.gz",
		Manifest: &mattermostModel.Manifest{
			Id:               "com.mattermost.demo-plugin",
			Name:             "Demo Plugin",
			Description:      "This plugin demonstrates the capabilities of a Mattermost plugin.",
			Version:          "0.1.0",
			MinServerVersion: "5.14.0",
		},
		Signature: "signature1",
	}

	demoPluginV2Min515 := &model.Plugin{
		HomepageURL: "https://github.com/mattermost/mattermost-plugin-demo",
		IconData:    "icon-data.svg",
		DownloadURL: "https://github.com/mattermost/mattermost-plugin-demo/releases/download/v0.2.0/com.mattermost.demo-plugin-0.2.0.tar.gz",
		Manifest: &mattermostModel.Manifest{
			Id:               "com.mattermost.demo-plugin",
			Name:             "Demo Plugin",
			Description:      "This plugin demonstrates the capabilities of a Mattermost plugin.",
			Version:          "0.2.0",
			MinServerVersion: "5.15.0",
		},
		Signature: "signature1",
	}

	starterPluginV1Min515 := &model.Plugin{
		HomepageURL: "https://github.com/mattermost/mattermost-plugin-starter-template",
		IconData:    "icon-data2.svg",
		DownloadURL: "https://github.com/mattermost/mattermost-plugin-starter-template/releases/download/v0.1.0/com.mattermost.plugin-starter-template-0.1.0.tar.gz",
		Manifest: &mattermostModel.Manifest{
			Id:               "com.mattermost.plugin-starter-template",
			Name:             "Plugin Starter Template",
			Description:      "This plugin serves as a starting point for writing a Mattermost plugin.",
			Version:          "0.1.0",
			MinServerVersion: "5.15.0",
		},
		Signature: "signature2",
	}

	data, err := json.Marshal([]*model.Plugin{
		demoPluginV1Min514,
		demoPluginV2Min515,
		starterPluginV1Min515,
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
		require.Equal(t, []*model.Plugin{demoPluginV2Min515}, actualPlugins)
	})

	t.Run("page 0, per page 10", func(t *testing.T) {
		actualPlugins, err := sqlStore.GetPlugins(&model.PluginFilter{
			Page:    0,
			PerPage: 10,
			Filter:  "",
		})
		require.NoError(t, err)
		require.Equal(t, []*model.Plugin{demoPluginV2Min515, starterPluginV1Min515}, actualPlugins)
	})

	t.Run("page 0, per page 1", func(t *testing.T) {
		actualPlugins, err := sqlStore.GetPlugins(&model.PluginFilter{
			Page:    0,
			PerPage: 1,
			Filter:  "",
		})
		require.NoError(t, err)
		require.Equal(t, []*model.Plugin{demoPluginV2Min515}, actualPlugins)
	})

	t.Run("page 0, per page 10", func(t *testing.T) {
		actualPlugins, err := sqlStore.GetPlugins(&model.PluginFilter{
			Page:    0,
			PerPage: 10,
			Filter:  "",
		})
		require.NoError(t, err)
		require.Equal(t, []*model.Plugin{demoPluginV2Min515, starterPluginV1Min515}, actualPlugins)
	})

	t.Run("default paging", func(t *testing.T) {
		actualPlugins, err := sqlStore.GetPlugins(&model.PluginFilter{PerPage: model.AllPerPage,
			Filter: "",
		})
		require.NoError(t, err)
		require.Equal(t, []*model.Plugin{demoPluginV2Min515, starterPluginV1Min515}, actualPlugins)
	})

	t.Run("filter spaces", func(t *testing.T) {
		actualPlugins, err := sqlStore.GetPlugins(&model.PluginFilter{PerPage: model.AllPerPage,
			Filter: "  ",
		})
		require.NoError(t, err)
		require.Equal(t, []*model.Plugin{demoPluginV2Min515, starterPluginV1Min515}, actualPlugins)
	})

	t.Run("id match, exact", func(t *testing.T) {
		actualPlugins, err := sqlStore.GetPlugins(&model.PluginFilter{PerPage: model.AllPerPage,
			Filter: "com.mattermost.demo-plugin",
		})
		require.NoError(t, err)
		require.Equal(t, []*model.Plugin{demoPluginV2Min515}, actualPlugins)
	})

	t.Run("id match, case-insensitive", func(t *testing.T) {
		actualPlugins, err := sqlStore.GetPlugins(&model.PluginFilter{PerPage: model.AllPerPage,
			Filter: "com.mattermost.demo-PLUGIN",
		})
		require.NoError(t, err)
		require.Equal(t, []*model.Plugin{demoPluginV2Min515}, actualPlugins)
	})

	t.Run("name match, exact", func(t *testing.T) {
		actualPlugins, err := sqlStore.GetPlugins(&model.PluginFilter{PerPage: model.AllPerPage,
			Filter: "Plugin Starter Template",
		})
		require.NoError(t, err)
		require.Equal(t, []*model.Plugin{starterPluginV1Min515}, actualPlugins)
	})

	t.Run("name match, partial", func(t *testing.T) {
		actualPlugins, err := sqlStore.GetPlugins(&model.PluginFilter{PerPage: model.AllPerPage,
			Filter: "Starter",
		})
		require.NoError(t, err)
		require.Equal(t, []*model.Plugin{starterPluginV1Min515}, actualPlugins)
	})

	t.Run("name match, case-insensitive", func(t *testing.T) {
		actualPlugins, err := sqlStore.GetPlugins(&model.PluginFilter{PerPage: model.AllPerPage,
			Filter: "TEMPLATE",
		})
		require.NoError(t, err)
		require.Equal(t, []*model.Plugin{starterPluginV1Min515}, actualPlugins)
	})

	t.Run("description match, partial", func(t *testing.T) {
		actualPlugins, err := sqlStore.GetPlugins(&model.PluginFilter{PerPage: model.AllPerPage,
			Filter: "capabilities",
		})
		require.NoError(t, err)
		require.Equal(t, []*model.Plugin{demoPluginV2Min515}, actualPlugins)
	})

	t.Run("description match, case-insensitive, multiple matches", func(t *testing.T) {
		actualPlugins, err := sqlStore.GetPlugins(&model.PluginFilter{PerPage: model.AllPerPage,
			Filter: "MATTERMOST",
		})
		require.NoError(t, err)
		require.Equal(t, []*model.Plugin{demoPluginV2Min515, starterPluginV1Min515}, actualPlugins)
	})

	t.Run("plugins that satisfy 5.15", func(t *testing.T) {
		actualPlugins, err := sqlStore.GetPlugins(&model.PluginFilter{PerPage: model.AllPerPage,
			Filter:        "MATTERMOST",
			ServerVersion: "5.15.0",
		})
		require.NoError(t, err)
		require.Equal(t, []*model.Plugin{demoPluginV2Min515, starterPluginV1Min515}, actualPlugins)
	})

	t.Run("plugins that satisfy 5.14", func(t *testing.T) {
		actualPlugins, err := sqlStore.GetPlugins(&model.PluginFilter{PerPage: model.AllPerPage,
			Filter:        "MATTERMOST",
			ServerVersion: "5.14.0",
		})
		require.NoError(t, err)
		require.Equal(t, []*model.Plugin{demoPluginV1Min514}, actualPlugins)
	})

	t.Run("with a server version that does not satisfy any plugin", func(t *testing.T) {
		actualPlugins, err := sqlStore.GetPlugins(&model.PluginFilter{PerPage: model.AllPerPage,
			ServerVersion: "5.13.0",
		})
		require.NoError(t, err)
		require.Nil(t, actualPlugins)
	})
}
