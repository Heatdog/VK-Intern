package film_service

import (
	"context"
	"log/slog"

	film_model "github.com/Heater_dog/Vk_Intern/internal/models/film"
	actor_repository "github.com/Heater_dog/Vk_Intern/internal/repository/actor"
	film_repository "github.com/Heater_dog/Vk_Intern/internal/repository/film"
	"github.com/google/uuid"
)

type FilmService interface {
	InsertFilm(ctx context.Context, film *film_model.FilmInsert) (string, error)
	UpdateFilm(ctx context.Context, film *film_model.UpdateFilm) error
	DeleteFilm(ctx context.Context, filmId film_model.Id) error
	GetFilms(ctx context.Context, order, orderType string) ([]film_model.Film, error)
}

type filmService struct {
	logger    *slog.Logger
	filmRepo  film_repository.FilmsRepository
	actorRepo actor_repository.ActorsRepository
}

func NewFilmService(logger *slog.Logger, filmrRepo film_repository.FilmsRepository,
	actorRepo actor_repository.ActorsRepository) FilmService {
	return &filmService{
		logger:    logger,
		filmRepo:  filmrRepo,
		actorRepo: actorRepo,
	}
}

func (s *filmService) InsertFilm(ctx context.Context, film *film_model.FilmInsert) (string, error) {
	s.logger.Info("film insert service")

	s.logger.Debug("check actors")
	for _, id := range film.UsersID {

		s.logger.Debug("Check actor", slog.String("id", id.ID.String()))
		if _, err := s.actorRepo.GetActor(ctx, id.ID.String()); err != nil {
			s.logger.Info("actor repo error", slog.Any("err", err))
			return "", err
		}
	}

	s.logger.Debug("insert film in repo")
	id, err := s.filmRepo.InsertFilm(ctx, film)
	if err != nil {
		s.logger.Error("film repo error", slog.Any("err", err))
		return "", err
	}

	s.logger.Debug("insert film in repo successfull", slog.String("id", id))
	return id, nil
}

func (s *filmService) UpdateFilm(ctx context.Context, film *film_model.UpdateFilm) error {
	s.logger.Info("film update service")

	s.logger.Debug("update film in repo", slog.String("id", film.ID.String()))
	if film.Description != "" {
		if err := s.updateDescription(ctx, film.ID, film.Description); err != nil {
			return err
		}
	}

	if film.Rating != "" {
		if err := s.updateRating(ctx, film.ID, film.Rating); err != nil {
			return err
		}
	}

	if film.Title != "" {
		if err := s.updateTitle(ctx, film.ID, film.Title); err != nil {
			return err
		}
	}

	if film.ReleaseDate != "" {
		if err := s.updateReleaseDate(ctx, film.ID, film.ReleaseDate); err != nil {
			return err
		}
	}

	if film.ActorsId != nil {
		if err := s.updateActors(ctx, film.ID, film.ActorsId); err != nil {
			return err
		}
	}

	return nil
}

func (s *filmService) updateDescription(ctx context.Context, id uuid.UUID, description string) error {
	s.logger.Debug("update film description", slog.String("id", id.String()))

	if err := s.filmRepo.UpdateFilmDescription(ctx, id, description); err != nil {
		s.logger.Info("update films description failed", slog.Any("err", err))
		return err
	}

	s.logger.Debug("update film description successful", slog.String("id", id.String()))
	return nil
}

func (s *filmService) updateRating(ctx context.Context, id uuid.UUID, ratingFilm string) error {
	s.logger.Debug("update film rating", slog.String("id", id.String()))

	if err := s.filmRepo.UpdateFilmRating(ctx, id, ratingFilm); err != nil {
		s.logger.Info("update films rating failed", slog.Any("err", err))
		return err
	}

	s.logger.Debug("update film rating successful", slog.String("id", id.String()))
	return nil
}

func (s *filmService) updateTitle(ctx context.Context, id uuid.UUID, titile string) error {
	s.logger.Debug("update film title", slog.String("id", id.String()))

	if err := s.filmRepo.UpdateFilmTitle(ctx, id, titile); err != nil {
		s.logger.Info("update films title failed", slog.Any("err", err))
		return err
	}

	s.logger.Debug("update film title successful", slog.String("id", id.String()))
	return nil
}

func (s *filmService) updateReleaseDate(ctx context.Context, id uuid.UUID, date string) error {
	s.logger.Debug("update film release date", slog.String("id", id.String()))

	if err := s.filmRepo.UpdateFilmReleaseDate(ctx, id, date); err != nil {
		s.logger.Info("update films date failed", slog.Any("err", err))
		return err
	}

	s.logger.Debug("update film date successful", slog.String("id", id.String()))
	return nil
}

func (s *filmService) updateActors(ctx context.Context, id uuid.UUID, actorsID []film_model.Id) error {
	s.logger.Debug("update film actors list", slog.String("id", id.String()))

	s.logger.Debug("check actors")
	for _, actor := range actorsID {
		s.logger.Debug("chechk actor", slog.String("actor id", actor.ID.String()))

		if _, err := s.actorRepo.GetActor(ctx, actor.ID.String()); err != nil {
			s.logger.Info("actor repo error", slog.Any("err", err))
			return err
		}
	}

	s.logger.Debug("insert new ators in film")
	if err := s.filmRepo.UpdateFilmActors(ctx, id, actorsID); err != nil {
		s.logger.Warn("film repo error")
		return err
	}

	s.logger.Debug("update film actors successful", slog.String("id", id.String()))
	return nil
}

func (s *filmService) DeleteFilm(ctx context.Context, filmId film_model.Id) error {
	s.logger.Info("film delete service")

	s.logger.Debug("delete film in repo", slog.String("id", filmId.ID.String()))
	return s.filmRepo.DeleteFilm(ctx, filmId.ID)
}

func (s *filmService) GetFilms(ctx context.Context, order, orderType string) ([]film_model.Film, error) {
	s.logger.Info("get films service")

	return s.filmRepo.GetFilms(ctx, order, orderType)
}
