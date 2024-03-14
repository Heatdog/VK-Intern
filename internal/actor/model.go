package actor

import (
	"time"

	"github.com/google/uuid"
)

type gender string

const (
	Male   gender = "Male"
	Female gender = "Female"
)

type Actor struct {
	id        uuid.UUID
	name      string
	gender    gender
	birthDate time.Time
}
