package actor

import "context"

//go:generate go run github.com/vektra/mockery/v2@v2.42.1 --name=ActorsRepository
type ActorsRepository interface {
	AddActor(ctx context.Context, actor ActorInsert) (id string, err error)
}
