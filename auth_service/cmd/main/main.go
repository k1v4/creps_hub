package main

import (
	"auth_service/internal/config"
	v1 "auth_service/internal/controller/http/v1"
	"auth_service/internal/usecase"
	"auth_service/internal/usecase/repository"
	"auth_service/pkg/DataBase/postgres"
	"auth_service/pkg/httpserver"
	"auth_service/pkg/logger"
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"os"
	"os/signal"
	"strconv"
	"syscall"
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

	url := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DBConfig.UserName,
		cfg.DBConfig.Password,
		cfg.DBConfig.Host,
		cfg.DBConfig.Port,
		cfg.DBConfig.DbName,
	)

	pg, err := postgres.New(url, postgres.MaxPoolSize(cfg.DBConfig.PoolMax))
	if err != nil {
		authLogger.Error(ctx, fmt.Sprintf("app - Run - postgres.New: %s", err))
	}
	defer pg.Close()

	authLogger.Info(ctx, "connected to database successfully")

	authRepo := repository.NewAuthRepository(pg)

	authUseCase := usecase.NewAuthUseCase(authRepo, cfg.TokenTTL, cfg.RefreshTokenTTL)

	handler := echo.New()

	v1.NewRouter(handler, authLogger, authUseCase)

	httpServer := httpserver.New(handler, httpserver.Port(strconv.Itoa(cfg.RestServerPort)))

	// signal for graceful shutdown
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		authLogger.Info(ctx, "app-Run-signal: "+s.String())
	case err = <-httpServer.Notify():
		authLogger.Error(ctx, fmt.Sprintf("app-Run-httpServer.Notify: %s", err))
	}

	// shutdown
	err = httpServer.Shutdown()
	if err != nil {
		authLogger.Error(ctx, fmt.Sprintf("app-Run-httpServer.Shutdown: %s", err))
	}

}
