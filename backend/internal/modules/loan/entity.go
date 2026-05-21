package loan

import (
	"basekarya-backend/internal/modules/user"
	"basekarya-backend/pkg/constants"
	"database/sql"
	"time"
)

type Loan struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	UserID    uint      `gorm:"not null" json:"user_id"`
	User      user.User `gorm:"foreignKey:UserID" json:"user,omitempty"`
	CompanyID uint      `gorm:"index;not null" json:"company_id"`

	EmployeeID uint          `gorm:"not null" json:"employee_id"`
	Employee   user.Employee `gorm:"foreignKey:EmployeeID" json:"employee,omitempty"`

	ApprovedBy *uint      `json:"approved_by"`
	Approver   *user.User `gorm:"foreignKey:ApprovedBy" json:"approver,omitempty"`

	TotalAmount       float64 `gorm:"type:decimal(15,2);not null" json:"total_amount"`
	InstallmentAmount float64 `gorm:"type:decimal(15,2);not null" json:"installment_amount"`
	RemainingAmount   float64 `gorm:"type:decimal(15,2);not null" json:"remaining_amount"`
	Reason            string  `gorm:"type:text" json:"reason"`

	Status          constants.LoanStatus `gorm:"type:enum('PENDING','APPROVED','REJECTED','PAID_OFF');default:'PENDING'" json:"status"`
	RejectionReason sql.NullString       `gorm:"type:text" json:"rejection_reason"`
}

func (Loan) TableName() string {
	return "loans"
}
