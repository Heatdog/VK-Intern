package postgre

import (
	"context"
	"fmt"
	"time"

	"github.com/Heater_dog/Vk_Intern/internal/config"
	"github.com/Heater_dog/Vk_Intern/pkg/client"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPostgreClient(ctx context.Context, cfg config.PostgreStorage) (client.Client, error) {
	time.Sleep(5 * time.Second)
	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Database)
	ctx, cancel := context.WithTimeout(ctx, 12*time.Second)
	defer cancel()

	conn, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, err
	}

	if err = conn.Ping(ctx); err != nil {
		return nil, err
	}
	return conn, nil
}
