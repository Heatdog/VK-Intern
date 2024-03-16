package actor_model

import (
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/google/uuid"
)

func init() {
	govalidator.TagMap["gender"] = govalidator.Validator(func(str string) bool {
		return str == "Male" || str == "Female"
	})
	govalidator.TagMap["date"] = govalidator.Validator(func(str string) bool {
		if _, err := time.Parse(time.DateOnly, str); err != nil {
			return false
		}
		return true
	})
}

type ActorInsert struct {
	Name      string `json:"name,omitempty" valid:",required"`
	Gender    string `json:"gender,omitempty" valid:"gender,required"`
	BirthDate string `json:"birth_date,omitempty" valid:"date,required"`
}

type UpdateActor struct {
	ID        uuid.UUID `json:"id" valid:"uuid,required"`
	Name      string    `json:"name"`
	BirthDate string    `json:"birth_date" valid:"date"`
	Gender    string    `json:"gender" valid:"gender"`
}
