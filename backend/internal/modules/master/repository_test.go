package master

import (
	"testing"

	"basekarya-backend/internal/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupMasterTestDB(t *testing.T) *testutil.TestDB {
	t.Helper()
	tdb := testutil.NewTestDB(
		&Department{},
		&Shift{},
		&LeaveType{},
	)
	t.Cleanup(tdb.Close)
	return tdb
}

func seedMasterTestData(t *testing.T, db *testutil.TestDB) {
	t.Helper()
	companyID := uint(1)

	require.NoError(t, db.DB.Create(&Department{ID: 1, Name: "Engineering", CompanyID: companyID}).Error)
	require.NoError(t, db.DB.Create(&Shift{ID: 1, Name: "Day", StartTime: "09:00", EndTime: "17:00", CompanyID: companyID}).Error)
	require.NoError(t, db.DB.Create(&LeaveType{ID: 1, Name: "Annual", DefaultQuota: 12, IsDeducted: true, CompanyID: companyID}).Error)
}

func TestRepoMaster_FindAllDepartments(t *testing.T) {
	tdb := setupMasterTestDB(t)
	repo := NewRepository(tdb.DB)
	ctx := testutil.CtxWithTenant(1, 1, false)

	seedMasterTestData(t, tdb)

	deps, err := repo.FindAllDepartments(ctx)

	require.NoError(t, err)
	assert.Len(t, deps, 1)
	assert.Equal(t, "Engineering", deps[0].Name)
}

func TestRepoMaster_FindAllShifts(t *testing.T) {
	tdb := setupMasterTestDB(t)
	repo := NewRepository(tdb.DB)
	ctx := testutil.CtxWithTenant(1, 1, false)

	seedMasterTestData(t, tdb)

	shifts, err := repo.FindAllShifts(ctx)

	require.NoError(t, err)
	assert.Len(t, shifts, 1)
	assert.Equal(t, "Day", shifts[0].Name)
}

func TestRepoMaster_FindAllLeaveTypes(t *testing.T) {
	tdb := setupMasterTestDB(t)
	repo := NewRepository(tdb.DB)
	ctx := testutil.CtxWithTenant(1, 1, false)

	seedMasterTestData(t, tdb)

	types, err := repo.FindAllLeaveTypes(ctx)

	require.NoError(t, err)
	assert.Len(t, types, 1)
	assert.Equal(t, "Annual", types[0].Name)
}

func TestRepoMaster_FindDepartmentByName(t *testing.T) {
	tdb := setupMasterTestDB(t)
	repo := NewRepository(tdb.DB)
	ctx := testutil.CtxWithTenant(1, 1, false)

	seedMasterTestData(t, tdb)

	tests := []struct {
		name    string
		find    string
		wantErr bool
	}{
		{name: "success", find: "Engineering", wantErr: false},
		{name: "not found", find: "Unknown", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dep, err := repo.FindDepartmentByName(ctx, tt.find)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.find, dep.Name)
			}
		})
	}
}

func TestRepoMaster_FindShiftByName(t *testing.T) {
	tdb := setupMasterTestDB(t)
	repo := NewRepository(tdb.DB)
	ctx := testutil.CtxWithTenant(1, 1, false)

	seedMasterTestData(t, tdb)

	tests := []struct {
		name    string
		find    string
		wantErr bool
	}{
		{name: "success", find: "Day", wantErr: false},
		{name: "not found", find: "Night", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			shift, err := repo.FindShiftByName(ctx, tt.find)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.find, shift.Name)
			}
		})
	}
}

func TestRepoMaster_SeedDefaults(t *testing.T) {
	tdb := setupMasterTestDB(t)
	repo := NewRepository(tdb.DB)
	ctx := testutil.CtxWithTenant(1, 1, false)

	err := repo.SeedDefaults(ctx, 1)
	require.NoError(t, err)

	deps, _ := repo.FindAllDepartments(ctx)
	assert.Len(t, deps, 1)
	assert.Equal(t, "Umum", deps[0].Name)

	shifts, _ := repo.FindAllShifts(ctx)
	assert.Len(t, shifts, 1)
	assert.Equal(t, "Regular", shifts[0].Name)

	types, _ := repo.FindAllLeaveTypes(ctx)
	assert.Len(t, types, 3)
}
