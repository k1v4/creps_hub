package v1

import (
	"auth_service/internal_rest/entity"
	"errors"
	"github.com/labstack/echo/v4"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
)

func errorResponse(c echo.Context, code int, msg string) error {
	return c.JSON(code, entity.ErrorResponse{Error: msg})
}
