package grpc

import (
	"context"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"user_service/pkg/logger"
)

func ContextWithLogger(l logger.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		l.Info(ctx, "request started", zap.String("method", info.FullMethod), zap.String("handler", info.FullMethod))
		return handler(ctx, req)
	}
}
