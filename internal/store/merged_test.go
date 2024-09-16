package store

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	mattermostModel "github.com/mattermost/mattermost-server/v6/model"

	"github.com/mattermost/mattermost-marketplace/internal/model"
	"github.com/mattermost/mattermost-marketplace/internal/testlib"
)

func TestMerged(t *testing.T) {
	plugin1V1 := &model.Plugin{
		HomepageURL: "https://github.com/mattermost/mattermost-plugin-demo",
		IconData:    "icon-data.svg",
		DownloadURL: "https://github.com/mattermost/mattermost-plugin-demo/releases/download/v0.1.0/com.mattermost.demo-plugin-0.1.0.tar.gz",
		Manifest: &mattermostModel.Manifest{
			Id:      "mattermost-plugin-demo",
			Name:    "mattermost-plugin-demo",
			Version: "0.1.0",
		},
		Signature: "signature1",
	}
	plugin1V2 := &model.Plugin{
		HomepageURL: "https://github.com/mattermost/mattermost-plugin-demo",
		IconData:    "icon-data.svg",
		DownloadURL: "https://github.com/mattermost/mattermost-plugin-demo/releases/download/v0.2.0/com.mattermost.demo-plugin-0.2.0.tar.gz",
		Manifest: &mattermostModel.Manifest{
			Id:      "mattermost-plugin-demo",
			Name:    "mattermost-plugin-demo",
			Version: "0.2.0",
		},
		Signature: "signature1",
	}
	plugin1V3 := &model.Plugin{
		HomepageURL: "https://github.com/mattermost/mattermost-plugin-demo",
		IconData:    "icon-data.svg",
		DownloadURL: "https://github.com/mattermost/mattermost-plugin-demo/releases/download/v0.3.0/com.mattermost.demo-plugin-0.3.0.tar.gz",
		Manifest: &mattermostModel.Manifest{
			Id:      "mattermost-plugin-demo",
			Name:    "mattermost-plugin-demo",
			Version: "0.3.0",
		},
		Signature: "signature1",
	}
	plugin2V1 := &model.Plugin{
		HomepageURL: "https://github.com/mattermost/mattermost-plugin-starter-template",
		IconData:    "icon-data2.svg",
		DownloadURL: "https://github.com/mattermost/mattermost-plugin-starter-template/releases/download/v0.1.0/com.mattermost.plugin-starter-template-0.1.0.tar.gz",
		Manifest: &mattermostModel.Manifest{
			Id:      "mattermost-plugin-starter-template",
			Name:    "mattermost-plugin-starter-template",
			Version: "0.1.0",
		},
		Signature: "signature2",
	}
	plugin3V1 := &model.Plugin{
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

	plugin3V2 := &model.Plugin{
		HomepageURL: "https://github.com/matterpoll/matterpoll",
		IconData:    "icon-data3.svg",
		DownloadURL: "https://github.com/matterpoll/matterpoll/releases/download/v1.2.0/com.github.matterpoll.matterpoll-1.2.0.tar.gz",
		Manifest: &mattermostModel.Manifest{
			Id:      "matterpoll",
			Name:    "matterpoll",
			Version: "1.2.0",
		},
		Signature: "signature3",
	}

	plugin3V3 := &model.Plugin{
		HomepageURL: "https://github.com/matterpoll/matterpoll",
		IconData:    "icon-data3.svg",
		DownloadURL: "https://github.com/matterpoll/matterpoll/releases/download/v1.3.0/com.github.matterpoll.matterpoll-1.3.0.tar.gz",
		Manifest: &mattermostModel.Manifest{
			Id:      "matterpoll",
			Name:    "matterpoll",
			Version: "1.3.0",
		},
		Signature: "signature3",
	}

	plugin4V1 := &model.Plugin{
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

	plugin4V1Later := &model.Plugin{
		HomepageURL: "fake_plugin",
		IconData:    "icon-data3-later.svg",
		DownloadURL: "fake_plugin.tar.gz",
		Manifest: &mattermostModel.Manifest{
			Id:      "fake_plugin",
			Name:    "Zfake_plugin",
			Version: "1.2.4",
		},
		Signature: "signature3",
	}

	t.Run("no stores", func(t *testing.T) {
		logger := testlib.MakeLogger(t)

		store := NewMerged(logger)
		require.NotNil(t, store)

		plugins, err := store.GetPlugins(&model.PluginFilter{
			Page:    0,
			PerPage: model.AllPerPage,
		})
		require.NoError(t, err)
		assert.Empty(t, plugins)
	})

	t.Run("single empty store", func(t *testing.T) {
		logger := testlib.MakeLogger(t)

		static1, err := NewStatic([]*model.Plugin{}, logger)
		require.NoError(t, err)

		store := NewMerged(logger, static1)
		assert.NoError(t, err)
		require.NotNil(t, store)

		plugins, err := store.GetPlugins(&model.PluginFilter{
			Page:    0,
			PerPage: model.AllPerPage,
		})
		require.NoError(t, err)
		assert.Empty(t, plugins)
	})

	t.Run("multiple empty stores", func(t *testing.T) {
		logger := testlib.MakeLogger(t)

		static1, err := NewStatic([]*model.Plugin{}, logger)
		require.NoError(t, err)
		static2, err := NewStatic([]*model.Plugin{}, logger)
		require.NoError(t, err)

		store := NewMerged(logger, static1, static2)
		assert.NoError(t, err)
		require.NotNil(t, store)

		plugins, err := store.GetPlugins(&model.PluginFilter{
			Page:    0,
			PerPage: model.AllPerPage,
		})
		require.NoError(t, err)
		assert.Empty(t, plugins)
	})

	t.Run("single, populated store", func(t *testing.T) {
		logger := testlib.MakeLogger(t)

		static1, err := NewStatic([]*model.Plugin{plugin1V3, plugin2V1, plugin3V3, plugin4V1}, logger)
		require.NoError(t, err)

		store := NewMerged(logger, static1)
		assert.NoError(t, err)
		require.NotNil(t, store)

		plugins, err := store.GetPlugins(&model.PluginFilter{
			Page:    0,
			PerPage: model.AllPerPage,
		})
		require.NoError(t, err)
		assert.Equal(t, []*model.Plugin{
			plugin1V3,
			plugin2V1,
			plugin3V3,
			plugin4V1,
		}, plugins)
	})

	t.Run("conflict-free merge", func(t *testing.T) {
		logger := testlib.MakeLogger(t)

		static1, err := NewStatic([]*model.Plugin{plugin1V1, plugin1V2, plugin1V3}, logger)
		require.NoError(t, err)
		static2, err := NewStatic([]*model.Plugin{plugin2V1, plugin3V1, plugin3V2, plugin3V3, plugin4V1}, logger)
		require.NoError(t, err)

		store := NewMerged(logger, static1, static2)
		assert.NoError(t, err)
		require.NotNil(t, store)

		plugins, err := store.GetPlugins(&model.PluginFilter{
			Page:    0,
			PerPage: model.AllPerPage,
		})
		require.NoError(t, err)
		assert.Equal(t, []*model.Plugin{
			plugin1V3,
			plugin2V1,
			plugin3V3,
			plugin4V1,
		}, plugins)
	})

	t.Run("newer versions win across stores", func(t *testing.T) {
		logger := testlib.MakeLogger(t)

		static1, err := NewStatic([]*model.Plugin{plugin1V1, plugin2V1, plugin3V1, plugin4V1}, logger)
		require.NoError(t, err)
		static2, err := NewStatic([]*model.Plugin{plugin1V3, plugin2V1, plugin3V3, plugin4V1}, logger)
		require.NoError(t, err)

		store := NewMerged(logger, static1, static2)
		assert.NoError(t, err)
		require.NotNil(t, store)

		plugins, err := store.GetPlugins(&model.PluginFilter{
			Page:    0,
			PerPage: model.AllPerPage,
		})
		require.NoError(t, err)
		assert.Equal(t, []*model.Plugin{
			plugin1V3,
			plugin2V1,
			plugin3V3,
			plugin4V1,
		}, plugins)
	})

	t.Run("later stores win across versions", func(t *testing.T) {
		logger := testlib.MakeLogger(t)

		static1, err := NewStatic([]*model.Plugin{plugin4V1}, logger)
		require.NoError(t, err)
		static2, err := NewStatic([]*model.Plugin{plugin4V1Later}, logger)
		require.NoError(t, err)
		static3, err := NewStatic([]*model.Plugin{plugin1V3}, logger)
		require.NoError(t, err)

		store := NewMerged(logger, static1, static2, static3)
		assert.NoError(t, err)
		require.NotNil(t, store)

		plugins, err := store.GetPlugins(&model.PluginFilter{
			Page:    0,
			PerPage: model.AllPerPage,
		})
		require.NoError(t, err)
		assert.Equal(t, []*model.Plugin{
			plugin1V3,
			plugin4V1Later,
		}, plugins)
	})
}
