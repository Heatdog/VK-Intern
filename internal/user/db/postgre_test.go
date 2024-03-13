package userDb

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"testing"

	"github.com/Heater_dog/Vk_Intern/internal/user"
	"github.com/google/uuid"
	"github.com/pashagolub/pgxmock/v2"
	"github.com/stretchr/testify/assert"
)

func TestFind(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	type mockBehavior func(id uuid.UUID, login, password, role string, err error)
	testingTables := []struct {
		name    string
		context context.Context
		login   string

		mockBehavior mockBehavior

		expectedUser  *user.User
		expectedError error
	}{
		{
			name:    "OK",
			context: context.Background(),
			login:   "John",

			mockBehavior: func(id uuid.UUID, login, password, role string, err error) {
				rows := pgxmock.NewRows([]string{"id", "login", "password", "role"}).
					AddRow(id, login, password, role)

				mock.ExpectQuery("SELECT id, login, password, role FROM Users").
					WithArgs(login).
					WillReturnRows(rows)
			},

			expectedUser: &user.User{
				ID:       uuid.New(),
				Login:    "John",
				Password: "123",
				Role:     "User",
			},
			expectedError: nil,
		},
		{
			name:    "Scan error",
			context: context.Background(),
			login:   "",

			mockBehavior: func(id uuid.UUID, login, password, role string, err error) {
				mock.ExpectQuery("SELECT id, login, password, role FROM Users").
					WithArgs(login).
					WillReturnError(err)
			},

			expectedUser:  nil,
			expectedError: fmt.Errorf("SQL error"),
		},
	}
	for _, testCase := range testingTables {
		t.Run(testCase.name, func(t *testing.T) {
			if testCase.expectedUser != nil {
				testCase.mockBehavior(testCase.expectedUser.ID, testCase.expectedUser.Login,
					testCase.expectedUser.Password, testCase.expectedUser.Role, testCase.expectedError)
			} else {
				testCase.mockBehavior(uuid.UUID{}, "", "", "", testCase.expectedError)
			}

			logger := slog.New(slog.NewJSONHandler(io.Discard, nil))
			repo := NewUserPostgreRepository(mock, logger)

			user, err := repo.Find(testCase.context, testCase.login)

			if !assert.Equal(t, testCase.expectedUser, user) {
				t.Errorf("user select test failed. Expected: %s, got %s", testCase.expectedUser, user)
			}
			if !assert.Equal(t, testCase.expectedError, err) {
				t.Errorf("error test failed. Expected: %s, got %s", testCase.expectedError, err)
			}
		})
	}
}
