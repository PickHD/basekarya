package loan

import (
	"database/sql"
	"testing"

	"basekarya-backend/internal/modules/master"
	"basekarya-backend/internal/modules/rbac"
	"basekarya-backend/internal/modules/user"
	"basekarya-backend/internal/testutil"
	"basekarya-backend/pkg/constants"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func setupLoanTestDB(t *testing.T) *testutil.TestDB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	require.NoError(t, err)

	require.NoError(t, db.AutoMigrate(
		&rbac.Role{},
		&master.Department{},
		&master.Shift{},
		&user.User{},
		&user.Employee{},
	))

	require.NoError(t, db.Exec(`CREATE TABLE IF NOT EXISTS loans (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		created_at DATETIME,
		updated_at DATETIME,
		user_id INTEGER NOT NULL,
		company_id INTEGER NOT NULL,
		employee_id INTEGER NOT NULL,
		approved_by INTEGER,
		total_amount REAL NOT NULL,
		installment_amount REAL NOT NULL,
		remaining_amount REAL NOT NULL,
		reason TEXT,
		status TEXT DEFAULT 'PENDING',
		rejection_reason TEXT
	)`).Error)

	t.Cleanup(func() {
		sqlDB, _ := db.DB()
		if sqlDB != nil {
			sqlDB.Close()
		}
	})

	return &testutil.TestDB{DB: db}
}

func seedLoanTestData(t *testing.T, db *testutil.TestDB) {
	t.Helper()
	companyID := uint(1)

	role := &rbac.Role{ID: 1, Name: "Admin", CompanyID: companyID}
	require.NoError(t, db.DB.Create(role).Error)

	dept := &master.Department{ID: 1, Name: "Engineering", CompanyID: companyID}
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

func TestRepo_Create(t *testing.T) {
	tdb := setupLoanTestDB(t)
	repo := NewRepository(tdb.DB)
	ctx := testutil.CtxWithTenant(1, 1, false)

	seedLoanTestData(t, tdb)

	tests := []struct {
		name    string
		loan    *Loan
		wantErr bool
	}{
		{
			name: "success",
			loan: &Loan{
				CompanyID:         1,
				UserID:            1,
				EmployeeID:        1,
				TotalAmount:       5000000,
				InstallmentAmount: 500000,
				RemainingAmount:   5000000,
				Status:            constants.LoanStatusPending,
				Reason:            "Emergency",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.Create(ctx, tt.loan)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.NotZero(t, tt.loan.ID)
			}
		})
	}
}

func TestRepo_FindByID(t *testing.T) {
	tdb := setupLoanTestDB(t)
	repo := NewRepository(tdb.DB)
	ctx := testutil.CtxWithTenant(1, 1, false)

	seedLoanTestData(t, tdb)

	loan := &Loan{
		CompanyID: 1, UserID: 1, EmployeeID: 1,
		TotalAmount: 5000000, InstallmentAmount: 500000, RemainingAmount: 5000000,
		Status: constants.LoanStatusPending, Reason: "Emergency",
	}
	require.NoError(t, repo.Create(ctx, loan))

	tests := []struct {
		name    string
		id      uint
		wantErr bool
	}{
		{name: "success", id: loan.ID, wantErr: false},
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
				assert.Equal(t, "John Doe", found.Employee.FullName)
			}
		})
	}
}

func TestRepo_FindActiveLoanByUserID(t *testing.T) {
	tdb := setupLoanTestDB(t)
	repo := NewRepository(tdb.DB)
	ctx := testutil.CtxWithTenant(1, 1, false)

	seedLoanTestData(t, tdb)

	loan := &Loan{
		CompanyID: 1, UserID: 1, EmployeeID: 1,
		TotalAmount: 5000000, InstallmentAmount: 500000, RemainingAmount: 5000000,
		Status: constants.LoanStatusApproved, Reason: "Emergency",
	}
	require.NoError(t, repo.Create(ctx, loan))

	tests := []struct {
		name     string
		userID   uint
		wantErr  bool
	}{
		{name: "success", userID: 1, wantErr: false},
		{name: "not found", userID: 99, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			found, err := repo.FindActiveLoanByUserID(ctx, tt.userID)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, loan.ID, found.ID)
			}
		})
	}
}

func TestRepo_FindAll(t *testing.T) {
	tdb := setupLoanTestDB(t)
	repo := NewRepository(tdb.DB)
	ctx := testutil.CtxWithTenant(1, 1, false)

	seedLoanTestData(t, tdb)

	for i := 0; i < 3; i++ {
		require.NoError(t, repo.Create(ctx, &Loan{
			CompanyID: 1, UserID: 1, EmployeeID: 1,
			TotalAmount: float64(i+1) * 1000000, InstallmentAmount: 500000, RemainingAmount: float64(i+1) * 1000000,
			Status: constants.LoanStatusPending, Reason: "Test",
		}))
	}

	tests := []struct {
		name      string
		filter    LoanFilter
		wantCount int
		wantTotal int64
		wantErr   bool
	}{
		{name: "all loans", filter: LoanFilter{Page: 1, Limit: 10}, wantCount: 3, wantTotal: 3, wantErr: false},
		{name: "paginated", filter: LoanFilter{Page: 1, Limit: 2}, wantCount: 2, wantTotal: 3, wantErr: false},
		{name: "filter by status", filter: LoanFilter{Page: 1, Limit: 10, Status: "PENDING"}, wantCount: 3, wantTotal: 3, wantErr: false},
		{name: "filter by user", filter: LoanFilter{Page: 1, Limit: 10, UserID: 1}, wantCount: 3, wantTotal: 3, wantErr: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			loans, total, err := repo.FindAll(ctx, tt.filter)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Len(t, loans, tt.wantCount)
				assert.Equal(t, tt.wantTotal, total)
			}
		})
	}
}

func TestRepo_Update(t *testing.T) {
	tdb := setupLoanTestDB(t)
	repo := NewRepository(tdb.DB)
	ctx := testutil.CtxWithTenant(1, 1, false)

	seedLoanTestData(t, tdb)

	loan := &Loan{
		CompanyID: 1, UserID: 1, EmployeeID: 1,
		TotalAmount: 5000000, InstallmentAmount: 500000, RemainingAmount: 5000000,
		Status: constants.LoanStatusPending, Reason: "Emergency",
	}
	require.NoError(t, repo.Create(ctx, loan))

	tests := []struct {
		name    string
		status  constants.LoanStatus
		wantErr bool
	}{
		{name: "approve loan", status: constants.LoanStatusApproved, wantErr: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			loan.Status = tt.status
			err := repo.Update(ctx, loan)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				found, _ := repo.FindByID(ctx, loan.ID)
				assert.Equal(t, tt.status, found.Status)
			}
		})
	}
}

func TestRepo_GetBulkActiveLoansByEmployeeIds(t *testing.T) {
	tdb := setupLoanTestDB(t)
	repo := NewRepository(tdb.DB)
	ctx := testutil.CtxWithTenant(1, 1, false)

	seedLoanTestData(t, tdb)

	loan := &Loan{
		CompanyID: 1, UserID: 1, EmployeeID: 1,
		TotalAmount: 5000000, InstallmentAmount: 500000, RemainingAmount: 5000000,
		Status: constants.LoanStatusApproved, Reason: "Emergency",
	}
	require.NoError(t, repo.Create(ctx, loan))

	tests := []struct {
		name      string
		ids       []uint
		wantCount int
		wantErr   bool
	}{
		{name: "found", ids: []uint{1}, wantCount: 1, wantErr: false},
		{name: "not found", ids: []uint{99}, wantCount: 0, wantErr: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := repo.GetBulkActiveLoansByEmployeeIds(ctx, tt.ids)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Len(t, result, tt.wantCount)
			}
		})
	}
}

func TestRepo_FindActiveLoanByUserID_Rejected(t *testing.T) {
	tdb := setupLoanTestDB(t)
	repo := NewRepository(tdb.DB)
	ctx := testutil.CtxWithTenant(1, 1, false)

	seedLoanTestData(t, tdb)

	loan := &Loan{
		CompanyID: 1, UserID: 1, EmployeeID: 1,
		TotalAmount: 5000000, InstallmentAmount: 500000, RemainingAmount: 5000000,
		Status: constants.LoanStatusRejected, Reason: "Test",
		RejectionReason: sql.NullString{String: "Not eligible", Valid: true},
	}
	require.NoError(t, repo.Create(ctx, loan))

	_, err := repo.FindActiveLoanByUserID(ctx, 1)
	require.Error(t, err)
}

func TestRepo_FindActiveLoanByUserID_PaidOff(t *testing.T) {
	tdb := setupLoanTestDB(t)
	repo := NewRepository(tdb.DB)
	ctx := testutil.CtxWithTenant(1, 1, false)

	seedLoanTestData(t, tdb)

	loan := &Loan{
		CompanyID: 1, UserID: 1, EmployeeID: 1,
		TotalAmount: 5000000, InstallmentAmount: 500000, RemainingAmount: 0,
		Status: constants.LoanStatusPaidOff, Reason: "Test",
	}
	require.NoError(t, repo.Create(ctx, loan))

	_, err := repo.FindActiveLoanByUserID(ctx, 1)
	require.Error(t, err)
}

func TestRepo_Update_WithRejection(t *testing.T) {
	tdb := setupLoanTestDB(t)
	repo := NewRepository(tdb.DB)
	ctx := testutil.CtxWithTenant(1, 1, false)

	seedLoanTestData(t, tdb)

	loan := &Loan{
		CompanyID: 1, UserID: 1, EmployeeID: 1,
		TotalAmount: 5000000, InstallmentAmount: 500000, RemainingAmount: 5000000,
		Status: constants.LoanStatusPending, Reason: "Emergency",
	}
	require.NoError(t, repo.Create(ctx, loan))

	loan.Status = constants.LoanStatusRejected
	loan.RejectionReason = sql.NullString{String: "Not eligible", Valid: true}
	require.NoError(t, repo.Update(ctx, loan))

	found, err := repo.FindByID(ctx, loan.ID)
	require.NoError(t, err)
	assert.Equal(t, constants.LoanStatusRejected, found.Status)
	assert.Equal(t, "Not eligible", found.RejectionReason.String)
}

func TestRepo_FindAll_EmptyStatus(t *testing.T) {
	tdb := setupLoanTestDB(t)
	repo := NewRepository(tdb.DB)
	ctx := testutil.CtxWithTenant(1, 1, false)

	seedLoanTestData(t, tdb)

	loans, total, err := repo.FindAll(ctx, LoanFilter{Page: 1, Limit: 10, Status: "APPROVED"})
	require.NoError(t, err)
	assert.Len(t, loans, 0)
	assert.Equal(t, int64(0), total)
}

func TestRepo_FindAll_SecondPage(t *testing.T) {
	tdb := setupLoanTestDB(t)
	repo := NewRepository(tdb.DB)
	ctx := testutil.CtxWithTenant(1, 1, false)

	seedLoanTestData(t, tdb)

	for i := 0; i < 5; i++ {
		require.NoError(t, repo.Create(ctx, &Loan{
			CompanyID: 1, UserID: 1, EmployeeID: 1,
			TotalAmount: 1000000, InstallmentAmount: 100000, RemainingAmount: 1000000,
			Status: constants.LoanStatusPending, Reason: "Test",
		}))
	}

	loans, total, err := repo.FindAll(ctx, LoanFilter{Page: 2, Limit: 2})
	require.NoError(t, err)
	assert.Len(t, loans, 2)
	assert.Equal(t, int64(5), total)
}
