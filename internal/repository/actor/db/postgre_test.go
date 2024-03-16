package actor_db

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"testing"
	"time"

	actor_model "github.com/Heater_dog/Vk_Intern/internal/models/actor"
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

	type mockBehavior func(id uuid.UUID, in actor_model.ActorInsert, err error)
	testTable := []struct {
		name string

		context     context.Context
		actorInsert actor_model.ActorInsert

		expectedId    uuid.UUID
		expectedError error

		mockBehavior mockBehavior
	}{
		{
			name:    "ok",
			context: context.Background(),
			actorInsert: actor_model.ActorInsert{
				Name:      "John Smith",
				Gender:    "Male",
				BirthDate: "2002-12-02",
			},
			mockBehavior: func(id uuid.UUID, in actor_model.ActorInsert, err error) {
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
			actorInsert: actor_model.ActorInsert{
				Name:      "John Smith",
				Gender:    "Male",
				BirthDate: "2002-12-02",
			},
			mockBehavior: func(id uuid.UUID, in actor_model.ActorInsert, err error) {

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

func TestGetActors(t *testing.T) {
	type mockBehavior func(actors []actor_model.Actor, err error)
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	testTable := []struct {
		name string

		context      context.Context
		mockBehavior mockBehavior

		expectedError  error
		expectedActors []actor_model.Actor
	}{
		{
			name: "ok",

			context: context.Background(),
			mockBehavior: func(actors []actor_model.Actor, err error) {
				rows := pgxmock.NewRows([]string{"id", "name", "gender", "birth_date"})
				for _, el := range actors {
					rows.AddRow(el.ID, el.Name, el.Gender, el.BirthDate)
				}

				mock.ExpectQuery("SELECT id, name, gender, birth_date FROM actors").
					WillReturnRows(rows)
			},

			expectedError: nil,
			expectedActors: []actor_model.Actor{
				{
					ID:        uuid.New(),
					Name:      "John",
					Gender:    "Male",
					BirthDate: time.Now(),
				},
				{
					ID:        uuid.New(),
					Name:      "Eric",
					Gender:    "Female",
					BirthDate: time.Now(),
				},
			},
		},
		{
			name: "query error",

			context: context.Background(),
			mockBehavior: func(actors []actor_model.Actor, err error) {
				rows := pgxmock.NewRows([]string{"id", "name", "gender", "birth_date"})
				for _, el := range actors {
					rows.AddRow(el.ID, el.Name, el.Gender, el.BirthDate)
				}

				mock.ExpectQuery("SELECT id, name, gender, birth_date FROM actors").
					WillReturnError(err)
			},

			expectedError:  fmt.Errorf("query error"),
			expectedActors: nil,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			logger := slog.New(slog.NewJSONHandler(io.Discard, nil))
			testCase.mockBehavior(testCase.expectedActors, testCase.expectedError)

			repo := NewActorPostgreRepository(mock, logger)

			actors, err := repo.GetActors(testCase.context)
			if !assert.Equal(t, testCase.expectedActors, actors) {
				t.Errorf("actor get test failed. Expected: %s, got %s", testCase.expectedActors, actors)
			}
			if !assert.Equal(t, testCase.expectedError, err) {
				t.Errorf("error test failed. Expected: %s, got %s", testCase.expectedError, err)
			}
		})
	}
}
