package actor

import (
	"context"
	"log/slog"

	"github.com/Heater_dog/Vk_Intern/internal/film"
	"github.com/google/uuid"
)

//go:generate go run github.com/vektra/mockery/v2@v2.42.1 --name=ActorsService
type ActorsService interface {
	InsertActor(context context.Context, actor ActorInsert) (string, error)
	GetActors(context.Context) ([]ActorFilms, error)
	DeleteActor(context context.Context, userID uuid.UUID) error
	UpdateActor(contex context.Context, fileds UpdateActor) error
}

type actorsService struct {
	logger        *slog.Logger
	storageActors ActorsRepository
	storageFilms  film.FilmsRepository
}

func NewActorsService(logger *slog.Logger, storageActors ActorsRepository,
	storageFilms film.FilmsRepository) ActorsService {
	return &actorsService{
		logger:        logger,
		storageActors: storageActors,
		storageFilms:  storageFilms,
	}
}

func (service *actorsService) InsertActor(context context.Context, actor ActorInsert) (string, error) {
	service.logger.Info("actor service insert")

	service.logger.Debug("insert actor in repositpry")
	return service.storageActors.AddActor(context, actor)
}

func (service *actorsService) GetActors(context context.Context) ([]ActorFilms, error) {
	service.logger.Info("actor service get")

	service.logger.Debug("get actors from repositpry")
	actros, err := service.storageActors.GetActors(context)
	if err != nil {
		service.logger.Error("actors storage error", slog.Any("err", err))
		return nil, err
	}

	var res []ActorFilms

	for _, actor := range actros {
		service.logger.Debug("get films with actor", slog.String("actor", actor.ID.String()))
		films, err := service.storageFilms.GetFilmsWithActor(context, actor.ID)
		if err != nil {
			service.logger.Error("films storage error", slog.Any("err", err))
			return nil, err
		}

		res = append(res, ActorFilms{
			Actor: actor,
			Films: films,
		})
	}

	return res, nil
}

func (service *actorsService) DeleteActor(context context.Context, actorId uuid.UUID) error {
	service.logger.Info("actor service delete", slog.String("id", actorId.String()))

	service.logger.Debug("delete actor in repository")
	return service.storageActors.DeleteActor(context, actorId)
}

func (service *actorsService) UpdateActor(contex context.Context, fileds UpdateActor) error {
	service.logger.Info("actor service update", slog.String("id", fileds.ID.String()))

	service.logger.Debug("update actor in repository")
	if fileds.Name != "" {
		service.logger.Debug("update actor name", slog.Any("id", fileds.ID))
		if err := service.storageActors.UpdateName(contex, fileds.ID, fileds.Name); err != nil {
			return err
		}
	}

	if fileds.Gender != "" {
		service.logger.Debug("update actor gender", slog.Any("id", fileds.ID))
		if err := service.storageActors.UpdateGender(contex, fileds.ID, fileds.Gender); err != nil {
			return err
		}
	}

	if fileds.BirthDate != "" {
		service.logger.Debug("update actor bith date", slog.Any("id", fileds.ID))
		if err := service.storageActors.UpdateBirthDate(contex, fileds.ID, fileds.BirthDate); err != nil {
			return err
		}
	}

	return nil
}
