package loan

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
			body: LoanRequest{
				TotalAmount:       5000000,
				InstallmentAmount: 500000,
				Reason:            "Emergency",
			},
			setupMocks: func(svc *mockService) {
				svc.On("Create", mock.Anything, mock.AnythingOfType("*loan.LoanRequest")).Return(nil)
			},
			setupContext: func(at *testutil.APITest) {
				at.WithAuthContext(&infrastructure.MyClaims{
					UserID:      1,
					CompanyID:   1,
					EmployeeID:  &employeeID,
					Permissions: []string{constants.VIEW_LOAN},
				})
			},
			wantStatus: http.StatusCreated,
		},
		{
			name: "service error",
			body: LoanRequest{
				TotalAmount:       5000000,
				InstallmentAmount: 500000,
				Reason:            "Emergency",
			},
			setupMocks: func(svc *mockService) {
				svc.On("Create", mock.Anything, mock.AnythingOfType("*loan.LoanRequest")).Return(errors.New("users still have loan"))
			},
			setupContext: func(at *testutil.APITest) {
				at.WithAuthContext(&infrastructure.MyClaims{
					UserID:      1,
					CompanyID:   1,
					EmployeeID:  &employeeID,
					Permissions: []string{constants.VIEW_LOAN},
				})
			},
			wantStatus: http.StatusInternalServerError,
		},
		{
			name: "validation error missing required fields",
			body: LoanRequest{
				Reason: "Emergency",
			},
			setupMocks: func(svc *mockService) {},
			setupContext: func(at *testutil.APITest) {
				at.WithAuthContext(&infrastructure.MyClaims{
					UserID:      1,
					CompanyID:   1,
					EmployeeID:  &employeeID,
					Permissions: []string{constants.VIEW_LOAN},
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

			at := testutil.NewAPITest(t, http.MethodPost, "/api/loan", tt.body)
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
			name:        "success with view loan permission",
			queryParams: "?page=1&limit=10",
			setupMocks: func(svc *mockService) {
				svc.On("GetLoans", mock.Anything, mock.AnythingOfType("loan.LoanFilter")).
					Return([]LoanListResponse{}, (*response.Meta)(nil), nil)
			},
			setupContext: func(at *testutil.APITest) {
				at.WithAuthContext(&infrastructure.MyClaims{
					UserID:      1,
					CompanyID:   1,
					Permissions: []string{constants.VIEW_LOAN},
				})
			},
			wantStatus: http.StatusOK,
		},
		{
			name:        "success with self loan permission only",
			queryParams: "?page=1&limit=10",
			setupMocks: func(svc *mockService) {
				svc.On("GetLoans", mock.Anything, mock.MatchedBy(func(f LoanFilter) bool {
					return f.UserID == 1
				})).Return([]LoanListResponse{}, (*response.Meta)(nil), nil)
			},
			setupContext: func(at *testutil.APITest) {
				at.WithAuthContext(&infrastructure.MyClaims{
					UserID:      1,
					CompanyID:   1,
					Permissions: []string{constants.VIEW_SELF_LOAN},
				})
			},
			wantStatus: http.StatusOK,
		},
		{
			name:        "success with default pagination",
			queryParams: "",
			setupMocks: func(svc *mockService) {
				svc.On("GetLoans", mock.Anything, mock.MatchedBy(func(f LoanFilter) bool {
					return f.Page == 1 && f.Limit == 10
				})).Return([]LoanListResponse{}, (*response.Meta)(nil), nil)
			},
			setupContext: func(at *testutil.APITest) {
				at.WithAuthContext(&infrastructure.MyClaims{
					UserID:      1,
					CompanyID:   1,
					Permissions: []string{constants.VIEW_LOAN},
				})
			},
			wantStatus: http.StatusOK,
		},
		{
			name:        "service error",
			queryParams: "?page=1&limit=10",
			setupMocks: func(svc *mockService) {
				svc.On("GetLoans", mock.Anything, mock.AnythingOfType("loan.LoanFilter")).
					Return([]LoanListResponse(nil), (*response.Meta)(nil), errors.New("db error"))
			},
			setupContext: func(at *testutil.APITest) {
				at.WithAuthContext(&infrastructure.MyClaims{
					UserID:      1,
					CompanyID:   1,
					Permissions: []string{constants.VIEW_LOAN},
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

			at := testutil.NewAPITest(t, http.MethodGet, "/api/loan"+tt.queryParams, nil)
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
				svc.On("GetLoanDetail", mock.Anything, uint(1)).Return(&LoanDetailResponse{ID: 1}, nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name:       "not found",
			pathParams: map[string]string{"id": "99"},
			setupMocks: func(svc *mockService) {
				svc.On("GetLoanDetail", mock.Anything, uint(99)).Return(nil, errors.New("not found"))
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

			at := testutil.NewAPITest(t, http.MethodGet, "/api/loan/:id", nil)
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
				Action: "APPROVE",
			},
			setupMocks: func(svc *mockService) {
				svc.On("ProcessAction", mock.Anything, mock.AnythingOfType("*loan.ActionRequest")).Return(nil)
			},
			setupContext: func(at *testutil.APITest) {
				at.WithAuthContext(&infrastructure.MyClaims{
					UserID:      10,
					CompanyID:   1,
					Permissions: []string{constants.VIEW_LOAN},
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
				svc.On("ProcessAction", mock.Anything, mock.AnythingOfType("*loan.ActionRequest")).Return(nil)
			},
			setupContext: func(at *testutil.APITest) {
				at.WithAuthContext(&infrastructure.MyClaims{
					UserID:      10,
					CompanyID:   1,
					Permissions: []string{constants.VIEW_LOAN},
				})
			},
			wantStatus: http.StatusOK,
		},
		{
			name:       "service error",
			pathParams: map[string]string{"id": "1"},
			body: ActionRequest{
				Action: "APPROVE",
			},
			setupMocks: func(svc *mockService) {
				svc.On("ProcessAction", mock.Anything, mock.AnythingOfType("*loan.ActionRequest")).Return(errors.New("cannot process loan"))
			},
			setupContext: func(at *testutil.APITest) {
				at.WithAuthContext(&infrastructure.MyClaims{
					UserID:      10,
					CompanyID:   1,
					Permissions: []string{constants.VIEW_LOAN},
				})
			},
			wantStatus: http.StatusInternalServerError,
		},
		{
			name:       "validation error missing action",
			pathParams: map[string]string{"id": "1"},
			body: map[string]interface{}{
				"rejection_reason": "",
			},
			setupMocks: func(svc *mockService) {},
			setupContext: func(at *testutil.APITest) {
				at.WithAuthContext(&infrastructure.MyClaims{
					UserID:      10,
					CompanyID:   1,
					Permissions: []string{constants.VIEW_LOAN},
				})
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "invalid id",
			pathParams: map[string]string{"id": "abc"},
			body: ActionRequest{
				Action: "APPROVE",
			},
			setupMocks: func(svc *mockService) {},
			setupContext: func(at *testutil.APITest) {
				at.WithAuthContext(&infrastructure.MyClaims{
					UserID:      10,
					CompanyID:   1,
					Permissions: []string{constants.VIEW_LOAN},
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

			at := testutil.NewAPITest(t, http.MethodPost, "/api/loan/:id/action", tt.body)
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
		queryParams  string
		setupMocks   func(*mockService)
		setupContext func(*testutil.APITest)
		wantStatus   int
	}{
		{
			name:        "success",
			queryParams: "?status=PENDING",
			setupMocks: func(svc *mockService) {
				svc.On("Export", mock.Anything, mock.AnythingOfType("loan.LoanFilter")).
					Return([]byte("fake-excel"), nil)
			},
			setupContext: func(at *testutil.APITest) {
				at.WithAuthContext(&infrastructure.MyClaims{
					UserID:      1,
					CompanyID:   1,
					Permissions: []string{constants.VIEW_LOAN},
				})
			},
			wantStatus: http.StatusOK,
		},
		{
			name:        "success with self loan permission",
			queryParams: "",
			setupMocks: func(svc *mockService) {
				svc.On("Export", mock.Anything, mock.MatchedBy(func(f LoanFilter) bool {
					return f.UserID == 1
				})).Return([]byte("fake-excel"), nil)
			},
			setupContext: func(at *testutil.APITest) {
				at.WithAuthContext(&infrastructure.MyClaims{
					UserID:      1,
					CompanyID:   1,
					Permissions: []string{constants.VIEW_SELF_LOAN},
				})
			},
			wantStatus: http.StatusOK,
		},
		{
			name:        "service error",
			queryParams: "",
			setupMocks: func(svc *mockService) {
				svc.On("Export", mock.Anything, mock.AnythingOfType("loan.LoanFilter")).
					Return(nil, errors.New("db error"))
			},
			setupContext: func(at *testutil.APITest) {
				at.WithAuthContext(&infrastructure.MyClaims{
					UserID:      1,
					CompanyID:   1,
					Permissions: []string{constants.VIEW_LOAN},
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

			at := testutil.NewAPITest(t, http.MethodGet, "/api/loan/export"+tt.queryParams, nil)
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
