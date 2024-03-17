package actor_transport_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	actor_model "github.com/Heater_dog/Vk_Intern/internal/models/actor"
	actor_mock "github.com/Heater_dog/Vk_Intern/internal/services/actor/mocks"
	token_mock "github.com/Heater_dog/Vk_Intern/internal/services/token/mocks"
	"github.com/Heater_dog/Vk_Intern/internal/transport"
	actor_transport "github.com/Heater_dog/Vk_Intern/internal/transport/actor"
	middleware_transport "github.com/Heater_dog/Vk_Intern/internal/transport/middleware"
	"github.com/stretchr/testify/assert"
)

func TestActorInsertHandler(t *testing.T) {
	type mockArgs struct {
		context     context.Context
		actorInsert actor_model.ActorInsert
		expectedId  string
		expectedErr error
	}
	type mockBehavior func(actorService *actor_mock.ActorsService, args mockArgs)
	testTable := []struct {
		name string

		method      string
		requestBody string
		context     context.Context

		mockBehavior mockBehavior

		expectedId         string
		expectedErr        error
		expectedStatusCode int
		expectedMessage    string
	}{
		{
			name: "ok",

			method:      "POST",
			requestBody: `{"name":"John", "gender":"Male", "birth_date":"2002-02-02"}`,
			context:     context.Background(),

			mockBehavior: func(actorService *actor_mock.ActorsService, ma mockArgs) {
				actorService.On("InsertActor", ma.context, ma.actorInsert).
					Return(ma.expectedId, ma.expectedErr)
			},

			expectedId:         "123",
			expectedStatusCode: http.StatusCreated,
			expectedErr:        nil,
			expectedMessage:    "123",
		},
		{
			name: "unmarshaling error",

			method:      "POST",
			requestBody: `{"12":"John", "gender":"Male", "birth_date":"2002-02-02"}`,
			context:     context.Background(),

			mockBehavior: func(actorService *actor_mock.ActorsService, args mockArgs) {},

			expectedId:         "123",
			expectedStatusCode: http.StatusBadRequest,
			expectedErr:        fmt.Errorf("name: non zero value required"),
			expectedMessage:    "name: non zero value required",
		},
		{
			name: "service error",

			method:      "POST",
			requestBody: `{"name":"John", "gender":"Male", "birth_date":"2002-02-02"}`,
			context:     context.Background(),

			mockBehavior: func(actorService *actor_mock.ActorsService, ma mockArgs) {
				actorService.On("InsertActor", ma.context, ma.actorInsert).
					Return(ma.expectedId, ma.expectedErr)
			},

			expectedId:         "123",
			expectedStatusCode: http.StatusInternalServerError,
			expectedErr:        fmt.Errorf("service error"),
			expectedMessage:    "service error",
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			logger := slog.New(slog.NewJSONHandler(io.Discard, nil))
			key := "123"
			reqActor := actor_model.ActorInsert{}
			json.Unmarshal([]byte(testCase.requestBody), &reqActor)

			authService := token_mock.NewTokenService(t)
			mid := middleware_transport.NewMiddleware(logger, authService, key)
			actorService := actor_mock.NewActorsService(t)
			testCase.mockBehavior(actorService, mockArgs{
				context:     testCase.context,
				actorInsert: reqActor,
				expectedId:  testCase.expectedId,
				expectedErr: testCase.expectedErr,
			})

			handler := actor_transport.NewActorsHandler(logger, actorService, mid)

			r := httptest.NewRequest(testCase.method, "/actor/insert", bytes.NewBufferString(testCase.requestBody))
			w := httptest.NewRecorder()
			handler.ActorsRouting(w, r)

			resp := w.Result()

			body, _ := io.ReadAll(resp.Body)

			message := transport.RespWriter{
				Text: testCase.expectedMessage,
			}
			expexctedMessage, _ := json.Marshal(&message)
			if !assert.Equal(t, expexctedMessage, body) {
				t.Errorf("messgae test case failed. Expected: %s, Got: %s", expexctedMessage, string(body))
			}
			if !assert.Equal(t, testCase.expectedStatusCode, resp.StatusCode) {
				t.Errorf("status code test case failed. Expected: %d, Got: %d",
					testCase.expectedStatusCode, resp.StatusCode)
			}
		})
	}
}
