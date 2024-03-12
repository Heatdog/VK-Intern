package userDb

import (
	"context"

	"github.com/Heater_dog/Vk_Intern/internal/user"
	"github.com/Heater_dog/Vk_Intern/pkg/client"
	"github.com/sirupsen/logrus"
)

type repository struct {
	dbClient client.Client
	logger   *logrus.Logger
}

func NewUserPostgreRepository(dbClient client.Client, logger *logrus.Logger) user.UserRepository {
	return &repository{
		dbClient: dbClient,
		logger:   logger,
	}
}

func (repo *repository) Find(ctx context.Context, login string) (*user.User, error) {
	repo.logger.Infof("USer repo find method for login: %s", login)
	q := `
			SELECT id, login, password, role
			FROM Users
			WHERE login = $1
	`
	row := repo.dbClient.QueryRow(ctx, q, login)

	var res user.User

	if err := row.Scan(&res); err != nil {
		repo.logger.Infof("SQL error: %s", err.Error())
		return nil, err
	}

	return &res, nil
}
