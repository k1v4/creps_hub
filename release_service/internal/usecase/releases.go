package usecase

import (
	"context"
	"fmt"
	"release_service/internal/entity"
	"release_service/pkg/translate"
	"strings"
	"time"
)

type ReleaseUseCase struct {
	repo IReleaseRepository
}

func NewReleaseUseCase(repo IReleaseRepository) *ReleaseUseCase {
	return &ReleaseUseCase{repo: repo}
}

func (r *ReleaseUseCase) AddRelease(ctx context.Context, name string, releaseDate time.Time) (int, error) {
	const op = "Usecase.AddRelease"

	releaseId, err := r.repo.AddRelease(ctx, name, releaseDate)
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

	release, err := r.repo.GetRelease(ctx, id)
	if err != nil {
		return entity.Release{}, fmt.Errorf("%s: %w", op, err)
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

	engMonth := translate.Translate(strings.ToLower(month))

	byMonth, err := r.repo.GetReleasesByMonth(ctx, engMonth)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return byMonth, nil
}
