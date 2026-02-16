package attendance

import (
	"hris-backend/pkg/constants"
	"hris-backend/pkg/response"
	"time"

	"gorm.io/gorm"
)

type Repository interface {
	GetTodayAttendance(employeeID uint) (*Attendance, error)
	Create(attendance *Attendance) error
	Update(attendance *Attendance) error
	GetHistory(employeeID uint, month, year, limit int, cursor string) ([]Attendance, *response.Cursor, error)
	FindAll(filter *FilterParams) ([]Attendance, *response.Cursor, error)
	CountByStatus(status constants.AttendanceStatus, todayDate string) (int64, error)
	CountAttendanceToday(todayDate string) (int64, error)
	GetBulkLateDuration(month, year int) (map[uint]int, error)
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

func (r *repository) GetHistory(employeeID uint, month, year, limit int, cursor string) ([]Attendance, *response.Cursor, error) {
	var logs []Attendance

	query := r.db.Model(&Attendance{}).
		Where("employee_id = ? ", employeeID).
		Order("created_at DESC, id DESC").
		Limit(limit + 1)

	if month > 0 {
		query = query.Where("MONTH(date) = ?", month)
	}

	if year > 0 {
		query = query.Where("YEAR(date) = ?", year)
	}

	if cursor != "" {
		var decoded *response.Cursor
		err := response.DecodeCursor(cursor, &decoded)

		if err == nil && decoded != nil {
			query = query.Where(
				"(created_at < ? ) OR (created_at = ? AND id < ?)",
				decoded.SortValue, decoded.SortValue, decoded.ID,
			)
		}
	}

	if err := query.Find(&logs).Error; err != nil {
		return nil, nil, err
	}

	var nextCursor *response.Cursor
	if len(logs) > limit {
		logs = logs[:limit]
		lastItem := logs[len(logs)-1]

		nextCursor = &response.Cursor{
			ID:        lastItem.ID,
			SortValue: lastItem.CreatedAt,
		}
	}

	return logs, nextCursor, nil
}

func (r *repository) FindAll(filter *FilterParams) ([]Attendance, *response.Cursor, error) {
	var logs []Attendance

	query := r.db.Model(&Attendance{}).
		Joins("JOIN employees ON employees.id = attendances.employee_id").
		Joins("JOIN ref_departments ON ref_departments.id = employees.department_id").
		Preload("Employee").
		Preload("Employee.Department").
		Preload("Shift").
		Order("attendances.created_at DESC, attendances.id DESC").
		Limit(filter.Limit + 1)

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

	// cursor logic implementation
	if filter.Cursor != "" {
		var decoded *response.Cursor
		err := response.DecodeCursor(filter.Cursor, &decoded)

		if err == nil && decoded != nil {
			query = query.Where(
				"(attendances.created_at < ?) OR (attendances.created_at = ? AND attendances.id < ?)",
				decoded.SortValue, decoded.SortValue, decoded.ID,
			)
		}
	}

	if err := query.Find(&logs).Error; err != nil {
		return nil, nil, err
	}

	var nextCursor *response.Cursor
	if len(logs) > filter.Limit {
		logs = logs[:filter.Limit]
		lastItem := logs[len(logs)-1]

		nextCursor = &response.Cursor{
			ID:        lastItem.ID,
			SortValue: lastItem.CreatedAt,
		}
	}

	return logs, nextCursor, nil
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

func (r *repository) GetBulkLateDuration(month, year int) (map[uint]int, error) {
	type Result struct {
		UserID      uint
		TotalMinute int
	}
	var results []Result

	err := r.db.Model(&Attendance{}).
		Select("employee_id, COALESCE(SUM(late_duration_minute), 0) as total_minute").
		Where("MONTH(check_in_time) = ? AND YEAR(check_in_time) = ?", month, year).
		Group("employee_id").
		Scan(&results).Error
	if err != nil {
		return nil, err
	}

	dataMap := make(map[uint]int)
	for _, res := range results {
		dataMap[res.UserID] = res.TotalMinute
	}

	return dataMap, err
}
