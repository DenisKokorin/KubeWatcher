package reviewservice

import (
	"context"
	"fmt"
	"k8s-mon/internal/models"
	"k8s-mon/internal/repository"
	"log/slog"

	"github.com/google/uuid"
)

type ReviewService struct {
	repository repository.ReviewRepositoryInterface
	log        *slog.Logger
}

func NewReviewService(rep repository.ReviewRepositoryInterface, log *slog.Logger) *ReviewService {
	return &ReviewService{
		repository: rep,
		log:        log,
	}
}

func (r *ReviewService) GetEmployee(ctx context.Context, id uuid.UUID) (*models.Employee, error) {
	const op = "service.GetEmployee"
	r.log.With(slog.String("op", op))
	r.log.Info("Getting employee", slog.String("id", id.String()))

	emp, err := r.repository.GetEmployee(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return emp, nil
}

func (r *ReviewService) GetAllEmployees(ctx context.Context) ([]*models.Employee, error) {
	const op = "service.GetAllEmployees"
	r.log.With(slog.String("op", op))
	r.log.Info("Getting all employees")

	emp, err := r.repository.GetAllEmployees(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return emp, nil
}

func (r *ReviewService) GetEmployees(ctx context.Context, query models.EmployeeQuery) (*models.EmployeeList, error) {
	const op = "service.GetEmployees"
	r.log.With(slog.String("op", op))
	r.log.Info("Getting filtered employees", slog.Any("query", query))

	emp, total, err := r.repository.GetEmployees(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &models.EmployeeList{Employees: emp, Total: total, Page: query.Page, Size: query.Size}, nil
}

func (r *ReviewService) GetEmployeesByTeam(ctx context.Context, team string) ([]*models.Employee, error) {
	const op = "service.GetEmployeesByTeam"
	r.log.With(slog.String("op", op))
	r.log.Info("Getting employee", slog.String("team", team))

	emp, err := r.repository.GetEmployeesByTeam(ctx, team)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return emp, nil
}

func (r *ReviewService) CreateReview(ctx context.Context, req *models.ReviewRequest) (*uuid.UUID, error) {
	const op = "service.CreateReview"
	r.log.With(slog.String("op", op))
	r.log.Info("Creating review", slog.Any("req", req))

	revID, err := r.repository.CreateReview(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return revID, nil
}

func (r *ReviewService) GetReview(ctx context.Context, id uuid.UUID) (*models.Review, error) {
	const op = "service.GetReview"
	r.log.With(slog.String("op", op))
	r.log.Info("Getting review", slog.String("id", id.String()))

	rev, err := r.repository.GetReview(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return rev, nil
}

func (r *ReviewService) GetEmployeeReviews(ctx context.Context, employeeID uuid.UUID) ([]*models.Review, error) {
	const op = "service.GetEmployeeReviews"
	r.log.With(slog.String("op", op))
	r.log.Info("Getting employee reviews", slog.String("employeeID", employeeID.String()))

	emp, err := r.repository.GetEmployeeReviews(ctx, employeeID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return emp, nil
}

func (r *ReviewService) GetAllReviews(ctx context.Context) ([]*models.Review, error) {
	const op = "service.GetAllReviews"
	r.log.With(slog.String("op", op))
	r.log.Info("Getting all reviews")

	rev, err := r.repository.GetAllReviews(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return rev, nil
}

func (r *ReviewService) GetTeamStats(ctx context.Context, team string) (*models.TeamStats, error) {
	const op = "service.GetTeamStats"
	r.log.With(slog.String("op", op))
	r.log.Info("Getting team stats", slog.String("team", team))

	stat, err := r.repository.GetTeamStats(ctx, team)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return stat, nil
}
