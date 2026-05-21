package storage

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var (
	ErrInvalidFileType = errors.New("invalid file type")
	ErrFileTooLarge    = errors.New("file too large")
	ErrFileNotFound    = errors.New("file not found")
)

const (
	MaxFileSize     = 10 * 1024 * 1024 // 10MB
	PresignedExpiry = 15 * time.Minute
)

var AllowedFileTypes = map[string]bool{
	".pdf":  true,
	".doc":  true,
	".docx": true,
	".txt":  true,
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".gif":  true,
}

type FileInfo struct {
	FileName  string
	FileSize  int64
	FileType  string
	ObjectKey string
	URL       string
}

type StorageService interface {
	UploadFile(ctx context.Context, file multipart.File, header *multipart.FileHeader, entityType string, entityID uuid.UUID) (*FileInfo, error)
	GetPresignedURL(ctx context.Context, objectKey string) (string, error)
	DeleteFile(ctx context.Context, objectKey string) error
	ValidateFile(file multipart.File, header *multipart.FileHeader) error
}

type MinIOStorage struct {
	client     *minio.Client
	bucketName string
	logger     *slog.Logger
}

type StorageConfig struct {
	Endpoint        string
	AccessKeyID     string
	SecretAccessKey string
	BucketName      string
	UseSSL          bool
}

func NewMinIOStorage(config StorageConfig, logger *slog.Logger) (*MinIOStorage, error) {
	client, err := minio.New(config.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.AccessKeyID, config.SecretAccessKey, ""),
		Secure: config.UseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create minio client: %w", err)
	}

	// Create bucket if it doesn't exist
	ctx := context.Background()
	exists, err := client.BucketExists(ctx, config.BucketName)
	if err != nil {
		return nil, fmt.Errorf("failed to check bucket existence: %w", err)
	}

	if !exists {
		err = client.MakeBucket(ctx, config.BucketName, minio.MakeBucketOptions{})
		if err != nil {
			return nil, fmt.Errorf("failed to create bucket: %w", err)
		}
		logger.Info("Created bucket", slog.String("bucket", config.BucketName))
	}

	return &MinIOStorage{
		client:     client,
		bucketName: config.BucketName,
		logger:     logger,
	}, nil
}

func (s *MinIOStorage) ValidateFile(file multipart.File, header *multipart.FileHeader) error {
	// Check file size
	if header.Size > MaxFileSize {
		return ErrFileTooLarge
	}

	// Check file type
	ext := strings.ToLower(filepath.Ext(header.Filename))
	if !AllowedFileTypes[ext] {
		return ErrInvalidFileType
	}

	return nil
}

func (s *MinIOStorage) UploadFile(ctx context.Context, file multipart.File, header *multipart.FileHeader, entityType string, entityID uuid.UUID) (*FileInfo, error) {
	const op = "storage.UploadFile"

	// Validate file
	if err := s.ValidateFile(file, header); err != nil {
		s.logger.Error("file validation failed", slog.String("error", err.Error()))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// Generate unique object key
	ext := filepath.Ext(header.Filename)
	objectKey := fmt.Sprintf("%s/%s/%s%s", entityType, entityID.String(), uuid.New().String(), ext)

	// Upload file
	_, err := s.client.PutObject(ctx, s.bucketName, objectKey, file, header.Size, minio.PutObjectOptions{
		ContentType: header.Header.Get("Content-Type"),
	})
	if err != nil {
		s.logger.Error("failed to upload file", slog.String("objectKey", objectKey), slog.String("error", err.Error()))
		return nil, fmt.Errorf("%s: failed to upload file: %w", op, err)
	}

	s.logger.Info("file uploaded successfully", slog.String("objectKey", objectKey))

	fileInfo := &FileInfo{
		FileName:  header.Filename,
		FileSize:  header.Size,
		FileType:  header.Header.Get("Content-Type"),
		ObjectKey: objectKey,
	}

	return fileInfo, nil
}

func (s *MinIOStorage) GetPresignedURL(ctx context.Context, objectKey string) (string, error) {
	const op = "storage.GetPresignedURL"

	// Check if object exists
	_, err := s.client.StatObject(ctx, s.bucketName, objectKey, minio.StatObjectOptions{})
	if err != nil {
		if minio.ToErrorResponse(err).Code == "NoSuchKey" {
			return "", fmt.Errorf("%s: %w", op, ErrFileNotFound)
		}
		s.logger.Error("failed to stat object", slog.String("objectKey", objectKey), slog.String("error", err.Error()))
		return "", fmt.Errorf("%s: failed to check object: %w", op, err)
	}

	// Generate presigned URL
	presignedURL, err := s.client.PresignedGetObject(ctx, s.bucketName, objectKey, PresignedExpiry, nil)
	if err != nil {
		s.logger.Error("failed to generate presigned URL", slog.String("objectKey", objectKey), slog.String("error", err.Error()))
		return "", fmt.Errorf("%s: failed to generate presigned URL: %w", op, err)
	}

	return presignedURL.String(), nil
}

func (s *MinIOStorage) DeleteFile(ctx context.Context, objectKey string) error {
	const op = "storage.DeleteFile"

	err := s.client.RemoveObject(ctx, s.bucketName, objectKey, minio.RemoveObjectOptions{})
	if err != nil {
		s.logger.Error("failed to delete file", slog.String("objectKey", objectKey), slog.String("error", err.Error()))
		return fmt.Errorf("%s: failed to delete file: %w", op, err)
	}

	s.logger.Info("file deleted successfully", slog.String("objectKey", objectKey))
	return nil
}
