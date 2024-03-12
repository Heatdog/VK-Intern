package redisStorage

import (
	"context"
	"fmt"
	"time"

	"github.com/Heater_dog/Vk_Intern/internal/config"
	"github.com/redis/go-redis/v9"
)

func NewRedisClient(ctx context.Context, cfg *config.RedisStorage) (*redis.Client, error) {
	time.Sleep(5 * time.Second)
	host := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)
	storage := redis.NewClient(&redis.Options{
		Addr:     host,
		Password: cfg.Password,
		DB:       0,
	})

	if _, err := storage.Ping(ctx).Result(); err != nil {
		return nil, err
	}
	return storage, nil
}
