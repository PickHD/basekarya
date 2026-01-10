package attendance

import (
	"hris-backend/pkg/constants"
	"time"

	"gorm.io/gorm"
)

type Repository interface {
	GetTodayAttendance(employeeID uint) (*Attendance, error)
	Create(attendance *Attendance) error
	Update(attendance *Attendance) error
	GetHistory(employeeID uint, month, year, page, limit int) ([]Attendance, int64, error)
	FindAll(filter *FilterParams) ([]Attendance, int64, error)
	CountByStatus(status constants.AttendanceStatus, todayDate string) (int64, error)
	CountAttendanceToday(todayDate string) (int64, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) GetTodayAttendance(employeeID uint) (*Attendance, error) {
	var att Attendance

	err := r.db.Where("employee_id = ? AND date = ?", employeeID, time.Now().Format("2006-01-02")).
		First(&att).Error
	if err != nil {
		return nil, err
	}

	return &att, nil
}

func (r *repository) Create(attendance *Attendance) error {
	return r.db.Create(attendance).Error
}

func (r *repository) Update(attendance *Attendance) error {
	return r.db.Save(attendance).Error
}

func (r *repository) GetHistory(employeeID uint, month, year, page, limit int) ([]Attendance, int64, error) {
	var logs []Attendance
	var total int64

	query := r.db.Model(&Attendance{}).
		Where("employee_id = ? AND MONTH(date) = ? AND YEAR(date) = ?", employeeID, month, year)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	err := query.
		Preload("Shift").
		Order("date DESC").
		Limit(limit).
		Offset(offset).
		Find(&logs).Error

	return logs, total, err
}

func (r *repository) FindAll(filter *FilterParams) ([]Attendance, int64, error) {
	var logs []Attendance
	var total int64

	query := r.db.Model(&Attendance{}).
		Joins("JOIN employees ON employees.id = attendances.employee_id").
		Joins("JOIN ref_departments ON ref_departments.id = employees.department_id").
		Preload("Employee").
		Preload("Employee.Department").
		Preload("Shift")

	// filter range date
	if filter.StartDate != "" && filter.EndDate != "" {
		query = query.Where("attendances.date BETWEEN ? AND ?", filter.StartDate, filter.EndDate)
	}

	// filter departments
	if filter.DepartmentID > 0 {
		query = query.Where("employees.department_id = ?", filter.DepartmentID)
	}

	// filter search by full name or NIK
	if filter.Search != "" {
		searchParam := "%" + filter.Search + "%"
		query = query.Where("LOWER(employees.full_name) LIKE LOWER(?) OR LOWER(employees.nik) LIKE LOWER(?)", searchParam, searchParam)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// check if there filter limit or not (pagination)
	if filter.Limit > 0 {
		offset := (filter.Page - 1) * filter.Limit
		query = query.Limit(filter.Limit).Offset(offset)
	}

	err := query.Order("attendances.date DESC").Find(&logs).Error

	return logs, total, err
}

func (r *repository) CountByStatus(status constants.AttendanceStatus, todayDate string) (int64, error) {
	var totalStatus int64
	if err := r.db.Model(&Attendance{}).
		Where("date = ? AND status = ?", todayDate, string(status)).
		Count(&totalStatus).Error; err != nil {
		return 0, err
	}

	return totalStatus, nil
}

func (r *repository) CountAttendanceToday(todayDate string) (int64, error) {
	var totalStatus int64
	if err := r.db.Model(&Attendance{}).
		Where("date = ?", todayDate).
		Count(&totalStatus).Error; err != nil {
		return 0, err
	}

	return totalStatus, nil
}
