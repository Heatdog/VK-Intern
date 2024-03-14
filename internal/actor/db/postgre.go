package dbActor

import (
	"context"
	"log/slog"

	"github.com/Heater_dog/Vk_Intern/internal/actor"
	"github.com/Heater_dog/Vk_Intern/pkg/client"
	"github.com/google/uuid"
)

type repository struct {
	dbClient client.Client
	logger   *slog.Logger
}

func NewActorPostgreRepository(dbClient client.Client, logger *slog.Logger) actor.ActorsRepository {
	return &repository{
		dbClient: dbClient,
		logger:   logger,
	}
}

func (repo *repository) AddActor(ctx context.Context, actor actor.ActorInsert) (string, error) {
	repo.logger.Info("insert actor in repo")
	q := `
			INSERT INTO actors (name, gender, birth_date)
			VALUES ($1, $2, $3) 
			RETURNING id
	`
	repo.logger.Debug("actor repo query", slog.String("query", q))
	row := repo.dbClient.QueryRow(ctx, q, actor.Name, actor.Gender, actor.BirthDate)

	var id uuid.UUID

	if err := row.Scan(&id); err != nil {
		repo.logger.Error("SQL error", slog.Any("error", err))
		return uuid.Nil.String(), err
	}

	return id.String(), nil
}
