package actor_repository

import (
	"context"

	actor_model "github.com/Heater_dog/Vk_Intern/internal/models/actor"
	"github.com/google/uuid"
)

//go:generate go run github.com/vektra/mockery/v2@v2.42.1 --name=ActorsRepository
type ActorsRepository interface {
	AddActor(ctx context.Context, actor actor_model.ActorInsert) (id string, err error)

	GetActors(ctx context.Context) ([]actor_model.Actor, error)
	GetActor(ctx context.Context, id string) (actor_model.Actor, error)
	GetActorsWithFilm(ctx context.Context, filmID string) ([]actor_model.Actor, error)
	SearchActors(ctx context.Context, searchQuery string) ([]actor_model.Actor, error)

	DeleteActor(ctx context.Context, id uuid.UUID) error

	UpdateName(ctx context.Context, id uuid.UUID, name string) error
	UpdateBirthDate(ctx context.Context, id uuid.UUID, birthDate string) error
	UpdateGender(ctx context.Context, id uuid.UUID, gender string) error
}
