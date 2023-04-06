package api_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ericbg27/RegistryAPI/db"
	mockdb "github.com/ericbg27/RegistryAPI/db/mock"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestCreateUser(t *testing.T) {
	user := db.User{
		FullName: "Test User",
		Phone:    "99989992",
		UserName: "testuser123",
		Password: "secret",
	}

	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(dbConnector *mockdb.MockDBConnector)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"full_name": user.FullName,
				"phone":     user.Phone,
				"user_name": user.UserName,
				"password":  user.Password,
			},
			buildStubs: func(dbConnnector *mockdb.MockDBConnector) {
				arg := db.CreateUserParams{
					FullName: user.FullName,
					Phone:    user.Phone,
					UserName: user.UserName,
					Password: user.Password,
				}

				dbConnnector.
					EXPECT().
					CreateUser(gomock.Eq(arg)).
					Times(1).
					Return(&user, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, recorder.Code)

				data, err := ioutil.ReadAll(recorder.Body)
				require.NoError(t, err)

				var bodyData map[string]any
				err = json.Unmarshal(data, &bodyData)
				require.NoError(t, err)

				message, ok := bodyData["message"]
				require.Equal(t, ok, true)

				message, ok = message.(string)
				require.Equal(t, ok, true)
				require.Equal(t, "User created successfully", message)
			},
		},
		{
			name: "Bad Request",
			body: gin.H{
				"full_name": user.FullName,
				"phone":     user.Phone,
				"user_name": user.UserName,
			},
			buildStubs: func(dbConnnector *mockdb.MockDBConnector) {
				dbConnnector.
					EXPECT().
					CreateUser(gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				validateErrorResponse(t, recorder, "BadRequest", "Incorrect parameters sent in request", http.StatusBadRequest)
			},
		},
		{
			name: "User Already Exists",
			body: gin.H{
				"full_name": user.FullName,
				"phone":     user.Phone,
				"user_name": user.UserName,
				"password":  user.Password,
			},
			buildStubs: func(dbConnnector *mockdb.MockDBConnector) {
				arg := db.CreateUserParams{
					FullName: user.FullName,
					Phone:    user.Phone,
					UserName: user.UserName,
					Password: user.Password,
				}

				dbConnnector.
					EXPECT().
					CreateUser(arg).
					Times(1).
					Return(nil, &db.BadInputError{
						Err: fmt.Errorf("An user with the provided information already exists"),
					})
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				validateErrorResponse(t, recorder, "AlreadyExists", "An user with the provided information already exists", http.StatusBadRequest)
			},
		},
		{
			name: "Internal Server Error When Executing DB Query",
			body: gin.H{
				"full_name": user.FullName,
				"phone":     user.Phone,
				"user_name": user.UserName,
				"password":  user.Password,
			},
			buildStubs: func(dbConnnector *mockdb.MockDBConnector) {
				arg := db.CreateUserParams{
					FullName: user.FullName,
					Phone:    user.Phone,
					UserName: user.UserName,
					Password: user.Password,
				}

				dbConnnector.
					EXPECT().
					CreateUser(arg).
					Times(1).
					Return(nil, fmt.Errorf("Error executing query"))
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				validateErrorResponse(t, recorder, "InternalServerError", "Unexpected server error. Try again later", http.StatusInternalServerError)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			dbConnector := mockdb.NewMockDBConnector(ctrl)
			tc.buildStubs(dbConnector)

			server := NewTestServer(t, dbConnector)
			recorder := httptest.NewRecorder()

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/users"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.Router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}
