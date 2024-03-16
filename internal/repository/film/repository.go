package film_repository

import (
	"context"

	film_model "github.com/Heater_dog/Vk_Intern/internal/models/film"
	"github.com/google/uuid"
)

//go:generate go run github.com/vektra/mockery/v2@v2.42.1 --name=FilmsRepository
type FilmsRepository interface {
	GetFilmsWithActor(ctx context.Context, userID uuid.UUID) ([]film_model.Film, error)
}
