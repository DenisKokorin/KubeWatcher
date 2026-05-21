package handler

import (
	"context"
	"k8s-mon/internal/models"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type ReviewHandler struct {
	ReviewService ReviewService
}

type ReviewService interface {
	GetEmployee(ctx context.Context, id uuid.UUID) (*models.Employee, error)
	GetAllEmployees(ctx context.Context) ([]*models.Employee, error)
	GetEmployees(ctx context.Context, query models.EmployeeQuery) (*models.EmployeeList, error)
	GetEmployeesByTeam(ctx context.Context, team string) ([]*models.Employee, error)
	CreateReview(ctx context.Context, req *models.ReviewRequest) (*uuid.UUID, error)
	GetReview(ctx context.Context, id uuid.UUID) (*models.Review, error)
	GetEmployeeReviews(ctx context.Context, employeeID uuid.UUID) ([]*models.Review, error)
	GetAllReviews(ctx context.Context) ([]*models.Review, error)
	GetTeamStats(ctx context.Context, team string) (*models.TeamStats, error)
}

func NewReviewHandler(r ReviewService) *ReviewHandler {
	return &ReviewHandler{
		ReviewService: r,
	}
}

func (r *ReviewHandler) GetEmployee(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{
			"error": "invalid id format",
		})
	}

	employee, err := r.ReviewService.GetEmployee(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]any{
			"error": "internal error",
		})
	}

	if employee == nil {
		return c.JSON(http.StatusNotFound, models.Response{
			Success: false,
			Error:   "Employee not found",
		})
	}

	return c.JSON(http.StatusOK, models.Response{
		Success: true,
		Data:    employee,
	})
}

func (r *ReviewHandler) GetAllEmployees(c echo.Context) error {
	query := models.EmployeeQuery{
		Search:    c.QueryParam("search"),
		Team:      c.QueryParam("team"),
		Position:  c.QueryParam("position"),
		Role:      c.QueryParam("role"),
		SortBy:    c.QueryParam("sortBy"),
		SortOrder: c.QueryParam("sortOrder"),
	}

	page, err := strconv.Atoi(c.QueryParam("page"))
	if err == nil {
		query.Page = page
	}
	size, err := strconv.Atoi(c.QueryParam("size"))
	if err == nil {
		query.Size = size
	}

	employeeList, err := r.ReviewService.GetEmployees(c.Request().Context(), query)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]any{
			"error": "internal error",
		})
	}

	return c.JSON(http.StatusOK, models.Response{
		Success: true,
		Data:    employeeList,
	})
}

func (r *ReviewHandler) CreateReview(c echo.Context) error {
	var req models.ReviewRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{
			"error": "bad request",
		})
	}

	// Validate that employee exists (using a placeholder UUID - this should be updated based on your actual employee ID format)
	// For now, we just create the review
	reviewID, err := r.ReviewService.CreateReview(c.Request().Context(), &req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]any{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, models.Response{
		Success: true,
		Message: "Review created successfully",
		Data:    reviewID,
	})
}

func (r *ReviewHandler) GetEmployeeWithReviews(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{
			"error": "invalid id format",
		})
	}

	ctx := c.Request().Context()
	employee, err := r.ReviewService.GetEmployee(ctx, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]any{
			"error": err.Error(),
		})
	}

	if employee == nil {
		return c.JSON(http.StatusNotFound, models.Response{
			Success: false,
			Error:   "Employee not found",
		})
	}

	reviews, err := r.ReviewService.GetEmployeeReviews(ctx, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]any{
			"error": err.Error(),
		})
	}

	employeeWithReviews := models.EmployeeWithReviews{
		Employee: employee,
		Reviews:  reviews,
	}

	return c.JSON(http.StatusOK, models.Response{
		Success: true,
		Data:    employeeWithReviews,
	})
}

func (r *ReviewHandler) GetTeamStats(c echo.Context) error {
	team := c.Param("team")

	stats, err := r.ReviewService.GetTeamStats(c.Request().Context(), team)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]any{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, models.Response{
		Success: true,
		Data:    stats,
	})
}

func (r *ReviewHandler) GetEmployeesByTeam(c echo.Context) error {
	team := c.Param("team")

	employees, err := r.ReviewService.GetEmployeesByTeam(c.Request().Context(), team)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]any{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, models.Response{
		Success: true,
		Data:    employees,
	})
}
