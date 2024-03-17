package film_model

import (
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/google/uuid"
)

func init() {
	govalidator.TagMap["date"] = govalidator.Validator(func(str string) bool {
		if _, err := time.Parse(time.DateOnly, str); err != nil {
			return false
		}
		return true
	})
}

type Film struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title" valid:",required"`
	Description string    `json:"description" valid:",required"`
	ReleaseDate time.Time `json:"release_date" valid:"date,required"`
	Rating      float32   `json:"rating" valid:"float,required,range(0.0|10.0)"`
}

type Id struct {
	ID uuid.UUID `json:"id" valid:",required"`
}

type FilmInsert struct {
	Title       string `json:"title" valid:"length(1|150),required"`
	Description string `json:"description" valid:"length(0|1000),required"`
	ReleaseDate string `json:"release_date" valid:"date,required"`
	Rating      string `json:"rating" valid:"float,required,range(0|10)"`

	UsersID []Id `json:"actors"`
}

type UpdateFilm struct {
	ID uuid.UUID `json:"id" valid:",required"`

	Title       string `json:"title" valid:"length(1|150)"`
	Description string `json:"description" valid:"length(0|1000)"`
	ReleaseDate string `json:"release_date" valid:"date"`
	Rating      string `json:"rating" valid:"float,range(0|10)"`

	ActorsId []Id `json:"actors"`
}

func ValidQuery(sort string, sortDir string) (string, string) {
	var resSort string
	switch sort {
	case "name":
		resSort = "title"
	case "date":
		resSort = "release_date"
	default:
		resSort = "rating"
	}

	if sortDir == "asc" {
		return resSort, "ASC"
	}
	return resSort, "DESC"
}
