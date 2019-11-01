package model

import (
	"bytes"
	"testing"

	mattermostModel "github.com/mattermost/mattermost-server/model"

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
			`{"homepage_url":"https://github.com/mattermost/mattermost-plugin-demo","icon_data":"icon-data.svg","download_url":"https://github.com/mattermost/mattermost-plugin-demo/releases/download/v0.1.0/com.mattermost.demo-plugin-0.1.0.tar.gz","download_signature":"c2lnbmF0dXJl","release_notes_url":"https://github.com/mattermost/mattermost-plugin-demo/releases/v0.1.0","manifest":{}}`,
		)))
		require.NoError(t, err)
		require.Equal(t, &Plugin{
			HomepageURL:       "https://github.com/mattermost/mattermost-plugin-demo",
			IconData:          "icon-data.svg",
			DownloadURL:       "https://github.com/mattermost/mattermost-plugin-demo/releases/download/v0.1.0/com.mattermost.demo-plugin-0.1.0.tar.gz",
			DownloadSignature: []byte("signature"),
			ReleaseNotesURL:   "https://github.com/mattermost/mattermost-plugin-demo/releases/v0.1.0",
			Manifest:          &mattermostModel.Manifest{},
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
			`[{"homepage_url":"https://github.com/mattermost/mattermost-plugin-demo","icon_data":"icon-data.svg","download_url":"https://github.com/mattermost/mattermost-plugin-demo/releases/download/v0.1.0/com.mattermost.demo-plugin-0.1.0.tar.gz","download_signature":"c2lnbmF0dXJl","release_notes_url":"https://github.com/mattermost/mattermost-plugin-demo/releases/v0.1.0","manifest":{}},{"homepage_url":"https://github.com/mattermost/mattermost-plugin-starter-template","icon_data":"icon-data2.svg","download_url":"https://github.com/mattermost/mattermost-plugin-starter-template/releases/download/v0.1.0/com.mattermost.plugin-starter-template-0.1.0.tar.gz","download_signature":"c2lnbmF0dXJlMg==","release_notes_url":"https://github.com/mattermost/mattermost-plugin-starter-template/releases/v0.1.0","manifest":{}}]`,
		)))
		require.NoError(t, err)
		require.Equal(t, []*Plugin{
			{
				HomepageURL:       "https://github.com/mattermost/mattermost-plugin-demo",
				IconData:          "icon-data.svg",
				DownloadURL:       "https://github.com/mattermost/mattermost-plugin-demo/releases/download/v0.1.0/com.mattermost.demo-plugin-0.1.0.tar.gz",
				DownloadSignature: []byte("signature"),
				ReleaseNotesURL:   "https://github.com/mattermost/mattermost-plugin-demo/releases/v0.1.0",
				Manifest:          &mattermostModel.Manifest{},
			},
			{
				HomepageURL:       "https://github.com/mattermost/mattermost-plugin-starter-template",
				IconData:          "icon-data2.svg",
				DownloadURL:       "https://github.com/mattermost/mattermost-plugin-starter-template/releases/download/v0.1.0/com.mattermost.plugin-starter-template-0.1.0.tar.gz",
				DownloadSignature: []byte("signature2"),
				ReleaseNotesURL:   "https://github.com/mattermost/mattermost-plugin-starter-template/releases/v0.1.0",
				Manifest:          &mattermostModel.Manifest{},
			},
		}, plugin)
	})
}
