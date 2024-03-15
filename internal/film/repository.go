package film

import (
	"context"

	"github.com/google/uuid"
)

//go:generate go run github.com/vektra/mockery/v2@v2.42.1 --name=FilmsRepository
type FilmsRepository interface {
	GetFilmsWithActor(ctx context.Context, userID uuid.UUID) ([]Film, error)
}
