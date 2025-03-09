package DataBase

import "errors"

var (
	ErrUserExists   = errors.New("user already exists")
	ErrShoeNotFound = errors.New("shoe not found")
)
