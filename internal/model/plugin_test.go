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
			`{"HomepageURL":"https://github.com/mattermost/mattermost-plugin-demo","IconURL":"http://example.com/icon.svg","DownloadURL":"https://github.com/mattermost/mattermost-plugin-demo/releases/download/v0.1.0/com.mattermost.demo-plugin-0.1.0.tar.gz","DownloadSignature":"c2lnbmF0dXJl","Manifest":{}}`,
		)))
		require.NoError(t, err)
		require.Equal(t, &Plugin{
			HomepageURL:       "https://github.com/mattermost/mattermost-plugin-demo",
			IconURL:           "http://example.com/icon.svg",
			DownloadURL:       "https://github.com/mattermost/mattermost-plugin-demo/releases/download/v0.1.0/com.mattermost.demo-plugin-0.1.0.tar.gz",
			DownloadSignature: []byte("signature"),
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
			`[{"HomepageURL":"https://github.com/mattermost/mattermost-plugin-demo","IconURL":"http://example.com/icon.svg","DownloadURL":"https://github.com/mattermost/mattermost-plugin-demo/releases/download/v0.1.0/com.mattermost.demo-plugin-0.1.0.tar.gz","DownloadSignature":"c2lnbmF0dXJl","Manifest":{}},{"HomepageURL":"https://github.com/mattermost/mattermost-plugin-starter-template","IconURL":"http://example.com/icon2.svg","DownloadURL":"https://github.com/mattermost/mattermost-plugin-starter-template/releases/download/v0.1.0/com.mattermost.plugin-starter-template-0.1.0.tar.gz","DownloadSignature":"c2lnbmF0dXJlMg==","Manifest":{}}]`,
		)))
		require.NoError(t, err)
		require.Equal(t, []*Plugin{
			{
				HomepageURL:       "https://github.com/mattermost/mattermost-plugin-demo",
				IconURL:           "http://example.com/icon.svg",
				DownloadURL:       "https://github.com/mattermost/mattermost-plugin-demo/releases/download/v0.1.0/com.mattermost.demo-plugin-0.1.0.tar.gz",
				DownloadSignature: []byte("signature"),
				Manifest:          &mattermostModel.Manifest{},
			},
			{
				HomepageURL:       "https://github.com/mattermost/mattermost-plugin-starter-template",
				IconURL:           "http://example.com/icon2.svg",
				DownloadURL:       "https://github.com/mattermost/mattermost-plugin-starter-template/releases/download/v0.1.0/com.mattermost.plugin-starter-template-0.1.0.tar.gz",
				DownloadSignature: []byte("signature2"),
				Manifest:          &mattermostModel.Manifest{},
			},
		}, plugin)
	})
}
