package user

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/Heater_dog/Vk_Intern/internal/auth"
	cryptohash "github.com/Heater_dog/Vk_Intern/pkg/cryptoHash"
)

type UserService struct {
	logger      *slog.Logger
	userRepo    UserRepository
	authService *auth.TokenService
}

func NewUserService(logger *slog.Logger, userRepo UserRepository, authService *auth.TokenService) *UserService {
	return &UserService{
		logger:      logger,
		userRepo:    userRepo,
		authService: authService,
	}
}

func (service *UserService) SignIn(ctx context.Context, user UserLogin) (string, string, time.Time, error) {
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
		accessToken, refreshToken, expire, err := service.authService.GenerateToken(ctx, auth.TokenFileds{
			ID:   res.ID.String(),
			Role: res.Role,
		})
		if err != nil {
			service.logger.Warn("jwt token generate failed", slog.Any("error", err))
			return "", "", time.Time{}, err
		}
		return accessToken, refreshToken, expire, nil
	} else {
		errStr := fmt.Sprint("wrong password", slog.Any("error", user.Login))
		service.logger.Info(errStr)
		return "", "", time.Time{}, fmt.Errorf(errStr)
	}
}
