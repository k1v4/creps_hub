package usecase

import (
	"auth_service/internal_rest/entity"
	"context"
)

type ISsoRepository interface {
	SaveUser(ctx context.Context, email string, password []byte) (int, error)
	GetUser(ctx context.Context, email string) (entity.User, error)
	GetUserById(ctx context.Context, id int) (entity.User, error)
	DeleteUser(ctx context.Context, id int) error
	UpdateUser(ctx context.Context, newUser entity.User) (entity.User, error)
}

type ISsoService interface{}
