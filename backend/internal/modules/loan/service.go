package loan

import (
	"basekarya-backend/internal/infrastructure"
	"basekarya-backend/pkg/constants"
	"basekarya-backend/pkg/response"
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type Service interface {
	Create(ctx context.Context, req *LoanRequest) error
	GetLoanDetail(ctx context.Context, id uint) (*LoanDetailResponse, error)
	GetLoans(ctx context.Context, filter LoanFilter) ([]LoanListResponse, *response.Meta, error)
	ProcessAction(ctx context.Context, req *ActionRequest) error
	Export(ctx context.Context, filter LoanFilter) ([]byte, error)
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

func (s *service) Create(ctx context.Context, req *LoanRequest) error {
	return s.transactionManager.RunInTransaction(ctx, func(ctx context.Context) error {
		if req.UserID == 0 && req.EmployeeID == 0 {
			return fmt.Errorf("user not found")
		}

		// check if users still have active loan or not
		exist, err := s.repo.FindActiveLoanByUserID(ctx, req.UserID)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		if exist != nil {
			return fmt.Errorf("users still have loan")
		}

		// validate maximum loan request amount
		if req.TotalAmount > constants.LoanMaximumTotalAmount {
			return fmt.Errorf("cannot exceed maximum loan request")
		}

		loan := &Loan{
			UserID:            req.UserID,
			EmployeeID:        req.EmployeeID,
			TotalAmount:       req.TotalAmount,
			InstallmentAmount: req.InstallmentAmount,
			RemainingAmount:   req.TotalAmount,
			Status:            constants.LoanStatusPending,
		}

		err = s.repo.Create(ctx, loan)
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
				string(constants.NotificationTypeLoanApprovalReq),
				"Pengajuan Kasbon Baru",
				fmt.Sprintf("Karyawan mengajukan kasbon dengan total Rp.%2.f", req.TotalAmount),
				loan.ID,
			)
		}()

		return nil
	})
}

func (s *service) GetLoanDetail(ctx context.Context, id uint) (*LoanDetailResponse, error) {
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

	return &LoanDetailResponse{
		ID:                detail.ID,
		EmployeeID:        detail.EmployeeID,
		EmployeeName:      detail.Employee.FullName,
		EmployeeNIK:       detail.Employee.NIK,
		TotalAmount:       detail.TotalAmount,
		InstallmentAmount: detail.InstallmentAmount,
		RemainingAmount:   detail.RemainingAmount,
		Reason:            detail.Reason,
		Status:            detail.Status,
		RejectionReason:   rejectionReason,
		CreatedAt:         detail.CreatedAt,
	}, nil
}

func (s *service) GetLoans(ctx context.Context, filter LoanFilter) ([]LoanListResponse, *response.Meta, error) {
	loans, total, err := s.repo.FindAll(ctx, filter)
	if err != nil {
		return []LoanListResponse{}, nil, nil
	}

	if len(loans) == 0 {
		return []LoanListResponse{}, nil, nil
	}

	var list []LoanListResponse
	for _, loan := range loans {
		list = append(list, LoanListResponse{
			ID:                loan.ID,
			EmployeeID:        loan.EmployeeID,
			EmployeeName:      loan.Employee.FullName,
			EmployeeNIK:       loan.Employee.NIK,
			TotalAmount:       loan.TotalAmount,
			InstallmentAmount: loan.InstallmentAmount,
			RemainingAmount:   loan.RemainingAmount,
			Status:            loan.Status,
			CreatedAt:         loan.CreatedAt,
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

		// check data only loan with status PENDING can do further process action
		if data.Status != constants.LoanStatusPending {
			return fmt.Errorf("cannot process loan with status %s", data.Status)
		}

		var (
			notificationType    constants.NotificationType
			notificationTitle   string
			notificationMessage string
		)
		switch constants.LoanAction(req.Action) {
		case constants.LoanActionApprove:
			data.Status = constants.LoanStatusApproved
			data.ApprovedBy = &req.SuperAdminID

			notificationType = constants.NotificationTypeApproved
			notificationTitle = "Permintaan Disetujui"
			notificationMessage = "Kasbon Anda telah disetujui oleh Admin."
		case constants.LoanActionReject:
			data.Status = constants.LoanStatusRejected

			// if rejected, reason become required
			if req.RejectionReason == "" {
				return fmt.Errorf("rejection reason is required")
			}

			data.RejectionReason.String = req.RejectionReason
			data.RejectionReason.Valid = true

			notificationType = constants.NotificationTypeRejected
			notificationTitle = "Permintaan Ditolak"
			notificationMessage = "Kasbon Anda telah ditolak oleh Admin."
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
				string(notificationType),
				notificationTitle,
				notificationMessage,
				data.ID,
			)
		}()

		return nil
	})
}

func (s *service) Export(ctx context.Context, filter LoanFilter) ([]byte, error) {
	filter.Page = 1
	filter.Limit = 999999

	loans, _, err := s.repo.FindAll(ctx, filter)
	if err != nil {
		return nil, err
	}

	headers := []string{
		"ID", "Karyawan", "Total", "Cicilan / bln", "Sisa", "Status", "Tanggal Pengajuan",
	}

	var rows [][]interface{}
	for _, loan := range loans {
		empName := "-"
		if loan.Employee.FullName != "" {
			empName = loan.Employee.FullName
		}

		row := []interface{}{
			loan.ID,
			empName,
			loan.TotalAmount,
			loan.InstallmentAmount,
			loan.RemainingAmount,
			loan.Status,
			loan.CreatedAt.Format("2006-01-02 15:04:05"),
		}
		rows = append(rows, row)
	}

	return s.excel.GenerateSimpleExcel("Loans", headers, rows)
}
