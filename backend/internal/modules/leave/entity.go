package leave

import (
	"basekarya-backend/internal/modules/master"
	"basekarya-backend/internal/modules/user"
	"basekarya-backend/pkg/constants"
	"time"

	"gorm.io/gorm"
)

type LeaveBalance struct {
	ID          uint `gorm:"primaryKey" json:"id"`
	EmployeeID  uint `gorm:"index" json:"employee_id"`
	LeaveTypeID uint `json:"leave_type_id"`
	Year        int  `json:"year"`
	QuotaTotal  int  `json:"quota_total"`
	QuotaUsed   int  `json:"quota_used"`
	QuotaLeft   int  `json:"quota_left"`

	Employee  *user.Employee    `gorm:"foreignKey:EmployeeID" json:"employee,omitempty"`
	LeaveType *master.LeaveType `gorm:"foreignKey:LeaveTypeID" json:"leave_type,omitempty"`
}

type LeaveRequest struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	UserID      uint `json:"user_id"`
	EmployeeID  uint `json:"employee_id"`
	LeaveTypeID uint `json:"leave_type_id"`

	StartDate time.Time `gorm:"type:date;not null" json:"start_date"`
	EndDate   time.Time `gorm:"type:date;not null" json:"end_date"`
	TotalDays int       `json:"total_days"`

	Reason        string `gorm:"type:text" json:"reason"`
	AttachmentURL string `json:"attachment_url"`

	Status constants.LeaveStatus `gorm:"type:varchar(20);default:'PENDING'" json:"status"`

	ApprovedBy      *uint  `json:"approved_by"`
	RejectionReason string `json:"rejection_reason"`

	User      user.User         `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Employee  *user.Employee    `gorm:"foreignKey:EmployeeID" json:"employee,omitempty"`
	LeaveType *master.LeaveType `gorm:"foreignKey:LeaveTypeID" json:"leave_type,omitempty"`
}
