package repository

import (
	"context"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"release_service/internal/entity"
	"release_service/pkg/DataBase/postgres"
	"time"
)

type ReleaseRepository struct {
	*postgres.Postgres
}

func NewReleaseRepository(pg *postgres.Postgres) *ReleaseRepository {
	return &ReleaseRepository{pg}
}

func (r *ReleaseRepository) AddRelease(ctx context.Context, name string, date time.Time) (int, error) {
	const op = "repository.AddRelease"

	s, args, err := r.Builder.Insert("releases").
		Columns("name", "date", "release").
		Values(name, date, date).Suffix("RETURNING id").
		ToSql()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	var id int
	err = r.Pool.QueryRow(ctx, s, args...).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (r *ReleaseRepository) UpdateRelease(ctx context.Context, release entity.Release) (entity.Release, error) {
	const op = "repository.UpdateRelease"

	s, args, err := r.Builder.Update("releases").
		Set("name", release.Name).
		Set("date", release.ReleaseDate).
		Where(sq.Eq{"id": release.Id}).
		ToSql()
	if err != nil {
		return entity.Release{}, fmt.Errorf("%s: %w", op, err)
	}

	_, err = r.Pool.Exec(ctx, s, args...)
	if err != nil {
		return entity.Release{}, fmt.Errorf("%s: %w", op, err)
	}

	return release, nil
}

func (r *ReleaseRepository) DeleteRelease(ctx context.Context, id int) error {
	const op = "repository.DeleteRelease"

	s, args, err := r.Builder.Delete("releases").
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = r.Pool.Exec(ctx, s, args...)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *ReleaseRepository) GetRelease(ctx context.Context, id int) (entity.Release, error) {
	const op = "repository.GetRelease"

	s, args, err := r.Builder.Select("*").
		From("releases").
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return entity.Release{}, fmt.Errorf("%s: %w", op, err)
	}

	var release entity.Release
	err = r.Pool.QueryRow(ctx, s, args...).Scan(&release.Id, &release.Name, &release.ReleaseDate)
	if err != nil {
		return entity.Release{}, fmt.Errorf("%s: %w", op, err)
	}

	return release, nil
}

func (r *ReleaseRepository) GetReleasesByMonth(ctx context.Context, month string) ([]entity.Release, error) {
	const op = "repository.GetReleasesByMonth"

	year := time.Now().Year()
	startOfMonth, err := time.Parse("January 2006", fmt.Sprintf("%s %d", month, year))
	if err != nil {
		return nil, fmt.Errorf("%s: invalid month name: %w", op, err)
	}
	endOfMonth := startOfMonth.AddDate(0, 1, -1)

	s, args, err := r.Builder.
		Select("id", "date", "name").
		From("releases").
		Where(sq.And{
			sq.GtOrEq{"date": startOfMonth},
			sq.LtOrEq{"date": endOfMonth},
		}).
		OrderBy("date").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	rows, err := r.Pool.Query(ctx, s, args...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var releases []entity.Release
	for rows.Next() {
		var release entity.Release
		if err = rows.Scan(&release.Id, &release.ReleaseDate, &release.Name); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		releases = append(releases, release)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return releases, nil
}
