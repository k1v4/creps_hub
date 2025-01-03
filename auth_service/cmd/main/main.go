package main

import (
	"auth_service/internal/config"
	"auth_service/internal/repository"
	"auth_service/internal/service"
	"auth_service/internal/transport/grpc"
	"auth_service/pkg/DataBase/postgres"
	"auth_service/pkg/logger"
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx := context.Background()

	authLogger := logger.NewLogger()
	ctx = context.WithValue(ctx, logger.LoggerKey, authLogger)

	cfg := config.MustLoadConfig()
	fmt.Println(cfg)
	if cfg == nil {
		panic("load config fail")
	}

	authLogger.Info(ctx, "read config successfully")

	storage, err := postgres.New(cfg.DBConfig)
	if err != nil {
		panic(err)
	}

	authRepo := repository.NewAuthRepository(storage)

	authServ := service.NewAuthService(authRepo, authRepo, authRepo, cfg.TokenTTL)

	grpcServer, err := grpc.NewServer(ctx, cfg.GRPCServerPort, cfg.RestServerPort, authServ)
	if err != nil {
		authLogger.Error(ctx, err.Error())
		return
	}

	graceCh := make(chan os.Signal, 1)
	signal.Notify(graceCh, syscall.SIGINT, syscall.SIGTERM)

	// запуск сервера
	go func() {
		if err = grpcServer.Start(ctx); err != nil {
			authLogger.Error(ctx, err.Error())
		}
	}()

	<-graceCh

	err = grpcServer.Stop(ctx)
	if err != nil {
		authLogger.Error(ctx, err.Error())
	}
	authLogger.Info(ctx, "Server stopped")
	fmt.Println("Server stopped")

	//TODO написать доки к функциям
}
