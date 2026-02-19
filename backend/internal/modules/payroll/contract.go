package payroll

import (
	"context"
	"hris-backend/internal/modules/company"
	"hris-backend/internal/modules/user"
)

type UserProvider interface {
	FindAllEmployeeActive(ctx context.Context) ([]user.Employee, error)
}

type AttendanceProvider interface {
	GetBulkLateDuration(ctx context.Context, month, year int) (map[uint]int, error)
}

type ReimbursementProvider interface {
	GetBulkApprovedAmount(ctx context.Context, month, year int) (map[uint]float64, error)
}

type CompanyProvider interface {
	FindByID(ctx context.Context, id uint) (*company.Company, error)
}

type NotificationProvider interface {
	SendNotification(userID uint,
		Type string,
		Title string,
		Message string, relatedID uint) error
}
