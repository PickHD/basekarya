package overtime

import (
	"testing"

	"basekarya-backend/internal/modules/master"
	"basekarya-backend/internal/modules/rbac"
	"basekarya-backend/internal/modules/user"
	"basekarya-backend/pkg/constants"
	"basekarya-backend/internal/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupOvertimeTestDB(t *testing.T) *testutil.TestDB {
	t.Helper()

	tdb := testutil.NewTestDB(
		&rbac.Role{},
		&master.Department{},
		&master.Shift{},
		&user.User{},
		&user.Employee{},
	)

	err := tdb.DB.Exec(`CREATE TABLE IF NOT EXISTS overtimes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		created_at DATETIME,
		updated_at DATETIME,
		deleted_at DATETIME,
		user_id INTEGER NOT NULL,
		company_id INTEGER NOT NULL,
		employee_id INTEGER NOT NULL,
		approved_by INTEGER,
		date DATE NOT NULL,
		start_time TIME NOT NULL,
		end_time TIME NOT NULL,
		duration_minutes INTEGER NOT NULL,
		reason TEXT,
		status TEXT DEFAULT 'PENDING',
		rejection_reason TEXT
	)`).Error
	require.NoError(t, err)

	t.Cleanup(tdb.Close)
	return tdb
}

func seedOvertimeTestData(t *testing.T, db *testutil.TestDB) {
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

func TestRepoOT_Create(t *testing.T) {
	tdb := setupOvertimeTestDB(t)
	repo := NewRepository(tdb.DB)
	ctx := testutil.CtxWithTenant(1, 1, false)

	seedOvertimeTestData(t, tdb)

	tests := []struct {
		name    string
		ot      *Overtime
		wantErr bool
	}{
		{
			name: "success",
			ot: &Overtime{
				CompanyID:       1,
				UserID:          1,
				EmployeeID:      1,
				Date:            "2026-06-01",
				StartTime:       "18:00",
				EndTime:         "20:00",
				DurationMinutes: 120,
				Reason:          "Project deadline",
				Status:          constants.OvertimeStatusPending,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.Create(ctx, tt.ot)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.NotZero(t, tt.ot.ID)
			}
		})
	}
}

func TestRepoOT_FindByID(t *testing.T) {
	tdb := setupOvertimeTestDB(t)
	repo := NewRepository(tdb.DB)
	ctx := testutil.CtxWithTenant(1, 1, false)

	seedOvertimeTestData(t, tdb)

	ot := &Overtime{
		CompanyID: 1, UserID: 1, EmployeeID: 1,
		Date: "2026-06-01", StartTime: "18:00", EndTime: "20:00",
		DurationMinutes: 120, Reason: "Test", Status: constants.OvertimeStatusPending,
	}
	require.NoError(t, repo.Create(ctx, ot))

	tests := []struct {
		name    string
		id      uint
		wantErr bool
	}{
		{name: "success", id: ot.ID, wantErr: false},
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

func TestRepoOT_FindAll(t *testing.T) {
	tdb := setupOvertimeTestDB(t)
	repo := NewRepository(tdb.DB)
	ctx := testutil.CtxWithTenant(1, 1, false)

	seedOvertimeTestData(t, tdb)

	for i := 0; i < 3; i++ {
		require.NoError(t, repo.Create(ctx, &Overtime{
			CompanyID: 1, UserID: 1, EmployeeID: 1,
			Date: "2026-06-01", StartTime: "18:00", EndTime: "20:00",
			DurationMinutes: 120, Reason: "Test", Status: constants.OvertimeStatusPending,
		}))
	}

	tests := []struct {
		name      string
		filter    OvertimeFilter
		wantCount int
		wantErr   bool
	}{
		{name: "all", filter: OvertimeFilter{Page: 1, Limit: 10}, wantCount: 3, wantErr: false},
		{name: "paginated", filter: OvertimeFilter{Page: 1, Limit: 2}, wantCount: 2, wantErr: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, total, err := repo.FindAll(ctx, tt.filter)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Len(t, results, tt.wantCount)
				assert.Equal(t, int64(3), total)
			}
		})
	}
}

func TestRepoOT_Update(t *testing.T) {
	tdb := setupOvertimeTestDB(t)
	repo := NewRepository(tdb.DB)
	ctx := testutil.CtxWithTenant(1, 1, false)

	seedOvertimeTestData(t, tdb)

	ot := &Overtime{
		CompanyID: 1, UserID: 1, EmployeeID: 1,
		Date: "2026-06-01", StartTime: "18:00", EndTime: "20:00",
		DurationMinutes: 120, Reason: "Test", Status: constants.OvertimeStatusPending,
	}
	require.NoError(t, repo.Create(ctx, ot))

	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "success", wantErr: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ot.Status = constants.OvertimeStatusApproved
			err := repo.Update(ctx, ot)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				var found Overtime
				tdb.DB.First(&found, ot.ID)
				assert.Equal(t, constants.OvertimeStatusApproved, found.Status)
			}
		})
	}
}
