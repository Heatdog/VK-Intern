package user_test

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"testing"
	"time"

	tokenMocks "github.com/Heater_dog/Vk_Intern/internal/auth/mocks"
	"github.com/Heater_dog/Vk_Intern/internal/user"
	userMocks "github.com/Heater_dog/Vk_Intern/internal/user/mocks"
	cryptohash "github.com/Heater_dog/Vk_Intern/pkg/cryptoHash"
	"github.com/Heater_dog/Vk_Intern/pkg/jwt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestSignIn(t *testing.T) {
	testingTables := []struct {
		name                  string
		context               context.Context
		userLogin             user.UserLogin
		expectedUser          *user.User
		expectedUserRepoError error

		expectedAccessToken       string
		expectedRefreshToken      string
		expectedTokenServiceError error

		pswdError error
	}{
		{
			name:    "OK",
			context: context.Background(),
			userLogin: user.UserLogin{
				Login:    "John",
				Password: "123",
			},
			expectedUser: &user.User{
				ID:       uuid.New(),
				Login:    "John",
				Role:     "User",
				Password: "",
			},
			expectedUserRepoError: nil,

			expectedAccessToken:       "123",
			expectedRefreshToken:      "456",
			expectedTokenServiceError: nil,
			pswdError:                 nil,
		},
		{
			name:    "User repo failed",
			context: context.Background(),
			userLogin: user.UserLogin{
				Login:    "John",
				Password: "123",
			},
			expectedUser:          nil,
			expectedUserRepoError: fmt.Errorf("not such user"),
			pswdError:             nil,
		},
		{
			name:    "Token service error",
			context: context.Background(),
			userLogin: user.UserLogin{
				Login:    "John",
				Password: "123",
			},
			expectedUser: &user.User{
				ID:       uuid.New(),
				Login:    "John",
				Role:     "User",
				Password: "",
			},
			expectedUserRepoError: nil,

			expectedAccessToken:       "",
			expectedRefreshToken:      "",
			expectedTokenServiceError: fmt.Errorf("token create failed"),
			pswdError:                 nil,
		},
		{
			name:    "Incorrect password",
			context: context.Background(),
			userLogin: user.UserLogin{
				Login:    "John",
				Password: "123",
			},
			expectedUser: &user.User{
				ID:       uuid.New(),
				Login:    "John",
				Role:     "User",
				Password: "456",
			},
			expectedUserRepoError: nil,
			pswdError:             fmt.Errorf("wrong password error=John"),
		},
	}
	for _, testCase := range testingTables {
		t.Run(testCase.name, func(t *testing.T) {
			logger := slog.New(slog.NewJSONHandler(io.Discard, nil))
			userRepo := userMocks.NewUserRepository(t)
			tokenService := tokenMocks.NewTokenService(t)

			userRepo.On("Find", testCase.context, testCase.userLogin.Login).
				Return(testCase.expectedUser, testCase.expectedUserRepoError)

			if testCase.expectedUserRepoError == nil && testCase.pswdError == nil {
				if testCase.expectedUser.Password == "" {
					pswd, _ := cryptohash.Hash(testCase.userLogin.Password)
					testCase.expectedUser.Password = string(pswd)
				}

				tokenService.On("GenerateToken", testCase.context, jwt.TokenFileds{
					ID:   testCase.expectedUser.ID.String(),
					Role: testCase.expectedUser.Role,
				}).Return(testCase.expectedAccessToken, testCase.expectedRefreshToken,
					time.Time{}, testCase.expectedTokenServiceError)
			}
			userService := user.NewUserService(logger, userRepo, tokenService)
			resAccessToken, resRefreshToken, _, err := userService.SignIn(testCase.context, testCase.userLogin)

			if !assert.Equal(t, testCase.expectedAccessToken, resAccessToken) {
				t.Errorf("access token test failed. Expected: %s, got %s",
					testCase.expectedAccessToken, resAccessToken)
			}
			if !assert.Equal(t, testCase.expectedRefreshToken, resRefreshToken) {
				t.Errorf("refresh token test failed. Expected: %s, got %s",
					testCase.expectedRefreshToken, resRefreshToken)
			}

			if testCase.expectedUserRepoError != nil {
				if !assert.Equal(t, testCase.expectedUserRepoError, err) {
					t.Errorf("error test failed. Expected: %s, got %s",
						testCase.expectedTokenServiceError, err)
				}
			} else {
				if testCase.pswdError != nil {
					if !assert.Equal(t, testCase.pswdError, err) {
						t.Errorf("error test failed. Expected: %s, got %s",
							testCase.pswdError, err)
					}
				} else {
					if !assert.Equal(t, testCase.expectedTokenServiceError, err) {
						t.Errorf("error test failed. Expected: %s, got %s",
							testCase.expectedTokenServiceError, err)
					}
				}
			}
		})
	}
}
