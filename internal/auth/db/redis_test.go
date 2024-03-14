package authDb_test

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"testing"
	"time"

	authDb "github.com/Heater_dog/Vk_Intern/internal/auth/db"
	"github.com/go-redis/redismock/v9"
	"github.com/stretchr/testify/assert"
)

func TestSetToken(t *testing.T) {
	db, mock := redismock.NewClientMock()
	defer db.Close()
	testingTables := []struct {
		name    string
		context context.Context

		userId string
		token  string
		expire time.Duration

		expectedError error
	}{
		{
			name:    "OK",
			context: context.Background(),

			userId: "123",
			token:  "456",
			expire: time.Duration(time.Second * 15),

			expectedError: nil,
		},
		{
			name:    "low expire time",
			context: context.Background(),

			userId: "123",
			token:  "456",
			expire: time.Duration(30),

			expectedError: fmt.Errorf("too low expire time"),
		},
		{
			name:    "internal redis error",
			context: context.Background(),

			userId: "123",
			token:  "456",
			expire: time.Duration(time.Second * 30),

			expectedError: fmt.Errorf("internal error"),
		},
	}
	for _, testCase := range testingTables {
		t.Run(testCase.name, func(t *testing.T) {
			logger := slog.New(slog.NewJSONHandler(io.Discard, nil))

			if testCase.expectedError == nil {
				mock.ExpectSet(testCase.userId, testCase.token, testCase.expire).
					SetVal("")
			} else {
				mock.ExpectSet(testCase.userId, testCase.token, testCase.expire).
					SetErr(testCase.expectedError)
			}
			storage := authDb.NewRedisTokenStorage(logger, db)

			err := storage.SetToken(testCase.context, testCase.userId, testCase.token, testCase.expire)
			if !assert.Equal(t, testCase.expectedError, err) {
				t.Errorf("error test failed. Expected %s, got %s", testCase.expectedError, err)
			}
		})
	}

}
