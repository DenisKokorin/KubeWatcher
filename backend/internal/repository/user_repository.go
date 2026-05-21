package repository

import (
	"context"
	"errors"
	"fmt"
	"k8s-mon/internal/models"
	"k8s-mon/pkg/db"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/redis/go-redis/v9"
)

var (
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrNoSessionFound    = errors.New("no session found")
)

type UserRepository struct {
	db      db.PostgresDB
	redisDB *redis.Client
}

type UserRepositoryInterface interface {
	CreateUser(ctx context.Context, user *models.Employee) (*uuid.UUID, error)
	User(ctx context.Context, uuid uuid.UUID) (*models.Employee, error)
	UserByEmail(ctx context.Context, email string) (*models.Employee, error)
	CreateJWTSession(ctx context.Context, uuid, refreshToken uuid.UUID, fingerprint, ip string, expiresIn int64, RefreshTTL time.Duration) error
	DeleteJWTSession(ctx context.Context, refreshToken uuid.UUID) error
	RefreshSession(ctx context.Context, oldRefresh uuid.UUID) (*models.RefreshSession, error)
	GetUser(ctx context.Context, uuid uuid.UUID) (*models.Employee, error)
	GetAllUsers(ctx context.Context) ([]models.Employee, error)
	UpdateUserRole(ctx context.Context, userID uuid.UUID, role string) error
	DeleteUser(ctx context.Context, id uuid.UUID) error
	UpdateUserProfileDocument(ctx context.Context, userID uuid.UUID, documentURL string) error
}

func NewUserRepository(db db.PostgresDB, redisDB *redis.Client) *UserRepository {
	return &UserRepository{db: db, redisDB: redisDB}
}

func (u *UserRepository) CreateUser(ctx context.Context, user *models.Employee) (*uuid.UUID, error) {
	const op = "user_repository.CreateUser"

	uuid := uuid.New()
	_, err := u.db.DB.Exec(ctx, "INSERT INTO users (uid, name, email, position, team, role, hire_date, profile_document_url, hashed_password) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9)", uuid.String(),
		user.Name, user.Email, user.Position, user.Team, user.Role, user.HireDate, (*string)(nil), user.HashedPassword)
	if err != nil {
		var pgxErr *pgconn.PgError
		if errors.As(err, &pgxErr) {
			if pgxErr.Code == "23505" {
				return nil, fmt.Errorf("%s: %w", op, ErrUserAlreadyExists)
			}
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &uuid, nil
}

func (u *UserRepository) User(ctx context.Context, uuid uuid.UUID) (*models.Employee, error) {
	const op = "user_repository.User"
	var user models.Employee
	var profileDocURL *string
	err := u.db.DB.QueryRow(ctx, "SELECT uid, name, email, position, team, role, hire_date, profile_document_url, hashed_password FROM users WHERE uid = $1", uuid.String()).Scan(&user.UUID, &user.Name,
		&user.Email, &user.Position, &user.Team, &user.Role, &user.HireDate, &profileDocURL, &user.HashedPassword)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	user.ProfileDocumentURL = profileDocURL
	return &user, nil
}

func (u *UserRepository) UserByEmail(ctx context.Context, email string) (*models.Employee, error) {
	const op = "repository.UserByEmail"
	var user models.Employee
	var profileDocURL *string
	err := u.db.DB.QueryRow(ctx, "SELECT uid, name, email, position, team, role, hire_date, profile_document_url, hashed_password FROM users WHERE email = $1", email).Scan(&user.UUID, &user.Name,
		&user.Email, &user.Position, &user.Team, &user.Role, &user.HireDate, &profileDocURL, &user.HashedPassword)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	user.ProfileDocumentURL = profileDocURL
	return &user, nil
}

func (u *UserRepository) CreateJWTSession(ctx context.Context, uuid, refreshToken uuid.UUID, fingerprint, ip string, expiresIn int64, RefreshTTL time.Duration) error {
	const op = "repository.CreateJWTSession"
	values := map[string]any{
		"uuid":          uuid.String(),
		"refresh_token": refreshToken.String(),
		"fingerprint":   fingerprint,
		"ip":            ip,
		"expires_in":    expiresIn,
	}
	err := u.redisDB.HSet(ctx, refreshToken.String(), values).Err()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	err = u.redisDB.Expire(ctx, refreshToken.String(), RefreshTTL).Err()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (u *UserRepository) DeleteJWTSession(ctx context.Context, refreshToken uuid.UUID) error {
	const op = "repository.DeleteJWTSession"
	err := u.redisDB.Del(ctx, refreshToken.String()).Err()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (u *UserRepository) RefreshSession(ctx context.Context, oldRefresh uuid.UUID) (*models.RefreshSession, error) {
	const op = "repository.RefreshSession"
	session, err := u.redisDB.HGetAll(ctx, oldRefresh.String()).Result()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	err = u.redisDB.Del(ctx, oldRefresh.String()).Err()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	exp, err := strconv.Atoi(session["expires_in"])
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	resSession := models.RefreshSession{
		UUID:         session["uuid"],
		RefreshToken: session["refresh_token"],
		Fingerprint:  session["fingerprint"],
		IP:           session["ip"],
		ExpiresIn:    exp,
	}
	return &resSession, nil
}

func (u *UserRepository) GetUser(ctx context.Context, id uuid.UUID) (*models.Employee, error) {
	return u.User(ctx, id)
}

func (u *UserRepository) GetAllUsers(ctx context.Context) ([]models.Employee, error) {
	const op = "repository.GetAllUsers"

	var users []models.Employee
	rows, err := u.db.DB.Query(ctx, "SELECT uid, name, email, position, team, role, hire_date, profile_document_url, hashed_password FROM users ORDER BY name ASC")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	for rows.Next() {
		var user models.Employee
		var profileDocURL *string
		err := rows.Scan(&user.UUID, &user.Name, &user.Email, &user.Position, &user.Team, &user.Role, &user.HireDate, &profileDocURL, &user.HashedPassword)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		user.ProfileDocumentURL = profileDocURL
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return users, nil
}

func (u *UserRepository) UpdateUserRole(ctx context.Context, userID uuid.UUID, role string) error {
	const op = "repository.UpdateUserRole"

	_, err := u.db.DB.Exec(ctx, "UPDATE users SET role = $1 WHERE uid = $2", role, userID.String())
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (u *UserRepository) DeleteUser(ctx context.Context, id uuid.UUID) error {
	const op = "repository.DeleteUser"

	_, err := u.db.DB.Exec(ctx, "DELETE FROM users WHERE uid = $1", id.String())
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (u *UserRepository) UpdateUserProfileDocument(ctx context.Context, userID uuid.UUID, documentURL string) error {
	const op = "repository.UpdateUserProfileDocument"

	var url *string
	if documentURL != "" {
		url = &documentURL
	}

	_, err := u.db.DB.Exec(ctx, "UPDATE users SET profile_document_url = $1 WHERE uid = $2", url, userID.String())
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
