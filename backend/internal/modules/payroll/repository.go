package payroll

import (
	"hris-backend/pkg/constants"
	"time"

	"gorm.io/gorm"
)

type Repository interface {
	CreateBulk(payroll *[]Payroll) error
	FindAll(filter *PayrollFilter) ([]Payroll, int64, error)
	FindByID(id uint) (*Payroll, error)
	GetExistingEmployeeID(month, year int) (map[uint]bool, error)
	UpdateStatus(id uint, status constants.PayrollStatus) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) CreateBulk(payroll *[]Payroll) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.CreateInBatches(payroll, 100).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *repository) FindAll(filter *PayrollFilter) ([]Payroll, int64, error) {
	var payrolls []Payroll
	var total int64

	query := r.db.Model(&Payroll{}).
		Joins("JOIN employees ON employees.id = payrolls.employee_id").
		Preload("Employee")

	if filter.Month > 0 && filter.Year > 0 {
		startDate := time.Date(filter.Year, time.Month(filter.Month), 1, 0, 0, 0, 0, time.Local)
		endDate := startDate.AddDate(0, 1, -1)
		query = query.Where("period_date BETWEEN ? AND ? ", startDate, endDate)
	}

	if filter.Keyword != "" {
		keywordParam := "%" + filter.Keyword + "%"
		query = query.Where("LOWER(employees.full_name) LIKE LOWER(?) OR LOWER(employees.nik) LIKE LOWER(?)", keywordParam, keywordParam)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (filter.Page - 1) * filter.Limit
	err := query.
		Order("payrolls.created_at DESC").
		Limit(filter.Limit).
		Offset(offset).
		Find(&payrolls).Error

	return payrolls, total, err
}

func (r *repository) FindByID(id uint) (*Payroll, error) {
	var payroll Payroll
	err := r.db.
		Preload("Employee").
		Preload("Details").
		First(&payroll, id).Error
	if err != nil {
		return nil, err
	}

	return &payroll, nil
}

func (r *repository) GetExistingEmployeeID(month, year int) (map[uint]bool, error) {
	var existingID []uint

	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.Local)
	endDate := startDate.AddDate(0, 1, -1)

	err := r.db.Model(&Payroll{}).
		Where("period_date BETWEEN ? AND ?", startDate, endDate).
		Pluck("employee_id", &existingID).Error
	if err != nil {
		return nil, err
	}

	existingMap := make(map[uint]bool)
	for _, id := range existingID {
		existingMap[id] = true
	}

	return existingMap, nil
}

func (r *repository) UpdateStatus(id uint, status constants.PayrollStatus) error {
	return r.db.Model(&Payroll{}).
		Where("id = ?", id).
		Update("status", status).Error
}
