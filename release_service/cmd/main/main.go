package main

import (
	"context"
	"fmt"
	uploaderv1 "github.com/k1v4/protos/gen/file_uploader"
	"github.com/labstack/echo/v4"
	"google.golang.org/grpc"
	"os"
	"os/signal"
	"release_service/internal/config"
	v1 "release_service/internal/controller/http/v1"
	"release_service/internal/usecase"
	"release_service/internal/usecase/repository"
	"release_service/pkg/DataBase/postgres"
	"release_service/pkg/DataBase/redis"
	"release_service/pkg/httpserver"
	"release_service/pkg/logger"
	"strconv"
	"syscall"
)

func main() {
	ctx := context.Background()

	releaseLogger := logger.NewLogger()
	ctx = context.WithValue(ctx, logger.LoggerKey, releaseLogger)

	cfg := config.MustLoadConfig()
	if cfg == nil {
		panic("load config fail")
	}

	releaseLogger.Info(ctx, "read config successfully")

	url := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DBConfig.UserName,
		cfg.DBConfig.Password,
		cfg.DBConfig.Host,
		cfg.DBConfig.Port,
		cfg.DBConfig.DbName,
	)

	pg, err := postgres.New(url, postgres.MaxPoolSize(cfg.DBConfig.PoolMax))
	if err != nil {
		releaseLogger.Error(ctx, fmt.Sprintf("app - Run - postgres.New: %s", err))
	}
	defer pg.Close()

	releaseLogger.Info(ctx, "connected to database successfully")

	authRepo := repository.NewReleaseRepository(pg)

	client, err := createUploaderClient(cfg.UploaderGRPCServerPort)
	if err != nil {
		releaseLogger.Error(ctx, err.Error())
	}

	clientRedis, err := redis.NewClient(ctx, cfg.RedisConfig)
	if err != nil {
		releaseLogger.Error(ctx, "redis client init fail")
	}

	authUseCase := usecase.NewReleaseUseCase(authRepo, client, clientRedis)

	handler := echo.New()

	v1.NewRouter(handler, releaseLogger, authUseCase)

	httpServer := httpserver.New(handler, httpserver.Port(strconv.Itoa(cfg.RestServerPort)))

	// signal for graceful shutdown
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		releaseLogger.Info(ctx, "app-Run-signal: "+s.String())
	case err = <-httpServer.Notify():
		releaseLogger.Error(ctx, fmt.Sprintf("app-Run-httpServer.Notify: %s", err))
	}

	// shutdown
	err = httpServer.Shutdown()
	if err != nil {
		releaseLogger.Error(ctx, fmt.Sprintf("app-Run-httpServer.Shutdown: %s", err))
	}

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
