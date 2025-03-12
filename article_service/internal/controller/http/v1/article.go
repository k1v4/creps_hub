package v1

import (
	"article_service/internal/entity"
	"article_service/internal/usecase"
	"article_service/pkg/jwtpkg"
	"article_service/pkg/logger"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
	"strings"
)

type containerRoutes struct {
	t usecase.IArticleService
	l logger.Logger
}

func newArticleRoutes(handler *echo.Group, t usecase.IArticleService, l logger.Logger) {
	r := &containerRoutes{t, l}

	// GET /api/v1/login
	handler.GET("/article/:id", r.ListArticle)

	// POST /api/buy/{item}
	handler.POST("/article", r.PostArticle)

	//// POST /api/sendCoin
	//handler.PUT("/users", r.UpdateUserInfo)
	//
	//// GET  /api/info
	//handler.DELETE("/users", r.DeleteAccount)
	//
	//// POST /api/refresh
	//handler.POST("/refresh", r.RefreshToken)
}

func (r *containerRoutes) PostArticle(c echo.Context) error {
	const op = "ArticleUseCase.PostArticle"

	ctx := c.Request().Context()

	token := jwtpkg.ExtractToken(c)
	if token == "" {
		errorResponse(c, http.StatusUnauthorized, "Unauthorized")

		return fmt.Errorf("%s: %s", op, "token is required")
	}

	// получаем user id
	userId, err := jwtpkg.ValidateTokenAndGetUserId(token)
	if err != nil {
		errorResponse(c, http.StatusUnauthorized, "Unauthorized")

		return fmt.Errorf("%s: %s", op, err)
	}

	u := new(entity.PostArticleRequest)
	if err = c.Bind(u); err != nil {
		errorResponse(c, http.StatusBadRequest, "bad request")

		return fmt.Errorf("%s: %w", op, err)
	}

	if len([]rune(u.Title)) == 0 || len([]rune(u.Title)) > 100 {
		errorResponse(c, http.StatusBadRequest, "bad request")

		return fmt.Errorf("%s: %w", op, errors.New("wrong len if title"))
	}

	articleId, err := r.t.AddArticle(ctx, userId, u.Title, u.Content)
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, "bad request")

		return fmt.Errorf("%s: %w", op, err)
	}

	return c.JSON(http.StatusOK, entity.PostArticleResponse{Id: articleId})
}

func (r *containerRoutes) ListArticle(c echo.Context) error {
	const op = "controller.ListArticle"

	ctx := c.Request().Context()

	token := jwtpkg.ExtractToken(c)
	if token == "" {
		errorResponse(c, http.StatusUnauthorized, "bad request")

		return fmt.Errorf("%s: %s", op, "token is required")
	}

	// получаем user id
	_, err := jwtpkg.ValidateTokenAndGetUserId(token)
	if err != nil {
		errorResponse(c, http.StatusUnauthorized, "bad request")

		return fmt.Errorf("%s: %s", op, err)
	}

	articleId := c.Param("id")

	if len(strings.TrimSpace(articleId)) == 0 {
		errorResponse(c, http.StatusBadRequest, "bad request")

		return fmt.Errorf("%s: %s", op, "item name is required")
	}

	id, err := strconv.Atoi(articleId)
	if err != nil {
		errorResponse(c, http.StatusBadRequest, "bad request")

		return fmt.Errorf("%s: %s", op, err)
	}

	findArticle, err := r.t.FindArticle(ctx, id)
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, "internal server error")

		return fmt.Errorf("%s: %s", op, err)
	}

	return c.JSON(http.StatusOK, findArticle)
}
