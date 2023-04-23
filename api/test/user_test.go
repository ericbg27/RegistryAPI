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
	"time"

	"github.com/ericbg27/RegistryAPI/db"
	mockdb "github.com/ericbg27/RegistryAPI/db/mock"
	"github.com/ericbg27/RegistryAPI/token"
	mocktoken "github.com/ericbg27/RegistryAPI/token/mock"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

const bearerStr = "Bearer "

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
		{
			name: "Bad Phone Number",
			body: gin.H{
				"full_name": user.FullName,
				"phone":     "9999999",
				"user_name": user.UserName,
				"password":  user.Password,
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
			name: "Bad Password",
			body: gin.H{
				"full_name": user.FullName,
				"phone":     user.Phone,
				"user_name": user.UserName,
				"password":  "test#",
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
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			dbConnector := mockdb.NewMockDBConnector(ctrl)
			tc.buildStubs(dbConnector)

			maker := mocktoken.NewMockMaker(ctrl)

			server := NewTestServer(t, dbConnector, maker)
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
		FullName:   "Test User",
		Phone:      "99989992",
		UserName:   "testuser123",
		Password:   "secret",
		LoginToken: "token",
	}
	user.ID = 0

	uuidToken, err := uuid.NewRandom()
	require.NoError(t, err)

	now := time.Now()

	issuedAt := now
	expiredAt := now.Add(time.Hour)

	tokenPayload := &token.Payload{
		ID:        uuidToken,
		Username:  user.UserName,
		IssuedAt:  issuedAt,
		ExpiredAt: expiredAt,
	}

	testCases := []struct {
		name          string
		body          gin.H
		token         string
		buildStubs    func(dbConnector *mockdb.MockDBConnector, maker *mocktoken.MockMaker)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"user_name": user.UserName,
			},
			token: user.LoginToken,
			buildStubs: func(dbConnector *mockdb.MockDBConnector, maker *mocktoken.MockMaker) {
				maker.
					EXPECT().
					VerifyToken(user.LoginToken).
					Times(1).
					Return(tokenPayload, nil)

				dbConnector.
					EXPECT().
					GetUser(gomock.Eq(user.UserName)).
					Times(2).
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
			name:  "Bad Request",
			body:  gin.H{},
			token: user.LoginToken,
			buildStubs: func(dbConnector *mockdb.MockDBConnector, maker *mocktoken.MockMaker) {
				maker.
					EXPECT().
					VerifyToken(user.LoginToken).
					Times(1).
					Return(tokenPayload, nil)

				dbConnector.
					EXPECT().
					GetUser(gomock.Any()).
					Times(1).
					Return(&user, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				validateErrorResponse(t, recorder, "BadRequest", "Incorrect parameters sent in request", http.StatusBadRequest)
			},
		},
		{
			name: "Wrong Token",
			body: gin.H{
				"user_name": user.UserName,
			},
			token: "wrongToken",
			buildStubs: func(dbConnector *mockdb.MockDBConnector, maker *mocktoken.MockMaker) {
				maker.
					EXPECT().
					VerifyToken("wrongToken").
					Times(1).
					Return(tokenPayload, nil)

				dbConnector.
					EXPECT().
					GetUser(gomock.Eq(user.UserName)).
					Times(1).
					Return(&user, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				validateErrorResponse(t, recorder, "Unauthorized", "User is not authorized to access this resource", http.StatusUnauthorized)
			},
		},
		{
			name: "Invalid Token",
			body: gin.H{
				"user_name": user.UserName,
			},
			token: "invalidToken",
			buildStubs: func(dbConnector *mockdb.MockDBConnector, maker *mocktoken.MockMaker) {
				maker.
					EXPECT().
					VerifyToken("invalidToken").
					Times(1).
					Return(nil, fmt.Errorf("Invalid token"))

				dbConnector.
					EXPECT().
					GetUser(gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				validateErrorResponse(t, recorder, "Unauthorized", "User is not authorized to access this resource", http.StatusUnauthorized)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			dbConnector := mockdb.NewMockDBConnector(ctrl)
			maker := mocktoken.NewMockMaker(ctrl)
			tc.buildStubs(dbConnector, maker)

			server := NewTestServer(t, dbConnector, maker)
			recorder := httptest.NewRecorder()

			url := "/v1/user/"
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			request.Header.Set("Authorization", bearerStr+tc.token)

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
	adminUser := db.User{
		FullName:   "Admin",
		Phone:      "91234567",
		UserName:   "adminuser",
		Password:   "secretadmin",
		LoginToken: "tokenadmin",
		Admin:      true,
	}

	uuidAdminToken, err := uuid.NewRandom()
	require.NoError(t, err)

	now := time.Now()

	issuedAt := now
	expiredAt := now.Add(time.Hour)

	adminTokenPayload := &token.Payload{
		ID:        uuidAdminToken,
		Username:  adminUser.UserName,
		IssuedAt:  issuedAt,
		ExpiredAt: expiredAt,
	}

	nonAdminUser := db.User{
		FullName:   "Non Admin",
		Phone:      "91234568",
		UserName:   "nonadminuser",
		Password:   "secretnonadmin",
		LoginToken: "tokennonadmin",
		Admin:      false,
	}

	uuidNonAdminToken, err := uuid.NewRandom()
	require.NoError(t, err)

	nonAdminTokenPayload := &token.Payload{
		ID:        uuidNonAdminToken,
		Username:  nonAdminUser.UserName,
		IssuedAt:  issuedAt,
		ExpiredAt: expiredAt,
	}

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
		token         string
		buildStubs    func(dbConnector *mockdb.MockDBConnector, maker *mocktoken.MockMaker)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"page":   0,
				"offset": 2,
			},
			token: adminUser.LoginToken,
			buildStubs: func(dbConnector *mockdb.MockDBConnector, maker *mocktoken.MockMaker) {
				maker.
					EXPECT().
					VerifyToken(gomock.Eq(adminUser.LoginToken)).
					Times(1).
					Return(adminTokenPayload, nil)

				dbConnector.
					EXPECT().
					GetUser(gomock.Eq(adminUser.UserName)).
					Times(1).
					Return(&adminUser, nil)

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
		{
			name: "Non Admin User Token",
			body: gin.H{
				"page":   0,
				"offset": 2,
			},
			token: nonAdminUser.LoginToken,
			buildStubs: func(dbConnector *mockdb.MockDBConnector, maker *mocktoken.MockMaker) {
				maker.
					EXPECT().
					VerifyToken(gomock.Eq(nonAdminUser.LoginToken)).
					Times(1).
					Return(nonAdminTokenPayload, nil)

				dbConnector.
					EXPECT().
					GetUser(gomock.Eq(nonAdminUser.UserName)).
					Times(1).
					Return(&nonAdminUser, nil)

				dbConnector.
					EXPECT().
					GetUsers(gomock.Any).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				validateErrorResponse(t, recorder, "Forbidden", "User is not allowed to access this resource", http.StatusForbidden)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			dbConnector := mockdb.NewMockDBConnector(ctrl)
			maker := mocktoken.NewMockMaker(ctrl)
			tc.buildStubs(dbConnector, maker)

			server := NewTestServer(t, dbConnector, maker)
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

			request.Header.Set("Authorization", bearerStr+tc.token)

			server.Router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func TestLoginUser(t *testing.T) {
	user := db.User{
		FullName:   "Test User",
		Phone:      "99989992",
		UserName:   "testuser123",
		Password:   "secret",
		LoginToken: "token",
	}

	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(dbConnector *mockdb.MockDBConnector, maker *mocktoken.MockMaker)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"user_name": user.UserName,
				"password":  user.Password,
			},
			buildStubs: func(dbConnector *mockdb.MockDBConnector, maker *mocktoken.MockMaker) {
				maker.
					EXPECT().
					CreateToken(gomock.Eq(user.UserName), gomock.Any()).
					Times(1).
					Return(user.LoginToken, nil)

				dbConnector.
					EXPECT().
					GetUser(gomock.Eq(user.UserName)).
					Times(1).
					Return(&user, nil)

				updateArgs := db.UpdateUserParams{
					ID:         0,
					FullName:   user.FullName,
					Phone:      user.Phone,
					Password:   user.Password,
					LoginToken: user.LoginToken,
				}

				dbConnector.
					EXPECT().
					UpdateUser(updateArgs).
					Times(1).
					Return(nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
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
				require.Equal(t, user.LoginToken, token)
			},
		},
		{
			name: "BadRequest",
			body: gin.H{
				"user_name": user.UserName,
			},
			buildStubs: func(dbConnector *mockdb.MockDBConnector, maker *mocktoken.MockMaker) {
				maker.
					EXPECT().
					CreateToken(gomock.Any(), gomock.Any()).
					Times(0)

				dbConnector.
					EXPECT().
					GetUser(gomock.Any).
					Times(0)

				dbConnector.
					EXPECT().
					UpdateUser(gomock.Any).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				validateErrorResponse(t, recorder, "BadRequest", "Incorrect parameters sent in request", http.StatusBadRequest)
			},
		},
		{
			name: "WrongPassword",
			body: gin.H{
				"user_name": user.UserName,
				"password":  "wrongpassword",
			},
			buildStubs: func(dbConnector *mockdb.MockDBConnector, maker *mocktoken.MockMaker) {
				maker.
					EXPECT().
					CreateToken(gomock.Any(), gomock.Any()).
					Times(0)

				dbConnector.
					EXPECT().
					GetUser(gomock.Eq(user.UserName)).
					Times(1).
					Return(&user, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
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
			maker := mocktoken.NewMockMaker(ctrl)
			tc.buildStubs(dbConnector, maker)

			server := NewTestServer(t, dbConnector, maker)
			recorder := httptest.NewRecorder()

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/v1/user/login"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.Router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func TestUpdateUser(t *testing.T) {
	user := db.User{
		FullName:   "Test user",
		Phone:      "99989992",
		UserName:   "testuser123",
		Password:   "secret",
		LoginToken: "token",
	}
	user.ID = 0

	uuidToken, err := uuid.NewRandom()
	require.NoError(t, err)

	now := time.Now()

	issuedAt := now
	expiredAt := now.Add(time.Hour)

	tokenPayload := &token.Payload{
		ID:        uuidToken,
		Username:  user.UserName,
		IssuedAt:  issuedAt,
		ExpiredAt: expiredAt,
	}

	testCases := []struct {
		name          string
		body          gin.H
		token         string
		buildStubs    func(dbConnector *mockdb.MockDBConnector, maker *mocktoken.MockMaker)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"full_name": "Test User",
				"phone":     "99989993",
			},
			token: user.LoginToken,
			buildStubs: func(dbConnector *mockdb.MockDBConnector, maker *mocktoken.MockMaker) {
				maker.
					EXPECT().
					VerifyToken(user.LoginToken).
					Times(1).
					Return(tokenPayload, nil)

				dbConnector.
					EXPECT().
					GetUser(gomock.Eq(user.UserName)).
					Times(1).
					Return(&user, nil)

				arg := db.UpdateUserParams{
					ID:         user.ID,
					FullName:   "Test User",
					Phone:      "99989993",
					Password:   user.Password,
					LoginToken: user.LoginToken,
				}

				dbConnector.
					EXPECT().
					UpdateUser(arg).
					Times(1).
					Return(nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNoContent, recorder.Code)
			},
		},
		{
			name: "Wrong Token",
			body: gin.H{
				"full_name": "Test User",
				"phone":     "99989993",
			},
			token: "wrongToken",
			buildStubs: func(dbConnector *mockdb.MockDBConnector, maker *mocktoken.MockMaker) {
				maker.
					EXPECT().
					VerifyToken("wrongToken").
					Times(1).
					Return(tokenPayload, nil)

				dbConnector.
					EXPECT().
					GetUser(gomock.Eq(user.UserName)).
					Times(1).
					Return(&user, nil)

				dbConnector.
					EXPECT().
					UpdateUser(gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				validateErrorResponse(t, recorder, "Unauthorized", "User is not authorized to access this resource", http.StatusUnauthorized)
			},
		},
		{
			name: "Invalid Token",
			body: gin.H{
				"full_name": "Test User",
				"phone":     "99989993",
			},
			token: "invalidToken",
			buildStubs: func(dbConnector *mockdb.MockDBConnector, maker *mocktoken.MockMaker) {
				maker.
					EXPECT().
					VerifyToken("invalidToken").
					Times(1).
					Return(nil, fmt.Errorf("Invalid token"))

				dbConnector.
					EXPECT().
					GetUser(gomock.Any()).
					Times(0)

				dbConnector.
					EXPECT().
					UpdateUser(gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				validateErrorResponse(t, recorder, "Unauthorized", "User is not authorized to access this resource", http.StatusUnauthorized)
			},
		},
		{
			name: "Bad Phone Number",
			body: gin.H{
				"full_name": user.FullName,
				"phone":     "9999999",
				"user_name": user.UserName,
				"password":  user.Password,
			},
			token: user.LoginToken,
			buildStubs: func(dbConnector *mockdb.MockDBConnector, maker *mocktoken.MockMaker) {
				maker.
					EXPECT().
					VerifyToken(user.LoginToken).
					Times(1).
					Return(tokenPayload, nil)

				dbConnector.
					EXPECT().
					GetUser(gomock.Eq(user.UserName)).
					Times(1).
					Return(&user, nil)

				dbConnector.
					EXPECT().
					UpdateUser(gomock.Any()).
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
			maker := mocktoken.NewMockMaker(ctrl)
			tc.buildStubs(dbConnector, maker)

			server := NewTestServer(t, dbConnector, maker)
			recorder := httptest.NewRecorder()

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/v1/user/"
			request, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(data))
			require.NoError(t, err)

			request.Header.Set("Authorization", bearerStr+tc.token)

			server.Router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}
