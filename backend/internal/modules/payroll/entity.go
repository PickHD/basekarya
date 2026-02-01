package payroll

import (
	"hris-backend/internal/modules/user"
	"hris-backend/pkg/constants"
	"time"

	"gorm.io/gorm"
)

type Payroll struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	EmployeeID uint           `json:"employee_id"`
	Employee   *user.Employee `gorm:"foreignKey:EmployeeID" json:"employee,omitempty"`
	PeriodDate time.Time      `gorm:"type:date;not null;index" json:"period_date"`

	BaseSalary     float64 `gorm:"type:decimal(15,2)" json:"base_salary"`
	TotalAllowance float64 `gorm:"type:decimal(15,2)" json:"total_allowance"`
	TotalDeduction float64 `gorm:"type:decimal(15,2)" json:"total_deduction"`
	NetSalary      float64 `gorm:"type:decimal(15,2)" json:"net_salary"`

	Status constants.PayrollStatus `gorm:"type:varchar(20);default:'DRAFT'" json:"status"`

	Notes string `gorm:"type:text" json:"notes"`

	Details []PayrollDetail `gorm:"foreignKey:PayrollID;constraint:OnDelete:CASCADE" json:"details,omitempty"`
}

type PayrollDetail struct {
	ID        uint `gorm:"primaryKey" json:"id"`
	PayrollID uint `json:"payroll_id"`

	Title string `gorm:"type:varchar(150);not null" json:"title"`

	Type constants.PayrollDetailType `gorm:"type:varchar(20);not null" json:"type"`

	Amount float64 `gorm:"type:decimal(15,2);not null" json:"amount"`
}
