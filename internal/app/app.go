package app

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Heater_dog/Vk_Intern/internal/auth"
	authDb "github.com/Heater_dog/Vk_Intern/internal/auth/db"
	"github.com/Heater_dog/Vk_Intern/internal/config"
	migrations "github.com/Heater_dog/Vk_Intern/internal/migration"
	"github.com/Heater_dog/Vk_Intern/internal/transport"
	"github.com/Heater_dog/Vk_Intern/internal/user"
	userDb "github.com/Heater_dog/Vk_Intern/internal/user/db"
	"github.com/Heater_dog/Vk_Intern/pkg/client/postgre"
	redisStorage "github.com/Heater_dog/Vk_Intern/pkg/client/redis"
	"github.com/sirupsen/logrus"
)

func App() {
	logger := logrus.New()
	ctx := context.Background()
	logger.Info("Reading server config files")
	cfg := config.NewConfigStorage(logger)

	logger.Info("Connecting to DataBase")
	dbClient, err := postgre.NewPostgreClient(ctx, cfg.Postgre)
	if err != nil {
		logger.Fatalf("Connection to PostgreSQL error: %s", err.Error())
	}
	defer dbClient.Close(ctx)

	if err := migrations.InitDb(dbClient); err != nil {
		logger.Fatalf("Init db error: %s", err.Error())
	}

	logger.Info("Connecting to TokenDataBase")
	redisClient, err := redisStorage.NewRedisClient(ctx, &cfg.Redis)
	if err != nil {
		logger.Fatalf("Connection to Redis error: %s", err.Error())
	}
	defer redisClient.Close()

	mux := http.NewServeMux()

	tokenStorage := authDb.NewRedisTokenStorage(logger, redisClient)
	tokenService := auth.NewTokenService(logger, tokenStorage, cfg.PasswordKey, cfg.Redis.TokenExparation)

	userRepo := userDb.NewUserPostgreRepository(dbClient, logger)
	userService := user.NewUserService(logger, userRepo, tokenService)
	transport.NewAuthHandler(logger, userService).Register(mux)

	dsn := fmt.Sprintf("%s:%s", cfg.Server.IP, cfg.Server.Port)

	logger.Infof("Listening server on: %s", dsn)
	if err := http.ListenAndServe(dsn, mux); err != nil {
		logger.Fatalf("Server listnenig error: %s", err.Error())
	}
}
