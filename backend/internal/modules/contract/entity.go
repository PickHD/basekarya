package contract

import (
	"basekarya-backend/internal/modules/user"
	"basekarya-backend/pkg/constants"
	"time"

	"gorm.io/gorm"
)

type Contract struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	EmployeeID uint `gorm:"unique;not null" json:"employee_id"`

	ContractType   constants.ContractType `gorm:"type:varchar(20);not null" json:"contract_type"`
	ContractNumber string                 `gorm:"type:varchar(50)" json:"contract_number"`

	StartDate time.Time  `gorm:"type:date;not null" json:"start_date"`
	EndDate   *time.Time `gorm:"type:date" json:"end_date"`

	Notes         string `gorm:"type:text" json:"notes"`
	AttachmentURL string `gorm:"type:varchar(255)" json:"attachment_url"`

	AlertedAt *time.Time `json:"alerted_at"`

	Employee *user.Employee `gorm:"foreignKey:EmployeeID" json:"employee,omitempty"`
}

func (Contract) TableName() string {
	return "contracts"
}
