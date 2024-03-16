package actor_model

import (
	"time"

	"github.com/google/uuid"
)

type Gender string

const (
	Male   Gender = "Male"
	Female Gender = "Female"
)

type Actor struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Gender    Gender    `json:"gender"`
	BirthDate time.Time `json:"birth_date"`
}
