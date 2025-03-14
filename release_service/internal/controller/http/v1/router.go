package v1

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
	"release_service/internal/usecase"
	"release_service/pkg/logger"
)

func NewRouter(handler *echo.Echo, l logger.Logger, t usecase.IReleaseUseCase) {
	// Middleware
	handler.Use(middleware.Logger())
	handler.Use(middleware.Recover())

	handler.GET("/api/releases/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})

	h := handler.Group("/api/v1")
	{
		newReleasesRoutes(h, t, l)
	}
}
