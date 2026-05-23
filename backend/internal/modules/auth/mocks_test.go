package auth

import (
	"context"
	"time"

	"basekarya-backend/internal/modules/company"
	"basekarya-backend/internal/modules/rbac"
	"basekarya-backend/internal/modules/user"
	"basekarya-backend/internal/testutil"

	"github.com/stretchr/testify/mock"
)

type mockUserProvider struct{ mock.Mock }

func (m *mockUserProvider) FindByUsername(ctx context.Context, username string) (*user.User, error) {
	args := m.Called(ctx, username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.User), args.Error(1)
}

func (m *mockUserProvider) FindEmployeeByEmail(ctx context.Context, email string) (*user.Employee, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.Employee), args.Error(1)
}

func (m *mockUserProvider) UpdatePasswordByEmail(ctx context.Context, email string, password string) error {
	return m.Called(ctx, email, password).Error(0)
}

func (m *mockUserProvider) CreateUser(ctx context.Context, u *user.User) error {
	return m.Called(ctx, u).Error(0)
}

type mockCacheProvider struct{ mock.Mock }

func (m *mockCacheProvider) Get(ctx context.Context, key string) (string, error) {
	args := m.Called(ctx, key)
	return args.String(0), args.Error(1)
}

func (m *mockCacheProvider) Set(ctx context.Context, key string, value any, expiration time.Duration) error {
	return m.Called(ctx, key, value, expiration).Error(0)
}

func (m *mockCacheProvider) Del(ctx context.Context, key string) error {
	return m.Called(ctx, key).Error(0)
}

func (m *mockCacheProvider) FlushDB(ctx context.Context) error {
	return m.Called(ctx).Error(0)
}

type mockCompanyProvider struct{ mock.Mock }

func (m *mockCompanyProvider) CreateCompany(ctx context.Context, c *company.Company) error {
	return m.Called(ctx, c).Error(0)
}

func (m *mockCompanyProvider) FindPlanIDBySlug(ctx context.Context, slug string) (uint, error) {
	args := m.Called(ctx, slug)
	return args.Get(0).(uint), args.Error(1)
}

type mockRoleProvider struct{ mock.Mock }

func (m *mockRoleProvider) Create(ctx context.Context, role *rbac.Role) error {
	return m.Called(ctx, role).Error(0)
}

func (m *mockRoleProvider) FindRoleByName(ctx context.Context, name string) (*rbac.Role, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*rbac.Role), args.Error(1)
}

func (m *mockRoleProvider) FindAllPermissionIDs(ctx context.Context) ([]uint, error) {
	args := m.Called(ctx)
	return args.Get(0).([]uint), args.Error(1)
}

func (m *mockRoleProvider) FindPermissionIDsByGroupNames(ctx context.Context, groupNames []string) ([]uint, error) {
	args := m.Called(ctx, groupNames)
	return args.Get(0).([]uint), args.Error(1)
}

func (m *mockRoleProvider) AssignPermissions(ctx context.Context, roleID uint, permissionIDs []uint, companyID uint) error {
	return m.Called(ctx, roleID, permissionIDs, companyID).Error(0)
}

type mockMasterProvider struct{ mock.Mock }

func (m *mockMasterProvider) SeedDefaults(ctx context.Context, companyID uint) error {
	return m.Called(ctx, companyID).Error(0)
}

type mockService struct{ mock.Mock }

func (m *mockService) Login(ctx context.Context, username, password string) (*LoginResponse, error) {
	args := m.Called(ctx, username, password)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*LoginResponse), args.Error(1)
}

func (m *mockService) RegisterCompany(ctx context.Context, req *RegisterCompanyRequest) (*RegisterCompanyResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*RegisterCompanyResponse), args.Error(1)
}

func (m *mockService) SendOrResendOTP(ctx context.Context, req *SendOrResendOTPRequest) error {
	return m.Called(ctx, req).Error(0)
}

func (m *mockService) VerifyOTP(ctx context.Context, req *VerifyOTPRequest) (*VerifyOTPResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*VerifyOTPResponse), args.Error(1)
}

func (m *mockService) ResetPassword(ctx context.Context, req *ResetPasswordRequest) error {
	return m.Called(ctx, req).Error(0)
}

func newTestAuthService() (Service, *mockUserProvider, *testutil.MockHasher, *testutil.MockTokenProvider, *mockCacheProvider, *testutil.MockEmailProvider, *mockCompanyProvider, *mockRoleProvider, *mockMasterProvider) {
	u := new(mockUserProvider)
	h := new(testutil.MockHasher)
	tok := new(testutil.MockTokenProvider)
	c := new(mockCacheProvider)
	e := new(testutil.MockEmailProvider)
	cp := new(mockCompanyProvider)
	rp := new(mockRoleProvider)
	mp := new(mockMasterProvider)

	return NewService(u, h, tok, c, e, cp, rp, mp), u, h, tok, c, e, cp, rp, mp
}
