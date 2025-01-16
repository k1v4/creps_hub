package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/lib/pq"
	"user_service/internal/models"
	DataBase "user_service/pkg/DB"
	"user_service/pkg/DB/postgres"
)

type UserRepository struct {
	db *postgres.DB
}

func NewUserRepository(db *postgres.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (u *UserRepository) SaveUser(ctx context.Context, id int64, name, surname, username string) (int64, error) {
	const op = "UserRepository.SaveUser"

	var result models.User

	err := sq.Insert("users").
		Columns("id", "name", "surname", "username").
		Values(id, name, surname, username).
		Suffix("RETURNING id").
		PlaceholderFormat(sq.Dollar).
		RunWith(u.db.Db).
		QueryRow().
		Scan(&result.Id)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return 0, fmt.Errorf("%s: %w", op, DataBase.ErrUserExists)
		}

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return result.Id, nil
}

func (u *UserRepository) UpdateUser(ctx context.Context, id int64, name, surname, username string) (*models.User, error) {
	const op = "UserRepository.UpdateUser"

	_, err := sq.Update("users").
		Set("name", name).
		Set("surname", surname).
		Set("username", username).
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar).
		RunWith(u.db.Db).
		Query()
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return &models.User{}, fmt.Errorf("%s: %w", op, DataBase.ErrUserExists)
		}

		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf(DataBase.ErrNoUser)
		}

		return &models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return &models.User{
		Id:       id,
		Name:     name,
		Surname:  surname,
		UserName: username,
	}, nil
}

func (u *UserRepository) DeleteUser(ctx context.Context, id int64) (bool, error) {
	const op = "UserRepository.DeleteUser"

	_, err := sq.Delete("users").
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar).
		RunWith(u.db.Db).
		Query()
	if err != nil {
		return false, fmt.Errorf("repository.DeleteOrder: %w", err)
	}

	return true, nil
}

func (u *UserRepository) GetUser(ctx context.Context, id int64) (*models.User, error) {
	const op = "UserRepository.GetUser"

	var user models.User

	err := sq.Select("*").
		From("users").
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar).
		RunWith(u.db.Db).
		QueryRow().
		Scan(&user.Id, &user.Name, &user.Surname, &user.UserName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &models.User{}, fmt.Errorf("%s: %w", op, DataBase.ErrUserNotFound)
		}

		return &models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return &user, nil
}
