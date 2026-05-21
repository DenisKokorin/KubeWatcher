package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"k8s-mon/internal/middleware"
	"k8s-mon/internal/models"
	"k8s-mon/pkg/jwt"
	"k8s-mon/pkg/storage"
	"mime/multipart"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type mockUserService struct {
	LoginFn   func(ctx context.Context, email, password, userAgent, ip string) (string, string, error)
	GetUserFn func(ctx context.Context, id uuid.UUID) (*models.Employee, error)
}

func (m *mockUserService) CreateUser(ctx context.Context, name, email, position, team, role, password string) (*uuid.UUID, error) {
	panic("not implemented")
}
func (m *mockUserService) Login(ctx context.Context, email, password, userAgent, ip string) (string, string, error) {
	return m.LoginFn(ctx, email, password, userAgent, ip)
}
func (m *mockUserService) RefreshToken(ctx context.Context, oldRefresh uuid.UUID, userAgent, ip string, oldAccessToken string) (*uuid.UUID, string, error) {
	panic("not implemented")
}
func (m *mockUserService) Logout(ctx context.Context, refreshToken string) error {
	panic("not implemented")
}
func (m *mockUserService) GetUser(ctx context.Context, id uuid.UUID) (*models.Employee, error) {
	return m.GetUserFn(ctx, id)
}
func (m *mockUserService) GetAllUsers(ctx context.Context) ([]models.Employee, error) {
	panic("not implemented")
}
func (m *mockUserService) UpdateUserRole(ctx context.Context, userID uuid.UUID, role string) error {
	panic("not implemented")
}
func (m *mockUserService) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	panic("not implemented")
}
func (m *mockUserService) UploadUserProfileDocument(ctx context.Context, userID uuid.UUID, file multipart.File, header *multipart.FileHeader) (*storage.FileInfo, error) {
	panic("not implemented")
}
func (m *mockUserService) GetUserDocumentURL(ctx context.Context, userID uuid.UUID) (string, error) {
	panic("not implemented")
}
func (m *mockUserService) DeleteUserDocument(ctx context.Context, userID uuid.UUID) error {
	panic("not implemented")
}

func TestUserHandler_LoginEndpoint(t *testing.T) {
	e := echo.New()
	service := &mockUserService{
		LoginFn: func(ctx context.Context, email, password, userAgent, ip string) (string, string, error) {
			return "access-token", "refresh-token", nil
		},
	}
	h := NewUserHandler(service)

	req := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(`{"email":"test@example.com","password":"secret"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if err := h.Login(c); err != nil {
		t.Fatal(err)
	}

	if rec.Code != http.StatusAccepted {
		t.Fatalf("expected status 202, got %d", rec.Code)
	}

	var body map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}

	if body["AccessToken"] != "access-token" {
		t.Fatalf("unexpected access token %q", body["AccessToken"])
	}
}

func TestUserHandler_GetMeEndpointWithAuth(t *testing.T) {
	e := echo.New()
	jwtSvc := jwt.NewUserJWTpkg("secret-key", time.Hour)
	middleware := middleware.NewJWTMiddleware(jwtSvc)

	user := models.Employee{
		UUID:  uuid.New(),
		Name:  "Test User",
		Email: "test@example.com",
		Role:  "user",
	}
	service := &mockUserService{
		GetUserFn: func(ctx context.Context, id uuid.UUID) (*models.Employee, error) {
			return &user, nil
		},
	}
	h := NewUserHandler(service)

	accessToken, err := jwtSvc.GenerateAccessToken(user)
	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest(http.MethodGet, "/auth/me", nil)
	req.Header.Set(echo.HeaderAuthorization, "Bearer "+accessToken)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	hFunc := middleware.Auth(h.GetMe)
	if err := hFunc(c); err != nil {
		t.Fatal(err)
	}

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	var body map[string]interface{}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}

	if body["email"] != "test@example.com" {
		t.Fatalf("expected email to match, got %v", body["email"])
	}
}
