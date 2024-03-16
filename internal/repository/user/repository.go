package user_repository

import (
	"context"

	user_model "github.com/Heater_dog/Vk_Intern/internal/models/user"
)

//go:generate go run github.com/vektra/mockery/v2@v2.42.1 --name=UserRepository
type UserRepository interface {
	Find(ctx context.Context, login string) (*user_model.User, error)
}
