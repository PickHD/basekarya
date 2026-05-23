package payroll

import (
	"encoding/json"
	"errors"
	"net/http"
	"testing"
	"time"

	"basekarya-backend/internal/infrastructure"
	"basekarya-backend/internal/testutil"
	"basekarya-backend/pkg/response"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestHandler_Generate(t *testing.T) {
	tests := []struct {
		name       string
		body       interface{}
		setupMocks func(*mockService)
		wantStatus int
	}{
		{
			name: "success",
			body: GenerateRequest{Month: 6, Year: 2025},
			setupMocks: func(svc *mockService) {
				svc.On("GenerateAll", mock.Anything, mock.AnythingOfType("*payroll.GenerateRequest")).Return(&GenerateResponse{
					SuccessCount: 1, Month: 6, Year: 2025,
				}, nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "service error",
			body: GenerateRequest{Month: 6, Year: 2025},
			setupMocks: func(svc *mockService) {
				svc.On("GenerateAll", mock.Anything, mock.AnythingOfType("*payroll.GenerateRequest")).Return(nil, errors.New("db error"))
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := new(mockService)
			tt.setupMocks(svc)
			handler := NewHandler(svc)

			at := testutil.NewAPITest(t, http.MethodPost, "/api/payroll/generate", tt.body)
			employeeID := uint(1)
			at.WithAuthContext(&infrastructure.MyClaims{
				UserID:     1,
				CompanyID:  1,
				EmployeeID: &employeeID,
			})

			rec, err := at.Execute(handler.Generate)
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

func TestHandler_GetList(t *testing.T) {
	tests := []struct {
		name       string
		query      string
		setupMocks func(*mockService)
		wantStatus int
	}{
		{
			name:  "success",
			query: "?month=6&year=2025&page=1&limit=10",
			setupMocks: func(svc *mockService) {
				svc.On("GetList", mock.Anything, mock.AnythingOfType("*payroll.PayrollFilter")).Return([]PayrollListResponse{
					{ID: 1, EmployeeName: "John Doe", NetSalary: 5000000, Status: "DRAFT"},
				}, &response.Meta{Page: 1, Limit: 10, TotalData: 1}, nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name:  "service error",
			query: "?month=6&year=2025&page=1&limit=10",
			setupMocks: func(svc *mockService) {
				svc.On("GetList", mock.Anything, mock.AnythingOfType("*payroll.PayrollFilter")).Return([]PayrollListResponse(nil), (*response.Meta)(nil), errors.New("db error"))
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := new(mockService)
			tt.setupMocks(svc)
			handler := NewHandler(svc)

			at := testutil.NewAPITest(t, http.MethodGet, "/api/payroll"+tt.query, nil)
			employeeID := uint(1)
			at.WithAuthContext(&infrastructure.MyClaims{
				UserID:     1,
				CompanyID:  1,
				EmployeeID: &employeeID,
			})

			rec, err := at.Execute(handler.GetList)
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
				svc.On("GetDetail", mock.Anything, uint(1)).Return(&PayrollDetailResponse{
					ID: 1, EmployeeName: "John Doe", NetSalary: 5000000, Status: "DRAFT",
				}, nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name:       "service error",
			pathParams: map[string]string{"id": "1"},
			setupMocks: func(svc *mockService) {
				svc.On("GetDetail", mock.Anything, uint(1)).Return(nil, errors.New("not found"))
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

			at := testutil.NewAPITest(t, http.MethodGet, "/api/payroll/:id", nil)
			at.WithPathParams(tt.pathParams)

			rec, err := at.Execute(handler.GetDetail)
			require.NoError(t, err)
			assert.Equal(t, tt.wantStatus, rec.Code)

			var resp map[string]interface{}
			require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
		})
	}
}

func TestHandler_MarkAsPaid(t *testing.T) {
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
				svc.On("MarkAsPaid", mock.Anything, uint(1)).Return(nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name:       "service error",
			pathParams: map[string]string{"id": "1"},
			setupMocks: func(svc *mockService) {
				svc.On("MarkAsPaid", mock.Anything, uint(1)).Return(errors.New("db error"))
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

			at := testutil.NewAPITest(t, http.MethodPut, "/api/payroll/:id/paid", nil)
			at.WithPathParams(tt.pathParams)

			rec, err := at.Execute(handler.MarkAsPaid)
			require.NoError(t, err)
			assert.Equal(t, tt.wantStatus, rec.Code)

			var resp map[string]interface{}
			require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
		})
	}
}

func TestHandler_DownloadPayslipPDF(t *testing.T) {
	tests := []struct {
		name       string
		pathParams map[string]string
		setupMocks func(*mockService)
		wantStatus int
	}{
		{
			name:       "invalid id",
			pathParams: map[string]string{"id": "abc"},
			setupMocks: func(svc *mockService) {},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "service error",
			pathParams: map[string]string{"id": "1"},
			setupMocks: func(svc *mockService) {
				svc.On("GeneratePayslipPDF", mock.Anything, uint(1)).Return(nil, nil, errors.New("not found"))
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := new(mockService)
			tt.setupMocks(svc)
			handler := NewHandler(svc)

			at := testutil.NewAPITest(t, http.MethodGet, "/api/payroll/:id/payslip-pdf", nil)
			at.WithPathParams(tt.pathParams)

			rec, err := at.Execute(handler.DownloadPayslipPDF)
			require.NoError(t, err)
			assert.Equal(t, tt.wantStatus, rec.Code)
		})
	}
}

func TestHandler_BlastPayslipEmail(t *testing.T) {
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
				svc.On("BlastPayslipEmail", mock.Anything, uint(1)).Return(nil)
			},
			wantStatus: http.StatusOK,
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

			at := testutil.NewAPITest(t, http.MethodPost, "/api/payroll/:id/blast-email", nil)
			at.WithPathParams(tt.pathParams)

			rec, err := at.Execute(handler.BlastPayslipEmail)
			require.NoError(t, err)
			assert.Equal(t, tt.wantStatus, rec.Code)

			time.Sleep(50 * time.Millisecond)
		})
	}
}
