package v1

import (
	"auth_service/internal_rest/usecase"
	"auth_service/pkg/logger"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func NewRouter(handler *echo.Echo, l logger.Logger, t usecase.ISsoService) {
	// Middleware
	handler.Use(middleware.Logger())
	handler.Use(middleware.Recover())

	h := handler.Group("/api")
	{
		newSsoRoutes(h, t, l)
	}
}
