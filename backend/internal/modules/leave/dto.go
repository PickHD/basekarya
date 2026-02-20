package leave

import (
	"basekarya-backend/internal/modules/master"
	"basekarya-backend/pkg/constants"
	"time"
)

type ApplyRequest struct {
	UserID           uint   `json:"-"`
	EmployeeID       uint   `json:"-"`
	LeaveTypeID      uint   `json:"leave_type_id" validate:"required"`
	StartDate        string `json:"start_date" validate:"required"`
	EndDate          string `json:"end_date" validate:"required"`
	Reason           string `json:"reason" validate:"required"`
	AttachmentBase64 string `json:"attachment_base64"`
}

type LeaveActionRequest struct {
	RequestID       uint   `json:"-"`
	ApproverID      uint   `json:"-"`
	Action          string `json:"action" validate:"required"`
	RejectionReason string `json:"rejection_reason" validate:"omitempty"`
}

type LeaveFilter struct {
	Page   int    `json:"page"`
	Limit  int    `json:"limit"`
	Status string `json:"status"`
	Search string `json:"search"`
	UserID uint   `json:"-"`
}

type LeaveRequestListResponse struct {
	ID           uint                            `json:"id"`
	EmployeeID   uint                            `json:"employee_id"`
	EmployeeName string                          `json:"employee_name"`
	EmployeeNIK  string                          `json:"employee_nik"`
	LeaveTypeID  uint                            `json:"leave_type_id"`
	LeaveType    *master.LookupLeaveTypeResponse `json:"leave_type"`
	TotalDays    int                             `json:"total_days"`
	StartDate    time.Time                       `json:"start_date"`
	EndDate      time.Time                       `json:"end_date"`
	Status       constants.LeaveStatus           `json:"status"`
	CreatedAt    time.Time                       `json:"created_at"`
}

type LeaveRequestDetailResponse struct {
	ID              uint                            `json:"id"`
	EmployeeID      uint                            `json:"employee_id"`
	EmployeeName    string                          `json:"employee_name"`
	EmployeeNIK     string                          `json:"employee_nik"`
	LeaveTypeID     uint                            `json:"leave_type_id"`
	LeaveType       *master.LookupLeaveTypeResponse `json:"leave_type"`
	StartDate       time.Time                       `json:"start_date"`
	EndDate         time.Time                       `json:"end_date"`
	TotalDays       int                             `json:"total_days"`
	Reason          string                          `json:"reason"`
	AttachmentURL   string                          `json:"attachment_url"`
	Status          constants.LeaveStatus           `json:"status"`
	RejectionReason string                          `json:"rejection_reason"`
	CreatedAt       time.Time                       `json:"created_at"`
}
