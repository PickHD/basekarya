package infrastructure

import (
	"context"
	"fmt"
	"hris-backend/internal/config"
	"hris-backend/pkg/logger"
	"io"
	"mime/multipart"
	"strings"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioStorageProvider struct {
	client       *minio.Client
	bucketName   string
	isSecure     bool
	publicDomain string
}

func NewMinioStorage(cfg *config.Config) *MinioStorageProvider {
	minioClient, err := minio.New(cfg.Minio.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.Minio.AccessKey, cfg.Minio.SecretKey, ""),
		Secure: cfg.Minio.IsSecure,
	})
	if err != nil {
		logger.Errorw("Failed to connect to MinIO:", err)
	}

	logger.Info("Connected to MinIO Object Storage")

	return &MinioStorageProvider{
		client:       minioClient,
		bucketName:   cfg.Minio.BucketName,
		isSecure:     cfg.Minio.IsSecure,
		publicDomain: cfg.Minio.PublicDomain,
	}
}

func (m *MinioStorageProvider) UploadFileMultipart(ctx context.Context, file *multipart.FileHeader, objectName string) (string, error) {
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	// Upload the file
	info, err := m.client.PutObject(ctx, m.bucketName, objectName, src, file.Size, minio.PutObjectOptions{
		ContentType: file.Header.Get("Content-Type"),
	})
	if err != nil {
		return "", err
	}

	// Generate permanent URL
	protocol := "http"
	if m.isSecure {
		protocol = "https"
	}

	// Clean endpoint to avoid double slashes
	endpoint := strings.TrimSuffix(m.publicDomain, "/")
	url := fmt.Sprintf("%s://%s/%s/%s", protocol, endpoint, m.bucketName, info.Key)

	return url, nil
}

func (m *MinioStorageProvider) UploadFileByte(ctx context.Context, objectName string, reader io.Reader, size int64, contentType string) (string, error) {
	// Upload the file
	_, err := m.client.PutObject(ctx, m.bucketName, objectName, reader, size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload file: %w", err)
	}

	// Generate permanent URL
	protocol := "http"
	if m.isSecure {
		protocol = "https"
	}

	// Clean endpoint to avoid double slashes
	endpoint := strings.TrimSuffix(m.publicDomain, "/")
	url := fmt.Sprintf("%s://%s/%s/%s", protocol, endpoint, m.bucketName, objectName)

	return url, nil
}
