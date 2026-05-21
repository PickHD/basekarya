package finance

import (
	"basekarya-backend/internal/modules/user"
	"basekarya-backend/pkg/constants"
	"database/sql"
	"time"
)

type FinanceCategory struct {
	ID          uint                  `gorm:"primaryKey" json:"id"`
	CreatedAt   time.Time             `json:"created_at"`
	UpdatedAt   time.Time             `json:"updated_at"`
	Name        string                `gorm:"type:varchar(255);not null" json:"name"`
	Type        constants.FinanceType `gorm:"type:enum('INCOME','EXPENSE');not null" json:"type"`
	Description sql.NullString        `gorm:"type:text" json:"description"`
}

func (FinanceCategory) TableName() string {
	return "finance_categories"
}

type FinanceTransaction struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	FinanceCategoryID uint            `gorm:"not null" json:"finance_category_id"`
	FinanceCategory   FinanceCategory `gorm:"foreignKey:FinanceCategoryID" json:"finance_category,omitempty"`

	CreatedBy  uint      `gorm:"not null" json:"created_by"`
	Creator    user.User `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`

	ApprovedBy *uint      `json:"approved_by"`
	Approver   *user.User `gorm:"foreignKey:ApprovedBy" json:"approver,omitempty"`

	Type            constants.FinanceType `gorm:"type:enum('INCOME','EXPENSE');not null" json:"type"`
	Amount          float64               `gorm:"type:decimal(15,2);not null" json:"amount"`
	Description     sql.NullString        `gorm:"type:text" json:"description"`
	TransactionDate time.Time             `gorm:"type:date;not null" json:"transaction_date"`
	ReferenceNumber sql.NullString        `gorm:"type:varchar(100)" json:"reference_number"`

	Status          constants.FinanceStatus `gorm:"type:enum('PENDING','APPROVED','REJECTED');default:'PENDING'" json:"status"`
	RejectionReason sql.NullString          `gorm:"type:text" json:"rejection_reason"`
}

func (FinanceTransaction) TableName() string {
	return "finance_transactions"
}
