package leave

import (
	"context"
	"io"
)

type StorageProvider interface {
	UploadFileByte(ctx context.Context, objectName string, reader io.Reader, size int64, contentType string) (string, error)
}

type NotificationProvider interface {
	SendNotification(userID uint,
		Type string,
		Title string,
		Message string, relatedID uint) error
}

type UserProvider interface {
	FindAdminID() (uint, error)
}
