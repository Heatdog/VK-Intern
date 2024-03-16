package actor_service

import (
	"context"
	"log/slog"

	actor_model "github.com/Heater_dog/Vk_Intern/internal/models/actor"
	actorfilm "github.com/Heater_dog/Vk_Intern/internal/models/actor_film"
	actor_repository "github.com/Heater_dog/Vk_Intern/internal/repository/actor"
	film_repository "github.com/Heater_dog/Vk_Intern/internal/repository/film"
	"github.com/google/uuid"
)

//go:generate go run github.com/vektra/mockery/v2@v2.42.1 --name=ActorsService
type ActorsService interface {
	InsertActor(context context.Context, actor actor_model.ActorInsert) (string, error)
	GetActors(context.Context) ([]actorfilm.ActorFilms, error)
	DeleteActor(context context.Context, userID uuid.UUID) error
	UpdateActor(contex context.Context, fileds actor_model.UpdateActor) error
}

type actorsService struct {
	logger     *slog.Logger
	repoActors actor_repository.ActorsRepository
	repoFilms  film_repository.FilmsRepository
}

func NewActorsService(logger *slog.Logger, repoActors actor_repository.ActorsRepository,
	repoFilms film_repository.FilmsRepository) ActorsService {
	return &actorsService{
		logger:     logger,
		repoActors: repoActors,
		repoFilms:  repoFilms,
	}
}

func (service *actorsService) InsertActor(context context.Context, actor actor_model.ActorInsert) (string, error) {
	service.logger.Info("actor service insert")

	service.logger.Debug("insert actor in repositpry")
	return service.repoActors.AddActor(context, actor)
}

func (service *actorsService) GetActors(context context.Context) ([]actorfilm.ActorFilms, error) {
	service.logger.Info("actor service get")

	service.logger.Debug("get actors from repositpry")
	actros, err := service.repoActors.GetActors(context)
	if err != nil {
		service.logger.Error("actors storage error", slog.Any("err", err))
		return nil, err
	}

	var res []actorfilm.ActorFilms

	for _, actor := range actros {
		service.logger.Debug("get films with actor", slog.String("actor", actor.ID.String()))
		films, err := service.repoFilms.GetFilmsWithActor(context, actor.ID)
		if err != nil {
			service.logger.Error("films storage error", slog.Any("err", err))
			return nil, err
		}

		res = append(res, actorfilm.ActorFilms{
			Actor: actor,
			Films: films,
		})
	}

	return res, nil
}

func (service *actorsService) DeleteActor(context context.Context, actorId uuid.UUID) error {
	service.logger.Info("actor service delete", slog.String("id", actorId.String()))

	service.logger.Debug("delete actor in repository")
	return service.repoActors.DeleteActor(context, actorId)
}

func (service *actorsService) UpdateActor(contex context.Context, fileds actor_model.UpdateActor) error {
	service.logger.Info("actor service update", slog.String("id", fileds.ID.String()))

	service.logger.Debug("update actor in repository")
	if fileds.Name != "" {
		service.logger.Debug("update actor name", slog.Any("id", fileds.ID))
		if err := service.repoActors.UpdateName(contex, fileds.ID, fileds.Name); err != nil {
			return err
		}
	}

	if fileds.Gender != "" {
		service.logger.Debug("update actor gender", slog.Any("id", fileds.ID))
		if err := service.repoActors.UpdateGender(contex, fileds.ID, fileds.Gender); err != nil {
			return err
		}
	}

	if fileds.BirthDate != "" {
		service.logger.Debug("update actor bith date", slog.Any("id", fileds.ID))
		if err := service.repoActors.UpdateBirthDate(contex, fileds.ID, fileds.BirthDate); err != nil {
			return err
		}
	}

	return nil
}
