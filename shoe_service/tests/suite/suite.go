package suite

import (
	"context"
	userv1 "github.com/k1v4/protos/gen/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net"
	"strconv"
	"testing"
	"time"
	"user_service/internal/config"
)

type Suite struct {
	*testing.T
	Cfg        *config.Config
	UserClient userv1.UserServiceClient
}

const (
	grpcHost = "localhost"
)

func New(t *testing.T) (context.Context, *Suite) {
	t.Helper()
	t.Parallel()

	cfg := config.MustLoadConfig()

	ctx, cancelCtx := context.WithTimeout(context.Background(), 100*time.Second)

	t.Cleanup(func() {
		t.Helper()
		cancelCtx()
	})

	cc, err := grpc.DialContext(
		context.Background(),
		grpcAddress(cfg),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		t.Fatalf("grpc server connection failed: %v", err)
	}

	return ctx, &Suite{
		T:          t,
		Cfg:        cfg,
		UserClient: userv1.NewUserServiceClient(cc),
	}
}

func grpcAddress(cfg *config.Config) string {
	return net.JoinHostPort(grpcHost, strconv.Itoa(cfg.GRPCServerPort))
}
