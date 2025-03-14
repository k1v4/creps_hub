package v1

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"release_service/internal/entity"
	"release_service/internal/usecase"
	"release_service/pkg/logger"
	"strconv"
	"strings"
)

type releasesRoutes struct {
	t usecase.IReleaseUseCase
	l logger.Logger
}

func newReleasesRoutes(handler *echo.Group, t usecase.IReleaseUseCase, l logger.Logger) {
	r := &releasesRoutes{t, l}

	// GET /api/v1/releases/{id}
	handler.GET("/releases/:id", r.getRelease)

	// DELETE /api/v1/releases/{id}
	handler.DELETE("/releases/:id", r.deleteRelease)

	// PUT /api/v1/releases/{id}
	handler.PUT("/releases/:id", r.updateRelease)
}

func (r *releasesRoutes) getRelease(c echo.Context) error {
	const op = "v1.getRelease"

	ctx := c.Request().Context()

	releaseID := c.Param("id")

	if len(strings.TrimSpace(releaseID)) == 0 {
		errorResponse(c, http.StatusBadRequest, "bad request")

		return fmt.Errorf("%s: %s", op, "releaseID is required")
	}

	releaseIdInt, err := strconv.Atoi(releaseID)
	if err != nil {
		errorResponse(c, http.StatusBadRequest, "bad request")

		return fmt.Errorf("%s: %s", op, err)
	}

	release, err := r.t.GetRelease(ctx, releaseIdInt)
	if err != nil {
		errorResponse(c, http.StatusNotFound, "release not found")

		return fmt.Errorf("%s: %s", op, err)
	}

	return c.JSON(http.StatusOK, release)
}

func (r *releasesRoutes) deleteRelease(c echo.Context) error {
	const op = "v1.deleteRelease"

	ctx := c.Request().Context()

	releaseID := c.Param("id")

	if len(strings.TrimSpace(releaseID)) == 0 {
		errorResponse(c, http.StatusBadRequest, "bad request")

		return fmt.Errorf("%s: %s", op, "releaseID is required")
	}

	releaseIdInt, err := strconv.Atoi(releaseID)
	if err != nil {
		errorResponse(c, http.StatusBadRequest, "bad request")

		return fmt.Errorf("%s: %s", op, err)
	}

	isDeleted, err := r.t.DeleteRelease(ctx, releaseIdInt)
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, "internal error")

		return fmt.Errorf("%s: %s", op, err)
	}

	return c.JSON(http.StatusOK, entity.DeleteRequest{
		IsDeleted: isDeleted,
	})
}

func (r *releasesRoutes) updateRelease(c echo.Context) error {
	const op = "v1.updateRelease"

	ctx := c.Request().Context()
}
