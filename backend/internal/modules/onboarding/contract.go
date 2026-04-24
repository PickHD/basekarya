package onboarding

import (
	"basekarya-backend/internal/modules/company"
	"basekarya-backend/internal/modules/master"
	"basekarya-backend/internal/modules/rbac"
	"basekarya-backend/internal/modules/user"
	"context"
	"io"
)

type NotificationProvider interface {
	SendNotification(userID uint, Type string, Title string, Message string, relatedID uint) error
	BlastNotification(userIDs []uint, Type string, Title string, Message string, relatedID uint) error
}

type UserProvider interface {
	FindApprovalUsers(ctx context.Context, permissionApprovalName string) ([]uint, error)
	CreateEmployee(ctx context.Context, req *user.CreateEmployeeRequest) error
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

type RoleProvider interface {
	FindRoleByName(ctx context.Context, name string) (*rbac.Role, error)
}

type MasterProvider interface {
	FindDepartmentByName(ctx context.Context, name string) (*master.Department, error)
	FindShiftByName(ctx context.Context, name string) (*master.Shift, error)
}
