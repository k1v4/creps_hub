package grpc

import (
	"context"
	userv1 "github.com/k1v4/protos/gen/user"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"user_service/internal/models"
)

type Service interface {
	AddUser(ctx context.Context, id int64, name, surname, username string) (int64, error)
	GetUser(ctx context.Context, username string) (*models.User, error)
	UpdateUser(ctx context.Context, id int64, name, surname, username string) (*models.User, error)
	DeleteUser(ctx context.Context, id int64) (bool, error)
}

type UserService struct {
	userv1.UnimplementedUserServiceServer
	service Service
}

func NewUserService(svc Service) *UserService {
	return &UserService{
		service: svc,
	}
}

func (s *UserService) AddUser(ctx context.Context, req *userv1.AddUserRequest) (*userv1.AddUserResponse, error) {
	name := req.GetName()
	if name == "" {
		return nil, status.Error(codes.InvalidArgument, "name is required")
	}

	surname := req.GetSurname()
	if surname == "" {
		return nil, status.Error(codes.InvalidArgument, "surname is required")
	}

	username := req.GetUsername()
	if username == "" {
		return nil, status.Error(codes.InvalidArgument, "username is required")
	}

	id := req.GetUserId()
	if id <= 0 {
		return nil, status.Error(codes.InvalidArgument, "userId is wrong")
	}

	userId, err := s.service.AddUser(ctx, id, name, surname, username)
	if err != nil {
		return nil, err
	}

}

func (s *UserService) GetUser(ctx context.Context, req *userv1.GetUserRequest) (*userv1.GetUserResponse, error) {
}

func (s *UserService) UpdateUser(ctx context.Context, req *userv1.UpdateUserRequest) (*userv1.UpdateUserResponse, error) {
}

func (s *UserService) DeleteUser(ctx context.Context, req *userv1.DeleteUserRequest) (*userv1.DeleteUserResponse, error) {
}
