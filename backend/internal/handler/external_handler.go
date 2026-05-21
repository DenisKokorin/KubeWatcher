package handler

import (
	"context"
	"k8s-mon/internal/models"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type ExternalHandler struct {
	ExternalService ExternalService
}

type ExternalService interface {
	FetchPosts(ctx context.Context, userID, limit int, search string) ([]models.ExternalPost, error)
}

func NewExternalHandler(service ExternalService) *ExternalHandler {
	return &ExternalHandler{
		ExternalService: service,
	}
}

func (h *ExternalHandler) GetPosts(c echo.Context) error {
	query := c.QueryParam("search")
	limit := 0
	userID := 0

	if limitParam := c.QueryParam("limit"); limitParam != "" {
		parsed, err := strconv.Atoi(limitParam)
		if err == nil && parsed > 0 {
			limit = parsed
		}
	}
	if userParam := c.QueryParam("userId"); userParam != "" {
		parsed, err := strconv.Atoi(userParam)
		if err == nil && parsed > 0 {
			userID = parsed
		}
	}

	posts, err := h.ExternalService.FetchPosts(c.Request().Context(), userID, limit, query)
	if err != nil {
		return c.JSON(http.StatusBadGateway, models.Response{
			Success: false,
			Error:   "External service unavailable",
		})
	}

	return c.JSON(http.StatusOK, models.Response{
		Success: true,
		Data:    posts,
	})
}
