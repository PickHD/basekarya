package reimbursement

import (
	"context"
	"mime/multipart"
)

type StorageProvider interface {
	UploadFileMultipart(ctx context.Context, file *multipart.FileHeader, objectName string) (string, error)
}

type NotificationProvider interface {
	SendNotification(userID uint,
		Type string,
		Title string,
		Message string, relatedID uint) error
}

type UserProvider interface {
	FindAdminID(ctx context.Context) (uint, error)
}
