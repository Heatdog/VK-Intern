package actorfilm

import (
	actor_model "github.com/Heater_dog/Vk_Intern/internal/models/actor"
	film_model "github.com/Heater_dog/Vk_Intern/internal/models/film"
)

type ActorFilms struct {
	Actor actor_model.Actor `json:"actor"`
	Films []film_model.Film `json:"films"`
}
