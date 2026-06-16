package tax

import (
	"time"

	"gorm.io/gorm"
)

type TERBracket struct {
	ID               uint           `gorm:"primaryKey" json:"id"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"`
	CompanyID        *uint          `gorm:"index" json:"company_id"`
	Category         string         `gorm:"type:char(1);not null" json:"category"`
	BracketNumber    int            `gorm:"not null" json:"bracket_number"`
	MinMonthlySalary float64        `gorm:"type:decimal(15,2);not null" json:"min_monthly_salary"`
	Rate             float64        `gorm:"type:decimal(5,4);not null" json:"rate"`
	EffectiveFrom    time.Time      `gorm:"type:date;not null" json:"effective_from"`
	EffectiveUntil   *time.Time     `gorm:"type:date" json:"effective_until"`
}

func (TERBracket) TableName() string { return "pph21_term_configs" }

type PTKPConfig struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
	Code          string         `gorm:"type:varchar(10);not null" json:"code"`
	AnnualAmount  float64        `gorm:"type:decimal(15,2);not null" json:"annual_amount"`
	EffectiveYear int            `gorm:"not null" json:"effective_year"`
}

func (PTKPConfig) TableName() string { return "ptkp_configs" }
