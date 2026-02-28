package overtime

import (
	"basekarya-backend/internal/modules/user"
	"basekarya-backend/pkg/constants"
	"database/sql"
	"time"
)

type Overtime struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	UserID uint      `gorm:"not null" json:"user_id"`
	User   user.User `gorm:"foreignKey:UserID" json:"user,omitempty"`

	EmployeeID uint          `gorm:"not null" json:"employee_id"`
	Employee   user.Employee `gorm:"foreignKey:EmployeeID" json:"employee,omitempty"`

	ApprovedBy *uint      `json:"approved_by"`
	Approver   *user.User `gorm:"foreignKey:ApprovedBy" json:"approver,omitempty"`

	Date            string `gorm:"type:date;not null" json:"date"`
	StartTime       string `gorm:"type:time;not null" json:"start_time"`
	EndTime         string `gorm:"type:time;not null" json:"end_time"`
	DurationMinutes int    `gorm:"type:int;not null" json:"duration_minutes"`
	Reason          string `gorm:"type:text" json:"reason"`

	Status          constants.OvertimeStatus `gorm:"type:enum('PENDING','APPROVED','REJECTED','PAID');default:'PENDING'" json:"status"`
	RejectionReason sql.NullString           `gorm:"type:text" json:"rejection_reason"`
}

func (Overtime) TableName() string {
	return "overtimes"
}
