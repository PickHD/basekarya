package recruitment

import (
	"basekarya-backend/internal/modules/onboarding"
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

type StorageProvider interface {
	UploadFileByte(ctx context.Context, objectName string, reader io.Reader, size int64, contentType string) (string, error)
}

type OnboardingProvider interface {
	CreateWorkflow(ctx context.Context, req *onboarding.CreateWorkflowRequest) error
}
