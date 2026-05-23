package rbac

import (
	"context"
	"time"

	"basekarya-backend/internal/infrastructure"
	"basekarya-backend/internal/testutil"

	"github.com/stretchr/testify/mock"
)

type mockRepo struct{ mock.Mock }

func (m *mockRepo) Create(ctx context.Context, role *Role) error {
	return m.Called(ctx, role).Error(0)
}

func (m *mockRepo) FindRoleByID(ctx context.Context, id uint) (*Role, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Role), args.Error(1)
}

func (m *mockRepo) FindRoleByName(ctx context.Context, name string) (*Role, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Role), args.Error(1)
}

func (m *mockRepo) ReplacingRolePermissions(ctx context.Context, roleID uint, permissionIDs []uint, companyID uint) error {
	return m.Called(ctx, roleID, permissionIDs, companyID).Error(0)
}

func (m *mockRepo) FindPermissionsByIDs(ctx context.Context, ids []uint) ([]Permission, error) {
	args := m.Called(ctx, ids)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]Permission), args.Error(1)
}

func (m *mockRepo) FindAllPermissions(ctx context.Context) ([]Permission, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]Permission), args.Error(1)
}

func (m *mockRepo) FindAllRoles(ctx context.Context) ([]Role, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]Role), args.Error(1)
}

func (m *mockRepo) FindAllPermissionIDs(ctx context.Context) ([]uint, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]uint), args.Error(1)
}

func (m *mockRepo) FindPermissionIDsByGroupNames(ctx context.Context, groupNames []string) ([]uint, error) {
	args := m.Called(ctx, groupNames)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]uint), args.Error(1)
}

func (m *mockRepo) FindAllPermissionsByGroupNames(ctx context.Context, groupNames []string) ([]Permission, error) {
	args := m.Called(ctx, groupNames)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]Permission), args.Error(1)
}

func (m *mockRepo) AssignPermissions(ctx context.Context, roleID uint, permissionIDs []uint, companyID uint) error {
	return m.Called(ctx, roleID, permissionIDs, companyID).Error(0)
}

func (m *mockRepo) FindRolesByCompanyID(ctx context.Context, companyID uint) ([]Role, error) {
	args := m.Called(ctx, companyID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]Role), args.Error(1)
}

func (m *mockRepo) FindRoleIDsByCompanyID(ctx context.Context, companyID uint) ([]uint, error) {
	args := m.Called(ctx, companyID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]uint), args.Error(1)
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

type mockPlanProvider struct{ mock.Mock }

func (m *mockPlanProvider) FindModulesByCompanyID(ctx context.Context, companyID uint) ([]string, error) {
	args := m.Called(ctx, companyID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]string), args.Error(1)
}

type mockService struct{ mock.Mock }

func (m *mockService) CreateRole(ctx context.Context, req *CreateRoleRequest) error {
	return m.Called(ctx, req).Error(0)
}

func (m *mockService) GetRolePermissions(ctx context.Context, roleID uint) (*RolePermissionsResponse, error) {
	args := m.Called(ctx, roleID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*RolePermissionsResponse), args.Error(1)
}

func (m *mockService) AssignPermissions(ctx context.Context, roleID uint, req *AssignPermissionsRequest) error {
	return m.Called(ctx, roleID, req).Error(0)
}

func (m *mockService) GetAllPermissions(ctx context.Context) ([]PermissionResponse, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]PermissionResponse), args.Error(1)
}

func (m *mockService) GetAllRoles(ctx context.Context) ([]RoleResponse, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]RoleResponse), args.Error(1)
}

func newTestRBACService() (Service, *mockRepo, *mockCacheProvider, *mockPlanProvider, infrastructure.TransactionManager) {
	repo := new(mockRepo)
	cache := new(mockCacheProvider)
	plan := new(mockPlanProvider)
	tm := testutil.NewMockTransactionManager()

	return NewService(repo, cache, plan, tm), repo, cache, plan, tm
}
