package infrastructure

import (
	"bytes"
	"context"
	"fmt"
	"hris-backend/internal/config"
	"hris-backend/pkg/logger"
	"io"
	"mime/multipart"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioStorageProvider struct {
	client       *minio.Client
	bucketName   string
	isSecure     bool
	publicDomain string
}

func NewMinioStorage(cfg *config.MinioConfig) *MinioStorageProvider {
	minioClient, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: cfg.IsSecure,
	})
	if err != nil {
		logger.Errorw("Failed to connect to MinIO:", err)
	}

	logger.Info("Connected to MinIO Object Storage")

	return &MinioStorageProvider{
		client:       minioClient,
		bucketName:   cfg.BucketName,
		isSecure:     cfg.IsSecure,
		publicDomain: cfg.PublicDomain,
	}
}

func (m *MinioStorageProvider) UploadFileMultipart(ctx context.Context, file *multipart.FileHeader, objectName string) (string, error) {
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	// Before upload, resize image & turn down the quality image till 75%
	img, err := imaging.Decode(src)
	if err != nil {
		return "", err
	}

	dstImage := imaging.Resize(img, 800, 0, imaging.Lanczos)

	var buf bytes.Buffer
	err = imaging.Encode(&buf, dstImage, imaging.JPEG, imaging.JPEGQuality(75))
	if err != nil {
		return "", err
	}

	// Upload the file
	info, err := m.client.PutObject(ctx, m.bucketName, objectName, &buf, int64(buf.Len()), minio.PutObjectOptions{
		ContentType: file.Header.Get("Content-Type"),
	})
	if err != nil {
		return "", err
	}

	return m.generateURL(objectName, info.Key), nil
}

func (m *MinioStorageProvider) UploadFileByte(ctx context.Context, objectName string, reader io.Reader, size int64, contentType string) (string, error) {
	// Before upload, resize image & turn down the quality image till 75%
	img, err := imaging.Decode(reader)
	if err != nil {
		return "", err
	}

	dstImage := imaging.Resize(img, 800, 0, imaging.Lanczos)

	var buf bytes.Buffer
	err = imaging.Encode(&buf, dstImage, imaging.JPEG, imaging.JPEGQuality(75))
	if err != nil {
		return "", err
	}

	// Upload the file
	info, err := m.client.PutObject(ctx, m.bucketName, objectName, &buf, int64(buf.Len()), minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload file: %w", err)
	}

	return m.generateURL(objectName, info.Key), nil
}

func (m *MinioStorageProvider) generateURL(objectName, key string) string {
	protocol := "http"
	if m.isSecure {
		protocol = "https"
	}

	finalKey := key
	if finalKey == "" {
		finalKey = objectName
	}

	endpoint := strings.TrimSuffix(m.publicDomain, "/")
	return fmt.Sprintf("%s://%s/%s/%s", protocol, endpoint, m.bucketName, finalKey)
}
