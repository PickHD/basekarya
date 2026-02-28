package overtime

import (
	"basekarya-backend/pkg/constants"
	"basekarya-backend/pkg/utils"
	"context"

	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, overtime *Overtime) error
	FindByID(ctx context.Context, id uint) (*Overtime, error)
	FindAll(ctx context.Context, filter OvertimeFilter) ([]Overtime, int64, error)
	GetBulkActiveOvertimesByEmployeeIds(ctx context.Context, month, year int, ids []uint) (map[uint]int, error)
	UpdateBulkStatusByEmployeeId(ctx context.Context, employeeID uint, periodMonth, periodYear int, status constants.OvertimeStatus) error
	Update(ctx context.Context, overtime *Overtime) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) Create(ctx context.Context, overtime *Overtime) error {
	db := utils.GetDBFromContext(ctx, r.db)
	return db.Create(overtime).Error
}

func (r *repository) FindByID(ctx context.Context, id uint) (*Overtime, error) {
	db := utils.GetDBFromContext(ctx, r.db)
	var overtime Overtime

	err := db.
		Preload("User").
		Preload("Employee").First(&overtime, id).Error
	if err != nil {
		return nil, err
	}

	return &overtime, nil
}

func (r *repository) FindAll(ctx context.Context, filter OvertimeFilter) ([]Overtime, int64, error) {
	db := utils.GetDBFromContext(ctx, r.db)
	var overtimes []Overtime
	var total int64

	query := db.Model(&Overtime{}).
		Joins("JOIN users ON users.id = overtimes.user_id").
		Joins("JOIN employees ON employees.id = overtimes.employee_id").
		Preload("User").
		Preload("Employee")

	if filter.UserID > 0 {
		query = query.Where("users.id = ?", filter.UserID)
	}

	if filter.Status != "" {
		query = query.Where("overtimes.status = ?", filter.Status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (filter.Page - 1) * filter.Limit
	err := query.
		Limit(filter.Limit).
		Offset(offset).
		Order("overtimes.created_at DESC").
		Find(&overtimes).Error

	return overtimes, total, err
}

func (r *repository) Update(ctx context.Context, overtime *Overtime) error {
	db := utils.GetDBFromContext(ctx, r.db)
	return db.Save(overtime).Error
}

func (r *repository) GetBulkActiveOvertimesByEmployeeIds(ctx context.Context, month, year int, ids []uint) (map[uint]int, error) {
	db := utils.GetDBFromContext(ctx, r.db)
	type Result struct {
		EmployeeID      uint
		DurationMinutes int
	}

	var results []Result

	err := db.Model(&Overtime{}).
		Select("employee_id, SUM(duration_minutes) as duration_minutes").
		Where("status = ?", string(constants.OvertimeStatusApproved)).
		Where("MONTH(date) = ? AND YEAR(date) = ?", month, year).
		Where("employee_id IN ?", ids).
		Group("employee_id").
		Scan(&results).Error
	if err != nil {
		return nil, err
	}

	dataMap := make(map[uint]int)
	for _, res := range results {
		dataMap[res.EmployeeID] = res.DurationMinutes
	}

	return dataMap, nil
}

func (r *repository) UpdateBulkStatusByEmployeeId(ctx context.Context, employeeID uint, periodMonth, periodYear int, status constants.OvertimeStatus) error {
	db := utils.GetDBFromContext(ctx, r.db)
	return db.Model(&Overtime{}).
		Where("employee_id = ?", employeeID).
		Where("status = ?", string(constants.OvertimeStatusApproved)).
		Where("MONTH(date) = ? AND YEAR(date) = ?", periodMonth, periodYear).
		Update("status", string(status)).Error
}
