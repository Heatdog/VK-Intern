package app

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	_ "github.com/Heater_dog/Vk_Intern/docs"

	"github.com/Heater_dog/Vk_Intern/internal/config"
	migrations "github.com/Heater_dog/Vk_Intern/internal/migration"
	actor_db "github.com/Heater_dog/Vk_Intern/internal/repository/actor/db"
	film_db "github.com/Heater_dog/Vk_Intern/internal/repository/film/db"
	token_db "github.com/Heater_dog/Vk_Intern/internal/repository/token/db"
	user_db "github.com/Heater_dog/Vk_Intern/internal/repository/user/db"
	actor_service "github.com/Heater_dog/Vk_Intern/internal/services/actor"
	token_service "github.com/Heater_dog/Vk_Intern/internal/services/token"
	user_service "github.com/Heater_dog/Vk_Intern/internal/services/user"
	actor_transport "github.com/Heater_dog/Vk_Intern/internal/transport/actor"
	auth_transport "github.com/Heater_dog/Vk_Intern/internal/transport/auth"
	middleware_transport "github.com/Heater_dog/Vk_Intern/internal/transport/middleware"
	"github.com/Heater_dog/Vk_Intern/pkg/client/postgre"
	redisStorage "github.com/Heater_dog/Vk_Intern/pkg/client/redis"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

// @title Фильмотека
// @description API server for Фильмотека

// @host localhost:8080
// @BasePath /

// @securityDefinitions.apiKey ApiKeyAuth
// @in header
// @name Authorization

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
	defer dbClient.Close()

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

	tokenStorage := token_db.NewRedisTokenRepository(logger, redisClient)
	tokenService := token_service.NewTokenService(logger, tokenStorage, cfg.PasswordKey, cfg.Redis.TokenExparation)

	userRepo := user_db.NewUserPostgreRepository(dbClient, logger)
	userService := user_service.NewUserService(logger, userRepo, tokenService)
	auth_transport.NewAuthHandler(logger, userService).Register(mux)

	mid := middleware_transport.NewMiddleware(logger, tokenService, cfg.PasswordKey)

	filmRepo := film_db.NewFilmsPostgreRepository(dbClient, logger)

	actorRepo := actor_db.NewActorPostgreRepository(dbClient, logger)
	actorService := actor_service.NewActorsService(logger, actorRepo, filmRepo)
	actor_transport.NewActorsHandler(logger, actorService, mid).Register(mux)

	logger.Info("adding swagger documentation")
	mux.HandleFunc("/swagger/", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
		httpSwagger.DeepLinking(true),
	))

	host := fmt.Sprintf("%s:%s", cfg.Server.IP, cfg.Server.Port)

	logger.Info("listening server", slog.String("host", host))
	if err := http.ListenAndServe(host, mux); err != nil {
		logger.Error("server listnenig failed", slog.Any("error", err))
	}
}
