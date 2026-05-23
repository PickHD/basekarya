package payroll

import (
	"testing"
	"time"

	"basekarya-backend/internal/modules/master"
	"basekarya-backend/internal/modules/rbac"
	"basekarya-backend/internal/modules/user"
	"basekarya-backend/internal/testutil"
	"basekarya-backend/pkg/constants"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupPayrollTestDB(t *testing.T) *testutil.TestDB {
	t.Helper()
	tdb := testutil.NewTestDB(
		&rbac.Role{},
		&master.Department{},
		&master.Shift{},
		&user.User{},
		&user.Employee{},
		&Payroll{},
		&PayrollDetail{},
	)
	t.Cleanup(tdb.Close)
	return tdb
}

func seedPayrollTestData(t *testing.T, db *testutil.TestDB) {
	t.Helper()
	companyID := uint(1)

	role := &rbac.Role{ID: 1, Name: "EMPLOYEE", CompanyID: companyID}
	require.NoError(t, db.DB.Create(role).Error)

	dept := &master.Department{ID: 1, Name: "Engineering", CompanyID: companyID}
	require.NoError(t, db.DB.Create(dept).Error)

	shift := &master.Shift{ID: 1, Name: "Day", StartTime: "09:00", EndTime: "17:00", CompanyID: companyID}
	require.NoError(t, db.DB.Create(shift).Error)

	usr := &user.User{ID: 1, Username: "john.doe", PasswordHash: "hashed", RoleID: 1, CompanyID: companyID, IsActive: true}
	require.NoError(t, db.DB.Create(usr).Error)

	emp := &user.Employee{
		ID: 1, UserID: 1, CompanyID: companyID, DepartmentID: 1, ShiftID: 1,
		NIK: "EMP001", FullName: "John Doe", Email: "john@example.com",
		Position: "Developer", BaseSalary: 5000000,
		BankName: "BCA", BankAccountNumber: "1234567890", BankAccountHolder: "John Doe",
	}
	require.NoError(t, db.DB.Create(emp).Error)

	usr2 := &user.User{ID: 2, Username: "jane.doe", PasswordHash: "hashed", RoleID: 1, CompanyID: companyID, IsActive: true}
	require.NoError(t, db.DB.Create(usr2).Error)

	emp2 := &user.Employee{
		ID: 2, UserID: 2, CompanyID: companyID, DepartmentID: 1, ShiftID: 1,
		NIK: "EMP002", FullName: "Jane Smith", Email: "jane@example.com",
		Position: "Designer", BaseSalary: 4000000,
	}
	require.NoError(t, db.DB.Create(emp2).Error)
}

func seedPayrollWithDetails(t *testing.T, db *testutil.TestDB, companyID uint) {
	t.Helper()
	periodDate := time.Date(2025, time.June, 1, 0, 0, 0, 0, time.Local)

	p := &Payroll{
		EmployeeID:     1,
		CompanyID:      companyID,
		PeriodDate:     periodDate,
		BaseSalary:     5000000,
		TotalAllowance: 5500000,
		TotalDeduction: 500000,
		NetSalary:      5000000,
		Status:         constants.PayrollStatusDraft,
		Notes:          "test payroll",
	}
	require.NoError(t, db.DB.Create(p).Error)

	details := []PayrollDetail{
		{PayrollID: p.ID, CompanyID: companyID, Title: "Base Salary", Type: constants.DetailTypeAllowance, Amount: 5000000},
		{PayrollID: p.ID, CompanyID: companyID, Title: "Reimbursement", Type: constants.DetailTypeAllowance, Amount: 500000},
		{PayrollID: p.ID, CompanyID: companyID, Title: "Potongan Terlambat (60 menit)", Type: constants.DetailTypeDeduction, Amount: 500000},
	}
	for _, d := range details {
		require.NoError(t, db.DB.Create(&d).Error)
	}
}

func TestRepo_CreateBulk(t *testing.T) {
	tdb := setupPayrollTestDB(t)
	repo := NewRepository(tdb.DB)
	ctx := testutil.CtxWithTenant(1, 1, false)

	seedPayrollTestData(t, tdb)

	periodDate := time.Date(2025, time.July, 1, 0, 0, 0, 0, time.Local)
	payrolls := []Payroll{
		{
			EmployeeID:     1,
			CompanyID:      1,
			PeriodDate:     periodDate,
			BaseSalary:     5000000,
			TotalAllowance: 5000000,
			TotalDeduction: 0,
			NetSalary:      5000000,
			Status:         constants.PayrollStatusDraft,
			Details: []PayrollDetail{
				{CompanyID: 1, Title: "Base Salary", Type: constants.DetailTypeAllowance, Amount: 5000000},
			},
		},
	}

	err := repo.CreateBulk(ctx, &payrolls)
	require.NoError(t, err)
	assert.NotZero(t, payrolls[0].ID)
}

func TestRepo_FindAll(t *testing.T) {
	tdb := setupPayrollTestDB(t)
	repo := NewRepository(tdb.DB)
	ctx := testutil.CtxWithTenant(1, 1, false)

	seedPayrollTestData(t, tdb)
	seedPayrollWithDetails(t, tdb, 1)

	tests := []struct {
		name      string
		filter    *PayrollFilter
		wantCount int
		wantErr   bool
	}{
		{
			name:      "find all no filter",
			filter:    &PayrollFilter{Page: 1, Limit: 10, Month: 0, Year: 0},
			wantCount: 1,
			wantErr:   false,
		},
		{
			name:      "filter by month and year",
			filter:    &PayrollFilter{Page: 1, Limit: 10, Month: 6, Year: 2025},
			wantCount: 1,
			wantErr:   false,
		},
		{
			name:      "filter by wrong month",
			filter:    &PayrollFilter{Page: 1, Limit: 10, Month: 1, Year: 2025},
			wantCount: 0,
			wantErr:   false,
		},
		{
			name:      "search by keyword",
			filter:    &PayrollFilter{Page: 1, Limit: 10, Month: 6, Year: 2025, Keyword: "John"},
			wantCount: 1,
			wantErr:   false,
		},
		{
			name:      "search by keyword no match",
			filter:    &PayrollFilter{Page: 1, Limit: 10, Month: 6, Year: 2025, Keyword: "Nobody"},
			wantCount: 0,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payrolls, total, err := repo.FindAll(ctx, tt.filter)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Len(t, payrolls, tt.wantCount)
				if tt.wantCount > 0 {
					assert.GreaterOrEqual(t, total, int64(1))
				}
			}
		})
	}
}

func TestRepo_FindByID(t *testing.T) {
	tdb := setupPayrollTestDB(t)
	repo := NewRepository(tdb.DB)
	ctx := testutil.CtxWithTenant(1, 1, false)

	seedPayrollTestData(t, tdb)
	seedPayrollWithDetails(t, tdb, 1)

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
			p, err := repo.FindByID(ctx, tt.id)
			if tt.wantErr {
				require.Error(t, err)
				assert.Nil(t, p)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.id, p.ID)
				assert.Equal(t, constants.PayrollStatusDraft, p.Status)
				assert.NotNil(t, p.Employee)
				assert.Equal(t, "John Doe", p.Employee.FullName)
				assert.NotEmpty(t, p.Details)
			}
		})
	}
}

func TestRepo_GetExistingEmployeeID(t *testing.T) {
	tdb := setupPayrollTestDB(t)
	repo := NewRepository(tdb.DB)
	ctx := testutil.CtxWithTenant(1, 1, false)

	seedPayrollTestData(t, tdb)
	seedPayrollWithDetails(t, tdb, 1)

	tests := []struct {
		name        string
		month       int
		year        int
		wantCount   int
		wantErr     bool
	}{
		{name: "found existing", month: 6, year: 2025, wantCount: 1, wantErr: false},
		{name: "no existing", month: 1, year: 2025, wantCount: 0, wantErr: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			existing, err := repo.GetExistingEmployeeID(ctx, tt.month, tt.year)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Len(t, existing, tt.wantCount)
			}
		})
	}
}

func TestRepo_UpdateStatus(t *testing.T) {
	tdb := setupPayrollTestDB(t)
	repo := NewRepository(tdb.DB)
	ctx := testutil.CtxWithTenant(1, 1, false)

	seedPayrollTestData(t, tdb)
	seedPayrollWithDetails(t, tdb, 1)

	err := repo.UpdateStatus(ctx, 1, constants.PayrollStatusPaid)
	require.NoError(t, err)

	p, err := repo.FindByID(ctx, 1)
	require.NoError(t, err)
	assert.Equal(t, constants.PayrollStatusPaid, p.Status)
}
