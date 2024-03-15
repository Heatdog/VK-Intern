package dbActor

import (
	"context"
	"fmt"
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

func (repo *repository) GetActors(ctx context.Context) ([]actor.Actor, error) {
	repo.logger.Info("get actors from repo")
	q := `
			SELECT id, name, gender, birth_date
			FROM actors
	`

	repo.logger.Debug("actor repo query", slog.String("query", q))
	rows, err := repo.dbClient.Query(ctx, q)
	if err != nil {
		repo.logger.Error("select actros from repo err", slog.Any("err", err))
		return nil, err
	}

	var res []actor.Actor

	for rows.Next() {
		var actor actor.Actor

		if err := rows.Scan(&actor.ID, &actor.Name, &actor.Gender, &actor.BirthDate); err != nil {
			repo.logger.Error("SQL Error", slog.Any("err", err))
			return nil, err
		}

		res = append(res, actor)
	}
	return res, nil
}

func (repo *repository) DeleteActor(ctx context.Context, id uuid.UUID) error {
	repo.logger.Info("delete actor from repo")
	q := `
			DELETE FROM actors
			WHERE id = $1
	`
	repo.logger.Debug("actor repo query", slog.String("query", q))
	commandTag, err := repo.dbClient.Exec(ctx, q, id)
	if err != nil {
		repo.logger.Error("SQL error", slog.Any("err", err))
		return err
	}

	if commandTag.RowsAffected() != 1 {
		return fmt.Errorf("No row found to delete")
	}

	return nil
}

func (repo *repository) UpdateName(ctx context.Context, id uuid.UUID, name string) error {
	repo.logger.Info("update name actor in repo")
	q := `
			UPDATE actors SET name = $1
			WHERE id = $2
	`

	repo.logger.Debug("actor repo query", slog.String("query", q))
	commandTag, err := repo.dbClient.Exec(ctx, q, name, id)
	if err != nil {
		repo.logger.Error("SQL error", slog.Any("err", err))
		return err
	}

	if commandTag.RowsAffected() != 1 {
		return fmt.Errorf("No row found to delete")
	}

	return nil
}

func (repo *repository) UpdateGender(ctx context.Context, id uuid.UUID, gender string) error {
	repo.logger.Info("update gender actor in repo")
	q := `
			UPDATE actors SET gender = $1
			WHERE id = $2
	`

	repo.logger.Debug("actor repo query", slog.String("query", q))
	commandTag, err := repo.dbClient.Exec(ctx, q, gender, id)
	if err != nil {
		repo.logger.Error("SQL error", slog.Any("err", err))
		return err
	}

	if commandTag.RowsAffected() != 1 {
		return fmt.Errorf("No row found to delete")
	}

	return nil
}

func (repo *repository) UpdateBirthDate(ctx context.Context, id uuid.UUID, birthDate string) error {
	repo.logger.Info("update birth date actor in repo")
	q := `
			UPDATE actors SET birth_date = $1
			WHERE id = $2
	`

	repo.logger.Debug("actor repo query", slog.String("query", q))
	commandTag, err := repo.dbClient.Exec(ctx, q, birthDate, id)
	if err != nil {
		repo.logger.Error("SQL error", slog.Any("err", err))
		return err
	}

	if commandTag.RowsAffected() != 1 {
		return fmt.Errorf("No row found to delete")
	}

	return nil
}
