package transport_test

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	authMock "github.com/Heater_dog/Vk_Intern/internal/auth/mocks"
	"github.com/Heater_dog/Vk_Intern/internal/transport"
)

func TestMiddleware(t *testing.T) {
	testTable := []struct {
		name string
	}{
		{
			name: "Access token exists",
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
			key := "123"

			authService := authMock.NewTokenService(t)
			mid := transport.NewMiddleware(logger, authService, key)

			nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			})
			handlerToTest := mid.Auth(nextHandler)

			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/test", nil)
			r.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTMwMjI2MDAsInJvbGUiOiJBZG1pbiIsInN1YiI6IjEwNDU5NmNjLTY4Y2QtNDA0Yi04YzgyLWFkNzM2N2U3ZGU1MiJ9.ntBFmQrIVerQzbEuTTxX5qnWBFB2e69jLyE1TW8G1F4")
			handlerToTest.ServeHTTP(w, r)
		})
	}
}
