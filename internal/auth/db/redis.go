package authDb

import (
	"context"
	"time"

	"github.com/Heater_dog/Vk_Intern/internal/auth"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

type RedisTokenStorage struct {
	client *redis.Client
	logger *logrus.Logger
}

func NewRedisTokenStorage(logger *logrus.Logger, storage *redis.Client) auth.TokenStorage {
	return &RedisTokenStorage{
		client: storage,
		logger: logger,
	}
}

func (storage *RedisTokenStorage) SetToken(ctx context.Context, userId, token string, expire time.Duration) error {
	if err := storage.client.Set(ctx, userId, token, expire).Err(); err != nil {
		storage.logger.Infof("Redis set token error: %s", err.Error())
		return err
	}
	return nil
}
