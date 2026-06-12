package asset

import (
	"basekarya-backend/pkg/utils"
	"context"

	"gorm.io/gorm"
)

type Repository interface {
	CreateCategory(ctx context.Context, category *AssetCategory) error
	FindCategoryByID(ctx context.Context, id uint) (*AssetCategory, error)
	FindAllCategories(ctx context.Context, filter AssetCategoryFilter) ([]AssetCategory, int64, error)
	UpdateCategory(ctx context.Context, category *AssetCategory) error
	DeleteCategory(ctx context.Context, id uint) error

	CreateAsset(ctx context.Context, asset *Asset) error
	FindAssetByID(ctx context.Context, id uint) (*Asset, error)
	FindAllAssets(ctx context.Context, filter AssetFilter) ([]Asset, int64, error)
	UpdateAsset(ctx context.Context, asset *Asset) error
	DeleteAsset(ctx context.Context, id uint) error

	CreateAssignment(ctx context.Context, assignment *AssetAssignment) error
	FindAssignmentByID(ctx context.Context, id uint) (*AssetAssignment, error)
	FindActiveAssignmentByAssetID(ctx context.Context, assetID uint) (*AssetAssignment, error)
	FindAllAssignments(ctx context.Context, filter AssetAssignmentFilter) ([]AssetAssignment, int64, error)
	UpdateAssignment(ctx context.Context, assignment *AssetAssignment) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) CreateCategory(ctx context.Context, category *AssetCategory) error {
	db := utils.TenantScope(ctx, utils.GetDBFromContext(ctx, r.db))
	return db.Create(category).Error
}

func (r *repository) FindCategoryByID(ctx context.Context, id uint) (*AssetCategory, error) {
	db := utils.TenantScope(ctx, utils.GetDBFromContext(ctx, r.db))
	var category AssetCategory
	err := db.First(&category, id).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *repository) FindAllCategories(ctx context.Context, filter AssetCategoryFilter) ([]AssetCategory, int64, error) {
	db := utils.GetDBFromContext(ctx, r.db)
	var categories []AssetCategory
	var total int64

	query := utils.TenantScope(ctx, db.Model(&AssetCategory{}))

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (filter.Page - 1) * filter.Limit
	err := query.
		Limit(filter.Limit).
		Offset(offset).
		Order("asset_categories.created_at DESC").
		Find(&categories).Error

	return categories, total, err
}

func (r *repository) UpdateCategory(ctx context.Context, category *AssetCategory) error {
	db := utils.TenantScope(ctx, utils.GetDBFromContext(ctx, r.db))
	return db.Save(category).Error
}

func (r *repository) DeleteCategory(ctx context.Context, id uint) error {
	db := utils.TenantScope(ctx, utils.GetDBFromContext(ctx, r.db))
	return db.Delete(&AssetCategory{}, id).Error
}

func (r *repository) CreateAsset(ctx context.Context, asset *Asset) error {
	db := utils.TenantScope(ctx, utils.GetDBFromContext(ctx, r.db))
	return db.Create(asset).Error
}

func (r *repository) FindAssetByID(ctx context.Context, id uint) (*Asset, error) {
	db := utils.TenantScope(ctx, utils.GetDBFromContext(ctx, r.db))
	var asset Asset
	err := db.Preload("AssetCategory").First(&asset, id).Error
	if err != nil {
		return nil, err
	}
	return &asset, nil
}

func (r *repository) FindAllAssets(ctx context.Context, filter AssetFilter) ([]Asset, int64, error) {
	db := utils.GetDBFromContext(ctx, r.db)
	var assets []Asset
	var total int64

	query := utils.TenantScope(ctx, db.Model(&Asset{})).
		Preload("AssetCategory")

	if filter.Status != "" {
		query = query.Where("assets.status = ?", filter.Status)
	}
	if filter.Condition != "" {
		query = query.Where("assets.condition = ?", filter.Condition)
	}
	if filter.CategoryID > 0 {
		query = query.Where("assets.asset_category_id = ?", filter.CategoryID)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (filter.Page - 1) * filter.Limit
	err := query.
		Limit(filter.Limit).
		Offset(offset).
		Order("assets.created_at DESC").
		Find(&assets).Error

	return assets, total, err
}

func (r *repository) UpdateAsset(ctx context.Context, asset *Asset) error {
	db := utils.TenantScope(ctx, utils.GetDBFromContext(ctx, r.db))
	return db.Save(asset).Error
}

func (r *repository) DeleteAsset(ctx context.Context, id uint) error {
	db := utils.TenantScope(ctx, utils.GetDBFromContext(ctx, r.db))
	return db.Delete(&Asset{}, id).Error
}

func (r *repository) CreateAssignment(ctx context.Context, assignment *AssetAssignment) error {
	db := utils.TenantScope(ctx, utils.GetDBFromContext(ctx, r.db))
	return db.Create(assignment).Error
}

func (r *repository) FindAssignmentByID(ctx context.Context, id uint) (*AssetAssignment, error) {
	db := utils.TenantScope(ctx, utils.GetDBFromContext(ctx, r.db))
	var assignment AssetAssignment
	err := db.
		Preload("Asset").
		Preload("User").
		Preload("Employee").
		First(&assignment, id).Error
	if err != nil {
		return nil, err
	}
	return &assignment, nil
}

func (r *repository) FindActiveAssignmentByAssetID(ctx context.Context, assetID uint) (*AssetAssignment, error) {
	db := utils.TenantScope(ctx, utils.GetDBFromContext(ctx, r.db))
	var assignment AssetAssignment
	err := db.
		Where("asset_id = ?", assetID).
		Where("status = ?", "ACTIVE").
		First(&assignment).Error
	if err != nil {
		return nil, err
	}
	return &assignment, nil
}

func (r *repository) FindAllAssignments(ctx context.Context, filter AssetAssignmentFilter) ([]AssetAssignment, int64, error) {
	db := utils.GetDBFromContext(ctx, r.db)
	var assignments []AssetAssignment
	var total int64

	query := utils.TenantScope(ctx, db.Model(&AssetAssignment{})).
		Joins("JOIN users ON users.id = asset_assignments.user_id").
		Joins("JOIN employees ON employees.id = asset_assignments.employee_id").
		Preload("Asset").
		Preload("User").
		Preload("Employee")

	if filter.UserID > 0 {
		query = query.Where("users.id = ?", filter.UserID)
	}
	if filter.Status != "" {
		query = query.Where("asset_assignments.status = ?", filter.Status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (filter.Page - 1) * filter.Limit
	err := query.
		Limit(filter.Limit).
		Offset(offset).
		Order("asset_assignments.created_at DESC").
		Find(&assignments).Error

	return assignments, total, err
}

func (r *repository) UpdateAssignment(ctx context.Context, assignment *AssetAssignment) error {
	db := utils.TenantScope(ctx, utils.GetDBFromContext(ctx, r.db))
	return db.Save(assignment).Error
}
