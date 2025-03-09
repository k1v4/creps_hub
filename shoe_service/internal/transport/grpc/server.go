package grpc

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	shoev1 "github.com/k1v4/protos/gen/shoe"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net"
	"net/http"
	"user_service/pkg/logger"
)

type Server struct {
	grpcServer *grpc.Server
	httpServer *http.Server
	listener   net.Listener
}

func NewServer(ctx context.Context, grpcPort, restPort int, shoeService IShoeService) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	opts := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(
			ContextWithLogger(logger.GetLoggerFromContext(ctx)),
		),
	}

	grpcServer := grpc.NewServer(opts...)
	shoev1.RegisterShoeServiceServer(grpcServer, NewShoeService(shoeService))

	conn, err := grpc.NewClient(
		fmt.Sprintf(":%d", grpcPort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to dial server: %w", err)
	}

	gwmux := runtime.NewServeMux()
	if err = shoev1.RegisterShoeServiceHandler(ctx, gwmux, conn); err != nil {
		return nil, fmt.Errorf("failed to register user gateway: %w", err)
	}
	if err = shoev1.RegisterShoeServiceHandler(ctx, gwmux, conn); err != nil {
		return nil, fmt.Errorf("failed to register shoe gateway: %w", err)
	}

	gwServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", restPort),
		Handler: gwmux,
	}

	return &Server{grpcServer, gwServer, listener}, nil
}

func (s *Server) Start(ctx context.Context) error {
	eg := errgroup.Group{}

	eg.Go(func() error {
		lg := logger.GetLoggerFromContext(ctx)
		if lg != nil {
			lg.Info(ctx, "starting grpc server", zap.Int("port", s.listener.Addr().(*net.TCPAddr).Port))
		}

		return s.grpcServer.Serve(s.listener)
	})

	eg.Go(func() error {
		lg := logger.GetLoggerFromContext(ctx)
		if lg != nil {
			lg.Info(ctx, "starting rest server", zap.String("port", s.httpServer.Addr))
		}

		return s.httpServer.ListenAndServe()
	})

	return eg.Wait()
}

func (s *Server) Stop(ctx context.Context) error {
	s.grpcServer.GracefulStop()

	l := logger.GetLoggerFromContext(ctx)
	if l != nil {
		l.Info(ctx, "grpc server stopped")
	}

	return s.httpServer.Shutdown(ctx)
}
