package transport

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"

	"github.com/Heater_dog/Vk_Intern/internal/actor"
	"github.com/asaskevich/govalidator"
)

type ActorsHandler struct {
	logger       *slog.Logger
	actorsServce actor.ActorsService
	middleware   *Middleware
}

func NewActorsHandler(logger *slog.Logger, actorsService actor.ActorsService, mid *Middleware) *ActorsHandler {
	return &ActorsHandler{
		logger:       logger,
		actorsServce: actorsService,
		middleware:   mid,
	}
}

const (
	addActor  = "/actor"
	getActors = "/actors"
)

func (handler *ActorsHandler) Register(router *http.ServeMux) {
	actorsHandler := http.HandlerFunc(handler.ActorsRouting)

	router.Handle(addActor, handler.middleware.Auth(handler.middleware.AdminAuth(actorsHandler)))
	router.Handle(getActors, handler.middleware.Auth(actorsHandler))
}

func (handler *ActorsHandler) ActorsRouting(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == addActor && r.Method == http.MethodPost {
		handler.AddActor(w, r)
		return
	}

	if r.URL.Path == getActors && r.Method == http.MethodGet {
		handler.GetActors(w, r)
		return
	}
	NewRespWriter(w, "", http.StatusNotFound, handler.logger)

}

func (handler *ActorsHandler) AddActor(w http.ResponseWriter, r *http.Request) {
	handler.logger.Info("add actor")

	handler.logger.Debug("read request body")
	data, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		handler.logger.Error("request body reading failed", slog.Any("error", err))
		NewRespWriter(w, err.Error(), http.StatusBadRequest, handler.logger)
		return
	}
	handler.logger.Debug("Request body", slog.String("body", string(data)))

	handler.logger.Debug("unmarshaling request body")
	var actor actor.ActorInsert
	if err = json.Unmarshal(data, &actor); err != nil {
		handler.logger.Warn("request body scheme error", slog.Any("error", err))
		NewRespWriter(w, err.Error(), http.StatusBadRequest, handler.logger)
		return
	}

	handler.logger.Debug("validate actor struct")
	_, err = govalidator.ValidateStruct(actor)
	if err != nil {
		handler.logger.Warn("struct validate failed", slog.Any("error", err))
		NewRespWriter(w, err.Error(), http.StatusBadRequest, handler.logger)
		return
	}

	handler.logger.Debug("insert actor service", slog.Any("actor", actor))
	id, err := handler.actorsServce.InsertActor(r.Context(), actor)
	if err != nil {
		handler.logger.Warn("actor service error", slog.Any("error", err))
		NewRespWriter(w, err.Error(), http.StatusInternalServerError, handler.logger)
		return
	}

	NewRespWriter(w, id, http.StatusCreated, handler.logger)
	handler.logger.Info("successful actor insert", slog.Any("actor", actor))
}

func (handler *ActorsHandler) GetActors(w http.ResponseWriter, r *http.Request) {
	handler.logger.Info("get actors")

	actors, err := handler.actorsServce.GetActors(r.Context())
	if err != nil {
		handler.logger.Warn("actor service error", slog.Any("error", err))
		NewRespWriter(w, err.Error(), http.StatusInternalServerError, handler.logger)
		return
	}

	handler.logger.Debug("marshal actors")
	res, err := json.Marshal(actors)
	if err != nil {
		handler.logger.Error("actor marshaling failed", slog.Any("error", err))
		NewRespWriter(w, err.Error(), http.StatusInternalServerError, handler.logger)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err = w.Write(res); err != nil {
		handler.logger.Error("writing in body failed", slog.Any("error", err))
		NewRespWriter(w, err.Error(), http.StatusInternalServerError, handler.logger)
		return
	}
	handler.logger.Info("successful user select")
}
