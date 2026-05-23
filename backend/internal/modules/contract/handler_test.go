package contract

import (
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

func TestHandler_Upsert(t *testing.T) {
	tests := []struct {
		name       string
		body       interface{}
		setupMocks func(*mockService)
		wantStatus int
	}{
		{
			name: "success",
			body: UpsertContractRequest{
				EmployeeID:     1,
				ContractType:   constants.ContractTypePKWT,
				ContractNumber: "CTR-001",
				StartDate:      "2026-01-01",
				EndDate:        "2026-12-31",
			},
			setupMocks: func(svc *mockService) {
				svc.On("Upsert", mock.Anything, mock.AnythingOfType("*contract.UpsertContractRequest")).Return(nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "service error",
			body: UpsertContractRequest{
				EmployeeID:     1,
				ContractType:   constants.ContractTypePKWT,
				ContractNumber: "CTR-001",
				StartDate:      "2026-01-01",
				EndDate:        "2026-12-31",
			},
			setupMocks: func(svc *mockService) {
				svc.On("Upsert", mock.Anything, mock.AnythingOfType("*contract.UpsertContractRequest")).Return(errors.New("db error"))
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := new(mockService)
			tt.setupMocks(svc)
			handler := NewHandler(svc)

			at := testutil.NewAPITest(t, http.MethodPost, "/api/contract", tt.body)
			at.WithAuthContext(&infrastructure.MyClaims{UserID: 1, CompanyID: 1})

			rec, err := at.Execute(handler.Upsert)
			require.NoError(t, err)
			assert.Equal(t, tt.wantStatus, rec.Code)
		})
	}
}

func TestHandler_GetAll(t *testing.T) {
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
				svc.On("GetList", mock.Anything, mock.AnythingOfType("*contract.ContractFilter")).
					Return([]ContractListResponse{}, (*response.Meta)(nil), nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name:        "service error",
			queryParams: "?page=1&limit=10",
			setupMocks: func(svc *mockService) {
				svc.On("GetList", mock.Anything, mock.AnythingOfType("*contract.ContractFilter")).
					Return([]ContractListResponse(nil), (*response.Meta)(nil), errors.New("db error"))
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := new(mockService)
			tt.setupMocks(svc)
			handler := NewHandler(svc)

			at := testutil.NewAPITest(t, http.MethodGet, "/api/contract"+tt.queryParams, nil)
			at.WithAuthContext(&infrastructure.MyClaims{UserID: 1, CompanyID: 1})

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
				svc.On("GetDetail", mock.Anything, uint(1)).Return(&ContractDetailResponse{ID: 1}, nil)
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := new(mockService)
			tt.setupMocks(svc)
			handler := NewHandler(svc)

			at := testutil.NewAPITest(t, http.MethodGet, "/api/contract/:id", nil)
			at.WithPathParams(tt.pathParams)

			rec, err := at.Execute(handler.GetDetail)
			require.NoError(t, err)
			assert.Equal(t, tt.wantStatus, rec.Code)
		})
	}
}

func TestHandler_GetByEmployee(t *testing.T) {
	tests := []struct {
		name       string
		pathParams map[string]string
		setupMocks func(*mockService)
		wantStatus int
	}{
		{
			name:       "success",
			pathParams: map[string]string{"employeeId": "1"},
			setupMocks: func(svc *mockService) {
				svc.On("GetByEmployeeID", mock.Anything, uint(1)).Return(&ContractDetailResponse{ID: 1}, nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name:       "not found",
			pathParams: map[string]string{"employeeId": "99"},
			setupMocks: func(svc *mockService) {
				svc.On("GetByEmployeeID", mock.Anything, uint(99)).Return(nil, errors.New("not found"))
			},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := new(mockService)
			tt.setupMocks(svc)
			handler := NewHandler(svc)

			at := testutil.NewAPITest(t, http.MethodGet, "/api/contract/employee/:employeeId", nil)
			at.WithPathParams(tt.pathParams)

			rec, err := at.Execute(handler.GetByEmployee)
			require.NoError(t, err)
			assert.Equal(t, tt.wantStatus, rec.Code)
		})
	}
}

func TestHandler_Delete(t *testing.T) {
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
				svc.On("Delete", mock.Anything, uint(1)).Return(nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name:       "error",
			pathParams: map[string]string{"id": "99"},
			setupMocks: func(svc *mockService) {
				svc.On("Delete", mock.Anything, uint(99)).Return(errors.New("not found"))
			},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := new(mockService)
			tt.setupMocks(svc)
			handler := NewHandler(svc)

			at := testutil.NewAPITest(t, http.MethodDelete, "/api/contract/:id", nil)
			at.WithPathParams(tt.pathParams)

			rec, err := at.Execute(handler.Delete)
			require.NoError(t, err)
			assert.Equal(t, tt.wantStatus, rec.Code)
		})
	}
}

func TestHandler_Export(t *testing.T) {
	tests := []struct {
		name        string
		queryParams string
		setupMocks  func(*mockService)
		wantStatus  int
	}{
		{
			name:        "success",
			queryParams: "?contract_type=PKWT",
			setupMocks: func(svc *mockService) {
				svc.On("Export", mock.Anything, mock.AnythingOfType("*contract.ContractFilter")).Return([]byte("fake-excel"), nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name:        "service error",
			queryParams: "",
			setupMocks: func(svc *mockService) {
				svc.On("Export", mock.Anything, mock.AnythingOfType("*contract.ContractFilter")).Return(nil, errors.New("db error"))
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := new(mockService)
			tt.setupMocks(svc)
			handler := NewHandler(svc)

			at := testutil.NewAPITest(t, http.MethodGet, "/api/contract/export"+tt.queryParams, nil)
			at.WithAuthContext(&infrastructure.MyClaims{UserID: 1, CompanyID: 1})

			rec, err := at.Execute(handler.Export)
			require.NoError(t, err)
			assert.Equal(t, tt.wantStatus, rec.Code)

			if tt.wantStatus == http.StatusOK {
				assert.Equal(t, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", rec.Header().Get("Content-Type"))
			}
		})
	}
}
