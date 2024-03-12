package transport

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/Heater_dog/Vk_Intern/internal/user"
	"github.com/sirupsen/logrus"
)

type AuthHandler struct {
	logger      *logrus.Logger
	userService *user.UserService
}

func NewAuthHandler(logger *logrus.Logger, userService *user.UserService) Handler {
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
	handler.logger.Infof("Sign in user")

	data, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		handler.logger.Errorf("Request body reading error: %s", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var user user.UserLogin
	if err = json.Unmarshal(data, &user); err != nil {
		handler.logger.Errorf("Request body scheme error: %s", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	token, err := handler.userService.SignIn(r.Context(), user)
	if err != nil {
		handler.logger.Infof("User service error: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Authorization", "Bearer "+token)
	w.WriteHeader(http.StatusOK)
	handler.logger.Infof("Successful auth: %s", user.Login)

}
