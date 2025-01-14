package grpc

import (
	"context"
	userv1 "github.com/k1v4/protos/gen/user"
	"user_service/internal/models"
)

type IShoeService interface {
	AddShoe(ctx context.Context, userID int64, name, imageUrl string) (int64, error)
	GetShoe(ctx context.Context, shoeId int64) (*models.Shoe, error)
	RemoveShoe(ctx context.Context, shoeID int64) (bool, error)
	UpdateShoe(ctx context.Context, shoeId, userId int64, name, imageUrl string) (*models.Shoe, error)
	GetShoes(ctx context.Context, userId int64) (*[]models.Shoe, error)
}

type ShoeService struct {
	userv1.UnimplementedShoeServiceServer
	service IShoeService
}

func NewShoeService(service IShoeService) *ShoeService {
	return &ShoeService{service: service}
}
