package DataBase

import "errors"

var (
	ErrUserExists   = errors.New("user already exists")
	ErrUserNotFound = errors.New("user not found")
	ErrNoUser       = "no user with this Id or username"
)
