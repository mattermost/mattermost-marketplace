package store

import (
	"bytes"
	"encoding/json"
	"testing"

	mattermostModel "github.com/mattermost/mattermost-server/v5/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mattermost/mattermost-marketplace/internal/model"
	"github.com/mattermost/mattermost-marketplace/internal/testlib"
)

func TestNewStatic(t *testing.T) {
	t.Run("empty stream", func(t *testing.T) {
		logger := testlib.MakeLogger(t)
		store, err := NewStatic([]*model.Plugin{}, logger)
		assert.NoError(t, err)
		require.NotNil(t, store)
		assert.Empty(t, store.plugins)
	})

	t.Run("missing manifest id", func(t *testing.T) {
		logger := testlib.MakeLogger(t)
		store, err := NewStatic([]*model.Plugin{
			{
				HomepageURL:     "https://github.com/mattermost/mattermost-plugin-demo",
				IconData:        "icon-data.svg",
				DownloadURL:     "https://github.com/mattermost/mattermost-plugin-demo/releases/download/v0.1.0/com.mattermost.demo-plugin-0.1.0.tar.gz",
				Signature:       "c2lnbmF0dXJl",
				ReleaseNotesURL: "https://github.com/mattermost/mattermost-plugin-demo/releases/v0.1.0",
				Manifest:        &mattermostModel.Manifest{},
			},
			{
				HomepageURL:     "https://github.com/mattermost/mattermost-plugin-starter-template",
				DownloadURL:     "https://github.com/mattermost/mattermost-plugin-starter-template/releases/download/v0.1.0/com.mattermost.plugin-starter-template-0.1.0.tar.gz",
				Signature:       "signature2",
				ReleaseNotesURL: "https://github.com/mattermost/mattermost-plugin-starter-template/releases/v0.1.0",
				Manifest:        &mattermostModel.Manifest{},
			},
		}, logger)
		assert.Error(t, err)
		assert.Nil(t, store)
	})

	t.Run("missing manifest version", func(t *testing.T) {
		logger := testlib.MakeLogger(t)
		store, err := NewStatic([]*model.Plugin{
			{
				HomepageURL:     "https://github.com/mattermost/mattermost-plugin-demo",
				IconData:        "icon-data.svg",
				DownloadURL:     "https://github.com/mattermost/mattermost-plugin-demo/releases/download/v0.1.0/com.mattermost.demo-plugin-0.1.0.tar.gz",
				Signature:       "c2lnbmF0dXJl",
				ReleaseNotesURL: "https://github.com/mattermost/mattermost-plugin-demo/releases/v0.1.0",
				Manifest: &mattermostModel.Manifest{
					Id: "test",
				},
			},
			{
				HomepageURL:     "https://github.com/mattermost/mattermost-plugin-starter-template",
				DownloadURL:     "https://github.com/mattermost/mattermost-plugin-starter-template/releases/download/v0.1.0/com.mattermost.plugin-starter-template-0.1.0.tar.gz",
				Signature:       "signature2",
				ReleaseNotesURL: "https://github.com/mattermost/mattermost-plugin-starter-template/releases/v0.1.0",
				Manifest: &mattermostModel.Manifest{
					Id: "test",
				},
			},
		}, logger)
		assert.Error(t, err)
		assert.Nil(t, store)
	})

	t.Run("missing min_server_version version is valid", func(t *testing.T) {
		logger := testlib.MakeLogger(t)
		store, err := NewStatic([]*model.Plugin{
			{
				HomepageURL:     "https://github.com/mattermost/mattermost-plugin-demo",
				IconData:        "icon-data.svg",
				DownloadURL:     "https://github.com/mattermost/mattermost-plugin-demo/releases/download/v0.1.0/com.mattermost.demo-plugin-0.1.0.tar.gz",
				Signature:       "c2lnbmF0dXJl",
				ReleaseNotesURL: "https://github.com/mattermost/mattermost-plugin-demo/releases/v0.1.0",
				Manifest: &mattermostModel.Manifest{
					Id:      "test",
					Version: "0.1.0",
				},
			},
			{
				HomepageURL:     "https://github.com/mattermost/mattermost-plugin-starter-template",
				DownloadURL:     "https://github.com/mattermost/mattermost-plugin-starter-template/releases/download/v0.1.0/com.mattermost.plugin-starter-template-0.1.0.tar.gz",
				Signature:       "signature2",
				ReleaseNotesURL: "https://github.com/mattermost/mattermost-plugin-starter-template/releases/v0.1.0",
				Manifest: &mattermostModel.Manifest{
					Id:      "test",
					Version: "0.1.0",
				},
			},
		}, logger)
		assert.NoError(t, err)
		assert.NotNil(t, store)
	})

	t.Run("valid stream", func(t *testing.T) {
		logger := testlib.MakeLogger(t)
		store, err := NewStatic([]*model.Plugin{
			{
				HomepageURL:     "https://github.com/mattermost/mattermost-plugin-demo",
				IconData:        "icon-data.svg",
				DownloadURL:     "https://github.com/mattermost/mattermost-plugin-demo/releases/download/v0.1.0/com.mattermost.demo-plugin-0.1.0.tar.gz",
				Signature:       "c2lnbmF0dXJl",
				ReleaseNotesURL: "https://github.com/mattermost/mattermost-plugin-demo/releases/v0.1.0",
				Manifest: &mattermostModel.Manifest{
					Id:               "test",
					Version:          "0.1.0",
					MinServerVersion: "5.23.0",
				},
			},
			{
				HomepageURL:     "https://github.com/mattermost/mattermost-plugin-starter-template",
				DownloadURL:     "https://github.com/mattermost/mattermost-plugin-starter-template/releases/download/v0.1.0/com.mattermost.plugin-starter-template-0.1.0.tar.gz",
				Signature:       "signature2",
				ReleaseNotesURL: "https://github.com/mattermost/mattermost-plugin-starter-template/releases/v0.1.0",
				Manifest: &mattermostModel.Manifest{
					Id:               "test",
					Version:          "0.1.0",
					MinServerVersion: "5.23.0",
				},
			},
		}, logger)
		assert.NoError(t, err)
		assert.NotNil(t, store)
	})
}

func TestNewStaticFromReader(t *testing.T) {
	t.Run("empty stream", func(t *testing.T) {
		logger := testlib.MakeLogger(t)
		store, err := NewStaticFromReader(bytes.NewReader([]byte{}), logger)
		assert.NoError(t, err)
		require.NotNil(t, store)
		assert.Empty(t, store.plugins)
	})

	t.Run("invalid stream", func(t *testing.T) {
		logger := testlib.MakeLogger(t)
		store, err := NewStaticFromReader(bytes.NewReader([]byte(`{"invalid":`)), logger)
		assert.EqualError(t, err, "failed to parse stream: unexpected EOF")
		assert.Nil(t, store)
	})

	t.Run("missing manifest id", func(t *testing.T) {
		logger := testlib.MakeLogger(t)
		store, err := NewStaticFromReader(bytes.NewReader([]byte(`[{"HomepageURL":"https://github.com/mattermost/mattermost-plugin-demo","IconData":"icon-data.svg","DownloadURL":"https://github.com/mattermost/mattermost-plugin-demo/releases/download/v0.1.0/com.mattermost.demo-plugin-0.1.0.tar.gz","Signature":"c2lnbmF0dXJl","ReleaseNotesURL":"https://github.com/mattermost/mattermost-plugin-demo/releases/v0.1.0","Manifest":{}},{"HomepageURL":"https://github.com/mattermost/mattermost-plugin-starter-template","DownloadURL":"https://github.com/mattermost/mattermost-plugin-starter-template/releases/download/v0.1.0/com.mattermost.plugin-starter-template-0.1.0.tar.gz","Signature":"signature2","ReleaseNotesURL":"https://github.com/mattermost/mattermost-plugin-starter-template/releases/v0.1.0","Manifest":{}}]`)), logger)
		assert.Error(t, err)
		assert.Nil(t, store)
	})

	t.Run("missing manifest version", func(t *testing.T) {
		logger := testlib.MakeLogger(t)
		store, err := NewStaticFromReader(bytes.NewReader([]byte(`[{"HomepageURL":"https://github.com/mattermost/mattermost-plugin-demo","IconData":"icon-data.svg","DownloadURL":"https://github.com/mattermost/mattermost-plugin-demo/releases/download/v0.1.0/com.mattermost.demo-plugin-0.1.0.tar.gz","Signature":"c2lnbmF0dXJl","ReleaseNotesURL":"https://github.com/mattermost/mattermost-plugin-demo/releases/v0.1.0","Manifest":{"id": "test"}},{"HomepageURL":"https://github.com/mattermost/mattermost-plugin-starter-template","DownloadURL":"https://github.com/mattermost/mattermost-plugin-starter-template/releases/download/v0.1.0/com.mattermost.plugin-starter-template-0.1.0.tar.gz","Signature":"signature2"],"ReleaseNotesURL":"https://github.com/mattermost/mattermost-plugin-starter-template/releases/v0.1.0","Manifest":{"id": "test"}}]`)), logger)
		assert.Error(t, err)
		assert.Nil(t, store)
	})

	t.Run("missing min_server_version version is valid", func(t *testing.T) {
		logger := testlib.MakeLogger(t)
		store, err := NewStaticFromReader(bytes.NewReader([]byte(`[{"HomepageURL":"https://github.com/mattermost/mattermost-plugin-demo","IconData":"icon-data.svg","DownloadURL":"https://github.com/mattermost/mattermost-plugin-demo/releases/download/v0.1.0/com.mattermost.demo-plugin-0.1.0.tar.gz","Signature":"c2lnbmF0dXJl","ReleaseNotesURL":"https://github.com/mattermost/mattermost-plugin-demo/releases/v0.1.0","Manifest":{"id": "test", "version": "0.1.0"}},{"HomepageURL":"https://github.com/mattermost/mattermost-plugin-starter-template","DownloadURL":"https://github.com/mattermost/mattermost-plugin-starter-template/releases/download/v0.1.0/com.mattermost.plugin-starter-template-0.1.0.tar.gz","Signature":"signature2","ReleaseNotesURL":"https://github.com/mattermost/mattermost-plugin-starter-template/releases/v0.1.0","Manifest":{"id": "test", "version": "0.1.0"}}]`)), logger)
		assert.NoError(t, err)
		assert.NotNil(t, store)
	})

	t.Run("valid stream", func(t *testing.T) {
		logger := testlib.MakeLogger(t)
		store, err := NewStaticFromReader(bytes.NewReader([]byte(`[{"HomepageURL":"https://github.com/mattermost/mattermost-plugin-demo","IconData":"icon-data.svg","DownloadURL":"https://github.com/mattermost/mattermost-plugin-demo/releases/download/v0.1.0/com.mattermost.demo-plugin-0.1.0.tar.gz","Signature":"c2lnbmF0dXJl","ReleaseNotesURL":"https://github.com/mattermost/mattermost-plugin-demo/releases/v0.1.0","Manifest":{"id": "test", "version": "0.1.0", "min_server_version":"5.23.0"}},{"HomepageURL":"https://github.com/mattermost/mattermost-plugin-starter-template","DownloadURL":"https://github.com/mattermost/mattermost-plugin-starter-template/releases/download/v0.1.0/com.mattermost.plugin-starter-template-0.1.0.tar.gz","Signature":"signature2","ReleaseNotesURL":"https://github.com/mattermost/mattermost-plugin-starter-template/releases/v0.1.0","Manifest":{"id": "test", "version": "0.1.0", "min_server_version":"5.23.0"}}]`)), logger)
		assert.NoError(t, err)
		assert.NotNil(t, store)
	})
}

func TestStaticGetPlugins(t *testing.T) {
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

	// earlier will never appear, since later instance with same version overrides
	starterPluginV1Min515Earlier := &model.Plugin{
		HomepageURL: "https://github.com/mattermost/mattermost-plugin-starter-template-earlier",
		IconData:    "icon-data2-earlier.svg",
		DownloadURL: "https://github.com/mattermost/mattermost-plugin-starter-template-earlier/releases/download/v0.1.0/com.mattermost.plugin-starter-template-0.1.0.tar.gz",
		Manifest: &mattermostModel.Manifest{
			Id:               "com.mattermost.plugin-starter-template",
			Name:             "Plugin Starter Template (Earlier)",
			Description:      "This plugin serves as a starting point for writing a Mattermost plugin.",
			Version:          "0.1.0",
			MinServerVersion: "5.15.0",
		},
		Signature: "signature2-earlier",
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
		starterPluginV1Min515Earlier,
		starterPluginV1Min515,
	})
	require.NoError(t, err)

	logger := testlib.MakeLogger(t)
	staticStore, err := NewStaticFromReader(bytes.NewReader(data), logger)
	require.NoError(t, err)

	t.Run("page 0, per page 0", func(t *testing.T) {
		actualPlugins, err := staticStore.GetPlugins(&model.PluginFilter{
			Page:    0,
			PerPage: 0,
			Filter:  "",
		})
		require.NoError(t, err)
		require.Empty(t, actualPlugins)
	})

	t.Run("page 0, per page 1", func(t *testing.T) {
		actualPlugins, err := staticStore.GetPlugins(&model.PluginFilter{
			Page:    0,
			PerPage: 1,
			Filter:  "",
		})
		require.NoError(t, err)
		require.Equal(t, []*model.Plugin{demoPluginV2Min515}, actualPlugins)
	})

	t.Run("page 0, per page 10", func(t *testing.T) {
		actualPlugins, err := staticStore.GetPlugins(&model.PluginFilter{
			Page:    0,
			PerPage: 10,
			Filter:  "",
		})
		require.NoError(t, err)
		require.Equal(t, []*model.Plugin{demoPluginV2Min515, starterPluginV1Min515}, actualPlugins)
	})

	t.Run("page 0, per page 1", func(t *testing.T) {
		actualPlugins, err := staticStore.GetPlugins(&model.PluginFilter{
			Page:    0,
			PerPage: 1,
			Filter:  "",
		})
		require.NoError(t, err)
		require.Equal(t, []*model.Plugin{demoPluginV2Min515}, actualPlugins)
	})

	t.Run("page 0, per page 10", func(t *testing.T) {
		actualPlugins, err := staticStore.GetPlugins(&model.PluginFilter{
			Page:    0,
			PerPage: 10,
			Filter:  "",
		})
		require.NoError(t, err)
		require.Equal(t, []*model.Plugin{demoPluginV2Min515, starterPluginV1Min515}, actualPlugins)
	})

	t.Run("default paging", func(t *testing.T) {
		actualPlugins, err := staticStore.GetPlugins(&model.PluginFilter{PerPage: model.AllPerPage,
			Filter: "",
		})
		require.NoError(t, err)
		require.Equal(t, []*model.Plugin{demoPluginV2Min515, starterPluginV1Min515}, actualPlugins)
	})

	t.Run("filter spaces", func(t *testing.T) {
		actualPlugins, err := staticStore.GetPlugins(&model.PluginFilter{PerPage: model.AllPerPage,
			Filter: "  ",
		})
		require.NoError(t, err)
		require.Equal(t, []*model.Plugin{demoPluginV2Min515, starterPluginV1Min515}, actualPlugins)
	})

	t.Run("id match, exact", func(t *testing.T) {
		actualPlugins, err := staticStore.GetPlugins(&model.PluginFilter{PerPage: model.AllPerPage,
			Filter: "com.mattermost.demo-plugin",
		})
		require.NoError(t, err)
		require.Equal(t, []*model.Plugin{demoPluginV2Min515}, actualPlugins)
	})

	t.Run("id match, case-insensitive", func(t *testing.T) {
		actualPlugins, err := staticStore.GetPlugins(&model.PluginFilter{PerPage: model.AllPerPage,
			Filter: "com.mattermost.demo-PLUGIN",
		})
		require.NoError(t, err)
		require.Equal(t, []*model.Plugin{demoPluginV2Min515}, actualPlugins)
	})

	t.Run("name match, exact", func(t *testing.T) {
		actualPlugins, err := staticStore.GetPlugins(&model.PluginFilter{PerPage: model.AllPerPage,
			Filter: "Plugin Starter Template",
		})
		require.NoError(t, err)
		require.Equal(t, []*model.Plugin{starterPluginV1Min515}, actualPlugins)
	})

	t.Run("name match, partial", func(t *testing.T) {
		actualPlugins, err := staticStore.GetPlugins(&model.PluginFilter{PerPage: model.AllPerPage,
			Filter: "Starter",
		})
		require.NoError(t, err)
		require.Equal(t, []*model.Plugin{starterPluginV1Min515}, actualPlugins)
	})

	t.Run("name match, case-insensitive", func(t *testing.T) {
		actualPlugins, err := staticStore.GetPlugins(&model.PluginFilter{PerPage: model.AllPerPage,
			Filter: "TEMPLATE",
		})
		require.NoError(t, err)
		require.Equal(t, []*model.Plugin{starterPluginV1Min515}, actualPlugins)
	})

	t.Run("description match, partial", func(t *testing.T) {
		actualPlugins, err := staticStore.GetPlugins(&model.PluginFilter{PerPage: model.AllPerPage,
			Filter: "capabilities",
		})
		require.NoError(t, err)
		require.Equal(t, []*model.Plugin{demoPluginV2Min515}, actualPlugins)
	})

	t.Run("description match, case-insensitive, multiple matches", func(t *testing.T) {
		actualPlugins, err := staticStore.GetPlugins(&model.PluginFilter{PerPage: model.AllPerPage,
			Filter: "MATTERMOST",
		})
		require.NoError(t, err)
		require.Equal(t, []*model.Plugin{demoPluginV2Min515, starterPluginV1Min515}, actualPlugins)
	})

	t.Run("plugins that satisfy 5.15", func(t *testing.T) {
		actualPlugins, err := staticStore.GetPlugins(&model.PluginFilter{PerPage: model.AllPerPage,
			Filter:        "MATTERMOST",
			ServerVersion: "5.15.0",
		})
		require.NoError(t, err)
		require.Equal(t, []*model.Plugin{demoPluginV2Min515, starterPluginV1Min515}, actualPlugins)
	})

	t.Run("plugins that satisfy 5.14", func(t *testing.T) {
		actualPlugins, err := staticStore.GetPlugins(&model.PluginFilter{PerPage: model.AllPerPage,
			Filter:        "MATTERMOST",
			ServerVersion: "5.14.0",
		})
		require.NoError(t, err)
		require.Equal(t, []*model.Plugin{demoPluginV1Min514}, actualPlugins)
	})

	t.Run("with a server version that does not satisfy any plugin", func(t *testing.T) {
		actualPlugins, err := staticStore.GetPlugins(&model.PluginFilter{PerPage: model.AllPerPage,
			ServerVersion: "5.13.0",
		})
		require.NoError(t, err)
		require.Nil(t, actualPlugins)
	})

	// Single plugin tests

	t.Run("page 0, per page 0", func(t *testing.T) {
		actualPlugins, err := staticStore.GetPlugins(&model.PluginFilter{
			Page:              0,
			PerPage:           0,
			Filter:            "",
			PluginID:          "com.mattermost.demo-plugin",
			ReturnAllVersions: true,
		})
		require.NoError(t, err)
		require.Empty(t, actualPlugins)
	})

	t.Run("page 0, per page 1", func(t *testing.T) {
		actualPlugins, err := staticStore.GetPlugins(&model.PluginFilter{
			Page:              0,
			PerPage:           1,
			Filter:            "",
			PluginID:          "com.mattermost.demo-plugin",
			ReturnAllVersions: true,
		})
		require.NoError(t, err)
		require.Equal(t, []*model.Plugin{demoPluginV2Min515}, actualPlugins)
	})

	t.Run("page 0, per page 10", func(t *testing.T) {
		actualPlugins, err := staticStore.GetPlugins(&model.PluginFilter{
			Page:              0,
			PerPage:           10,
			Filter:            "",
			PluginID:          "com.mattermost.demo-plugin",
			ReturnAllVersions: true,
		})
		require.NoError(t, err)
		require.Equal(t, []*model.Plugin{demoPluginV2Min515, demoPluginV1Min514}, actualPlugins)
	})

	t.Run("page 0, per page 1", func(t *testing.T) {
		actualPlugins, err := staticStore.GetPlugins(&model.PluginFilter{
			Page:              0,
			PerPage:           1,
			Filter:            "",
			PluginID:          "com.mattermost.demo-plugin",
			ReturnAllVersions: true,
		})
		require.NoError(t, err)
		require.Equal(t, []*model.Plugin{demoPluginV2Min515}, actualPlugins)
	})

	t.Run("default paging", func(t *testing.T) {
		actualPlugins, err := staticStore.GetPlugins(&model.PluginFilter{PerPage: model.AllPerPage,
			Filter:            "",
			PluginID:          "com.mattermost.demo-plugin",
			ReturnAllVersions: true,
		})
		require.NoError(t, err)
		require.Equal(t, []*model.Plugin{demoPluginV2Min515, demoPluginV1Min514}, actualPlugins)
	})

	t.Run("plugins that satisfy 5.15", func(t *testing.T) {
		actualPlugins, err := staticStore.GetPlugins(&model.PluginFilter{PerPage: model.AllPerPage,
			ServerVersion:     "5.15.0",
			PluginID:          "com.mattermost.demo-plugin",
			ReturnAllVersions: true,
		})
		require.NoError(t, err)
		require.Equal(t, []*model.Plugin{demoPluginV2Min515, demoPluginV1Min514}, actualPlugins)
	})

	t.Run("plugins that satisfy 5.14", func(t *testing.T) {
		actualPlugins, err := staticStore.GetPlugins(&model.PluginFilter{PerPage: model.AllPerPage,
			ServerVersion:     "5.14.0",
			PluginID:          "com.mattermost.demo-plugin",
			ReturnAllVersions: true,
		})
		require.NoError(t, err)
		require.Equal(t, []*model.Plugin{demoPluginV1Min514}, actualPlugins)
	})

	t.Run("with a server version that does not satisfy any plugin", func(t *testing.T) {
		actualPlugins, err := staticStore.GetPlugins(&model.PluginFilter{PerPage: model.AllPerPage,
			ServerVersion:     "5.13.0",
			PluginID:          "com.mattermost.demo-plugin",
			ReturnAllVersions: true,
		})
		require.NoError(t, err)
		require.Nil(t, actualPlugins)
	})
}
