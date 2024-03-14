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
}

func NewActorsHandler(logger *slog.Logger, actorsService actor.ActorsService) Handler {
	return &ActorsHandler{
		logger:       logger,
		actorsServce: actorsService,
	}
}

const (
	addActor = "/actor"
)

func (handler *ActorsHandler) Register(router *http.ServeMux) {
	router.HandleFunc(addActor, handler.actorsRouting)
}

func (handler *ActorsHandler) actorsRouting(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		handler.AddActor(w, r)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
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

	var actor actor.ActorInsert
	handler.logger.Debug("unmarshaling request body")
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
	if err = handler.actorsServce.InsertActor(r.Context(), actor); err != nil {
		handler.logger.Warn("actor service error", slog.Any("error", err))
		NewRespWriter(w, err.Error(), http.StatusInternalServerError, handler.logger)
		return
	}
	NewRespWriter(w, "successfull actor insert", http.StatusCreated, handler.logger)
	handler.logger.Info("successful actor insert", slog.Any("actor", actor))
}
