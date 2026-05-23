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
		&OnboardingTemplate{},
		&OnboardingTemplateItem{},
		&OnboardingWorkflow{},
		&OnboardingTask{},
		&user.User{},
		&user.Employee{},
	)
	t.Cleanup(tdb.Close)
	return tdb
}

func seedOnboardingData(t *testing.T, db *testutil.TestDB) {
	t.Helper()

	tmpl := &OnboardingTemplate{
		CompanyID:  1,
		Name:       "IT Setup",
		Department: "IT",
		Items: []OnboardingTemplateItem{
			{CompanyID: 1, TaskName: "Create Email", Description: "Setup email", SortOrder: 1},
			{CompanyID: 1, TaskName: "Setup Laptop", Description: "Provision laptop", SortOrder: 2},
		},
	}
	require.NoError(t, db.DB.Create(tmpl).Error)

	usr := &user.User{ID: 1, Username: "admin", PasswordHash: "hash", RoleID: 1, CompanyID: 1, IsActive: true}
	require.NoError(t, db.DB.Create(usr).Error)

	emp := &user.Employee{
		ID: 1, UserID: 1, CompanyID: 1, DepartmentID: 1, ShiftID: 1,
		NIK: "EMP001", FullName: "Admin User", Email: "admin@test.com", Position: "IT Admin",
	}
	require.NoError(t, db.DB.Create(emp).Error)
}

func TestRepo_CreateTemplate(t *testing.T) {
	tdb := setupOnboardingTestDB(t)
	repo := NewRepository(tdb.DB)
	ctx := testutil.CtxWithTenant(1, 1, false)

	tests := []struct {
		name    string
		tmpl    *OnboardingTemplate
		wantErr bool
	}{
		{
			name: "success",
			tmpl: &OnboardingTemplate{
				CompanyID:  1,
				Name:       "HR Onboarding",
				Department: "HR",
				Items: []OnboardingTemplateItem{
					{CompanyID: 1, TaskName: "Sign NDA", SortOrder: 1},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.CreateTemplate(ctx, tt.tmpl)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.NotZero(t, tt.tmpl.ID)
			}
		})
	}
}

func TestRepo_FindAllTemplates(t *testing.T) {
	tdb := setupOnboardingTestDB(t)
	repo := NewRepository(tdb.DB)
	ctx := testutil.CtxWithTenant(1, 1, false)

	seedOnboardingData(t, tdb)

	templates, err := repo.FindAllTemplates(ctx)
	require.NoError(t, err)
	assert.Len(t, templates, 1)
	assert.Equal(t, "IT Setup", templates[0].Name)
	assert.Len(t, templates[0].Items, 2)
}

func TestRepo_FindTemplateByID(t *testing.T) {
	tdb := setupOnboardingTestDB(t)
	repo := NewRepository(tdb.DB)
	ctx := testutil.CtxWithTenant(1, 1, false)

	seedOnboardingData(t, tdb)

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
			tmpl, err := repo.FindTemplateByID(ctx, tt.id)
			if tt.wantErr {
				require.Error(t, err)
				assert.Nil(t, tmpl)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.id, tmpl.ID)
				assert.Len(t, tmpl.Items, 2)
			}
		})
	}
}

func TestRepo_UpdateTemplate(t *testing.T) {
	tdb := setupOnboardingTestDB(t)
	repo := NewRepository(tdb.DB)
	ctx := testutil.CtxWithTenant(1, 1, true)

	seedOnboardingData(t, tdb)

	tmpl, err := repo.FindTemplateByID(ctx, 1)
	require.NoError(t, err)

	tmpl.Name = "IT Setup Updated"
	tmpl.Items = []OnboardingTemplateItem{
		{CompanyID: 1, TemplateID: 1, TaskName: "Create Email Account", SortOrder: 1},
	}

	err = repo.UpdateTemplate(ctx, tmpl)
	require.NoError(t, err)

	updated, err := repo.FindTemplateByID(ctx, 1)
	require.NoError(t, err)
	assert.Equal(t, "IT Setup Updated", updated.Name)
	assert.Len(t, updated.Items, 1)
	assert.Equal(t, "Create Email Account", updated.Items[0].TaskName)
}

func TestRepo_DeleteTemplate(t *testing.T) {
	tdb := setupOnboardingTestDB(t)
	repo := NewRepository(tdb.DB)
	ctx := testutil.CtxWithTenant(1, 1, false)

	seedOnboardingData(t, tdb)

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
			err := repo.DeleteTemplate(ctx, tt.id)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
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
			Department:           "IT",
			SortOrder:            1,
		},
		{
			CompanyID:            1,
			OnboardingWorkflowID: wf.ID,
			TaskName:             "Sign NDA",
			Department:           "HR",
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
		Department:           "IT",
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

	seedOnboardingData(t, tdb)

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
		Department:           "IT",
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
		Department:           "IT",
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

	seedOnboardingData(t, tdb)

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
		Department:           "IT",
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
		{CompanyID: 1, OnboardingWorkflowID: wf.ID, TaskName: "Task 1", Department: "IT", IsCompleted: false},
		{CompanyID: 1, OnboardingWorkflowID: wf.ID, TaskName: "Task 2", Department: "HR", IsCompleted: true},
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
		{CompanyID: 1, OnboardingWorkflowID: wf.ID, TaskName: "Task 1", Department: "IT"},
		{CompanyID: 1, OnboardingWorkflowID: wf.ID, TaskName: "Task 2", Department: "HR"},
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
