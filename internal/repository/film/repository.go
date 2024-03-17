package film_repository

import (
	"context"

	film_model "github.com/Heater_dog/Vk_Intern/internal/models/film"
	"github.com/google/uuid"
)

//go:generate go run github.com/vektra/mockery/v2@v2.42.1 --name=FilmsRepository
type FilmsRepository interface {
	GetFilmsWithActor(ctx context.Context, userID uuid.UUID) ([]film_model.Film, error)
	InsertFilm(ctx context.Context, film *film_model.FilmInsert) (string, error)

	UpdateFilmDescription(ctx context.Context, filmId uuid.UUID, description string) error
	UpdateFilmRating(ctx context.Context, filmId uuid.UUID, rating string) error
	UpdateFilmTitle(ctx context.Context, filmId uuid.UUID, title string) error
	UpdateFilmReleaseDate(ctx context.Context, filmId uuid.UUID, date string) error
	UpdateFilmActors(ctx context.Context, filmId uuid.UUID, actorsID []film_model.Id) error

	DeleteActorsFormFilm(ctx context.Context, filmId uuid.UUID) error
	DeleteFilm(ctx context.Context, filmId uuid.UUID) error

	GetFilms(ctx context.Context, order, orderType string) ([]film_model.Film, error)
	SearchFilms(ctx context.Context, searchQuery string) ([]film_model.Film, error)
}
