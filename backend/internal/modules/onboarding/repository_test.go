package onboarding

import (
	"testing"
	"time"

	"basekarya-backend/internal/modules/user"
	"basekarya-backend/internal/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupOnboardingTestDB(t *testing.T) *testutil.TestDB {
	t.Helper()
	tdb := testutil.NewTestDB(
		&OnboardingWorkflow{},
		&OnboardingTask{},
		&user.User{},
		&user.Employee{},
	)
	t.Cleanup(tdb.Close)
	return tdb
}

func TestRepo_CreateWorkflow(t *testing.T) {
	tdb := setupOnboardingTestDB(t)
	repo := NewRepository(tdb.DB)
	ctx := testutil.CtxWithTenant(1, 1, false)

	startDate := time.Now().Add(24 * time.Hour)

	wf := &OnboardingWorkflow{
		CompanyID:    1,
		NewHireName:  "Jane Smith",
		NewHireEmail: "jane@example.com",
		Position:     "Developer",
		Department:   "Engineering",
		StartDate:    &startDate,
		Status:       WorkflowStatusInProgress,
	}

	err := repo.CreateWorkflow(ctx, wf)
	require.NoError(t, err)
	assert.NotZero(t, wf.ID)
}

func TestRepo_CreateTasks(t *testing.T) {
	tdb := setupOnboardingTestDB(t)
	repo := NewRepository(tdb.DB)
	ctx := testutil.CtxWithTenant(1, 1, false)

	startDate := time.Now()
	wf := &OnboardingWorkflow{
		CompanyID:    1,
		NewHireName:  "Jane Smith",
		NewHireEmail: "jane@example.com",
		Position:     "Developer",
		Department:   "Engineering",
		StartDate:    &startDate,
		Status:       WorkflowStatusInProgress,
	}
	require.NoError(t, repo.CreateWorkflow(ctx, wf))

	tasks := []OnboardingTask{
		{
			CompanyID:            1,
			OnboardingWorkflowID: wf.ID,
			TaskName:             "Create Email",
			SortOrder:            1,
		},
		{
			CompanyID:            1,
			OnboardingWorkflowID: wf.ID,
			TaskName:             "Sign NDA",
			SortOrder:            2,
		},
	}

	err := repo.CreateTasks(ctx, tasks)
	require.NoError(t, err)
	assert.NotZero(t, tasks[0].ID)
	assert.NotZero(t, tasks[1].ID)
}

func TestRepo_CreateTasks_Empty(t *testing.T) {
	tdb := setupOnboardingTestDB(t)
	repo := NewRepository(tdb.DB)
	ctx := testutil.CtxWithTenant(1, 1, false)

	err := repo.CreateTasks(ctx, nil)
	require.NoError(t, err)
}

func TestRepo_FindAllWorkflows(t *testing.T) {
	tdb := setupOnboardingTestDB(t)
	repo := NewRepository(tdb.DB)
	ctx := testutil.CtxWithTenant(1, 1, false)

	startDate := time.Now()
	wf := &OnboardingWorkflow{
		CompanyID:    1,
		NewHireName:  "Jane Smith",
		NewHireEmail: "jane@example.com",
		Position:     "Developer",
		Department:   "Engineering",
		StartDate:    &startDate,
		Status:       WorkflowStatusInProgress,
	}
	require.NoError(t, repo.CreateWorkflow(ctx, wf))

	task := &OnboardingTask{
		CompanyID:            1,
		OnboardingWorkflowID: wf.ID,
		TaskName:             "Create Email",
		SortOrder:            1,
	}
	require.NoError(t, tdb.DB.Create(task).Error)

	tests := []struct {
		name       string
		filter     *WorkflowFilter
		wantCount  int
		wantTotal  int64
		wantErr    bool
	}{
		{
			name:      "all workflows",
			filter:    &WorkflowFilter{Page: 1, Limit: 10},
			wantCount: 1,
			wantTotal: 1,
			wantErr:   false,
		},
		{
			name:      "filter by status",
			filter:    &WorkflowFilter{Status: "IN_PROGRESS", Page: 1, Limit: 10},
			wantCount: 1,
			wantTotal: 1,
			wantErr:   false,
		},
		{
			name:      "filter by search",
			filter:    &WorkflowFilter{Search: "Jane", Page: 1, Limit: 10},
			wantCount: 1,
			wantTotal: 1,
			wantErr:   false,
		},
		{
			name:      "search no match",
			filter:    &WorkflowFilter{Search: "Nobody", Page: 1, Limit: 10},
			wantCount: 0,
			wantTotal: 0,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			workflows, total, err := repo.FindAllWorkflows(ctx, tt.filter)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Len(t, workflows, tt.wantCount)
				assert.Equal(t, tt.wantTotal, total)
			}
		})
	}
}

func TestRepo_FindWorkflowByID(t *testing.T) {
	tdb := setupOnboardingTestDB(t)
	repo := NewRepository(tdb.DB)
	ctx := testutil.CtxWithTenant(1, 1, false)

	usr := &user.User{ID: 1, Username: "admin", PasswordHash: "hash", RoleID: 1, CompanyID: 1, IsActive: true}
	require.NoError(t, tdb.DB.Create(usr).Error)
	emp := &user.Employee{
		ID: 1, UserID: 1, CompanyID: 1, DepartmentID: 1, ShiftID: 1,
		NIK: "EMP001", FullName: "Admin User", Email: "admin@test.com", Position: "IT Admin",
	}
	require.NoError(t, tdb.DB.Create(emp).Error)

	startDate := time.Now()
	wf := &OnboardingWorkflow{
		CompanyID:    1,
		NewHireName:  "Jane Smith",
		NewHireEmail: "jane@example.com",
		Position:     "Developer",
		Department:   "Engineering",
		StartDate:    &startDate,
		Status:       WorkflowStatusInProgress,
	}
	require.NoError(t, repo.CreateWorkflow(ctx, wf))

	task := &OnboardingTask{
		CompanyID:            1,
		OnboardingWorkflowID: wf.ID,
		TaskName:             "Create Email",
		SortOrder:            1,
	}
	require.NoError(t, tdb.DB.Create(task).Error)

	tests := []struct {
		name    string
		id      uint
		wantErr bool
	}{
		{name: "success", id: wf.ID, wantErr: false},
		{name: "not found", id: 999, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			found, err := repo.FindWorkflowByID(ctx, tt.id)
			if tt.wantErr {
				require.Error(t, err)
				assert.Nil(t, found)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.id, found.ID)
				assert.Len(t, found.Tasks, 1)
			}
		})
	}
}

func TestRepo_MarkWorkflowEmailSent(t *testing.T) {
	tdb := setupOnboardingTestDB(t)
	repo := NewRepository(tdb.DB)
	ctx := testutil.CtxWithTenant(1, 1, false)

	wf := &OnboardingWorkflow{
		CompanyID:    1,
		NewHireName:  "Jane Smith",
		NewHireEmail: "jane@example.com",
		Status:       WorkflowStatusInProgress,
	}
	require.NoError(t, repo.CreateWorkflow(ctx, wf))

	err := repo.MarkWorkflowEmailSent(ctx, wf.ID)
	require.NoError(t, err)

	found, err := repo.FindWorkflowByID(ctx, wf.ID)
	require.NoError(t, err)
	assert.True(t, found.WelcomeEmailSent)
}

func TestRepo_FindTaskByID(t *testing.T) {
	tdb := setupOnboardingTestDB(t)
	repo := NewRepository(tdb.DB)
	ctx := testutil.CtxWithTenant(1, 1, false)

	wf := &OnboardingWorkflow{
		CompanyID:    1,
		NewHireName:  "Jane Smith",
		NewHireEmail: "jane@example.com",
		Status:       WorkflowStatusInProgress,
	}
	require.NoError(t, repo.CreateWorkflow(ctx, wf))

	task := &OnboardingTask{
		CompanyID:            1,
		OnboardingWorkflowID: wf.ID,
		TaskName:             "Create Email",
		SortOrder:            1,
	}
	require.NoError(t, tdb.DB.Create(task).Error)

	tests := []struct {
		name    string
		id      uint
		wantErr bool
	}{
		{name: "success", id: task.ID, wantErr: false},
		{name: "not found", id: 999, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			found, err := repo.FindTaskByID(ctx, tt.id)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, "Create Email", found.TaskName)
			}
		})
	}
}

func TestRepo_CompleteTask(t *testing.T) {
	tdb := setupOnboardingTestDB(t)
	repo := NewRepository(tdb.DB)
	ctx := testutil.CtxWithTenant(1, 1, false)

	usr := &user.User{ID: 1, Username: "admin", PasswordHash: "hash", RoleID: 1, CompanyID: 1, IsActive: true}
	require.NoError(t, tdb.DB.Create(usr).Error)
	emp := &user.Employee{
		ID: 1, UserID: 1, CompanyID: 1, DepartmentID: 1, ShiftID: 1,
		NIK: "EMP001", FullName: "Admin User", Email: "admin@test.com", Position: "IT Admin",
	}
	require.NoError(t, tdb.DB.Create(emp).Error)

	wf := &OnboardingWorkflow{
		CompanyID:    1,
		NewHireName:  "Jane Smith",
		NewHireEmail: "jane@example.com",
		Status:       WorkflowStatusInProgress,
	}
	require.NoError(t, repo.CreateWorkflow(ctx, wf))

	task := &OnboardingTask{
		CompanyID:            1,
		OnboardingWorkflowID: wf.ID,
		TaskName:             "Create Email",
		SortOrder:            1,
	}
	require.NoError(t, tdb.DB.Create(task).Error)

	err := repo.CompleteTask(ctx, task.ID, 1, "Done")
	require.NoError(t, err)

	found, err := repo.FindTaskByID(ctx, task.ID)
	require.NoError(t, err)
	assert.True(t, found.IsCompleted)
	assert.Equal(t, uint(1), *found.CompletedBy)
}

func TestRepo_CountPendingTasks(t *testing.T) {
	tdb := setupOnboardingTestDB(t)
	repo := NewRepository(tdb.DB)
	ctx := testutil.CtxWithTenant(1, 1, false)

	wf := &OnboardingWorkflow{
		CompanyID:    1,
		NewHireName:  "Jane Smith",
		NewHireEmail: "jane@example.com",
		Status:       WorkflowStatusInProgress,
	}
	require.NoError(t, repo.CreateWorkflow(ctx, wf))

	tasks := []OnboardingTask{
		{CompanyID: 1, OnboardingWorkflowID: wf.ID, TaskName: "Task 1", IsCompleted: false},
		{CompanyID: 1, OnboardingWorkflowID: wf.ID, TaskName: "Task 2", IsCompleted: true},
	}
	require.NoError(t, tdb.DB.Create(&tasks).Error)

	count, err := repo.CountPendingTasks(ctx, wf.ID)
	require.NoError(t, err)
	assert.Equal(t, int64(1), count)
}

func TestRepo_CountTotalTasks(t *testing.T) {
	tdb := setupOnboardingTestDB(t)
	repo := NewRepository(tdb.DB)
	ctx := testutil.CtxWithTenant(1, 1, false)

	wf := &OnboardingWorkflow{
		CompanyID:    1,
		NewHireName:  "Jane Smith",
		NewHireEmail: "jane@example.com",
		Status:       WorkflowStatusInProgress,
	}
	require.NoError(t, repo.CreateWorkflow(ctx, wf))

	tasks := []OnboardingTask{
		{CompanyID: 1, OnboardingWorkflowID: wf.ID, TaskName: "Task 1"},
		{CompanyID: 1, OnboardingWorkflowID: wf.ID, TaskName: "Task 2"},
	}
	require.NoError(t, tdb.DB.Create(&tasks).Error)

	count, err := repo.CountTotalTasks(ctx, wf.ID)
	require.NoError(t, err)
	assert.Equal(t, int64(2), count)
}

func TestRepo_MarkWorkflowCompleted(t *testing.T) {
	tdb := setupOnboardingTestDB(t)
	repo := NewRepository(tdb.DB)
	ctx := testutil.CtxWithTenant(1, 1, false)

	wf := &OnboardingWorkflow{
		CompanyID:    1,
		NewHireName:  "Jane Smith",
		NewHireEmail: "jane@example.com",
		Status:       WorkflowStatusInProgress,
	}
	require.NoError(t, repo.CreateWorkflow(ctx, wf))

	err := repo.MarkWorkflowCompleted(ctx, wf.ID)
	require.NoError(t, err)

	found, err := repo.FindWorkflowByID(ctx, wf.ID)
	require.NoError(t, err)
	assert.Equal(t, WorkflowStatusCompleted, found.Status)
}
