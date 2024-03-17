package film_transport

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"

	actorfilm "github.com/Heater_dog/Vk_Intern/internal/models/actor_film"
	film_model "github.com/Heater_dog/Vk_Intern/internal/models/film"
	actor_service "github.com/Heater_dog/Vk_Intern/internal/services/actor"
	film_service "github.com/Heater_dog/Vk_Intern/internal/services/film"
	"github.com/Heater_dog/Vk_Intern/internal/transport"
	middleware_transport "github.com/Heater_dog/Vk_Intern/internal/transport/middleware"
	"github.com/asaskevich/govalidator"
)

type FilmsHandler struct {
	logger       *slog.Logger
	filmsService film_service.FilmService
	actorService actor_service.ActorsService
	middleware   *middleware_transport.Middleware
}

func NewFilmsHandler(logger *slog.Logger, filmsService film_service.FilmService,
	actorService actor_service.ActorsService, mid *middleware_transport.Middleware) *FilmsHandler {
	return &FilmsHandler{
		logger:       logger,
		filmsService: filmsService,
		actorService: actorService,
		middleware:   mid,
	}
}

const (
	addFilm    = "/film/insert"
	updateFilm = "/film/update"
	deleteFilm = "/film/delete"
	getFilms   = "/films"

	getFilmsSearch = "/films/search"
)

func (handler *FilmsHandler) Register(router *http.ServeMux) {
	filmsHandler := http.HandlerFunc(handler.FilmsRouting)

	router.Handle(addFilm, handler.middleware.Auth(handler.middleware.AdminAuth(filmsHandler)))
	router.Handle(updateFilm, handler.middleware.Auth(handler.middleware.AdminAuth(filmsHandler)))
	router.Handle(deleteFilm, handler.middleware.Auth(handler.middleware.AdminAuth(filmsHandler)))
	router.Handle(getFilms, handler.middleware.Auth(filmsHandler))
	router.Handle(getFilmsSearch, handler.middleware.Auth(filmsHandler))
}

func (handler *FilmsHandler) FilmsRouting(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == addFilm && r.Method == http.MethodPost {
		handler.InsertFilm(w, r)
		return
	}

	if r.URL.Path == updateFilm && r.Method == http.MethodPut {
		handler.UpdateFilm(w, r)
		return
	}

	if r.URL.Path == deleteFilm && r.Method == http.MethodDelete {
		handler.DeleteFilm(w, r)
		return
	}

	if r.URL.Path == getFilms && r.Method == http.MethodGet {
		handler.GetFilms(w, r)
		return
	}

	if r.URL.Path == getFilmsSearch && r.Method == http.MethodGet {
		handler.GetFilmsSearch(w, r)
		return
	}

	transport.NewRespWriter(w, "not found", http.StatusNotFound, handler.logger)

}

// Добавление фильма
// @Summary InsertFilm
// @Security ApiKeyAuth
// @Description Добавление фильма в систему. Все параметры являются обязтельными, кроме списка актеров.
// @Tags films
// @ID add-film
// @Accept json
// @Produce json
// @Param input body film_model.FilmInsert true "actor fields"
// @Success 201 {object} transport.RespWriter ID созданного фильма
// @Failure 400 {object} transport.RespWriter Некорректный ввод данных
// @Failure 401 {object} transport.RespWriter Отсутсвие токенов авторизации
// @Failure 403 {object} transport.RespWriter Отсутсвии токена Администратора
// @Failure 500 {object} transport.RespWriter Внутренняя ошибка сервера
// @Router /film/insert [post]
func (handler *FilmsHandler) InsertFilm(w http.ResponseWriter, r *http.Request) {
	handler.logger.Info("insert film handler")

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
	film := film_model.FilmInsert{}
	if err = json.Unmarshal(data, &film); err != nil {
		handler.logger.Warn("request body scheme error", slog.Any("error", err))
		transport.NewRespWriter(w, err.Error(), http.StatusBadRequest, handler.logger)
		return
	}

	handler.logger.Debug("validate film struct", slog.Any("struct", film))
	_, err = govalidator.ValidateStruct(film)
	if err != nil {
		handler.logger.Warn("struct validate failed", slog.Any("error", err))
		transport.NewRespWriter(w, err.Error(), http.StatusBadRequest, handler.logger)
		return
	}

	handler.logger.Info("insert film")
	id, err := handler.filmsService.InsertFilm(r.Context(), &film)
	if err != nil {
		handler.logger.Warn("insert film failed", slog.Any("error", err))
		transport.NewRespWriter(w, err.Error(), http.StatusInternalServerError, handler.logger)
		return
	}

	transport.NewRespWriter(w, id, http.StatusCreated, handler.logger)
	handler.logger.Info("successful film insert", slog.String("film", id))
}

// Обновление фильма
// @Summary UpdateFilm
// @Security ApiKeyAuth
// @Description Обновление фильма в системе. Возомжно обновлять как частично, так и полностью все поля.
// @Description Если указать список id актеров, список актеров в фильме заменится переданным.
// @Tags films
// @ID update-film
// @Accept json
// @Produce json
// @Param input body film_model.UpdateFilm true "actor fields"
// @Success 200 {object} nil Успешное обновление фильма
// @Failure 400 {object} transport.RespWriter Некорректный ввод данных
// @Failure 401 {object} transport.RespWriter Отсутсвие токенов авторизации
// @Failure 403 {object} transport.RespWriter Отсутсвии токена Администратора
// @Failure 500 {object} transport.RespWriter Внутренняя ошибка сервера
// @Router /film/update [put]
func (handler *FilmsHandler) UpdateFilm(w http.ResponseWriter, r *http.Request) {
	handler.logger.Info("update film handler")

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
	updateFilm := film_model.UpdateFilm{}

	if err = json.Unmarshal(data, &updateFilm); err != nil {
		handler.logger.Warn("request body scheme error", slog.Any("error", err))
		transport.NewRespWriter(w, err.Error(), http.StatusBadRequest, handler.logger)
		return
	}

	handler.logger.Debug("validate film struct", slog.Any("struct", updateFilm))
	_, err = govalidator.ValidateStruct(updateFilm)
	if err != nil {
		handler.logger.Warn("struct validate failed", slog.Any("error", err))
		transport.NewRespWriter(w, err.Error(), http.StatusBadRequest, handler.logger)
		return
	}

	handler.logger.Info("update film", slog.Any("id", updateFilm.ID))
	if err := handler.filmsService.UpdateFilm(r.Context(), &updateFilm); err != nil {
		handler.logger.Info("film service failed", slog.Any("error", err))
		transport.NewRespWriter(w, err.Error(), http.StatusInternalServerError, handler.logger)
		return
	}

	w.WriteHeader(http.StatusOK)
	handler.logger.Info("successful film update", slog.String("film", updateFilm.ID.String()))
}

// Удаление фильма
// @Summary DeleteFilm
// @Security ApiKeyAuth
// @Description Удаление фильма в системе. Id передается в теле запроса и представлен UUID.
// @Tags films
// @ID delete-film
// @Accept json
// @Produce json
// @Param input body film_model.Id true "actor fields"
// @Success 200 {object} nil Успешное удаление фильма
// @Failure 400 {object} transport.RespWriter Некорректный ввод данных
// @Failure 401 {object} transport.RespWriter Отсутсвие токенов авторизации
// @Failure 403 {object} transport.RespWriter Отсутсвии токена Администратора
// @Failure 500 {object} transport.RespWriter Внутренняя ошибка сервера
// @Router /film/delete [delete]
func (handler *FilmsHandler) DeleteFilm(w http.ResponseWriter, r *http.Request) {
	handler.logger.Info("delete film handler")

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
	deletedFilm := film_model.Id{}

	if err = json.Unmarshal(data, &deletedFilm); err != nil {
		handler.logger.Warn("request body scheme error", slog.Any("error", err))
		transport.NewRespWriter(w, err.Error(), http.StatusBadRequest, handler.logger)
		return
	}

	handler.logger.Debug("validate film struct", slog.Any("struct", deletedFilm))
	_, err = govalidator.ValidateStruct(deletedFilm)
	if err != nil {
		handler.logger.Warn("struct validate failed", slog.Any("error", err))
		transport.NewRespWriter(w, err.Error(), http.StatusBadRequest, handler.logger)
		return
	}

	if err := handler.filmsService.DeleteFilm(r.Context(), deletedFilm); err != nil {
		handler.logger.Info("film service failed", slog.Any("error", err))
		transport.NewRespWriter(w, err.Error(), http.StatusInternalServerError, handler.logger)
		return
	}

	w.WriteHeader(http.StatusOK)
	handler.logger.Info("successful film delete", slog.String("film", deletedFilm.ID.String()))
}

// Получение списка фильмов
// @Summary GetFilms
// @Security ApiKeyAuth
// @Description Получение фильмов списка с актерами, участвующими в данном фильме.
// @Description Сортировка задается параметрами URL: order и type.
// @Description Если не задать эти параметры, то будет сортировка по убыванию рейтинга
// @Tags films
// @ID get-films
// @Accept json
// @Produce json
// @Param order query string false "type of order"
// @Param type query string false "asc or desc"
// @Success 200 {object} []actorfilm.FilmActors Список фильмов с актерами, участвующими в них.
// @Failure 400 {object} transport.RespWriter  Некорректный ввод данных
// @Failure 401 {object} transport.RespWriter Отсутсвие токенов авторизации
// @Failure 500 {object} transport.RespWriter Внутренняя ошибка сервера
// @Router /films [get]
func (handler *FilmsHandler) GetFilms(w http.ResponseWriter, r *http.Request) {
	handler.logger.Info("get films handler")

	values := r.URL.Query()
	order, orderType := film_model.ValidQuery(values.Get("order"), values.Get("type"))

	handler.logger.Debug("get film", slog.String("sort", order), slog.String("direction", orderType))

	films, err := handler.filmsService.GetFilms(r.Context(), order, orderType)
	if err != nil {
		handler.logger.Info("film service failed", slog.Any("error", err))
		transport.NewRespWriter(w, err.Error(), http.StatusInternalServerError, handler.logger)
		return
	}

	res, err := json.Marshal(films)
	if err != nil {
		handler.logger.Info("json marshal failed", slog.Any("error", err))
		transport.NewRespWriter(w, err.Error(), http.StatusInternalServerError, handler.logger)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(res); err != nil {
		handler.logger.Info("films write failed", slog.Any("error", err))
		transport.NewRespWriter(w, err.Error(), http.StatusInternalServerError, handler.logger)
		return
	}

	handler.logger.Info("get films successful")
}

// Поиск фильма
// @Summary GetFilmsSearch
// @Security ApiKeyAuth
// @Description Поиск фильмов по названию и имени актеров.
// @Description Выводиться сначала список фильмов, вместе с актерами, принимающих участие.
// @Description Далее выводится список актеров с фильмами, в кторых они снимались.
// @Description Запрос передается в параметре search. Список ранжируетя по более точному совпадению.
// @Tags films
// @ID search-films
// @Accept json
// @Produce json
// @Param search query string true "search query"
// @Success 200 {object} actorfilm.SearchStruct  Список фильмов и актеров
// @Failure 400 {object} transport.RespWriter Некорректный ввод данных
// @Failure 401 {object} transport.RespWriter Отсутсвие токенов авторизации
// @Failure 500 {object} transport.RespWriter Внутренняя ошибка сервера
// @Router /films/search [get]
func (handler *FilmsHandler) GetFilmsSearch(w http.ResponseWriter, r *http.Request) {
	handler.logger.Info("get films handler")

	searchQuery := r.URL.Query().Get("search")

	handler.logger.Debug("get film", slog.String("search query", searchQuery))

	films, err := handler.filmsService.SearchFilms(r.Context(), searchQuery)
	if err != nil {
		handler.logger.Info("film service failed", slog.Any("error", err))
		transport.NewRespWriter(w, err.Error(), http.StatusInternalServerError, handler.logger)
		return
	}

	actors, err := handler.actorService.SearchActor(r.Context(), searchQuery)
	if err != nil {
		handler.logger.Info("actor service failed", slog.Any("error", err))
		transport.NewRespWriter(w, err.Error(), http.StatusInternalServerError, handler.logger)
		return
	}

	actorsAndFilms := &actorfilm.SearchStruct{
		Actors: actors,
		Films:  films,
	}

	res, err := json.Marshal(actorsAndFilms)
	if err != nil {
		handler.logger.Info("json marshal failed", slog.Any("error", err))
		transport.NewRespWriter(w, err.Error(), http.StatusInternalServerError, handler.logger)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(res); err != nil {
		handler.logger.Info("films write failed", slog.Any("error", err))
		transport.NewRespWriter(w, err.Error(), http.StatusInternalServerError, handler.logger)
		return
	}

	handler.logger.Info("get films successful")
}
