package reimbursement

import (
	"database/sql"
	"hris-backend/internal/modules/user"
	"hris-backend/pkg/constants"
	"time"

	"gorm.io/gorm"
)

type Reimbursement struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	UserID uint      `gorm:"not null" json:"user_id"`
	User   user.User `gorm:"foreignKey:UserID" json:"user,omitempty"`

	ApprovedBy *uint      `json:"approved_by"`
	Approver   *user.User `gorm:"foreignKey:ApprovedBy" json:"approver,omitempty"`

	Title           string                        `gorm:"type:varchar(255);not null" json:"title"`
	Description     string                        `gorm:"type:text" json:"description"`
	Amount          float64                       `gorm:"type:decimal(15,2);not null" json:"amount"`
	DateOfExpense   time.Time                     `gorm:"type:date;not null" json:"date_of_expense"`
	ProofFileURL    string                        `gorm:"type:varchar(255);not null" json:"proof_file_url"`
	Status          constants.ReimbursementStatus `gorm:"type:enum('PENDING','APPROVED','REJECTED');default:'PENDING'" json:"status"`
	RejectionReason sql.NullString                `gorm:"type:text" json:"rejection_reason"`
}
