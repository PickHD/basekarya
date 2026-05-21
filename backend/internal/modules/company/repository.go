package company

import (
	"context"
	"encoding/json"
	"basekarya-backend/pkg/utils"

	"gorm.io/gorm"
)

type Repository interface {
	FindByID(ctx context.Context, id uint) (*Company, error)
	Update(ctx context.Context, company *Company) error
	CreateCompany(ctx context.Context, c *Company) error
	FindPlanIDBySlug(ctx context.Context, slug string) (uint, error)
	FindPlanByCompanyID(ctx context.Context, companyID uint) (string, int, string, error)
	FindModulesByCompanyID(ctx context.Context, companyID uint) ([]string, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) FindByID(ctx context.Context, id uint) (*Company, error) {
	var company Company
	err := utils.GetDBFromContext(ctx, r.db).
		First(&company, id).Error

	return &company, err
}

func (r *repository) Update(ctx context.Context, company *Company) error {
	return utils.GetDBFromContext(ctx, r.db).Save(company).Error
}

func (r *repository) CreateCompany(ctx context.Context, c *Company) error {
	return utils.GetDBFromContext(ctx, r.db).Create(c).Error
}

func (r *repository) FindPlanIDBySlug(ctx context.Context, slug string) (uint, error) {
	var id uint
	err := r.db.Table("subscription_plans").Where("slug = ? AND is_active = ?", slug, true).Select("id").Scan(&id).Error
	return id, err
}

type planInfo struct {
	Name         string
	MaxEmployees int
	Features     string
}

func (r *repository) FindPlanByCompanyID(ctx context.Context, companyID uint) (string, int, string, error) {
	var plan planInfo
	err := r.db.Table("subscription_plans").
		Select("subscription_plans.name, subscription_plans.max_employees, subscription_plans.features").
		Joins("JOIN companies ON companies.subscription_plan_id = subscription_plans.id").
		Where("companies.id = ?", companyID).
		Scan(&plan).Error
	return plan.Name, plan.MaxEmployees, plan.Features, err
}

func (r *repository) FindModulesByCompanyID(ctx context.Context, companyID uint) ([]string, error) {
	var featuresJSON string
	err := r.db.Table("subscription_plans").
		Select("subscription_plans.features").
		Joins("JOIN companies ON companies.subscription_plan_id = subscription_plans.id").
		Where("companies.id = ?", companyID).
		Scan(&featuresJSON).Error
	if err != nil {
		return nil, err
	}

	type featureStruct struct {
		Modules []string `json:"modules"`
	}

	var features featureStruct
	if err := json.Unmarshal([]byte(featuresJSON), &features); err != nil {
		return nil, err
	}

	return features.Modules, nil
}
