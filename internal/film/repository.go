package film

import (
	"context"

	"github.com/google/uuid"
)

type FilmsRepository interface {
	GetFilmsWithActor(ctx context.Context, userID uuid.UUID) ([]Film, error)
}
