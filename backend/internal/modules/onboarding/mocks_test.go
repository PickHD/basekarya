package onboarding

import (
	"context"
	"io"

	"basekarya-backend/internal/modules/company"
	"basekarya-backend/internal/modules/department"
	"basekarya-backend/internal/modules/master"
	"basekarya-backend/internal/modules/rbac"
	"basekarya-backend/internal/modules/user"
	"basekarya-backend/internal/testutil"
	"basekarya-backend/pkg/response"

	"github.com/stretchr/testify/mock"
	"mime/multipart"
)

type mockRepo struct{ mock.Mock }

func (m *mockRepo) CreateTemplate(ctx context.Context, t *OnboardingTemplate) error {
	return m.Called(ctx, t).Error(0)
}

func (m *mockRepo) FindAllTemplates(ctx context.Context) ([]OnboardingTemplate, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]OnboardingTemplate), args.Error(1)
}

func (m *mockRepo) FindTemplateByID(ctx context.Context, id uint) (*OnboardingTemplate, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*OnboardingTemplate), args.Error(1)
}

func (m *mockRepo) UpdateTemplate(ctx context.Context, t *OnboardingTemplate) error {
	return m.Called(ctx, t).Error(0)
}

func (m *mockRepo) DeleteTemplate(ctx context.Context, id uint) error {
	return m.Called(ctx, id).Error(0)
}

func (m *mockRepo) CreateWorkflow(ctx context.Context, w *OnboardingWorkflow) error {
	return m.Called(ctx, w).Error(0)
}

func (m *mockRepo) CreateTasks(ctx context.Context, tasks []OnboardingTask) error {
	return m.Called(ctx, tasks).Error(0)
}

func (m *mockRepo) FindAllWorkflows(ctx context.Context, filter *WorkflowFilter) ([]OnboardingWorkflow, int64, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).([]OnboardingWorkflow), args.Get(1).(int64), args.Error(2)
}

func (m *mockRepo) FindWorkflowByID(ctx context.Context, id uint) (*OnboardingWorkflow, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*OnboardingWorkflow), args.Error(1)
}

func (m *mockRepo) MarkWorkflowEmailSent(ctx context.Context, id uint) error {
	return m.Called(ctx, id).Error(0)
}

func (m *mockRepo) FindTaskByID(ctx context.Context, id uint) (*OnboardingTask, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*OnboardingTask), args.Error(1)
}

func (m *mockRepo) CompleteTask(ctx context.Context, id uint, completedBy uint, notes string) error {
	return m.Called(ctx, id, completedBy, notes).Error(0)
}

func (m *mockRepo) CountPendingTasks(ctx context.Context, workflowID uint) (int64, error) {
	args := m.Called(ctx, workflowID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *mockRepo) CountTotalTasks(ctx context.Context, workflowID uint) (int64, error) {
	args := m.Called(ctx, workflowID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *mockRepo) MarkWorkflowCompleted(ctx context.Context, id uint) error {
	return m.Called(ctx, id).Error(0)
}

type mockNotificationProvider struct{ mock.Mock }

func (m *mockNotificationProvider) SendNotification(ctx context.Context, userID uint, Type string, Title string, Message string, relatedID uint) error {
	return m.Called(ctx, userID, Type, Title, Message, relatedID).Error(0)
}

func (m *mockNotificationProvider) BlastNotification(ctx context.Context, userIDs []uint, Type string, Title string, Message string, relatedID uint) error {
	return m.Called(ctx, userIDs, Type, Title, Message, relatedID).Error(0)
}

type mockUserProvider struct{ mock.Mock }

func (m *mockUserProvider) FindApprovalUsers(ctx context.Context, permissionApprovalName string) ([]uint, error) {
	args := m.Called(ctx, permissionApprovalName)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]uint), args.Error(1)
}

func (m *mockUserProvider) CreateEmployee(ctx context.Context, req *user.CreateEmployeeRequest) (*user.CreateEmployeeResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.CreateEmployeeResponse), args.Error(1)
}

type mockEmailProvider struct{ mock.Mock }

func (m *mockEmailProvider) Send(to string, subject string, htmlBody string) error {
	return m.Called(to, subject, htmlBody).Error(0)
}

type mockStorageProvider struct{ mock.Mock }

func (m *mockStorageProvider) UploadFileByte(ctx context.Context, objectName string, reader io.Reader, size int64, contentType string) (string, error) {
	args := m.Called(ctx, objectName, reader, size, contentType)
	return args.String(0), args.Error(1)
}

type mockCompanyProvider struct{ mock.Mock }

func (m *mockCompanyProvider) FindByID(ctx context.Context, id uint) (*company.Company, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*company.Company), args.Error(1)
}

type mockRoleProvider struct{ mock.Mock }

func (m *mockRoleProvider) FindRoleByName(ctx context.Context, name string) (*rbac.Role, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*rbac.Role), args.Error(1)
}

type mockDepartmentProvider struct{ mock.Mock }

func (m *mockDepartmentProvider) FindByName(ctx context.Context, name string) (*department.Department, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*department.Department), args.Error(1)
}

type mockMasterProvider struct{ mock.Mock }

func (m *mockMasterProvider) FindShiftByName(ctx context.Context, name string) (*master.Shift, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*master.Shift), args.Error(1)
}

type mockService struct{ mock.Mock }

func (m *mockService) CreateTemplate(ctx context.Context, req *CreateTemplateRequest) error {
	return m.Called(ctx, req).Error(0)
}

func (m *mockService) GetTemplates(ctx context.Context) ([]TemplateResponse, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]TemplateResponse), args.Error(1)
}

func (m *mockService) GetTemplateByID(ctx context.Context, id uint) (*TemplateResponse, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*TemplateResponse), args.Error(1)
}

func (m *mockService) UpdateTemplate(ctx context.Context, id uint, req *UpdateTemplateRequest) error {
	return m.Called(ctx, id, req).Error(0)
}

func (m *mockService) DeleteTemplate(ctx context.Context, id uint) error {
	return m.Called(ctx, id).Error(0)
}

func (m *mockService) CreateWorkflow(ctx context.Context, req *CreateWorkflowRequest) error {
	return m.Called(ctx, req).Error(0)
}

func (m *mockService) GetWorkflows(ctx context.Context, filter *WorkflowFilter) ([]WorkflowListResponse, *response.Meta, error) {
	args := m.Called(ctx, filter)
	var meta *response.Meta
	if args.Get(1) != nil {
		meta = args.Get(1).(*response.Meta)
	}
	if args.Get(0) == nil {
		return nil, meta, args.Error(2)
	}
	return args.Get(0).([]WorkflowListResponse), meta, args.Error(2)
}

func (m *mockService) GetWorkflowDetail(ctx context.Context, id uint) (*WorkflowDetailResponse, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*WorkflowDetailResponse), args.Error(1)
}

func (m *mockService) CompleteTask(ctx context.Context, taskID uint, completedByID uint, req *CompleteTaskRequest) error {
	return m.Called(ctx, taskID, completedByID, req).Error(0)
}

func newTestOnboardingService() (Service, *mockRepo, *mockNotificationProvider, *mockUserProvider, *mockEmailProvider, *mockCompanyProvider, *mockRoleProvider, *mockDepartmentProvider, *mockMasterProvider, *testutil.MockTransactionManager) {
	repo := new(mockRepo)
	notif := new(mockNotificationProvider)
	userProv := new(mockUserProvider)
	emailProv := new(mockEmailProvider)
	companyProv := new(mockCompanyProvider)
	roleProv := new(mockRoleProvider)
	deptProv := new(mockDepartmentProvider)
	masterProv := new(mockMasterProvider)
	tm := testutil.NewMockTransactionManager()

	svc := NewService(repo, notif, userProv, emailProv, companyProv, roleProv, deptProv, masterProv, tm)
	return svc, repo, notif, userProv, emailProv, companyProv, roleProv, deptProv, masterProv, tm
}

var _ company.StorageProvider = (*companyMockStorage)(nil)

type companyMockStorage struct{ mock.Mock }

func (m *companyMockStorage) UploadFileMultipart(ctx context.Context, file *multipart.FileHeader, objectName string) (string, error) {
	args := m.Called(ctx, file, objectName)
	return args.String(0), args.Error(1)
}
