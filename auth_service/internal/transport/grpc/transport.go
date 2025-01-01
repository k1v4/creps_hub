package grpc

import (
	"context"
	ssov1 "github.com/k1v4/protos/gen/sso"
)

type Service interface {
	Login(ctx context.Context, email string, password string, appId int64) (string, error)
	Register(ctx context.Context, email string, password string) (int64, error)
	IsAdmin(ctx context.Context, userID int64) (bool, error)
}

type AuthService struct {
	ssov1.UnimplementedAuthServer
	service Service
}

func NewAuthService(service Service) *AuthService {
	return &AuthService{service: service}
}
