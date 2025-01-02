package repository

import (
	"auth_service/internal/models"
	"auth_service/pkg/DataBase/postgres"
	"context"
	"fmt"
	sq "github.com/Masterminds/squirrel"
)

type AuthRepository struct {
	db *postgres.DB
}

func NewAuthRepository(db *postgres.DB) *AuthRepository {
	return &AuthRepository{db}
}

func (a *AuthRepository) SaveUser(ctx context.Context, email string, password []byte) (int64, error) {
	const op = "repository.SaveUser"

	var result models.User

	err := sq.Insert("users").
		Columns("email", "pass_hash").
		Values(email, password).
		Suffix("RETURNING id").
		PlaceholderFormat(sq.Dollar).
		RunWith(a.db.Db).
		QueryRow().
		Scan(&result.ID)
	if err != nil {
		return 0, fmt.Errorf("%s : %w", op, err)
	}

	return result.ID, nil
}

func (a *AuthRepository) GetUser(ctx context.Context, email string) (*models.User, error) {
	const op = "repository.GetUser"

	var result models.User
	var isAdmin bool

	err := sq.Select("*").
		From("users").
		Where(sq.Eq{"email": email}).
		PlaceholderFormat(sq.Dollar).
		RunWith(a.db.Db).
		QueryRow().
		Scan(&result.ID, &result.Email, &result.PassHash, &isAdmin)
	if err != nil {
		return nil, fmt.Errorf("%s : %w", op, err)
	}

	return &result, nil
}

func (a *AuthRepository) IsAdmin(ctx context.Context, id int64) (bool, error) {
	const op = "repository.IsAdmin"

	var isAdmin bool

	err := sq.Select("is_admin").
		From("users").
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar).
		RunWith(a.db.Db).
		QueryRow().
		Scan(&isAdmin)
	if err != nil {
		return false, fmt.Errorf("%s : %w", op, err)
	}

	return isAdmin, nil
}

func (a *AuthRepository) GetApp(ctx context.Context, id int64) (*models.App, error) {
	const op = "repository.App"
	var app models.App

	err := sq.Select("*").
		From("apps").
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar).
		RunWith(a.db.Db).
		QueryRow().
		Scan(&app.ID, &app.Name, &app.Secret)
	if err != nil {
		return nil, fmt.Errorf("%s : %w", op, err)
	}

	return &app, nil
}
