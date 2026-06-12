package reimbursement

import (
	"database/sql"
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

func setupReimbursementTestDB(t *testing.T) *testutil.TestDB {
	t.Helper()

	tdb := testutil.NewTestDB(
		&rbac.Role{},
		&department.Department{},
		&master.Shift{},
		&user.User{},
		&user.Employee{},
	)

	err := tdb.DB.Exec(`CREATE TABLE IF NOT EXISTS reimbursements (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		created_at DATETIME,
		updated_at DATETIME,
		deleted_at DATETIME,
		user_id INTEGER NOT NULL,
		company_id INTEGER NOT NULL,
		approved_by INTEGER,
		title VARCHAR(255) NOT NULL,
		description TEXT,
		amount REAL NOT NULL,
		date_of_expense DATE NOT NULL,
		proof_file_url VARCHAR(255) NOT NULL,
		status TEXT DEFAULT 'PENDING',
		rejection_reason TEXT
	)`).Error
	require.NoError(t, err)

	t.Cleanup(tdb.Close)
	return tdb
}

func seedReimbursementTestData(t *testing.T, db *testutil.TestDB) {
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

func TestRepo_Create(t *testing.T) {
	tdb := setupReimbursementTestDB(t)
	repo := NewRepository(tdb.DB)
	ctx := testutil.CtxWithTenant(1, 1, false)

	seedReimbursementTestData(t, tdb)

	tests := []struct {
		name           string
		reimbursement  *Reimbursement
		wantErr        bool
	}{
		{
			name: "success",
			reimbursement: &Reimbursement{
				CompanyID:     1,
				UserID:        1,
				Title:         "Office Supplies",
				Description:   "Purchased office supplies",
				Amount:        50000,
				DateOfExpense: time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC),
				ProofFileURL:  "https://storage.example.com/receipt.jpg",
				Status:        constants.ReimbursementStatusPending,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.Create(ctx, tt.reimbursement)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.NotZero(t, tt.reimbursement.ID)
			}
		})
	}
}

func TestRepo_FindByID(t *testing.T) {
	tdb := setupReimbursementTestDB(t)
	repo := NewRepository(tdb.DB)
	ctx := testutil.CtxWithTenant(1, 1, false)

	seedReimbursementTestData(t, tdb)

	r := &Reimbursement{
		CompanyID: 1, UserID: 1,
		Title: "Test", Amount: 100000,
		DateOfExpense: time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC),
		ProofFileURL: "https://storage.example.com/file.jpg",
		Status:       constants.ReimbursementStatusPending,
	}
	require.NoError(t, repo.Create(ctx, r))

	tests := []struct {
		name    string
		id      uint
		wantErr bool
	}{
		{name: "success", id: r.ID, wantErr: false},
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
				assert.Equal(t, "john", found.User.Username)
			}
		})
	}
}

func TestRepo_FindAll(t *testing.T) {
	tdb := setupReimbursementTestDB(t)
	repo := NewRepository(tdb.DB)
	ctx := testutil.CtxWithTenant(1, 1, false)

	seedReimbursementTestData(t, tdb)

	for i := 0; i < 3; i++ {
		require.NoError(t, repo.Create(ctx, &Reimbursement{
			CompanyID: 1, UserID: 1,
			Title:         "Test Reimbursement",
			Amount:        50000,
			DateOfExpense: time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC),
			ProofFileURL:  "https://storage.example.com/file.jpg",
			Status:        constants.ReimbursementStatusPending,
		}))
	}

	require.NoError(t, repo.Create(ctx, &Reimbursement{
		CompanyID: 1, UserID: 1,
		Title:         "Approved Reimbursement",
		Amount:        75000,
		DateOfExpense: time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC),
		ProofFileURL:  "https://storage.example.com/file2.jpg",
		Status:        constants.ReimbursementStatusApproved,
	}))

	tests := []struct {
		name      string
		filter    ReimbursementFilter
		wantCount int
		wantTotal int64
		wantErr   bool
	}{
		{
			name:      "all reimbursements",
			filter:    ReimbursementFilter{Page: 1, Limit: 10},
			wantCount: 4,
			wantTotal: 4,
			wantErr:   false,
		},
		{
			name:      "filter by status pending",
			filter:    ReimbursementFilter{Status: "PENDING", Page: 1, Limit: 10},
			wantCount: 3,
			wantTotal: 3,
			wantErr:   false,
		},
		{
			name:      "filter by user id",
			filter:    ReimbursementFilter{UserID: 1, Page: 1, Limit: 10},
			wantCount: 4,
			wantTotal: 4,
			wantErr:   false,
		},
		{
			name:      "paginated",
			filter:    ReimbursementFilter{Page: 1, Limit: 2},
			wantCount: 2,
			wantTotal: 4,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, total, err := repo.FindAll(ctx, tt.filter)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Len(t, results, tt.wantCount)
				assert.Equal(t, tt.wantTotal, total)
			}
		})
	}
}

func TestRepo_Update(t *testing.T) {
	tdb := setupReimbursementTestDB(t)
	repo := NewRepository(tdb.DB)
	ctx := testutil.CtxWithTenant(1, 1, false)

	seedReimbursementTestData(t, tdb)

	r := &Reimbursement{
		CompanyID: 1, UserID: 1,
		Title:         "Test",
		Amount:        50000,
		DateOfExpense: time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC),
		ProofFileURL:  "https://storage.example.com/file.jpg",
		Status:        constants.ReimbursementStatusPending,
	}
	require.NoError(t, repo.Create(ctx, r))

	tests := []struct {
		name    string
		status  constants.ReimbursementStatus
		wantErr bool
	}{
		{name: "approve success", status: constants.ReimbursementStatusApproved, wantErr: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r.Status = tt.status
			adminID := uint(10)
			r.ApprovedBy = &adminID
			err := repo.Update(ctx, r)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				found, err := repo.FindByID(ctx, r.ID)
				require.NoError(t, err)
				assert.Equal(t, tt.status, found.Status)
			}
		})
	}
}

func TestRepo_Update_Reject(t *testing.T) {
	tdb := setupReimbursementTestDB(t)
	repo := NewRepository(tdb.DB)
	ctx := testutil.CtxWithTenant(1, 1, false)

	seedReimbursementTestData(t, tdb)

	r := &Reimbursement{
		CompanyID: 1, UserID: 1,
		Title:         "Test",
		Amount:        50000,
		DateOfExpense: time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC),
		ProofFileURL:  "https://storage.example.com/file.jpg",
		Status:        constants.ReimbursementStatusPending,
	}
	require.NoError(t, repo.Create(ctx, r))

	r.Status = constants.ReimbursementStatusRejected
	r.RejectionReason = sql.NullString{String: "Not eligible", Valid: true}
	require.NoError(t, repo.Update(ctx, r))

	found, err := repo.FindByID(ctx, r.ID)
	require.NoError(t, err)
	assert.Equal(t, constants.ReimbursementStatusRejected, found.Status)
	assert.Equal(t, "Not eligible", found.RejectionReason.String)
}
