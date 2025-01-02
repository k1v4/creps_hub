package service

import (
	"auth_service/internal/lib/jwt"
	"auth_service/internal/models"
	"auth_service/pkg/DataBase"
	"context"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"time"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidAppId       = errors.New("invalid app id")
	ErrUserExist          = errors.New("user exist")
)

type AuthService struct {
	UsSaver  UserSaver
	UsProv   UserProvider
	AppProv  AppProvider
	TokenTTL time.Duration
}

type UserSaver interface {
	SaveUser(ctx context.Context, email string, password []byte) (int64, error)
}

type UserProvider interface {
	GetUser(ctx context.Context, email string) (*models.User, error)
	IsAdmin(ctx context.Context, id int64) (bool, error)
}

type AppProvider interface {
	GetApp(ctx context.Context, id int64) (*models.App, error)
}

func NewAuthService(usSaver UserSaver, usProv UserProvider, appProv AppProvider, tokenTTL time.Duration) *AuthService {
	return &AuthService{
		UsSaver:  usSaver,
		UsProv:   usProv,
		AppProv:  appProv,
		TokenTTL: tokenTTL,
	}
}

func (s *AuthService) Login(ctx context.Context, email string, password string, appId int64) (string, error) {
	const op = "service.Login"

	user, err := s.UsProv.GetUser(ctx, email)
	if err != nil {
		//TODO доп проверки
		return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	if err = bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	app, err := s.AppProv.GetApp(ctx, appId)
	if err != nil {
		//TODO доп проверки
		return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	token, err := jwt.NewToken(*user, *app, s.TokenTTL)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return token, nil
}

func (s *AuthService) Register(ctx context.Context, email string, password string) (int64, error) {
	const op = "service.Register"

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := s.UsSaver.SaveUser(ctx, email, passHash)
	if err != nil {
		//TODO доп проверки
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

// IsAdmin checks if user is admin.
func (s *AuthService) IsAdmin(ctx context.Context, userID int64) (bool, error) {
	const op = "service.IsAdmin"

	isAdmin, err := s.UsProv.IsAdmin(ctx, userID)
	if err != nil {
		if errors.Is(err, DataBase.ErrAppNotFound) {

			return false, fmt.Errorf("%s: %w", op, ErrInvalidAppId)
		}

		return false, fmt.Errorf("%s: %w", op, err)
	}

	return isAdmin, nil
}
