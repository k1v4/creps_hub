package v1

import (
	"auth_service/internal/entity"
	"auth_service/internal/usecase"
	jwtPkg "auth_service/pkg/jwtpkg"
	"auth_service/pkg/logger"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
)

type containerRoutes struct {
	t usecase.ISsoService
	l logger.Logger
}

func newSsoRoutes(handler *echo.Group, t usecase.ISsoService, l logger.Logger) {
	r := &containerRoutes{t, l}

	// POST /api/v1/login
	handler.POST("/login", r.Auth)

	// POST /api/v1/register
	handler.POST("/register", r.Register)

	// PUT /api/v1/users
	handler.PUT("/users", r.UpdateUserInfo)

	// DELETE  /api/v1/users
	handler.DELETE("/users", r.DeleteAccount)

	// POST /api/v1/refresh
	handler.POST("/refresh", r.RefreshToken)
}

func (r *containerRoutes) Auth(c echo.Context) error {
	const op = "controller.Auth"

	ctx := c.Request().Context()

	u := new(entity.LoginRequest)
	if err := c.Bind(u); err != nil {
		errorResponse(c, http.StatusInternalServerError, "internal error")

		return fmt.Errorf("%s: %w", op, err)
	}

	if len(u.Password) == 0 || len([]rune(u.Email)) == 0 {
		errorResponse(c, http.StatusBadRequest, "bad request")

		return fmt.Errorf("%s: %w", op, errors.New("bad request"))
	}

	access_id, accessToken, refreshToken, err := r.t.Login(ctx, u.Email, u.Password)
	if err != nil {
		if errors.Is(err, usecase.ErrInvalidCredentials) {
			errorResponse(c, http.StatusUnauthorized, "invalid credentials")

			return fmt.Errorf("%s: %w", op, err)
		}

		if errors.Is(err, usecase.ErrNoUser) {
			errorResponse(c, http.StatusUnauthorized, "no user")

			return fmt.Errorf("%s: %w", op, err)
		}

		errorResponse(c, http.StatusInternalServerError, "internal error")

		return fmt.Errorf("%s: %w", op, err)
	}

	return c.JSON(http.StatusOK, entity.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		AccessId:     access_id,
	})
}

func (r *containerRoutes) Register(c echo.Context) error {
	const op = "controller.Register"

	ctx := c.Request().Context()

	u := new(entity.RegisterRequest)
	if err := c.Bind(u); err != nil {
		errorResponse(c, http.StatusInternalServerError, "internal error")

		return fmt.Errorf("%s: %w", op, err)
	}

	if len([]rune(u.Password)) < 10 {
		errorResponse(c, http.StatusBadRequest, "password must be equal or longer than 10")

		return fmt.Errorf("%s: %w", op, errors.New("password must be equal or longer than 10"))
	}

	if len([]rune(u.Email)) == 0 {
		errorResponse(c, http.StatusBadRequest, "email is required")

		return fmt.Errorf("%s: %w", op, errors.New("email is required"))
	}

	if len([]rune(u.Username)) == 0 {
		errorResponse(c, http.StatusBadRequest, "username is required")

		return fmt.Errorf("%s: %w", op, errors.New("username is required"))
	}

	register, err := r.t.Register(ctx, u.Email, u.Password, u.Username)
	if err != nil {
		if errors.Is(err, usecase.ErrUserExist) {
			errorResponse(c, http.StatusUnauthorized, "email or username is exist")

			return fmt.Errorf("%s: %w", op, err)
		}

		errorResponse(c, http.StatusInternalServerError, "internal error")

		return fmt.Errorf("%s: %w", op, err)
	}

	// TODO подумать про добавление автоматической авторизации
	return c.JSON(http.StatusOK, map[string]interface{}{
		"user_id": register,
	})
}

func (r *containerRoutes) UpdateUserInfo(c echo.Context) error {
	const op = "controller.UpdateUserInfo"

	// достаём access token
	token := jwtPkg.ExtractToken(c)
	if token == "" {
		errorResponse(c, http.StatusUnauthorized, "token is required")

		return fmt.Errorf("%s: %s", op, "token is required")
	}

	// валидируем токен и достаём id пользователя
	userId, err := jwtPkg.ValidateTokenAndGetUserId(token)
	if err != nil {
		errorResponse(c, http.StatusUnauthorized, "wrong token")

		return fmt.Errorf("%s: %s", op, err)
	}

	ctx := c.Request().Context()

	u := new(entity.UpdateUserRequest)
	if err = c.Bind(u); err != nil {
		errorResponse(c, http.StatusInternalServerError, "internal error")

		return fmt.Errorf("%s: %w", op, err)
	}

	if len(u.Password) < 10 || len(u.Email) == 0 {
		errorResponse(c, http.StatusBadRequest, "bad request")

		return fmt.Errorf("%s: %w", op, errors.New("bad request"))
	}

	_, err = r.t.UpdateUserInfo(ctx, userId, u.Email, u.Password, u.Name, u.Surname, u.Username)
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, "internal_old error")

		return fmt.Errorf("%s: %w", op, err)
	}

	return c.JSON(http.StatusOK, entity.UpdateUserResponse{
		User: entity.User{
			ID:       userId,
			Email:    u.Email,
			Name:     u.Name,
			Surname:  u.Surname,
			Username: u.Username,
		},
	})
}

func (r *containerRoutes) DeleteAccount(c echo.Context) error {
	const op = "controller.DeleteAccount"

	// достаём access token
	token := jwtPkg.ExtractToken(c)
	if token == "" {
		errorResponse(c, http.StatusBadRequest, "bad request")

		return fmt.Errorf("%s: %s", op, "token is required")
	}

	// валидируем токен и достаём id пользователя
	userId, err := jwtPkg.ValidateTokenAndGetUserId(token)
	if err != nil {
		errorResponse(c, http.StatusUnauthorized, "bad request")

		return fmt.Errorf("%s: %s", op, err)
	}

	ctx := c.Request().Context()

	isSucceed, err := r.t.DeleteAccount(ctx, userId)
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, "internal_old error")

		return fmt.Errorf("%s: %w", op, err)
	}

	return c.JSON(http.StatusOK, entity.DeleteUserResponse{
		IsSuccessfully: isSucceed,
	})
}

func (r *containerRoutes) RefreshToken(c echo.Context) error {
	const op = "controller.RefreshToken"

	refreshTokenOld := jwtPkg.ExtractToken(c)
	ctx := c.Request().Context()

	accessToken, refreshToken, err := r.t.RefreshToken(ctx, refreshTokenOld)
	if err != nil {
		errorResponse(c, http.StatusUnauthorized, "token error")

		return fmt.Errorf("%s: %w", op, err)
	}

	return c.JSON(http.StatusOK, entity.RefreshTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}
