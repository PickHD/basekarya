package bpjs

import (
	"basekarya-backend/pkg/utils"
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
)

type Repository interface {
	FindActiveByType(ctx context.Context, bpjsType string, effectiveDate time.Time) ([]BPJSRateConfig, error)
	FindAllActive(ctx context.Context, effectiveDate time.Time) ([]BPJSRateConfig, error)
	Create(ctx context.Context, config *BPJSRateConfig) error
	FindByID(ctx context.Context, id uint) (*BPJSRateConfig, error)
	Update(ctx context.Context, config *BPJSRateConfig) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, filter BPJSRateConfigFilter) ([]BPJSRateConfig, int64, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) scopedDB(ctx context.Context) *gorm.DB {
	return utils.TenantScope(ctx, utils.GetDBFromContext(ctx, r.db))
}

func (r *repository) FindActiveByType(ctx context.Context, bpjsType string, effectiveDate time.Time) ([]BPJSRateConfig, error) {
	var configs []BPJSRateConfig
	query := utils.GetDBFromContext(ctx, r.db).
		Where("type = ? AND is_active = ? AND effective_from <= ? AND (effective_until IS NULL OR effective_until >= ?)", bpjsType, true, effectiveDate, effectiveDate)

	companyID := utils.GetCompanyIDFromCtx(ctx)
	if companyID > 0 {
		query = query.Where("company_id IS NULL OR company_id = ?", companyID)
	} else {
		query = query.Where("company_id IS NULL")
	}

	err := query.Order("company_id ASC").Find(&configs).Error
	return configs, err
}

func (r *repository) FindAllActive(ctx context.Context, effectiveDate time.Time) ([]BPJSRateConfig, error) {
	var configs []BPJSRateConfig
	query := utils.GetDBFromContext(ctx, r.db).
		Where("is_active = ? AND effective_from <= ? AND (effective_until IS NULL OR effective_until >= ?)", true, effectiveDate, effectiveDate)

	companyID := utils.GetCompanyIDFromCtx(ctx)
	if companyID > 0 {
		query = query.Where("company_id IS NULL OR company_id = ?", companyID)
	} else {
		query = query.Where("company_id IS NULL")
	}

	err := query.Order("type ASC, company_id ASC").Find(&configs).Error
	return configs, err
}

func (r *repository) Create(ctx context.Context, config *BPJSRateConfig) error {
	return utils.GetDBFromContext(ctx, r.db).Create(config).Error
}

func (r *repository) FindByID(ctx context.Context, id uint) (*BPJSRateConfig, error) {
	var config BPJSRateConfig
	err := r.scopedDB(ctx).Where("id = ?", id).First(&config).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("BPJS rate config not found")
		}
		return nil, err
	}
	return &config, nil
}

func (r *repository) Update(ctx context.Context, config *BPJSRateConfig) error {
	return r.scopedDB(ctx).Save(config).Error
}

func (r *repository) Delete(ctx context.Context, id uint) error {
	return r.scopedDB(ctx).Model(&BPJSRateConfig{}).Where("id = ?", id).Update("deleted_at", time.Now()).Error
}

func (r *repository) List(ctx context.Context, filter BPJSRateConfigFilter) ([]BPJSRateConfig, int64, error) {
	var configs []BPJSRateConfig
	var total int64
	query := utils.GetDBFromContext(ctx, r.db).Model(&BPJSRateConfig{})

	companyID := utils.GetCompanyIDFromCtx(ctx)
	if companyID > 0 {
		query = query.Where("company_id IS NULL OR company_id = ?", companyID)
	} else {
		query = query.Where("company_id IS NULL")
	}

	if filter.Type != "" {
		query = query.Where("type = ?", filter.Type)
	}
	if filter.IsActive != nil {
		query = query.Where("is_active = ?", *filter.IsActive)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (filter.Page - 1) * filter.Limit
	if err := query.Offset(offset).Limit(filter.Limit).Order("type ASC, company_id ASC").Find(&configs).Error; err != nil {
		return nil, 0, err
	}
	return configs, total, nil
}
