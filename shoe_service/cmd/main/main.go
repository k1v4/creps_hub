package main

import (
	"context"
	"fmt"
	uploaderv1 "github.com/k1v4/protos/gen/file_uploader"
	"google.golang.org/grpc"
	"os"
	"os/signal"
	"syscall"
	"user_service/internal/config"
	"user_service/internal/repository"
	"user_service/internal/service"
	grpc_transport "user_service/internal/transport/grpc"
	"user_service/pkg/DB/postgres"
	"user_service/pkg/logger"
)

func main() {
	ctx := context.Background()

	shoeLogger := logger.NewLogger()
	ctx = context.WithValue(ctx, logger.LoggerKey, shoeLogger)

	cfg := config.MustLoadConfig()
	if cfg == nil {
		panic("load config fail")
	}

	shoeLogger.Info(ctx, "read config successfully")

	storage, err := postgres.New(cfg.DBConfig)
	if err != nil {
		panic(err)
	}

	shoeRepo := repository.NewShoeRepository(storage)

	client, err := createUploaderClient(cfg.UploaderGRPCServerPort)
	if err != nil {
		shoeLogger.Error(ctx, err.Error())
	}

	shoeServ := service.NewShoeService(shoeRepo, client)

	grpcServer, err := grpc_transport.NewServer(ctx, cfg.GRPCServerPort, cfg.RestServerPort, shoeServ)
	if err != nil {
		shoeLogger.Error(ctx, err.Error())
		return
	}

	graceCh := make(chan os.Signal, 1)
	signal.Notify(graceCh, syscall.SIGINT, syscall.SIGTERM)

	// запуск сервера
	go func() {
		if err = grpcServer.Start(ctx); err != nil {
			shoeLogger.Error(ctx, err.Error())
		}
	}()

	<-graceCh

	err = grpcServer.Stop(ctx)
	if err != nil {
		shoeLogger.Error(ctx, err.Error())
	}
	shoeLogger.Info(ctx, "Server stopped")
	fmt.Println("Server stopped")

	//TODO написать доки к функциям
}

func createUploaderClient(port int) (uploaderv1.FileUploaderClient, error) {
	// TODO в конфиг uploader
	conn, err := grpc.Dial(fmt.Sprintf("uploader:%d", port), grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	client := uploaderv1.NewFileUploaderClient(conn)

	return client, nil
}
