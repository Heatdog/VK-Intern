package db

import (
	"context"
	"log/slog"

	"github.com/Heater_dog/Vk_Intern/internal/film"
	"github.com/Heater_dog/Vk_Intern/pkg/client"
	"github.com/google/uuid"
)

type repository struct {
	dbClient client.Client
	logger   *slog.Logger
}

func NewFilmsPostgreRepository(dbClient client.Client, logger *slog.Logger) film.FilmsRepository {
	return &repository{
		dbClient: dbClient,
		logger:   logger,
	}
}

func (repo *repository) GetFilmsWithActor(ctx context.Context, userID uuid.UUID) ([]film.Film, error) {
	repo.logger.Info("get films from repo")
	q := `
			SELECT f.id, f.title, f.description, f.release_date, f.rating
			FROM films f
			LEFT JOIN actors_to_films af ON af.film_id = f.id
			WHERE af.user_id = $1
			
	`

	repo.logger.Debug("films repo query", slog.String("query", q))
	rows, err := repo.dbClient.Query(ctx, q, userID)
	if err != nil {
		repo.logger.Error("select films from repo err", slog.Any("err", err))
		return nil, err
	}

	var res []film.Film
	for rows.Next() {
		var film film.Film

		if err := rows.Scan(&film.ID, &film.Title, &film.Description, &film.ReleaseDate, &film.Rating); err != nil {
			repo.logger.Error("SQL Error", slog.Any("err", err))
			return nil, err
		}

		res = append(res, film)
	}

	return res, nil
}
