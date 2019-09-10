package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHealthCheck(t *testing.T) {
	router := mux.NewRouter()

	Register(router, &Context{
		Logger: logrus.New(),
	})

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/api/v1/health", nil)
	router.ServeHTTP(w, r)

	result := w.Result()
	require.NotNil(t, result)
	defer result.Body.Close()

	respose := &healthCheckResponse{}
	err := json.NewDecoder(result.Body).Decode(&respose)
	require.NoError(t, err)
	require.NotNil(t, respose)

	assert.Equal(t, respose.Status, "pass")
	assert.Equal(t, respose.Version, "1")
	assert.Equal(t, respose.ReleaseID, "") // This needs to be changed, when the first tag is cut
	assert.Len(t, respose.Notes, 1)
	assert.NotEmpty(t, respose.Notes[0])
	assert.NotEmpty(t, respose.Description)
}
