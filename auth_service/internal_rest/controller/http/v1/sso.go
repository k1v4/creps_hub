package v1

import (
	"auth_service/internal_rest/usecase"
	"auth_service/pkg/logger"
	"github.com/labstack/echo/v4"
)

type conatainerRoutes struct {
	t usecase.ISsoService
	l logger.Logger
}

func newSsoRoutes(handler *echo.Group, t usecase.ISsoService, l logger.Logger) {
	r := &conatainerRoutes{t, l}
	_ = r
	// POST /api/auth
	//handler.POST("/auth", r.Auth)

	//GET /api/buy/{item}
	//handler.GET("/buy/:item", r.Buy)

	//POST /api/sendCoin"
	//handler.POST("/sendCoin", r.SendCoins)

	//GET  /api/info
	//handler.GET("/info", r.Info)
}
