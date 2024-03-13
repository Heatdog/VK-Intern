package auth_test

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"testing"
	"time"

	"github.com/Heater_dog/Vk_Intern/internal/auth"
	"github.com/Heater_dog/Vk_Intern/internal/auth/mocks"
	"github.com/stretchr/testify/assert"
)

func TestTokenService_GenerateToken(t *testing.T) {
	type mockTokens func(tokenFileds auth.TokenFileds, key string) (accessToken, refreshToken string)
	testTable := []struct {
		name        string
		context     context.Context
		tokenFields auth.TokenFileds

		mockTokens    mockTokens
		expectedError error
	}{
		{
			name:    "OK",
			context: context.Background(),
			tokenFields: auth.TokenFileds{
				ID:   "123",
				Role: "User",
			},

			mockTokens: func(tokenFields auth.TokenFileds, key string) (accessToken string, refreshToken string) {
				generateAccessToken, _ := auth.GenerateToken(tokenFields, key)
				generateRefreshToken, _ := auth.GenerateRefreshToken()
				return generateAccessToken, generateRefreshToken
			},
			expectedError: nil,
		},
		{
			name:    "set in storage failed",
			context: context.Background(),
			tokenFields: auth.TokenFileds{
				ID:   "123",
				Role: "User",
			},
			mockTokens: func(tokenFileds auth.TokenFileds, key string) (accessToken string, refreshToken string) {
				return "", ""
			},
			expectedError: fmt.Errorf("set error"),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			logger := slog.New(slog.NewJSONHandler(io.Discard, nil))
			key := "123"
			expire := time.Duration(30) * time.Hour * 24
			tokenStorage := mocks.NewTokenStorage(t)

			generateAccessToken, _ := auth.GenerateToken(testCase.tokenFields, key)
			generateRefreshToken, _ := auth.GenerateRefreshToken()

			tokenStorage.On("SetToken", testCase.context, testCase.tokenFields.ID,
				generateRefreshToken, expire).
				Return(testCase.expectedError)

			tokenService := auth.NewTokenService(logger, tokenStorage, key, 30)
			accessToken, refreshToken, _, err := tokenService.GenerateToken(testCase.context, testCase.tokenFields)
			if err != nil {
				generateAccessToken = ""
				generateRefreshToken = ""
			}
			if !assert.Equal(t, generateAccessToken, accessToken) {
				t.Errorf("access token test failed. Expected: %s, got %s", generateAccessToken, accessToken)
			}
			if !assert.Equal(t, generateRefreshToken, refreshToken) {
				t.Errorf("refresh token test failed. Expected: %s, got %s", generateRefreshToken, refreshToken)
			}
			if !assert.Equal(t, testCase.expectedError, err) {
				t.Errorf("error test failed. Expected: %s, got %s", testCase.expectedError, err)
			}
		})
	}
}
