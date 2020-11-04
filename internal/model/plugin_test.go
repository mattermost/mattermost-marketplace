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
			`{"homepage_url":"https://github.com/mattermost/mattermost-plugin-demo","icon_data":"icon-data.svg","download_url":"https://github.com/mattermost/mattermost-plugin-demo/releases/download/v0.1.0/com.mattermost.demo-plugin-0.1.0.tar.gz","signature":"signature1","repo_name":"mattermost-plugin-demo","release_notes_url":"https://github.com/mattermost/mattermost-plugin-demo/releases/v0.1.0","manifest":{}}`,
		)))
		require.NoError(t, err)
		require.Equal(t, &Plugin{
			RepoName:        "mattermost-plugin-demo",
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
			`[{"repo_name":"mattermost-plugin-demo","homepage_url":"https://github.com/mattermost/mattermost-plugin-demo","icon_data":"icon-data.svg","download_url":"https://github.com/mattermost/mattermost-plugin-demo/releases/download/v0.1.0/com.mattermost.demo-plugin-0.1.0.tar.gz","signature":"signature1","release_notes_url":"https://github.com/mattermost/mattermost-plugin-demo/releases/v0.1.0","manifest":{},"platforms":{"linux-amd64":{"download_url":"https://github.com/mattermost/mattermost-plugin-demo/releases/download/v0.1.0/com.mattermost.demo-plugin-0.1.0-linux-amd64.tar.gz","signature":"signature1 for linux"},"darwin-amd64":{"download_url":"https://github.com/mattermost/mattermost-plugin-demo/releases/download/v0.1.0/com.mattermost.demo-plugin-0.1.0-darwin-amd64.tar.gz","signature":"signature1 for darwin"},"windows-amd64":{"download_url":"https://github.com/mattermost/mattermost-plugin-demo/releases/download/v0.1.0/com.mattermost.demo-plugin-0.1.0-windows-amd64.tar.gz","signature":"signature1 for windows"}}},{"repo_name":"mattermost-plugin-starter-template","homepage_url":"https://github.com/mattermost/mattermost-plugin-starter-template","icon_data":"icon-data2.svg","download_url":"https://github.com/mattermost/mattermost-plugin-starter-template/releases/download/v0.1.0/com.mattermost.plugin-starter-template-0.1.0.tar.gz","signature":"signature2","release_notes_url":"https://github.com/mattermost/mattermost-plugin-starter-template/releases/v0.1.0","manifest":{}}]`,
		)))
		require.NoError(t, err)
		require.Equal(t, []*Plugin{
			{
				RepoName:        "mattermost-plugin-demo",
				HomepageURL:     "https://github.com/mattermost/mattermost-plugin-demo",
				IconData:        "icon-data.svg",
				DownloadURL:     "https://github.com/mattermost/mattermost-plugin-demo/releases/download/v0.1.0/com.mattermost.demo-plugin-0.1.0.tar.gz",
				Signature:       "signature1",
				ReleaseNotesURL: "https://github.com/mattermost/mattermost-plugin-demo/releases/v0.1.0",
				Manifest:        &mattermostModel.Manifest{},
				Platforms: PlatformBundles{
					LinuxAmd64: PlatformBundleMetadata{
						DownloadURL: "https://github.com/mattermost/mattermost-plugin-demo/releases/download/v0.1.0/com.mattermost.demo-plugin-0.1.0-linux-amd64.tar.gz",
						Signature:   "signature1 for linux",
					},
					DarwinAmd64: PlatformBundleMetadata{
						DownloadURL: "https://github.com/mattermost/mattermost-plugin-demo/releases/download/v0.1.0/com.mattermost.demo-plugin-0.1.0-darwin-amd64.tar.gz",
						Signature:   "signature1 for darwin",
					},
					WindowsAmd64: PlatformBundleMetadata{
						DownloadURL: "https://github.com/mattermost/mattermost-plugin-demo/releases/download/v0.1.0/com.mattermost.demo-plugin-0.1.0-windows-amd64.tar.gz",
						Signature:   "signature1 for windows",
					},
				},
			},
			{
				RepoName:        "mattermost-plugin-starter-template",
				HomepageURL:     "https://github.com/mattermost/mattermost-plugin-starter-template",
				IconData:        "icon-data2.svg",
				DownloadURL:     "https://github.com/mattermost/mattermost-plugin-starter-template/releases/download/v0.1.0/com.mattermost.plugin-starter-template-0.1.0.tar.gz",
				Signature:       "signature2",
				ReleaseNotesURL: "https://github.com/mattermost/mattermost-plugin-starter-template/releases/v0.1.0",
				Manifest:        &mattermostModel.Manifest{},
				Platforms:       PlatformBundles{},
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
			RepoName:        "mattermost-plugin-demo",
			HomepageURL:     "https://github.com/mattermost/mattermost-plugin-demo",
			IconData:        "icon-data.svg",
			DownloadURL:     "https://github.com/mattermost/mattermost-plugin-demo/releases/download/v0.1.0/com.mattermost.demo-plugin-0.1.0.tar.gz",
			Signature:       "signature1",
			ReleaseNotesURL: "https://github.com/mattermost/mattermost-plugin-demo/releases/v0.1.0",
			Maintainer:      Mattermost,
			Stage:           Production,
			Manifest: &mattermostModel.Manifest{
				Id:      "demo",
				Version: "1.0.0",
			},
			Platforms: PlatformBundles{
				LinuxAmd64: PlatformBundleMetadata{
					DownloadURL: "https://github.com/mattermost/mattermost-plugin-demo/releases/download/v0.1.0/com.mattermost.plugin.demo-plugin-0.1.0-linux-amd64.tar.gz",
					Signature:   "signature1 for linux",
				},
				DarwinAmd64: PlatformBundleMetadata{
					DownloadURL: "https://github.com/mattermost/mattermost-plugin-demo/releases/download/v0.1.0/com.mattermost.plugin.demo-plugin-0.1.0-darwin-amd64.tar.gz",
					Signature:   "signature1 for darwin",
				},
				WindowsAmd64: PlatformBundleMetadata{
					DownloadURL: "https://github.com/mattermost/mattermost-plugin-demo/releases/download/v0.1.0/com.mattermost.plugin.demo-plugin-0.1.0-windows-amd64.tar.gz",
					Signature:   "signature1 for windows",
				},
			},
			Hosting: OnPrem,
		}, {
			RepoName:        "mattermost-plugin-starter-template",
			HomepageURL:     "https://github.com/mattermost/mattermost-plugin-starter-template",
			IconData:        "icon-data2.svg",
			DownloadURL:     "https://github.com/mattermost/mattermost-plugin-starter-template/releases/download/v0.1.0/com.mattermost.plugin-starter-template-0.1.0.tar.gz",
			Signature:       "signature2",
			ReleaseNotesURL: "https://github.com/mattermost/mattermost-plugin-starter-template/releases/v0.1.0",
			Maintainer:      Community,
			Stage:           Beta,
			Manifest: &mattermostModel.Manifest{
				Id:      "template",
				Version: "2.0.0",
			},
			Platforms: PlatformBundles{},
			Hosting:   Cloud,
		}}

		err := PluginsToWriter(&b, p)
		expectedResult := `[
  {
    "homepage_url": "https://github.com/mattermost/mattermost-plugin-demo",
    "icon_data": "icon-data.svg",
    "download_url": "https://github.com/mattermost/mattermost-plugin-demo/releases/download/v0.1.0/com.mattermost.demo-plugin-0.1.0.tar.gz",
    "release_notes_url": "https://github.com/mattermost/mattermost-plugin-demo/releases/v0.1.0",
    "hosting": "on-prem",
    "maintainer": "mattermost",
    "stage": "production",
    "enterprise": false,
    "signature": "signature1",
    "repo_name": "mattermost-plugin-demo",
    "manifest": {
      "id": "demo",
      "version": "1.0.0"
    },
    "platforms": {
      "linux-amd64": {
        "download_url": "https://github.com/mattermost/mattermost-plugin-demo/releases/download/v0.1.0/com.mattermost.plugin.demo-plugin-0.1.0-linux-amd64.tar.gz",
        "signature": "signature1 for linux"
      },
      "darwin-amd64": {
        "download_url": "https://github.com/mattermost/mattermost-plugin-demo/releases/download/v0.1.0/com.mattermost.plugin.demo-plugin-0.1.0-darwin-amd64.tar.gz",
        "signature": "signature1 for darwin"
      },
      "windows-amd64": {
        "download_url": "https://github.com/mattermost/mattermost-plugin-demo/releases/download/v0.1.0/com.mattermost.plugin.demo-plugin-0.1.0-windows-amd64.tar.gz",
        "signature": "signature1 for windows"
      }
    },
    "updated_at": "0001-01-01T00:00:00Z"
  },
  {
    "homepage_url": "https://github.com/mattermost/mattermost-plugin-starter-template",
    "icon_data": "icon-data2.svg",
    "download_url": "https://github.com/mattermost/mattermost-plugin-starter-template/releases/download/v0.1.0/com.mattermost.plugin-starter-template-0.1.0.tar.gz",
    "release_notes_url": "https://github.com/mattermost/mattermost-plugin-starter-template/releases/v0.1.0",
    "hosting": "cloud",
    "maintainer": "community",
    "stage": "beta",
    "enterprise": false,
    "signature": "signature2",
    "repo_name": "mattermost-plugin-starter-template",
    "manifest": {
      "id": "template",
      "version": "2.0.0"
    },
    "platforms": {
      "linux-amd64": {},
      "darwin-amd64": {},
      "windows-amd64": {}
    },
    "updated_at": "0001-01-01T00:00:00Z"
  }
]
`
		require.NoError(t, err)
		assert.Equal(t, expectedResult, b.String())
	})
}
