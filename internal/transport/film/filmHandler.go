package film_transport

import (
	"io"
	"log/slog"
	"net/http"

	film_service "github.com/Heater_dog/Vk_Intern/internal/services/film"
	"github.com/Heater_dog/Vk_Intern/internal/transport"
	middleware_transport "github.com/Heater_dog/Vk_Intern/internal/transport/middleware"
)

type FilmsHandler struct {
	logger       *slog.Logger
	filmsService film_service.FilmService
	middleware   *middleware_transport.Middleware
}

func NewFilmsHandler(logger *slog.Logger, filmsService film_service.FilmService,
	mid *middleware_transport.Middleware) *FilmsHandler {
	return &FilmsHandler{
		logger:       logger,
		filmsService: filmsService,
		middleware:   mid,
	}
}

const (
	addFilm = "/film/insert"
)

func (handler *FilmsHandler) Register(router *http.ServeMux) {
	filmsHandler := http.HandlerFunc(handler.ActorsRouting)

	router.Handle(addFilm, handler.middleware.Auth(handler.middleware.AdminAuth(filmsHandler)))
}

func (handler *FilmsHandler) ActorsRouting(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == addFilm && r.Method == http.MethodPost {
		handler.AddFilm(w, r)
		return
	}

	transport.NewRespWriter(w, "not found", http.StatusNotFound, handler.logger)

}

func (handler *FilmsHandler) AddFilm(w http.ResponseWriter, r *http.Request) {
	handler.logger.Info("add film handler")

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

}
