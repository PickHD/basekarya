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

func (r *repository) dbFromCtx(ctx context.Context) *gorm.DB {
	return utils.GetDBFromContext(ctx, r.db)
}

func (r *repository) FindActiveByType(ctx context.Context, bpjsType string, effectiveDate time.Time) ([]BPJSRateConfig, error) {
	var configs []BPJSRateConfig
	err := r.dbFromCtx(ctx).
		Where("type = ? AND is_active = ? AND effective_from <= ? AND (effective_until IS NULL OR effective_until >= ?)", bpjsType, true, effectiveDate, effectiveDate).
		Find(&configs).Error
	return configs, err
}

func (r *repository) FindAllActive(ctx context.Context, effectiveDate time.Time) ([]BPJSRateConfig, error) {
	var configs []BPJSRateConfig
	err := r.dbFromCtx(ctx).
		Where("is_active = ? AND effective_from <= ? AND (effective_until IS NULL OR effective_until >= ?)", true, effectiveDate, effectiveDate).
		Find(&configs).Error
	return configs, err
}

func (r *repository) Create(ctx context.Context, config *BPJSRateConfig) error {
	return r.dbFromCtx(ctx).Create(config).Error
}

func (r *repository) FindByID(ctx context.Context, id uint) (*BPJSRateConfig, error) {
	var config BPJSRateConfig
	err := r.dbFromCtx(ctx).Where("id = ?", id).First(&config).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("BPJS rate config not found")
		}
		return nil, err
	}
	return &config, nil
}

func (r *repository) Update(ctx context.Context, config *BPJSRateConfig) error {
	return r.dbFromCtx(ctx).Save(config).Error
}

func (r *repository) Delete(ctx context.Context, id uint) error {
	return r.dbFromCtx(ctx).Model(&BPJSRateConfig{}).Where("id = ?", id).Update("deleted_at", time.Now()).Error
}

func (r *repository) List(ctx context.Context, filter BPJSRateConfigFilter) ([]BPJSRateConfig, int64, error) {
	var configs []BPJSRateConfig
	var total int64
	query := r.dbFromCtx(ctx).Model(&BPJSRateConfig{})
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
	if err := query.Offset(offset).Limit(filter.Limit).Order("type ASC").Find(&configs).Error; err != nil {
		return nil, 0, err
	}
	return configs, total, nil
}
