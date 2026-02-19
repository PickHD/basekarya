package company

import (
	"context"
	"hris-backend/pkg/utils"

	"gorm.io/gorm"
)

type Repository interface {
	FindByID(ctx context.Context, id uint) (*Company, error)
	Update(ctx context.Context, company *Company) error
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
