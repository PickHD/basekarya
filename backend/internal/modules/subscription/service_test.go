package subscription

import (
	stdcontext "context"
	"errors"
	"testing"

	"basekarya-backend/internal/modules/company"
	"basekarya-backend/internal/testutil"
	"basekarya-backend/pkg/constants"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestService_ListPlans(t *testing.T) {
	tests := []struct {
		name       string
		setupMocks func(*mockRepo)
		wantLen    int
		wantErr    bool
	}{
		{
			name: "success",
			setupMocks: func(repo *mockRepo) {
				repo.On("FindAllPlans", mock.Anything).Return([]SubscriptionPlan{
					{ID: 1, Name: "Basic", Slug: "basic", PriceMonthly: 0},
					{ID: 2, Name: "Pro", Slug: "pro", PriceMonthly: 99},
				}, nil)
			},
			wantLen: 2,
			wantErr: false,
		},
		{
			name: "repo error",
			setupMocks: func(repo *mockRepo) {
				repo.On("FindAllPlans", mock.Anything).Return(nil, errors.New("db error"))
			},
			wantLen: 0,
			wantErr: true,
		},
		{
			name: "empty plans",
			setupMocks: func(repo *mockRepo) {
				repo.On("FindAllPlans", mock.Anything).Return([]SubscriptionPlan{}, nil)
			},
			wantLen: 0,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, repo, _, _, _, _ := newTestSubscriptionService()
			tt.setupMocks(repo)

			result, err := svc.ListPlans(testutil.CtxWithTenant(1, 1, false))

			if tt.wantErr {
				require.Error(t, err)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				assert.Len(t, result, tt.wantLen)
			}
		})
	}
}

func TestService_RequestUpgrade(t *testing.T) {
	proPlanID := uint(2)
	basicPlanID := uint(1)

	tests := []struct {
		name       string
		ctx        func() stdcontext.Context
		req        *UpgradeRequest
		setupMocks func(*mockRepo, *mockCompanyRepo)
		wantErr    bool
		errMsg     string
	}{
		{
			name: "success with no current plan",
			ctx: func() stdcontext.Context {
				return testutil.CtxWithTenant(1, 1, false)
			},
			req: &UpgradeRequest{PlanSlug: "pro"},
			setupMocks: func(repo *mockRepo, companyRepo *mockCompanyRepo) {
				repo.On("FindPlanBySlug", mock.Anything, "pro").Return(&SubscriptionPlan{ID: 2, Name: "Pro", Slug: "pro", PriceMonthly: 99}, nil)
				companyRepo.On("FindByID", mock.Anything, uint(1)).Return(&company.Company{ID: 1, SubscriptionPlanID: nil}, nil)
				repo.On("FindPendingRequestByCompanyID", mock.Anything, uint(1)).Return(nil, errors.New("not found"))
				repo.On("CreateRequest", mock.Anything, mock.AnythingOfType("*subscription.SubscriptionRequest")).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "success with existing plan upgrade",
			ctx: func() stdcontext.Context {
				return testutil.CtxWithTenant(1, 1, false)
			},
			req: &UpgradeRequest{PlanSlug: "pro"},
			setupMocks: func(repo *mockRepo, companyRepo *mockCompanyRepo) {
				repo.On("FindPlanBySlug", mock.Anything, "pro").Return(&SubscriptionPlan{ID: proPlanID, Name: "Pro", Slug: "pro", PriceMonthly: 99}, nil)
				companyRepo.On("FindByID", mock.Anything, uint(1)).Return(&company.Company{ID: 1, SubscriptionPlanID: &basicPlanID}, nil)
				repo.On("FindPlanByID", mock.Anything, basicPlanID).Return(&SubscriptionPlan{ID: basicPlanID, Name: "Basic", PriceMonthly: 0}, nil)
				repo.On("FindPendingRequestByCompanyID", mock.Anything, uint(1)).Return(nil, errors.New("not found"))
				repo.On("CreateRequest", mock.Anything, mock.AnythingOfType("*subscription.SubscriptionRequest")).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "error company not in context",
			ctx: func() stdcontext.Context {
				return testutil.CtxWithTenant(0, 1, false)
			},
			req:        &UpgradeRequest{PlanSlug: "pro"},
			setupMocks: func(repo *mockRepo, companyRepo *mockCompanyRepo) {},
			wantErr:    true,
			errMsg:     "company not found in context",
		},
		{
			name: "error plan not found",
			ctx: func() stdcontext.Context {
				return testutil.CtxWithTenant(1, 1, false)
			},
			req: &UpgradeRequest{PlanSlug: "pro"},
			setupMocks: func(repo *mockRepo, companyRepo *mockCompanyRepo) {
				repo.On("FindPlanBySlug", mock.Anything, "pro").Return(nil, errors.New("not found"))
			},
			wantErr: true,
			errMsg:  "plan not found",
		},
		{
			name: "error company not found",
			ctx: func() stdcontext.Context {
				return testutil.CtxWithTenant(1, 1, false)
			},
			req: &UpgradeRequest{PlanSlug: "pro"},
			setupMocks: func(repo *mockRepo, companyRepo *mockCompanyRepo) {
				repo.On("FindPlanBySlug", mock.Anything, "pro").Return(&SubscriptionPlan{ID: 2, Name: "Pro", PriceMonthly: 99}, nil)
				companyRepo.On("FindByID", mock.Anything, uint(1)).Return(nil, errors.New("not found"))
			},
			wantErr: true,
			errMsg:  "company not found",
		},
		{
			name: "error already on this plan",
			ctx: func() stdcontext.Context {
				return testutil.CtxWithTenant(1, 1, false)
			},
			req: &UpgradeRequest{PlanSlug: "pro"},
			setupMocks: func(repo *mockRepo, companyRepo *mockCompanyRepo) {
				repo.On("FindPlanBySlug", mock.Anything, "pro").Return(&SubscriptionPlan{ID: proPlanID, Name: "Pro", PriceMonthly: 99}, nil)
				companyRepo.On("FindByID", mock.Anything, uint(1)).Return(&company.Company{ID: 1, SubscriptionPlanID: &proPlanID}, nil)
			},
			wantErr: true,
			errMsg:  "already on this plan",
		},
		{
			name: "error can only upgrade to higher plan",
			ctx: func() stdcontext.Context {
				return testutil.CtxWithTenant(1, 1, false)
			},
			req: &UpgradeRequest{PlanSlug: "basic"},
			setupMocks: func(repo *mockRepo, companyRepo *mockCompanyRepo) {
				repo.On("FindPlanBySlug", mock.Anything, "basic").Return(&SubscriptionPlan{ID: basicPlanID, Name: "Basic", PriceMonthly: 0}, nil)
				companyRepo.On("FindByID", mock.Anything, uint(1)).Return(&company.Company{ID: 1, SubscriptionPlanID: &proPlanID}, nil)
				repo.On("FindPlanByID", mock.Anything, proPlanID).Return(&SubscriptionPlan{ID: proPlanID, Name: "Pro", PriceMonthly: 99}, nil)
			},
			wantErr: true,
			errMsg:  "can only upgrade to a higher plan",
		},
		{
			name: "error pending request exists",
			ctx: func() stdcontext.Context {
				return testutil.CtxWithTenant(1, 1, false)
			},
			req: &UpgradeRequest{PlanSlug: "pro"},
			setupMocks: func(repo *mockRepo, companyRepo *mockCompanyRepo) {
				repo.On("FindPlanBySlug", mock.Anything, "pro").Return(&SubscriptionPlan{ID: proPlanID, Name: "Pro", PriceMonthly: 99}, nil)
				companyRepo.On("FindByID", mock.Anything, uint(1)).Return(&company.Company{ID: 1, SubscriptionPlanID: &basicPlanID}, nil)
				repo.On("FindPlanByID", mock.Anything, basicPlanID).Return(&SubscriptionPlan{ID: basicPlanID, Name: "Basic", PriceMonthly: 0}, nil)
				repo.On("FindPendingRequestByCompanyID", mock.Anything, uint(1)).Return(&SubscriptionRequest{ID: 10}, nil)
			},
			wantErr: true,
			errMsg:  "you already have a pending upgrade request",
		},
		{
			name: "error create request fails",
			ctx: func() stdcontext.Context {
				return testutil.CtxWithTenant(1, 1, false)
			},
			req: &UpgradeRequest{PlanSlug: "pro"},
			setupMocks: func(repo *mockRepo, companyRepo *mockCompanyRepo) {
				repo.On("FindPlanBySlug", mock.Anything, "pro").Return(&SubscriptionPlan{ID: 2, Name: "Pro", PriceMonthly: 99}, nil)
				companyRepo.On("FindByID", mock.Anything, uint(1)).Return(&company.Company{ID: 1, SubscriptionPlanID: nil}, nil)
				repo.On("FindPendingRequestByCompanyID", mock.Anything, uint(1)).Return(nil, errors.New("not found"))
				repo.On("CreateRequest", mock.Anything, mock.AnythingOfType("*subscription.SubscriptionRequest")).Return(errors.New("db error"))
			},
			wantErr: true,
			errMsg:  "failed to create upgrade request: db error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, repo, companyRepo, _, _, _ := newTestSubscriptionService()
			tt.setupMocks(repo, companyRepo)

			result, err := svc.RequestUpgrade(tt.ctx(), tt.req)

			if tt.wantErr {
				require.Error(t, err)
				assert.Equal(t, tt.errMsg, err.Error())
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
			}
		})
	}
}

func TestService_ListPendingRequests(t *testing.T) {
	tests := []struct {
		name       string
		setupMocks func(*mockRepo)
		wantLen    int
		wantErr    bool
	}{
		{
			name: "success",
			setupMocks: func(repo *mockRepo) {
				repo.On("FindAllPendingRequests", mock.Anything).Return([]SubscriptionRequestResponse{
					{ID: 1, CompanyName: "Acme", Status: "PENDING"},
				}, nil)
			},
			wantLen: 1,
			wantErr: false,
		},
		{
			name: "repo error",
			setupMocks: func(repo *mockRepo) {
				repo.On("FindAllPendingRequests", mock.Anything).Return(nil, errors.New("db error"))
			},
			wantLen: 0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, repo, _, _, _, _ := newTestSubscriptionService()
			tt.setupMocks(repo)

			result, err := svc.ListPendingRequests(testutil.CtxWithTenant(1, 1, false))

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Len(t, result, tt.wantLen)
			}
		})
	}
}

func TestService_ReviewRequest(t *testing.T) {
	proPlanID := uint(2)
	reviewerID := uint(10)

	tests := []struct {
		name       string
		requestID  uint
		ctx        func() stdcontext.Context
		req        *ReviewRequest
		setupMocks func(*mockRepo, *mockCompanyRepo, *mockRole, *mockUser, *mockCache)
		wantErr    bool
		errMsg     string
	}{
		{
			name:      "approve success",
			requestID: 1,
			ctx: func() stdcontext.Context {
				return testutil.CtxWithTenant(1, reviewerID, true)
			},
			req: &ReviewRequest{Status: constants.SubReqStatusApproved, Notes: "ok"},
			setupMocks: func(repo *mockRepo, companyRepo *mockCompanyRepo, role *mockRole, user *mockUser, cache *mockCache) {
				repo.On("FindRequestByID", mock.Anything, uint(1)).Return(&SubscriptionRequest{
					ID: 1, CompanyID: 1, CurrentPlanID: 1, RequestedPlanID: proPlanID, Status: constants.SubReqStatusPending,
				}, nil)
				repo.On("UpdateRequest", mock.Anything, mock.AnythingOfType("*subscription.SubscriptionRequest")).Return(nil)
				companyRepo.On("FindByID", mock.Anything, uint(1)).Return(&company.Company{ID: 1, SubscriptionPlanID: nil}, nil)
				companyRepo.On("Update", mock.Anything, mock.AnythingOfType("*company.Company")).Return(nil)
				repo.On("FindPlanByID", mock.Anything, proPlanID).Return(&SubscriptionPlan{ID: proPlanID, Slug: "pro"}, nil)
				role.On("FindPermissionIDsByGroupNames", mock.Anything, mock.Anything).Return([]uint{1, 2}, nil)
				role.On("FindRoleIDsByCompanyID", mock.Anything, uint(1)).Return([]uint{1}, nil)
				role.On("AssignPermissions", mock.Anything, uint(1), []uint{1, 2}, uint(1)).Return(nil)
				cache.On("Del", mock.Anything, "subscription:features:1").Return(nil)
				cache.On("Del", mock.Anything, "company:profile:1").Return(nil)
				user.On("ForceResetPasswordByCompanyID", mock.Anything, uint(1)).Return(nil)
			},
			wantErr: false,
		},
		{
			name:      "reject success",
			requestID: 1,
			ctx: func() stdcontext.Context {
				return testutil.CtxWithTenant(1, reviewerID, true)
			},
			req: &ReviewRequest{Status: constants.SubReqStatusRejected, Notes: "not eligible"},
			setupMocks: func(repo *mockRepo, companyRepo *mockCompanyRepo, role *mockRole, user *mockUser, cache *mockCache) {
				repo.On("FindRequestByID", mock.Anything, uint(1)).Return(&SubscriptionRequest{
					ID: 1, CompanyID: 1, Status: constants.SubReqStatusPending,
				}, nil)
				repo.On("UpdateRequest", mock.Anything, mock.AnythingOfType("*subscription.SubscriptionRequest")).Return(nil)
			},
			wantErr: false,
		},
		{
			name:      "error request not found",
			requestID: 99,
			ctx: func() stdcontext.Context {
				return testutil.CtxWithTenant(1, reviewerID, true)
			},
			req: &ReviewRequest{Status: constants.SubReqStatusApproved},
			setupMocks: func(repo *mockRepo, companyRepo *mockCompanyRepo, role *mockRole, user *mockUser, cache *mockCache) {
				repo.On("FindRequestByID", mock.Anything, uint(99)).Return(nil, errors.New("not found"))
			},
			wantErr: true,
			errMsg:  "request not found",
		},
		{
			name:      "error request already reviewed",
			requestID: 1,
			ctx: func() stdcontext.Context {
				return testutil.CtxWithTenant(1, reviewerID, true)
			},
			req: &ReviewRequest{Status: constants.SubReqStatusApproved},
			setupMocks: func(repo *mockRepo, companyRepo *mockCompanyRepo, role *mockRole, user *mockUser, cache *mockCache) {
				repo.On("FindRequestByID", mock.Anything, uint(1)).Return(&SubscriptionRequest{
					ID: 1, Status: constants.SubReqStatusApproved,
				}, nil)
			},
			wantErr: true,
			errMsg:  "request already reviewed",
		},
		{
			name:      "error update request fails",
			requestID: 1,
			ctx: func() stdcontext.Context {
				return testutil.CtxWithTenant(1, reviewerID, true)
			},
			req: &ReviewRequest{Status: constants.SubReqStatusRejected},
			setupMocks: func(repo *mockRepo, companyRepo *mockCompanyRepo, role *mockRole, user *mockUser, cache *mockCache) {
				repo.On("FindRequestByID", mock.Anything, uint(1)).Return(&SubscriptionRequest{
					ID: 1, Status: constants.SubReqStatusPending,
				}, nil)
				repo.On("UpdateRequest", mock.Anything, mock.AnythingOfType("*subscription.SubscriptionRequest")).Return(errors.New("db error"))
			},
			wantErr: true,
			errMsg:  "failed to update request: db error",
		},
		{
			name:      "approve error company not found",
			requestID: 1,
			ctx: func() stdcontext.Context {
				return testutil.CtxWithTenant(1, reviewerID, true)
			},
			req: &ReviewRequest{Status: constants.SubReqStatusApproved},
			setupMocks: func(repo *mockRepo, companyRepo *mockCompanyRepo, role *mockRole, user *mockUser, cache *mockCache) {
				repo.On("FindRequestByID", mock.Anything, uint(1)).Return(&SubscriptionRequest{
					ID: 1, CompanyID: 1, RequestedPlanID: proPlanID, Status: constants.SubReqStatusPending,
				}, nil)
				repo.On("UpdateRequest", mock.Anything, mock.AnythingOfType("*subscription.SubscriptionRequest")).Return(nil)
				companyRepo.On("FindByID", mock.Anything, uint(1)).Return(nil, errors.New("not found"))
			},
			wantErr: true,
			errMsg:  "company not found",
		},
		{
			name:      "approve error update company fails",
			requestID: 1,
			ctx: func() stdcontext.Context {
				return testutil.CtxWithTenant(1, reviewerID, true)
			},
			req: &ReviewRequest{Status: constants.SubReqStatusApproved},
			setupMocks: func(repo *mockRepo, companyRepo *mockCompanyRepo, role *mockRole, user *mockUser, cache *mockCache) {
				repo.On("FindRequestByID", mock.Anything, uint(1)).Return(&SubscriptionRequest{
					ID: 1, CompanyID: 1, RequestedPlanID: proPlanID, Status: constants.SubReqStatusPending,
				}, nil)
				repo.On("UpdateRequest", mock.Anything, mock.AnythingOfType("*subscription.SubscriptionRequest")).Return(nil)
				companyRepo.On("FindByID", mock.Anything, uint(1)).Return(&company.Company{ID: 1}, nil)
				companyRepo.On("Update", mock.Anything, mock.AnythingOfType("*company.Company")).Return(errors.New("db error"))
			},
			wantErr: true,
			errMsg:  "failed to update company plan: db error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, repo, companyRepo, role, user, cache := newTestSubscriptionService()
			tt.setupMocks(repo, companyRepo, role, user, cache)

			err := svc.ReviewRequest(tt.ctx(), tt.requestID, tt.req)

			if tt.wantErr {
				require.Error(t, err)
				assert.Equal(t, tt.errMsg, err.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestService_ListCompanies(t *testing.T) {
	tests := []struct {
		name       string
		search     string
		setupMocks func(*mockRepo)
		wantLen    int
		wantErr    bool
	}{
		{
			name:   "success",
			search: "",
			setupMocks: func(repo *mockRepo) {
				repo.On("FindAllCompanies", mock.Anything, "").Return([]CompanyListItem{
					{ID: 1, Name: "Acme"},
				}, nil)
			},
			wantLen: 1,
			wantErr: false,
		},
		{
			name:   "with search",
			search: "acme",
			setupMocks: func(repo *mockRepo) {
				repo.On("FindAllCompanies", mock.Anything, "acme").Return([]CompanyListItem{
					{ID: 1, Name: "Acme"},
				}, nil)
			},
			wantLen: 1,
			wantErr: false,
		},
		{
			name:   "repo error",
			search: "",
			setupMocks: func(repo *mockRepo) {
				repo.On("FindAllCompanies", mock.Anything, "").Return(nil, errors.New("db error"))
			},
			wantLen: 0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, repo, _, _, _, _ := newTestSubscriptionService()
			tt.setupMocks(repo)

			result, err := svc.ListCompanies(testutil.CtxWithTenant(1, 1, true), tt.search)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Len(t, result, tt.wantLen)
			}
		})
	}
}

func TestService_GetCompanyDetail(t *testing.T) {
	tests := []struct {
		name       string
		id         uint
		setupMocks func(*mockRepo)
		wantErr    bool
	}{
		{
			name: "success",
			id:   1,
			setupMocks: func(repo *mockRepo) {
				repo.On("FindCompanyDetailByID", mock.Anything, uint(1)).Return(&CompanyDetail{ID: 1, Name: "Acme"}, nil)
			},
			wantErr: false,
		},
		{
			name: "not found",
			id:   99,
			setupMocks: func(repo *mockRepo) {
				repo.On("FindCompanyDetailByID", mock.Anything, uint(99)).Return(nil, errors.New("not found"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, repo, _, _, _, _ := newTestSubscriptionService()
			tt.setupMocks(repo)

			result, err := svc.GetCompanyDetail(testutil.CtxWithTenant(1, 1, true), tt.id)

			if tt.wantErr {
				require.Error(t, err)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
				assert.Equal(t, tt.id, result.ID)
			}
		})
	}
}

func TestService_UpdateCompanyStatus(t *testing.T) {
	tests := []struct {
		name       string
		companyID  uint
		req        *UpdateCompanyStatusRequest
		setupMocks func(*mockRepo)
		setupCache func(*mockCache)
		wantErr    bool
	}{
		{
			name:      "success",
			companyID: 1,
			req:       &UpdateCompanyStatusRequest{SubscriptionStatus: constants.SubStatusActive},
			setupMocks: func(repo *mockRepo) {
				repo.On("UpdateCompanyStatus", mock.Anything, uint(1), constants.SubStatusActive).Return(nil)
			},
			setupCache: func(cache *mockCache) {
				cache.On("Del", mock.Anything, "subscription:features:1").Return(nil)
				cache.On("Del", mock.Anything, "company:profile:1").Return(nil)
			},
			wantErr: false,
		},
		{
			name:      "repo error",
			companyID: 1,
			req:       &UpdateCompanyStatusRequest{SubscriptionStatus: constants.SubStatusExpired},
			setupMocks: func(repo *mockRepo) {
				repo.On("UpdateCompanyStatus", mock.Anything, uint(1), constants.SubStatusExpired).Return(errors.New("db error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, repo, _, _, _, cache := newTestSubscriptionService()
			tt.setupMocks(repo)
			if tt.setupCache != nil {
				tt.setupCache(cache)
			}

			err := svc.UpdateCompanyStatus(testutil.CtxWithTenant(1, 1, true), tt.companyID, tt.req)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestService_GetDashboardStats(t *testing.T) {
	tests := []struct {
		name       string
		setupMocks func(*mockRepo)
		wantErr    bool
	}{
		{
			name: "success",
			setupMocks: func(repo *mockRepo) {
				repo.On("GetDashboardStats", mock.Anything).Return(&DashboardStatsResponse{
					TotalCompanies:      10,
					ActiveSubscriptions: 8,
					TotalRevenue:        800,
				}, nil)
			},
			wantErr: false,
		},
		{
			name: "repo error",
			setupMocks: func(repo *mockRepo) {
				repo.On("GetDashboardStats", mock.Anything).Return(nil, errors.New("db error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, repo, _, _, _, _ := newTestSubscriptionService()
			tt.setupMocks(repo)

			result, err := svc.GetDashboardStats(testutil.CtxWithTenant(1, 1, true))

			if tt.wantErr {
				require.Error(t, err)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
			}
		})
	}
}
