package subscription

type PlanResponse struct {
	ID           uint    `json:"id"`
	Name         string  `json:"name"`
	Slug         string  `json:"slug"`
	MaxEmployees int     `json:"max_employees"`
	PriceMonthly float64 `json:"price_monthly"`
	Features     string  `json:"features"`
}

type UpgradeRequest struct {
	PlanSlug string `json:"plan_slug" validate:"required,oneof=basic pro"`
}

type UpgradeResponse struct {
	ID              uint   `json:"id"`
	RequestedPlanID uint   `json:"requested_plan_id"`
	Status          string `json:"status"`
}

type ReviewRequest struct {
	Status string `json:"status" validate:"required,oneof=APPROVED REJECTED"`
	Notes  string `json:"notes"`
}

type SubscriptionRequestResponse struct {
	ID                 uint    `json:"id"`
	CompanyName        string  `json:"company_name"`
	CurrentPlanName    string  `json:"current_plan_name"`
	RequestedPlanName  string  `json:"requested_plan_name"`
	CurrentPlanPrice   float64 `json:"current_plan_price"`
	RequestedPlanPrice float64 `json:"requested_plan_price"`
	PriceDifference    float64 `json:"price_difference"`
	Status             string  `json:"status"`
	RequestedByName    string  `json:"requested_by_name"`
	RequestedByEmail   string  `json:"requested_by_email"`
	Notes              string  `json:"notes"`
	CreatedAt          string  `json:"created_at"`
}

type CompanyListItem struct {
	ID                    uint    `json:"id"`
	Name                  string  `json:"name"`
	Email                 string  `json:"email"`
	PhoneNumber           string  `json:"phone_number"`
	PlanName              string  `json:"plan_name"`
	PlanSlug              string  `json:"plan_slug"`
	SubscriptionStatus    string  `json:"subscription_status"`
	SubscriptionExpiresAt string  `json:"subscription_expires_at"`
	EmployeeCount         int     `json:"employee_count"`
	CreatedAt             string  `json:"created_at"`
}

type CompanyDetail struct {
	ID                    uint    `json:"id"`
	Name                  string  `json:"name"`
	Email                 string  `json:"email"`
	PhoneNumber           string  `json:"phone_number"`
	Address               string  `json:"address"`
	PlanName              string  `json:"plan_name"`
	PlanSlug              string  `json:"plan_slug"`
	MaxEmployees          int     `json:"max_employees"`
	PriceMonthly          float64 `json:"price_monthly"`
	SubscriptionStatus    string  `json:"subscription_status"`
	SubscriptionExpiresAt string  `json:"subscription_expires_at"`
	EmployeeCount         int     `json:"employee_count"`
	CreatedAt             string  `json:"created_at"`
}

type UpdateCompanyStatusRequest struct {
	SubscriptionStatus    string `json:"subscription_status" validate:"required,oneof=ACTIVE PENDING_PAYMENT EXPIRED"`
}

type DashboardStatsResponse struct {
	TotalCompanies    int                `json:"total_companies"`
	ActiveSubscriptions int              `json:"active_subscriptions"`
	PendingPayments   int                `json:"pending_payments"`
	TotalRevenue      float64            `json:"total_revenue"`
	PlanDistribution  []PlanDistribution `json:"plan_distribution"`
}

type PlanDistribution struct {
	PlanName string `json:"plan_name"`
	PlanSlug string `json:"plan_slug"`
	Count    int    `json:"count"`
	Revenue  float64 `json:"revenue"`
}
