package service

import (
	"context"
	"errors"
	"fmt"
	"user_service/internal/models"
	DataBase "user_service/pkg/DB"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserExist          = errors.New("user exist")
)

type UserService struct {
	UsPr UserProvider
	UsSa UserSaver
	UsBa UserBanner
}

type UserProvider interface {
	UpdateUser(ctx context.Context, id int64, name, surname, username string) (*models.User, error)
	GetUser(ctx context.Context, id int64) (*models.User, error)
}

type UserSaver interface {
	SaveUser(ctx context.Context, id int64, name, surname, username string) (int64, error)
}

type UserBanner interface {
	DeleteUser(ctx context.Context, id int64) (bool, error)
}

func NewUserService(usPr UserProvider, usSa UserSaver, usBa UserBanner) *UserService {
	return &UserService{
		UsPr: usPr,
		UsSa: usSa,
		UsBa: usBa,
	}
}

func (us *UserService) AddUser(ctx context.Context, id int64, name, surname, username string) (int64, error) {
	const op = "UserService.Register"

	userId, err := us.UsSa.SaveUser(ctx, id, name, surname, username)
	if err != nil {
		if errors.Is(err, DataBase.ErrUserExists) {
			return 0, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return userId, nil
}

func (us *UserService) GetUser(ctx context.Context, id int64) (*models.User, error) {
	const op = "UserService.GetUser"

	user, err := us.UsPr.GetUser(ctx, id)
	if err != nil {
		if errors.Is(err, DataBase.ErrUserExists) {
			return nil, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (us *UserService) UpdateUser(ctx context.Context, id int64, name, surname, username string) (*models.User, error) {
	const op = "UserService.UpdateUser"

	user, err := us.UsPr.UpdateUser(ctx, id, name, surname, username)
	if err != nil {
		if errors.Is(err, DataBase.ErrUserExists) {
			return nil, fmt.Errorf("%s: %w", op, ErrUserExist)
		}

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (us *UserService) DeleteUser(ctx context.Context, id int64) (bool, error) {
	const op = "UserService.DeleteUser"

	deleteRes, err := us.UsBa.DeleteUser(ctx, id)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return deleteRes, nil
}
