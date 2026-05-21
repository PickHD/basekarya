package company

import "time"

type Company struct {
	ID                    uint       `gorm:"primaryKey" json:"id"`
	Name                  string     `gorm:"type:varchar(255);not null" json:"name"`
	Address               string     `gorm:"type:text" json:"address"`
	Email                 string     `gorm:"type:varchar(255)" json:"email"`
	PhoneNumber           string     `gorm:"type:varchar(50)" json:"phone_number"`
	Website               string     `json:"website"`
	TaxNumber             string     `json:"tax_number"`
	LogoURL               string     `json:"logo_url"`
	SubscriptionPlanID    *uint      `json:"subscription_plan_id"`
	SubscriptionStatus    string     `gorm:"type:varchar(20);default:'ACTIVE'" json:"subscription_status"`
	SubscriptionExpiresAt *time.Time `json:"subscription_expires_at"`
	OwnerUserID           *uint      `json:"owner_user_id"`
	CreatedAt             time.Time  `json:"created_at"`
	UpdatedAt             time.Time  `json:"updated_at"`
}

func (Company) TableName() string {
	return "companies"
}
