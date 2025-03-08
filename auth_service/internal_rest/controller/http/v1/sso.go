package v1

import (
	"auth_service/internal_rest/entity"
	"auth_service/internal_rest/usecase"
	"auth_service/pkg/logger"
	"errors"
	"fmt"
	"github.com/k1v4/avito_shop/pkg/jwtPkg"
	"github.com/labstack/echo/v4"
	"net/http"
)

type containerRoutes struct {
	t usecase.ISsoService
	l logger.Logger
}

func newSsoRoutes(handler *echo.Group, t usecase.ISsoService, l logger.Logger) {
	r := &containerRoutes{t, l}

	// POST /api/auth
	handler.POST("/auth", r.Auth)

	//GET /api/buy/{item}
	handler.POST("/register", r.Register)

	//POST /api/sendCoin"
	//handler.POST("/sendCoin", r.SendCoins)

	//GET  /api/info
	//handler.GET("/info", r.Info)
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

	accessToken, refreshToken, err := r.t.Login(ctx, u.Email, u.Password)
	if err != nil {
		if errors.Is(err, usecase.ErrInvalidCredentials) {
			errorResponse(c, http.StatusUnauthorized, "bad request")

			return fmt.Errorf("%s: %w", op, err)
		}

		errorResponse(c, http.StatusInternalServerError, "internal error")

		return fmt.Errorf("%s: %w", op, err)
	}

	return c.JSON(http.StatusOK, entity.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}

func (r *containerRoutes) Register(c echo.Context) error {
	const op = "controller.Register"

	ctx := c.Request().Context()

	return c.JSON(http.StatusOK, "")
}

func (r *containerRoutes) UpdateUserInfo(c echo.Context) error {
	const op = "controller.UpdateUserInfo"

	token := jwtPkg.ExtractToken(c)
	if token == "" {
		errorResponse(c, http.StatusBadRequest, "bad request")

		return fmt.Errorf("%s: %s", op, "token is required")
	}
}
