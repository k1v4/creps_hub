package main

import (
	"article_service/internal/config"
	v1 "article_service/internal/controller/http/v1"
	"article_service/internal/usecase"
	"article_service/internal/usecase/repository"
	"article_service/pkg/DataBase/postgres"
	"article_service/pkg/httpserver"
	"article_service/pkg/logger"
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

	articleLogger := logger.NewLogger()
	ctx = context.WithValue(ctx, logger.LoggerKey, articleLogger)

	cfg := config.MustLoadConfig()
	if cfg == nil {
		panic("load config fail")
	}

	articleLogger.Info(ctx, "read config successfully")

	url := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DBConfig.UserName,
		cfg.DBConfig.Password,
		cfg.DBConfig.Host,
		cfg.DBConfig.Port,
		cfg.DBConfig.DbName,
	)

	pg, err := postgres.New(url, postgres.MaxPoolSize(cfg.DBConfig.PoolMax))
	if err != nil {
		articleLogger.Error(ctx, fmt.Sprintf("app - Run - postgres.New: %s", err))
	}
	defer pg.Close()

	articleLogger.Info(ctx, "connected to database successfully")

	authRepo := repository.NewArticleRepository(pg)

	authUseCase := usecase.NewArticleUseCase(authRepo)

	handler := echo.New()

	v1.NewRouter(handler, articleLogger, authUseCase)

	httpServer := httpserver.New(handler, httpserver.Port(strconv.Itoa(cfg.RestServerPort)))

	// signal for graceful shutdown
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		articleLogger.Info(ctx, "app-Run-signal: "+s.String())
	case err = <-httpServer.Notify():
		articleLogger.Error(ctx, fmt.Sprintf("app-Run-httpServer.Notify: %s", err))
	}

	// shutdown
	err = httpServer.Shutdown()
	if err != nil {
		articleLogger.Error(ctx, fmt.Sprintf("app-Run-httpServer.Shutdown: %s", err))
	}

}
