package onboarding

import (
	"errors"
	"testing"
	"time"

	"basekarya-backend/internal/modules/company"
	"basekarya-backend/internal/modules/department"
	"basekarya-backend/internal/modules/master"
	"basekarya-backend/internal/modules/rbac"
	"basekarya-backend/internal/modules/user"
	"basekarya-backend/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestService_CreateWorkflow(t *testing.T) {
	ctx := testutil.CtxWithTenant(1, 1, false)

	tests := []struct {
		name       string
		req        *CreateWorkflowRequest
		setupMocks func(*mockRepo, *mockUserProvider, *mockEmailProvider, *mockCompanyProvider)
		wantErr    bool
		errMsg     string
	}{
		{
			name: "success",
			req: &CreateWorkflowRequest{
				NewHireName:  "Jane Smith",
				NewHireEmail: "jane@example.com",
				Position:     "Developer",
				Department:   "Engineering",
				StartDate:    "2025-01-15",
				Tasks: []WorkflowTaskRequest{
					{TaskName: "Create Email", Description: "Setup email", SortOrder: 1},
				},
			},
			setupMocks: func(repo *mockRepo, userProv *mockUserProvider, emailProv *mockEmailProvider, companyProv *mockCompanyProvider) {
				repo.On("CreateWorkflow", mock.Anything, mock.AnythingOfType("*onboarding.OnboardingWorkflow")).Return(nil)
				repo.On("CreateTasks", mock.Anything, mock.Anything).Return(nil)
				repo.On("MarkWorkflowEmailSent", mock.Anything, mock.Anything).Return(nil)
				companyProv.On("FindByID", mock.Anything, uint(1)).Return(&company.Company{ID: 1, Name: "TestCo"}, nil)
				emailProv.On("Send", "jane@example.com", mock.Anything, mock.Anything).Return(nil)
				userProv.On("FindApprovalUsers", mock.Anything, mock.Anything).Return([]uint{}, nil)
			},
			wantErr: false,
		},
		{
			name: "create workflow error",
			req: &CreateWorkflowRequest{
				NewHireName:  "Jane Smith",
				NewHireEmail: "jane@example.com",
			},
			setupMocks: func(repo *mockRepo, userProv *mockUserProvider, emailProv *mockEmailProvider, companyProv *mockCompanyProvider) {
				repo.On("CreateWorkflow", mock.Anything, mock.AnythingOfType("*onboarding.OnboardingWorkflow")).Return(errors.New("db error"))
			},
			wantErr: true,
			errMsg:  "db error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, repo, _, userProv, emailProv, companyProv, _, _, _, _ := newTestOnboardingService()
			tt.setupMocks(repo, userProv, emailProv, companyProv)

			err := svc.CreateWorkflow(ctx, tt.req)

			if tt.wantErr {
				require.Error(t, err)
				assert.Equal(t, tt.errMsg, err.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestService_GetWorkflows(t *testing.T) {
	tests := []struct {
		name       string
		filter     *WorkflowFilter
		setupMocks func(*mockRepo)
		wantLen    int
		wantErr    bool
	}{
		{
			name:   "success with data",
			filter: &WorkflowFilter{Page: 1, Limit: 10},
			setupMocks: func(repo *mockRepo) {
				repo.On("FindAllWorkflows", mock.Anything, mock.AnythingOfType("*onboarding.WorkflowFilter")).Return(
					[]OnboardingWorkflow{
						{
							ID: 1, NewHireName: "Jane", NewHireEmail: "jane@example.com",
							Status: WorkflowStatusInProgress,
							Tasks:  []OnboardingTask{{IsCompleted: true}, {IsCompleted: false}},
						},
					}, int64(1), nil)
			},
			wantLen: 1,
			wantErr: false,
		},
		{
			name:   "default pagination",
			filter: &WorkflowFilter{Page: 0, Limit: 0},
			setupMocks: func(repo *mockRepo) {
				repo.On("FindAllWorkflows", mock.Anything, mock.AnythingOfType("*onboarding.WorkflowFilter")).Return(
					[]OnboardingWorkflow{}, int64(0), nil)
			},
			wantLen: 0,
			wantErr: false,
		},
		{
			name:   "repo error",
			filter: &WorkflowFilter{Page: 1, Limit: 10},
			setupMocks: func(repo *mockRepo) {
				repo.On("FindAllWorkflows", mock.Anything, mock.AnythingOfType("*onboarding.WorkflowFilter")).Return(
					[]OnboardingWorkflow(nil), int64(0), errors.New("db error"))
			},
			wantLen: 0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, repo, _, _, _, _, _, _, _, _ := newTestOnboardingService()
			tt.setupMocks(repo)

			result, meta, err := svc.GetWorkflows(testutil.CtxWithTenant(1, 1, false), tt.filter)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Len(t, result, tt.wantLen)
				if tt.wantLen > 0 {
					assert.NotNil(t, meta)
					assert.Equal(t, 50, result[0].Progress)
				}
			}
		})
	}
}

func TestService_GetWorkflowDetail(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name       string
		id         uint
		setupMocks func(*mockRepo)
		wantErr    bool
		errMsg     string
	}{
		{
			name: "success with tasks",
			id:   1,
			setupMocks: func(repo *mockRepo) {
				repo.On("FindWorkflowByID", mock.Anything, uint(1)).Return(&OnboardingWorkflow{
					ID: 1, NewHireName: "Jane", NewHireEmail: "jane@example.com",
					Status: WorkflowStatusInProgress, CreatedAt: now,
					Tasks: []OnboardingTask{
						{ID: 1, TaskName: "Create Email", IsCompleted: true, CompletedBy: uintPtr(1), CompletedAt: &now, CompletedByUser: &user.User{Employee: &user.Employee{FullName: "Admin"}}},
						{ID: 2, TaskName: "Sign NDA", IsCompleted: false},
						{ID: 3, TaskName: "Other Task", IsCompleted: false},
					},
				}, nil)
			},
			wantErr: false,
		},
		{
			name: "not found",
			id:   999,
			setupMocks: func(repo *mockRepo) {
				repo.On("FindWorkflowByID", mock.Anything, uint(999)).Return(nil, errors.New("not found"))
			},
			wantErr: true,
			errMsg:  "workflow not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, repo, _, _, _, _, _, _, _, _ := newTestOnboardingService()
			tt.setupMocks(repo)

			result, err := svc.GetWorkflowDetail(testutil.CtxWithTenant(1, 1, false), tt.id)

			if tt.wantErr {
				require.Error(t, err)
				assert.Equal(t, tt.errMsg, err.Error())
			} else {
				require.NoError(t, err)
				assert.NotNil(t, result)
				assert.Len(t, result.Tasks, 3)
				assert.Equal(t, 33, result.Progress)
			}
		})
	}
}

func TestService_CompleteTask(t *testing.T) {
	ctx := testutil.CtxWithTenant(1, 1, false)

	tests := []struct {
		name       string
		taskID     uint
		userID     uint
		req        *CompleteTaskRequest
		setupMocks func(*mockRepo, *mockRoleProvider, *mockMasterProvider, *mockDepartmentProvider, *mockUserProvider)
		wantErr    bool
		errMsg     string
	}{
		{
			name:   "success with remaining tasks",
			taskID: 1,
			userID: 1,
			req:    &CompleteTaskRequest{Notes: "Done"},
			setupMocks: func(repo *mockRepo, roleProv *mockRoleProvider, masterProv *mockMasterProvider, deptProv *mockDepartmentProvider, userProv *mockUserProvider) {
				repo.On("FindTaskByID", mock.Anything, uint(1)).Return(&OnboardingTask{
					ID:                   1,
					OnboardingWorkflowID: 1,
					IsCompleted:          false,
				}, nil)
				repo.On("CompleteTask", mock.Anything, uint(1), uint(1), "Done").Return(nil)
				repo.On("CountPendingTasks", mock.Anything, uint(1)).Return(int64(1), nil)
			},
			wantErr: false,
		},
		{
			name:   "success workflow completed",
			taskID: 1,
			userID: 1,
			req:    &CompleteTaskRequest{Notes: "Final task"},
			setupMocks: func(repo *mockRepo, roleProv *mockRoleProvider, masterProv *mockMasterProvider, deptProv *mockDepartmentProvider, userProv *mockUserProvider) {
				repo.On("FindTaskByID", mock.Anything, uint(1)).Return(&OnboardingTask{
					ID:                   1,
					OnboardingWorkflowID: 1,
					IsCompleted:          false,
				}, nil)
				repo.On("CompleteTask", mock.Anything, uint(1), uint(1), "Final task").Return(nil)
				repo.On("CountPendingTasks", mock.Anything, uint(1)).Return(int64(0), nil)
				repo.On("MarkWorkflowCompleted", mock.Anything, uint(1)).Return(nil)
				repo.On("FindWorkflowByID", mock.Anything, uint(1)).Return(&OnboardingWorkflow{
					ID: 1, NewHireName: "Jane Smith", NewHireEmail: "jane@example.com", Position: "Dev",
				}, nil)
				roleProv.On("FindRoleByName", mock.Anything, "EMPLOYEE").Return(&rbac.Role{ID: 2}, nil)
				deptProv.On("FindByName", mock.Anything, "Umum").Return(&department.Department{ID: 1}, nil)
				masterProv.On("FindShiftByName", mock.Anything, "Regular").Return(&master.Shift{ID: 1}, nil)
				userProv.On("CreateEmployee", mock.Anything, mock.AnythingOfType("*user.CreateEmployeeRequest")).Return(&user.CreateEmployeeResponse{Username: "janesmith"}, nil)
			},
			wantErr: false,
		},
		{
			name:   "task not found",
			taskID: 999,
			userID: 1,
			req:    &CompleteTaskRequest{Notes: ""},
			setupMocks: func(repo *mockRepo, roleProv *mockRoleProvider, masterProv *mockMasterProvider, deptProv *mockDepartmentProvider, userProv *mockUserProvider) {
				repo.On("FindTaskByID", mock.Anything, uint(999)).Return(nil, errors.New("not found"))
			},
			wantErr: true,
			errMsg:  "task not found",
		},
		{
			name:   "task already completed",
			taskID: 1,
			userID: 1,
			req:    &CompleteTaskRequest{Notes: ""},
			setupMocks: func(repo *mockRepo, roleProv *mockRoleProvider, masterProv *mockMasterProvider, deptProv *mockDepartmentProvider, userProv *mockUserProvider) {
				repo.On("FindTaskByID", mock.Anything, uint(1)).Return(&OnboardingTask{
					ID:          1,
					IsCompleted: true,
				}, nil)
			},
			wantErr: true,
			errMsg:  "task is already completed",
		},
		{
			name:   "complete task repo error",
			taskID: 1,
			userID: 1,
			req:    &CompleteTaskRequest{Notes: "Done"},
			setupMocks: func(repo *mockRepo, roleProv *mockRoleProvider, masterProv *mockMasterProvider, deptProv *mockDepartmentProvider, userProv *mockUserProvider) {
				repo.On("FindTaskByID", mock.Anything, uint(1)).Return(&OnboardingTask{
					ID:                   1,
					OnboardingWorkflowID: 1,
					IsCompleted:          false,
				}, nil)
				repo.On("CompleteTask", mock.Anything, uint(1), uint(1), "Done").Return(errors.New("db error"))
			},
			wantErr: true,
			errMsg:  "db error",
		},
		{
			name:   "workflow completed but mark fails",
			taskID: 1,
			userID: 1,
			req:    &CompleteTaskRequest{Notes: "Done"},
			setupMocks: func(repo *mockRepo, roleProv *mockRoleProvider, masterProv *mockMasterProvider, deptProv *mockDepartmentProvider, userProv *mockUserProvider) {
				repo.On("FindTaskByID", mock.Anything, uint(1)).Return(&OnboardingTask{
					ID:                   1,
					OnboardingWorkflowID: 1,
					IsCompleted:          false,
				}, nil)
				repo.On("CompleteTask", mock.Anything, uint(1), uint(1), "Done").Return(nil)
				repo.On("CountPendingTasks", mock.Anything, uint(1)).Return(int64(0), nil)
				repo.On("MarkWorkflowCompleted", mock.Anything, uint(1)).Return(errors.New("mark error"))
			},
			wantErr: true,
			errMsg:  "mark error",
		},
		{
			name:   "workflow completed but role not found",
			taskID: 1,
			userID: 1,
			req:    &CompleteTaskRequest{Notes: "Done"},
			setupMocks: func(repo *mockRepo, roleProv *mockRoleProvider, masterProv *mockMasterProvider, deptProv *mockDepartmentProvider, userProv *mockUserProvider) {
				repo.On("FindTaskByID", mock.Anything, uint(1)).Return(&OnboardingTask{
					ID:                   1,
					OnboardingWorkflowID: 1,
					IsCompleted:          false,
				}, nil)
				repo.On("CompleteTask", mock.Anything, uint(1), uint(1), "Done").Return(nil)
				repo.On("CountPendingTasks", mock.Anything, uint(1)).Return(int64(0), nil)
				repo.On("MarkWorkflowCompleted", mock.Anything, uint(1)).Return(nil)
				repo.On("FindWorkflowByID", mock.Anything, uint(1)).Return(&OnboardingWorkflow{
					ID: 1, NewHireName: "Jane Smith", NewHireEmail: "jane@example.com", Position: "Dev",
				}, nil)
				roleProv.On("FindRoleByName", mock.Anything, "EMPLOYEE").Return(nil, errors.New("role not found"))
			},
			wantErr: true,
			errMsg:  "role not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, repo, _, userProv, _, _, roleProv, deptProv, masterProv, _ := newTestOnboardingService()
			tt.setupMocks(repo, roleProv, masterProv, deptProv, userProv)

			err := svc.CompleteTask(ctx, tt.taskID, tt.userID, tt.req)

			if tt.wantErr {
				require.Error(t, err)
				assert.Equal(t, tt.errMsg, err.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestService_GetWorkflows_ProgressCalculation(t *testing.T) {
	svc, repo, _, _, _, _, _, _, _, _ := newTestOnboardingService()

	repo.On("FindAllWorkflows", mock.Anything, mock.AnythingOfType("*onboarding.WorkflowFilter")).Return(
		[]OnboardingWorkflow{
			{
				ID:     1,
				Status: WorkflowStatusInProgress,
				Tasks:  []OnboardingTask{{IsCompleted: true}, {IsCompleted: true}, {IsCompleted: true}, {IsCompleted: false}},
			},
		}, int64(1), nil)

	result, meta, err := svc.GetWorkflows(testutil.CtxWithTenant(1, 1, false), &WorkflowFilter{Page: 1, Limit: 10})

	require.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, 75, result[0].Progress)
	assert.NotNil(t, meta)
}

func TestService_GetWorkflowDetail_ProgressZeroTasks(t *testing.T) {
	svc, repo, _, _, _, _, _, _, _, _ := newTestOnboardingService()

	repo.On("FindWorkflowByID", mock.Anything, uint(1)).Return(&OnboardingWorkflow{
		ID: 1, NewHireName: "Jane", Status: WorkflowStatusInProgress,
		Tasks: []OnboardingTask{},
	}, nil)

	result, err := svc.GetWorkflowDetail(testutil.CtxWithTenant(1, 1, false), 1)

	require.NoError(t, err)
	assert.Equal(t, 0, result.Progress)
	assert.Empty(t, result.Tasks)
}

func uintPtr(v uint) *uint {
	return &v
}

func TestService_GetWorkflows_PaginationDefaults(t *testing.T) {
	svc, repo, _, _, _, _, _, _, _, _ := newTestOnboardingService()

	repo.On("FindAllWorkflows", mock.Anything, mock.AnythingOfType("*onboarding.WorkflowFilter")).Return(
		[]OnboardingWorkflow{}, int64(0), nil)

	result, meta, err := svc.GetWorkflows(testutil.CtxWithTenant(1, 1, false), &WorkflowFilter{Page: 0, Limit: 0})

	require.NoError(t, err)
	assert.Empty(t, result)
	assert.NotNil(t, meta)
	assert.Equal(t, 1, meta.Page)
	assert.Equal(t, 10, meta.Limit)
}

func TestService_GetWorkflows_MetaResponse(t *testing.T) {
	svc, repo, _, _, _, _, _, _, _, _ := newTestOnboardingService()

	repo.On("FindAllWorkflows", mock.Anything, mock.AnythingOfType("*onboarding.WorkflowFilter")).Return(
		[]OnboardingWorkflow{
			{ID: 1, Status: WorkflowStatusInProgress, Tasks: []OnboardingTask{}},
		}, int64(15), nil)

	result, meta, err := svc.GetWorkflows(testutil.CtxWithTenant(1, 1, false), &WorkflowFilter{Page: 1, Limit: 10})

	require.NoError(t, err)
	assert.Len(t, result, 1)
	require.NotNil(t, meta)
	assert.Equal(t, int64(15), meta.TotalData)
	assert.Equal(t, int64(2), meta.TotalPage)
}

func TestService_UpdateWorkflowTasks(t *testing.T) {
	ctx := testutil.CtxWithTenant(1, 1, false)

	tests := []struct {
		name       string
		workflowID uint
		req        *UpdateWorkflowTasksRequest
		setupMocks func(*mockRepo)
		wantErr    bool
		errMsg     string
	}{
		{
			name:       "success",
			workflowID: 1,
			req: &UpdateWorkflowTasksRequest{
				Tasks: []WorkflowTaskRequest{
					{TaskName: "Setup Email", Description: "Create email", SortOrder: 1},
				},
			},
			setupMocks: func(repo *mockRepo) {
				repo.On("FindWorkflowByID", mock.Anything, uint(1)).Return(&OnboardingWorkflow{ID: 1}, nil)
				repo.On("DeletePendingTasks", mock.Anything, uint(1)).Return(nil)
				repo.On("CreateTasks", mock.Anything, mock.Anything).Return(nil)
			},
			wantErr: false,
		},
		{
			name:       "workflow not found",
			workflowID: 999,
			req: &UpdateWorkflowTasksRequest{
				Tasks: []WorkflowTaskRequest{
					{TaskName: "Task", SortOrder: 1},
				},
			},
			setupMocks: func(repo *mockRepo) {
				repo.On("FindWorkflowByID", mock.Anything, uint(999)).Return(nil, errors.New("not found"))
			},
			wantErr: true,
			errMsg:  "workflow not found",
		},
		{
			name:       "delete pending tasks error",
			workflowID: 1,
			req: &UpdateWorkflowTasksRequest{
				Tasks: []WorkflowTaskRequest{
					{TaskName: "Task", SortOrder: 1},
				},
			},
			setupMocks: func(repo *mockRepo) {
				repo.On("FindWorkflowByID", mock.Anything, uint(1)).Return(&OnboardingWorkflow{ID: 1}, nil)
				repo.On("DeletePendingTasks", mock.Anything, uint(1)).Return(errors.New("db error"))
			},
			wantErr: true,
			errMsg:  "db error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, repo, _, _, _, _, _, _, _, _ := newTestOnboardingService()
			tt.setupMocks(repo)

			err := svc.UpdateWorkflowTasks(ctx, tt.workflowID, tt.req)

			if tt.wantErr {
				require.Error(t, err)
				assert.Equal(t, tt.errMsg, err.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}
