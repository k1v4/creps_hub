package repository

import "auth_service/pkg/DB/postgres"

type AuthRepository struct {
	db *postgres.DB
}

func NewAuthRepository(db *postgres.DB) *AuthRepository {
	return &AuthRepository{db}
}
