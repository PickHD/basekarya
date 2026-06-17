package tax

import (
	"basekarya-backend/pkg/utils"
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
)

type Repository interface {
	FindTERBrackets(ctx context.Context, category string, effectiveDate time.Time) ([]TERBracket, error)
	CreateTERBracket(ctx context.Context, bracket *TERBracket) error
	FindTERBracketByID(ctx context.Context, id uint) (*TERBracket, error)
	UpdateTERBracket(ctx context.Context, bracket *TERBracket) error
	DeleteTERBracket(ctx context.Context, id uint) error
	ListTERBrackets(ctx context.Context, filter TERBracketFilter) ([]TERBracket, int64, error)
	FindPTKPByYear(ctx context.Context, year int) ([]PTKPConfig, error)
	CreatePTKPConfig(ctx context.Context, ptkp *PTKPConfig) error
	FindPTKPConfigByID(ctx context.Context, id uint) (*PTKPConfig, error)
	UpdatePTKPConfig(ctx context.Context, ptkp *PTKPConfig) error
	DeletePTKPConfig(ctx context.Context, id uint) error
	ListPTKPConfigs(ctx context.Context, year int) ([]PTKPConfig, int64, error)
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

func (r *repository) FindTERBrackets(ctx context.Context, category string, effectiveDate time.Time) ([]TERBracket, error) {
	var brackets []TERBracket
	query := utils.GetDBFromContext(ctx, r.db).
		Where("category = ? AND effective_from <= ? AND (effective_until IS NULL OR effective_until >= ?)", category, effectiveDate, effectiveDate)

	companyID := utils.GetCompanyIDFromCtx(ctx)
	if companyID > 0 {
		query = query.Where("company_id IS NULL OR company_id = ?", companyID)
	} else {
		query = query.Where("company_id IS NULL")
	}

	err := query.Order("bracket_number ASC, company_id ASC").Find(&brackets).Error
	return brackets, err
}

func (r *repository) CreateTERBracket(ctx context.Context, bracket *TERBracket) error {
	return utils.GetDBFromContext(ctx, r.db).Create(bracket).Error
}

func (r *repository) FindTERBracketByID(ctx context.Context, id uint) (*TERBracket, error) {
	var bracket TERBracket
	err := r.scopedDB(ctx).Where("id = ?", id).First(&bracket).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("TER bracket not found")
		}
		return nil, err
	}
	return &bracket, nil
}

func (r *repository) UpdateTERBracket(ctx context.Context, bracket *TERBracket) error {
	return r.scopedDB(ctx).Save(bracket).Error
}

func (r *repository) DeleteTERBracket(ctx context.Context, id uint) error {
	return r.scopedDB(ctx).Model(&TERBracket{}).Where("id = ?", id).Update("deleted_at", time.Now()).Error
}

func (r *repository) ListTERBrackets(ctx context.Context, filter TERBracketFilter) ([]TERBracket, int64, error) {
	var brackets []TERBracket
	var total int64
	query := utils.GetDBFromContext(ctx, r.db).Model(&TERBracket{})

	companyID := utils.GetCompanyIDFromCtx(ctx)
	if companyID > 0 {
		query = query.Where("company_id IS NULL OR company_id = ?", companyID)
	} else {
		query = query.Where("company_id IS NULL")
	}

	if filter.Category != "" {
		query = query.Where("category = ?", filter.Category)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (filter.Page - 1) * filter.Limit
	if err := query.Offset(offset).Limit(filter.Limit).Order("category ASC, bracket_number ASC, company_id ASC").Find(&brackets).Error; err != nil {
		return nil, 0, err
	}
	return brackets, total, nil
}

func (r *repository) FindPTKPByYear(ctx context.Context, year int) ([]PTKPConfig, error) {
	var configs []PTKPConfig
	err := utils.GetDBFromContext(ctx, r.db).Where("effective_year = ?", year).Find(&configs).Error
	return configs, err
}

func (r *repository) CreatePTKPConfig(ctx context.Context, ptkp *PTKPConfig) error {
	return utils.GetDBFromContext(ctx, r.db).Create(ptkp).Error
}

func (r *repository) FindPTKPConfigByID(ctx context.Context, id uint) (*PTKPConfig, error) {
	var config PTKPConfig
	err := utils.GetDBFromContext(ctx, r.db).Where("id = ?", id).First(&config).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("PTKP config not found")
		}
		return nil, err
	}
	return &config, nil
}

func (r *repository) UpdatePTKPConfig(ctx context.Context, ptkp *PTKPConfig) error {
	return utils.GetDBFromContext(ctx, r.db).Save(ptkp).Error
}

func (r *repository) DeletePTKPConfig(ctx context.Context, id uint) error {
	return utils.GetDBFromContext(ctx, r.db).Model(&PTKPConfig{}).Where("id = ?", id).Update("deleted_at", time.Now()).Error
}

func (r *repository) ListPTKPConfigs(ctx context.Context, year int) ([]PTKPConfig, int64, error) {
	var configs []PTKPConfig
	var total int64
	query := utils.GetDBFromContext(ctx, r.db).Model(&PTKPConfig{})
	if year > 0 {
		query = query.Where("effective_year = ?", year)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := query.Order("effective_year DESC, code ASC").Find(&configs).Error; err != nil {
		return nil, 0, err
	}
	return configs, total, nil
}
