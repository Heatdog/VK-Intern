package actor

import (
	"time"

	"github.com/Heater_dog/Vk_Intern/internal/film"
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

type ActorFilms struct {
	Actor Actor       `json:"actor"`
	Films []film.Film `json:"films"`
}

type UpdateActor struct {
	ID        uuid.UUID `json:"id" valid:",required"`
	Name      string    `json:"name"`
	BirthDate string    `json:"birth_date" valid:"date"`
	Gender    string    `json:"gender" valid:"gender"`
}
