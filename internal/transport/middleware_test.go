package transport_test

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	authMock "github.com/Heater_dog/Vk_Intern/internal/auth/mocks"
	"github.com/Heater_dog/Vk_Intern/internal/transport"
	"github.com/stretchr/testify/assert"
)

func TestMiddleware(t *testing.T) {
	type mockBehavior func(srvice *authMock.TokenService, ctx context.Context, token string)
	type formRequest func(r *http.Request, refreshToken string) *http.Request
	testTable := []struct {
		name string

		context context.Context

		refreshToken string
		formRequest  formRequest
		mockBehavior mockBehavior

		expectedHttpStatus int
	}{
		{
			name: "Access token exists",

			context: context.Background(),

			refreshToken: "",
			formRequest: func(r *http.Request, refreshToken string) *http.Request {
				r.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTMwMjI2MDAsInJvbGUiOiJBZG1pbiIsInN1YiI6IjEwNDU5NmNjLTY4Y2QtNDA0Yi04YzgyLWFkNzM2N2U3ZGU1MiJ9.ntBFmQrIVerQzbEuTTxX5qnWBFB2e69jLyE1TW8G1F4")
				return r
			},
			mockBehavior: func(ts *authMock.TokenService, ctx context.Context, token string) {},

			expectedHttpStatus: http.StatusOK,
		},
		{
			name: "Empty token header",

			context: context.Background(),

			refreshToken: "",
			formRequest: func(r *http.Request, refreshToken string) *http.Request {
				r.Header.Set("Authorization", "")
				return r
			},
			mockBehavior: func(ts *authMock.TokenService, ctx context.Context, token string) {},

			expectedHttpStatus: http.StatusUnauthorized,
		},
		{
			name: "Wrong sheame",

			context: context.Background(),

			refreshToken: "",
			formRequest: func(r *http.Request, refreshToken string) *http.Request {
				r.Header.Set("Authorization", "123 456")
				return r
			},
			mockBehavior: func(srvice *authMock.TokenService, ctx context.Context, token string) {},

			expectedHttpStatus: http.StatusUnauthorized,
		},
		{
			name: "Wrong token",

			context: context.Background(),

			refreshToken: "",
			formRequest: func(r *http.Request, refreshToken string) *http.Request {
				r.Header.Set("Authorization", "Bearer 456")
				return r
			},
			mockBehavior: func(srvice *authMock.TokenService, ctx context.Context, token string) {},

			expectedHttpStatus: http.StatusUnauthorized,
		},
		{
			name: "no tokens",

			context: context.Background(),

			refreshToken: "",
			formRequest: func(r *http.Request, refreshToken string) *http.Request {
				return r
			},
			mockBehavior: func(srvice *authMock.TokenService, ctx context.Context, token string) {},

			expectedHttpStatus: http.StatusUnauthorized,
		},
		{
			name: "auth service error",

			context: context.Background(),

			refreshToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTMwMjI2MDAsInJvbGUiOiJBZG1pbiIsInN1YiI6IjEwNDU5NmNjLTY4Y2QtNDA0Yi04YzgyLWFkNzM2N2U3ZGU1MiJ9.ntBFmQrIVerQzbEuTTxX5qnWBFB2e69jLyE1TW8G1F4",
			formRequest: func(r *http.Request, refreshToken string) *http.Request {
				r.AddCookie(&http.Cookie{
					Name:    "token",
					Value:   refreshToken,
					Expires: time.Now().Add(30 * time.Hour),
					Secure:  true,
				})
				return r
			},
			mockBehavior: func(ts *authMock.TokenService, ctx context.Context, token string) {
				ts.On("VerifyToken", ctx, token).Return("", "", time.Time{}, fmt.Errorf("auth service error"))
			},

			expectedHttpStatus: http.StatusUnauthorized,
		},
		{
			name: "have refresh token",

			context: context.Background(),

			refreshToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTMwMjI2MDAsInJvbGUiOiJBZG1pbiIsInN1YiI6IjEwNDU5NmNjLTY4Y2QtNDA0Yi04YzgyLWFkNzM2N2U3ZGU1MiJ9.ntBFmQrIVerQzbEuTTxX5qnWBFB2e69jLyE1TW8G1F4",
			formRequest: func(r *http.Request, refreshToken string) *http.Request {
				r.AddCookie(&http.Cookie{
					Name:    "token",
					Value:   refreshToken,
					Expires: time.Now().Add(30 * time.Hour),
					Secure:  true,
				})
				return r
			},
			mockBehavior: func(ts *authMock.TokenService, ctx context.Context, token string) {
				ts.On("VerifyToken", ctx, token).Return("123", "345", time.Now().Add(30*time.Hour), nil)
			},

			expectedHttpStatus: http.StatusOK,
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			logger := slog.New(slog.NewJSONHandler(io.Discard, nil))
			key := "123"

			authService := authMock.NewTokenService(t)
			testCase.mockBehavior(authService, testCase.context, testCase.refreshToken)

			mid := transport.NewMiddleware(logger, authService, key)

			nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})
			handlerToTest := mid.Auth(nextHandler)

			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/test", nil)

			r = testCase.formRequest(r, testCase.refreshToken)

			handlerToTest.ServeHTTP(w, r)

			code := w.Result().StatusCode

			if !assert.Equal(t, testCase.expectedHttpStatus, code) {
				t.Errorf("status code test case failed. Expected: %d, Got: %d", testCase.expectedHttpStatus, code)
			}
		})
	}
}

func TestAdminAuth(t *testing.T) {
	type formRequest func(r *http.Request) *http.Request
	testTable := []struct {
		name string

		formRequest formRequest

		expectedHttpStatus int
	}{
		{
			name: "Access token exists",

			formRequest: func(r *http.Request) *http.Request {
				r.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTMwMjI2MDAsInJvbGUiOiJBZG1pbiIsInN1YiI6IjEwNDU5NmNjLTY4Y2QtNDA0Yi04YzgyLWFkNzM2N2U3ZGU1MiJ9.ntBFmQrIVerQzbEuTTxX5qnWBFB2e69jLyE1TW8G1F4")
				return r
			},

			expectedHttpStatus: http.StatusOK,
		},
		{
			name: "Empty token header",

			formRequest: func(r *http.Request) *http.Request {
				r.Header.Set("Authorization", "")
				return r
			},

			expectedHttpStatus: http.StatusUnauthorized,
		},
		{
			name: "Wrong sheame",

			formRequest: func(r *http.Request) *http.Request {
				r.Header.Set("Authorization", "123 456")
				return r
			},

			expectedHttpStatus: http.StatusUnauthorized,
		},
		{
			name: "Permission denied",

			formRequest: func(r *http.Request) *http.Request {
				r.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTMwMjI2MDAsInJvbGUiOiJVc2VyIiwic3ViIjoiMTA0NTk2Y2MtNjhjZC00MDRiLThjODItYWQ3MzY3ZTdkZTUyIn0.fsPfxh3YU8W_3-sL7nkSgzLBLw5eoksQch4AslonpFI")
				return r
			},

			expectedHttpStatus: http.StatusForbidden,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			logger := slog.New(slog.NewJSONHandler(io.Discard, nil))
			key := "123"

			authService := authMock.NewTokenService(t)

			mid := transport.NewMiddleware(logger, authService, key)

			nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})
			handlerToTest := mid.AdminAuth(nextHandler)

			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/test", nil)

			r = testCase.formRequest(r)

			handlerToTest.ServeHTTP(w, r)

			code := w.Result().StatusCode

			if !assert.Equal(t, testCase.expectedHttpStatus, code) {
				t.Errorf("status code test case failed. Expected: %d, Got: %d", testCase.expectedHttpStatus, code)
			}
		})
	}
}
