package leave

import (
	"context"
	"errors"
	"hris-backend/internal/modules/attendance"
	"hris-backend/internal/modules/master"
	"hris-backend/pkg/constants"
	"hris-backend/pkg/utils"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository interface {
	CreateRequest(ctx context.Context, req *LeaveRequest) error
	FindRequestByID(ctx context.Context, id uint) (*LeaveRequest, error)
	FindAllRequests(ctx context.Context, filter *LeaveFilter) ([]LeaveRequest, int64, error)

	GetBalance(ctx context.Context, employeeID, leaveTypeID uint, year int) (*LeaveBalance, error)

	ApproveRequest(ctx context.Context, requestID uint, approverID uint, attendanceRecords []attendance.Attendance, shouldDeduct bool, days int) error
	RejectRequest(ctx context.Context, requestID uint, approverID uint, reason string) error

	// For Initial Balance Generation
	FindAllLeaveTypes(ctx context.Context) ([]master.LeaveType, error)
	CreateLeaveBalances(ctx context.Context, balances []LeaveBalance) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) CreateRequest(ctx context.Context, req *LeaveRequest) error {
	db := utils.GetDBFromContext(ctx, r.db)
	return db.Create(req).Error
}

func (r *repository) FindRequestByID(ctx context.Context, id uint) (*LeaveRequest, error) {
	db := utils.GetDBFromContext(ctx, r.db)
	var req LeaveRequest
	err := db.
		Preload("User").
		Preload("Employee").
		Preload("LeaveType").
		First(&req, id).Error

	return &req, err
}

func (r *repository) FindAllRequests(ctx context.Context, filter *LeaveFilter) ([]LeaveRequest, int64, error) {
	db := utils.GetDBFromContext(ctx, r.db)
	var requests []LeaveRequest
	var total int64

	offset := (filter.Page - 1) * filter.Limit

	query := db.Model(&LeaveRequest{}).
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

func (r *repository) GetBalance(ctx context.Context, employeeID, leaveTypeID uint, year int) (*LeaveBalance, error) {
	db := utils.GetDBFromContext(ctx, r.db)
	var balance LeaveBalance
	err := db.
		Where("employee_id = ? AND leave_type_id = ? AND year = ?",
			employeeID, leaveTypeID, year).First(&balance).Error

	return &balance, err
}

func (r *repository) ApproveRequest(ctx context.Context, requestID uint, approverID uint, attendanceRecords []attendance.Attendance, shouldDeduct bool, days int) error {
	db := utils.GetDBFromContext(ctx, r.db)

	if err := db.Model(&LeaveRequest{}).Where("id = ?", requestID).
		Updates(map[string]interface{}{
			"status":      constants.LeaveStatusApproved,
			"approved_by": approverID,
		}).Error; err != nil {
		return err
	}

	if len(attendanceRecords) > 0 {
		if err := db.Clauses(clause.OnConflict{
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
		if err := db.First(&req, requestID).Error; err != nil {
			return err
		}

		if err := db.Model(&LeaveBalance{}).
			Where("employee_id = ? AND leave_type_id = ? AND year = ?", req.EmployeeID, req.LeaveTypeID, time.Now().Year()).
			Updates(map[string]interface{}{
				"quota_used": gorm.Expr("quota_used + ?", days),
				"quota_left": gorm.Expr("quota_left - ?", days),
			}).Error; err != nil {
			return err
		}
	}

	return nil
}

func (r *repository) RejectRequest(ctx context.Context, requestID uint, approverID uint, reason string) error {
	db := utils.GetDBFromContext(ctx, r.db)
	if reason == "" {
		return errors.New("reject reason required")
	}

	return db.Model(&LeaveRequest{}).Where("id = ?", requestID).
		Updates(map[string]interface{}{
			"status":           constants.LeaveStatusRejected,
			"approved_by":      approverID,
			"rejection_reason": reason,
		}).Error
}

func (r *repository) FindAllLeaveTypes(ctx context.Context) ([]master.LeaveType, error) {
	db := utils.GetDBFromContext(ctx, r.db)
	var leaveTypes []master.LeaveType
	if err := db.Find(&leaveTypes).Error; err != nil {
		return nil, err
	}
	return leaveTypes, nil
}

func (r *repository) CreateLeaveBalances(ctx context.Context, balances []LeaveBalance) error {
	db := utils.GetDBFromContext(ctx, r.db)
	if len(balances) == 0 {
		return nil
	}
	return db.Create(&balances).Error
}
