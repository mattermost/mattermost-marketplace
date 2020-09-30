package model

import (
	"bytes"
	"testing"

	mattermostModel "github.com/mattermost/mattermost-server/v5/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPluginFromReader(t *testing.T) {
	t.Run("empty request", func(t *testing.T) {
		plugin, err := PluginFromReader(bytes.NewReader([]byte(
			``,
		)))
		require.NoError(t, err)
		require.Equal(t, &Plugin{}, plugin)
	})

	t.Run("invalid request", func(t *testing.T) {
		plugin, err := PluginFromReader(bytes.NewReader([]byte(
			`{test`,
		)))
		require.Error(t, err)
		require.Nil(t, plugin)
	})

	t.Run("request", func(t *testing.T) {
		plugin, err := PluginFromReader(bytes.NewReader([]byte(
			`{"homepage_url":"https://github.com/mattermost/mattermost-plugin-demo","icon_data":"icon-data.svg","download_url":"https://github.com/mattermost/mattermost-plugin-demo/releases/download/v0.1.0/com.mattermost.demo-plugin-0.1.0.tar.gz","signature":"signature1", "release_notes_url":"https://github.com/mattermost/mattermost-plugin-demo/releases/v0.1.0","manifest":{}}`,
		)))
		require.NoError(t, err)
		require.Equal(t, &Plugin{
			HomepageURL:     "https://github.com/mattermost/mattermost-plugin-demo",
			IconData:        "icon-data.svg",
			DownloadURL:     "https://github.com/mattermost/mattermost-plugin-demo/releases/download/v0.1.0/com.mattermost.demo-plugin-0.1.0.tar.gz",
			Signature:       "signature1",
			ReleaseNotesURL: "https://github.com/mattermost/mattermost-plugin-demo/releases/v0.1.0",
			Manifest:        &mattermostModel.Manifest{},
		}, plugin)
	})
}

func TestPluginsFromReader(t *testing.T) {
	t.Run("empty request", func(t *testing.T) {
		plugins, err := PluginsFromReader(bytes.NewReader([]byte(
			``,
		)))
		require.NoError(t, err)
		require.Equal(t, []*Plugin{}, plugins)
	})

	t.Run("invalid request", func(t *testing.T) {
		plugins, err := PluginsFromReader(bytes.NewReader([]byte(
			`{test`,
		)))
		require.Error(t, err)
		require.Nil(t, plugins)
	})

	t.Run("request", func(t *testing.T) {
		plugin, err := PluginsFromReader(bytes.NewReader([]byte(
			`[{"homepage_url":"https://github.com/mattermost/mattermost-plugin-demo","icon_data":"icon-data.svg","download_url":"https://github.com/mattermost/mattermost-plugin-demo/releases/download/v0.1.0/com.mattermost.demo-plugin-0.1.0.tar.gz","signature":"signature1","release_notes_url":"https://github.com/mattermost/mattermost-plugin-demo/releases/v0.1.0","manifest":{},"platform_bundles":{"linux-amd64":{"download_url":"https://github.com/mattermost/mattermost-plugin-demo/releases/download/v0.1.0/com.mattermost.demo-plugin-0.1.0-linux-amd64.tar.gz","signature":"signature for linux"}}},{"homepage_url":"https://github.com/mattermost/mattermost-plugin-starter-template","icon_data":"icon-data2.svg","download_url":"https://github.com/mattermost/mattermost-plugin-starter-template/releases/download/v0.1.0/com.mattermost.plugin-starter-template-0.1.0.tar.gz","signature":"signature2","release_notes_url":"https://github.com/mattermost/mattermost-plugin-starter-template/releases/v0.1.0","manifest":{}}]`,
		)))
		require.NoError(t, err)
		require.Equal(t, []*Plugin{
			{
				HomepageURL:     "https://github.com/mattermost/mattermost-plugin-demo",
				IconData:        "icon-data.svg",
				DownloadURL:     "https://github.com/mattermost/mattermost-plugin-demo/releases/download/v0.1.0/com.mattermost.demo-plugin-0.1.0.tar.gz",
				Signature:       "signature1",
				ReleaseNotesURL: "https://github.com/mattermost/mattermost-plugin-demo/releases/v0.1.0",
				Manifest:        &mattermostModel.Manifest{},
				PlatformBundles: PlatformBundles{
					LinuxAmd64: &PlatformBundleMetadata{
						DownloadURL: "https://github.com/mattermost/mattermost-plugin-demo/releases/download/v0.1.0/com.mattermost.demo-plugin-0.1.0-linux-amd64.tar.gz",
						Signature:   "signature for linux",
					},
				},
			},
			{
				HomepageURL:     "https://github.com/mattermost/mattermost-plugin-starter-template",
				IconData:        "icon-data2.svg",
				DownloadURL:     "https://github.com/mattermost/mattermost-plugin-starter-template/releases/download/v0.1.0/com.mattermost.plugin-starter-template-0.1.0.tar.gz",
				Signature:       "signature2",
				ReleaseNotesURL: "https://github.com/mattermost/mattermost-plugin-starter-template/releases/v0.1.0",
				Manifest:        &mattermostModel.Manifest{},
			},
		}, plugin)
	})
}

func TestPluginsToWriter(t *testing.T) {
	t.Run("nil request", func(t *testing.T) {
		var b bytes.Buffer

		err := PluginsToWriter(&b, nil)

		require.NoError(t, err)
		assert.Equal(t, "null\n", b.String())
	})

	t.Run("empty request", func(t *testing.T) {
		var b bytes.Buffer
		p := []*Plugin{}

		err := PluginsToWriter(&b, p)

		require.NoError(t, err)
		assert.Equal(t, "[]\n", b.String())
	})

	t.Run("empty request", func(t *testing.T) {
		var b bytes.Buffer
		p := []*Plugin{{
			HomepageURL:     "https://github.com/mattermost/mattermost-plugin-demo",
			IconData:        "icon-data.svg",
			DownloadURL:     "https://github.com/mattermost/mattermost-plugin-demo/releases/download/v0.1.0/com.mattermost.demo-plugin-0.1.0.tar.gz",
			Signature:       "signature1",
			ReleaseNotesURL: "https://github.com/mattermost/mattermost-plugin-demo/releases/v0.1.0",
			Manifest: &mattermostModel.Manifest{
				Id:      "demo",
				Version: "1.0.0",
			},
		}, {
			HomepageURL:     "https://github.com/mattermost/mattermost-plugin-starter-template",
			IconData:        "icon-data2.svg",
			DownloadURL:     "https://github.com/mattermost/mattermost-plugin-starter-template/releases/download/v0.1.0/com.mattermost.plugin-starter-template-0.1.0.tar.gz",
			Signature:       "signature2",
			ReleaseNotesURL: "https://github.com/mattermost/mattermost-plugin-starter-template/releases/v0.1.0",
			Manifest: &mattermostModel.Manifest{
				Id:      "template",
				Version: "2.0.0",
			},
		}}

		err := PluginsToWriter(&b, p)
		expectedResult := `[
  {
    "homepage_url": "https://github.com/mattermost/mattermost-plugin-demo",
    "icon_data": "icon-data.svg",
    "download_url": "https://github.com/mattermost/mattermost-plugin-demo/releases/download/v0.1.0/com.mattermost.demo-plugin-0.1.0.tar.gz",
    "release_notes_url": "https://github.com/mattermost/mattermost-plugin-demo/releases/v0.1.0",
    "signature": "signature1",
    "manifest": {
      "id": "demo",
      "version": "1.0.0"
    },
    "enterprise": false,
    "updated_at": "0001-01-01T00:00:00Z",
    "platform_bundles": {}
  },
  {
    "homepage_url": "https://github.com/mattermost/mattermost-plugin-starter-template",
    "icon_data": "icon-data2.svg",
    "download_url": "https://github.com/mattermost/mattermost-plugin-starter-template/releases/download/v0.1.0/com.mattermost.plugin-starter-template-0.1.0.tar.gz",
    "release_notes_url": "https://github.com/mattermost/mattermost-plugin-starter-template/releases/v0.1.0",
    "signature": "signature2",
    "manifest": {
      "id": "template",
      "version": "2.0.0"
    },
    "enterprise": false,
    "updated_at": "0001-01-01T00:00:00Z",
    "platform_bundles": {}
  }
]
`
		require.NoError(t, err)
		assert.Equal(t, expectedResult, b.String())
	})
}
