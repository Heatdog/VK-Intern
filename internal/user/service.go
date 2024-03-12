package user

import (
	"context"
	"fmt"

	cryptohash "github.com/Heater_dog/Vk_Intern/pkg/cryptoHash"
	jwttoken "github.com/Heater_dog/Vk_Intern/pkg/jwtToken"
	"github.com/sirupsen/logrus"
)

type UserService struct {
	logger      *logrus.Logger
	userRepo    UserRepository
	passwordKey string
}

func NewUserService(logger *logrus.Logger, userRepo UserRepository) *UserService {
	return &UserService{
		logger:   logger,
		userRepo: userRepo,
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
		token, err := jwttoken.GenerateToken(jwttoken.TokenFileds{
			ID:   res.ID,
			Role: res.Role,
		}, service.passwordKey)
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
