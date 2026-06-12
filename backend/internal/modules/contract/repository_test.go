package contract

import (
	"testing"
	"time"

	"basekarya-backend/internal/modules/department"
	"basekarya-backend/internal/modules/master"
	"basekarya-backend/internal/modules/rbac"
	"basekarya-backend/internal/modules/user"
	"basekarya-backend/internal/testutil"
	"basekarya-backend/pkg/constants"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupContractTestDB(t *testing.T) *testutil.TestDB {
	t.Helper()
	tdb := testutil.NewTestDB(
		&rbac.Role{},
		&department.Department{},
		&master.Shift{},
		&user.User{},
		&user.Employee{},
		&Contract{},
	)
	t.Cleanup(tdb.Close)
	return tdb
}

func seedContractTestData(t *testing.T, db *testutil.TestDB) {
	t.Helper()
	companyID := uint(1)

	role := &rbac.Role{ID: 1, Name: "Admin", CompanyID: companyID}
	require.NoError(t, db.DB.Create(role).Error)

	dept := &department.Department{ID: 1, Name: "Engineering", CompanyID: companyID}
	require.NoError(t, db.DB.Create(dept).Error)

	shift := &master.Shift{ID: 1, Name: "Day", StartTime: "09:00", EndTime: "17:00", CompanyID: companyID}
	require.NoError(t, db.DB.Create(shift).Error)

	usr := &user.User{ID: 1, Username: "john", RoleID: 1, CompanyID: companyID}
	require.NoError(t, db.DB.Create(usr).Error)

	emp := &user.Employee{
		ID: 1, UserID: 1, CompanyID: companyID, DepartmentID: 1, ShiftID: 1,
		NIK: "EMP001", FullName: "John Doe",
	}
	require.NoError(t, db.DB.Create(emp).Error)
}

func TestRepo_Upsert(t *testing.T) {
	tdb := setupContractTestDB(t)
	repo := NewRepository(tdb.DB)
	ctx := testutil.CtxWithTenant(1, 1, false)

	seedContractTestData(t, tdb)

	startDate := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2026, 12, 31, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name    string
		contract *Contract
		wantErr bool
	}{
		{
			name: "success PKWT",
			contract: &Contract{
				CompanyID:      1,
				EmployeeID:     1,
				ContractType:   constants.ContractTypePKWT,
				ContractNumber: "CTR-001",
				StartDate:      startDate,
				EndDate:        &endDate,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.Upsert(ctx, tt.contract)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.NotZero(t, tt.contract.ID)
			}
		})
	}
}

func TestRepo_Upsert_UpdateExisting(t *testing.T) {
	tdb := setupContractTestDB(t)
	repo := NewRepository(tdb.DB)
	ctx := testutil.CtxWithTenant(1, 1, false)

	seedContractTestData(t, tdb)

	startDate := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2026, 12, 31, 0, 0, 0, 0, time.UTC)

	c := &Contract{
		CompanyID:      1,
		EmployeeID:     1,
		ContractType:   constants.ContractTypePKWT,
		ContractNumber: "CTR-001",
		StartDate:      startDate,
		EndDate:        &endDate,
	}
	require.NoError(t, repo.Upsert(ctx, c))
	originalID := c.ID

	newEndDate := time.Date(2027, 6, 30, 0, 0, 0, 0, time.UTC)
	c2 := &Contract{
		CompanyID:      1,
		EmployeeID:     1,
		ContractType:   constants.ContractTypePKWT,
		ContractNumber: "CTR-002",
		StartDate:      startDate,
		EndDate:        &newEndDate,
	}
	require.NoError(t, repo.Upsert(ctx, c2))
	assert.Equal(t, originalID, c2.ID)

	found, err := repo.FindByEmployeeID(ctx, 1)
	require.NoError(t, err)
	assert.Equal(t, "CTR-002", found.ContractNumber)
}

func TestRepo_FindByID(t *testing.T) {
	tdb := setupContractTestDB(t)
	repo := NewRepository(tdb.DB)
	ctx := testutil.CtxWithTenant(1, 1, false)

	seedContractTestData(t, tdb)

	startDate := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2026, 12, 31, 0, 0, 0, 0, time.UTC)
	c := &Contract{
		CompanyID: 1, EmployeeID: 1, ContractType: constants.ContractTypePKWT,
		ContractNumber: "CTR-001", StartDate: startDate, EndDate: &endDate,
	}
	require.NoError(t, repo.Upsert(ctx, c))

	tests := []struct {
		name    string
		id      uint
		wantErr bool
	}{
		{name: "success", id: c.ID, wantErr: false},
		{name: "not found", id: 999, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			found, err := repo.FindByID(ctx, tt.id)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.id, found.ID)
				assert.NotNil(t, found.Employee)
				assert.Equal(t, "John Doe", found.Employee.FullName)
			}
		})
	}
}

func TestRepo_FindByEmployeeID(t *testing.T) {
	tdb := setupContractTestDB(t)
	repo := NewRepository(tdb.DB)
	ctx := testutil.CtxWithTenant(1, 1, false)

	seedContractTestData(t, tdb)

	startDate := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2026, 12, 31, 0, 0, 0, 0, time.UTC)
	require.NoError(t, repo.Upsert(ctx, &Contract{
		CompanyID: 1, EmployeeID: 1, ContractType: constants.ContractTypePKWT,
		ContractNumber: "CTR-001", StartDate: startDate, EndDate: &endDate,
	}))

	tests := []struct {
		name       string
		employeeID uint
		wantErr    bool
	}{
		{name: "success", employeeID: 1, wantErr: false},
		{name: "not found", employeeID: 99, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			found, err := repo.FindByEmployeeID(ctx, tt.employeeID)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.employeeID, found.EmployeeID)
			}
		})
	}
}

func TestRepo_FindAll(t *testing.T) {
	tdb := setupContractTestDB(t)
	repo := NewRepository(tdb.DB)
	ctx := testutil.CtxWithTenant(1, 1, false)

	seedContractTestData(t, tdb)

	startDate := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2026, 12, 31, 0, 0, 0, 0, time.UTC)
	require.NoError(t, repo.Upsert(ctx, &Contract{
		CompanyID: 1, EmployeeID: 1, ContractType: constants.ContractTypePKWT,
		ContractNumber: "CTR-001", StartDate: startDate, EndDate: &endDate,
	}))

	tests := []struct {
		name      string
		filter    *ContractFilter
		wantCount int
		wantTotal int64
		wantErr   bool
	}{
		{
			name:      "all contracts",
			filter:    &ContractFilter{Page: 1, Limit: 10},
			wantCount: 1,
			wantTotal: 1,
			wantErr:   false,
		},
		{
			name:      "filter by contract type PKWT",
			filter:    &ContractFilter{ContractType: "PKWT", Page: 1, Limit: 10},
			wantCount: 1,
			wantTotal: 1,
			wantErr:   false,
		},
		{
			name:      "filter by contract type PKWTT empty",
			filter:    &ContractFilter{ContractType: "PKWTT", Page: 1, Limit: 10},
			wantCount: 0,
			wantTotal: 0,
			wantErr:   false,
		},
		{
			name:      "paginated",
			filter:    &ContractFilter{Page: 1, Limit: 1},
			wantCount: 1,
			wantTotal: 1,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			contracts, total, err := repo.FindAll(ctx, tt.filter)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Len(t, contracts, tt.wantCount)
				assert.Equal(t, tt.wantTotal, total)
			}
		})
	}
}

func TestRepo_SoftDelete(t *testing.T) {
	tdb := setupContractTestDB(t)
	repo := NewRepository(tdb.DB)
	ctx := testutil.CtxWithTenant(1, 1, false)

	seedContractTestData(t, tdb)

	startDate := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2026, 12, 31, 0, 0, 0, 0, time.UTC)
	c := &Contract{
		CompanyID: 1, EmployeeID: 1, ContractType: constants.ContractTypePKWT,
		ContractNumber: "CTR-001", StartDate: startDate, EndDate: &endDate,
	}
	require.NoError(t, repo.Upsert(ctx, c))

	err := repo.SoftDelete(ctx, c.ID)
	require.NoError(t, err)

	_, err = repo.FindByID(ctx, c.ID)
	require.Error(t, err)
}
