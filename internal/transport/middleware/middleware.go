package middleware_transport

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	token_service "github.com/Heater_dog/Vk_Intern/internal/services/token"
	"github.com/Heater_dog/Vk_Intern/internal/transport"
	"github.com/Heater_dog/Vk_Intern/pkg/jwt"
)

type Middleware struct {
	logger      *slog.Logger
	authService token_service.TokenService
	key         string
}

func NewMiddleware(logger *slog.Logger, authService token_service.TokenService, key string) *Middleware {
	return &Middleware{
		logger:      logger,
		authService: authService,
		key:         key,
	}
}

func (mid *Middleware) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mid.logger.Debug("verify access token")

		header := r.Header.Get("authorization")
		if header == "" {
			mid.logger.Debug("header is empty")
			newAccessToken, newRefreshToken, expire, err := mid.emptyAccessTokenHeader(r)
			if err != nil {
				mid.logger.Info("access header error")
				transport.NewRespWriter(w, err.Error(), http.StatusUnauthorized, mid.logger)
				return
			}

			mid.logger.Debug("set new refresh token", slog.String("token", newRefreshToken))
			http.SetCookie(w, &http.Cookie{
				Name:     "token",
				Value:    newRefreshToken,
				HttpOnly: true,
				Expires:  expire,
				Path:     "/",
			})
			r.AddCookie(&http.Cookie{
				Name:     "token",
				Value:    newRefreshToken,
				HttpOnly: true,
				Expires:  expire,
				Path:     "/",
			})
			mid.logger.Debug("set new access token", slog.String("token", newAccessToken))
			r.Header.Set("authorization", "Bearer "+newAccessToken)
			w.Header().Set("authorization", "Bearer "+newAccessToken)

			next.ServeHTTP(w, r)
			return
		}

		mid.logger.Debug("got access token", slog.String("token", header))
		_, err := mid.verifyTokenHeader(header)
		if err != nil {
			mid.logger.Warn("auth header err", slog.Any("err", err))
			transport.NewRespWriter(w, err.Error(), http.StatusUnauthorized, mid.logger)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (mid *Middleware) AdminAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mid.logger.Debug("verify admin token")
		header := r.Header.Get("authorization")
		mid.logger.Debug("Get token", slog.String("token", header))
		if header == "" {
			mid.logger.Debug("access header empty")
			transport.NewRespWriter(w, "access header empty", http.StatusUnauthorized, mid.logger)
			return
		}

		mid.logger.Debug("got access token")
		fileds, err := mid.verifyTokenHeader(header)
		if err != nil {
			mid.logger.Warn("auth header err", slog.Any("err", err))
			transport.NewRespWriter(w, err.Error(), http.StatusUnauthorized, mid.logger)
			return
		}
		if fileds.Role != "Admin" {
			mid.logger.Warn("permission denied")
			transport.NewRespWriter(w, "permission denied", http.StatusForbidden, mid.logger)
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
	oldRefreshToken, err := mid.getRefreshToken(r)
	if err != nil {
		mid.logger.Info("empty refresh token cookie")
		return "", "", time.Time{}, err
	}

	newAccessToken, newRefreshToken, expire, err := mid.authService.
		VerifyToken(r.Context(), oldRefreshToken)

	if err != nil {
		mid.logger.Warn("auth serrvice err", slog.Any("err", err))
		return "", "", time.Time{}, err
	}
	return newAccessToken, newRefreshToken, expire, nil
}

func (mid *Middleware) verifyTokenHeader(header string) (*jwt.TokenFileds, error) {
	mid.logger.Debug("check number of fields", slog.String("header", header))

	headers := strings.Split(header, " ")
	if len(headers) != 2 {
		err := fmt.Errorf("wrong scheame of auth header")
		mid.logger.Warn("auth header err", slog.Any("err", err))
		return nil, err
	}

	mid.logger.Debug("check scheame")

	if headers[0] != "Bearer" {
		err := fmt.Errorf("wrong scheame of auth header")
		mid.logger.Warn("auth header err", slog.Any("err", err))
		return nil, err
	}

	mid.logger.Debug("verify token", slog.String("token", string(header[1])))

	fields, err := jwt.VerifyToken(string(headers[1]), mid.key)
	if err != nil {
		mid.logger.Warn("auth header err", slog.Any("err", err))
		return nil, err
	}
	return fields, nil
}
