package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"user_service/internal/models"
	DataBase "user_service/pkg/DB"
	"user_service/pkg/DB/postgres"
)

type ShoeRepository struct {
	db *postgres.DB
}

func NewShoeRepository(db *postgres.DB) *ShoeRepository {
	return &ShoeRepository{
		db: db,
	}
}

func (s *ShoeRepository) AddShoe(ctx context.Context, userID int64, name, imageUrl string) (int64, error) {
	const op = "ShoeRepository.AddShoe"

	var shoeId int64

	err := sq.Insert("shoes").
		Columns("name", "image_url", "user_id").
		Values(name, imageUrl, userID).
		Suffix("RETURNING id").
		PlaceholderFormat(sq.Dollar).
		RunWith(s.db.Db).
		QueryRow().
		Scan(&shoeId)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return shoeId, nil
}

func (s *ShoeRepository) GetShoe(ctx context.Context, shoeId int64) (*models.Shoe, error) {
	const op = "ShoeRepository.GetShoes"

	var shoe models.Shoe

	err := sq.Select("*").
		From("shoes").
		Where(sq.Eq{"id": shoeId}).
		RunWith(s.db.Db).
		QueryRow().
		Scan(&shoe.Id, &shoe.Name, &shoe.ImageUrl, &shoe.UserId)
	if err != nil {
		//TODO доп проверки
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &shoe, nil
}

func (s *ShoeRepository) RemoveShoe(ctx context.Context, shoeID int64) (bool, error) {
	const op = "ShoeRepository.RemoveShoe"

	_, err := sq.Delete("shoes").
		Where(sq.Eq{"id": shoeID}).
		PlaceholderFormat(sq.Dollar).
		RunWith(s.db.Db).
		Query()
	if err != nil {
		//TODO доп проверки

		return false, fmt.Errorf("%s: %w", op, err)
	}

	return true, nil
}

func (s *ShoeRepository) UpdateShoe(ctx context.Context, shoeId, userId int64, name, imageUrl string) (*models.Shoe, error) {
	const op = "ShoeRepository.UpdateShoes"

	_, err := sq.Update("shoes").
		Set("name", name).
		Set("image_url", imageUrl).
		Where(sq.Eq{"id": shoeId}).
		Where(sq.Eq{"user_id": userId}).
		PlaceholderFormat(sq.Dollar).
		RunWith(s.db.Db).
		Query()
	if err != nil {
		//TODO доп проверки
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf(DataBase.ErrNoUser)
		}

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &models.Shoe{
		Id:       shoeId,
		Name:     name,
		ImageUrl: imageUrl,
		UserId:   userId,
	}, nil
}

func (s *ShoeRepository) GetShoes(ctx context.Context, userId int64) (*[]models.Shoe, error) {
	const op = "ShoeRepository.GetShoes"

	var shoes []models.Shoe

	rows, err := sq.Select("*").
		From("shoes").
		Where(sq.Eq{"user_id": userId}).
		PlaceholderFormat(sq.Dollar).
		RunWith(s.db.Db).
		Query()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	for rows.Next() {
		var shoe models.Shoe
		err = rows.Scan(&shoe.Id, &shoe.Name, &shoe.ImageUrl, &shoe.UserId)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		shoes = append(shoes, shoe)
	}

	return &shoes, nil
}
