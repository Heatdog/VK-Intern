package actor

import (
	"context"
	"log/slog"
)

//go:generate go run github.com/vektra/mockery/v2@v2.42.1 --name=ActorsService
type ActorsService interface {
	InsertActor(context context.Context, actor ActorInsert) (string, error)
}

type actorsService struct {
	logger  *slog.Logger
	storage ActorsRepository
}

func NewActorsService(logger *slog.Logger, storage ActorsRepository) ActorsService {
	return &actorsService{
		logger:  logger,
		storage: storage,
	}
}

func (service *actorsService) InsertActor(context context.Context, actor ActorInsert) (string, error) {
	service.logger.Info("get actor service")

	service.logger.Debug("get actor from repositpry")
	return service.storage.AddActor(context, actor)
}
