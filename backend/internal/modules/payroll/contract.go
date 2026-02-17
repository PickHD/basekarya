package payroll

import "hris-backend/internal/modules/user"

type UserProvider interface {
	FindAllEmployeeActive() ([]user.Employee, error)
}

type AttendanceProvider interface {
	GetBulkLateDuration(month, year int) (map[uint]int, error)
}

type ReimbursementProvider interface {
	GetBulkApprovedAmount(month, year int) (map[uint]float64, error)
}
