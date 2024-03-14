package transport_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/Heater_dog/Vk_Intern/internal/transport"
	"github.com/Heater_dog/Vk_Intern/internal/user"
	"github.com/Heater_dog/Vk_Intern/internal/user/mocks"
	"github.com/stretchr/testify/assert"
)

func TestSignInHandler(t *testing.T) {
	type serviceArgs struct {
		context              context.Context
		reqUser              user.UserLogin
		expectedAccessToken  string
		expectedRefreshToken string
		expectedError        error
	}
	type mockUserService func(args *serviceArgs) *mocks.UserService
	testingTables := []struct {
		name                 string
		userService          mockUserService
		context              context.Context
		requestBody          string
		expectedStatusCode   int
		expectedAccessToken  string
		expectedRefreshToken string
		expectedMessage      string
		expectedError        error

		method string
	}{
		{
			name: "OK",
			userService: func(args *serviceArgs) *mocks.UserService {
				userService := mocks.NewUserService(t)
				userService.On("SignIn", args.context, args.reqUser).
					Return(args.expectedAccessToken, args.expectedRefreshToken,
						time.Time{}, args.expectedError)
				return userService
			},
			context:              context.Background(),
			requestBody:          `{"login":"Admin", "password":"Admin"}`,
			expectedStatusCode:   http.StatusOK,
			expectedAccessToken:  "123",
			expectedRefreshToken: "456",
			expectedMessage:      "",
			expectedError:        nil,

			method: "POST",
		},
		{
			name:                 "Bad Request scheame",
			userService:          nil,
			context:              context.Background(),
			requestBody:          `{"123":"56"}`,
			expectedStatusCode:   http.StatusBadRequest,
			expectedAccessToken:  "",
			expectedRefreshToken: "",
			expectedMessage:      "login: non zero value required;password: non zero value required",
			expectedError:        fmt.Errorf("login: non zero value required;password: non zero value required"),

			method: "POST",
		},
		{
			name:                 "Empty Request body",
			userService:          nil,
			context:              context.Background(),
			requestBody:          ``,
			expectedStatusCode:   http.StatusBadRequest,
			expectedAccessToken:  "",
			expectedRefreshToken: "",
			expectedMessage:      "unexpected end of JSON input",
			expectedError:        fmt.Errorf("unexpected end of JSON input"),

			method: "POST",
		},
		{
			name: "Service error",
			userService: func(args *serviceArgs) *mocks.UserService {
				userService := mocks.NewUserService(t)
				userService.On("SignIn", args.context, args.reqUser).
					Return(args.expectedAccessToken, args.expectedRefreshToken,
						time.Time{}, args.expectedError)
				return userService
			},
			context:              context.Background(),
			requestBody:          `{"login":"Admin", "password":"Admin"}`,
			expectedStatusCode:   http.StatusInternalServerError,
			expectedAccessToken:  "",
			expectedRefreshToken: "",
			expectedMessage:      "service error",
			expectedError:        fmt.Errorf("service error"),

			method: "POST",
		},
		{
			name:   "Routin failed",
			method: "GET",

			userService: nil,

			context:            context.Background(),
			expectedStatusCode: http.StatusNotFound,
		},
	}

	for _, el := range testingTables {
		t.Run(el.name, func(t *testing.T) {
			reqUser := user.UserLogin{}
			json.Unmarshal([]byte(el.requestBody), &reqUser)

			var userService *mocks.UserService
			if el.userService == nil {
				userService = nil
			} else {
				userService = el.userService(&serviceArgs{
					context:              el.context,
					reqUser:              reqUser,
					expectedAccessToken:  el.expectedAccessToken,
					expectedRefreshToken: el.expectedRefreshToken,
					expectedError:        el.expectedError,
				})
			}

			logger := slog.New(slog.NewJSONHandler(io.Discard, nil))

			authHandler := transport.NewAuthHandler(logger, userService)

			r := httptest.NewRequest(el.method, "/login", bytes.NewBufferString(el.requestBody))
			w := httptest.NewRecorder()

			authHandler.LoginRouting(w, r)
			resp := w.Result()

			header := strings.Split(resp.Header.Get("Authorization"), " ")
			accesToken := header[len(header)-1]
			var refreshToken string
			for _, cookie := range resp.Cookies() {
				if cookie.Name == "token" {
					refreshToken = cookie.Value
					break
				}
			}
			body, _ := io.ReadAll(resp.Body)

			message := transport.RespWriter{
				Text: el.expectedMessage,
			}

			expexctedMessage, _ := json.Marshal(&message)
			if !assert.Equal(t, body, expexctedMessage) {
				t.Errorf("messgae test case failed. Expected: %s, Got: %s", el.expectedMessage, string(body))
			}
			if !assert.Equal(t, el.expectedStatusCode, resp.StatusCode) {
				t.Errorf("status code test case failed. Expected: %d, Got: %d", el.expectedStatusCode, resp.StatusCode)
			}
			if !assert.Equal(t, el.expectedAccessToken, accesToken) {
				t.Errorf("access token test case failed. Expected: %s, Got: %s", el.expectedAccessToken, accesToken)
			}
			if !assert.Equal(t, el.expectedRefreshToken, refreshToken) {
				t.Errorf("refresh token test case failed. Expected: %s, Got: %s", el.expectedRefreshToken, refreshToken)
			}

		})

	}
}
