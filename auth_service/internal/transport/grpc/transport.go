package grpc

import (
	"auth_service/internal/service"
	"context"
	"errors"
	"fmt"
	ssov1 "github.com/k1v4/protos/gen/sso"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidAppId       = errors.New("invalid app id")
	ErrUserExist          = errors.New("user exist")
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

func (s *AuthService) Login(ctx context.Context, req *ssov1.LoginRequest) (*ssov1.LoginResponse, error) {
	email := req.GetEmail()
	if len(email) == 0 {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}

	password := req.GetPassword()
	if len(password) == 0 {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}

	appId := int64(req.GetAppId())
	if appId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "app_id is invalid")
	}

	token, err := s.service.Login(ctx, email, password, appId)
	if err != nil {
		fmt.Println(err)
		if errors.Is(err, service.ErrInvalidCredentials) {
			return nil, status.Error(codes.InvalidArgument, "invalid credentials")
		}
		return nil, status.Error(codes.Internal, "internal server error")
	}

	return &ssov1.LoginResponse{Token: token}, err
}

func (s *AuthService) Register(ctx context.Context, req *ssov1.RegisterRequest) (*ssov1.RegisterResponse, error) {
	fmt.Println("transport.Register")

	email := req.GetEmail()
	if len(email) == 0 {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}

	password := req.GetPassword()
	if len(password) == 0 {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}

	register, err := s.service.Register(ctx, email, password)
	if err != nil {
		if errors.Is(err, service.ErrUserExist) {
			return nil, status.Error(codes.InvalidArgument, "user already exists")
		}

		return nil, status.Error(codes.Internal, "internal server error")
	}

	return &ssov1.RegisterResponse{UserId: register}, nil
}

func (s *AuthService) IsAdmin(ctx context.Context, req *ssov1.IsAdminRequest) (*ssov1.IsAdminResponse, error) {
	userId := int64(req.GetUserId())
	if userId < 0 {
		return nil, status.Error(codes.InvalidArgument, "user id is invalid")
	}

	isAdmin, err := s.service.IsAdmin(ctx, userId)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal server error")
	}

	return &ssov1.IsAdminResponse{IsAdmin: isAdmin}, nil
}
