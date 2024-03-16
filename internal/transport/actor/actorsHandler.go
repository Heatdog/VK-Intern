package actor_transport

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"

	actor_model "github.com/Heater_dog/Vk_Intern/internal/models/actor"
	_ "github.com/Heater_dog/Vk_Intern/internal/models/actor_film"
	actor_service "github.com/Heater_dog/Vk_Intern/internal/services/actor"
	"github.com/Heater_dog/Vk_Intern/internal/transport"
	middleware_transport "github.com/Heater_dog/Vk_Intern/internal/transport/middleware"
	"github.com/asaskevich/govalidator"
	"github.com/google/uuid"
)

type ActorsHandler struct {
	logger       *slog.Logger
	actorsServce actor_service.ActorsService
	middleware   *middleware_transport.Middleware
}

func NewActorsHandler(logger *slog.Logger, actorsService actor_service.ActorsService,
	mid *middleware_transport.Middleware) *ActorsHandler {
	return &ActorsHandler{
		logger:       logger,
		actorsServce: actorsService,
		middleware:   mid,
	}
}

const (
	addActor    = "/actor/insert"
	getActors   = "/actors"
	deleteActor = "/actor/delete"
	updateActor = "/actor/update"
)

func (handler *ActorsHandler) Register(router *http.ServeMux) {
	actorsHandler := http.HandlerFunc(handler.ActorsRouting)

	router.Handle(addActor, handler.middleware.Auth(handler.middleware.AdminAuth(actorsHandler)))
	router.Handle(getActors, handler.middleware.Auth(actorsHandler))
	router.Handle(deleteActor, handler.middleware.Auth(handler.middleware.AdminAuth(actorsHandler)))
	router.Handle(updateActor, handler.middleware.Auth(handler.middleware.AdminAuth(actorsHandler)))
}

func (handler *ActorsHandler) ActorsRouting(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == addActor && r.Method == http.MethodPost {
		handler.AddActor(w, r)
		return
	}

	if r.URL.Path == deleteActor && r.Method == http.MethodDelete {
		handler.DeleteActor(w, r)
		return
	}

	if r.URL.Path == getActors && r.Method == http.MethodGet {
		handler.GetActors(w, r)
		return
	}

	if r.URL.Path == updateActor && r.Method == http.MethodPut {
		handler.UpdateActor(w, r)
		return
	}

	transport.NewRespWriter(w, "not found", http.StatusNotFound, handler.logger)

}

// Добавление актера
// @Summary AddActor
// @Security ApiKeyAuth
// @Description Добавление актера в систему
// @Tags actors
// @ID add-actor
// @Accept json
// @Produce json
// @Param input body actor_model.ActorInsert true "actor info"
// @Success 201 {object} transport.RespWriter
// @Failure 400 {object} transport.RespWriter
// @Failure 401 {object} transport.RespWriter
// @Failure 403 {object} transport.RespWriter
// @Failure 500 {object} transport.RespWriter
// @Router /actor/insert [post]
func (handler *ActorsHandler) AddActor(w http.ResponseWriter, r *http.Request) {
	handler.logger.Info("add actor handler")

	handler.logger.Debug("read request body")
	data, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		handler.logger.Error("request body reading failed", slog.Any("error", err))
		transport.NewRespWriter(w, err.Error(), http.StatusBadRequest, handler.logger)
		return
	}
	handler.logger.Debug("Request body", slog.String("body", string(data)))

	handler.logger.Debug("unmarshaling request body")
	var actor actor_model.ActorInsert
	if err = json.Unmarshal(data, &actor); err != nil {
		handler.logger.Warn("request body scheme error", slog.Any("error", err))
		transport.NewRespWriter(w, err.Error(), http.StatusBadRequest, handler.logger)
		return
	}

	handler.logger.Debug("validate actor struct")
	_, err = govalidator.ValidateStruct(actor)
	if err != nil {
		handler.logger.Warn("struct validate failed", slog.Any("error", err))
		transport.NewRespWriter(w, err.Error(), http.StatusBadRequest, handler.logger)
		return
	}

	handler.logger.Debug("insert actor service", slog.Any("actor", actor))
	id, err := handler.actorsServce.InsertActor(r.Context(), actor)
	if err != nil {
		handler.logger.Warn("actor service error", slog.Any("error", err))
		transport.NewRespWriter(w, err.Error(), http.StatusInternalServerError, handler.logger)
		return
	}

	transport.NewRespWriter(w, id, http.StatusCreated, handler.logger)
	handler.logger.Info("successful actor insert", slog.Any("actor", actor))
}

// Получение списка актеров
// @Summary GetActors
// @Security ApiKeyAuth
// @Description Получение актеров из системы. Вместе с актерами выводится список фильмов.
// @Tags actors
// @ID get-actors
// @Accept json
// @Produce json
// @Success 200 {object} []actorfilm.ActorFilms
// @Failure 400 {object} transport.RespWriter
// @Failure 401 {object} transport.RespWriter
// @Failure 403 {object} transport.RespWriter
// @Failure 500 {object} transport.RespWriter
// @Router /actors [get]
func (handler *ActorsHandler) GetActors(w http.ResponseWriter, r *http.Request) {
	handler.logger.Info("get actors handler")

	actors, err := handler.actorsServce.GetActors(r.Context())
	if err != nil {
		handler.logger.Warn("actor service error", slog.Any("error", err))
		transport.NewRespWriter(w, err.Error(), http.StatusInternalServerError, handler.logger)
		return
	}

	handler.logger.Debug("marshal actors")
	res, err := json.Marshal(actors)
	if err != nil {
		handler.logger.Error("actor marshaling failed", slog.Any("error", err))
		transport.NewRespWriter(w, err.Error(), http.StatusInternalServerError, handler.logger)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err = w.Write(res); err != nil {
		handler.logger.Error("writing in body failed", slog.Any("error", err))
		transport.NewRespWriter(w, err.Error(), http.StatusInternalServerError, handler.logger)
		return
	}
	handler.logger.Info("successful user select")
}

type ActorId struct {
	ID uuid.UUID `json:"id"`
}

// удаление актера
// @Summary DeleteActor
// @Security ApiKeyAuth
// @Description Удаление актера из системы
// @Tags actors
// @ID delete-actor
// @Accept json
// @Produce json
// @Param input body ActorId true "actor id"
// @Success 200 {object} transport.RespWriter
// @Failure 400 {object} transport.RespWriter
// @Failure 401 {object} transport.RespWriter
// @Failure 403 {object} transport.RespWriter
// @Failure 500 {object} transport.RespWriter
// @Router /actor/delete [delete]
func (handler *ActorsHandler) DeleteActor(w http.ResponseWriter, r *http.Request) {
	handler.logger.Info("delete actor handler")

	handler.logger.Debug("read request body")
	data, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		handler.logger.Error("request body reading failed", slog.Any("error", err))
		transport.NewRespWriter(w, err.Error(), http.StatusBadRequest, handler.logger)
		return
	}

	handler.logger.Debug("unmasrhaling data", slog.String("data", string(data)))
	reqId := ActorId{}
	if err = json.Unmarshal(data, &reqId); err != nil {
		handler.logger.Warn("request body scheme error", slog.Any("error", err))
		transport.NewRespWriter(w, err.Error(), http.StatusBadRequest, handler.logger)
		return
	}

	handler.logger.Info("delete actor", slog.Any("id", reqId.ID))
	if err = handler.actorsServce.DeleteActor(r.Context(), reqId.ID); err != nil {
		handler.logger.Warn("actor serevice error")
		transport.NewRespWriter(w, "actor serevice error", http.StatusInternalServerError, handler.logger)
		return
	}
	w.WriteHeader(http.StatusOK)
	handler.logger.Info("successful user delete")
}

// Обновление информации актера
// @Summary UpdateActor
// @Security ApiKeyAuth
// @Description Обновление информации об акетре в системе
// @Tags actors
// @ID update-actor
// @Accept json
// @Produce json
// @Param input body actor_model.UpdateActor true "actor fields"
// @Success 200 {object} transport.RespWriter
// @Failure 400 {object} transport.RespWriter
// @Failure 401 {object} transport.RespWriter
// @Failure 403 {object} transport.RespWriter
// @Failure 500 {object} transport.RespWriter
// @Router /actor/update [put]
func (handler *ActorsHandler) UpdateActor(w http.ResponseWriter, r *http.Request) {
	handler.logger.Info("update actor handler")

	handler.logger.Debug("read request body")
	data, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		handler.logger.Error("request body reading failed", slog.Any("error", err))
		transport.NewRespWriter(w, err.Error(), http.StatusBadRequest, handler.logger)
		return
	}

	handler.logger.Debug("unmasrhaling data", slog.String("data", string(data)))
	updatedFields := actor_model.UpdateActor{}
	if err = json.Unmarshal(data, &updatedFields); err != nil {
		handler.logger.Warn("request body scheme error", slog.Any("error", err))
		transport.NewRespWriter(w, err.Error(), http.StatusBadRequest, handler.logger)
		return
	}
	handler.logger.Debug("validate actor struct")
	_, err = govalidator.ValidateStruct(updatedFields)
	if err != nil {
		handler.logger.Warn("struct validate failed", slog.Any("error", err))
		transport.NewRespWriter(w, err.Error(), http.StatusBadRequest, handler.logger)
		return
	}

	handler.logger.Info("update actor", slog.Any("id", updatedFields.ID))
	if err = handler.actorsServce.UpdateActor(r.Context(), updatedFields); err != nil {
		handler.logger.Warn("actor serevice error")
		transport.NewRespWriter(w, "actor serevice error", http.StatusInternalServerError, handler.logger)
		return
	}
	w.WriteHeader(http.StatusOK)
	handler.logger.Info("successful user update")
}
