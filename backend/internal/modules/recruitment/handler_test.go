package recruitment

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

func TestHandler_CreateRequisition(t *testing.T) {
	employeeID := uint(1)

	tests := []struct {
		name         string
		body         interface{}
		setupMocks   func(*mockService)
		setupContext func(*testutil.APITest)
		wantStatus   int
	}{
		{
			name: "success",
			body: CreateRequisitionRequest{
				DepartmentID:   1,
				Title:          "Senior Go Developer",
				Quantity:       2,
				EmploymentType: "PKWTT",
				Priority:       "HIGH",
			},
			setupMocks: func(svc *mockService) {
				svc.On("CreateRequisition", mock.Anything, uint(1), mock.AnythingOfType("*recruitment.CreateRequisitionRequest")).Return(nil)
			},
			setupContext: func(at *testutil.APITest) {
				at.WithAuthContext(&infrastructure.MyClaims{
					UserID: 1, CompanyID: 1, EmployeeID: &employeeID,
				})
			},
			wantStatus: http.StatusCreated,
		},
		{
			name: "service error",
			body: CreateRequisitionRequest{
				DepartmentID:   1,
				Title:          "Dev",
				Quantity:       1,
				EmploymentType: "PKWTT",
				Priority:       "MEDIUM",
			},
			setupMocks: func(svc *mockService) {
				svc.On("CreateRequisition", mock.Anything, uint(1), mock.AnythingOfType("*recruitment.CreateRequisitionRequest")).Return(errors.New("db error"))
			},
			setupContext: func(at *testutil.APITest) {
				at.WithAuthContext(&infrastructure.MyClaims{
					UserID: 1, CompanyID: 1, EmployeeID: &employeeID,
				})
			},
			wantStatus: http.StatusInternalServerError,
		},
		{
			name: "validation error missing required fields",
			body: CreateRequisitionRequest{},
			setupMocks: func(svc *mockService) {},
			setupContext: func(at *testutil.APITest) {
				at.WithAuthContext(&infrastructure.MyClaims{
					UserID: 1, CompanyID: 1, EmployeeID: &employeeID,
				})
			},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := new(mockService)
			tt.setupMocks(svc)
			handler := NewHandler(svc)

			at := testutil.NewAPITest(t, http.MethodPost, "/api/recruitment/requisitions", tt.body)
			tt.setupContext(at)

			rec, err := at.Execute(handler.CreateRequisition)
			require.NoError(t, err)
			assert.Equal(t, tt.wantStatus, rec.Code)
		})
	}
}

func TestHandler_SubmitRequisition(t *testing.T) {
	employeeID := uint(1)

	tests := []struct {
		name         string
		pathParams   map[string]string
		setupMocks   func(*mockService)
		setupContext func(*testutil.APITest)
		wantStatus   int
	}{
		{
			name:       "success",
			pathParams: map[string]string{"id": "1"},
			setupMocks: func(svc *mockService) {
				svc.On("SubmitRequisition", mock.Anything, uint(1), uint(1)).Return(nil)
			},
			setupContext: func(at *testutil.APITest) {
				at.WithAuthContext(&infrastructure.MyClaims{
					UserID: 1, CompanyID: 1, EmployeeID: &employeeID,
				})
			},
			wantStatus: http.StatusOK,
		},
		{
			name:       "service error",
			pathParams: map[string]string{"id": "1"},
			setupMocks: func(svc *mockService) {
				svc.On("SubmitRequisition", mock.Anything, uint(1), uint(1)).Return(errors.New("not requester"))
			},
			setupContext: func(at *testutil.APITest) {
				at.WithAuthContext(&infrastructure.MyClaims{
					UserID: 1, CompanyID: 1, EmployeeID: &employeeID,
				})
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "invalid id",
			pathParams: map[string]string{"id": "abc"},
			setupMocks: func(svc *mockService) {},
			setupContext: func(at *testutil.APITest) {
				at.WithAuthContext(&infrastructure.MyClaims{
					UserID: 1, CompanyID: 1, EmployeeID: &employeeID,
				})
			},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := new(mockService)
			tt.setupMocks(svc)
			handler := NewHandler(svc)

			at := testutil.NewAPITest(t, http.MethodPost, "/api/recruitment/requisitions/:id/submit", nil)
			at.WithPathParams(tt.pathParams)
			tt.setupContext(at)

			rec, err := at.Execute(handler.SubmitRequisition)
			require.NoError(t, err)
			assert.Equal(t, tt.wantStatus, rec.Code)
		})
	}
}

func TestHandler_RequisitionAction(t *testing.T) {
	employeeID := uint(1)

	tests := []struct {
		name         string
		pathParams   map[string]string
		body         interface{}
		setupMocks   func(*mockService)
		setupContext func(*testutil.APITest)
		wantStatus   int
	}{
		{
			name:       "approve success",
			pathParams: map[string]string{"id": "1"},
			body: RequisitionActionRequest{
				Action: "APPROVE",
			},
			setupMocks: func(svc *mockService) {
				svc.On("RequisitionAction", mock.Anything, uint(1), uint(1), mock.AnythingOfType("*recruitment.RequisitionActionRequest")).Return(nil)
			},
			setupContext: func(at *testutil.APITest) {
				at.WithAuthContext(&infrastructure.MyClaims{
					UserID: 1, CompanyID: 1, EmployeeID: &employeeID,
				})
			},
			wantStatus: http.StatusOK,
		},
		{
			name:       "validation error missing action",
			pathParams: map[string]string{"id": "1"},
			body:       map[string]interface{}{},
			setupMocks: func(svc *mockService) {},
			setupContext: func(at *testutil.APITest) {
				at.WithAuthContext(&infrastructure.MyClaims{
					UserID: 1, CompanyID: 1, EmployeeID: &employeeID,
				})
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "invalid id",
			pathParams: map[string]string{"id": "abc"},
			body:       RequisitionActionRequest{Action: "APPROVE"},
			setupMocks: func(svc *mockService) {},
			setupContext: func(at *testutil.APITest) {
				at.WithAuthContext(&infrastructure.MyClaims{
					UserID: 1, CompanyID: 1, EmployeeID: &employeeID,
				})
			},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := new(mockService)
			tt.setupMocks(svc)
			handler := NewHandler(svc)

			at := testutil.NewAPITest(t, http.MethodPost, "/api/recruitment/requisitions/:id/action", tt.body)
			at.WithPathParams(tt.pathParams)
			tt.setupContext(at)

			rec, err := at.Execute(handler.RequisitionAction)
			require.NoError(t, err)
			assert.Equal(t, tt.wantStatus, rec.Code)
		})
	}
}

func TestHandler_GetRequisitions(t *testing.T) {
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
				svc.On("GetRequisitions", mock.Anything, mock.AnythingOfType("*recruitment.RequisitionFilter")).
					Return([]RequisitionListResponse{}, (*response.Meta)(nil), nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name:        "service error",
			queryParams: "?page=1&limit=10",
			setupMocks: func(svc *mockService) {
				svc.On("GetRequisitions", mock.Anything, mock.AnythingOfType("*recruitment.RequisitionFilter")).
					Return([]RequisitionListResponse(nil), (*response.Meta)(nil), errors.New("db error"))
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := new(mockService)
			tt.setupMocks(svc)
			handler := NewHandler(svc)

			at := testutil.NewAPITest(t, http.MethodGet, "/api/recruitment/requisitions"+tt.queryParams, nil)

			rec, err := at.Execute(handler.GetRequisitions)
			require.NoError(t, err)
			assert.Equal(t, tt.wantStatus, rec.Code)
		})
	}
}

func TestHandler_GetRequisitionDetail(t *testing.T) {
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
				svc.On("GetRequisitionDetail", mock.Anything, uint(1)).Return(&RequisitionDetailResponse{ID: 1}, nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name:       "not found",
			pathParams: map[string]string{"id": "99"},
			setupMocks: func(svc *mockService) {
				svc.On("GetRequisitionDetail", mock.Anything, uint(99)).Return(nil, errors.New("requisition not found"))
			},
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "invalid id",
			pathParams: map[string]string{"id": "abc"},
			setupMocks: func(svc *mockService) {},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := new(mockService)
			tt.setupMocks(svc)
			handler := NewHandler(svc)

			at := testutil.NewAPITest(t, http.MethodGet, "/api/recruitment/requisitions/:id", nil)
			at.WithPathParams(tt.pathParams)

			rec, err := at.Execute(handler.GetRequisitionDetail)
			require.NoError(t, err)
			assert.Equal(t, tt.wantStatus, rec.Code)
		})
	}
}

func TestHandler_CloseRequisition(t *testing.T) {
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
				svc.On("CloseRequisition", mock.Anything, uint(1)).Return(nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name:       "service error",
			pathParams: map[string]string{"id": "1"},
			setupMocks: func(svc *mockService) {
				svc.On("CloseRequisition", mock.Anything, uint(1)).Return(errors.New("already closed"))
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "invalid id",
			pathParams: map[string]string{"id": "abc"},
			setupMocks: func(svc *mockService) {},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := new(mockService)
			tt.setupMocks(svc)
			handler := NewHandler(svc)

			at := testutil.NewAPITest(t, http.MethodPost, "/api/recruitment/requisitions/:id/close", nil)
			at.WithPathParams(tt.pathParams)

			rec, err := at.Execute(handler.CloseRequisition)
			require.NoError(t, err)
			assert.Equal(t, tt.wantStatus, rec.Code)
		})
	}
}

func TestHandler_DeleteRequisition(t *testing.T) {
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
				svc.On("DeleteRequisition", mock.Anything, uint(1)).Return(nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name:       "service error",
			pathParams: map[string]string{"id": "1"},
			setupMocks: func(svc *mockService) {
				svc.On("DeleteRequisition", mock.Anything, uint(1)).Return(errors.New("not found"))
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "invalid id",
			pathParams: map[string]string{"id": "abc"},
			setupMocks: func(svc *mockService) {},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := new(mockService)
			tt.setupMocks(svc)
			handler := NewHandler(svc)

			at := testutil.NewAPITest(t, http.MethodDelete, "/api/recruitment/requisitions/:id", nil)
			at.WithPathParams(tt.pathParams)

			rec, err := at.Execute(handler.DeleteRequisition)
			require.NoError(t, err)
			assert.Equal(t, tt.wantStatus, rec.Code)
		})
	}
}

func TestHandler_AddApplicant(t *testing.T) {
	tests := []struct {
		name         string
		pathParams   map[string]string
		body         interface{}
		setupMocks   func(*mockService)
		wantStatus   int
	}{
		{
			name:       "success",
			pathParams: map[string]string{"id": "1"},
			body: CreateApplicantRequest{
				FullName: "Jane Smith", Email: "jane@example.com", PhoneNumber: "0812345678",
			},
			setupMocks: func(svc *mockService) {
				svc.On("AddApplicant", mock.Anything, uint(1), mock.AnythingOfType("*recruitment.CreateApplicantRequest")).Return(nil)
			},
			wantStatus: http.StatusCreated,
		},
		{
			name:       "service error",
			pathParams: map[string]string{"id": "1"},
			body: CreateApplicantRequest{
				FullName: "Jane Smith", Email: "jane@example.com",
			},
			setupMocks: func(svc *mockService) {
				svc.On("AddApplicant", mock.Anything, uint(1), mock.AnythingOfType("*recruitment.CreateApplicantRequest")).Return(errors.New("requisition not found"))
			},
			wantStatus: http.StatusInternalServerError,
		},
		{
			name:       "validation error missing required fields",
			pathParams: map[string]string{"id": "1"},
			body:       CreateApplicantRequest{},
			setupMocks: func(svc *mockService) {},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "invalid requisition id",
			pathParams: map[string]string{"id": "abc"},
			body: CreateApplicantRequest{
				FullName: "Jane", Email: "jane@example.com",
			},
			setupMocks: func(svc *mockService) {},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := new(mockService)
			tt.setupMocks(svc)
			handler := NewHandler(svc)

			at := testutil.NewAPITest(t, http.MethodPost, "/api/recruitment/requisitions/:id/applicants", tt.body)
			at.WithPathParams(tt.pathParams)

			rec, err := at.Execute(handler.AddApplicant)
			require.NoError(t, err)
			assert.Equal(t, tt.wantStatus, rec.Code)
		})
	}
}

func TestHandler_UpdateApplicantStage(t *testing.T) {
	employeeID := uint(1)

	tests := []struct {
		name         string
		pathParams   map[string]string
		body         interface{}
		setupMocks   func(*mockService)
		setupContext func(*testutil.APITest)
		wantStatus   int
	}{
		{
			name:       "success",
			pathParams: map[string]string{"id": "1"},
			body: UpdateApplicantStageRequest{
				Stage: "INTERVIEW",
			},
			setupMocks: func(svc *mockService) {
				svc.On("UpdateStage", mock.Anything, uint(1), uint(1), mock.AnythingOfType("*recruitment.UpdateApplicantStageRequest")).Return(nil)
			},
			setupContext: func(at *testutil.APITest) {
				at.WithAuthContext(&infrastructure.MyClaims{
					UserID: 1, CompanyID: 1, EmployeeID: &employeeID,
				})
			},
			wantStatus: http.StatusOK,
		},
		{
			name:       "validation error missing stage",
			pathParams: map[string]string{"id": "1"},
			body:       map[string]interface{}{},
			setupMocks: func(svc *mockService) {},
			setupContext: func(at *testutil.APITest) {
				at.WithAuthContext(&infrastructure.MyClaims{
					UserID: 1, CompanyID: 1, EmployeeID: &employeeID,
				})
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "invalid id",
			pathParams: map[string]string{"id": "abc"},
			body:       UpdateApplicantStageRequest{Stage: "INTERVIEW"},
			setupMocks: func(svc *mockService) {},
			setupContext: func(at *testutil.APITest) {
				at.WithAuthContext(&infrastructure.MyClaims{
					UserID: 1, CompanyID: 1, EmployeeID: &employeeID,
				})
			},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := new(mockService)
			tt.setupMocks(svc)
			handler := NewHandler(svc)

			at := testutil.NewAPITest(t, http.MethodPut, "/api/recruitment/applicants/:id/stage", tt.body)
			at.WithPathParams(tt.pathParams)
			tt.setupContext(at)

			rec, err := at.Execute(handler.UpdateApplicantStage)
			require.NoError(t, err)
			assert.Equal(t, tt.wantStatus, rec.Code)
		})
	}
}

func TestHandler_GetApplicants(t *testing.T) {
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
				svc.On("GetApplicantsByRequisition", mock.Anything, uint(1)).Return(&KanbanBoardResponse{}, nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name:       "service error",
			pathParams: map[string]string{"id": "1"},
			setupMocks: func(svc *mockService) {
				svc.On("GetApplicantsByRequisition", mock.Anything, uint(1)).Return(nil, errors.New("db error"))
			},
			wantStatus: http.StatusInternalServerError,
		},
		{
			name:       "invalid id",
			pathParams: map[string]string{"id": "abc"},
			setupMocks: func(svc *mockService) {},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := new(mockService)
			tt.setupMocks(svc)
			handler := NewHandler(svc)

			at := testutil.NewAPITest(t, http.MethodGet, "/api/recruitment/requisitions/:id/applicants", nil)
			at.WithPathParams(tt.pathParams)

			rec, err := at.Execute(handler.GetApplicants)
			require.NoError(t, err)
			assert.Equal(t, tt.wantStatus, rec.Code)

			if tt.wantStatus == http.StatusOK {
				var resp map[string]interface{}
				require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
			}
		})
	}
}

func TestHandler_GetApplicantDetail(t *testing.T) {
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
				svc.On("GetApplicantDetail", mock.Anything, uint(1)).Return(&ApplicantDetailResponse{ID: 1}, nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name:       "not found",
			pathParams: map[string]string{"id": "99"},
			setupMocks: func(svc *mockService) {
				svc.On("GetApplicantDetail", mock.Anything, uint(99)).Return(nil, errors.New("applicant not found"))
			},
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "invalid id",
			pathParams: map[string]string{"id": "abc"},
			setupMocks: func(svc *mockService) {},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := new(mockService)
			tt.setupMocks(svc)
			handler := NewHandler(svc)

			at := testutil.NewAPITest(t, http.MethodGet, "/api/recruitment/applicants/:id", nil)
			at.WithPathParams(tt.pathParams)

			rec, err := at.Execute(handler.GetApplicantDetail)
			require.NoError(t, err)
			assert.Equal(t, tt.wantStatus, rec.Code)
		})
	}
}
