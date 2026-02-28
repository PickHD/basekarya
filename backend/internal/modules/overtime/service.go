package overtime

import (
	"basekarya-backend/internal/infrastructure"
	"basekarya-backend/pkg/constants"
	"basekarya-backend/pkg/response"
	"context"
	"fmt"
	"time"
)

type Service interface {
	Create(ctx context.Context, req *OvertimeRequest) error
	GetDetail(ctx context.Context, id uint) (*OvertimeDetailResponse, error)
	GetList(ctx context.Context, filter OvertimeFilter) ([]OvertimeListResponse, *response.Meta, error)
	ProcessAction(ctx context.Context, req *ActionRequest) error
	Export(ctx context.Context, filter OvertimeFilter) ([]byte, error)
}

type service struct {
	repo               Repository
	notification       NotificationProvider
	user               UserProvider
	transactionManager infrastructure.TransactionManager
	excel              infrastructure.ExcelProvider
}

func NewService(repo Repository, notification NotificationProvider, user UserProvider, transactionManager infrastructure.TransactionManager, excel infrastructure.ExcelProvider) Service {
	return &service{repo, notification, user, transactionManager, excel}
}

func (s *service) Create(ctx context.Context, req *OvertimeRequest) error {
	return s.transactionManager.RunInTransaction(ctx, func(ctx context.Context) error {
		if req.UserID == 0 && req.EmployeeID == 0 {
			return fmt.Errorf("user not found")
		}

		// Calculate duration
		start, err := time.Parse(constants.ShiftHourFormat, req.StartTime)
		if err != nil {
			return fmt.Errorf("invalid start time format %s", constants.ShiftHourFormat)
		}
		end, err := time.Parse(constants.ShiftHourFormat, req.EndTime)
		if err != nil {
			return fmt.Errorf("invalid end time format %s", constants.ShiftHourFormat)
		}

		duration := end.Sub(start)
		if duration < 0 {
			// Crossed midnight
			duration += 24 * time.Hour
		}

		durationMinutes := int(duration.Minutes())

		if durationMinutes <= 0 {
			return fmt.Errorf("duration must be greater than 0")
		}

		overtime := &Overtime{
			UserID:          req.UserID,
			EmployeeID:      req.EmployeeID,
			Date:            req.Date,
			StartTime:       req.StartTime,
			EndTime:         req.EndTime,
			DurationMinutes: durationMinutes,
			Reason:          req.Reason,
			Status:          constants.OvertimeStatusPending,
		}

		err = s.repo.Create(ctx, overtime)
		if err != nil {
			return err
		}

		adminID, err := s.user.FindAdminID(ctx)
		if err != nil {
			return err
		}

		go func() {
			_ = s.notification.SendNotification(
				adminID,
				string(constants.NotificationTypeOvertimeApprovalReq),
				"Pengajuan Lembur Baru",
				fmt.Sprintf("Karyawan mengajukan lembur selama %d menit", durationMinutes),
				overtime.ID,
			)
		}()

		return nil
	})
}

func (s *service) GetDetail(ctx context.Context, id uint) (*OvertimeDetailResponse, error) {
	detail, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if detail.User.ID == 0 && detail.Employee.ID == 0 {
		return nil, fmt.Errorf("data user not found")
	}

	rejectionReason := ""
	if detail.RejectionReason.Valid {
		rejectionReason = detail.RejectionReason.String
	}

	return &OvertimeDetailResponse{
		ID:              detail.ID,
		EmployeeID:      detail.EmployeeID,
		EmployeeName:    detail.Employee.FullName,
		EmployeeNIK:     detail.Employee.NIK,
		Date:            detail.Date,
		StartTime:       detail.StartTime,
		EndTime:         detail.EndTime,
		DurationMinutes: detail.DurationMinutes,
		Reason:          detail.Reason,
		Status:          detail.Status,
		RejectionReason: rejectionReason,
		CreatedAt:       detail.CreatedAt,
	}, nil
}

func (s *service) GetList(ctx context.Context, filter OvertimeFilter) ([]OvertimeListResponse, *response.Meta, error) {
	overtimes, total, err := s.repo.FindAll(ctx, filter)
	if err != nil {
		return []OvertimeListResponse{}, nil, nil
	}

	if len(overtimes) == 0 {
		return []OvertimeListResponse{}, nil, nil
	}

	var list []OvertimeListResponse
	for _, overtime := range overtimes {
		list = append(list, OvertimeListResponse{
			ID:              overtime.ID,
			EmployeeID:      overtime.EmployeeID,
			EmployeeName:    overtime.Employee.FullName,
			EmployeeNIK:     overtime.Employee.NIK,
			Date:            overtime.Date,
			StartTime:       overtime.StartTime,
			EndTime:         overtime.EndTime,
			DurationMinutes: overtime.DurationMinutes,
			Status:          overtime.Status,
			CreatedAt:       overtime.CreatedAt,
		})
	}

	meta := response.NewMetaOffset(filter.Page, filter.Limit, total)
	return list, meta, nil
}

func (s *service) ProcessAction(ctx context.Context, req *ActionRequest) error {
	return s.transactionManager.RunInTransaction(ctx, func(ctx context.Context) error {
		data, err := s.repo.FindByID(ctx, req.ID)
		if err != nil {
			return err
		}

		if data.Status != constants.OvertimeStatusPending {
			return fmt.Errorf("cannot process overtime with status %s", data.Status)
		}

		var (
			notificationType    string
			notificationTitle   string
			notificationMessage string
		)
		switch constants.OvertimeAction(req.Action) {
		case constants.OvertimeActionApprove:
			data.Status = constants.OvertimeStatusApproved
			data.ApprovedBy = &req.SuperAdminID

			notificationType = string(constants.NotificationTypeApproved)
			notificationTitle = "Lembur Disetujui"
			notificationMessage = "Lembur Anda telah disetujui oleh Admin."
		case constants.OvertimeActionReject:
			data.Status = constants.OvertimeStatusRejected

			if req.RejectionReason == "" {
				return fmt.Errorf("rejection reason is required")
			}

			data.RejectionReason.String = req.RejectionReason
			data.RejectionReason.Valid = true

			notificationType = string(constants.NotificationTypeRejected)
			notificationTitle = "Lembur Ditolak"
			notificationMessage = "Lembur Anda telah ditolak oleh Admin."
		default:
			return fmt.Errorf("invalid action: %s", req.Action)
		}

		err = s.repo.Update(ctx, data)
		if err != nil {
			return err
		}

		go func() {
			_ = s.notification.SendNotification(
				data.UserID,
				notificationType,
				notificationTitle,
				notificationMessage,
				data.ID,
			)
		}()

		return nil
	})
}

func (s *service) Export(ctx context.Context, filter OvertimeFilter) ([]byte, error) {
	filter.Page = 1
	filter.Limit = 999999

	overtimes, _, err := s.repo.FindAll(ctx, filter)
	if err != nil {
		return nil, err
	}

	headers := []string{
		"ID", "Karyawan", "Tanggal", "Mulai", "Selesai", "Durasi (Menit)", "Alasan", "Status", "Dibuat Pada",
	}

	var rows [][]interface{}
	for _, ot := range overtimes {
		empName := "-"
		if ot.Employee.FullName != "" {
			empName = ot.Employee.FullName
		}

		row := []interface{}{
			ot.ID,
			empName,
			ot.Date,
			ot.StartTime,
			ot.EndTime,
			ot.DurationMinutes,
			ot.Reason,
			ot.Status,
			ot.CreatedAt.Format("2006-01-02 15:04:05"),
		}
		rows = append(rows, row)
	}

	return s.excel.GenerateSimpleExcel("Overtimes", headers, rows)
}
