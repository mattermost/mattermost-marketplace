package client_test

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/mattermost/mattermost-marketplace/client"
	"github.com/mattermost/mattermost-marketplace/internal/api"
	"github.com/mattermost/mattermost-marketplace/internal/model"
	"github.com/mattermost/mattermost-marketplace/internal/store"
	"github.com/mattermost/mattermost-marketplace/internal/testlib"
	"github.com/stretchr/testify/require"
)

func SetupApi(t *testing.T, plugins []*model.Plugin) (*client.Client, string, func()) {
	logger := testlib.MakeLogger(t)

	data, err := json.Marshal(plugins)
	require.NoError(t, err)
	store, err := store.New(bytes.NewReader(data), logger)
	require.NoError(t, err)

	router := mux.NewRouter()
	api.Register(router, &api.Context{
		Store:  store,
		Logger: logger,
	})
	ts := httptest.NewServer(router)

	return client.NewClient(ts.URL), ts.URL, func() {
		ts.Close()
	}
}
