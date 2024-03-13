package transport

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"

	"github.com/Heater_dog/Vk_Intern/internal/user"
)

type AuthHandler struct {
	logger      *slog.Logger
	userService *user.UserService
}

func NewAuthHandler(logger *slog.Logger, userService *user.UserService) Handler {
	return &AuthHandler{
		logger:      logger,
		userService: userService,
	}
}

const (
	signInURL = "/login"
)

func (handler *AuthHandler) Register(router *http.ServeMux) {
	router.HandleFunc(signInURL, handler.signInHandle)
}

func (handler *AuthHandler) signInHandle(w http.ResponseWriter, r *http.Request) {
	handler.logger.Info("sign in user")

	handler.logger.Debug("read request body")
	data, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		handler.logger.Error("request body reading failed", slog.Any("error", err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var user user.UserLogin
	handler.logger.Debug("unmarshaling request body")
	if err = json.Unmarshal(data, &user); err != nil {
		handler.logger.Error("request body scheme error", slog.Any("error", err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	handler.logger.Debug("sign in service", slog.String("user", user.Login))
	accessToken, refreshToken, expire, err := handler.userService.SignIn(r.Context(), user)
	if err != nil {
		handler.logger.Warn("user service error", slog.Any("error", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "refresh token",
		Value:   refreshToken,
		Expires: expire,
	})
	w.Header().Set("Authorization", "Bearer "+accessToken)
	w.WriteHeader(http.StatusOK)
	handler.logger.Info("successful auth", slog.String("user", user.Login))

}
