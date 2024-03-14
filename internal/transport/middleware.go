package transport

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/Heater_dog/Vk_Intern/internal/auth"
	"github.com/Heater_dog/Vk_Intern/pkg/jwt"
)

type Middleware struct {
	logger      *slog.Logger
	authService auth.TokenService
	key         string
}

func NewMiddleware(logger *slog.Logger, authService auth.TokenService, key string) *Middleware {
	return &Middleware{
		logger:      logger,
		authService: authService,
		key:         key,
	}
}

func (mid *Middleware) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mid.logger.Debug("verify access token")

		header := r.Header.Get("Authorization")
		if header == "" {
			newAccessToken, newRefreshToken, expire, err := mid.emptyAccessTokenHeader(r)
			if err != nil {
				mid.logger.Info("access header error")
				NewRespWriter(w, err.Error(), http.StatusUnauthorized, mid.logger)
				return
			}

			http.SetCookie(w, &http.Cookie{
				Name:     "token",
				Value:    newRefreshToken,
				HttpOnly: true,
				Expires:  expire,
				Secure:   true,
			})
			w.Header().Set("Authorization", "Bearer "+newAccessToken)

			next.ServeHTTP(w, r)
		}

		mid.logger.Debug("got access token")
		err := mid.verifyTokenHeader(header)
		if err != nil {
			mid.logger.Warn("auth header err", slog.Any("err", err))
			NewRespWriter(w, err.Error(), http.StatusUnauthorized, mid.logger)
			return
		}

		next.ServeHTTP(w, r)
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

func (mid *Middleware) emptyAccessTokenHeader(r *http.Request) (string, string, time.Time, error) {
	mid.logger.Debug("empty token header")
	refreshToken, err := mid.getRefreshToken(r)
	if err != nil {
		mid.logger.Info("empty refresh token cookie")
		return "", "", time.Time{}, err
	}

	newAccessToken, newRefreshToken, expire, err := mid.authService.
		VerifyToken(r.Context(), refreshToken)

	if err != nil {
		mid.logger.Warn("auth serrvice err", slog.Any("err", err))
		return "", "", time.Time{}, err
	}
	return newAccessToken, newRefreshToken, expire, nil
}

func (mid *Middleware) verifyTokenHeader(header string) error {
	headers := strings.Split(header, " ")
	if len(headers) != 2 {
		err := fmt.Errorf("wrong scheame of auth header")
		mid.logger.Warn("auth header err", slog.Any("err", err))
		return err
	}

	if headers[0] != "Bearer" {
		err := fmt.Errorf("wrong scheame of auth header")
		mid.logger.Warn("auth header err", slog.Any("err", err))
		return err
	}

	if _, err := jwt.VerifyToken(headers[1], mid.key); err != nil {
		mid.logger.Warn("auth header err", slog.Any("err", err))
		return err
	}
	return nil
}
