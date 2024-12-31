package main

import (
	"auth_service/pkg/logger"
	"context"
)

func main() {
	ctx := context.Background()

	authLogger := logger.NewLogger()
	ctx = context.WithValue(ctx, logger.LoggerKey, authLogger)

	//TODO config

	//TODO storage

	//TODO repository

	//TODO service

	//TODO server

	//TODO graceful shutdown

}
