package auth

import (
	"context"
	"log/slog"
	"time"
)

//go:generate go run github.com/vektra/mockery/v2@v2.42.1 --name=TokenService
type TokenService interface {
	GenerateToken(ctx context.Context, tokenFields TokenFileds) (accessToken string, refreshToken string,
		expire time.Time, err error)
}

type tokenService struct {
	logger      *slog.Logger
	secretKey   string
	storage     TokenStorage
	tokenExpire int
}

func NewTokenService(logger *slog.Logger, storage TokenStorage, secretKey string, tokeExpire int) TokenService {
	return &tokenService{
		logger:      logger,
		storage:     storage,
		secretKey:   secretKey,
		tokenExpire: tokeExpire,
	}
}

func (service *tokenService) GenerateToken(ctx context.Context, tokenFields TokenFileds) (string, string, time.Time, error) {
	service.logger.Info("generate tokens", slog.Any("user", tokenFields.ID))
	service.logger.Debug("generate access token", slog.Any("user", tokenFields.ID))
	accessToken, err := GenerateToken(tokenFields, service.secretKey)
	if err != nil {
		service.logger.Error("generate access token failed", slog.Any("error", err))
		return "", "", time.Time{}, err
	}
	service.logger.Debug("generate refresh token", slog.Any("user", tokenFields.ID))
	refreshToken, err := GenerateRefreshToken()
	if err != nil {
		service.logger.Error("generate refresh token failed", slog.Any("error", err))
		return "", "", time.Time{}, err
	}
	expire := time.Duration(service.tokenExpire) * time.Hour * 24
	service.logger.Debug("set token in storage",
		slog.Any("user", tokenFields.ID),
		slog.String("token", refreshToken),
		slog.Any("expire", expire))
	if err = service.storage.SetToken(ctx, tokenFields.ID, refreshToken, expire); err != nil {
		service.logger.Error("set token in repo failed", slog.Any("error", err))
		return "", "", time.Time{}, err
	}
	return accessToken, refreshToken, time.Now().Add(expire), nil
}
