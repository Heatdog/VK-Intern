package auth

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
)

type TokenService struct {
	logger      *logrus.Logger
	secretKey   string
	storage     TokenStorage
	tokenExpire int
}

func NewTokenService(logger *logrus.Logger, storage TokenStorage, secretKey string, tokeExpire int) *TokenService {
	return &TokenService{
		logger:      logger,
		storage:     storage,
		secretKey:   secretKey,
		tokenExpire: tokeExpire,
	}
}

func (service *TokenService) GenerateToken(ctx context.Context, tokenFields TokenFileds) (string, error) {
	service.logger.Infof("Generate token for user: %s", tokenFields.ID)
	accessToken, err := GenerateToken(tokenFields, service.secretKey)
	if err != nil {
		service.logger.Errorf("Generate token error: %s", err.Error())
		return "", err
	}
	refreshToken, err := GenerateRefreshToken()
	if err != nil {
		service.logger.Errorf("Generate refresh token error: %s", err.Error())
		return "", err
	}
	expire := time.Hour * 60 * time.Duration(service.tokenExpire)
	if err = service.storage.SetToken(ctx, tokenFields.ID, refreshToken, expire); err != nil {
		service.logger.Errorf("Set token in repo error: %s", err.Error())
		return "", err
	}
	return accessToken, nil
}
