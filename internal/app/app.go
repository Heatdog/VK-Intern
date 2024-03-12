package app

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Heater_dog/Vk_Intern/internal/config"
	migrations "github.com/Heater_dog/Vk_Intern/internal/migration"
	"github.com/Heater_dog/Vk_Intern/internal/transport"
	"github.com/Heater_dog/Vk_Intern/internal/user"
	userDb "github.com/Heater_dog/Vk_Intern/internal/user/db"
	"github.com/Heater_dog/Vk_Intern/pkg/client/postgre"
	"github.com/sirupsen/logrus"
)

func App() {
	logger := logrus.New()
	ctx := context.Background()
	logger.Info("Reading server config files")
	cfg := config.NewConfigStorage(logger)

	logger.Info("Connecting to DataBase")
	dbClient, err := postgre.NewClient(ctx, cfg.Postgre)
	if err != nil {
		logger.Fatalf("Connection to PostgreSQL error: %s", err.Error())
	}
	defer dbClient.Close(ctx)

	if err := migrations.InitDb(dbClient); err != nil {
		logger.Fatalf("Init db error: %s", err.Error())
	}

	mux := http.NewServeMux()

	userRepo := userDb.NewUserPostgreRepository(dbClient, logger)
	userService := user.NewUserService(logger, userRepo)
	transport.NewAuthHandler(logger, userService).Register(mux)

	dsn := fmt.Sprintf("%s:%s", cfg.Server.IP, cfg.Server.Port)

	logger.Infof("Listening server on: %s", dsn)
	if err := http.ListenAndServe(dsn, mux); err != nil {
		logger.Fatalf("Server listnenig error: %s", err.Error())
	}
}
