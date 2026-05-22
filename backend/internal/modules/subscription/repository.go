package subscription

import (
	"basekarya-backend/pkg/constants"
	"context"

	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) FindAllPlans(ctx context.Context) ([]SubscriptionPlan, error) {
	var plans []SubscriptionPlan
	err := r.db.Where("is_active = ?", true).Order("price_monthly ASC").Find(&plans).Error
	return plans, err
}

func (r *repository) FindPlanBySlug(ctx context.Context, slug string) (*SubscriptionPlan, error) {
	var plan SubscriptionPlan
	err := r.db.Where("slug = ? AND is_active = ?", slug, true).First(&plan).Error
	return &plan, err
}

func (r *repository) FindPlanByID(ctx context.Context, id uint) (*SubscriptionPlan, error) {
	var plan SubscriptionPlan
	err := r.db.First(&plan, id).Error
	return &plan, err
}

func (r *repository) CreateRequest(ctx context.Context, req *SubscriptionRequest) error {
	return r.db.WithContext(ctx).Create(req).Error
}

func (r *repository) FindPendingRequestByCompanyID(ctx context.Context, companyID uint) (*SubscriptionRequest, error) {
	var req SubscriptionRequest
	err := r.db.WithContext(ctx).Where("company_id = ? AND status = ?", companyID, constants.SubReqStatusPending).First(&req).Error
	return &req, err
}

func (r *repository) FindRequestByID(ctx context.Context, id uint) (*SubscriptionRequest, error) {
	var req SubscriptionRequest
	err := r.db.WithContext(ctx).First(&req, id).Error
	return &req, err
}

func (r *repository) FindAllPendingRequests(ctx context.Context) ([]SubscriptionRequestResponse, error) {
	var results []SubscriptionRequestResponse
	err := r.db.Table("subscription_requests").
		Select(`subscription_requests.id,
			companies.name as company_name,
			cp.name as current_plan_name,
			cp.price_monthly as current_plan_price,
			rp.name as requested_plan_name,
			rp.price_monthly as requested_plan_price,
			(rp.price_monthly - cp.price_monthly) as price_difference,
			subscription_requests.status,
			u.username as requested_by_name,
			u.username as requested_by_email,
			subscription_requests.notes,
			subscription_requests.created_at`).
		Joins("JOIN companies ON companies.id = subscription_requests.company_id").
		Joins("JOIN subscription_plans cp ON cp.id = subscription_requests.current_plan_id").
		Joins("JOIN subscription_plans rp ON rp.id = subscription_requests.requested_plan_id").
		Joins("LEFT JOIN users u ON u.id = subscription_requests.requested_by").
		Where("subscription_requests.status = ?", constants.SubReqStatusPending).
		Order("subscription_requests.created_at DESC").
		Scan(&results).Error
	return results, err
}

func (r *repository) UpdateRequest(ctx context.Context, req *SubscriptionRequest) error {
	return r.db.WithContext(ctx).Save(req).Error
}

func (r *repository) FindAllRequests(ctx context.Context) ([]SubscriptionRequestResponse, error) {
	var results []SubscriptionRequestResponse
	err := r.db.Table("subscription_requests").
		Select(`subscription_requests.id,
			companies.name as company_name,
			cp.name as current_plan_name,
			cp.price_monthly as current_plan_price,
			rp.name as requested_plan_name,
			rp.price_monthly as requested_plan_price,
			(rp.price_monthly - cp.price_monthly) as price_difference,
			subscription_requests.status,
			u.username as requested_by_name,
			u.username as requested_by_email,
			subscription_requests.notes,
			subscription_requests.created_at`).
		Joins("JOIN companies ON companies.id = subscription_requests.company_id").
		Joins("JOIN subscription_plans cp ON cp.id = subscription_requests.current_plan_id").
		Joins("JOIN subscription_plans rp ON rp.id = subscription_requests.requested_plan_id").
		Joins("LEFT JOIN users u ON u.id = subscription_requests.requested_by").
		Order("subscription_requests.created_at DESC").
		Scan(&results).Error
	return results, err
}

func (r *repository) FindAllCompanies(ctx context.Context, search string) ([]CompanyListItem, error) {
	var results []CompanyListItem
	q := r.db.Table("companies").
		Select(`companies.id, companies.name, companies.email, companies.phone_number,
			sp.name as plan_name, sp.slug as plan_slug,
			companies.subscription_status, companies.subscription_expires_at,
			companies.created_at,
			(SELECT COUNT(*) FROM users WHERE users.company_id = companies.id AND users.is_active = 1 AND users.is_platform_admin = 0) as employee_count`).
		Joins("LEFT JOIN subscription_plans sp ON sp.id = companies.subscription_plan_id").
		Where("companies.id != 0")

	if search != "" {
		q = q.Where("companies.name LIKE ? OR companies.email LIKE ?", "%"+search+"%", "%"+search+"%")
	}

	err := q.Order("companies.created_at DESC").Scan(&results).Error
	return results, err
}

func (r *repository) FindCompanyDetailByID(ctx context.Context, id uint) (*CompanyDetail, error) {
	var result CompanyDetail
	err := r.db.Table("companies").
		Select(`companies.id, companies.name, companies.email, companies.phone_number, companies.address,
			sp.name as plan_name, sp.slug as plan_slug,
			sp.max_employees, sp.price_monthly,
			companies.subscription_status, companies.subscription_expires_at,
			companies.created_at,
			(SELECT COUNT(*) FROM users WHERE users.company_id = companies.id AND users.is_active = 1 AND users.is_platform_admin = 0) as employee_count`).
		Joins("LEFT JOIN subscription_plans sp ON sp.id = companies.subscription_plan_id").
		Where("companies.id = ?", id).
		Scan(&result).Error
	return &result, err
}

func (r *repository) UpdateCompanyStatus(ctx context.Context, companyID uint, status string) error {
	return r.db.Table("companies").
		Where("id = ?", companyID).
		Update("subscription_status", status).Error
}

func (r *repository) GetDashboardStats(ctx context.Context) (*DashboardStatsResponse, error) {
	var stats DashboardStatsResponse

	type baseStats struct {
		TotalCompanies      int
		ActiveSubscriptions int
		PendingPayments     int
	}
	var base baseStats

	err := r.db.Raw(`
		SELECT
			(SELECT COUNT(*) FROM companies WHERE id != 0) as total_companies,
			(SELECT COUNT(*) FROM companies WHERE subscription_status = ? AND id != 0) as active_subscriptions,
			(SELECT COUNT(*) FROM companies WHERE subscription_status = ? AND id != 0) as pending_payments
	`, constants.SubStatusActive, constants.SubStatusPendingPayment).Scan(&base).Error
	if err != nil {
		return nil, err
	}

	stats.TotalCompanies = base.TotalCompanies
	stats.ActiveSubscriptions = base.ActiveSubscriptions
	stats.PendingPayments = base.PendingPayments

	var distributions []PlanDistribution
	err = r.db.Raw(`
		SELECT sp.name as plan_name, sp.slug as plan_slug,
			COUNT(c.id) as count,
			COUNT(c.id) * sp.price_monthly as revenue
		FROM subscription_plans sp
		LEFT JOIN companies c ON c.subscription_plan_id = sp.id AND c.subscription_status = ?
		GROUP BY sp.id, sp.name, sp.slug
		ORDER BY sp.price_monthly ASC
	`, constants.SubStatusActive).Scan(&distributions).Error
	if err != nil {
		return nil, err
	}

	stats.PlanDistribution = distributions

	var totalRevenue float64
	for _, d := range distributions {
		totalRevenue += d.Revenue
	}
	stats.TotalRevenue = totalRevenue

	return &stats, nil
}
