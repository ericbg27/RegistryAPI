package api_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/ericbg27/RegistryAPI/db"
	mockdb "github.com/ericbg27/RegistryAPI/db/mock"
	"github.com/ericbg27/RegistryAPI/token"
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
			buildStubs: func(dbConnector *mockdb.MockDBConnector) {
				arg := db.CreateUserParams{
					FullName: user.FullName,
					Phone:    user.Phone,
					UserName: user.UserName,
					Password: user.Password,
				}

				dbConnector.
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
				require.Equal(t, true, ok)

				message, ok = message.(string)
				require.Equal(t, true, ok)
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
			buildStubs: func(dbConnector *mockdb.MockDBConnector) {
				dbConnector.
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
			buildStubs: func(dbConnector *mockdb.MockDBConnector) {
				arg := db.CreateUserParams{
					FullName: user.FullName,
					Phone:    user.Phone,
					UserName: user.UserName,
					Password: user.Password,
				}

				dbConnector.
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
			buildStubs: func(dbConnector *mockdb.MockDBConnector) {
				arg := db.CreateUserParams{
					FullName: user.FullName,
					Phone:    user.Phone,
					UserName: user.UserName,
					Password: user.Password,
				}

				dbConnector.
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

			url := "/v1/user/"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.Router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func TestGetUser(t *testing.T) {
	user := db.User{
		FullName: "Test User",
		Phone:    "99989992",
		UserName: "testuser123",
		Password: "secret",
	}
	user.ID = 0

	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(dbConnector *mockdb.MockDBConnector)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"user_name": user.UserName,
			},
			buildStubs: func(dbConnector *mockdb.MockDBConnector) {
				dbConnector.
					EXPECT().
					GetUser(gomock.Eq(user.UserName)).
					Times(1).
					Return(&user, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)

				data, err := ioutil.ReadAll(recorder.Body)
				require.NoError(t, err)

				var bodyData map[string]any
				err = json.Unmarshal(data, &bodyData)
				require.NoError(t, err)

				fullName, ok := bodyData["full_name"]
				require.Equal(t, true, ok)

				fullName, ok = fullName.(string)
				require.Equal(t, true, ok)
				require.Equal(t, user.FullName, fullName)

				phone, ok := bodyData["phone"]
				require.Equal(t, true, ok)

				phone, ok = phone.(string)
				require.Equal(t, true, ok)
				require.Equal(t, user.Phone, phone)

				userName, ok := bodyData["user_name"]
				require.Equal(t, true, ok)

				userName, ok = userName.(string)
				require.Equal(t, true, ok)
				require.Equal(t, user.UserName, userName)
			},
		},
		{
			name: "Bad Request",
			body: gin.H{},
			buildStubs: func(dbConnector *mockdb.MockDBConnector) {
				dbConnector.
					EXPECT().
					GetUser(gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				validateErrorResponse(t, recorder, "BadRequest", "Incorrect parameters sent in request", http.StatusBadRequest)
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

			url := "/v1/user/"
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			q := request.URL.Query()
			for k, v := range tc.body {
				value, ok := v.(string)
				require.Equal(t, true, ok)

				q.Add(k, value)
			}

			request.URL.RawQuery = q.Encode()

			server.Router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func TestGetUsers(t *testing.T) {
	users := []db.User{}

	for i := 0; i < 5; i++ {
		user := db.User{
			FullName: "Test User " + strconv.Itoa(i),
			Phone:    "9998999" + strconv.Itoa(i),
			UserName: "testuser" + strconv.Itoa(i),
			Password: "secret" + strconv.Itoa(i),
		}

		users = append(users, user)
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
				"page":   0,
				"offset": 2,
			},
			buildStubs: func(dbConnector *mockdb.MockDBConnector) {
				args := db.GetUsersParams{
					PageIndex: 0,
					Offset:    2,
				}

				minIndex := (args.PageIndex * args.Offset)
				maxIndex := minIndex + args.Offset
				if maxIndex > len(users)-1 {
					maxIndex = len(users) - 1
				}

				dbConnector.
					EXPECT().
					GetUsers(gomock.Eq(args)).
					Times(1).
					Return(users[minIndex:maxIndex], nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)

				data, err := ioutil.ReadAll(recorder.Body)
				require.NoError(t, err)

				var bodyData map[string]any
				err = json.Unmarshal(data, &bodyData)
				require.NoError(t, err)

				usersRes, ok := bodyData["users"]
				require.Equal(t, true, ok)

				usersArr, ok := usersRes.([]interface{})
				require.Equal(t, true, ok)

				require.Equal(t, 2, len(usersArr))

				for i, userObj := range usersArr {
					userRes, ok := userObj.(map[string]interface{})
					require.Equal(t, true, ok)

					fullName, ok := userRes["full_name"]
					require.Equal(t, true, ok)

					fullName, ok = fullName.(string)
					require.Equal(t, true, ok)
					require.Equal(t, users[i].FullName, fullName)

					phone, ok := userRes["phone"]
					require.Equal(t, true, ok)

					phone, ok = phone.(string)
					require.Equal(t, true, ok)
					require.Equal(t, users[i].Phone, phone)

					userName, ok := userRes["user_name"]
					require.Equal(t, true, ok)

					userName, ok = userName.(string)
					require.Equal(t, true, ok)
					require.Equal(t, users[i].UserName, userName)
				}
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

			url := "/v1/users/"
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			q := request.URL.Query()
			for k, v := range tc.body {
				value, ok := v.(int)
				require.Equal(t, true, ok)

				q.Add(k, strconv.Itoa(value))
			}

			request.URL.RawQuery = q.Encode()

			server.Router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func TestLoginUser(t *testing.T) {
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
		checkResponse func(recorder *httptest.ResponseRecorder, maker token.Maker)
	}{
		{
			name: "OK",
			body: gin.H{
				"user_name": user.UserName,
				"password":  user.Password,
			},
			buildStubs: func(dbConnector *mockdb.MockDBConnector) {
				dbConnector.
					EXPECT().
					GetUser(gomock.Eq(user.UserName)).
					Times(1).
					Return(&user, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder, maker token.Maker) {
				require.Equal(t, http.StatusOK, recorder.Code)

				data, err := ioutil.ReadAll(recorder.Body)
				require.NoError(t, err)

				var bodyData map[string]any
				err = json.Unmarshal(data, &bodyData)
				require.NoError(t, err)

				tokenValue, ok := bodyData["token"]
				require.Equal(t, true, ok)

				token, ok := tokenValue.(string)
				require.Equal(t, true, ok)

				payload, err := maker.VerifyToken(token)
				require.NoError(t, err)

				require.Equal(t, "testuser123", payload.Username)
			},
		},
		{
			name: "BadRequest",
			body: gin.H{
				"user_name": user.UserName,
			},
			buildStubs: func(dbConnector *mockdb.MockDBConnector) {
				dbConnector.
					EXPECT().
					GetUser(gomock.Any).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder, maker token.Maker) {
				validateErrorResponse(t, recorder, "BadRequest", "Incorrect parameters sent in request", http.StatusBadRequest)
			},
		},
		{
			name: "WrongPassword",
			body: gin.H{
				"user_name": user.UserName,
				"password":  "wrongpassword",
			},
			buildStubs: func(dbConnector *mockdb.MockDBConnector) {
				dbConnector.
					EXPECT().
					GetUser(gomock.Eq(user.UserName)).
					Times(1).
					Return(&user, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder, maker token.Maker) {
				validateErrorResponse(t, recorder, "Unauthorized", "Wrong password sent in request", http.StatusUnauthorized)
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

			url := "/v1/user/login"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.Router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder, server.Maker)
		})
	}
}
