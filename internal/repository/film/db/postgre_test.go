package film_db

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"testing"
	"time"

	film_model "github.com/Heater_dog/Vk_Intern/internal/models/film"
	"github.com/google/uuid"
	"github.com/pashagolub/pgxmock/v2"
	"github.com/stretchr/testify/assert"
)

func TestGetActors(t *testing.T) {
	type mockBehavior func(films []film_model.Film, actorId uuid.UUID, err error)
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	testTable := []struct {
		name string

		context      context.Context
		mockBehavior mockBehavior
		actorId      uuid.UUID

		expectedError error
		expectedFilms []film_model.Film
	}{
		{
			name: "ok",

			context: context.Background(),
			mockBehavior: func(films []film_model.Film, actorId uuid.UUID, err error) {
				rows := pgxmock.NewRows([]string{"id", "title", "description", "release_date", "rating"})
				for _, el := range films {
					rows.AddRow(el.ID, el.Title, el.Description, el.ReleaseDate, el.Rating)
				}

				mock.ExpectQuery("SELECT f.id, f.title, f.description, f.release_date, f.rating FROM films f").
					WithArgs(actorId).
					WillReturnRows(rows)
			},

			actorId:       uuid.New(),
			expectedError: nil,
			expectedFilms: []film_model.Film{
				{
					ID:          uuid.New(),
					Title:       "1323",
					Description: "231123",
					ReleaseDate: time.Now(),
					Rating:      1.2,
				},
				{
					ID:          uuid.New(),
					Title:       "123213",
					Description: "23432",
					ReleaseDate: time.Now(),
					Rating:      2.4,
				},
			},
		},

		{
			name: "query error",

			context: context.Background(),
			mockBehavior: func(films []film_model.Film, actorId uuid.UUID, err error) {
				rows := pgxmock.NewRows([]string{"id", "title", "description", "release_date", "rating"})
				for _, el := range films {
					rows.AddRow(el.ID, el.Title, el.Description, el.ReleaseDate, el.Rating)
				}

				mock.ExpectQuery("SELECT f.id, f.title, f.description, f.release_date, f.rating FROM films f").
					WithArgs(actorId).
					WillReturnError(err)
			},

			expectedError: fmt.Errorf("query error"),
			expectedFilms: nil,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			logger := slog.New(slog.NewJSONHandler(io.Discard, nil))
			testCase.mockBehavior(testCase.expectedFilms, testCase.actorId, testCase.expectedError)

			repo := NewFilmsPostgreRepository(mock, logger)

			films, err := repo.GetFilmsWithActor(testCase.context, testCase.actorId)
			if !assert.Equal(t, testCase.expectedFilms, films) {
				t.Errorf("actor get test failed.")
			}
			if !assert.Equal(t, testCase.expectedError, err) {
				t.Errorf("error test failed. Expected: %s, got %s", testCase.expectedError, err)
			}
		})
	}
}
