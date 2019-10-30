package store

import (
	"bytes"
	"testing"

	"github.com/mattermost/mattermost-marketplace/internal/testlib"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	t.Run("empty stream", func(t *testing.T) {
		logger := testlib.MakeLogger(t)
		store, err := New(bytes.NewReader([]byte{}), logger)
		require.NoError(t, err)
		require.NotNil(t, store)
		require.Empty(t, store.plugins)
	})

	t.Run("invalid stream", func(t *testing.T) {
		logger := testlib.MakeLogger(t)
		store, err := New(bytes.NewReader([]byte(`{"invalid":`)), logger)
		require.EqualError(t, err, "failed to parse stream: unexpected EOF")
		require.Nil(t, store)
	})

	t.Run("valid stream", func(t *testing.T) {
		logger := testlib.MakeLogger(t)
		store, err := New(bytes.NewReader([]byte(`[{"HomepageURL":"https://github.com/mattermost/mattermost-plugin-demo","IconData":"icon-data.svg","DownloadURL":"https://github.com/mattermost/mattermost-plugin-demo/releases/download/v0.1.0/com.mattermost.demo-plugin-0.1.0.tar.gz","Manifest":{},"signatures":[{"signature":"signature1","public_key_hash":"hash1"}]},{"HomepageURL":"https://github.com/mattermost/mattermost-plugin-starter-template","DownloadURL":"https://github.com/mattermost/mattermost-plugin-starter-template/releases/download/v0.1.0/com.mattermost.plugin-starter-template-0.1.0.tar.gz","Manifest":{},"signatures":[{"signature":"signature2","public_key_hash":"hash2"}]}]`)), logger)
		require.NoError(t, err)
		require.NotNil(t, store)
	})
}
