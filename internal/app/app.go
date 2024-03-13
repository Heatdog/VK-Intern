package app

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/Heater_dog/Vk_Intern/internal/auth"
	authDb "github.com/Heater_dog/Vk_Intern/internal/auth/db"
	"github.com/Heater_dog/Vk_Intern/internal/config"
	migrations "github.com/Heater_dog/Vk_Intern/internal/migration"
	"github.com/Heater_dog/Vk_Intern/internal/transport"
	"github.com/Heater_dog/Vk_Intern/internal/user"
	userDb "github.com/Heater_dog/Vk_Intern/internal/user/db"
	"github.com/Heater_dog/Vk_Intern/pkg/client/postgre"
	redisStorage "github.com/Heater_dog/Vk_Intern/pkg/client/redis"
)

func App() {
	opt := &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, opt))
	slog.SetDefault(logger)

	ctx := context.Background()
	logger.Info("reading server config files")
	cfg := config.NewConfigStorage(logger)

	logger.Info("connecting to DataBase")
	dbClient, err := postgre.NewPostgreClient(ctx, cfg.Postgre)
	if err != nil {
		logger.Error("connection to PostgreSQL failed", slog.Any("error", err))
	}
	defer dbClient.Close(ctx)

	logger.Info("init database with users")
	if err := migrations.InitDb(dbClient); err != nil {
		logger.Error("init db failed", slog.Any("error", err))
	}

	logger.Info("connecting to TokenDataBase")
	redisClient, err := redisStorage.NewRedisClient(ctx, &cfg.Redis)
	if err != nil {
		logger.Error("connection to Redis failed", slog.Any("error", err))
	}
	defer redisClient.Close()

	mux := http.NewServeMux()

	tokenStorage := authDb.NewRedisTokenStorage(logger, redisClient)
	tokenService := auth.NewTokenService(logger, tokenStorage, cfg.PasswordKey, cfg.Redis.TokenExparation)

	userRepo := userDb.NewUserPostgreRepository(dbClient, logger)
	userService := user.NewUserService(logger, userRepo, tokenService)
	transport.NewAuthHandler(logger, userService).Register(mux)

	host := fmt.Sprintf("%s:%s", cfg.Server.IP, cfg.Server.Port)

	logger.Info("listening server", slog.String("host", host))
	if err := http.ListenAndServe(host, mux); err != nil {
		logger.Error("server listnenig failed", slog.Any("error", err))
	}
}
