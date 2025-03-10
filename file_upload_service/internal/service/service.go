package service

import (
	"context"
	"fmt"
)

type IUploadRepository interface {
	UploadImage(ctx context.Context, fileName string, imageData []byte) (string, error)
	DeleteImage(ctx context.Context, fileName string) error
}

type UploadServer struct {
	upRepo IUploadRepository
}

func NewUploadService(upRepo IUploadRepository) *UploadServer {
	return &UploadServer{
		upRepo: upRepo,
	}
}

func (s *UploadServer) UploadImage(ctx context.Context, fileName string, imageData []byte) (string, error) {
	const op = "service.GetShoes"

	urlImage, err := s.upRepo.UploadImage(ctx, fileName, imageData)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return urlImage, nil
}

func (s *UploadServer) DeleteImage(ctx context.Context, fileName string) (bool, error) {
	const op = "service.DeleteImage"

	err := s.upRepo.DeleteImage(ctx, fileName)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return true, nil
}
