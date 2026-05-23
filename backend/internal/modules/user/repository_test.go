package user

import (
	"testing"

	"basekarya-backend/internal/modules/master"
	"basekarya-backend/internal/modules/rbac"
	"basekarya-backend/internal/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupUserTestDB(t *testing.T) *testutil.TestDB {
	t.Helper()
	tdb := testutil.NewTestDB(
		&rbac.Role{},
		&master.Department{},
		&master.Shift{},
		&User{},
		&Employee{},
	)
	t.Cleanup(tdb.Close)
	return tdb
}

func seedUserTestData(t *testing.T, db *testutil.TestDB) {
	t.Helper()
	companyID := uint(1)

	role := &rbac.Role{ID: 1, Name: "EMPLOYEE", CompanyID: companyID}
	require.NoError(t, db.DB.Create(role).Error)

	dept := &master.Department{ID: 1, Name: "Engineering", CompanyID: companyID}
	require.NoError(t, db.DB.Create(dept).Error)

	shift := &master.Shift{ID: 1, Name: "Day", StartTime: "09:00", EndTime: "17:00", CompanyID: companyID}
	require.NoError(t, db.DB.Create(shift).Error)

	usr := &User{ID: 1, Username: "john.doe", PasswordHash: "hashedpassword", RoleID: 1, CompanyID: companyID, IsActive: true}
	require.NoError(t, db.DB.Create(usr).Error)

	emp := &Employee{
		ID: 1, UserID: 1, CompanyID: companyID, DepartmentID: 1, ShiftID: 1,
		NIK: "EMP001", FullName: "John Doe", Email: "john@example.com",
		Position: "Developer", BaseSalary: 5000000,
	}
	require.NoError(t, db.DB.Create(emp).Error)
}

func TestRepo_FindByID(t *testing.T) {
	tdb := setupUserTestDB(t)
	repo := NewRepository(tdb.DB)
	ctx := testutil.CtxWithTenant(1, 1, false)

	seedUserTestData(t, tdb)

	tests := []struct {
		name    string
		id      uint
		wantErr bool
	}{
		{name: "success", id: 1, wantErr: false},
		{name: "not found", id: 999, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := repo.FindByID(ctx, tt.id)
			if tt.wantErr {
				require.Error(t, err)
				assert.Nil(t, user)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.id, user.ID)
				assert.Equal(t, "john.doe", user.Username)
				assert.NotNil(t, user.Employee)
				assert.Equal(t, "John Doe", user.Employee.FullName)
				assert.NotNil(t, user.Role)
				assert.Equal(t, "EMPLOYEE", user.Role.Name)
			}
		})
	}
}

func TestRepo_FindByUsername(t *testing.T) {
	tdb := setupUserTestDB(t)
	repo := NewRepository(tdb.DB)
	ctx := testutil.CtxWithTenant(1, 1, false)

	seedUserTestData(t, tdb)

	tests := []struct {
		name     string
		username string
		wantErr  bool
	}{
		{name: "success", username: "john.doe", wantErr: false},
		{name: "not found", username: "nonexistent", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := repo.FindByUsername(ctx, tt.username)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.username, user.Username)
			}
		})
	}
}

func TestRepo_CreateUser(t *testing.T) {
	tdb := setupUserTestDB(t)
	repo := NewRepository(tdb.DB)
	ctx := testutil.CtxWithTenant(1, 1, false)

	seedUserTestData(t, tdb)

	tests := []struct {
		name    string
		user    *User
		wantErr bool
	}{
		{
			name:    "success",
			user:    &User{Username: "jane.doe", PasswordHash: "hash", RoleID: 1, CompanyID: 1, MustChangePassword: true},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.CreateUser(ctx, tt.user)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.NotZero(t, tt.user.ID)
			}
		})
	}
}

func TestRepo_CreateEmployee(t *testing.T) {
	tdb := setupUserTestDB(t)
	repo := NewRepository(tdb.DB)
	ctx := testutil.CtxWithTenant(1, 1, false)

	seedUserTestData(t, tdb)

	usr := &User{Username: "jane.doe", PasswordHash: "hash", RoleID: 1, CompanyID: 1}
	require.NoError(t, tdb.DB.Create(usr).Error)

	tests := []struct {
		name    string
		emp     *Employee
		wantErr bool
	}{
		{
			name: "success",
			emp: &Employee{
				UserID: usr.ID, CompanyID: 1, DepartmentID: 1, ShiftID: 1,
				NIK: "EMP002", FullName: "Jane Doe", Email: "jane@example.com", Position: "Designer",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.CreateEmployee(ctx, tt.emp)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.NotZero(t, tt.emp.ID)
			}
		})
	}
}

func TestRepo_UpdateEmployee(t *testing.T) {
	tdb := setupUserTestDB(t)
	repo := NewRepository(tdb.DB)
	ctx := testutil.CtxWithTenant(1, 1, false)

	seedUserTestData(t, tdb)

	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "success", wantErr: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			emp, err := repo.FindEmployeeByID(ctx, 1)
			require.NoError(t, err)

			emp.FullName = "John Updated"
			err = repo.UpdateEmployee(ctx, emp)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				found, _ := repo.FindEmployeeByID(ctx, 1)
				assert.Equal(t, "John Updated", found.FullName)
			}
		})
	}
}

func TestRepo_UpdateUser(t *testing.T) {
	tdb := setupUserTestDB(t)
	repo := NewRepository(tdb.DB)
	ctx := testutil.CtxWithTenant(1, 1, false)

	seedUserTestData(t, tdb)

	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "success", wantErr: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := repo.FindByID(ctx, 1)
			require.NoError(t, err)

			user.MustChangePassword = false
			err = repo.UpdateUser(ctx, user)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				found, _ := repo.FindByID(ctx, 1)
				assert.False(t, found.MustChangePassword)
			}
		})
	}
}

func TestRepo_DeleteUser(t *testing.T) {
	tdb := setupUserTestDB(t)
	repo := NewRepository(tdb.DB)
	ctx := testutil.CtxWithTenant(1, 1, false)

	seedUserTestData(t, tdb)

	tests := []struct {
		name    string
		id      uint
		wantErr bool
	}{
		{name: "success", id: 1, wantErr: false},
		{name: "not found", id: 999, wantErr: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.DeleteUser(ctx, tt.id)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestRepo_FindEmployeeByID(t *testing.T) {
	tdb := setupUserTestDB(t)
	repo := NewRepository(tdb.DB)
	ctx := testutil.CtxWithTenant(1, 1, false)

	seedUserTestData(t, tdb)

	tests := []struct {
		name    string
		id      uint
		wantErr bool
	}{
		{name: "success", id: 1, wantErr: false},
		{name: "not found", id: 999, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			emp, err := repo.FindEmployeeByID(ctx, tt.id)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.id, emp.ID)
				assert.Equal(t, "John Doe", emp.FullName)
			}
		})
	}
}

func TestRepo_FindAllEmployees(t *testing.T) {
	tdb := setupUserTestDB(t)
	repo := NewRepository(tdb.DB)
	ctx := testutil.CtxWithTenant(1, 1, false)

	seedUserTestData(t, tdb)

	tests := []struct {
		name      string
		page      int
		limit     int
		search    string
		wantCount int
		wantErr   bool
	}{
		{name: "all employees", page: 1, limit: 10, search: "", wantCount: 1, wantErr: false},
		{name: "search by name", page: 1, limit: 10, search: "John", wantCount: 1, wantErr: false},
		{name: "search no match", page: 1, limit: 10, search: "Nobody", wantCount: 0, wantErr: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			users, total, err := repo.FindAllEmployees(ctx, tt.page, tt.limit, tt.search)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Len(t, users, tt.wantCount)
				if tt.wantCount > 0 {
					assert.Equal(t, int64(1), total)
				}
			}
		})
	}
}

func TestRepo_FindRoleByID(t *testing.T) {
	tdb := setupUserTestDB(t)
	repo := NewRepository(tdb.DB)
	ctx := testutil.CtxWithTenant(1, 1, false)

	seedUserTestData(t, tdb)

	tests := []struct {
		name    string
		id      uint
		wantErr bool
	}{
		{name: "success", id: 1, wantErr: false},
		{name: "not found", id: 999, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			role, err := repo.FindRoleByID(ctx, tt.id)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.id, role.ID)
				assert.Equal(t, "EMPLOYEE", role.Name)
			}
		})
	}
}

func TestRepo_FindAllUserIDs(t *testing.T) {
	tdb := setupUserTestDB(t)
	repo := NewRepository(tdb.DB)
	ctx := testutil.CtxWithTenant(1, 1, false)

	seedUserTestData(t, tdb)

	tests := []struct {
		name      string
		wantCount int
		wantErr   bool
	}{
		{name: "success", wantCount: 1, wantErr: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ids, err := repo.FindAllUserIDs(ctx)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Len(t, ids, tt.wantCount)
			}
		})
	}
}
