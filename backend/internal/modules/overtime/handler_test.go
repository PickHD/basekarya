package overtime

import (
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"basekarya-backend/internal/infrastructure"
	"basekarya-backend/internal/testutil"
	"basekarya-backend/pkg/constants"
	"basekarya-backend/pkg/response"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestHandler_Create(t *testing.T) {
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
			body: OvertimeRequest{
				Date:      "2026-06-01",
				StartTime: "18:00",
				EndTime:   "20:00",
				Reason:    "Project deadline",
			},
			setupMocks: func(svc *mockService) {
				svc.On("Create", mock.Anything, mock.AnythingOfType("*overtime.OvertimeRequest")).Return(nil)
			},
			setupContext: func(at *testutil.APITest) {
				at.WithAuthContext(&infrastructure.MyClaims{
					UserID:      1,
					CompanyID:   1,
					EmployeeID:  &employeeID,
					Permissions: []string{constants.VIEW_OVERTIME},
				})
			},
			wantStatus: http.StatusCreated,
		},
		{
			name: "service error",
			body: OvertimeRequest{
				Date:      "2026-06-01",
				StartTime: "18:00",
				EndTime:   "20:00",
				Reason:    "Project deadline",
			},
			setupMocks: func(svc *mockService) {
				svc.On("Create", mock.Anything, mock.AnythingOfType("*overtime.OvertimeRequest")).Return(errors.New("db error"))
			},
			setupContext: func(at *testutil.APITest) {
				at.WithAuthContext(&infrastructure.MyClaims{
					UserID:      1,
					CompanyID:   1,
					EmployeeID:  &employeeID,
					Permissions: []string{constants.VIEW_OVERTIME},
				})
			},
			wantStatus: http.StatusInternalServerError,
		},
		{
			name: "missing required fields",
			body: OvertimeRequest{
				Date: "2026-06-01",
			},
			setupMocks: func(svc *mockService) {},
			setupContext: func(at *testutil.APITest) {
				at.WithAuthContext(&infrastructure.MyClaims{
					UserID:      1,
					CompanyID:   1,
					EmployeeID:  &employeeID,
					Permissions: []string{constants.VIEW_OVERTIME},
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

			at := testutil.NewAPITest(t, http.MethodPost, "/api/overtime", tt.body)
			tt.setupContext(at)

			rec, err := at.Execute(handler.Create)
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

func TestHandler_GetAll(t *testing.T) {
	tests := []struct {
		name         string
		queryParams  string
		setupMocks   func(*mockService)
		setupContext func(*testutil.APITest)
		wantStatus   int
	}{
		{
			name:        "success",
			queryParams: "?page=1&limit=10",
			setupMocks: func(svc *mockService) {
				svc.On("GetList", mock.Anything, mock.AnythingOfType("overtime.OvertimeFilter")).
					Return([]OvertimeListResponse{}, (*response.Meta)(nil), nil)
			},
			setupContext: func(at *testutil.APITest) {
				at.WithAuthContext(&infrastructure.MyClaims{
					UserID:      1,
					CompanyID:   1,
					Permissions: []string{constants.VIEW_OVERTIME},
				})
			},
			wantStatus: http.StatusOK,
		},
		{
			name:        "self only permission",
			queryParams: "?page=1&limit=10",
			setupMocks: func(svc *mockService) {
				svc.On("GetList", mock.Anything, mock.MatchedBy(func(f OvertimeFilter) bool {
					return f.UserID == 1
				})).Return([]OvertimeListResponse{}, (*response.Meta)(nil), nil)
			},
			setupContext: func(at *testutil.APITest) {
				at.WithAuthContext(&infrastructure.MyClaims{
					UserID:      1,
					CompanyID:   1,
					Permissions: []string{constants.VIEW_SELF_OVERTIME},
				})
			},
			wantStatus: http.StatusOK,
		},
		{
			name:        "service error",
			queryParams: "?page=1&limit=10",
			setupMocks: func(svc *mockService) {
				svc.On("GetList", mock.Anything, mock.AnythingOfType("overtime.OvertimeFilter")).
					Return([]OvertimeListResponse(nil), (*response.Meta)(nil), errors.New("db error"))
			},
			setupContext: func(at *testutil.APITest) {
				at.WithAuthContext(&infrastructure.MyClaims{
					UserID:      1,
					CompanyID:   1,
					Permissions: []string{constants.VIEW_OVERTIME},
				})
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := new(mockService)
			tt.setupMocks(svc)
			handler := NewHandler(svc)

			at := testutil.NewAPITest(t, http.MethodGet, "/api/overtime"+tt.queryParams, nil)
			tt.setupContext(at)

			rec, err := at.Execute(handler.GetAll)
			require.NoError(t, err)
			assert.Equal(t, tt.wantStatus, rec.Code)
		})
	}
}

func TestHandler_GetDetail(t *testing.T) {
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
				svc.On("GetDetail", mock.Anything, uint(1)).Return(&OvertimeDetailResponse{ID: 1}, nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name:       "not found",
			pathParams: map[string]string{"id": "99"},
			setupMocks: func(svc *mockService) {
				svc.On("GetDetail", mock.Anything, uint(99)).Return(nil, errors.New("not found"))
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

			at := testutil.NewAPITest(t, http.MethodGet, "/api/overtime/:id", nil)
			at.WithPathParams(tt.pathParams)

			rec, err := at.Execute(handler.GetDetail)
			require.NoError(t, err)
			assert.Equal(t, tt.wantStatus, rec.Code)
		})
	}
}

func TestHandler_ProcessAction(t *testing.T) {
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
			body: ActionRequest{
				Action:          "APPROVE",
				RejectionReason: "",
			},
			setupMocks: func(svc *mockService) {
				svc.On("ProcessAction", mock.Anything, mock.AnythingOfType("*overtime.ActionRequest")).Return(nil)
			},
			setupContext: func(at *testutil.APITest) {
				at.WithAuthContext(&infrastructure.MyClaims{
					UserID:      10,
					CompanyID:   1,
					Permissions: []string{constants.VIEW_OVERTIME},
				})
			},
			wantStatus: http.StatusOK,
		},
		{
			name:       "reject success",
			pathParams: map[string]string{"id": "2"},
			body: ActionRequest{
				Action:          "REJECT",
				RejectionReason: "Not eligible",
			},
			setupMocks: func(svc *mockService) {
				svc.On("ProcessAction", mock.Anything, mock.AnythingOfType("*overtime.ActionRequest")).Return(nil)
			},
			setupContext: func(at *testutil.APITest) {
				at.WithAuthContext(&infrastructure.MyClaims{
					UserID:      10,
					CompanyID:   1,
					Permissions: []string{constants.VIEW_OVERTIME},
				})
			},
			wantStatus: http.StatusOK,
		},
		{
			name:       "service error",
			pathParams: map[string]string{"id": "3"},
			body: ActionRequest{
				Action: "APPROVE",
			},
			setupMocks: func(svc *mockService) {
				svc.On("ProcessAction", mock.Anything, mock.AnythingOfType("*overtime.ActionRequest")).Return(errors.New("not pending"))
			},
			setupContext: func(at *testutil.APITest) {
				at.WithAuthContext(&infrastructure.MyClaims{
					UserID:      10,
					CompanyID:   1,
					Permissions: []string{constants.VIEW_OVERTIME},
				})
			},
			wantStatus: http.StatusInternalServerError,
		},
		{
			name:       "missing action",
			pathParams: map[string]string{"id": "1"},
			body:       ActionRequest{},
			setupMocks: func(svc *mockService) {},
			setupContext: func(at *testutil.APITest) {
				at.WithAuthContext(&infrastructure.MyClaims{
					UserID:      10,
					CompanyID:   1,
					Permissions: []string{constants.VIEW_OVERTIME},
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

			at := testutil.NewAPITest(t, http.MethodPost, "/api/overtime/:id/action", tt.body)
			at.WithPathParams(tt.pathParams)
			if tt.setupContext != nil {
				tt.setupContext(at)
			}

			rec, err := at.Execute(handler.ProcessAction)
			require.NoError(t, err)
			assert.Equal(t, tt.wantStatus, rec.Code)
		})
	}
}

func TestHandler_Export(t *testing.T) {
	tests := []struct {
		name         string
		setupMocks   func(*mockService)
		setupContext func(*testutil.APITest)
		wantStatus   int
	}{
		{
			name: "success",
			setupMocks: func(svc *mockService) {
				svc.On("Export", mock.Anything, mock.AnythingOfType("overtime.OvertimeFilter")).
					Return([]byte("fake-excel"), nil)
			},
			setupContext: func(at *testutil.APITest) {
				at.WithAuthContext(&infrastructure.MyClaims{
					UserID:      1,
					CompanyID:   1,
					Permissions: []string{constants.VIEW_OVERTIME},
				})
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "service error",
			setupMocks: func(svc *mockService) {
				svc.On("Export", mock.Anything, mock.AnythingOfType("overtime.OvertimeFilter")).
					Return([]byte(nil), errors.New("db error"))
			},
			setupContext: func(at *testutil.APITest) {
				at.WithAuthContext(&infrastructure.MyClaims{
					UserID:      1,
					CompanyID:   1,
					Permissions: []string{constants.VIEW_OVERTIME},
				})
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := new(mockService)
			tt.setupMocks(svc)
			handler := NewHandler(svc)

			at := testutil.NewAPITest(t, http.MethodGet, "/api/overtime/export", nil)
			tt.setupContext(at)

			rec, err := at.Execute(handler.Export)
			require.NoError(t, err)
			assert.Equal(t, tt.wantStatus, rec.Code)

			if tt.wantStatus == http.StatusOK {
				assert.Equal(t, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", rec.Header().Get("Content-Type"))
			}
		})
	}
}
