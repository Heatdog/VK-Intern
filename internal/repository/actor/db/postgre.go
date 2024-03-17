package actor_db

import (
	"context"
	"fmt"
	"log/slog"

	actor_model "github.com/Heater_dog/Vk_Intern/internal/models/actor"
	actor_repository "github.com/Heater_dog/Vk_Intern/internal/repository/actor"
	"github.com/Heater_dog/Vk_Intern/pkg/client"
	"github.com/google/uuid"
)

type repository struct {
	dbClient client.Client
	logger   *slog.Logger
}

func NewActorPostgreRepository(dbClient client.Client, logger *slog.Logger) actor_repository.ActorsRepository {
	return &repository{
		dbClient: dbClient,
		logger:   logger,
	}
}

func (repo *repository) AddActor(ctx context.Context, actor actor_model.ActorInsert) (string, error) {
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

func (repo *repository) GetActors(ctx context.Context) ([]actor_model.Actor, error) {
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

	var res []actor_model.Actor

	for rows.Next() {
		var actor actor_model.Actor

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

	return repo.updateActor(ctx, q, id, name)
}

func (repo *repository) UpdateGender(ctx context.Context, id uuid.UUID, gender string) error {
	repo.logger.Info("update gender actor in repo")
	q := `
			UPDATE actors SET gender = $1
			WHERE id = $2
	`

	return repo.updateActor(ctx, q, id, gender)
}

func (repo *repository) UpdateBirthDate(ctx context.Context, id uuid.UUID, birthDate string) error {
	repo.logger.Info("update birth date actor in repo")
	q := `
			UPDATE actors SET birth_date = $1
			WHERE id = $2
	`

	return repo.updateActor(ctx, q, id, birthDate)
}

func (repo *repository) updateActor(ctx context.Context, query string, id uuid.UUID, field string) error {
	repo.logger.Debug("actor repo query", slog.String("query", query))
	commandTag, err := repo.dbClient.Exec(ctx, query, field, id)

	if err != nil {
		repo.logger.Error("SQL error", slog.Any("err", err))
		return err
	}

	if commandTag.RowsAffected() != 1 {
		return fmt.Errorf("No row found to update")
	}

	return nil
}

func (repo *repository) GetActor(ctx context.Context, id string) (actor_model.Actor, error) {
	repo.logger.Info("get actor from repo", slog.String("id", id))

	q := `
		SELECT id, name, gender, birth_date
		FROM actors
		WHERE id = $1
	`
	repo.logger.Debug("actor repo query", slog.String("query", q))
	row := repo.dbClient.QueryRow(ctx, q, id)

	var actor actor_model.Actor

	if err := row.Scan(&actor.ID, &actor.Name, &actor.Gender, &actor.BirthDate); err != nil {
		repo.logger.Error("SQL error", slog.Any("error", err))
		return actor_model.Actor{}, err
	}

	return actor, nil
}

func (repo *repository) GetActorsWithFilm(ctx context.Context, filmID string) ([]actor_model.Actor, error) {
	repo.logger.Info("get actor from repo with film", slog.String("film_id", filmID))

	q := `
		SELECT a.id, a.name, a.gender, a.birth_date
		FROM actors a
		LEFT JOIN actors_to_films af ON af.actor_id = a.id
		WHERE af.film_id = $1
	`
	repo.logger.Debug("actor repo query", slog.String("query", q))
	rows, err := repo.dbClient.Query(ctx, q, filmID)

	if err != nil {
		repo.logger.Error("select actors from repo err", slog.Any("err", err))
		return nil, err
	}

	var res []actor_model.Actor

	for rows.Next() {
		var actor actor_model.Actor

		if err := rows.Scan(&actor.ID, &actor.Name, &actor.Gender, &actor.BirthDate); err != nil {
			repo.logger.Error("SQL Error", slog.Any("err", err))
			return nil, err
		}

		res = append(res, actor)
	}

	return res, nil
}

func (repo *repository) SearchActors(ctx context.Context, searchQuery string) ([]actor_model.Actor, error) {
	repo.logger.Info("search actors in repo")
	q := `
		SELECT id, name, gender, birth_date
		FROM actors
		WHERE to_tsvector(name) @@ to_tsquery($1)
		ORDER BY ts_rank(to_tsvector(name), to_tsquery($1)) DESC
	`

	repo.logger.Debug("SQL query", slog.String("query", q))
	rows, err := repo.dbClient.Query(ctx, q, searchQuery)

	if err != nil {
		repo.logger.Error("SQL error", slog.Any("err", err))
		return nil, err
	}

	var actors []actor_model.Actor
	for rows.Next() {

		var actor actor_model.Actor
		if err = rows.Scan(&actor.ID, &actor.Name, &actor.Gender, &actor.BirthDate); err != nil {
			repo.logger.Error("row scan error", slog.Any("err", err))
			return nil, err
		}

		actors = append(actors, actor)
	}

	repo.logger.Info("successful select")
	return actors, nil
}
