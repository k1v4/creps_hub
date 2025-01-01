package main

import (
	"auth_service/internal/config"
	"auth_service/internal/repository"
	"auth_service/pkg/DB/postgres"
	"auth_service/pkg/logger"
	"context"
)

func main() {
	ctx := context.Background()

	authLogger := logger.NewLogger()
	ctx = context.WithValue(ctx, logger.LoggerKey, authLogger)

	cfg := config.MustLoadConfig()
	if cfg == nil {
		panic("load config fail")
	}

	authLogger.Info(ctx, "read config successfully")

	//TODO storage
	postgres, err := postgres.New(cfg.DBConfig)
	if err != nil {
		panic(err)
	}

	//TODO repository
	repo := repository.NewAuthRepository(postgres)

	//TODO service

	//TODO server

	//TODO graceful shutdown

}
