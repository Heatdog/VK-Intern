package transport

import (
	"log/slog"
	"net/http"

	"github.com/Heater_dog/Vk_Intern/internal/auth"
)

type Middleware struct {
	logger      *slog.Logger
	authService auth.TokenService
	key         []byte
}

func (mid *Middleware) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mid.logger.Debug("verify access token")

		header := r.Header.Get("Authorization")
		if header == "" {
			mid.logger.Debug("empty token header")
			refreshToken, err := mid.getRefreshToken(r)
			if err != nil {
				mid.logger.Info("empty refresh token cookie")
				NewRespWriter(w, err.Error(), http.StatusUnauthorized, mid.logger)
				return
			} else {

			}
		}
	})
}

func (mid *Middleware) getRefreshToken(r *http.Request) (string, error) {
	mid.logger.Debug("check refresh token")
	cookie, err := r.Cookie("token")
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}
