package finance

import (
	"basekarya-backend/internal/infrastructure"
	"basekarya-backend/pkg/constants"
	"basekarya-backend/pkg/response"
	"basekarya-backend/pkg/utils"
	"context"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type Service interface {
	CreateTransaction(ctx context.Context, req *CreateTransactionRequest) error
	GetTransactionDetail(ctx context.Context, id uint) (*TransactionDetailResponse, error)
	GetTransactions(ctx context.Context, filter TransactionFilter) ([]TransactionListResponse, *response.Meta, error)
	ProcessAction(ctx context.Context, req *ActionRequest) error
	ExportTransactions(ctx context.Context, filter TransactionFilter) ([]byte, error)

	CreateCategory(ctx context.Context, req *CategoryRequest) error
	GetCategories(ctx context.Context, catType string) ([]CategoryResponse, error)
	UpdateCategory(ctx context.Context, id uint, req *CategoryRequest) error
	DeleteCategory(ctx context.Context, id uint) error

	GetDashboard(ctx context.Context, startDate, endDate string) (*DashboardResponse, error)
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

func (s *service) CreateTransaction(ctx context.Context, req *CreateTransactionRequest) error {
	return s.transactionManager.RunInTransaction(ctx, func(ctx context.Context) error {
		if req.CreatedBy == 0 {
			return fmt.Errorf("user not found")
		}

		txDate, err := time.Parse("2006-01-02", req.TransactionDate)
		if err != nil {
			return fmt.Errorf("invalid transaction_date format, use YYYY-MM-DD")
		}

		_, err = s.repo.FindCategoryByID(ctx, req.FinanceCategoryID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("finance category not found")
			}
			return err
		}

		tx := &FinanceTransaction{
			CompanyID:          utils.GetCompanyIDFromCtx(ctx),
			FinanceCategoryID: req.FinanceCategoryID,
			CreatedBy:         req.CreatedBy,
			Type:              constants.FinanceType(req.Type),
			Amount:            req.Amount,
			TransactionDate:   txDate,
			Status:            constants.FinanceStatusPending,
		}

		if req.Description != "" {
			tx.Description.String = req.Description
			tx.Description.Valid = true
		}

		if req.ReferenceNumber != "" {
			tx.ReferenceNumber.String = req.ReferenceNumber
			tx.ReferenceNumber.Valid = true
		}

		err = s.repo.CreateTransaction(ctx, tx)
		if err != nil {
			return err
		}

		approvalUserIDs, err := s.user.FindApprovalUsers(ctx, string(constants.APPROVAL_FINANCE))
		if err != nil {
			return err
		}

		go func() {
			_ = s.notification.BlastNotification(
				ctx,
				approvalUserIDs,
				string(constants.NotificationTypeFinanceApprovalReq),
				"Pengajuan Transaksi Keuangan Baru",
				fmt.Sprintf("Transaksi keuangan %s sebesar Rp%.0f memerlukan persetujuan", req.Type, req.Amount),
				tx.ID,
			)
		}()

		return nil
	})
}

func (s *service) GetTransactionDetail(ctx context.Context, id uint) (*TransactionDetailResponse, error) {
	data, err := s.repo.FindTransactionByID(ctx, id)
	if err != nil {
		return nil, err
	}

	description := ""
	if data.Description.Valid {
		description = data.Description.String
	}

	referenceNumber := ""
	if data.ReferenceNumber.Valid {
		referenceNumber = data.ReferenceNumber.String
	}

	rejectionReason := ""
	if data.RejectionReason.Valid {
		rejectionReason = data.RejectionReason.String
	}

	creatorName := ""
	if data.Creator.Employee != nil {
		creatorName = data.Creator.Employee.FullName
	}

	approverName := ""
	if data.Approver != nil && data.Approver.Employee != nil {
		approverName = data.Approver.Employee.FullName
	}

	return &TransactionDetailResponse{
		ID:              data.ID,
		CreatorName:     creatorName,
		CategoryName:    data.FinanceCategory.Name,
		CategoryType:    data.FinanceCategory.Type,
		Type:            data.Type,
		Amount:          data.Amount,
		Description:     description,
		TransactionDate: data.TransactionDate,
		ReferenceNumber: referenceNumber,
		Status:          data.Status,
		RejectionReason: rejectionReason,
		ApprovedBy:      data.ApprovedBy,
		ApproverName:    approverName,
		CreatedAt:       data.CreatedAt,
	}, nil
}

func (s *service) GetTransactions(ctx context.Context, filter TransactionFilter) ([]TransactionListResponse, *response.Meta, error) {
	transactions, total, err := s.repo.FindAllTransactions(ctx, filter)
	if err != nil {
		return []TransactionListResponse{}, nil, nil
	}

	if len(transactions) == 0 {
		return []TransactionListResponse{}, nil, nil
	}

	var list []TransactionListResponse
	for _, tx := range transactions {
		referenceNumber := ""
		if tx.ReferenceNumber.Valid {
			referenceNumber = tx.ReferenceNumber.String
		}

		creatorName := ""
		if tx.Creator.Employee != nil {
			creatorName = tx.Creator.Employee.FullName
		}

		list = append(list, TransactionListResponse{
			ID:              tx.ID,
			CreatorName:     creatorName,
			CategoryName:    tx.FinanceCategory.Name,
			Type:            tx.Type,
			Amount:          tx.Amount,
			TransactionDate: tx.TransactionDate,
			ReferenceNumber: referenceNumber,
			Status:          tx.Status,
			CreatedAt:       tx.CreatedAt,
		})
	}

	meta := response.NewMetaOffset(filter.Page, filter.Limit, total)
	return list, meta, nil
}

func (s *service) ProcessAction(ctx context.Context, req *ActionRequest) error {
	return s.transactionManager.RunInTransaction(ctx, func(ctx context.Context) error {
		data, err := s.repo.FindTransactionByID(ctx, req.ID)
		if err != nil {
			return err
		}

		if data.Status != constants.FinanceStatusPending {
			return fmt.Errorf("cannot process transaction with status %s", data.Status)
		}

		var (
			notificationType    constants.NotificationType
			notificationTitle   string
			notificationMessage string
		)

		switch constants.FinanceAction(req.Action) {
		case constants.FinanceActionApprove:
			data.Status = constants.FinanceStatusApproved
			data.ApprovedBy = &req.SuperAdminID

			notificationType = constants.NotificationTypeApproved
			notificationTitle = "Transaksi Keuangan Disetujui"
			notificationMessage = fmt.Sprintf("Transaksi keuangan %s Anda telah disetujui.", data.Type)
		case constants.FinanceActionReject:
			data.Status = constants.FinanceStatusRejected

			if req.RejectionReason == "" {
				return fmt.Errorf("rejection reason is required")
			}

			data.RejectionReason.String = req.RejectionReason
			data.RejectionReason.Valid = true

			notificationType = constants.NotificationTypeRejected
			notificationTitle = "Transaksi Keuangan Ditolak"
			notificationMessage = fmt.Sprintf("Transaksi keuangan %s Anda telah ditolak.", data.Type)
		default:
			return fmt.Errorf("invalid action: %s", req.Action)
		}

		err = s.repo.UpdateTransaction(ctx, data)
		if err != nil {
			return err
		}

		go func() {
			_ = s.notification.SendNotification(
				ctx,
				data.CreatedBy,
				string(notificationType),
				notificationTitle,
				notificationMessage,
				data.ID,
			)
		}()

		return nil
	})
}

func (s *service) ExportTransactions(ctx context.Context, filter TransactionFilter) ([]byte, error) {
	filter.Page = 1
	filter.Limit = 999999

	transactions, _, err := s.repo.FindAllTransactions(ctx, filter)
	if err != nil {
		return nil, err
	}

	headers := []string{
		"ID", "Dibuat Oleh", "Kategori", "Tipe", "Jumlah", "Tanggal Transaksi", "No Referensi", "Status", "Tanggal Dibuat",
	}

	var rows [][]interface{}
	for _, tx := range transactions {
		referenceNumber := "-"
		if tx.ReferenceNumber.Valid {
			referenceNumber = tx.ReferenceNumber.String
		}

		creatorName := "-"
		if tx.Creator.Employee != nil {
			creatorName = tx.Creator.Employee.FullName
		}

		row := []interface{}{
			tx.ID,
			creatorName,
			tx.FinanceCategory.Name,
			tx.Type,
			tx.Amount,
			tx.TransactionDate.Format("2006-01-02"),
			referenceNumber,
			tx.Status,
			tx.CreatedAt.Format("2006-01-02 15:04:05"),
		}
		rows = append(rows, row)
	}

	return s.excel.GenerateSimpleExcel("Finance Transactions", headers, rows)
}

func (s *service) CreateCategory(ctx context.Context, req *CategoryRequest) error {
	cat := &FinanceCategory{
		CompanyID: utils.GetCompanyIDFromCtx(ctx),
		Name:      req.Name,
		Type: constants.FinanceType(req.Type),
	}

	if req.Description != "" {
		cat.Description.String = req.Description
		cat.Description.Valid = true
	}

	return s.repo.CreateCategory(ctx, cat)
}

func (s *service) GetCategories(ctx context.Context, catType string) ([]CategoryResponse, error) {
	categories, err := s.repo.FindAllCategories(ctx, catType)
	if err != nil {
		return []CategoryResponse{}, err
	}

	var list []CategoryResponse
	for _, cat := range categories {
		description := ""
		if cat.Description.Valid {
			description = cat.Description.String
		}

		list = append(list, CategoryResponse{
			ID:          cat.ID,
			Name:        cat.Name,
			Type:        cat.Type,
			Description: description,
			CreatedAt:   cat.CreatedAt,
			UpdatedAt:   cat.UpdatedAt,
		})
	}

	return list, nil
}

func (s *service) UpdateCategory(ctx context.Context, id uint, req *CategoryRequest) error {
	cat, err := s.repo.FindCategoryByID(ctx, id)
	if err != nil {
		return err
	}

	cat.Name = req.Name
	cat.Type = constants.FinanceType(req.Type)

	if req.Description != "" {
		cat.Description.String = req.Description
		cat.Description.Valid = true
	} else {
		cat.Description.Valid = false
	}

	return s.repo.UpdateCategory(ctx, cat)
}

func (s *service) DeleteCategory(ctx context.Context, id uint) error {
	return s.repo.DeleteCategory(ctx, id)
}

func (s *service) GetDashboard(ctx context.Context, startDate, endDate string) (*DashboardResponse, error) {
	return s.repo.GetDashboardSummary(ctx, startDate, endDate)
}
