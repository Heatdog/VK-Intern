package user

import "context"

type UserRepository interface {
	Find(ctx context.Context, login string) (*User, error)
}
