package contract

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository interface {
	Upsert(ctx context.Context, contract *Contract) error
	FindByID(ctx context.Context, id uint) (*Contract, error)
	FindByEmployeeID(ctx context.Context, employeeID uint) (*Contract, error)
	FindAll(ctx context.Context, filter *ContractFilter) ([]Contract, int64, error)
	FindExpiringContracts(ctx context.Context, withinDays int) ([]Contract, error)
	MarkAlerted(ctx context.Context, ids []uint) error
	SoftDelete(ctx context.Context, id uint) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) Upsert(ctx context.Context, contract *Contract) error {
	return r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "employee_id"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"contract_type",
			"contract_number",
			"start_date",
			"end_date",
			"notes",
			"attachment_url",
			"alerted_at",
			"updated_at",
		}),
	}).Create(contract).Error
}

func (r *repository) FindByID(ctx context.Context, id uint) (*Contract, error) {
	var contract Contract
	if err := r.db.WithContext(ctx).Preload("Employee").First(&contract, id).Error; err != nil {
		return nil, err
	}
	return &contract, nil
}

func (r *repository) FindByEmployeeID(ctx context.Context, employeeID uint) (*Contract, error) {
	var contract Contract
	if err := r.db.WithContext(ctx).Preload("Employee").Where("employee_id = ?", employeeID).First(&contract).Error; err != nil {
		return nil, err
	}
	return &contract, nil
}

func (r *repository) FindAll(ctx context.Context, filter *ContractFilter) ([]Contract, int64, error) {
	var contracts []Contract
	var total int64
	query := r.db.WithContext(ctx).Model(&Contract{}).Preload("Employee")

	if filter.ContractType != "" {
		query = query.Where("contracts.contract_type = ?", filter.ContractType)
	}

	if filter.Search != "" {
		query = query.Joins("JOIN employees ON employees.id = contracts.employee_id").
			Where("employees.full_name LIKE ? OR employees.nik LIKE ?", "%"+filter.Search+"%", "%"+filter.Search+"%")
	}

	if filter.ExpiringWithinDays > 0 {
		query = query.Where("contracts.contract_type = ? AND contracts.end_date <= DATE_ADD(NOW(), INTERVAL ? DAY) AND contracts.end_date >= NOW()", "PKWT", filter.ExpiringWithinDays)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (filter.Page - 1) * filter.Limit
	if err := query.Order("contracts.created_at DESC").Offset(offset).Limit(filter.Limit).Find(&contracts).Error; err != nil {
		return nil, 0, err
	}

	return contracts, total, nil
}

func (r *repository) FindExpiringContracts(ctx context.Context, withinDays int) ([]Contract, error) {
	var contracts []Contract
	err := r.db.WithContext(ctx).
		Where("contract_type = ? AND end_date <= DATE_ADD(NOW(), INTERVAL ? DAY) AND end_date >= CURDATE() AND alerted_at IS NULL", "PKWT", withinDays).
		Find(&contracts).Error
	return contracts, err
}

func (r *repository) MarkAlerted(ctx context.Context, ids []uint) error {
	if len(ids) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Model(&Contract{}).Where("id IN ?", ids).Update("alerted_at", gorm.Expr("NOW()")).Error
}

func (r *repository) SoftDelete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&Contract{}, id).Error
}
