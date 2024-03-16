package token_service_test

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"testing"
	"time"

	"github.com/Heater_dog/Vk_Intern/internal/repository/token/mocks"
	token_service "github.com/Heater_dog/Vk_Intern/internal/services/token"
	"github.com/Heater_dog/Vk_Intern/pkg/jwt"
	"github.com/stretchr/testify/assert"

	innerJwt "github.com/golang-jwt/jwt"
)

func TestTokenService_GenerateToken(t *testing.T) {
	type mockTokens func(tokenFileds jwt.TokenFileds, key string, expire time.Duration) (accessToken,
		refreshToken string)
	testTable := []struct {
		name        string
		context     context.Context
		tokenFields jwt.TokenFileds

		mockTokens    mockTokens
		expectedError error
	}{
		{
			name:    "OK",
			context: context.Background(),
			tokenFields: jwt.TokenFileds{
				ID:   "123",
				Role: "User",
			},

			mockTokens: func(tokenFields jwt.TokenFileds, key string, expire time.Duration) (accessToken string,
				refreshToken string) {
				generateAccessToken, _ := jwt.GenerateToken(tokenFields, key)
				generateRefreshToken, _ := jwt.GenerateRefreshToken(tokenFields, key, expire)
				return generateAccessToken, generateRefreshToken
			},
			expectedError: nil,
		},
		{
			name:    "set in storage failed",
			context: context.Background(),
			tokenFields: jwt.TokenFileds{
				ID: "123",
			},
			mockTokens: func(tokenFileds jwt.TokenFileds, key string, expire time.Duration) (accessToken string,
				refreshToken string) {
				return "", ""
			},
			expectedError: fmt.Errorf("set error"),
		},
		{
			name:    "generate acces token error",
			context: context.Background(),
			tokenFields: jwt.TokenFileds{
				ID:   "123",
				Role: "456",
			},
			mockTokens: func(tokenFileds jwt.TokenFileds, key string, expire time.Duration) (accessToken string, refreshToken string) {
				return "", ""
			},
			expectedError: fmt.Errorf("567"),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			logger := slog.New(slog.NewJSONHandler(io.Discard, nil))
			key := "123"
			expire := time.Duration(30) * time.Hour * 24
			tokenStorage := mocks.NewTokenRepository(t)

			generateAccessToken, _ := jwt.GenerateToken(testCase.tokenFields, key)
			generateRefreshToken, _ := jwt.GenerateRefreshToken(testCase.tokenFields, key, expire)

			tokenStorage.On("SetToken", testCase.context, testCase.tokenFields.ID,
				generateRefreshToken, expire).
				Return(testCase.expectedError)

			tokenService := token_service.NewTokenService(logger, tokenStorage, key, 30)
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

func TestTokenService_VerifyToken(t *testing.T) {
	type mockArgs struct {
		context      context.Context
		tokenStorage *mocks.TokenRepository
		refToken     string
		fields       jwt.TokenFileds
		errGet       error
		errSet       error
	}

	type mockBehavior func(mockArgs) (string, string)
	testTable := []struct {
		name         string
		context      context.Context
		refreshToken string

		expectedRefreshToken string
		errGet               error
		errSet               error
		expectedError        error

		tokenFields  jwt.TokenFileds
		mockBehavior mockBehavior
	}{
		{
			name:         "OK",
			context:      context.Background(),
			refreshToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTMwMjI2MDAsInJvbGUiOiJBZG1pbiIsInN1YiI6IjEwNDU5NmNjLTY4Y2QtNDA0Yi04YzgyLWFkNzM2N2U3ZGU1MiJ9.ntBFmQrIVerQzbEuTTxX5qnWBFB2e69jLyE1TW8G1F4",

			expectedRefreshToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTMwMjI2MDAsInJvbGUiOiJBZG1pbiIsInN1YiI6IjEwNDU5NmNjLTY4Y2QtNDA0Yi04YzgyLWFkNzM2N2U3ZGU1MiJ9.ntBFmQrIVerQzbEuTTxX5qnWBFB2e69jLyE1TW8G1F4",
			errGet:               nil,
			errSet:               nil,
			expectedError:        nil,

			tokenFields: jwt.TokenFileds{
				ID:   "104596cc-68cd-404b-8c82-ad7367e7de52",
				Role: "Admin",
			},

			mockBehavior: func(args mockArgs) (string, string) {
				args.tokenStorage.On("GetToken", args.context, args.fields.ID).
					Return(args.refToken, args.errGet)

				key := "123"
				expire := time.Duration(30) * time.Hour * 24
				generateAccessToken, _ := jwt.GenerateToken(args.fields, key)
				generateRefreshToken, _ := jwt.GenerateRefreshToken(args.fields, key, expire)

				args.tokenStorage.On("SetToken", args.context, args.fields.ID, generateRefreshToken,
					expire).Return(args.errSet)

				return generateAccessToken, generateRefreshToken
			},
		},
		{
			name:         "incorrect token",
			context:      context.Background(),
			refreshToken: "23",

			expectedRefreshToken: "45",
			errGet:               nil,
			errSet:               nil,
			expectedError:        innerJwt.NewValidationError("token contains an invalid number of segments", 1),
			tokenFields: jwt.TokenFileds{
				ID:   "123",
				Role: "64",
			},

			mockBehavior: func(ma mockArgs) (string, string) {
				return "", ""
			},
		},
		{
			name:         "storage error",
			context:      context.Background(),
			refreshToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTMwMjI2MDAsInJvbGUiOiJBZG1pbiIsInN1YiI6IjEwNDU5NmNjLTY4Y2QtNDA0Yi04YzgyLWFkNzM2N2U3ZGU1MiJ9.ntBFmQrIVerQzbEuTTxX5qnWBFB2e69jLyE1TW8G1F4",

			expectedRefreshToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTMwMjI2MDAsInJvbGUiOiJBZG1pbiIsInN1YiI6IjEwNDU5NmNjLTY4Y2QtNDA0Yi04YzgyLWFkNzM2N2U3ZGU1MiJ9.ntBFmQrIVerQzbEuTTxX5qnWBFB2e69jLyE1TW8G1F4",

			tokenFields: jwt.TokenFileds{
				ID:   "104596cc-68cd-404b-8c82-ad7367e7de52",
				Role: "Admin",
			},

			mockBehavior: func(args mockArgs) (string, string) {
				args.tokenStorage.On("GetToken", args.context, args.fields.ID).
					Return(args.refToken, args.errGet)
				return "", ""
			},

			errGet:        fmt.Errorf("storage error"),
			errSet:        nil,
			expectedError: fmt.Errorf("storage error"),
		},
		{
			name:    "not equal tokens",
			context: context.Background(),

			refreshToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTMwMjI2MDAsInJvbGUiOiJBZG1pbiIsInN1YiI6IjEwNDU5NmNjLTY4Y2QtNDA0Yi04YzgyLWFkNzM2N2U3ZGU1MiJ9.ntBFmQrIVerQzbEuTTxX5qnWBFB2e69jLyE1TW8G1F4",

			expectedRefreshToken: "123",

			tokenFields: jwt.TokenFileds{
				ID:   "104596cc-68cd-404b-8c82-ad7367e7de52",
				Role: "Admin",
			},

			mockBehavior: func(ma mockArgs) (string, string) {
				ma.tokenStorage.On("GetToken", ma.context, ma.fields.ID).
					Return("1233", ma.errGet)
				return "", ""
			},

			errGet:        nil,
			errSet:        nil,
			expectedError: fmt.Errorf("tokens are not equal"),
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			logger := slog.New(slog.NewJSONHandler(io.Discard, nil))
			tokenStorage := mocks.NewTokenRepository(t)

			generateAccessToken, generateRefreshToken := testCase.mockBehavior(mockArgs{
				context:      testCase.context,
				tokenStorage: tokenStorage,
				refToken:     testCase.refreshToken,
				fields:       testCase.tokenFields,

				errGet: testCase.errGet,
				errSet: testCase.errSet,
			})

			key := "123"
			tokenService := token_service.NewTokenService(logger, tokenStorage, key, 30)
			acToken, refToken, _, err := tokenService.VerifyToken(testCase.context, testCase.refreshToken)

			if !assert.Equal(t, generateAccessToken, acToken) {
				t.Errorf("access token test failed. Expected %s, got %s", generateAccessToken, acToken)
			}
			if !assert.Equal(t, generateRefreshToken, refToken) {
				t.Errorf("refresh token test failed. Expected %s, got %s", generateRefreshToken, refToken)
			}
			if !assert.Equal(t, testCase.expectedError, err) {
				t.Errorf("error test failed. Expected %s, got %s", testCase.expectedError, err)
			}
		})
	}
}
