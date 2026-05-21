package repository

import (
	"context"
	"errors"
	"fmt"
	"k8s-mon/internal/models"
	"k8s-mon/pkg/db"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type ReviewRepository struct {
	db db.PostgresDB
}

type ReviewRepositoryInterface interface {
	GetEmployee(ctx context.Context, id uuid.UUID) (*models.Employee, error)
	GetAllEmployees(ctx context.Context) ([]*models.Employee, error)
	GetEmployees(ctx context.Context, query models.EmployeeQuery) ([]*models.Employee, int, error)
	GetEmployeesByTeam(ctx context.Context, team string) ([]*models.Employee, error)
	CreateReview(ctx context.Context, req *models.ReviewRequest) (*uuid.UUID, error)
	GetReview(ctx context.Context, id uuid.UUID) (*models.Review, error)
	GetEmployeeReviews(ctx context.Context, employeeID uuid.UUID) ([]*models.Review, error)
	GetAllReviews(ctx context.Context) ([]*models.Review, error)
	GetTeamStats(ctx context.Context, team string) (*models.TeamStats, error)
}

func NewReviewRepository(db db.PostgresDB) *ReviewRepository {
	return &ReviewRepository{db: db}
}

func (r *ReviewRepository) GetEmployee(ctx context.Context, id uuid.UUID) (*models.Employee, error) {
	const op = "review_repository.GetEmployee"

	var employee models.Employee
	err := r.db.DB.QueryRow(
		ctx,
		`SELECT uid, name, email, position, team, role, hire_date, profile_document_url, hashed_password 
		 FROM users WHERE uid = $1`,
		id.String(),
	).Scan(
		&employee.UUID, &employee.Name, &employee.Email,
		&employee.Position, &employee.Team, &employee.Role, &employee.HireDate,
		&employee.ProfileDocumentURL, &employee.HashedPassword,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: employee not found", op)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &employee, nil
}

func (r *ReviewRepository) GetAllEmployees(ctx context.Context) ([]*models.Employee, error) {
	const op = "review_repository.GetAllEmployees"

	rows, err := r.db.DB.Query(
		ctx,
		`SELECT uid, name, email, position, team, role, hire_date, profile_document_url, hashed_password 
		 FROM users ORDER BY name ASC`,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var employees []*models.Employee
	for rows.Next() {
		var employee models.Employee
		err := rows.Scan(
			&employee.UUID, &employee.Name, &employee.Email,
			&employee.Position, &employee.Team, &employee.Role, &employee.HireDate,
			&employee.ProfileDocumentURL, &employee.HashedPassword,
		)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		employees = append(employees, &employee)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return employees, nil
}

func (r *ReviewRepository) GetEmployees(ctx context.Context, query models.EmployeeQuery) ([]*models.Employee, int, error) {
	const op = "review_repository.GetEmployees"

	if query.Size <= 0 {
		query.Size = 10
	}
	if query.Page < 0 {
		query.Page = 0
	}

	where := "WHERE TRUE"
	args := []interface{}{}
	argPos := 1

	if query.Search != "" {
		where += fmt.Sprintf(" AND (name ILIKE $%d OR email ILIKE $%d OR position ILIKE $%d OR team ILIKE $%d)", argPos, argPos, argPos, argPos)
		args = append(args, "%"+query.Search+"%")
		argPos++
	}
	if query.Team != "" {
		where += fmt.Sprintf(" AND team = $%d", argPos)
		args = append(args, query.Team)
		argPos++
	}
	if query.Position != "" {
		where += fmt.Sprintf(" AND position = $%d", argPos)
		args = append(args, query.Position)
		argPos++
	}
	if query.Role != "" {
		where += fmt.Sprintf(" AND role = $%d", argPos)
		args = append(args, query.Role)
		argPos++
	}

	sortBy := "name"
	sortOrder := "ASC"
	switch query.SortBy {
	case "email":
		sortBy = "email"
	case "position":
		sortBy = "position"
	case "team":
		sortBy = "team"
	case "hire_date":
		sortBy = "hire_date"
	}
	if query.SortOrder == "desc" {
		sortOrder = "DESC"
	}

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM users %s", where)
	var total int
	if err := r.db.DB.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("%s: %w", op, err)
	}

	limitOffsetQuery := fmt.Sprintf(
		"SELECT uid, name, email, position, team, role, hire_date, profile_document_url, hashed_password FROM users %s ORDER BY %s %s LIMIT $%d OFFSET $%d",
		where,
		sortBy,
		sortOrder,
		argPos,
		argPos+1,
	)
	args = append(args, query.Size, query.Page*query.Size)

	rows, err := r.db.DB.Query(ctx, limitOffsetQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var employees []*models.Employee
	for rows.Next() {
		var employee models.Employee
		err := rows.Scan(
			&employee.UUID, &employee.Name, &employee.Email,
			&employee.Position, &employee.Team, &employee.Role, &employee.HireDate,
			&employee.ProfileDocumentURL, &employee.HashedPassword,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("%s: %w", op, err)
		}
		employees = append(employees, &employee)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("%s: %w", op, err)
	}

	return employees, total, nil
}

func (r *ReviewRepository) GetEmployeesByTeam(ctx context.Context, team string) ([]*models.Employee, error) {
	const op = "review_repository.GetEmployeesByTeam"

	rows, err := r.db.DB.Query(
		ctx,
		`SELECT uid, name, email, position, team, role, hire_date, profile_document_url, hashed_password 
		 FROM users WHERE team = $1 ORDER BY name ASC`,
		team,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var employees []*models.Employee
	for rows.Next() {
		var employee models.Employee
		err := rows.Scan(
			&employee.UUID, &employee.Name, &employee.Email,
			&employee.Position, &employee.Team, &employee.Role, &employee.HireDate,
			&employee.ProfileDocumentURL, &employee.HashedPassword,
		)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		employees = append(employees, &employee)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return employees, nil
}

func (r *ReviewRepository) CreateReview(ctx context.Context, req *models.ReviewRequest) (*uuid.UUID, error) {
	const op = "review_repository.CreateReview"

	reviewID := uuid.New()
	now := time.Now()

	_, err := r.db.DB.Exec(
		ctx,
		`INSERT INTO reviews (id, employee_id, reviewer_id, period, rating, comments, goals, strengths, areas_for_improvement, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`,
		reviewID.String(),
		fmt.Sprintf("%s", req.EmployeeID),
		fmt.Sprintf("%s", req.ReviewerID),
		req.Period,
		req.Rating,
		req.Comments,
		req.Goals,
		req.Strengths,
		req.AreasForImprovement,
		now,
		now,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &reviewID, nil
}

func (r *ReviewRepository) GetReview(ctx context.Context, id uuid.UUID) (*models.Review, error) {
	const op = "review_repository.GetReview"

	var review models.Review
	err := r.db.DB.QueryRow(
		ctx,
		`SELECT id, employee_id, reviewer_id, period, rating, comments, goals, strengths, areas_for_improvement, created_at, updated_at
		 FROM reviews WHERE id = $1`,
		id.String(),
	).Scan(
		&review.ID, &review.EmployeeID, &review.ReviewerID,
		&review.Period, &review.Rating, &review.Comments,
		&review.Goals, &review.Strengths, &review.AreasForImprovement,
		&review.CreatedAt, &review.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: review not found", op)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &review, nil
}

func (r *ReviewRepository) GetEmployeeReviews(ctx context.Context, employeeID uuid.UUID) ([]*models.Review, error) {
	const op = "review_repository.GetEmployeeReviews"

	rows, err := r.db.DB.Query(
		ctx,
		`SELECT id, employee_id, reviewer_id, period, rating, comments, goals, strengths, areas_for_improvement, created_at, updated_at
		 FROM reviews WHERE employee_id = $1 ORDER BY created_at DESC`,
		employeeID.String(),
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var reviews []*models.Review
	for rows.Next() {
		var review models.Review
		err := rows.Scan(
			&review.ID, &review.EmployeeID, &review.ReviewerID,
			&review.Period, &review.Rating, &review.Comments,
			&review.Goals, &review.Strengths, &review.AreasForImprovement,
			&review.CreatedAt, &review.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		reviews = append(reviews, &review)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return reviews, nil
}

func (r *ReviewRepository) GetAllReviews(ctx context.Context) ([]*models.Review, error) {
	const op = "review_repository.GetAllReviews"

	rows, err := r.db.DB.Query(
		ctx,
		`SELECT id, employee_id, reviewer_id, period, rating, comments, goals, strengths, areas_for_improvement, created_at, updated_at
		 FROM reviews ORDER BY created_at DESC`,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var reviews []*models.Review
	for rows.Next() {
		var review models.Review
		err := rows.Scan(
			&review.ID, &review.EmployeeID, &review.ReviewerID,
			&review.Period, &review.Rating, &review.Comments,
			&review.Goals, &review.Strengths, &review.AreasForImprovement,
			&review.CreatedAt, &review.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		reviews = append(reviews, &review)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return reviews, nil
}

func (r *ReviewRepository) GetTeamStats(ctx context.Context, team string) (*models.TeamStats, error) {
	const op = "review_repository.GetTeamStats"

	var stats models.TeamStats
	err := r.db.DB.QueryRow(
		ctx,
		`SELECT 
			$1,
			COUNT(DISTINCT u.uid) as employee_count,
			COALESCE(AVG(rv.rating), 0) as average_rating,
			COUNT(rv.id) as reviews_count
		FROM users u
		LEFT JOIN reviews rv ON u.uid::text = rv.employee_id::text
		WHERE u.team = $1
		GROUP BY u.team`,
		team,
	).Scan(
		&stats.Team,
		&stats.EmployeeCount,
		&stats.AverageRating,
		&stats.ReviewsCount,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &models.TeamStats{
				Team:          team,
				EmployeeCount: 0,
				AverageRating: 0,
				ReviewsCount:  0,
			}, nil
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &stats, nil
}
