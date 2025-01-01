package main

import (
	"auth_service/internal/config"
	"auth_service/internal/repository"
	"auth_service/internal/service"
	"auth_service/pkg/DataBase/postgres"
	"auth_service/pkg/logger"
	"context"
)

func main() {
	ctx := context.Background()

	authLogger := logger.NewLogger()
	ctx = context.WithValue(ctx, logger.LoggerKey, authLogger)

	//TODO добавить время жизни токена в кфг
	cfg := config.MustLoadConfig()
	if cfg == nil {
		panic("load config fail")
	}

	authLogger.Info(ctx, "read config successfully")

	storage, err := postgres.New(cfg.DBConfig)
	if err != nil {
		panic(err)
	}

	repo := repository.NewAuthRepository(storage)

	service := service.NewAuthService(repo, repo, repo, 0)

	//TODO server

	//TODO graceful shutdown

	//TODO написать доки к функциям
}
