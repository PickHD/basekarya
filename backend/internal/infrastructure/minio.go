package infrastructure

import (
	"bytes"
	"context"
	"fmt"
	"basekarya-backend/internal/config"
	"basekarya-backend/pkg/logger"
	"io"
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var allowedMimeTypes = map[string]bool{
	"image/jpeg":      true,
	"image/jpg":       true,
	"image/png":       true,
	"image/webp":      true,
	"application/pdf": true,
}

var allowedExtensions = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".webp": true,
	".pdf":  true,
}

func validateFileType(fileName string, contentType string) error {
	ext := strings.ToLower(filepath.Ext(fileName))
	if !allowedExtensions[ext] {
		return fmt.Errorf("file type not allowed: %s", ext)
	}

	detectedType, err := detectContentType(fileName, contentType)
	if err != nil {
		return err
	}

	if !allowedMimeTypes[detectedType] {
		return fmt.Errorf("file content type not allowed: %s", detectedType)
	}

	return nil
}

func detectContentType(fileName string, contentType string) (string, error) {
	if contentType != "" {
		cleaned := strings.ToLower(strings.Split(contentType, ";")[0])
		cleaned = strings.TrimSpace(cleaned)
		if cleaned != "" {
			return cleaned, nil
		}
	}

	ext := strings.ToLower(filepath.Ext(fileName))
	switch ext {
	case ".jpg", ".jpeg":
		return "image/jpeg", nil
	case ".png":
		return "image/png", nil
	case ".webp":
		return "image/webp", nil
	case ".pdf":
		return "application/pdf", nil
	default:
		return "", fmt.Errorf("unknown file extension: %s", ext)
	}
}

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
	if err := validateFileType(file.Filename, file.Header.Get("Content-Type")); err != nil {
		return "", err
	}

	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	uploadReader := io.Reader(src)
	uploadSize := file.Size
	contentType := file.Header.Get("Content-Type")

	// Before upload, resize image & turn down the quality image till 75%
	if strings.HasPrefix(contentType, "image/") {
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

		uploadReader = &buf
		uploadSize = int64(buf.Len())
	}

	// Upload the file
	info, err := m.client.PutObject(ctx, m.bucketName, objectName, uploadReader, uploadSize, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", err
	}

	return m.generateURL(objectName, info.Key), nil
}

func (m *MinioStorageProvider) UploadFileByte(ctx context.Context, objectName string, reader io.Reader, size int64, contentType string) (string, error) {
	uploadReader := io.Reader(reader)
	uploadSize := size

	// Before upload, resize image & turn down the quality image till 75%
	if strings.HasPrefix(contentType, "image/") {
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

		uploadReader = &buf
		uploadSize = int64(buf.Len())
	}

	// Upload the file
	info, err := m.client.PutObject(ctx, m.bucketName, objectName, uploadReader, uploadSize, minio.PutObjectOptions{
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
