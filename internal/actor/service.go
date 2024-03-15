package actor

import (
	"context"
	"log/slog"

	"github.com/Heater_dog/Vk_Intern/internal/film"
)

//go:generate go run github.com/vektra/mockery/v2@v2.42.1 --name=ActorsService
type ActorsService interface {
	InsertActor(context context.Context, actor ActorInsert) (string, error)
	GetActors(context.Context) ([]ActorFilms, error)
}

type actorsService struct {
	logger        *slog.Logger
	storageActors ActorsRepository
	storageFilms  film.FilmsRepository
}

func NewActorsService(logger *slog.Logger, storageActors ActorsRepository) ActorsService {
	return &actorsService{
		logger:        logger,
		storageActors: storageActors,
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
