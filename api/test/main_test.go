package api_test

import (
	"encoding/json"
	"io/ioutil"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/ericbg27/RegistryAPI/api"
	"github.com/ericbg27/RegistryAPI/db"
	"github.com/ericbg27/RegistryAPI/util"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func NewTestServer(t *testing.T, dbConnector db.DBConnector) *api.Server {
	config := util.Config{
		AccessTokenDuration: 15 * time.Minute,
		TokenSymmetricKey:   util.RandomString(32),
	}

	server, err := api.NewServer(dbConnector, config)
	require.NoError(t, err)

	return server
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	os.Exit(m.Run())
}

func validateErrorResponse(t *testing.T, recorder *httptest.ResponseRecorder, expectedName string, expectedMessage string, expectedStatusCode int) {
	require.Equal(t, expectedStatusCode, recorder.Code)

	data, err := ioutil.ReadAll(recorder.Body)
	require.NoError(t, err)

	var bodyData map[string]any
	err = json.Unmarshal(data, &bodyData)
	require.NoError(t, err)

	name, ok := bodyData["name"]
	require.Equal(t, ok, true)

	name, ok = name.(string)
	require.Equal(t, ok, true)
	require.Equal(t, expectedName, name)

	message, ok := bodyData["message"]
	require.Equal(t, ok, true)

	message, ok = message.(string)
	require.Equal(t, ok, true)
	require.Equal(t, expectedMessage, message)
}
