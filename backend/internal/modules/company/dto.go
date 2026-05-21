package company

type CompanyProfileResponse struct {
	ID                   uint   `json:"id"`
	Name                 string `json:"name"`
	Address              string `json:"address"`
	Email                string `json:"email"`
	PhoneNumber          string `json:"phone_number"`
	Website              string `json:"website"`
	TaxNumber            string `json:"tax_number"`
	LogoURL              string `json:"logo_url"`
	SubscriptionPlanName string `json:"subscription_plan_name"`
	SubscriptionStatus   string `json:"subscription_status"`
	SubscriptionExpiresAt string `json:"subscription_expires_at,omitempty"`
	MaxEmployees         int    `json:"max_employees"`
	PlanModules          string `json:"plan_modules"`
}

type UpdateCompanyProfileRequest struct {
	Name        string `form:"name" validate:"required"`
	Address     string `form:"address"`
	Email       string `form:"email"`
	PhoneNumber string `form:"phone_number"`
	Website     string `form:"website"`
	TaxNumber   string `form:"tax_number"`
}
