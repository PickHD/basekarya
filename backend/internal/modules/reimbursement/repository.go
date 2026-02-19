package reimbursement

import (
	"context"
	"hris-backend/pkg/utils"

	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, reimbursement *Reimbursement) error
	FindByID(ctx context.Context, id uint) (*Reimbursement, error)
	FindAll(ctx context.Context, filter ReimbursementFilter) ([]Reimbursement, int64, error)
	Update(ctx context.Context, reimbursement *Reimbursement) error
	GetBulkApprovedAmount(ctx context.Context, month, year int) (map[uint]float64, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) Create(ctx context.Context, reimbursement *Reimbursement) error {
	db := utils.GetDBFromContext(ctx, r.db)
	return db.Create(reimbursement).Error
}

func (r *repository) FindByID(ctx context.Context, id uint) (*Reimbursement, error) {
	db := utils.GetDBFromContext(ctx, r.db)
	var reimburstment Reimbursement

	err := db.Preload("User").First(&reimburstment, id).Error
	if err != nil {
		return nil, err
	}

	return &reimburstment, nil
}

func (r *repository) FindAll(ctx context.Context, filter ReimbursementFilter) ([]Reimbursement, int64, error) {
	db := utils.GetDBFromContext(ctx, r.db)
	var reimbursements []Reimbursement
	var total int64

	query := db.Model(&Reimbursement{})

	if filter.UserID > 0 {
		query = query.Where("user_id = ?", filter.UserID)
	}

	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (filter.Page - 1) * filter.Limit
	err := query.
		Limit(filter.Limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&reimbursements).Error

	return reimbursements, total, err
}

func (r *repository) Update(ctx context.Context, reimbursement *Reimbursement) error {
	db := utils.GetDBFromContext(ctx, r.db)
	return db.Save(reimbursement).Error
}

func (r *repository) GetBulkApprovedAmount(ctx context.Context, month, year int) (map[uint]float64, error) {
	db := utils.GetDBFromContext(ctx, r.db)
	type Result struct {
		UserID      uint
		TotalAmount float64
	}
	var results []Result

	err := db.Model(&Reimbursement{}).
		Select("user_id, COALESCE(SUM(amount), 0) as total_amount").
		Where("status = ?", "APPROVED").
		Where("MONTH(date_of_expense) = ? AND YEAR(date_of_expense) = ?", month, year).
		Group("user_id").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	dataMap := make(map[uint]float64)
	for _, res := range results {
		dataMap[res.UserID] = res.TotalAmount
	}

	return dataMap, err
}
