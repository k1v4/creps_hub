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

	// GET /api/v1/articles/{id}
	handler.GET("/articles/:id", r.GetArticle)

	// POST /api/v1/articles
	handler.POST("/articles", r.PostArticle)

	// DELETE /api/v1/articles/{id}
	handler.DELETE("/articles/:id", r.DeleteArticle)

	// GET /api/v1/articles?limit=5&offset=0
	handler.GET("/articles", r.ListArticles)

	// GET /api/v1/user_articles
	handler.GET("/user_articles", r.GetArticlesByUser)
}

func (r *containerRoutes) ListArticles(c echo.Context) error {
	const op = "v1.ListArticles"

	ctx := c.Request().Context()

	limitStr := c.QueryParam("limit")
	offsetStr := c.QueryParam("offset")

	// преобразуем limit и offset в числа
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 10 // значение по умолчанию
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0 // значение по умолчанию
	}

	articles, err := r.t.FindAllArticle(ctx, limit, offset)
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, "internal server error")

		return fmt.Errorf("%s: %w", op, err)
	}

	return c.JSON(http.StatusOK, entity.PaginatedResponse{
		Items:  articles,
		Total:  limit + offset,
		Limit:  limit,
		Offset: offset,
	})
}

func (r *containerRoutes) DeleteArticle(c echo.Context) error {
	const op = "v1.DeleteArticle"

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

	articleId := c.Param("id")

	if len(strings.TrimSpace(articleId)) == 0 {
		errorResponse(c, http.StatusBadRequest, "bad request")

		return fmt.Errorf("%s: %s", op, "item name is required")
	}

	articleIdInt, err := strconv.Atoi(articleId)
	if err != nil {
		errorResponse(c, http.StatusBadRequest, "bad request")

		return fmt.Errorf("%s: %s", op, err)
	}

	deleteResp, err := r.t.DeleteArticle(ctx, articleIdInt, userId)
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, "internal server error")

		return fmt.Errorf("%s: %s", op, err)
	}

	return c.JSON(http.StatusOK, entity.DeleteArticleResponse{IsDeleted: deleteResp})
}

func (r *containerRoutes) PostArticle(c echo.Context) error {
	const op = "v1.PostArticle"

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

	articleId, err := r.t.AddArticle(ctx, userId, u.Title, u.Content, u.ImageName, u.ImageData)
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, "bad request")

		return fmt.Errorf("%s: %w", op, err)
	}

	return c.JSON(http.StatusOK, entity.PostArticleResponse{Id: articleId})
}

func (r *containerRoutes) GetArticle(c echo.Context) error {
	const op = "v1.GetArticle"

	ctx := c.Request().Context()

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

func (r *containerRoutes) GetArticlesByUser(c echo.Context) error {
	const op = "v1.GetArticlesByUser"

	ctx := c.Request().Context()

	token := jwtpkg.ExtractToken(c)
	if token == "" {
		errorResponse(c, http.StatusUnauthorized, "Unauthorized")

		return fmt.Errorf("%s: %s", op, "token is required")
	}

	// получаем user id
	userId, err := jwtpkg.ValidateTokenAndGetUserId(token)
	if err != nil {
		fmt.Println(err, err.Error())
		if err.Error() == "token expired" {
			errorResponse(c, http.StatusUnauthorized, "token expired")

			return fmt.Errorf("%s: %s", op, err)
		}

		errorResponse(c, http.StatusUnauthorized, "Unauthorized")

		return fmt.Errorf("%s: %s", op, err)
	}

	allArticleByUser, err := r.t.FindAllArticleByUser(ctx, userId)
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, "bad request")

		return fmt.Errorf("%s: %w", op, err)
	}

	return c.JSON(http.StatusOK, allArticleByUser)
}
