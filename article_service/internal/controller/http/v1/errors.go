package v1

import (
	"article_service/internal/entity"
	"errors"
	"github.com/labstack/echo/v4"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
)

func errorResponse(c echo.Context, code int, msg string) error {
	return c.JSON(code, entity.ErrorResponse{Error: msg})
}
