package v1

import (
	"article_service/internal/usecase"
	"article_service/pkg/logger"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
)

func NewRouter(handler *echo.Echo, l logger.Logger, t usecase.IArticleService) {
	// Middleware
	handler.Use(middleware.Logger())
	handler.Use(middleware.Recover())

	handler.GET("/api/article/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})

	h := handler.Group("/api/v1")
	{
		newArticleRoutes(h, t, l)
	}
}
