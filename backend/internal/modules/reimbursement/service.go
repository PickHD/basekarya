package reimbursement

import (
	"context"
	"fmt"
	"hris-backend/pkg/constants"
	"hris-backend/pkg/response"
	"path/filepath"
	"time"

	"github.com/google/uuid"
)

type Service interface {
	Create(ctx context.Context, req *ReimbursementRequest) error
	ProcessAction(ctx context.Context, req *ActionRequest) error
	GetReimburseDetail(ctx context.Context, id uint) (*ReimbursementDetailResponse, error)
	GetReimbursements(ctx context.Context, filter ReimbursementFilter) ([]ReimbursementListResponse, *response.Meta, error)
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

func (s *service) Create(ctx context.Context, req *ReimbursementRequest) error {
	if req.UserID == 0 {
		return fmt.Errorf("user id is invalid")
	}

	ext := filepath.Ext(req.File.Filename)
	newFileName := fmt.Sprintf("%s%s", uuid.New().String(), ext)

	now := time.Now()
	objectName := fmt.Sprintf("reimbursements/%d/%02d/%s", now.Year(), now.Month(), newFileName)

	fileURL, err := s.storage.UploadFileMultipart(ctx, req.File, objectName)
	if err != nil {
		return fmt.Errorf("failed to upload proof: %w", err)
	}

	dateExpense, err := time.Parse(constants.DefaultTimeFormat, req.Date)
	if err != nil {
		return fmt.Errorf("invalid date format: %w", err)
	}

	reimburstment := &Reimbursement{
		UserID:        req.UserID,
		Title:         req.Title,
		Description:   req.Description,
		Amount:        req.Amount,
		DateOfExpense: dateExpense,
		ProofFileURL:  fileURL,
		Status:        constants.ReimbursementStatusPending,
	}

	err = s.repo.Create(reimburstment)
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
			string(constants.NotificationTypeReimburseApprovalReq),
			"Pengajuan Reimbursement Baru",
			fmt.Sprintf("Karyawan mengajukan reimburse pada tanggal %s", req.Date),
			reimburstment.ID,
		)
	}()

	return nil
}

func (s *service) GetReimburseDetail(ctx context.Context, id uint) (*ReimbursementDetailResponse, error) {
	detail, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if detail.User.ID == 0 {
		return nil, fmt.Errorf("data user not found")
	}

	rejectionReason := ""
	if detail.RejectionReason.Valid {
		rejectionReason = detail.RejectionReason.String
	}

	return &ReimbursementDetailResponse{
		ID:              detail.ID,
		Title:           detail.Title,
		Description:     detail.Description,
		Amount:          detail.Amount,
		DateOfExpense:   detail.DateOfExpense,
		ProofFileURL:    detail.ProofFileURL,
		Status:          string(detail.Status),
		RejectionReason: &rejectionReason,
		RequesterName:   detail.User.Username,
	}, nil
}

func (s *service) GetReimbursements(ctx context.Context, filter ReimbursementFilter) ([]ReimbursementListResponse, *response.Meta, error) {
	reimbursements, total, err := s.repo.FindAll(filter)
	if err != nil {
		return []ReimbursementListResponse{}, nil, nil
	}

	if len(reimbursements) == 0 {
		return []ReimbursementListResponse{}, nil, nil
	}

	var list []ReimbursementListResponse
	for _, rem := range reimbursements {
		list = append(list, ReimbursementListResponse{
			ID:            rem.ID,
			Title:         rem.Title,
			Amount:        rem.Amount,
			DateOfExpense: rem.DateOfExpense,
			ProofFileURL:  rem.ProofFileURL,
			Status:        string(rem.Status),
		})
	}

	meta := response.NewMetaOffset(filter.Page, filter.Limit, total)
	return list, meta, nil
}

func (s *service) ProcessAction(ctx context.Context, req *ActionRequest) error {
	data, err := s.repo.FindByID(req.ID)
	if err != nil {
		return err
	}

	// check data only reimburstment with status PENDING can do further process action
	if data.Status != constants.ReimbursementStatusPending {
		return fmt.Errorf("cannot process reimburstment with status %s", data.Status)
	}

	var (
		notificationType    constants.NotificationType
		notificationTitle   string
		notificationMessage string
	)
	switch constants.ReimbursementAction(req.Action) {
	case constants.ReimbursementActionApprove:
		data.Status = constants.ReimbursementStatusApproved
		data.ApprovedBy = &req.SuperAdminID

		notificationType = constants.NotificationTypeApproved
		notificationTitle = "Permintaan Disetujui"
		notificationMessage = "Reimburse Anda telah disetujui oleh Admin."
	case constants.ReimbursementActionReject:
		data.Status = constants.ReimbursementStatusRejected

		// if rejected, reason become required
		if req.RejectionReason == "" {
			return fmt.Errorf("rejection reason is required")
		}

		data.RejectionReason.String = req.RejectionReason
		data.RejectionReason.Valid = true

		notificationType = constants.NotificationTypeRejected
		notificationTitle = "Permintaan Ditolak"
		notificationMessage = "Reimburse Anda telah ditolak oleh Admin."
	default:
		return fmt.Errorf("invalid action: %s", req.Action)
	}

	err = s.repo.Update(data)
	if err != nil {
		return err
	}

	// send notification to requester
	go func() {
		_ = s.notification.SendNotification(
			data.User.ID,
			string(notificationType),
			notificationTitle,
			notificationMessage,
			data.ID,
		)
	}()

	return nil
}
