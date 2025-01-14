package grpc

import (
	"context"
	"github.com/AlekSi/pointer"
	userv1 "github.com/k1v4/protos/gen/user"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"user_service/internal/models"
)

type IUserService interface {
	AddUser(ctx context.Context, id int64, name, surname, username string) (int64, error)
	GetUser(ctx context.Context, id int64) (*models.User, error)
	UpdateUser(ctx context.Context, id int64, name, surname, username string) (*models.User, error)
	DeleteUser(ctx context.Context, id int64) (bool, error)
}

type UserService struct {
	userv1.UnimplementedUserServiceServer
	service IUserService
}

func NewUserService(svc IUserService) *UserService {
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
		//TODO добавить доп проверки
		return nil, err
	}

	return &userv1.AddUserResponse{UserId: userId}, nil
}

func (s *UserService) GetUser(ctx context.Context, req *userv1.GetUserRequest) (*userv1.GetUserResponse, error) {
	userID := req.GetUserId()
	if userID <= 0 {
		return nil, status.Error(codes.InvalidArgument, "userId is wrong")
	}

	user, err := s.service.GetUser(ctx, userID)
	if err != nil {
		//TODO добавить доп проверки
		return nil, err
	}

	r := pointer.Get(user)

	//TODO добавить обращение к сервису для кроссовок
	return &userv1.GetUserResponse{
		User: &userv1.User{
			Id:       r.Id,
			Name:     r.Name,
			Surname:  r.Surname,
			Username: r.UserName,
			Shoes:    nil,
		},
	}, nil
}

func (s *UserService) UpdateUser(ctx context.Context, req *userv1.UpdateUserRequest) (*userv1.UpdateUserResponse, error) {
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

	user, err := s.service.UpdateUser(ctx, id, name, surname, username)
	if err != nil {
		//TODO добавить доп проверки
		return nil, err
	}

	r := pointer.Get(user)

	//TODO добавить обращение к сервису для кроссовок
	return &userv1.UpdateUserResponse{
		User: &userv1.User{
			Id:       r.Id,
			Name:     r.Name,
			Surname:  r.Surname,
			Username: r.UserName,
			Shoes:    nil,
		},
	}, nil
}

func (s *UserService) DeleteUser(ctx context.Context, req *userv1.DeleteUserRequest) (*userv1.DeleteUserResponse, error) {
	id := req.GetUserId()
	if id <= 0 {
		return nil, status.Error(codes.InvalidArgument, "userId is wrong")
	}

	IsSuccessfully, err := s.service.DeleteUser(ctx, id)
	if err != nil {
		//TODO добавить доп проверки
		return nil, err
	}

	return &userv1.DeleteUserResponse{
		IsSuccessfully: IsSuccessfully,
	}, nil
}
