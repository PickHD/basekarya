package payroll

import (
	"hris-backend/internal/modules/company"
	"hris-backend/internal/modules/user"
)

type UserProvider interface {
	FindAllEmployeeActive() ([]user.Employee, error)
}

type AttendanceProvider interface {
	GetBulkLateDuration(month, year int) (map[uint]int, error)
}

type ReimbursementProvider interface {
	GetBulkApprovedAmount(month, year int) (map[uint]float64, error)
}

type CompanyProvider interface {
	FindByID(id uint) (*company.Company, error)
}

type NotificationProvider interface {
	SendNotification(userID uint,
		Type string,
		Title string,
		Message string, relatedID uint) error
}
