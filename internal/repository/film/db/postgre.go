package film_db

import (
	"context"
	"fmt"
	"log/slog"

	film_model "github.com/Heater_dog/Vk_Intern/internal/models/film"
	film_repository "github.com/Heater_dog/Vk_Intern/internal/repository/film"
	"github.com/Heater_dog/Vk_Intern/pkg/client"
	"github.com/google/uuid"
)

type repository struct {
	dbClient client.Client
	logger   *slog.Logger
}

func NewFilmsPostgreRepository(dbClient client.Client, logger *slog.Logger) film_repository.FilmsRepository {
	return &repository{
		dbClient: dbClient,
		logger:   logger,
	}
}

func (repo *repository) GetFilmsWithActor(ctx context.Context, userID uuid.UUID) ([]film_model.Film, error) {
	repo.logger.Info("get films from repo")
	q := `
			SELECT f.id, f.title, f.description, f.release_date, f.rating
			FROM films f
			LEFT JOIN actors_to_films af ON af.film_id = f.id
			WHERE af.actor_id = $1
	`

	repo.logger.Debug("films repo query", slog.String("query", q))
	rows, err := repo.dbClient.Query(ctx, q, userID)
	if err != nil {
		repo.logger.Error("select films from repo err", slog.Any("err", err))
		return nil, err
	}

	var res []film_model.Film
	for rows.Next() {
		var film film_model.Film

		if err := rows.Scan(&film.ID, &film.Title, &film.Description, &film.ReleaseDate, &film.Rating); err != nil {
			repo.logger.Error("SQL Error", slog.Any("err", err))
			return nil, err
		}

		res = append(res, film)
	}

	return res, nil
}

func (repo *repository) InsertFilm(ctx context.Context, film *film_model.FilmInsert) (string, error) {
	repo.logger.Info("insert films in repo")

	repo.logger.Debug("begin transaction")
	transaction, err := repo.dbClient.Begin(ctx)
	if err != nil {
		repo.logger.Error("transaction begin failed")
		return "", err
	}

	repo.logger.Debug("insert film")
	q := `
		INSERT INTO films (title, description, release_date, rating)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`

	repo.logger.Debug("films repo query", slog.String("query", q))
	row := repo.dbClient.QueryRow(ctx, q, film.Title, film.Description, film.ReleaseDate, film.Rating)

	var id uuid.UUID

	if err := row.Scan(&id); err != nil {
		repo.logger.Error("SQL error", slog.Any("error", err))

		if err = transaction.Rollback(ctx); err != nil {
			repo.logger.Error("SQL rollback error", slog.Any("err", err))
		}
		return "", err
	}

	for _, actorId := range film.UsersID {
		if err := repo.insertFilmActor(ctx, id.String(), actorId.ID.String()); err != nil {
			repo.logger.Error("SQL error", slog.Any("err", err))

			if err = transaction.Rollback(ctx); err != nil {
				repo.logger.Error("SQL rollback error", slog.Any("err", err))
			}
			return "", err
		}
	}

	if err = transaction.Commit(ctx); err != nil {
		repo.logger.Error("commit transaction failed", slog.Any("err", err))
		return "", err
	}

	return id.String(), nil
}

func (repo *repository) insertFilmActor(ctx context.Context, filmId string, actorId string) error {
	repo.logger.Info("insert into films_actors", slog.String("actorId", actorId), slog.String("filmId", filmId))

	q := `
		INSERT INTO actors_to_films (actor_id, film_id)
		VALUES ($1, $2)
	`
	repo.logger.Debug("films repo query", slog.String("query", q))

	tag, err := repo.dbClient.Exec(ctx, q, actorId, filmId)
	if err != nil {
		repo.logger.Error("SQL error", slog.Any("err", err))
		return err
	}

	if tag.RowsAffected() != 1 {
		repo.logger.Error("rows affected failed")
		return fmt.Errorf("rows affected error")
	}

	repo.logger.Info("insert successfull")
	return nil
}

func (repo *repository) UpdateFilmDescription(ctx context.Context, filmId uuid.UUID, description string) error {
	repo.logger.Info("update film description in repo")
	q := `
			UPDATE films SET description = $1
			WHERE id = $2
	`

	return repo.filmUpdate(ctx, q, filmId, description)
}

func (repo *repository) UpdateFilmRating(ctx context.Context, filmId uuid.UUID, rating string) error {
	repo.logger.Info("update film rating in repo")
	q := `
			UPDATE films SET rating = $1
			WHERE id = $2
	`

	return repo.filmUpdate(ctx, q, filmId, rating)
}

func (repo *repository) UpdateFilmTitle(ctx context.Context, filmId uuid.UUID, title string) error {
	repo.logger.Info("update film title in repo")
	q := `
			UPDATE films SET title = $1
			WHERE id = $2
	`

	return repo.filmUpdate(ctx, q, filmId, title)
}

func (repo *repository) UpdateFilmReleaseDate(ctx context.Context, filmId uuid.UUID, date string) error {
	repo.logger.Info("update film date in repo")
	q := `
			UPDATE films SET release_date = $1
			WHERE id = $2
	`

	return repo.filmUpdate(ctx, q, filmId, date)
}

func (repo *repository) filmUpdate(ctx context.Context, qeury string, filmId uuid.UUID, newField string) error {
	repo.logger.Debug("films repo query", slog.String("query", qeury))
	commandTag, err := repo.dbClient.Exec(ctx, qeury, newField, filmId)

	if err != nil {
		repo.logger.Error("SQL error", slog.Any("err", err))
		return err
	}

	if commandTag.RowsAffected() != 1 {
		return fmt.Errorf("No row found to update")
	}
	return nil
}

func (repo *repository) DeleteActorsFormFilm(ctx context.Context, filmId uuid.UUID) error {
	repo.logger.Info("delete actors from film", slog.String("id", filmId.String()))
	q := `
		DELETE FROM actors_to_films
		WHERE film_id = $1
	`
	repo.logger.Debug("films repo query", slog.String("query", q))
	if _, err := repo.dbClient.Exec(ctx, q, filmId); err != nil {
		repo.logger.Error("SQL error", slog.Any("err", err))
		return err
	}
	return nil
}

func (repo *repository) UpdateFilmActors(ctx context.Context, filmId uuid.UUID,
	actorsID []film_model.Id) error {

	repo.logger.Info("update film actors in repo")
	transaction, err := repo.dbClient.Begin(ctx)
	if err != nil {
		repo.logger.Error("transaction start error", slog.Any("err", err))
		return err
	}

	if err := repo.DeleteActorsFormFilm(ctx, filmId); err != nil {
		repo.logger.Error("SQL error", slog.Any("err", err))

		if err = transaction.Rollback(ctx); err != nil {
			repo.logger.Error("transaction rollback error")
		}
		return err
	}

	q := `
		INSERT INTO actors_to_films (film_id, actor_id)
		VALUES ($1, $2)
	`
	for _, actor := range actorsID {
		repo.logger.Debug("films repo query", slog.String("query", q),
			slog.String("film id", filmId.String()), slog.String("actor id", actor.ID.String()))

		commandTag, err := repo.dbClient.Exec(ctx, q, filmId, actor.ID)

		if err != nil {
			repo.logger.Error("SQL error", slog.Any("err", err))

			if err = transaction.Rollback(ctx); err != nil {
				repo.logger.Error("transaction rollback error")
			}
			return err
		}

		if commandTag.RowsAffected() != 1 {
			repo.logger.Error("SQL error", slog.Any("err", err))

			if err = transaction.Rollback(ctx); err != nil {
				repo.logger.Error("transaction rollback error")
			}
			return err
		}

		repo.logger.Debug("Successful insert")
	}

	if err := transaction.Commit(ctx); err != nil {
		repo.logger.Error("commit transaction error")
		return err
	}

	repo.logger.Info("successful update")
	return nil
}

func (repo *repository) DeleteFilm(ctx context.Context, filmId uuid.UUID) error {
	repo.logger.Info("delete film from repo", slog.String("id", filmId.String()))
	q := `
		DELETE FROM films
		WHERE id = $1
	`

	repo.logger.Debug("SQL query", slog.String("query", q))
	tag, err := repo.dbClient.Exec(ctx, q, filmId)

	if err != nil {
		repo.logger.Error("SQL error", slog.Any("err", err))
		return err
	}

	if tag.RowsAffected() != 1 {
		repo.logger.Warn("row affected not equal 1")
		return fmt.Errorf("row affected not equal 1")
	}

	repo.logger.Info("succesful film delete", slog.String("id", filmId.String()))
	return nil
}

func (repo *repository) GetFilms(ctx context.Context, order, orderType string) ([]film_model.Film, error) {
	repo.logger.Info("get films from repo")
	q := fmt.Sprintf(`
		SELECT (id, titile, description, rating, release_date)
		FROM films
		ORDER BY %s %s 
	`, order, orderType)

	repo.logger.Debug("SQL query", slog.String("query", q))
	rows, err := repo.dbClient.Query(ctx, q)

	if err != nil {
		repo.logger.Error("SQL error", slog.Any("err", err))
		return nil, err
	}

	var films []film_model.Film
	for rows.Next() {

		var film film_model.Film
		if err = rows.Scan(&film.ID, &film.Title, &film.Description, &film.Rating, &film.ReleaseDate); err != nil {
			repo.logger.Error("row scan error", slog.Any("err", err))
			return nil, err
		}

		films = append(films, film)
	}

	repo.logger.Info("successful select")
	return films, nil
}
