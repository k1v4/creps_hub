package v1

import (
	"errors"
	"github.com/labstack/echo/v4"
	"release_service/internal/entity"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
)

func errorResponse(c echo.Context, code int, msg string) error {
	return c.JSON(code, entity.ErrorResponse{Error: msg})
}
