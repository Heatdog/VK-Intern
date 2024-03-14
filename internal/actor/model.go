package actor

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
	id        uuid.UUID
	name      string
	gender    Gender
	birthDate time.Time
}
