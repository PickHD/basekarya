package overtime

import (
	"basekarya-backend/pkg/constants"
	"time"
)

type OvertimeFilter struct {
	UserID uint
	Status string
	Page   int
	Limit  int
}

type OvertimeRequest struct {
	UserID     uint   `json:"-"`
	EmployeeID uint   `json:"-"`
	Date       string `json:"date" validate:"required"`
	StartTime  string `json:"start_time" validate:"required"`
	EndTime    string `json:"end_time" validate:"required"`
	Reason     string `json:"reason" validate:"required"`
}

type ActionRequest struct {
	ID              uint   `json:"-"`
	SuperAdminID    uint   `json:"-"`
	Action          string `json:"action" validate:"required"`
	RejectionReason string `json:"rejection_reason" validate:"omitempty"`
}

type OvertimeListResponse struct {
	ID              uint                     `json:"id"`
	EmployeeID      uint                     `json:"employee_id"`
	EmployeeName    string                   `json:"employee_name"`
	EmployeeNIK     string                   `json:"employee_nik"`
	Date            string                   `json:"date"`
	StartTime       string                   `json:"start_time"`
	EndTime         string                   `json:"end_time"`
	DurationMinutes int                      `json:"duration_minutes"`
	Status          constants.OvertimeStatus `json:"status"`
	CreatedAt       time.Time                `json:"created_at"`
}

type OvertimeDetailResponse struct {
	ID              uint                     `json:"id"`
	EmployeeID      uint                     `json:"employee_id"`
	EmployeeName    string                   `json:"employee_name"`
	EmployeeNIK     string                   `json:"employee_nik"`
	Date            string                   `json:"date"`
	StartTime       string                   `json:"start_time"`
	EndTime         string                   `json:"end_time"`
	DurationMinutes int                      `json:"duration_minutes"`
	Reason          string                   `json:"reason"`
	Status          constants.OvertimeStatus `json:"status"`
	RejectionReason string                   `json:"rejection_reason"`
	CreatedAt       time.Time                `json:"created_at"`
}
