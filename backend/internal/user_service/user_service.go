package service

import (
	"context"
	"errors"
	"fmt"
	"k8s-mon/internal/models"
	"k8s-mon/internal/repository"
	"k8s-mon/pkg/hash"
	"k8s-mon/pkg/jwt"
	"k8s-mon/pkg/storage"
	"log/slog"
	"mime/multipart"
	"time"

	"github.com/google/uuid"
)

var (
	ErrUserAlreadyExists     = errors.New("user already exists")
	ErrPasswordMismatch      = errors.New("password mismatch")
	ErrTokenExpired          = errors.New("token expired")
	ErrInvalidRefreshSession = errors.New("invalid refresh session")
)

type UserService struct {
	userRepo   repository.UserRepositoryInterface
	storageSvc storage.StorageService
	log        *slog.Logger
	jwtTkn     jwt.TokenService
	RefreshTTL time.Duration
}

type UserServiceInterface interface {
	CreateUser(ctx context.Context, name, email, position, team, role, password string) (*uuid.UUID, error)
	Login(ctx context.Context, email, password, userAgent, ip string) (string, string, error)
	RefreshToken(ctx context.Context, oldRefresh uuid.UUID, userAgent, ip string, oldAccessToken string) (*uuid.UUID, string, error)
	Logout(ctx context.Context, refreshToken string) error
	GetUser(ctx context.Context, id uuid.UUID) (*models.Employee, error)
	GetAllUsers(ctx context.Context) ([]models.Employee, error)
	UpdateUserRole(ctx context.Context, userID uuid.UUID, role string) error
	DeleteUser(ctx context.Context, userID uuid.UUID) error
	UploadUserProfileDocument(ctx context.Context, userID uuid.UUID, file multipart.File, header *multipart.FileHeader) (*storage.FileInfo, error)
	GetUserDocumentURL(ctx context.Context, userID uuid.UUID) (string, error)
	DeleteUserDocument(ctx context.Context, userID uuid.UUID) error
}

func NewUserService(userRepo repository.UserRepositoryInterface, storageSvc storage.StorageService, log *slog.Logger, jwtTkn jwt.TokenService, RefreshTTL time.Duration) *UserService {
	return &UserService{userRepo: userRepo, storageSvc: storageSvc, log: log, jwtTkn: jwtTkn, RefreshTTL: RefreshTTL}
}

func (u *UserService) CreateUser(ctx context.Context, name, email, position, team, role, password string) (*uuid.UUID, error) {
	const op = "service.CreateUser"
	u.log.With(slog.String("op", op))
	u.log.Info("Creating user", slog.String("email", email))

	if role == "" {
		role = "user"
	}

	hashPassword, err := hash.HashPassword(password)
	if err != nil {
		slog.Error("failed to generate hashpassword", slog.String("error", err.Error()))
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	user := models.Employee{
		Name:           name,
		Email:          email,
		Position:       position,
		Team:           team,
		Role:           role,
		HashedPassword: hashPassword,
		HireDate:       time.Now(),
	}

	uid, err := u.userRepo.CreateUser(ctx, &user)
	if err != nil {
		if errors.Is(err, repository.ErrUserAlreadyExists) {
			u.log.Error("user already exists error", slog.String("email", user.Email))
			return nil, fmt.Errorf("%s: %w", op, ErrUserAlreadyExists)
		}
		u.log.Error("creating user error", slog.String("email", user.Email), slog.String("error", err.Error()))
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	u.log.Info("created user", slog.String("email", user.Email))
	return uid, nil
}

func (u *UserService) Login(ctx context.Context, email, password, userAgent, ip string) (string, string, error) {
	const op = "service.Login"
	u.log.With(slog.String("op", op))
	u.log.Info("logging the user", slog.String("email", email))

	user, err := u.userRepo.UserByEmail(ctx, email)
	if err != nil {
		u.log.Error("failed to get user", slog.String("error", err.Error()))
		return "", "", fmt.Errorf("%s: %w", op, err)
	}
	if !hash.CheckPassword(user.HashedPassword, password) {
		u.log.Warn("passwords mismatch")
		return "", "", fmt.Errorf("%s: %w", op, ErrPasswordMismatch)
	}
	AccessToken, err := u.jwtTkn.GenerateAccessToken(*user)
	if err != nil {
		u.log.Error("failed to create access token", slog.String("error", err.Error()))
		return "", "", fmt.Errorf("%s: %w", op, err)
	}
	RefreshToken := uuid.New()

	fingerprint := ip + userAgent

	err = u.userRepo.CreateJWTSession(ctx, user.UUID, RefreshToken, fingerprint, ip, time.Now().Add(u.RefreshTTL).Unix(), u.RefreshTTL)
	if err != nil {
		u.log.Error("failed to create jwt session", slog.String("error", err.Error()))
		return "", "", fmt.Errorf("%s: %w", op, err)
	}
	return AccessToken, RefreshToken.String(), nil
}

func (u *UserService) Logout(ctx context.Context, refreshToken string) error {
	const op = "service.Logout"
	u.log.With(slog.String("op", op))
	u.log.Info("logging out the user")

	err := u.userRepo.DeleteJWTSession(ctx, uuid.MustParse(refreshToken))
	if err != nil {
		u.log.Error("failed to delete jwt session", slog.String("error", err.Error()))
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (u *UserService) RefreshToken(ctx context.Context, oldRefresh uuid.UUID, userAgent, ip string, oldAccessToken string) (*uuid.UUID, string, error) {
	const op = "service.RefreshToken"
	log := u.log.With(slog.String("op", op))
	log.Info("refreshing token")
	session, err := u.userRepo.RefreshSession(ctx, oldRefresh)
	if err != nil {
		log.Error("getting refresh token error", slog.String("error", err.Error()))
		return nil, "", fmt.Errorf("%s: %w", op, err)
	}
	if int64(session.ExpiresIn) < time.Now().Unix() {
		return nil, "", ErrTokenExpired
	}
	fingerprint := ip + userAgent

	if session.Fingerprint != fingerprint {
		return nil, "", ErrInvalidRefreshSession
	}
	newRefresh := uuid.New()
	err = u.userRepo.CreateJWTSession(ctx, uuid.MustParse(session.UUID), newRefresh, session.Fingerprint,
		session.IP, time.Now().Add(u.RefreshTTL).Unix(), u.RefreshTTL)
	if err != nil {
		log.Error("creating jwt session error", slog.String("error", err.Error()))
		return nil, "", fmt.Errorf("%s: %w", op, err)
	}

	token, err := u.jwtTkn.RegenerateToken(oldAccessToken)
	if err != nil {
		log.Error("creating jwt token error", slog.String("error", err.Error()))
		return nil, "", fmt.Errorf("%s: %w", op, err)
	}
	return &newRefresh, token, nil
}

func (u *UserService) GetUser(ctx context.Context, id uuid.UUID) (*models.Employee, error) {
	const op = "service.GetUser"
	log := u.log.With(slog.String("op", op))
	log.Info("getting user", slog.String("uuid", id.String()))

	user, err := u.userRepo.GetUser(ctx, id)
	if err != nil {
		log.Error("failed to get user", slog.String("error", err.Error()))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (u *UserService) GetAllUsers(ctx context.Context) ([]models.Employee, error) {
	const op = "service.GetAllUsers"
	log := u.log.With(slog.String("op", op))
	log.Info("getting all users")

	users, err := u.userRepo.GetAllUsers(ctx)
	if err != nil {
		log.Error("failed to get all users", slog.String("error", err.Error()))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return users, nil
}

func (u *UserService) UpdateUserRole(ctx context.Context, userID uuid.UUID, role string) error {
	const op = "service.UpdateUserRole"
	log := u.log.With(slog.String("op", op))
	log.Info("updating user role", slog.String("uuid", userID.String()), slog.String("role", role))

	err := u.userRepo.UpdateUserRole(ctx, userID, role)
	if err != nil {
		log.Error("failed to update user role", slog.String("error", err.Error()))
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (u *UserService) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	const op = "service.DeleteUser"
	log := u.log.With(slog.String("op", op))
	log.Info("deleting user", slog.String("uuid", userID.String()))

	err := u.userRepo.DeleteUser(ctx, userID)
	if err != nil {
		log.Error("failed to delete user", slog.String("error", err.Error()))
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (u *UserService) UploadUserProfileDocument(ctx context.Context, userID uuid.UUID, file multipart.File, header *multipart.FileHeader) (*storage.FileInfo, error) {
	const op = "service.UploadUserProfileDocument"
	log := u.log.With(slog.String("op", op))
	log.Info("uploading user document", slog.String("uuid", userID.String()), slog.String("filename", header.Filename))

	// Upload file to storage
	fileInfo, err := u.storageSvc.UploadFile(ctx, file, header, "users", userID)
	if err != nil {
		log.Error("failed to upload file to storage", slog.String("error", err.Error()))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// Update user record with object key
	err = u.userRepo.UpdateUserProfileDocument(ctx, userID, fileInfo.ObjectKey)
	if err != nil {
		log.Error("failed to update user document metadata", slog.String("error", err.Error()))
		// Try to clean up uploaded file if metadata update fails
		if cleanupErr := u.storageSvc.DeleteFile(ctx, fileInfo.ObjectKey); cleanupErr != nil {
			log.Error("failed to cleanup uploaded file after metadata update failure", slog.String("error", cleanupErr.Error()))
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("user document uploaded successfully", slog.String("objectKey", fileInfo.ObjectKey))
	return fileInfo, nil
}

func (u *UserService) GetUserDocumentURL(ctx context.Context, userID uuid.UUID) (string, error) {
	const op = "service.GetUserDocumentURL"
	log := u.log.With(slog.String("op", op))
	log.Info("getting user document URL", slog.String("uuid", userID.String()))

	// Get user to find document object key
	user, err := u.userRepo.GetUser(ctx, userID)
	if err != nil {
		log.Error("failed to get user", slog.String("error", err.Error()))
		return "", fmt.Errorf("%s: %w", op, err)
	}

	if user.ProfileDocumentURL == nil || *user.ProfileDocumentURL == "" {
		return "", fmt.Errorf("%s: no document found for user", op)
	}

	// Generate presigned URL
	presignedURL, err := u.storageSvc.GetPresignedURL(ctx, *user.ProfileDocumentURL)
	if err != nil {
		log.Error("failed to generate presigned URL", slog.String("error", err.Error()))
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return presignedURL, nil
}

func (u *UserService) DeleteUserDocument(ctx context.Context, userID uuid.UUID) error {
	const op = "service.DeleteUserDocument"
	log := u.log.With(slog.String("op", op))
	log.Info("deleting user document", slog.String("uuid", userID.String()))

	// Get user to find document object key
	user, err := u.userRepo.GetUser(ctx, userID)
	if err != nil {
		log.Error("failed to get user", slog.String("error", err.Error()))
		return fmt.Errorf("%s: %w", op, err)
	}

	if user.ProfileDocumentURL == nil || *user.ProfileDocumentURL == "" {
		return fmt.Errorf("%s: no document found for user", op)
	}

	// Delete file from storage
	err = u.storageSvc.DeleteFile(ctx, *user.ProfileDocumentURL)
	if err != nil {
		log.Error("failed to delete file from storage", slog.String("error", err.Error()))
		return fmt.Errorf("%s: %w", op, err)
	}

	// Clear document URL from user record
	err = u.userRepo.UpdateUserProfileDocument(ctx, userID, "")
	if err != nil {
		log.Error("failed to clear document metadata", slog.String("error", err.Error()))
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("user document deleted successfully")
	return nil
}
