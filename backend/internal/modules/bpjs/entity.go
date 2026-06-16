package bpjs

import (
	"time"

	"gorm.io/gorm"
)

type BPJSRateConfig struct {
	ID                uint           `gorm:"primaryKey" json:"id"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"-"`
	CompanyID         *uint          `gorm:"index" json:"company_id"`
	Type              string         `gorm:"type:varchar(20);not null" json:"type"`
	EmployeeRate      float64        `gorm:"type:decimal(5,4);not null;default:0" json:"employee_rate"`
	EmployerRate      float64        `gorm:"type:decimal(5,4);not null;default:0" json:"employer_rate"`
	MaxSalaryCap      *float64       `gorm:"type:decimal(15,2)" json:"max_salary_cap"`
	IndustryRiskLevel *string        `gorm:"type:varchar(10)" json:"industry_risk_level"`
	IsActive          bool           `gorm:"default:true" json:"is_active"`
	EffectiveFrom     time.Time      `gorm:"type:date;not null" json:"effective_from"`
	EffectiveUntil    *time.Time     `gorm:"type:date" json:"effective_until"`
}

func (BPJSRateConfig) TableName() string { return "bpjs_rate_configs" }
