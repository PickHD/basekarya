package recruitment

import (
	"context"
	"io"

	"basekarya-backend/internal/modules/onboarding"
	"basekarya-backend/internal/testutil"
	"basekarya-backend/pkg/response"

	"github.com/stretchr/testify/mock"
)

type mockRepo struct{ mock.Mock }

func (m *mockRepo) CreateRequisition(ctx context.Context, req *JobRequisition) error {
	return m.Called(ctx, req).Error(0)
}

func (m *mockRepo) FindRequisitionByID(ctx context.Context, id uint) (*JobRequisition, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*JobRequisition), args.Error(1)
}

func (m *mockRepo) FindAllRequisitions(ctx context.Context, filter *RequisitionFilter) ([]JobRequisition, int64, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).([]JobRequisition), args.Get(1).(int64), args.Error(2)
}

func (m *mockRepo) UpdateRequisitionStatus(ctx context.Context, id uint, status string, approvedBy *uint, rejectionReason string) error {
	return m.Called(ctx, id, status, approvedBy, rejectionReason).Error(0)
}

func (m *mockRepo) SoftDeleteRequisition(ctx context.Context, id uint) error {
	return m.Called(ctx, id).Error(0)
}

func (m *mockRepo) CreateApplicant(ctx context.Context, applicant *Applicant) error {
	return m.Called(ctx, applicant).Error(0)
}

func (m *mockRepo) FindApplicantByID(ctx context.Context, id uint) (*Applicant, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Applicant), args.Error(1)
}

func (m *mockRepo) FindApplicantsByRequisitionID(ctx context.Context, requisitionID uint) ([]Applicant, error) {
	args := m.Called(ctx, requisitionID)
	return args.Get(0).([]Applicant), args.Error(1)
}

func (m *mockRepo) UpdateApplicantStage(ctx context.Context, id uint, stage string, stageOrder int, notes, rejectionReason string) error {
	return m.Called(ctx, id, stage, stageOrder, notes, rejectionReason).Error(0)
}

func (m *mockRepo) CreateStageHistory(ctx context.Context, history *ApplicantStageHistory) error {
	return m.Called(ctx, history).Error(0)
}

func (m *mockRepo) CountApplicantsByRequisitionAndStage(ctx context.Context, requisitionID uint, stage string) (int64, error) {
	args := m.Called(ctx, requisitionID, stage)
	return args.Get(0).(int64), args.Error(1)
}

func (m *mockRepo) SoftDeleteApplicant(ctx context.Context, id uint) error {
	return m.Called(ctx, id).Error(0)
}

type mockStorage struct{ mock.Mock }

func (m *mockStorage) UploadFileByte(ctx context.Context, objectName string, reader io.Reader, size int64, contentType string) (string, error) {
	args := m.Called(ctx, objectName, reader, size, contentType)
	return args.String(0), args.Error(1)
}

type mockNotification struct{ mock.Mock }

func (m *mockNotification) SendNotification(ctx context.Context, userID uint, notifType string, title string, message string, relatedID uint) error {
	return m.Called(ctx, userID, notifType, title, message, relatedID).Error(0)
}

func (m *mockNotification) BlastNotification(ctx context.Context, userIDs []uint, notifType string, title string, message string, relatedID uint) error {
	return m.Called(ctx, userIDs, notifType, title, message, relatedID).Error(0)
}

type mockUserProvider struct{ mock.Mock }

func (m *mockUserProvider) FindApprovalUsers(ctx context.Context, permissionApprovalName string) ([]uint, error) {
	args := m.Called(ctx, permissionApprovalName)
	return args.Get(0).([]uint), args.Error(1)
}

type mockOnboarding struct{ mock.Mock }

func (m *mockOnboarding) CreateWorkflow(ctx context.Context, req *onboarding.CreateWorkflowRequest) error {
	return m.Called(ctx, req).Error(0)
}

type mockService struct{ mock.Mock }

func (m *mockService) CreateRequisition(ctx context.Context, requesterID uint, req *CreateRequisitionRequest) error {
	return m.Called(ctx, requesterID, req).Error(0)
}

func (m *mockService) SubmitRequisition(ctx context.Context, id uint, requesterID uint) error {
	return m.Called(ctx, id, requesterID).Error(0)
}

func (m *mockService) RequisitionAction(ctx context.Context, id uint, approverID uint, req *RequisitionActionRequest) error {
	return m.Called(ctx, id, approverID, req).Error(0)
}

func (m *mockService) GetRequisitions(ctx context.Context, filter *RequisitionFilter) ([]RequisitionListResponse, *response.Meta, error) {
	args := m.Called(ctx, filter)
	var meta *response.Meta
	if args.Get(1) != nil {
		meta = args.Get(1).(*response.Meta)
	}
	return args.Get(0).([]RequisitionListResponse), meta, args.Error(2)
}

func (m *mockService) GetRequisitionDetail(ctx context.Context, id uint) (*RequisitionDetailResponse, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*RequisitionDetailResponse), args.Error(1)
}

func (m *mockService) CloseRequisition(ctx context.Context, id uint) error {
	return m.Called(ctx, id).Error(0)
}

func (m *mockService) DeleteRequisition(ctx context.Context, id uint) error {
	return m.Called(ctx, id).Error(0)
}

func (m *mockService) AddApplicant(ctx context.Context, requisitionID uint, req *CreateApplicantRequest) error {
	return m.Called(ctx, requisitionID, req).Error(0)
}

func (m *mockService) UpdateStage(ctx context.Context, id uint, changedByID uint, req *UpdateApplicantStageRequest) error {
	return m.Called(ctx, id, changedByID, req).Error(0)
}

func (m *mockService) GetApplicantsByRequisition(ctx context.Context, requisitionID uint) (*KanbanBoardResponse, error) {
	args := m.Called(ctx, requisitionID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*KanbanBoardResponse), args.Error(1)
}

func (m *mockService) GetApplicantDetail(ctx context.Context, id uint) (*ApplicantDetailResponse, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ApplicantDetailResponse), args.Error(1)
}

func newTestService() (Service, *mockRepo, *mockStorage, *mockNotification, *mockUserProvider, *mockOnboarding) {
	repo := new(mockRepo)
	storage := new(mockStorage)
	notif := new(mockNotification)
	userProv := new(mockUserProvider)
	onboard := new(mockOnboarding)
	tm := testutil.NewMockTransactionManager()
	svc := NewService(repo, storage, notif, userProv, onboard, tm)
	return svc, repo, storage, notif, userProv, onboard
}
