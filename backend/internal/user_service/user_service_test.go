package service

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"testing"
	"time"

	"k8s-mon/internal/models"
	"k8s-mon/internal/repository"
	"k8s-mon/pkg/hash"
	"k8s-mon/pkg/jwt"

	"github.com/google/uuid"
)

type mockUserRepo struct {
	CreateUserFn           func(ctx context.Context, user *models.Employee) (*uuid.UUID, error)
	UserByEmailFn          func(ctx context.Context, email string) (*models.Employee, error)
	CreateJWTSessionFn     func(ctx context.Context, uuid, refreshToken uuid.UUID, fingerprint, ip string, expiresIn int64, RefreshTTL time.Duration) error
	RefreshSessionFn       func(ctx context.Context, oldRefresh uuid.UUID) (*models.RefreshSession, error)
	DeleteJWTSessionFn     func(ctx context.Context, refreshToken uuid.UUID) error
	GetUserFn              func(ctx context.Context, uuid uuid.UUID) (*models.Employee, error)
	UpdateUserProfileDocFn func(ctx context.Context, userID uuid.UUID, documentURL string) error
	GetAllUsersFn          func(ctx context.Context) ([]models.Employee, error)
	UpdateUserRoleFn       func(ctx context.Context, userID uuid.UUID, role string) error
	DeleteUserFn           func(ctx context.Context, id uuid.UUID) error
}

func (m *mockUserRepo) CreateUser(ctx context.Context, user *models.Employee) (*uuid.UUID, error) {
	if m.CreateUserFn != nil {
		return m.CreateUserFn(ctx, user)
	}
	return nil, nil
}

func (m *mockUserRepo) User(ctx context.Context, uuid uuid.UUID) (*models.Employee, error) {
	return nil, nil
}

func (m *mockUserRepo) UserByEmail(ctx context.Context, email string) (*models.Employee, error) {
	if m.UserByEmailFn != nil {
		return m.UserByEmailFn(ctx, email)
	}
	return nil, errors.New("not implemented")
}

func (m *mockUserRepo) CreateJWTSession(ctx context.Context, uuid, refreshToken uuid.UUID, fingerprint, ip string, expiresIn int64, RefreshTTL time.Duration) error {
	if m.CreateJWTSessionFn != nil {
		return m.CreateJWTSessionFn(ctx, uuid, refreshToken, fingerprint, ip, expiresIn, RefreshTTL)
	}
	return nil
}

func (m *mockUserRepo) DeleteJWTSession(ctx context.Context, refreshToken uuid.UUID) error {
	if m.DeleteJWTSessionFn != nil {
		return m.DeleteJWTSessionFn(ctx, refreshToken)
	}
	return nil
}

func (m *mockUserRepo) RefreshSession(ctx context.Context, oldRefresh uuid.UUID) (*models.RefreshSession, error) {
	if m.RefreshSessionFn != nil {
		return m.RefreshSessionFn(ctx, oldRefresh)
	}
	return nil, errors.New("not implemented")
}

func (m *mockUserRepo) GetUser(ctx context.Context, uuid uuid.UUID) (*models.Employee, error) {
	if m.GetUserFn != nil {
		return m.GetUserFn(ctx, uuid)
	}
	return nil, errors.New("not implemented")
}

func (m *mockUserRepo) GetAllUsers(ctx context.Context) ([]models.Employee, error) {
	if m.GetAllUsersFn != nil {
		return m.GetAllUsersFn(ctx)
	}
	return nil, nil
}

func (m *mockUserRepo) UpdateUserRole(ctx context.Context, userID uuid.UUID, role string) error {
	if m.UpdateUserRoleFn != nil {
		return m.UpdateUserRoleFn(ctx, userID, role)
	}
	return nil
}

func (m *mockUserRepo) DeleteUser(ctx context.Context, id uuid.UUID) error {
	if m.DeleteUserFn != nil {
		return m.DeleteUserFn(ctx, id)
	}
	return nil
}

func (m *mockUserRepo) UpdateUserProfileDocument(ctx context.Context, userID uuid.UUID, documentURL string) error {
	if m.UpdateUserProfileDocFn != nil {
		return m.UpdateUserProfileDocFn(ctx, userID, documentURL)
	}
	return nil
}

func newLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

func TestUserService_LoginSuccess(t *testing.T) {
	hashPassword, err := hash.HashPassword("secret123")
	if err != nil {
		t.Fatal(err)
	}

	called := false
	repo := &mockUserRepo{
		UserByEmailFn: func(ctx context.Context, email string) (*models.Employee, error) {
			return &models.Employee{
				UUID:           uuid.New(),
				Name:           "Test User",
				Email:          email,
				Role:           "user",
				HashedPassword: hashPassword,
			}, nil
		},
		CreateJWTSessionFn: func(ctx context.Context, uuid, refreshToken uuid.UUID, fingerprint, ip string, expiresIn int64, RefreshTTL time.Duration) error {
			called = true
			return nil
		},
	}

	jwtSvc := jwt.NewUserJWTpkg("secret-key", 15*time.Minute)
	service := NewUserService(repo, nil, newLogger(), jwtSvc, 24*time.Hour)

	accessToken, refreshToken, err := service.Login(context.Background(), "test@example.com", "secret123", "agent", "127.0.0.1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if accessToken == "" || refreshToken == "" {
		t.Fatal("expected non-empty tokens")
	}
	if !called {
		t.Fatal("expected CreateJWTSession to be called")
	}

	claims, err := jwtSvc.ParseToken(accessToken)
	if err != nil {
		t.Fatalf("failed to parse token: %v", err)
	}
	if claims.Email != "test@example.com" {
		t.Fatalf("expected email claim to match, got %q", claims.Email)
	}
}

func TestUserService_LoginPasswordMismatch(t *testing.T) {
	hashPassword, err := hash.HashPassword("secret123")
	if err != nil {
		t.Fatal(err)
	}

	repo := &mockUserRepo{
		UserByEmailFn: func(ctx context.Context, email string) (*models.Employee, error) {
			return &models.Employee{
				UUID:           uuid.New(),
				Email:          email,
				Role:           "user",
				HashedPassword: hashPassword,
			}, nil
		},
	}

	service := NewUserService(repo, nil, newLogger(), jwt.NewUserJWTpkg("secret-key", 15*time.Minute), 24*time.Hour)

	_, _, err = service.Login(context.Background(), "test@example.com", "wrong-password", "agent", "127.0.0.1")
	if !errors.Is(err, ErrPasswordMismatch) {
		t.Fatalf("expected ErrPasswordMismatch, got %v", err)
	}
}

func TestUserService_CreateUser_UserAlreadyExists(t *testing.T) {
	repo := &mockUserRepo{
		CreateUserFn: func(ctx context.Context, user *models.Employee) (*uuid.UUID, error) {
			return nil, fmt.Errorf("repo: %w", repository.ErrUserAlreadyExists)
		},
	}

	service := NewUserService(repo, nil, newLogger(), jwt.NewUserJWTpkg("secret-key", 15*time.Minute), 24*time.Hour)

	_, err := service.CreateUser(context.Background(), "Test", "test@example.com", "Engineer", "Dev", "user", "secret123")
	if !errors.Is(err, ErrUserAlreadyExists) {
		t.Fatalf("expected ErrUserAlreadyExists, got %v", err)
	}
}

func TestUserService_RefreshTokenSuccess(t *testing.T) {
	existingUser := models.Employee{UUID: uuid.New(), Name: "Test User", Email: "test@example.com", Role: "user"}
	jwtSvc := jwt.NewUserJWTpkg("secret-key", 15*time.Minute)
	accessToken, err := jwtSvc.GenerateAccessToken(existingUser)
	if err != nil {
		t.Fatal(err)
	}

	oldRefresh := uuid.New()
	created := false
	repo := &mockUserRepo{
		RefreshSessionFn: func(ctx context.Context, oldRefreshID uuid.UUID) (*models.RefreshSession, error) {
			return &models.RefreshSession{
				UUID:         existingUser.UUID.String(),
				RefreshToken: oldRefresh.String(),
				Fingerprint:  "127.0.0.1agent",
				IP:           "127.0.0.1",
				ExpiresIn:    int(time.Now().Add(time.Hour).Unix()),
			}, nil
		},
		CreateJWTSessionFn: func(ctx context.Context, uuid, refreshToken uuid.UUID, fingerprint, ip string, expiresIn int64, RefreshTTL time.Duration) error {
			created = true
			return nil
		},
	}

	service := NewUserService(repo, nil, newLogger(), jwtSvc, 24*time.Hour)

	newRefresh, newAccess, err := service.RefreshToken(context.Background(), oldRefresh, "agent", "127.0.0.1", accessToken)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if newRefresh == nil || newAccess == "" {
		t.Fatal("expected refreshed credentials")
	}
	if !created {
		t.Fatal("expected CreateJWTSession to be called")
	}
}

func TestUserService_RefreshTokenInvalidFingerprint(t *testing.T) {
	existingUser := models.Employee{UUID: uuid.New(), Name: "Test User", Email: "test@example.com", Role: "user"}
	jwtSvc := jwt.NewUserJWTpkg("secret-key", 15*time.Minute)
	accessToken, err := jwtSvc.GenerateAccessToken(existingUser)
	if err != nil {
		t.Fatal(err)
	}

	oldRefresh := uuid.New()
	repo := &mockUserRepo{
		RefreshSessionFn: func(ctx context.Context, oldRefreshID uuid.UUID) (*models.RefreshSession, error) {
			return &models.RefreshSession{
				UUID:         existingUser.UUID.String(),
				RefreshToken: oldRefresh.String(),
				Fingerprint:  "wrong-fingerprint",
				IP:           "127.0.0.1",
				ExpiresIn:    int(time.Now().Add(time.Hour).Unix()),
			}, nil
		},
	}

	service := NewUserService(repo, nil, newLogger(), jwtSvc, 24*time.Hour)

	_, _, err = service.RefreshToken(context.Background(), oldRefresh, "agent", "127.0.0.1", accessToken)
	if !errors.Is(err, ErrInvalidRefreshSession) {
		t.Fatalf("expected ErrInvalidRefreshSession, got %v", err)
	}
}
