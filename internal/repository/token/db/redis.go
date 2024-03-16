package token_db

import (
	"context"
	"log/slog"
	"time"

	token_repository "github.com/Heater_dog/Vk_Intern/internal/repository/token"
	"github.com/redis/go-redis/v9"
)

type RedisTokenRepository struct {
	client *redis.Client
	logger *slog.Logger
}

func NewRedisTokenRepository(logger *slog.Logger, client *redis.Client) token_repository.TokenRepository {
	return &RedisTokenRepository{
		client: client,
		logger: logger,
	}
}

func (repo *RedisTokenRepository) SetToken(ctx context.Context, userId, token string, expire time.Duration) error {
	repo.logger.Debug("redis set token", slog.String("id", userId), slog.String("token", token))
	if err := repo.client.Set(ctx, userId, token, expire).Err(); err != nil {
		repo.logger.Error("redis set token failed", slog.Any("error", err))
		return err
	}
	return nil
}

func (repo *RedisTokenRepository) GetToken(ctx context.Context, userId string) (string, error) {
	repo.logger.Debug("redis get token", slog.String("id", userId))
	return repo.client.Get(ctx, userId).Result()
}
