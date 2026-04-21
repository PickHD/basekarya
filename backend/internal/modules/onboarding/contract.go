package onboarding

import (
	"basekarya-backend/internal/modules/company"
	"context"
	"io"
)

type NotificationProvider interface {
	SendNotification(userID uint, Type string, Title string, Message string, relatedID uint) error
	BlastNotification(userIDs []uint, Type string, Title string, Message string, relatedID uint) error
}

type UserProvider interface {
	FindApprovalUsers(ctx context.Context, permissionApprovalName string) ([]uint, error)
}

type EmailProvider interface {
	Send(to string, subject string, htmlBody string) error
}

type StorageProvider interface {
	UploadFileByte(ctx context.Context, objectName string, reader io.Reader, size int64, contentType string) (string, error)
}

type CompanyProvider interface {
	FindByID(ctx context.Context, id uint) (*company.Company, error)
}
