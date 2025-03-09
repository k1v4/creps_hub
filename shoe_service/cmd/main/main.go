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

	userLogger := logger.NewLogger()
	ctx = context.WithValue(ctx, logger.LoggerKey, userLogger)

	cfg := config.MustLoadConfig()
	if cfg == nil {
		panic("load config fail")
	}

	userLogger.Info(ctx, "read config successfully")

	storage, err := postgres.New(cfg.DBConfig)
	if err != nil {
		panic(err)
	}

	shoeRepo := repository.NewShoeRepository(storage)

	shoeServ := service.NewShoeService(shoeRepo)

	grpcServer, err := grpc.NewServer(ctx, cfg.GRPCServerPort, cfg.RestServerPort, shoeServ)
	if err != nil {
		userLogger.Error(ctx, err.Error())
		return
	}

	graceCh := make(chan os.Signal, 1)
	signal.Notify(graceCh, syscall.SIGINT, syscall.SIGTERM)

	// запуск сервера
	go func() {
		if err = grpcServer.Start(ctx); err != nil {
			userLogger.Error(ctx, err.Error())
		}
	}()

	<-graceCh

	err = grpcServer.Stop(ctx)
	if err != nil {
		userLogger.Error(ctx, err.Error())
	}
	userLogger.Info(ctx, "Server stopped")
	fmt.Println("Server stopped")

	//TODO написать доки к функциям
}
