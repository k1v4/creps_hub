package usecase

import (
	"context"
	"release_service/internal/entity"
	"time"
)

type IReleaseRepository interface {
	GetReleasesByMonth(ctx context.Context, month string) ([]entity.Release, error)
	GetRelease(ctx context.Context, id int) (entity.Release, error)
	DeleteRelease(ctx context.Context, id int) error
	UpdateRelease(ctx context.Context, release entity.Release) (entity.Release, error)
	AddRelease(ctx context.Context, name string, date time.Time, imageUrl string) (int, error)
}

type IReleaseUseCase interface {
	AddRelease(ctx context.Context, name string, releaseDate time.Time, imageName string, imageData []byte) (int, error)
	DeleteRelease(ctx context.Context, id int) (bool, error)
	GetRelease(ctx context.Context, id int) (entity.Release, error)
	UpdateRelease(ctx context.Context, id int, name string, releaseDate time.Time) (entity.Release, error)
	GetReleasesByMonth(ctx context.Context, month string) ([]entity.Release, error)
}
