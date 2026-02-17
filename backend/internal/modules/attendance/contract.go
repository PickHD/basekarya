package attendance

import (
	"context"
	"hris-backend/internal/modules/user"
	"io"
)

type StorageProvider interface {
	UploadFileByte(ctx context.Context, objectName string, reader io.Reader, size int64, contentType string) (string, error)
}

type LocationFetcher interface {
	GetAddressFromCoords(lat, long float64) string
}

type UserProvider interface {
	FindByID(id uint) (*user.User, error)
	CountActiveEmployee() (int64, error)
}
