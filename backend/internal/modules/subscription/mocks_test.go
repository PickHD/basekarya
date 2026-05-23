package subscription

import (
	"context"

	"basekarya-backend/internal/modules/company"

	"github.com/stretchr/testify/mock"
)

type mockRepo struct{ mock.Mock }

func (m *mockRepo) FindAllPlans(ctx context.Context) ([]SubscriptionPlan, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]SubscriptionPlan), args.Error(1)
}

func (m *mockRepo) FindPlanBySlug(ctx context.Context, slug string) (*SubscriptionPlan, error) {
	args := m.Called(ctx, slug)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*SubscriptionPlan), args.Error(1)
}

func (m *mockRepo) FindPlanByID(ctx context.Context, id uint) (*SubscriptionPlan, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*SubscriptionPlan), args.Error(1)
}

func (m *mockRepo) CreateRequest(ctx context.Context, req *SubscriptionRequest) error {
	return m.Called(ctx, req).Error(0)
}

func (m *mockRepo) FindPendingRequestByCompanyID(ctx context.Context, companyID uint) (*SubscriptionRequest, error) {
	args := m.Called(ctx, companyID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*SubscriptionRequest), args.Error(1)
}

func (m *mockRepo) FindRequestByID(ctx context.Context, id uint) (*SubscriptionRequest, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*SubscriptionRequest), args.Error(1)
}

func (m *mockRepo) FindAllPendingRequests(ctx context.Context) ([]SubscriptionRequestResponse, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]SubscriptionRequestResponse), args.Error(1)
}

func (m *mockRepo) FindAllRequests(ctx context.Context) ([]SubscriptionRequestResponse, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]SubscriptionRequestResponse), args.Error(1)
}

func (m *mockRepo) UpdateRequest(ctx context.Context, req *SubscriptionRequest) error {
	return m.Called(ctx, req).Error(0)
}

func (m *mockRepo) FindAllCompanies(ctx context.Context, search string) ([]CompanyListItem, error) {
	args := m.Called(ctx, search)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]CompanyListItem), args.Error(1)
}

func (m *mockRepo) FindCompanyDetailByID(ctx context.Context, id uint) (*CompanyDetail, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*CompanyDetail), args.Error(1)
}

func (m *mockRepo) UpdateCompanyStatus(ctx context.Context, companyID uint, status string) error {
	return m.Called(ctx, companyID, status).Error(0)
}

func (m *mockRepo) GetDashboardStats(ctx context.Context) (*DashboardStatsResponse, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*DashboardStatsResponse), args.Error(1)
}

type mockCompanyRepo struct{ mock.Mock }

func (m *mockCompanyRepo) FindByID(ctx context.Context, id uint) (*company.Company, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*company.Company), args.Error(1)
}

func (m *mockCompanyRepo) Update(ctx context.Context, c *company.Company) error {
	return m.Called(ctx, c).Error(0)
}

func (m *mockCompanyRepo) CreateCompany(ctx context.Context, c *company.Company) error {
	return m.Called(ctx, c).Error(0)
}

func (m *mockCompanyRepo) FindPlanIDBySlug(ctx context.Context, slug string) (uint, error) {
	args := m.Called(ctx, slug)
	return args.Get(0).(uint), args.Error(1)
}

func (m *mockCompanyRepo) FindPlanByCompanyID(ctx context.Context, companyID uint) (string, int, string, error) {
	args := m.Called(ctx, companyID)
	return args.String(0), args.Get(1).(int), args.String(2), args.Error(3)
}

func (m *mockCompanyRepo) FindModulesByCompanyID(ctx context.Context, companyID uint) ([]string, error) {
	args := m.Called(ctx, companyID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]string), args.Error(1)
}

type mockRole struct{ mock.Mock }

func (m *mockRole) FindPermissionIDsByGroupNames(ctx context.Context, groupNames []string) ([]uint, error) {
	args := m.Called(ctx, groupNames)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]uint), args.Error(1)
}

func (m *mockRole) FindRoleIDsByCompanyID(ctx context.Context, companyID uint) ([]uint, error) {
	args := m.Called(ctx, companyID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]uint), args.Error(1)
}

func (m *mockRole) AssignPermissions(ctx context.Context, roleID uint, permissionIDs []uint, companyID uint) error {
	return m.Called(ctx, roleID, permissionIDs, companyID).Error(0)
}

type mockUser struct{ mock.Mock }

func (m *mockUser) ForceResetPasswordByCompanyID(ctx context.Context, companyID uint) error {
	return m.Called(ctx, companyID).Error(0)
}

type mockCache struct{ mock.Mock }

func (m *mockCache) FlushDB(ctx context.Context) error {
	return m.Called(ctx).Error(0)
}

type mockService struct{ mock.Mock }

func (m *mockService) ListPlans(ctx context.Context) ([]PlanResponse, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]PlanResponse), args.Error(1)
}

func (m *mockService) RequestUpgrade(ctx context.Context, req *UpgradeRequest) (*UpgradeResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*UpgradeResponse), args.Error(1)
}

func (m *mockService) ListPendingRequests(ctx context.Context) ([]SubscriptionRequestResponse, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]SubscriptionRequestResponse), args.Error(1)
}

func (m *mockService) ListAllRequests(ctx context.Context) ([]SubscriptionRequestResponse, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]SubscriptionRequestResponse), args.Error(1)
}

func (m *mockService) ReviewRequest(ctx context.Context, requestID uint, req *ReviewRequest) error {
	return m.Called(ctx, requestID, req).Error(0)
}

func (m *mockService) ListCompanies(ctx context.Context, search string) ([]CompanyListItem, error) {
	args := m.Called(ctx, search)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]CompanyListItem), args.Error(1)
}

func (m *mockService) GetCompanyDetail(ctx context.Context, id uint) (*CompanyDetail, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*CompanyDetail), args.Error(1)
}

func (m *mockService) UpdateCompanyStatus(ctx context.Context, companyID uint, req *UpdateCompanyStatusRequest) error {
	return m.Called(ctx, companyID, req).Error(0)
}

func (m *mockService) GetDashboardStats(ctx context.Context) (*DashboardStatsResponse, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*DashboardStatsResponse), args.Error(1)
}

func newTestSubscriptionService() (Service, *mockRepo, *mockCompanyRepo, *mockRole, *mockUser, *mockCache) {
	repo := new(mockRepo)
	companyRepo := new(mockCompanyRepo)
	role := new(mockRole)
	user := new(mockUser)
	cache := new(mockCache)

	svc := NewService(repo, companyRepo, role, user, cache)
	return svc, repo, companyRepo, role, user, cache
}
