package user_db

import (
	"context"
	"log/slog"

	user_model "github.com/Heater_dog/Vk_Intern/internal/models/user"
	user_repository "github.com/Heater_dog/Vk_Intern/internal/repository/user"
	"github.com/Heater_dog/Vk_Intern/pkg/client"
)

type repository struct {
	dbClient client.Client
	logger   *slog.Logger
}

func NewUserPostgreRepository(dbClient client.Client, logger *slog.Logger) user_repository.UserRepository {
	return &repository{
		dbClient: dbClient,
		logger:   logger,
	}
}

func (repo *repository) Find(ctx context.Context, login string) (*user_model.User, error) {
	repo.logger.Info("find user in repo", slog.Any("login", login))
	q := `
			SELECT id, login, password, role
			FROM users
			WHERE login = $1
	`
	repo.logger.Debug("user repo query", slog.String("query", q))
	row := repo.dbClient.QueryRow(ctx, q, login)

	var res user_model.User

	if err := row.Scan(&res.ID, &res.Login, &res.Password, &res.Role); err != nil {
		repo.logger.Error("SQL error", slog.Any("error", err))
		return nil, err
	}

	return &res, nil
}
