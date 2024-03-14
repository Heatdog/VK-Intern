package userDb

import (
	"context"
	"log/slog"

	"github.com/Heater_dog/Vk_Intern/internal/user"
	"github.com/Heater_dog/Vk_Intern/pkg/client"
)

type repository struct {
	dbClient client.Client
	logger   *slog.Logger
}

func NewUserPostgreRepository(dbClient client.Client, logger *slog.Logger) user.UserRepository {
	return &repository{
		dbClient: dbClient,
		logger:   logger,
	}
}

func (repo *repository) Find(ctx context.Context, login string) (*user.User, error) {
	repo.logger.Info("find user in repo", slog.Any("login", login))
	q := `
			SELECT id, login, password, role
			FROM users
			WHERE login = $1
	`
	repo.logger.Debug("user repo query", slog.String("query", q))
	row := repo.dbClient.QueryRow(ctx, q, login)

	var res user.User

	if err := row.Scan(&res.ID, &res.Login, &res.Password, &res.Role); err != nil {
		repo.logger.Error("SQL error", slog.Any("error", err))
		return nil, err
	}

	return &res, nil
}
