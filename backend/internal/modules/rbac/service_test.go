package rbac

import (
	"errors"
	"testing"

	"basekarya-backend/internal/testutil"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestService_CreateRole(t *testing.T) {
	ctx := testutil.CtxWithTenant(1, 1, false)

	tests := []struct {
		name       string
		req        *CreateRoleRequest
		setupMocks func(*mockRepo, *mockCacheProvider)
		wantErr    bool
		errMsg     string
	}{
		{
			name: "success",
			req:  &CreateRoleRequest{Name: "MANAGER"},
			setupMocks: func(repo *mockRepo, cache *mockCacheProvider) {
				repo.On("FindRoleByName", mock.Anything, "MANAGER").Return(nil, errors.New("not found"))
				repo.On("Create", mock.Anything, mock.AnythingOfType("*rbac.Role")).Return(nil)
				cache.On("Del", mock.Anything, mock.AnythingOfType("string")).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "role already exists",
			req:  &CreateRoleRequest{Name: "SUPERADMIN"},
			setupMocks: func(repo *mockRepo, cache *mockCacheProvider) {
				repo.On("FindRoleByName", mock.Anything, "SUPERADMIN").Return(&Role{ID: 1, Name: "SUPERADMIN"}, nil)
			},
			wantErr: true,
			errMsg:  "role already exists",
		},
		{
			name: "create fails",
			req:  &CreateRoleRequest{Name: "MANAGER"},
			setupMocks: func(repo *mockRepo, cache *mockCacheProvider) {
				repo.On("FindRoleByName", mock.Anything, "MANAGER").Return(nil, errors.New("not found"))
				repo.On("Create", mock.Anything, mock.AnythingOfType("*rbac.Role")).Return(errors.New("db error"))
			},
			wantErr: true,
			errMsg:  "failed to create role",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, repo, cache, _, _ := newTestRBACService()
			tt.setupMocks(repo, cache)

			err := svc.CreateRole(ctx, tt.req)

			if tt.wantErr {
				require.Error(t, err)
				assert.Equal(t, tt.errMsg, err.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestService_GetRolePermissions(t *testing.T) {
	ctx := testutil.CtxWithTenant(1, 1, false)

	tests := []struct {
		name       string
		roleID     uint
		setupMocks func(*mockRepo, *mockCacheProvider)
		wantErr    bool
	}{
		{
			name: "from db on cache miss",
			roleID: 1,
			setupMocks: func(repo *mockRepo, cache *mockCacheProvider) {
				cache.On("Get", mock.Anything, mock.AnythingOfType("string")).Return("", redis.Nil)
				repo.On("FindRoleByID", mock.Anything, uint(1)).Return(&Role{
					ID:          1,
					Name:        "SUPERADMIN",
					Permissions: []Permission{{ID: 1, Name: "VIEW_EMPLOYEE"}},
				}, nil)
				cache.On("Set", mock.Anything, mock.AnythingOfType("string"), mock.Anything, mock.Anything).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "role not found",
			roleID: 99,
			setupMocks: func(repo *mockRepo, cache *mockCacheProvider) {
				cache.On("Get", mock.Anything, mock.AnythingOfType("string")).Return("", redis.Nil)
				repo.On("FindRoleByID", mock.Anything, uint(99)).Return(nil, errors.New("not found"))
			},
			wantErr: true,
		},
		{
			name:   "from cache",
			roleID: 1,
			setupMocks: func(repo *mockRepo, cache *mockCacheProvider) {
				cache.On("Get", mock.Anything, mock.AnythingOfType("string")).Return(`{"role_id":1,"role_name":"SUPERADMIN","permissions":[]}`, nil)
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, repo, cache, _, _ := newTestRBACService()
			tt.setupMocks(repo, cache)

			resp, err := svc.GetRolePermissions(ctx, tt.roleID)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, resp)
			}
		})
	}
}

func TestService_AssignPermissions(t *testing.T) {
	ctx := testutil.CtxWithTenant(1, 1, false)

	tests := []struct {
		name       string
		roleID     uint
		req        *AssignPermissionsRequest
		setupMocks func(*mockRepo, *mockCacheProvider)
		wantErr    bool
		errMsg     string
	}{
		{
			name:   "success",
			roleID: 1,
			req:    &AssignPermissionsRequest{PermissionIDs: []uint{1, 2}},
			setupMocks: func(repo *mockRepo, cache *mockCacheProvider) {
				repo.On("FindRoleByID", mock.Anything, uint(1)).Return(&Role{ID: 1, Name: "SUPERADMIN"}, nil)
				repo.On("FindPermissionsByIDs", mock.Anything, []uint{1, 2}).Return([]Permission{
					{ID: 1, Name: "VIEW_EMPLOYEE"},
					{ID: 2, Name: "CREATE_EMPLOYEE"},
				}, nil)
				repo.On("ReplacingRolePermissions", mock.Anything, uint(1), []uint{1, 2}, uint(1)).Return(nil)
				cache.On("Del", mock.Anything, mock.AnythingOfType("string")).Return(nil)
			},
			wantErr: false,
		},
		{
			name:   "role not found",
			roleID: 99,
			req:    &AssignPermissionsRequest{PermissionIDs: []uint{1}},
			setupMocks: func(repo *mockRepo, cache *mockCacheProvider) {
				repo.On("FindRoleByID", mock.Anything, uint(99)).Return(nil, errors.New("not found"))
			},
			wantErr: true,
			errMsg:  "role not found",
		},
		{
			name:   "invalid permission ids",
			roleID: 1,
			req:    &AssignPermissionsRequest{PermissionIDs: []uint{1, 99}},
			setupMocks: func(repo *mockRepo, cache *mockCacheProvider) {
				repo.On("FindRoleByID", mock.Anything, uint(1)).Return(&Role{ID: 1}, nil)
				repo.On("FindPermissionsByIDs", mock.Anything, []uint{1, 99}).Return([]Permission{
					{ID: 1, Name: "VIEW_EMPLOYEE"},
				}, nil)
			},
			wantErr: true,
			errMsg:  "one or more permissions are invalid",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, repo, cache, _, _ := newTestRBACService()
			tt.setupMocks(repo, cache)

			err := svc.AssignPermissions(ctx, tt.roleID, tt.req)

			if tt.wantErr {
				require.Error(t, err)
				assert.Equal(t, tt.errMsg, err.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestService_GetAllPermissions(t *testing.T) {
	ctx := testutil.CtxWithTenant(1, 1, false)

	tests := []struct {
		name       string
		setupMocks func(*mockRepo, *mockCacheProvider, *mockPlanProvider)
		wantErr    bool
	}{
		{
			name: "from db on cache miss",
			setupMocks: func(repo *mockRepo, cache *mockCacheProvider, plan *mockPlanProvider) {
				plan.On("FindModulesByCompanyID", mock.Anything, uint(1)).Return([]string{"employee_management"}, nil)
				cache.On("Get", mock.Anything, mock.AnythingOfType("string")).Return("", redis.Nil)
				repo.On("FindAllPermissionsByGroupNames", mock.Anything, mock.Anything).Return([]Permission{
					{ID: 1, Name: "VIEW_EMPLOYEE", DisplayName: "View", PermissionGroup: PermissionGroup{ID: 1, Name: "employee"}},
				}, nil)
				cache.On("Set", mock.Anything, mock.AnythingOfType("string"), mock.Anything, mock.Anything).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "from cache",
			setupMocks: func(repo *mockRepo, cache *mockCacheProvider, plan *mockPlanProvider) {
				plan.On("FindModulesByCompanyID", mock.Anything, uint(1)).Return([]string{"employee_management"}, nil)
				cache.On("Get", mock.Anything, mock.AnythingOfType("string")).Return(`[{"group":{"id":1,"name":"employee"},"permissions":[{"id":1,"name":"VIEW_EMPLOYEE","display_name":"View","description":""}]}]`, nil)
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, repo, cache, plan, _ := newTestRBACService()
			tt.setupMocks(repo, cache, plan)

			_, err := svc.GetAllPermissions(ctx)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestService_GetAllRoles(t *testing.T) {
	ctx := testutil.CtxWithTenant(1, 1, false)

	tests := []struct {
		name       string
		setupMocks func(*mockRepo, *mockCacheProvider)
		wantErr    bool
	}{
		{
			name: "from db on cache miss",
			setupMocks: func(repo *mockRepo, cache *mockCacheProvider) {
				cache.On("Get", mock.Anything, mock.AnythingOfType("string")).Return("", redis.Nil)
				repo.On("FindAllRoles", mock.Anything).Return([]Role{
					{ID: 1, Name: "SUPERADMIN"},
					{ID: 2, Name: "EMPLOYEE"},
				}, nil)
				cache.On("Set", mock.Anything, mock.AnythingOfType("string"), mock.Anything, mock.Anything).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "from cache",
			setupMocks: func(repo *mockRepo, cache *mockCacheProvider) {
				cache.On("Get", mock.Anything, mock.AnythingOfType("string")).Return(`[{"id":1,"name":"SUPERADMIN"}]`, nil)
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, repo, cache, _, _ := newTestRBACService()
			tt.setupMocks(repo, cache)

			_, err := svc.GetAllRoles(ctx)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
