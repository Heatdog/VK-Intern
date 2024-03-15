package film

import (
	"time"

	"github.com/google/uuid"
)

type Film struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	ReleaseDate time.Time `json:"release_date"`
	Rating      float32   `json:"rating"`
}
