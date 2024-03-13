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

// SignIn	godoc
// @Summary SignIn
// @Tags auth
// @Description sign in web site
// @ID sign-in
// @Accept json
// @Produce json
// @Param input body user.UserLogin true "account info"
// @Success 200 {integer} integer 1
// @Failure 400 {object} respWriter
// @Failure 500 {object} respWriter
// @Router /login [post]
func (handler *AuthHandler) signInHandle(w http.ResponseWriter, r *http.Request) {
	handler.logger.Info("sign in user")

	handler.logger.Debug("read request body")
	data, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		handler.logger.Error("request body reading failed", slog.Any("error", err))
		NewRespWriter(w, err.Error(), http.StatusBadRequest, handler.logger)
		return
	}
	handler.logger.Debug("Request body", slog.String("body", string(data)))
	var user user.UserLogin
	handler.logger.Debug("unmarshaling request body")
	if err = json.Unmarshal(data, &user); err != nil {
		handler.logger.Error("request body scheme error", slog.Any("error", err))
		NewRespWriter(w, err.Error(), http.StatusBadRequest, handler.logger)
		return
	}

	handler.logger.Debug("sign in service", slog.String("user", user.Login))
	accessToken, refreshToken, expire, err := handler.userService.SignIn(r.Context(), user)
	if err != nil {
		handler.logger.Warn("user service error", slog.Any("error", err))
		NewRespWriter(w, err.Error(), http.StatusInternalServerError, handler.logger)
		return
	}

	handler.logger.Info("user tokens set", slog.String("user", user.Login),
		slog.String("access token", accessToken), slog.String("refresh token", refreshToken))
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Authorization", "Bearer "+accessToken)
	http.SetCookie(w, &http.Cookie{
		Name:    "refresh token",
		Value:   refreshToken,
		Expires: expire,
	})
	handler.logger.Info("successful auth", slog.String("user", user.Login))
}
