package authDb

import (
	"context"
	"log/slog"
	"time"

	"github.com/Heater_dog/Vk_Intern/internal/auth"

	"github.com/redis/go-redis/v9"
)

type RedisTokenStorage struct {
	client *redis.Client
	logger *slog.Logger
}

func NewRedisTokenStorage(logger *slog.Logger, storage *redis.Client) auth.TokenStorage {
	return &RedisTokenStorage{
		client: storage,
		logger: logger,
	}
}

func (storage *RedisTokenStorage) SetToken(ctx context.Context, userId, token string, expire time.Duration) error {
	storage.logger.Debug("redis set toke", slog.String("id", userId), slog.String("token", token))
	if err := storage.client.Set(ctx, userId, token, expire).Err(); err != nil {
		storage.logger.Error("redis set token failed", slog.Any("error", err))
		return err
	}
	return nil
}
