package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"user_service/internal/config"
	"user_service/internal/repository"
	"user_service/internal/service"
	"user_service/internal/transport/grpc"
	"user_service/pkg/DB/postgres"
	"user_service/pkg/logger"
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

	storage, err := postgres.New(cfg.DBConfig)
	if err != nil {
		panic(err)
	}

	userRepo := repository.NewUserRepository(storage)

	authServ := service.NewUserService(userRepo, userRepo, userRepo)

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
