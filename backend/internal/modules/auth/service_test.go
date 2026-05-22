package auth

import (
	"basekarya-backend/internal/modules/company"
	"basekarya-backend/internal/modules/rbac"
	"basekarya-backend/internal/modules/user"
	"basekarya-backend/internal/testutil"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// === Mocks for auth service interfaces ===

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

// === Tests ===

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

func TestLogin_Success(t *testing.T) {
	svc, userProv, hasher, tokenProv, _, _, _, _, _ := newTestAuthService()
	ctx := context.Background()

	userProv.On("FindByUsername", ctx, "admin").Return(&user.User{
		ID: 1, Username: "admin", PasswordHash: "hashed", Role: &rbac.Role{Name: "SUPERADMIN"}, IsPlatformAdmin: true,
	}, nil)
	hasher.On("CheckPasswordHash", "pass123", "hashed").Return(true)
	tokenProv.On("GenerateToken", uint(1), uint(0), true, "SUPERADMIN", (*uint)(nil), []string(nil)).Return("jwt-token", nil)

	resp, err := svc.Login(ctx, "admin", "pass123")

	require.NoError(t, err)
	assert.Equal(t, "jwt-token", resp.Token)
	assert.False(t, resp.MustChangePassword)
}

func TestLogin_UserNotFound(t *testing.T) {
	svc, userProv, _, _, _, _, _, _, _ := newTestAuthService()
	ctx := context.Background()

	userProv.On("FindByUsername", ctx, "unknown").Return(nil, errors.New("not found"))

	resp, err := svc.Login(ctx, "unknown", "any")

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, "invalid credentials", err.Error())
}

func TestLogin_WrongPassword(t *testing.T) {
	svc, userProv, hasher, _, _, _, _, _, _ := newTestAuthService()
	ctx := context.Background()

	userProv.On("FindByUsername", ctx, "admin").Return(&user.User{
		ID: 1, Username: "admin", PasswordHash: "hashed", Role: &rbac.Role{Name: "SUPERADMIN"},
	}, nil)
	hasher.On("CheckPasswordHash", "wrong", "hashed").Return(false)

	resp, err := svc.Login(ctx, "admin", "wrong")

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, "invalid credentials", err.Error())
}

func TestVerifyOTP_Valid(t *testing.T) {
	svc, _, _, _, cache, _, _, _, _ := newTestAuthService()
	ctx := context.Background()

	cache.On("Get", ctx, "123456").Return("user@email.com", nil)

	resp, err := svc.VerifyOTP(ctx, &VerifyOTPRequest{Code: "123456"})

	require.NoError(t, err)
	assert.True(t, resp.IsValid)
}

func TestVerifyOTP_Invalid(t *testing.T) {
	svc, _, _, _, cache, _, _, _, _ := newTestAuthService()
	ctx := context.Background()

	cache.On("Get", ctx, "000000").Return("", errors.New("not found"))

	resp, err := svc.VerifyOTP(ctx, &VerifyOTPRequest{Code: "000000"})

	require.NoError(t, err)
	assert.False(t, resp.IsValid)
}

func TestResetPassword_Success(t *testing.T) {
	svc, userProv, hasher, _, cache, _, _, _, _ := newTestAuthService()
	ctx := context.Background()

	cache.On("Get", ctx, "123456").Return("user@email.com", nil)
	hasher.On("HashPassword", "newpass123").Return("newhash", nil)
	userProv.On("UpdatePasswordByEmail", ctx, "user@email.com", "newhash").Return(nil)

	err := svc.ResetPassword(ctx, &ResetPasswordRequest{Code: "123456", Password: "newpass123"})

	require.NoError(t, err)
}

func TestResetPassword_InvalidOTP(t *testing.T) {
	svc, _, _, _, cache, _, _, _, _ := newTestAuthService()
	ctx := context.Background()

	cache.On("Get", ctx, "000000").Return("", errors.New("not found"))

	err := svc.ResetPassword(ctx, &ResetPasswordRequest{Code: "000000", Password: "newpass123"})

	assert.Error(t, err)
	assert.Equal(t, "invalid OTP", err.Error())
}
