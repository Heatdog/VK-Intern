package actor_service_test

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"testing"

	actor_model "github.com/Heater_dog/Vk_Intern/internal/models/actor"
	actorMock "github.com/Heater_dog/Vk_Intern/internal/repository/actor/mocks"
	filmMock "github.com/Heater_dog/Vk_Intern/internal/repository/film/mocks"
	actor_service "github.com/Heater_dog/Vk_Intern/internal/services/actor"
	"github.com/stretchr/testify/assert"
)

func TestAddActor(t *testing.T) {
	testTable := []struct {
		name string

		context     context.Context
		actorInsert actor_model.ActorInsert

		expectedId    string
		expectedError error
	}{
		{
			name: "ok",

			context: context.Background(),
			actorInsert: actor_model.ActorInsert{
				Name:      "John",
				Gender:    "Male",
				BirthDate: "2001-05-02",
			},

			expectedId:    "123",
			expectedError: nil,
		},
		{
			name: "err",

			context:     context.Background(),
			actorInsert: actor_model.ActorInsert{},

			expectedId:    "",
			expectedError: fmt.Errorf("storage err"),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			logger := slog.New(slog.NewJSONHandler(io.Discard, nil))

			filmRepo := filmMock.NewFilmsRepository(t)

			actorRepo := actorMock.NewActorsRepository(t)
			actorRepo.On("AddActor", testCase.context, testCase.actorInsert).
				Return(testCase.expectedId, testCase.expectedError)

			actorService := actor_service.NewActorsService(logger, actorRepo, filmRepo)

			id, err := actorService.InsertActor(testCase.context, testCase.actorInsert)
			if !assert.Equal(t, testCase.expectedId, id) {
				t.Errorf("actor insert test failed. Expected: %s, got %s", testCase.expectedId, id)
			}
			if !assert.Equal(t, testCase.expectedError, err) {
				t.Errorf("error test failed. Expected: %s, got %s", testCase.expectedError, err)
			}
		})
	}
}

func TestGetActors(t *testing.T) {

}
