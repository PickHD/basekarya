package leave

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"hris-backend/internal/modules/attendance"
	"hris-backend/internal/modules/master"
	"hris-backend/internal/modules/user"
	"hris-backend/pkg/constants"
	"hris-backend/pkg/response"
	"hris-backend/pkg/utils"
	"time"

	"gorm.io/gorm"
)

type Service interface {
	Apply(ctx context.Context, req *ApplyRequest) error
	RequestAction(ctx context.Context, req *LeaveActionRequest) error
	GetList(ctx context.Context, filter *LeaveFilter) ([]LeaveRequestListResponse, *response.Meta, error)
	GetDetail(ctx context.Context, id uint) (*LeaveRequestDetailResponse, error)
	GenerateInitialBalance(ctx context.Context, txObj interface{}, employeeID uint) (*gorm.DB, error)
	GenerateAnnualBalance(ctx context.Context) error
}

type service struct {
	repo         Repository
	storage      StorageProvider
	notification NotificationProvider
	user         UserProvider
}

func NewService(repo Repository, storage StorageProvider, notification NotificationProvider, user UserProvider) Service {
	return &service{repo, storage, notification, user}
}

func (s *service) Apply(ctx context.Context, req *ApplyRequest) error {
	// parse date
	start, err := time.Parse(constants.DefaultTimeFormat, req.StartDate)
	if err != nil {
		return errors.New("invalid start date format")
	}
	end, err := time.Parse(constants.DefaultTimeFormat, req.EndDate)
	if err != nil {
		return errors.New("invalid end date format")
	}

	if end.Before(start) {
		return errors.New("end date must be after start date")
	}

	// calculate total days
	totalDays := int(end.Sub(start).Hours()/24) + 1

	// check balance still available or not
	balance, err := s.repo.GetBalance(req.EmployeeID, req.LeaveTypeID, start.Year())
	if err != nil {
		return err
	}

	if balance.QuotaLeft < totalDays {
		return errors.New("insufficient leave balance")
	}

	attachmentUrl := ""
	if req.AttachmentBase64 != "" {
		// construct attachment_base64 if not empty
		imgBytes, err := utils.DecodeBase64Image(req.AttachmentBase64)
		if err != nil {
			return errors.New("invalid image")
		}

		imageReader := bytes.NewReader(imgBytes)
		now := time.Now()
		todayString := now.Format(constants.DefaultTimeFormat)
		fileName := fmt.Sprintf("leaves/%d/%s-%d.jpg", req.EmployeeID, todayString, now.Unix())

		attachmentUrl, err = s.storage.UploadFileByte(ctx, fmt.Sprintf("in-%s", fileName), imageReader, int64(len(imgBytes)), "image/jpg")
		if err != nil {
			return err
		}
	}

	// construct leave request and save it to db
	leaveReq := &LeaveRequest{
		UserID:        req.UserID,
		EmployeeID:    req.EmployeeID,
		LeaveTypeID:   req.LeaveTypeID,
		StartDate:     start,
		EndDate:       end,
		TotalDays:     totalDays,
		Reason:        req.Reason,
		AttachmentURL: attachmentUrl,
		Status:        constants.LeaveStatusPending,
	}

	err = s.repo.CreateRequest(leaveReq)
	if err != nil {
		return err
	}

	adminID, err := s.user.FindAdminID()
	if err != nil {
		return err
	}

	// send notification to admin
	go func() {
		_ = s.notification.SendNotification(
			adminID,
			string(constants.NotificationTypeLeaveApprovalReq),
			"Pengajuan Cuti Baru",
			fmt.Sprintf("Karyawan mengajukan cuti pada tanggal %s s.d %s", req.StartDate, req.EndDate),
			leaveReq.ID,
		)
	}()

	return nil
}

func (s *service) RequestAction(ctx context.Context, req *LeaveActionRequest) error {
	leaveRequest, err := s.repo.FindRequestByID(req.RequestID)
	if err != nil {
		return err
	}

	if leaveRequest.Status != constants.LeaveStatusPending {
		return errors.New("request is not pending")
	}

	var (
		notificationType    constants.NotificationType
		notificationTitle   string
		notificationMessage string
	)
	switch constants.LeaveAction(req.Action) {
	case constants.LeaveActionApprove:
		shouldDeduct := leaveRequest.LeaveType.IsDeducted
		if shouldDeduct {
			balance, err := s.repo.GetBalance(leaveRequest.EmployeeID, leaveRequest.LeaveTypeID, leaveRequest.StartDate.Year())
			if err != nil {
				return errors.New("balance record not found for this employee/year")
			}
			if balance.QuotaLeft < leaveRequest.TotalDays {
				return errors.New("insufficient leave balance quota")
			}
		}

		var attendanceRecords []attendance.Attendance

		currentDate := leaveRequest.StartDate
		for !currentDate.After(leaveRequest.EndDate) {

			if currentDate.Weekday() == time.Saturday || currentDate.Weekday() == time.Sunday {
				currentDate = currentDate.AddDate(0, 0, 1)
				continue
			}

			status := constants.AttendanceStatusExcused
			if leaveRequest.LeaveType.Name == "Sick" {
				status = constants.AttendanceStatusSick
			}

			attendance := attendance.Attendance{
				EmployeeID:         leaveRequest.EmployeeID,
				ShiftID:            leaveRequest.Employee.ShiftID,
				Date:               currentDate,
				CheckInTime:        currentDate,
				CheckInLat:         0,
				CheckInLong:        0,
				CheckInAddress:     "SYSTEM_GENERATED",
				CheckInImageURL:    "",
				Status:             string(status),
				Notes:              "",
				LateDurationMinute: 0,
				IsSuspicious:       false,
			}
			attendanceRecords = append(attendanceRecords, attendance)

			currentDate = currentDate.AddDate(0, 0, 1)
		}

		err := s.repo.ApproveRequest(req.RequestID, req.ApproverID, attendanceRecords, shouldDeduct, leaveRequest.TotalDays)
		if err != nil {
			return err
		}

		notificationType = constants.NotificationTypeApproved
		notificationTitle = "Permintaan Disetujui"
		notificationMessage = "Cuti Anda telah disetujui oleh Admin."
	case constants.LeaveActionReject:
		if req.RejectionReason == "" {
			return errors.New("rejection reason required")
		}

		err := s.repo.RejectRequest(req.RequestID, req.ApproverID, req.RejectionReason)
		if err != nil {
			return err
		}

		notificationType = constants.NotificationTypeRejected
		notificationTitle = "Permintaan Ditolak"
		notificationMessage = "Cuti Anda telah ditolak oleh Admin."
	default:
		return errors.New("invalid action")
	}

	// send notification to requester
	go func() {
		_ = s.notification.SendNotification(
			leaveRequest.User.ID,
			string(notificationType),
			notificationTitle,
			notificationMessage,
			leaveRequest.ID,
		)
	}()

	return nil
}

func (s *service) GetList(ctx context.Context, filter *LeaveFilter) ([]LeaveRequestListResponse, *response.Meta, error) {
	requests, total, err := s.repo.FindAllRequests(filter)
	if err != nil {
		return []LeaveRequestListResponse{}, nil, nil
	}

	if len(requests) == 0 {
		return []LeaveRequestListResponse{}, nil, nil
	}

	var list []LeaveRequestListResponse
	for _, req := range requests {
		empName := "-"
		empNik := "-"
		leaveTypeResp := &master.LookupLeaveTypeResponse{}

		if req.Employee != nil {
			empName = req.Employee.FullName
			empNik = req.Employee.NIK
		}

		if req.LeaveType != nil {
			leaveTypeResp.ID = req.LeaveTypeID
			leaveTypeResp.Name = req.LeaveType.Name
			leaveTypeResp.DefaultQuota = req.LeaveType.DefaultQuota
			leaveTypeResp.IsDeducted = req.LeaveType.IsDeducted
		}

		list = append(list, LeaveRequestListResponse{
			ID:           req.ID,
			StartDate:    req.StartDate,
			EndDate:      req.EndDate,
			Status:       req.Status,
			EmployeeID:   req.EmployeeID,
			EmployeeName: empName,
			EmployeeNIK:  empNik,
			TotalDays:    req.TotalDays,
			LeaveTypeID:  req.LeaveTypeID,
			LeaveType:    leaveTypeResp,
			CreatedAt:    req.CreatedAt,
		})

	}

	meta := response.NewMetaOffset(filter.Page, filter.Limit, total)
	return list, meta, nil
}

func (s *service) GetDetail(ctx context.Context, id uint) (*LeaveRequestDetailResponse, error) {
	empName := "-"
	empNik := "-"
	leaveTypeResp := &master.LookupLeaveTypeResponse{}

	detail, err := s.repo.FindRequestByID(id)
	if err != nil {
		return nil, err
	}

	if detail.Employee != nil {
		empName = detail.Employee.FullName
		empNik = detail.Employee.NIK
	}

	if detail.LeaveType != nil {
		leaveTypeResp.ID = detail.LeaveTypeID
		leaveTypeResp.Name = detail.LeaveType.Name
		leaveTypeResp.DefaultQuota = detail.LeaveType.DefaultQuota
		leaveTypeResp.IsDeducted = detail.LeaveType.IsDeducted
	}

	return &LeaveRequestDetailResponse{
		ID:              id,
		StartDate:       detail.StartDate,
		EndDate:         detail.EndDate,
		Status:          detail.Status,
		EmployeeID:      detail.EmployeeID,
		EmployeeName:    empName,
		EmployeeNIK:     empNik,
		LeaveTypeID:     detail.LeaveTypeID,
		LeaveType:       leaveTypeResp,
		TotalDays:       detail.TotalDays,
		Reason:          detail.Reason,
		AttachmentURL:   detail.AttachmentURL,
		RejectionReason: detail.RejectionReason,
		CreatedAt:       detail.CreatedAt,
	}, nil
}

func (s *service) GenerateInitialBalance(ctx context.Context, txObj interface{}, employeeID uint) (*gorm.DB, error) {
	tx, ok := txObj.(*gorm.DB)
	if !ok {
		return nil, errors.New("invalid transaction type")
	}

	var leaveTypes []master.LeaveType
	if err := tx.Find(&leaveTypes).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	currentYear := time.Now().Year()
	var balances []LeaveBalance

	for _, lt := range leaveTypes {
		quota := lt.DefaultQuota

		if lt.Name == "Annual" {
			currentMonth := int(time.Now().Month())
			remainingMonths := 12 - currentMonth + 1
			quota = (lt.DefaultQuota * remainingMonths) / 12
		}

		balances = append(balances, LeaveBalance{
			EmployeeID:  employeeID,
			LeaveTypeID: lt.ID,
			Year:        currentYear,
			QuotaTotal:  quota,
			QuotaUsed:   0,
			QuotaLeft:   quota,
		})
	}

	if len(balances) > 0 {
		if err := tx.Create(&balances).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	return tx, nil
}

func (s *service) GenerateAnnualBalance(ctx context.Context) error {
	tx := s.repo.StartTX()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	currentYear := time.Now().Year()

	var employees []user.Employee
	if err := tx.Find(&employees).Error; err != nil {
		tx.Rollback()
		return err
	}

	var leaveTypes []master.LeaveType
	if err := tx.Find(&leaveTypes).Error; err != nil {
		tx.Rollback()
		return err
	}

	var newBalances []LeaveBalance

	for _, emp := range employees {
		for _, lt := range leaveTypes {
			var count int64
			tx.Model(&LeaveBalance{}).
				Where("employee_id = ? AND leave_type_id = ? AND year = ?", emp.ID, lt.ID, currentYear).
				Count(&count)

			if count == 0 {
				newBalances = append(newBalances, LeaveBalance{
					EmployeeID:  emp.ID,
					LeaveTypeID: lt.ID,
					Year:        currentYear,
					QuotaTotal:  lt.DefaultQuota,
					QuotaUsed:   0,
					QuotaLeft:   lt.DefaultQuota,
				})
			}
		}
	}

	if len(newBalances) > 0 {
		if err := tx.CreateInBatches(newBalances, 100).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}
