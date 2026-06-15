package onboarding

import (
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"basekarya-backend/internal/infrastructure"
	"basekarya-backend/internal/testutil"
	"basekarya-backend/pkg/response"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestHandler_CreateWorkflow(t *testing.T) {
	tests := []struct {
		name       string
		body       interface{}
		setupMocks func(*mockService)
		wantStatus int
	}{
		{
			name: "success",
			body: CreateWorkflowRequest{
				NewHireName:  "Jane Smith",
				NewHireEmail: "jane@example.com",
				Position:     "Developer",
				Department:   "Engineering",
			},
			setupMocks: func(svc *mockService) {
				svc.On("CreateWorkflow", mock.Anything, mock.AnythingOfType("*onboarding.CreateWorkflowRequest")).Return(nil)
			},
			wantStatus: http.StatusCreated,
		},
		{
			name: "invalid body - missing email",
			body: CreateWorkflowRequest{
				NewHireName: "Jane Smith",
			},
			setupMocks: func(svc *mockService) {},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "service error",
			body: CreateWorkflowRequest{
				NewHireName:  "Jane Smith",
				NewHireEmail: "jane@example.com",
			},
			setupMocks: func(svc *mockService) {
				svc.On("CreateWorkflow", mock.Anything, mock.AnythingOfType("*onboarding.CreateWorkflowRequest")).Return(errors.New("db error"))
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := new(mockService)
			tt.setupMocks(svc)
			handler := NewHandler(svc)

			at := testutil.NewAPITest(t, http.MethodPost, "/api/onboarding/workflows", tt.body)
			at.WithAuthContext(&infrastructure.MyClaims{UserID: 1, CompanyID: 1})

			rec, err := at.Execute(handler.CreateWorkflow)
			require.NoError(t, err)
			assert.Equal(t, tt.wantStatus, rec.Code)

			var resp map[string]interface{}
			require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
			if tt.wantStatus < 400 {
				assert.Nil(t, resp["error"])
			}
		})
	}
}

func TestHandler_GetWorkflows(t *testing.T) {
	tests := []struct {
		name        string
		queryParams string
		setupMocks  func(*mockService)
		wantStatus  int
	}{
		{
			name:        "success",
			queryParams: "?page=1&limit=10",
			setupMocks: func(svc *mockService) {
				svc.On("GetWorkflows", mock.Anything, mock.AnythingOfType("*onboarding.WorkflowFilter")).
					Return([]WorkflowListResponse{}, (*response.Meta)(nil), nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name:        "service error",
			queryParams: "?page=1&limit=10",
			setupMocks: func(svc *mockService) {
				svc.On("GetWorkflows", mock.Anything, mock.AnythingOfType("*onboarding.WorkflowFilter")).
					Return([]WorkflowListResponse(nil), (*response.Meta)(nil), errors.New("db error"))
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := new(mockService)
			tt.setupMocks(svc)
			handler := NewHandler(svc)

			at := testutil.NewAPITest(t, http.MethodGet, "/api/onboarding/workflows"+tt.queryParams, nil)
			at.WithAuthContext(&infrastructure.MyClaims{UserID: 1, CompanyID: 1})

			rec, err := at.Execute(handler.GetWorkflows)
			require.NoError(t, err)
			assert.Equal(t, tt.wantStatus, rec.Code)

			var resp map[string]interface{}
			require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
			if tt.wantStatus < 400 {
				assert.Nil(t, resp["error"])
			}
		})
	}
}

func TestHandler_GetWorkflowDetail(t *testing.T) {
	tests := []struct {
		name       string
		pathParams map[string]string
		setupMocks func(*mockService)
		wantStatus int
	}{
		{
			name:       "success",
			pathParams: map[string]string{"id": "1"},
			setupMocks: func(svc *mockService) {
				svc.On("GetWorkflowDetail", mock.Anything, uint(1)).Return(&WorkflowDetailResponse{
					ID: 1, NewHireName: "Jane", Status: WorkflowStatusInProgress,
				}, nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name:       "invalid id",
			pathParams: map[string]string{"id": "abc"},
			setupMocks: func(svc *mockService) {},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "not found",
			pathParams: map[string]string{"id": "999"},
			setupMocks: func(svc *mockService) {
				svc.On("GetWorkflowDetail", mock.Anything, uint(999)).Return(nil, errors.New("not found"))
			},
			wantStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := new(mockService)
			tt.setupMocks(svc)
			handler := NewHandler(svc)

			at := testutil.NewAPITest(t, http.MethodGet, "/api/onboarding/workflows/:id", nil)
			at.WithPathParams(tt.pathParams)
			at.WithAuthContext(&infrastructure.MyClaims{UserID: 1, CompanyID: 1})

			rec, err := at.Execute(handler.GetWorkflowDetail)
			require.NoError(t, err)
			assert.Equal(t, tt.wantStatus, rec.Code)
		})
	}
}

func TestHandler_CompleteTask(t *testing.T) {
	tests := []struct {
		name       string
		pathParams map[string]string
		body       interface{}
		setupMocks func(*mockService)
		wantStatus int
	}{
		{
			name:       "success",
			pathParams: map[string]string{"id": "1"},
			body:       CompleteTaskRequest{Notes: "Done"},
			setupMocks: func(svc *mockService) {
				svc.On("CompleteTask", mock.Anything, uint(1), uint(1), mock.AnythingOfType("*onboarding.CompleteTaskRequest")).Return(nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name:       "invalid id",
			pathParams: map[string]string{"id": "abc"},
			body:       CompleteTaskRequest{Notes: "Done"},
			setupMocks: func(svc *mockService) {},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "service error",
			pathParams: map[string]string{"id": "1"},
			body:       CompleteTaskRequest{Notes: "Done"},
			setupMocks: func(svc *mockService) {
				svc.On("CompleteTask", mock.Anything, uint(1), uint(1), mock.AnythingOfType("*onboarding.CompleteTaskRequest")).Return(errors.New("not found"))
			},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := new(mockService)
			tt.setupMocks(svc)
			handler := NewHandler(svc)

			at := testutil.NewAPITest(t, http.MethodPost, "/api/onboarding/tasks/:id/complete", tt.body)
			at.WithPathParams(tt.pathParams)
			at.WithAuthContext(&infrastructure.MyClaims{UserID: 1, CompanyID: 1})

			rec, err := at.Execute(handler.CompleteTask)
			require.NoError(t, err)
			assert.Equal(t, tt.wantStatus, rec.Code)
		})
	}
}

func TestHandler_UpdateWorkflowTasks(t *testing.T) {
	tests := []struct {
		name       string
		pathParams map[string]string
		body       interface{}
		setupMocks func(*mockService)
		wantStatus int
	}{
		{
			name:       "success",
			pathParams: map[string]string{"id": "1"},
			body: UpdateWorkflowTasksRequest{
				Tasks: []WorkflowTaskRequest{
					{TaskName: "Setup Email", SortOrder: 1},
				},
			},
			setupMocks: func(svc *mockService) {
				svc.On("UpdateWorkflowTasks", mock.Anything, uint(1), mock.AnythingOfType("*onboarding.UpdateWorkflowTasksRequest")).Return(nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name:       "invalid id",
			pathParams: map[string]string{"id": "abc"},
			body: UpdateWorkflowTasksRequest{
				Tasks: []WorkflowTaskRequest{
					{TaskName: "Task", SortOrder: 1},
				},
			},
			setupMocks: func(svc *mockService) {},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "empty tasks",
			pathParams: map[string]string{"id": "1"},
			body:       UpdateWorkflowTasksRequest{Tasks: []WorkflowTaskRequest{}},
			setupMocks: func(svc *mockService) {},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "service error",
			pathParams: map[string]string{"id": "1"},
			body: UpdateWorkflowTasksRequest{
				Tasks: []WorkflowTaskRequest{
					{TaskName: "Task", SortOrder: 1},
				},
			},
			setupMocks: func(svc *mockService) {
				svc.On("UpdateWorkflowTasks", mock.Anything, uint(1), mock.AnythingOfType("*onboarding.UpdateWorkflowTasksRequest")).Return(errors.New("not found"))
			},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := new(mockService)
			tt.setupMocks(svc)
			handler := NewHandler(svc)

			at := testutil.NewAPITest(t, http.MethodPut, "/api/onboarding/workflows/:id/tasks", tt.body)
			at.WithPathParams(tt.pathParams)
			at.WithAuthContext(&infrastructure.MyClaims{UserID: 1, CompanyID: 1})

			rec, err := at.Execute(handler.UpdateWorkflowTasks)
			require.NoError(t, err)
			assert.Equal(t, tt.wantStatus, rec.Code)

			var resp map[string]interface{}
			require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
			if tt.wantStatus < 400 {
				assert.Nil(t, resp["error"])
			}
		})
	}
}
