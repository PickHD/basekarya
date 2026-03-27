package loan

import (
	"basekarya-backend/pkg/constants"
	"basekarya-backend/pkg/utils"
	"context"

	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, loan *Loan) error
	FindByID(ctx context.Context, id uint) (*Loan, error)
	FindActiveLoanByUserID(ctx context.Context, userID uint) (*Loan, error)
	FindAll(ctx context.Context, filter LoanFilter) ([]Loan, int64, error)
	Update(ctx context.Context, loan *Loan) error
	GetBulkActiveLoansByEmployeeIds(ctx context.Context, ids []uint) (map[uint]Loan, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) Create(ctx context.Context, loan *Loan) error {
	db := utils.GetDBFromContext(ctx, r.db)
	return db.Create(loan).Error
}

func (r *repository) FindByID(ctx context.Context, id uint) (*Loan, error) {
	db := utils.GetDBFromContext(ctx, r.db)
	var loan Loan

	err := db.
		Preload("User").
		Preload("Employee").First(&loan, id).Error
	if err != nil {
		return nil, err
	}

	return &loan, nil
}

func (r *repository) FindActiveLoanByUserID(ctx context.Context, userID uint) (*Loan, error) {
	db := utils.GetDBFromContext(ctx, r.db)
	var loan Loan

	err := db.
		Preload("User").
		Preload("Employee").
		Where("user_id = ?", userID).
		Where("remaining_amount > 0").
		Where("status != ?", string(constants.LoanStatusPaidOff)).
		Where("status != ?", string(constants.LoanStatusRejected)).
		First(&loan).Error
	if err != nil {
		return nil, err
	}

	return &loan, nil
}

func (r *repository) FindAll(ctx context.Context, filter LoanFilter) ([]Loan, int64, error) {
	db := utils.GetDBFromContext(ctx, r.db)
	var loans []Loan
	var total int64

	query := db.Model(&Loan{}).
		Joins("JOIN users ON users.id = loans.user_id").
		Joins("JOIN employees ON employees.id = loans.employee_id").
		Preload("User").
		Preload("Employee")

	if filter.UserID > 0 {
		query = query.Where("users.id = ?", filter.UserID)
	}

	if filter.Status != "" {
		query = query.Where("loans.status = ?", filter.Status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (filter.Page - 1) * filter.Limit
	err := query.
		Limit(filter.Limit).
		Offset(offset).
		Order("loans.created_at DESC").
		Find(&loans).Error

	return loans, total, err
}

func (r *repository) Update(ctx context.Context, loan *Loan) error {
	db := utils.GetDBFromContext(ctx, r.db)
	return db.Save(loan).Error
}

func (r *repository) GetBulkActiveLoansByEmployeeIds(ctx context.Context, ids []uint) (map[uint]Loan, error) {
	db := utils.GetDBFromContext(ctx, r.db)
	type Result struct {
		EmployeeID uint
		Loan       Loan
	}
	var loans []Loan

	err := db.Model(&Loan{}).
		Where("status = ?", string(constants.LoanStatusApproved)).
		Where("remaining_amount > 0").
		Where("employee_id IN ?", ids).
		Find(&loans).Error
	if err != nil {
		return nil, err
	}

	dataMap := make(map[uint]Loan)
	for _, loan := range loans {
		dataMap[loan.EmployeeID] = loan
	}

	return dataMap, nil
}
