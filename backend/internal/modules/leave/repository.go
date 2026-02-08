package leave

import (
	"errors"
	"hris-backend/internal/modules/attendance"
	"hris-backend/pkg/constants"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository interface {
	CreateRequest(req *LeaveRequest) error
	FindRequestByID(id uint) (*LeaveRequest, error)
	FindAllRequests(filter *LeaveFilter) ([]LeaveRequest, int64, error)

	GetBalance(employeeID, leaveTypeID uint, year int) (*LeaveBalance, error)

	ApproveRequest(requestID uint, approverID uint, attendanceRecords []attendance.Attendance, shouldDeduct bool, days int) error
	RejectRequest(requestID uint, approverID uint, reason string) error
	StartTX() *gorm.DB
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) CreateRequest(req *LeaveRequest) error {
	return r.db.Create(req).Error
}

func (r *repository) FindRequestByID(id uint) (*LeaveRequest, error) {
	var req LeaveRequest
	err := r.db.
		Preload("User").
		Preload("Employee").
		Preload("LeaveType").
		First(&req, id).Error

	return &req, err
}

func (r *repository) FindAllRequests(filter *LeaveFilter) ([]LeaveRequest, int64, error) {
	var requests []LeaveRequest
	var total int64

	offset := (filter.Page - 1) * filter.Limit

	query := r.db.Model(&LeaveRequest{}).
		Joins("JOIN employees ON employees.id = leave_requests.employee_id").
		Joins("JOIN ref_leave_types ON ref_leave_types.id = leave_requests.leave_type_id").
		Preload("Employee").
		Preload("LeaveType")

	if filter.Status != "" {
		query = query.Where("leave_requests.status = ?", filter.Status)
	}
	if filter.UserID > 0 {
		query = query.Where("leave_requests.user_id = ?", filter.UserID)
	}
	if filter.Search != "" {
		searchParam := "%" + filter.Search + "%"
		query = query.Where("LOWER(employees.full_name) LIKE LOWER(?) OR LOWER(employees.nik) LIKE LOWER(?)", searchParam, searchParam)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.
		Order("leave_requests.created_at DESC").
		Limit(filter.Limit).
		Offset(offset).
		Find(&requests).Error

	return requests, total, err
}

func (r *repository) GetBalance(employeeID, leaveTypeID uint, year int) (*LeaveBalance, error) {
	var balance LeaveBalance
	err := r.db.
		Where("employee_id = ? AND leave_type_id = ? AND year = ?",
			employeeID, leaveTypeID, year).First(&balance).Error

	return &balance, err
}

func (r *repository) ApproveRequest(requestID uint, approverID uint, attendanceRecords []attendance.Attendance, shouldDeduct bool, days int) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&LeaveRequest{}).Where("id = ?", requestID).
			Updates(map[string]interface{}{
				"status":      constants.LeaveStatusApproved,
				"approved_by": approverID,
			}).Error; err != nil {
			return err
		}

		if len(attendanceRecords) > 0 {
			if err := tx.Clauses(clause.OnConflict{
				Columns: []clause.Column{{Name: "employee_id"}, {Name: "date"}},
				DoUpdates: clause.AssignmentColumns([]string{
					"shift_id",
					"check_in_time",
					"check_in_lat",
					"check_in_long",
					"check_in_address",
					"check_in_image_url",
					"status",
					"notes",
					"late_duration_minute",
					"is_suspicious",
				}),
			}).CreateInBatches(attendanceRecords, 31).Error; err != nil {
				return err
			}
		}

		if shouldDeduct {
			var req LeaveRequest
			if err := tx.First(&req, requestID).Error; err != nil {
				return err
			}

			if err := tx.Model(&LeaveBalance{}).
				Where("employee_id = ? AND leave_type_id = ? AND year = ?", req.EmployeeID, req.LeaveTypeID, time.Now().Year()).
				Updates(map[string]interface{}{
					"quota_used": gorm.Expr("quota_used + ?", days),
					"quota_left": gorm.Expr("quota_left - ?", days),
				}).Error; err != nil {
				return err
			}
		}

		return nil
	})

}

func (r *repository) RejectRequest(requestID uint, approverID uint, reason string) error {
	if reason == "" {
		return errors.New("reject reason required")
	}

	return r.db.Model(&LeaveRequest{}).Where("id = ?", requestID).
		Updates(map[string]interface{}{
			"status":           constants.LeaveStatusRejected,
			"approved_by":      approverID,
			"rejection_reason": reason,
		}).Error
}

func (r *repository) StartTX() *gorm.DB {
	return r.db.Begin()
}
