package auth_transport

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"

	user_model "github.com/Heater_dog/Vk_Intern/internal/models/user"
	user_service "github.com/Heater_dog/Vk_Intern/internal/services/user"
	"github.com/Heater_dog/Vk_Intern/internal/transport"
	"github.com/asaskevich/govalidator"
)

type AuthHandler struct {
	logger      *slog.Logger
	userService user_service.UserService
}

func NewAuthHandler(logger *slog.Logger, userService user_service.UserService) *AuthHandler {
	return &AuthHandler{
		logger:      logger,
		userService: userService,
	}
}

const (
	signInURL = "/login"
)

func (handler *AuthHandler) Register(router *http.ServeMux) {
	router.HandleFunc(signInURL, handler.LoginRouting)
}

func (handler *AuthHandler) LoginRouting(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		handler.SignInHandle(w, r)
	} else {
		transport.NewRespWriter(w, "", http.StatusNotFound, handler.logger)
	}
}

// Вход в систему
// @Summary SignIn
// @Tags auth
// @Description Вход в систему. При успешном входе выдаются refresh и access токены.
// @Description Refresh токен храниться в http-only cookie
// @ID sign-in
// @Accept json
// @Produce json
// @Param input body user_model.UserLogin true "account info"
// @Success 200 {object} nil Успешный вход
// @Failure 400 {object} transport.RespWriter Некооректные входные данные
// @Failure 500 {object} transport.RespWriter Внутренняя ошибка сервера
// @Router /login [post]
func (handler *AuthHandler) SignInHandle(w http.ResponseWriter, r *http.Request) {
	handler.logger.Info("sign in user")

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
	var user user_model.UserLogin
	if err = json.Unmarshal(data, &user); err != nil {
		handler.logger.Warn("request body scheme error", slog.Any("error", err))
		transport.NewRespWriter(w, err.Error(), http.StatusBadRequest, handler.logger)
		return
	}

	handler.logger.Debug("validate user struct")
	_, err = govalidator.ValidateStruct(user)
	if err != nil {
		handler.logger.Warn("struct validate failed", slog.Any("error", err))
		transport.NewRespWriter(w, err.Error(), http.StatusBadRequest, handler.logger)
		return
	}

	handler.logger.Debug("sign in service", slog.String("user", user.Login))
	accessToken, refreshToken, expire, err := handler.userService.SignIn(r.Context(), user)
	if err != nil {
		handler.logger.Warn("user service error", slog.Any("error", err))
		transport.NewRespWriter(w, err.Error(), http.StatusInternalServerError, handler.logger)
		return
	}

	handler.logger.Info("user tokens set", slog.String("user", user.Login),
		slog.String("access token", accessToken),
		slog.Group("refresh token", slog.String("token", refreshToken), slog.Any("expire", expire)))
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    refreshToken,
		HttpOnly: true,
		Expires:  expire,
		Path:     "/",
	})
	w.Header().Set("authorization", "Bearer "+accessToken)
	w.WriteHeader(http.StatusOK)
	handler.logger.Info("successful auth", slog.String("user", user.Login))
}
