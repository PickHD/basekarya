package user

import (
	"context"
	"mime/multipart"
	"time"

	"basekarya-backend/internal/modules/rbac"
	"basekarya-backend/internal/testutil"
	"basekarya-backend/pkg/response"

	"github.com/stretchr/testify/mock"
)

type mockRepo struct{ mock.Mock }

func (m *mockRepo) FindByUsername(ctx context.Context, username string) (*User, error) {
	args := m.Called(ctx, username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*User), args.Error(1)
}

func (m *mockRepo) FindByID(ctx context.Context, id uint) (*User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*User), args.Error(1)
}

func (m *mockRepo) UpdateEmployee(ctx context.Context, emp *Employee) error {
	return m.Called(ctx, emp).Error(0)
}

func (m *mockRepo) UpdateUser(ctx context.Context, user *User) error {
	return m.Called(ctx, user).Error(0)
}

func (m *mockRepo) FindAllEmployees(ctx context.Context, page, limit int, search string) ([]User, int64, error) {
	args := m.Called(ctx, page, limit, search)
	return args.Get(0).([]User), args.Get(1).(int64), args.Error(2)
}

func (m *mockRepo) CreateUser(ctx context.Context, user *User) error {
	return m.Called(ctx, user).Error(0)
}

func (m *mockRepo) CreateEmployee(ctx context.Context, emp *Employee) error {
	return m.Called(ctx, emp).Error(0)
}

func (m *mockRepo) DeleteUser(ctx context.Context, id uint) error {
	return m.Called(ctx, id).Error(0)
}

func (m *mockRepo) FindEmployeeByID(ctx context.Context, id uint) (*Employee, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Employee), args.Error(1)
}

func (m *mockRepo) FindEmployeeByEmail(ctx context.Context, email string) (*Employee, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Employee), args.Error(1)
}

func (m *mockRepo) UpdatePasswordByEmail(ctx context.Context, email string, password string) error {
	return m.Called(ctx, email, password).Error(0)
}

func (m *mockRepo) CountActiveEmployee(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

func (m *mockRepo) FindAllEmployeeActive(ctx context.Context) ([]Employee, error) {
	args := m.Called(ctx)
	return args.Get(0).([]Employee), args.Error(1)
}

func (m *mockRepo) FindApprovalUsers(ctx context.Context, permissionApprovalName string) ([]uint, error) {
	args := m.Called(ctx, permissionApprovalName)
	return args.Get(0).([]uint), args.Error(1)
}

func (m *mockRepo) FindRoleByID(ctx context.Context, id uint) (*rbac.Role, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*rbac.Role), args.Error(1)
}

func (m *mockRepo) FindAllUserIDs(ctx context.Context) ([]uint, error) {
	args := m.Called(ctx)
	return args.Get(0).([]uint), args.Error(1)
}

func (m *mockRepo) ForceResetPasswordByCompanyID(ctx context.Context, companyID uint) error {
	return m.Called(ctx, companyID).Error(0)
}

type mockHasher struct{ mock.Mock }

func (m *mockHasher) HashPassword(password string) (string, error) {
	args := m.Called(password)
	return args.String(0), args.Error(1)
}

func (m *mockHasher) CheckPasswordHash(password, hash string) bool {
	args := m.Called(password, hash)
	return args.Bool(0)
}

type mockStorage struct{ mock.Mock }

func (m *mockStorage) UploadFileMultipart(ctx context.Context, file *multipart.FileHeader, objectName string) (string, error) {
	args := m.Called(ctx, file, objectName)
	return args.String(0), args.Error(1)
}

type mockCache struct{ mock.Mock }

func (m *mockCache) Get(ctx context.Context, key string) (string, error) {
	args := m.Called(ctx, key)
	return args.String(0), args.Error(1)
}

func (m *mockCache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	args := m.Called(ctx, key, value, expiration)
	return args.Error(0)
}

func (m *mockCache) Del(ctx context.Context, key string) error {
	args := m.Called(ctx, key)
	return args.Error(0)
}

type mockLeaveGen struct{ mock.Mock }

func (m *mockLeaveGen) GenerateInitialBalance(ctx context.Context, employeeID uint) error {
	return m.Called(ctx, employeeID).Error(0)
}

type mockSubscription struct{ mock.Mock }

func (m *mockSubscription) CheckEmployeeLimit(ctx context.Context) (bool, error) {
	args := m.Called(ctx)
	return args.Bool(0), args.Error(1)
}

type mockService struct{ mock.Mock }

func (m *mockService) GetProfile(userID uint) (*UserProfileResponse, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*UserProfileResponse), args.Error(1)
}

func (m *mockService) UpdateProfile(ctx context.Context, userID uint, req *UpdateProfileRequest, file *multipart.FileHeader) error {
	return m.Called(ctx, userID, req, file).Error(0)
}

func (m *mockService) ChangePassword(ctx context.Context, userID uint, req *ChangePasswordRequest) error {
	return m.Called(ctx, userID, req).Error(0)
}

func (m *mockService) GetAllEmployees(ctx context.Context, page, limit int, search string) ([]EmployeeListResponse, *response.Meta, error) {
	args := m.Called(ctx, page, limit, search)
	var meta *response.Meta
	if args.Get(1) != nil {
		meta = args.Get(1).(*response.Meta)
	}
	return args.Get(0).([]EmployeeListResponse), meta, args.Error(2)
}

func (m *mockService) CreateEmployee(ctx context.Context, req *CreateEmployeeRequest) (*CreateEmployeeResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*CreateEmployeeResponse), args.Error(1)
}

func (m *mockService) UpdateEmployee(ctx context.Context, id uint, req *UpdateEmployeeRequest) error {
	return m.Called(ctx, id, req).Error(0)
}

func (m *mockService) DeleteEmployee(ctx context.Context, id uint) error {
	return m.Called(ctx, id).Error(0)
}

func (m *mockService) FindApprovalUsers(ctx context.Context, permissionApprovalName string) ([]uint, error) {
	args := m.Called(ctx, permissionApprovalName)
	return args.Get(0).([]uint), args.Error(1)
}

func newTestUserService() (Service, *mockRepo, *mockHasher, *mockStorage, *mockCache, *mockLeaveGen, *testutil.MockTransactionManager, *mockSubscription) {
	repo := new(mockRepo)
	hasher := new(mockHasher)
	storage := new(mockStorage)
	cache := new(mockCache)
	leaveGen := new(mockLeaveGen)
	tm := testutil.NewMockTransactionManager()
	sub := new(mockSubscription)

	svc := NewService(repo, hasher, storage, cache, leaveGen, tm, sub)
	return svc, repo, hasher, storage, cache, leaveGen, tm, sub
}
