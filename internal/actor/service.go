package actor

import "context"

type ActorsService interface {
	InsertActor(context context.Context, actor ActorInsert) error
}
