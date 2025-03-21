package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	uploaderv1 "github.com/k1v4/protos/gen/file_uploader"
	"github.com/redis/go-redis/v9"
	"shoe_service/internal/models"
	DataBase "shoe_service/pkg/DB"
	"time"
)

var (
	ErrShoeNotFound = errors.New("shoe not found")
)

type ShoeService struct {
	ShoeProv ShoeProvider
	client   uploaderv1.FileUploaderClient
	cache    *redis.Client
}

type ShoeProvider interface {
	AddShoe(ctx context.Context, userID int64, name, imageUrl string) (int64, error)
	GetShoe(ctx context.Context, shoeId int64) (*models.Shoe, error)
	RemoveShoe(ctx context.Context, shoeID int64) (bool, error)
	UpdateShoe(ctx context.Context, shoeId, userId int64, name, imageUrl string) (*models.Shoe, error)
	GetShoes(ctx context.Context, userId int64) (*[]models.Shoe, error)
}

func NewShoeService(shoeProvider ShoeProvider, client uploaderv1.FileUploaderClient, cache *redis.Client) *ShoeService {
	return &ShoeService{
		ShoeProv: shoeProvider,
		client:   client,
		cache:    cache,
	}
}

func (s *ShoeService) AddShoe(ctx context.Context, userID int64, shoeName, fileName string, imageData []byte) (int64, error) {
	const op = "ShoeService.AddShoe"

	image, err := s.client.UploadFile(ctx, &uploaderv1.ImageUploadRequest{
		ImageData: imageData,
		FileName:  fileName,
	})
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	imageUrl := image.GetUrl()
	if len(imageUrl) == 0 {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	shoeId, err := s.ShoeProv.AddShoe(ctx, userID, shoeName, imageUrl)
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
		if errors.Is(err, DataBase.ErrShoeNotFound) {
			return nil, ErrShoeNotFound
		}

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return shoe, nil
}

func (s *ShoeService) DeleteShoe(ctx context.Context, shoeID int64) (bool, error) {
	const op = "ShoeService.DeleteShoe"

	shoe, err := s.ShoeProv.GetShoe(ctx, shoeID)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	_, err = s.client.DeleteFile(ctx, &uploaderv1.ImageDeleteRequest{
		Url: shoe.ImageUrl,
	})
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	ok, err := s.ShoeProv.RemoveShoe(ctx, shoeID)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return ok, nil
}

func (s *ShoeService) UpdateShoe(ctx context.Context, shoeId, userId int64, name string, imageData []byte) (*models.Shoe, error) {
	const op = "ShoeService.UpdateShoe"

	// текущее состояние для удаления фото
	getShoe, err := s.ShoeProv.GetShoe(ctx, shoeId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// удаление фото
	_, err = s.client.DeleteFile(ctx, &uploaderv1.ImageDeleteRequest{
		Url: getShoe.ImageUrl,
	})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// загрузка фото
	uploadFile, err := s.client.UploadFile(ctx, &uploaderv1.ImageUploadRequest{
		ImageData: imageData,
		FileName:  name,
	})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	image := uploadFile.GetUrl()
	if len(image) == 0 {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	shoe, err := s.ShoeProv.UpdateShoe(ctx, shoeId, userId, name, image)
	if err != nil {
		if errors.Is(err, DataBase.ErrShoeNotFound) {
			return nil, ErrShoeNotFound
		}

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return shoe, nil
}

func (s *ShoeService) GetShoes(ctx context.Context, userId int64) (*[]models.Shoe, error) {
	const op = "ShoeService.GetShoes"
	var shoesRedis []models.Shoe

	data, err := s.cache.Get(ctx, fmt.Sprintf("%d", userId)).Bytes()
	if err == nil {
		errMarshal := json.Unmarshal(data, &shoesRedis)
		if errMarshal != nil {
			return nil, fmt.Errorf("failed to unmarshal articles: %w", err)
		}

		return &shoesRedis, nil
	} else if err != nil && !errors.Is(err, redis.Nil) {
		return nil, fmt.Errorf("failed to get articles from Redis: %w", err)
	}

	shoes, err := s.ShoeProv.GetShoes(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	dataSet, err := json.Marshal(shoes)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal articles: %w", err)
	}

	s.cache.Set(ctx, fmt.Sprintf("%d", userId), dataSet, 10*time.Minute)

	return shoes, nil
}
