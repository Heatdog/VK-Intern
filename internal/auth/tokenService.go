package auth

import (
	"context"
	"log/slog"
	"time"
)

type TokenService struct {
	logger      *slog.Logger
	secretKey   string
	storage     TokenStorage
	tokenExpire int
}

func NewTokenService(logger *slog.Logger, storage TokenStorage, secretKey string, tokeExpire int) *TokenService {
	return &TokenService{
		logger:      logger,
		storage:     storage,
		secretKey:   secretKey,
		tokenExpire: tokeExpire,
	}
}

func (service *TokenService) GenerateToken(ctx context.Context, tokenFields TokenFileds) (string, string, time.Time, error) {
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
	expire := time.Hour * 60 * time.Duration(service.tokenExpire)
	service.logger.Debug("set token in storage",
		slog.Any("user", tokenFields.ID),
		slog.String("token", refreshToken))
	if err = service.storage.SetToken(ctx, tokenFields.ID, refreshToken, expire); err != nil {
		service.logger.Error("set token in repo failed", slog.Any("error", err))
		return "", "", time.Time{}, err
	}
	return accessToken, refreshToken, time.Now().Add(expire), nil
}
