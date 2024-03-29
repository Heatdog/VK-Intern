package user_service

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	user_model "github.com/Heater_dog/Vk_Intern/internal/models/user"
	user_repository "github.com/Heater_dog/Vk_Intern/internal/repository/user"
	token_service "github.com/Heater_dog/Vk_Intern/internal/services/token"
	cryptohash "github.com/Heater_dog/Vk_Intern/pkg/cryptoHash"
	"github.com/Heater_dog/Vk_Intern/pkg/jwt"
)

//go:generate go run github.com/vektra/mockery/v2@v2.42.1 --name=UserService
type UserService interface {
	SignIn(ctx context.Context, user user_model.UserLogin) (accessToken string, refreshToken string, expire time.Time, err error)
}

type userService struct {
	logger      *slog.Logger
	userRepo    user_repository.UserRepository
	authService token_service.TokenService
}

func NewUserService(logger *slog.Logger, userRepo user_repository.UserRepository,
	authService token_service.TokenService) UserService {
	return &userService{
		logger:      logger,
		userRepo:    userRepo,
		authService: authService,
	}
}

func (service *userService) SignIn(ctx context.Context, user user_model.UserLogin) (string, string, time.Time, error) {
	service.logger.Info("sign in", slog.String("user", user.Login))
	service.logger.Debug("get user from repo", slog.String("user", user.Login))
	res, err := service.userRepo.Find(ctx, user.Login)
	if err != nil {
		service.logger.Warn("user repo error", slog.Any("error", err))
		return "", "", time.Time{}, err
	}

	service.logger.Debug("verify password", slog.String("user", user.Login))
	if cryptohash.VerifyHash([]byte(res.Password), user.Password) {
		service.logger.Debug("generate tokens", slog.String("user", user.Login))
		accessToken, refreshToken, expire, err := service.authService.GenerateToken(ctx, jwt.TokenFileds{
			ID:   res.ID.String(),
			Role: res.Role,
		})
		if err != nil {
			service.logger.Warn("jwt token generate failed", slog.Any("error", err))
			return "", "", time.Time{}, err
		}
		return accessToken, refreshToken, expire, nil
	}

	errStr := fmt.Sprint("wrong password ", slog.Any("error", user.Login))
	service.logger.Info(errStr)
	return "", "", time.Time{}, fmt.Errorf(errStr)

}
