package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mattermost/mattermost-marketplace/internal/model"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetLabels(t *testing.T) {
	router := mux.NewRouter()

	Register(router, &Context{
		Logger: logrus.New(),
	})

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/api/v1/labels", nil)
	router.ServeHTTP(w, r)

	result := w.Result()
	require.NotNil(t, result)
	defer result.Body.Close()

	var respose []model.Label
	err := json.NewDecoder(result.Body).Decode(&respose)
	require.NoError(t, err)
	assert.Equal(t, model.AllLabels, respose)
}
