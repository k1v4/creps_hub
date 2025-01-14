package service

import (
	"context"
	"user_service/internal/models"
)

type ShoeService struct {
	ShoeProv ShoeProvider
}

type ShoeProvider interface {
	AddShoe(ctx context.Context, userID int64, name, imageUrl string) (int64, error)
	GetShoe(ctx context.Context, shoeId int64) (*models.Shoe, error)
	RemoveShoe(ctx context.Context, shoeID int64) (bool, error)
	UpdateShoe(ctx context.Context, shoeId, userId int64, name, imageUrl string) (*models.Shoe, error)
	GetShoes(ctx context.Context, userId int64) (*[]models.Shoe, error)
}

func NewShoeService(shoeProvider ShoeProvider) *ShoeService {
	return &ShoeService{shoeProvider}
}

func (s *ShoeService) AddShoe(ctx context.Context, userID int64, name, imageUrl string) (int64, error) {
}

func (s *ShoeService) GetShoe(ctx context.Context, shoeId int64) (*models.Shoe, error) {}

func (s *ShoeService) RemoveShoe(ctx context.Context, shoeID int64) (bool, error) {}

func (s *ShoeService) UpdateShoe(ctx context.Context, shoeId, userId int64, name, imageUrl string) (*models.Shoe, error) {
}

func (s *ShoeService) GetShoes(ctx context.Context, userId int64) (*[]models.Shoe, error) {}
