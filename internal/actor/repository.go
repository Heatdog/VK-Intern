package actor

import (
	"context"

	"github.com/google/uuid"
)

//go:generate go run github.com/vektra/mockery/v2@v2.42.1 --name=ActorsRepository
type ActorsRepository interface {
	AddActor(ctx context.Context, actor ActorInsert) (id string, err error)
	GetActors(ctx context.Context) ([]Actor, error)
	DeleteActor(ctx context.Context, id uuid.UUID) error
	UpdateName(ctx context.Context, id uuid.UUID, name string) error
	UpdateBirthDate(ctx context.Context, id uuid.UUID, birthDate string) error
	UpdateGender(ctx context.Context, id uuid.UUID, gender string) error
}
