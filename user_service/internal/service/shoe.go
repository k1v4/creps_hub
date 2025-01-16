package service

import (
	"context"
	"fmt"
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

//TODO когда будет сделано, добавить проверку авторизации

func (s *ShoeService) AddShoe(ctx context.Context, userID int64, name, imageUrl string) (int64, error) {
	const op = "ShoeService.AddShoe"

	shoeId, err := s.ShoeProv.AddShoe(ctx, userID, name, imageUrl)
	if err != nil {
		// TODO доп проверки

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return shoeId, nil
}

func (s *ShoeService) GetShoe(ctx context.Context, shoeId int64) (*models.Shoe, error) {
	const op = "ShoeService.GetShoe"

	shoe, err := s.ShoeProv.GetShoe(ctx, shoeId)
	if err != nil {
		// TODO доп проверки

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return shoe, nil
}

func (s *ShoeService) DeleteShoe(ctx context.Context, shoeID int64) (bool, error) {
	const op = "ShoeService.DeleteShoe"

	ok, err := s.ShoeProv.RemoveShoe(ctx, shoeID)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return ok, nil
}

func (s *ShoeService) UpdateShoe(ctx context.Context, shoeId, userId int64, name, imageUrl string) (*models.Shoe, error) {
	const op = "ShoeService.UpdateShoe"

	shoe, err := s.ShoeProv.UpdateShoe(ctx, shoeId, userId, name, imageUrl)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return shoe, nil
}

func (s *ShoeService) GetShoes(ctx context.Context, userId int64) (*[]models.Shoe, error) {
	const op = "ShoeService.GetShoes"

	shoes, err := s.ShoeProv.GetShoes(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return shoes, nil
}
