package subscription

import "time"

type SubscriptionPlan struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Name         string    `gorm:"type:varchar(100);not null" json:"name"`
	Slug         string    `gorm:"type:varchar(50);unique;not null" json:"slug"`
	MaxEmployees int       `gorm:"not null;default:0" json:"max_employees"`
	PriceMonthly float64   `gorm:"type:decimal(10,2);not null;default:0" json:"price_monthly"`
	Features     string    `gorm:"type:json" json:"features"`
	IsActive     bool      `gorm:"default:true" json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (SubscriptionPlan) TableName() string {
	return "subscription_plans"
}

type SubscriptionRequest struct {
	ID              uint       `gorm:"primaryKey" json:"id"`
	CompanyID       uint       `gorm:"not null" json:"company_id"`
	CurrentPlanID   uint       `gorm:"not null" json:"current_plan_id"`
	RequestedPlanID uint       `gorm:"not null" json:"requested_plan_id"`
	Status          string     `gorm:"type:varchar(20);not null;default:'pending'" json:"status"`
	RequestedBy     *uint      `json:"requested_by"`
	ReviewedBy      *uint      `json:"reviewed_by"`
	ReviewedAt      *time.Time `json:"reviewed_at"`
	Notes           string     `gorm:"type:text" json:"notes"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

func (SubscriptionRequest) TableName() string {
	return "subscription_requests"
}
