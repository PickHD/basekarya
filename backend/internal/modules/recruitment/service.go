package recruitment

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"basekarya-backend/internal/infrastructure"
	"basekarya-backend/pkg/constants"
	"basekarya-backend/pkg/logger"
	"basekarya-backend/pkg/response"
	"basekarya-backend/pkg/utils"
)

type Service interface {
	// Requisition
	CreateRequisition(ctx context.Context, requesterID uint, req *CreateRequisitionRequest) error
	SubmitRequisition(ctx context.Context, id uint, requesterID uint) error
	RequisitionAction(ctx context.Context, id uint, approverID uint, req *RequisitionActionRequest) error
	GetRequisitions(ctx context.Context, filter *RequisitionFilter) ([]RequisitionListResponse, *response.Meta, error)
	GetRequisitionDetail(ctx context.Context, id uint) (*RequisitionDetailResponse, error)
	CloseRequisition(ctx context.Context, id uint) error
	DeleteRequisition(ctx context.Context, id uint) error

	// Applicant
	AddApplicant(ctx context.Context, requisitionID uint, req *CreateApplicantRequest) error
	UpdateStage(ctx context.Context, id uint, changedByID uint, req *UpdateApplicantStageRequest) error
	GetApplicantsByRequisition(ctx context.Context, requisitionID uint) (*KanbanBoardResponse, error)
	GetApplicantDetail(ctx context.Context, id uint) (*ApplicantDetailResponse, error)
}

type service struct {
	repo               Repository
	storage            StorageProvider
	notification       NotificationProvider
	user               UserProvider
	transactionManager infrastructure.TransactionManager
}

func NewService(repo Repository, storage StorageProvider, notification NotificationProvider, user UserProvider, transactionManager infrastructure.TransactionManager) Service {
	return &service{repo, storage, notification, user, transactionManager}
}

func (s *service) CreateRequisition(ctx context.Context, requesterID uint, req *CreateRequisitionRequest) error {
	return s.transactionManager.RunInTransaction(ctx, func(ctx context.Context) error {
		jr := &JobRequisition{
			RequesterID:    requesterID,
			DepartmentID:   req.DepartmentID,
			Title:          req.Title,
			Description:    req.Description,
			Quantity:       req.Quantity,
			EmploymentType: req.EmploymentType,
			Priority:       req.Priority,
			Status:         constants.RequisitionStatusDraft,
		}

		if jr.Quantity < 1 {
			jr.Quantity = 1
		}

		if req.TargetDate != "" {
			t, err := time.Parse(constants.DefaultTimeFormat, req.TargetDate)
			if err != nil {
				return errors.New("invalid target_date format (expected YYYY-MM-DD)")
			}
			jr.TargetDate = &t
		}

		err := s.repo.CreateRequisition(ctx, jr)
		if err != nil {
			return err
		}

		return nil
	})
}

func (s *service) SubmitRequisition(ctx context.Context, id uint, requesterID uint) error {
	return s.transactionManager.RunInTransaction(ctx, func(ctx context.Context) error {
		jr, err := s.repo.FindRequisitionByID(ctx, id)
		if err != nil {
			return errors.New("requisition not found")
		}

		if jr.RequesterID != requesterID {
			return errors.New("you are not the requester of this requisition")
		}
		if jr.Status != constants.RequisitionStatusDraft {
			return fmt.Errorf("requisition cannot be submitted from status '%s'", jr.Status)
		}

		if err := s.repo.UpdateRequisitionStatus(ctx, id, constants.RequisitionStatusPending, nil, ""); err != nil {
			return err
		}

		approverIDs, err := s.user.FindApprovalUsers(ctx, constants.APPROVAL_REQUISITION)
		if err != nil {
			return err
		}

		go func() {
			_ = s.notification.BlastNotification(
				approverIDs,
				string(constants.NotificationTypeRequisitionApprovalReq),
				"Pengajuan Lowongan Baru",
				fmt.Sprintf("Pengajuan lowongan '%s' membutuhkan persetujuan Anda.", jr.Title),
				id,
			)
		}()

		return nil
	})
}

func (s *service) RequisitionAction(ctx context.Context, id uint, approverID uint, req *RequisitionActionRequest) error {
	return s.transactionManager.RunInTransaction(ctx, func(ctx context.Context) error {
		jr, err := s.repo.FindRequisitionByID(ctx, id)
		if err != nil {
			return errors.New("requisition not found")
		}

		if jr.Status != constants.RequisitionStatusPending {
			return fmt.Errorf("requisition is not in PENDING status (current: %s)", jr.Status)
		}

		var (
			newStatus       string
			rejectionReason string
		)
		switch constants.RequisitionAction(req.Action) {
		case constants.RequisitionActionApprove:
			newStatus = constants.RequisitionStatusApproved
		case constants.RequisitionActionReject:
			if req.RejectionReason == "" {
				return fmt.Errorf("rejection reason is required")
			}
			newStatus = constants.RequisitionStatusRejected
			rejectionReason = req.RejectionReason
		default:
			return fmt.Errorf("invalid action: %s", req.Action)
		}

		if err := s.repo.UpdateRequisitionStatus(ctx, id, newStatus, &approverID, rejectionReason); err != nil {
			return err
		}

		go func() {
			_ = s.notification.SendNotification(
				jr.RequesterID,
				string(constants.NotificationTypeRequisitionApprovalReq),
				fmt.Sprintf("Requisition %s", strings.Title(strings.ToLower(newStatus))),
				fmt.Sprintf("Your requisition '%s' has been %s.", jr.Title, newStatus),
				id,
			)
		}()

		return nil
	})
}

func (s *service) GetRequisitions(ctx context.Context, filter *RequisitionFilter) ([]RequisitionListResponse, *response.Meta, error) {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 {
		filter.Limit = 10
	}

	items, total, err := s.repo.FindAllRequisitions(ctx, filter)
	if err != nil {
		return nil, nil, err
	}

	var result []RequisitionListResponse
	for _, jr := range items {
		r := RequisitionListResponse{
			ID:             jr.ID,
			Title:          jr.Title,
			DepartmentID:   jr.DepartmentID,
			EmploymentType: jr.EmploymentType,
			Quantity:       jr.Quantity,
			Priority:       jr.Priority,
			Status:         jr.Status,
			RequesterID:    jr.RequesterID,
			TargetDate:     jr.TargetDate,
			CreatedAt:      jr.CreatedAt,
		}
		if jr.Requester != nil && jr.Requester.Employee != nil {
			r.RequesterName = jr.Requester.Employee.FullName
		}
		if jr.Department != nil {
			r.DepartmentName = jr.Department.Name
		}
		result = append(result, r)
	}

	meta := response.NewMetaOffset(filter.Page, filter.Limit, total)

	return result, meta, nil
}

func (s *service) GetRequisitionDetail(ctx context.Context, id uint) (*RequisitionDetailResponse, error) {
	jr, err := s.repo.FindRequisitionByID(ctx, id)
	if err != nil {
		return nil, errors.New("requisition not found")
	}

	result := &RequisitionDetailResponse{
		ID:              jr.ID,
		Title:           jr.Title,
		Description:     jr.Description,
		DepartmentID:    jr.DepartmentID,
		EmploymentType:  jr.EmploymentType,
		Quantity:        jr.Quantity,
		Priority:        jr.Priority,
		Status:          jr.Status,
		RequesterID:     jr.RequesterID,
		ApprovedBy:      jr.ApprovedBy,
		RejectionReason: jr.RejectionReason,
		TargetDate:      jr.TargetDate,
		CreatedAt:       jr.CreatedAt,
		UpdatedAt:       jr.UpdatedAt,
	}

	if jr.Requester != nil && jr.Requester.Employee != nil {
		result.RequesterName = jr.Requester.Employee.FullName
	}
	if jr.Approver != nil && jr.Approver.Employee != nil {
		result.ApproverName = jr.Approver.Employee.FullName
	}
	if jr.Department != nil {
		result.DepartmentName = jr.Department.Name
	}

	return result, nil
}

func (s *service) CloseRequisition(ctx context.Context, id uint) error {
	return s.transactionManager.RunInTransaction(ctx, func(ctx context.Context) error {
		jr, err := s.repo.FindRequisitionByID(ctx, id)
		if err != nil {
			return errors.New("requisition not found")
		}
		if jr.Status == constants.RequisitionStatusClosed {
			return errors.New("requisition is already closed")
		}
		err = s.repo.UpdateRequisitionStatus(ctx, id, constants.RequisitionStatusClosed, jr.ApprovedBy, jr.RejectionReason)
		if err != nil {
			return err
		}

		return nil
	})
}

func (s *service) DeleteRequisition(ctx context.Context, id uint) error {
	return s.transactionManager.RunInTransaction(ctx, func(ctx context.Context) error {
		_, err := s.repo.FindRequisitionByID(ctx, id)
		if err != nil {
			return errors.New("requisition not found")
		}
		err = s.repo.SoftDeleteRequisition(ctx, id)
		if err != nil {
			return err
		}

		return nil
	})
}

func (s *service) AddApplicant(ctx context.Context, requisitionID uint, req *CreateApplicantRequest) error {
	return s.transactionManager.RunInTransaction(ctx, func(ctx context.Context) error {
		_, err := s.repo.FindRequisitionByID(ctx, requisitionID)
		if err != nil {
			return errors.New("requisition not found")
		}

		count, err := s.repo.CountApplicantsByRequisitionAndStage(ctx, requisitionID, constants.ApplicantStageScreening)
		if err != nil {
			return err
		}

		applicant := &Applicant{
			JobRequisitionID: requisitionID,
			FullName:         req.FullName,
			Email:            req.Email,
			PhoneNumber:      req.PhoneNumber,
			Stage:            constants.ApplicantStageScreening,
			StageOrder:       int(count),
		}

		if req.ResumeBase64 != "" {
			fileBytes, err := utils.DecodeBase64Image(req.ResumeBase64)
			if err == nil && len(fileBytes) > 0 {
				objectName := fmt.Sprintf("resumes/%d-%s.pdf", requisitionID, req.FullName)
				url, uploadErr := s.storage.UploadFileByte(ctx, objectName, bytes.NewReader(fileBytes), int64(len(fileBytes)), "application/pdf")
				if uploadErr == nil {
					applicant.ResumeURL = url
				} else {
					logger.Errorf("resume upload failed: %v", uploadErr)
				}
			}
		}

		err = s.repo.CreateApplicant(ctx, applicant)
		if err != nil {
			return err
		}

		return nil
	})
}

func (s *service) UpdateStage(ctx context.Context, id uint, changedByID uint, req *UpdateApplicantStageRequest) error {
	return s.transactionManager.RunInTransaction(ctx, func(ctx context.Context) error {
		applicant, err := s.repo.FindApplicantByID(ctx, id)
		if err != nil {
			return errors.New("applicant not found")
		}

		fromStage := applicant.Stage

		count, err := s.repo.CountApplicantsByRequisitionAndStage(ctx, applicant.JobRequisitionID, req.Stage)
		if err != nil {
			return err
		}

		if err := s.repo.UpdateApplicantStage(ctx, id, req.Stage, int(count), req.Notes, req.RejectionReason); err != nil {
			return err
		}

		history := &ApplicantStageHistory{
			ApplicantID: id,
			FromStage:   fromStage,
			ToStage:     req.Stage,
			ChangedBy:   changedByID,
			Notes:       req.Notes,
		}

		err = s.repo.CreateStageHistory(ctx, history)
		if err != nil {
			return err
		}

		return nil
	})
}

func (s *service) GetApplicantsByRequisition(ctx context.Context, requisitionID uint) (*KanbanBoardResponse, error) {
	applicants, err := s.repo.FindApplicantsByRequisitionID(ctx, requisitionID)
	if err != nil {
		return nil, err
	}

	board := &KanbanBoardResponse{
		Screening: []ApplicantListResponse{},
		Interview: []ApplicantListResponse{},
		Offering:  []ApplicantListResponse{},
		Hired:     []ApplicantListResponse{},
		Rejected:  []ApplicantListResponse{},
	}

	for _, a := range applicants {
		item := ApplicantListResponse{
			ID:               a.ID,
			JobRequisitionID: a.JobRequisitionID,
			FullName:         a.FullName,
			Email:            a.Email,
			PhoneNumber:      a.PhoneNumber,
			ResumeURL:        a.ResumeURL,
			Stage:            a.Stage,
			StageOrder:       a.StageOrder,
			CreatedAt:        a.CreatedAt,
		}
		switch a.Stage {
		case constants.ApplicantStageScreening:
			board.Screening = append(board.Screening, item)
		case constants.ApplicantStageInterview:
			board.Interview = append(board.Interview, item)
		case constants.ApplicantStageOffering:
			board.Offering = append(board.Offering, item)
		case constants.ApplicantStageHired:
			board.Hired = append(board.Hired, item)
		case constants.ApplicantStageRejected:
			board.Rejected = append(board.Rejected, item)
		}
	}

	return board, nil
}

func (s *service) GetApplicantDetail(ctx context.Context, id uint) (*ApplicantDetailResponse, error) {
	applicant, err := s.repo.FindApplicantByID(ctx, id)
	if err != nil {
		return nil, errors.New("applicant not found")
	}

	result := &ApplicantDetailResponse{
		ID:               applicant.ID,
		JobRequisitionID: applicant.JobRequisitionID,
		FullName:         applicant.FullName,
		Email:            applicant.Email,
		PhoneNumber:      applicant.PhoneNumber,
		ResumeURL:        applicant.ResumeURL,
		Stage:            applicant.Stage,
		Notes:            applicant.Notes,
		RejectionReason:  applicant.RejectionReason,
		CreatedAt:        applicant.CreatedAt,
		StageHistories:   []StageHistoryResponse{},
	}

	for _, h := range applicant.StageHistories {
		hr := StageHistoryResponse{
			ID:        h.ID,
			FromStage: h.FromStage,
			ToStage:   h.ToStage,
			Notes:     h.Notes,
			CreatedAt: h.CreatedAt,
		}
		if h.ChangedByUser != nil && h.ChangedByUser.Employee != nil {
			hr.ChangedByName = h.ChangedByUser.Employee.FullName
		}
		result.StageHistories = append(result.StageHistories, hr)
	}

	return result, nil
}
