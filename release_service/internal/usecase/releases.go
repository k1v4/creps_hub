package usecase

import (
	"context"
	"errors"
	"fmt"
	uploaderv1 "github.com/k1v4/protos/gen/file_uploader"
	"github.com/redis/go-redis/v9"
	"release_service/internal/entity"
	"strings"
	"time"
)

type ReleaseUseCase struct {
	repo   IReleaseRepository
	client uploaderv1.FileUploaderClient
	cache  *redis.Client
}

func NewReleaseUseCase(repo IReleaseRepository, client uploaderv1.FileUploaderClient, cache *redis.Client) *ReleaseUseCase {
	return &ReleaseUseCase{
		repo:   repo,
		client: client,
		cache:  cache,
	}
}

func (r *ReleaseUseCase) AddRelease(ctx context.Context, name string, releaseDate time.Time, imageName string, imageData []byte) (int, error) {
	const op = "Usecase.AddRelease"

	image, err := r.client.UploadFile(ctx, &uploaderv1.ImageUploadRequest{
		ImageData: imageData,
		FileName:  imageName,
	})
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	releaseId, err := r.repo.AddRelease(ctx, name, releaseDate, image.GetUrl())
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return releaseId, nil
}

func (r *ReleaseUseCase) DeleteRelease(ctx context.Context, id int) (bool, error) {
	const op = "Usecase.DeleteRelease"

	err := r.repo.DeleteRelease(ctx, id)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return true, nil
}

func (r *ReleaseUseCase) GetRelease(ctx context.Context, id int) (entity.Release, error) {
	const op = "Usecase.GetRelease"
	var release entity.Release

	err := r.cache.Get(ctx, fmt.Sprintf("%d", id)).Scan(&release)
	if err == nil {
		return release, nil
	} else if err != nil && !errors.Is(err, redis.Nil) {
		return entity.Release{}, fmt.Errorf("%s: %w", op, err)
	}

	release, err = r.repo.GetRelease(ctx, id)
	if err != nil {
		return entity.Release{}, fmt.Errorf("%s: %w", op, err)
	}

	statusCmd := r.cache.Set(ctx, fmt.Sprintf("%d", id), &release, 1*time.Hour)
	if statusCmd.Err() != nil {
		return entity.Release{}, fmt.Errorf("%s: %w", op, statusCmd.Err())
	}

	return release, nil
}

func (r *ReleaseUseCase) UpdateRelease(ctx context.Context, id int, name string, releaseDate time.Time) (entity.Release, error) {
	const op = "Usecase.UpdateRelease"

	newRelease, err := r.repo.UpdateRelease(ctx, entity.Release{
		Id:          id,
		Name:        name,
		ReleaseDate: releaseDate},
	)
	if err != nil {

		return entity.Release{}, fmt.Errorf("%s: %w", op, err)
	}

	return newRelease, nil
}

func (r *ReleaseUseCase) GetReleasesByMonth(ctx context.Context, month string) ([]entity.Release, error) {
	const op = "Usecase.GetReleasesByMonth"

	engMonth := strings.ToLower(month)

	byMonth, err := r.repo.GetReleasesByMonth(ctx, engMonth)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return byMonth, nil
}
