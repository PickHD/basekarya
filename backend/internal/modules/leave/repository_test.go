package leave

import (
	"testing"
	"time"

	"basekarya-backend/internal/modules/department"
	"basekarya-backend/internal/modules/master"
	"basekarya-backend/internal/modules/rbac"
	"basekarya-backend/internal/modules/user"
	"basekarya-backend/pkg/constants"
	"basekarya-backend/internal/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupLeaveTestDB(t *testing.T) *testutil.TestDB {
	t.Helper()
	// Note: attendance.Attendance uses MySQL enum which is incompatible with SQLite.
	// Only include entities needed for leave repo methods being tested.
	tdb := testutil.NewTestDB(
		&rbac.Role{},
		&department.Department{},
		&master.Shift{},
		&master.LeaveType{},
		&user.User{},
		&user.Employee{},
		&LeaveBalance{},
		&LeaveRequest{},
	)
	t.Cleanup(tdb.Close)
	return tdb
}

func seedLeaveTestData(t *testing.T, db *testutil.TestDB) {
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

	lt := &master.LeaveType{ID: 1, Name: "Annual", DefaultQuota: 12, IsDeducted: true, CompanyID: companyID}
	require.NoError(t, db.DB.Create(lt).Error)
}

func TestRepo_CreateRequest(t *testing.T) {
	tdb := setupLeaveTestDB(t)
	repo := NewRepository(tdb.DB)
	ctx := testutil.CtxWithTenant(1, 1, false)

	seedLeaveTestData(t, tdb)

	tests := []struct {
		name    string
		req     *LeaveRequest
		wantErr bool
	}{
		{
			name: "success",
			req: &LeaveRequest{
				CompanyID:   1,
				UserID:      1,
				EmployeeID:  1,
				LeaveTypeID: 1,
				StartDate:   time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC),
				EndDate:     time.Date(2026, 6, 2, 0, 0, 0, 0, time.UTC),
				TotalDays:   2,
				Reason:      "Family event",
				Status:      constants.LeaveStatusPending,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.CreateRequest(ctx, tt.req)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.NotZero(t, tt.req.ID)
			}
		})
	}
}

func TestRepo_FindRequestByID(t *testing.T) {
	tdb := setupLeaveTestDB(t)
	repo := NewRepository(tdb.DB)
	ctx := testutil.CtxWithTenant(1, 1, false)

	seedLeaveTestData(t, tdb)

	// Create a request first
	req := &LeaveRequest{
		CompanyID: 1, UserID: 1, EmployeeID: 1, LeaveTypeID: 1,
		StartDate: time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2026, 6, 2, 0, 0, 0, 0, time.UTC),
		TotalDays: 2, Reason: "Test", Status: constants.LeaveStatusPending,
	}
	require.NoError(t, repo.CreateRequest(ctx, req))

	tests := []struct {
		name    string
		id      uint
		wantErr bool
	}{
		{name: "success", id: req.ID, wantErr: false},
		{name: "not found", id: 999, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			found, err := repo.FindRequestByID(ctx, tt.id)
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

func TestRepo_GetBalance(t *testing.T) {
	tdb := setupLeaveTestDB(t)
	repo := NewRepository(tdb.DB)
	ctx := testutil.CtxWithTenant(1, 1, false)

	seedLeaveTestData(t, tdb)

	// Seed a balance
	balance := &LeaveBalance{
		CompanyID: 1, EmployeeID: 1, LeaveTypeID: 1,
		Year: 2026, QuotaTotal: 12, QuotaUsed: 2, QuotaLeft: 10,
	}
	require.NoError(t, tdb.DB.Create(balance).Error)

	tests := []struct {
		name         string
		employeeID   uint
		leaveTypeID  uint
		year         int
		wantQuotaLeft int
		wantErr      bool
	}{
		{name: "success", employeeID: 1, leaveTypeID: 1, year: 2026, wantQuotaLeft: 10, wantErr: false},
		{name: "not found", employeeID: 99, leaveTypeID: 1, year: 2026, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bal, err := repo.GetBalance(ctx, tt.employeeID, tt.leaveTypeID, tt.year)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.wantQuotaLeft, bal.QuotaLeft)
			}
		})
	}
}

func TestRepo_FindAllRequests(t *testing.T) {
	tdb := setupLeaveTestDB(t)
	repo := NewRepository(tdb.DB)
	ctx := testutil.CtxWithTenant(1, 1, false)

	seedLeaveTestData(t, tdb)

	// Seed requests
	for i := 0; i < 3; i++ {
		require.NoError(t, repo.CreateRequest(ctx, &LeaveRequest{
			CompanyID: 1, UserID: 1, EmployeeID: 1, LeaveTypeID: 1,
			StartDate: time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC),
			EndDate:   time.Date(2026, 6, 2, 0, 0, 0, 0, time.UTC),
			TotalDays: 2, Reason: "Test", Status: constants.LeaveStatusPending,
		}))
	}

	tests := []struct {
		name      string
		filter    *LeaveFilter
		wantCount int
		wantErr   bool
	}{
		{name: "all requests", filter: &LeaveFilter{Page: 1, Limit: 10}, wantCount: 3, wantErr: false},
		{name: "paginated", filter: &LeaveFilter{Page: 1, Limit: 2}, wantCount: 2, wantErr: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			requests, total, err := repo.FindAllRequests(ctx, tt.filter)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Len(t, requests, tt.wantCount)
				assert.Equal(t, int64(3), total)
			}
		})
	}
}

func TestRepo_RejectRequest(t *testing.T) {
	tdb := setupLeaveTestDB(t)
	repo := NewRepository(tdb.DB)
	ctx := testutil.CtxWithTenant(1, 1, false)

	seedLeaveTestData(t, tdb)

	req := &LeaveRequest{
		CompanyID: 1, UserID: 1, EmployeeID: 1, LeaveTypeID: 1,
		StartDate: time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2026, 6, 2, 0, 0, 0, 0, time.UTC),
		TotalDays: 2, Reason: "Test", Status: constants.LeaveStatusPending,
	}
	require.NoError(t, repo.CreateRequest(ctx, req))

	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "success", wantErr: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.RejectRequest(ctx, req.ID, 10, "Not eligible")
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				var found LeaveRequest
				tdb.DB.First(&found, req.ID)
				assert.Equal(t, constants.LeaveStatusRejected, found.Status)
			}
		})
	}
}

func TestRepo_FindAllLeaveTypes(t *testing.T) {
	tdb := setupLeaveTestDB(t)
	repo := NewRepository(tdb.DB)
	ctx := testutil.CtxWithTenant(1, 1, false)

	seedLeaveTestData(t, tdb)

	tests := []struct {
		name      string
		wantCount int
		wantErr   bool
	}{
		{name: "success", wantCount: 1, wantErr: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			types, err := repo.FindAllLeaveTypes(ctx)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Len(t, types, tt.wantCount)
			}
		})
	}
}

func TestRepo_CreateLeaveBalances(t *testing.T) {
	tdb := setupLeaveTestDB(t)
	repo := NewRepository(tdb.DB)
	ctx := testutil.CtxWithTenant(1, 1, false)

	seedLeaveTestData(t, tdb)

	balances := []LeaveBalance{
		{CompanyID: 1, EmployeeID: 1, LeaveTypeID: 1, Year: 2026, QuotaTotal: 12, QuotaUsed: 0, QuotaLeft: 12},
	}

	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "success", wantErr: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.CreateLeaveBalances(ctx, balances)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				var count int64
				tdb.DB.Model(&LeaveBalance{}).Count(&count)
				assert.Equal(t, int64(1), count)
			}
		})
	}
}
