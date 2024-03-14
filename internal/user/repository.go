package user

import "context"

//go:generate go run github.com/vektra/mockery/v2@v2.42.1 --name=UserRepository
type UserRepository interface {
	Find(ctx context.Context, login string) (*User, error)
}
