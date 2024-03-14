package dbActor

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"testing"

	"github.com/Heater_dog/Vk_Intern/internal/actor"
	"github.com/google/uuid"
	"github.com/pashagolub/pgxmock/v2"
	"github.com/stretchr/testify/assert"
)

func TestAddActor(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	type mockBehavior func(id uuid.UUID, in actor.ActorInsert, err error)
	testTable := []struct {
		name string

		context     context.Context
		actorInsert actor.ActorInsert

		expectedId    uuid.UUID
		expectedError error

		mockBehavior mockBehavior
	}{
		{
			name:    "ok",
			context: context.Background(),
			actorInsert: actor.ActorInsert{
				Name:      "John Smith",
				Gender:    "Male",
				BirthDate: "2002-12-02",
			},
			mockBehavior: func(id uuid.UUID, in actor.ActorInsert, err error) {
				rows := pgxmock.NewRows([]string{"id"}).AddRow(id)

				mock.ExpectQuery("INSERT INTO actors").
					WithArgs(in.Name, in.Gender, in.BirthDate).
					WillReturnRows(rows)
			},
			expectedId:    uuid.New(),
			expectedError: nil,
		},
		{
			name:    "sql error",
			context: context.Background(),
			actorInsert: actor.ActorInsert{
				Name:      "John Smith",
				Gender:    "Male",
				BirthDate: "2002-12-02",
			},
			mockBehavior: func(id uuid.UUID, in actor.ActorInsert, err error) {

				mock.ExpectQuery("INSERT INTO actors").
					WithArgs(in.Name, in.Gender, in.BirthDate).
					WillReturnError(err)
			},
			expectedId:    uuid.Nil,
			expectedError: fmt.Errorf("sql error"),
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior(testCase.expectedId, testCase.actorInsert, testCase.expectedError)
			logger := slog.New(slog.NewJSONHandler(io.Discard, nil))

			repo := NewActorPostgreRepository(mock, logger)

			id, err := repo.AddActor(testCase.context, testCase.actorInsert)
			if !assert.Equal(t, testCase.expectedId.String(), id) {
				t.Errorf("actor insert test failed. Expected: %s, got %s", testCase.expectedId, id)
			}
			if !assert.Equal(t, testCase.expectedError, err) {
				t.Errorf("error test failed. Expected: %s, got %s", testCase.expectedError, err)
			}
		})
	}
}
