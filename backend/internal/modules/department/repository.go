package department

import (
	"context"
	"fmt"

	"basekarya-backend/pkg/utils"

	"gorm.io/gorm"
)

type Repository interface {
	FindAll(ctx context.Context) ([]Department, error)
	FindByID(ctx context.Context, id uint) (*Department, error)
	FindByName(ctx context.Context, name string) (*Department, error)
	Create(ctx context.Context, dept *Department) error
	Update(ctx context.Context, dept *Department) error
	Delete(ctx context.Context, id uint) error
	CountEmployees(ctx context.Context, departmentID uint) (int64, error)
	ExistsByName(ctx context.Context, name string, excludeID uint) (bool, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) FindAll(ctx context.Context) ([]Department, error) {
	var deps []Department
	db := utils.TenantScope(ctx, utils.GetDBFromContext(ctx, r.db))
	if err := db.Model(&Department{}).Find(&deps).Error; err != nil {
		return nil, err
	}
	return deps, nil
}

func (r *repository) FindByID(ctx context.Context, id uint) (*Department, error) {
	var dept Department
	db := utils.TenantScope(ctx, utils.GetDBFromContext(ctx, r.db))
	if err := db.Model(&Department{}).Where("id = ?", id).First(&dept).Error; err != nil {
		return nil, err
	}
	return &dept, nil
}

func (r *repository) FindByName(ctx context.Context, name string) (*Department, error) {
	var department Department
	db := utils.TenantScope(ctx, utils.GetDBFromContext(ctx, r.db))
	if err := db.Model(&Department{}).Where("name = ?", name).First(&department).Error; err != nil {
		return nil, err
	}
	return &department, nil
}

func (r *repository) Create(ctx context.Context, dept *Department) error {
	db := utils.TenantScope(ctx, utils.GetDBFromContext(ctx, r.db))
	return db.Create(dept).Error
}

func (r *repository) Update(ctx context.Context, dept *Department) error {
	db := utils.TenantScope(ctx, utils.GetDBFromContext(ctx, r.db))
	return db.Save(dept).Error
}

func (r *repository) Delete(ctx context.Context, id uint) error {
	db := utils.TenantScope(ctx, utils.GetDBFromContext(ctx, r.db))

	var count int64
	if err := db.Session(&gorm.Session{}).Table("employees").Where("department_id = ?", id).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return fmt.Errorf("cannot delete, department still in use by %d employees", count)
	}

	return db.Delete(&Department{}, id).Error
}

func (r *repository) CountEmployees(ctx context.Context, departmentID uint) (int64, error) {
	db := utils.TenantScope(ctx, utils.GetDBFromContext(ctx, r.db))
	var count int64
	if err := db.Table("employees").Where("department_id = ?", departmentID).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *repository) ExistsByName(ctx context.Context, name string, excludeID uint) (bool, error) {
	db := utils.TenantScope(ctx, utils.GetDBFromContext(ctx, r.db))
	var count int64
	query := db.Model(&Department{}).Where("name = ?", name)
	if excludeID != 0 {
		query = query.Where("id != ?", excludeID)
	}
	if err := query.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
