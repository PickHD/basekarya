package rbac

import (
	"testing"

	"basekarya-backend/internal/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupRBACTestDB(t *testing.T) *testutil.TestDB {
	t.Helper()
	tdb := testutil.NewTestDB(
		&PermissionGroup{},
		&Permission{},
		&Role{},
		&RolePermission{},
	)
	t.Cleanup(tdb.Close)
	return tdb
}

func seedRBACTestData(t *testing.T, db *testutil.TestDB) {
	t.Helper()
	companyID := uint(1)

	pg := &PermissionGroup{ID: 1, Name: "employee"}
	require.NoError(t, db.DB.Create(pg).Error)

	p1 := &Permission{ID: 1, Name: "VIEW_EMPLOYEE", DisplayName: "View Employee", PermissionGroupID: 1}
	require.NoError(t, db.DB.Create(p1).Error)

	p2 := &Permission{ID: 2, Name: "CREATE_EMPLOYEE", DisplayName: "Create Employee", PermissionGroupID: 1}
	require.NoError(t, db.DB.Create(p2).Error)

	role := &Role{ID: 1, Name: "SUPERADMIN", CompanyID: companyID}
	require.NoError(t, db.DB.Create(role).Error)

	rp1 := &RolePermission{RoleID: 1, PermissionID: 1, CompanyID: companyID}
	require.NoError(t, db.DB.Create(rp1).Error)
}

func TestRepoRBAC_Create(t *testing.T) {
	tdb := setupRBACTestDB(t)
	repo := NewRepository(tdb.DB)
	ctx := testutil.CtxWithTenant(1, 1, false)

	tests := []struct {
		name    string
		role    *Role
		wantErr bool
	}{
		{
			name:    "success",
			role:    &Role{Name: "MANAGER", CompanyID: 1},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.Create(ctx, tt.role)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.NotZero(t, tt.role.ID)
			}
		})
	}
}

func TestRepoRBAC_FindRoleByID(t *testing.T) {
	tdb := setupRBACTestDB(t)
	repo := NewRepository(tdb.DB)
	ctx := testutil.CtxWithTenant(1, 1, false)

	seedRBACTestData(t, tdb)

	tests := []struct {
		name    string
		id      uint
		wantErr bool
	}{
		{name: "success", id: 1, wantErr: false},
		{name: "not found", id: 99, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			role, err := repo.FindRoleByID(ctx, tt.id)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.id, role.ID)
				assert.Equal(t, "SUPERADMIN", role.Name)
				assert.Len(t, role.Permissions, 1)
			}
		})
	}
}

func TestRepoRBAC_FindRoleByName(t *testing.T) {
	tdb := setupRBACTestDB(t)
	repo := NewRepository(tdb.DB)
	ctx := testutil.CtxWithTenant(1, 1, false)

	seedRBACTestData(t, tdb)

	tests := []struct {
		name     string
		roleName string
		wantErr  bool
	}{
		{name: "success", roleName: "SUPERADMIN", wantErr: false},
		{name: "not found", roleName: "NONEXISTENT", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			role, err := repo.FindRoleByName(ctx, tt.roleName)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.roleName, role.Name)
			}
		})
	}
}

func TestRepoRBAC_FindAllRoles(t *testing.T) {
	tdb := setupRBACTestDB(t)
	repo := NewRepository(tdb.DB)
	ctx := testutil.CtxWithTenant(1, 1, false)

	seedRBACTestData(t, tdb)

	roles, err := repo.FindAllRoles(ctx)

	require.NoError(t, err)
	assert.Len(t, roles, 1)
	assert.Equal(t, "SUPERADMIN", roles[0].Name)
}

func TestRepoRBAC_FindAllPermissions(t *testing.T) {
	tdb := setupRBACTestDB(t)
	repo := NewRepository(tdb.DB)
	ctx := testutil.CtxWithTenant(1, 1, false)

	seedRBACTestData(t, tdb)

	perms, err := repo.FindAllPermissions(ctx)

	require.NoError(t, err)
	assert.Len(t, perms, 2)
}

func TestRepoRBAC_FindPermissionsByIDs(t *testing.T) {
	tdb := setupRBACTestDB(t)
	repo := NewRepository(tdb.DB)
	ctx := testutil.CtxWithTenant(1, 1, false)

	seedRBACTestData(t, tdb)

	tests := []struct {
		name    string
		ids     []uint
		wantLen int
		wantErr bool
	}{
		{name: "success", ids: []uint{1, 2}, wantLen: 2, wantErr: false},
		{name: "partial match", ids: []uint{1, 99}, wantLen: 1, wantErr: false},
		{name: "empty ids", ids: []uint{}, wantLen: 0, wantErr: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			perms, err := repo.FindPermissionsByIDs(ctx, tt.ids)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Len(t, perms, tt.wantLen)
			}
		})
	}
}

func TestRepoRBAC_ReplacingRolePermissions(t *testing.T) {
	tdb := setupRBACTestDB(t)
	repo := NewRepository(tdb.DB)
	ctx := testutil.CtxWithTenant(1, 1, false)

	seedRBACTestData(t, tdb)

	err := repo.ReplacingRolePermissions(ctx, 1, []uint{1, 2}, 1)
	require.NoError(t, err)

	role, err := repo.FindRoleByID(ctx, 1)
	require.NoError(t, err)
	assert.Len(t, role.Permissions, 2)
}

func TestRepoRBAC_FindAllPermissionIDs(t *testing.T) {
	tdb := setupRBACTestDB(t)
	repo := NewRepository(tdb.DB)
	ctx := testutil.CtxWithTenant(1, 1, false)

	seedRBACTestData(t, tdb)

	ids, err := repo.FindAllPermissionIDs(ctx)
	require.NoError(t, err)
	assert.Len(t, ids, 2)
}

func TestRepoRBAC_AssignPermissions(t *testing.T) {
	tdb := setupRBACTestDB(t)
	repo := NewRepository(tdb.DB)
	ctx := testutil.CtxWithTenant(1, 1, false)

	seedRBACTestData(t, tdb)

	err := repo.AssignPermissions(ctx, 1, []uint{1, 2}, 1)
	require.NoError(t, err)

	role, err := repo.FindRoleByID(ctx, 1)
	require.NoError(t, err)
	assert.Len(t, role.Permissions, 2)
}

func TestRepoRBAC_FindRolesByCompanyID(t *testing.T) {
	tdb := setupRBACTestDB(t)
	repo := NewRepository(tdb.DB)
	ctx := testutil.CtxWithTenant(1, 1, false)

	seedRBACTestData(t, tdb)

	roles, err := repo.FindRolesByCompanyID(ctx, 1)
	require.NoError(t, err)
	assert.Len(t, roles, 1)
}

func TestRepoRBAC_FindRoleIDsByCompanyID(t *testing.T) {
	tdb := setupRBACTestDB(t)
	repo := NewRepository(tdb.DB)
	ctx := testutil.CtxWithTenant(1, 1, false)

	seedRBACTestData(t, tdb)

	ids, err := repo.FindRoleIDsByCompanyID(ctx, 1)
	require.NoError(t, err)
	assert.Len(t, ids, 1)
}
