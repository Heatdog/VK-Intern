package user

import (
	"context"
	"fmt"

	"github.com/Heater_dog/Vk_Intern/internal/auth"
	cryptohash "github.com/Heater_dog/Vk_Intern/pkg/cryptoHash"
	"github.com/sirupsen/logrus"
)

type UserService struct {
	logger      *logrus.Logger
	userRepo    UserRepository
	authService *auth.TokenService
	passwordKey string
}

func NewUserService(logger *logrus.Logger, userRepo UserRepository, authService *auth.TokenService) *UserService {
	return &UserService{
		logger:      logger,
		userRepo:    userRepo,
		authService: authService,
	}
}

func (service *UserService) SignIn(ctx context.Context, user UserLogin) (string, error) {
	service.logger.Infof("User %s sign in started", user.Login)
	res, err := service.userRepo.Find(ctx, user.Login)
	if err != nil {
		service.logger.Infof("User repo error: %s", err.Error())
		return "", err
	}

	if cryptohash.VerifyHash([]byte(res.Password), user.Password) {
		token, err := service.authService.GenerateToken(ctx, auth.TokenFileds{
			ID:   res.ID,
			Role: res.Role,
		})
		if err != nil {
			service.logger.Errorf("jwt token generate error: %s", err.Error())
			return "", err
		}
		return token, nil
	} else {
		errStr := fmt.Sprintf("Wrong password for user %s", user.Login)
		service.logger.Info(errStr)
		return "", fmt.Errorf(errStr)
	}
}
