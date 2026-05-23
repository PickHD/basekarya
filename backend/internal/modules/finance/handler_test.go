package finance

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

func TestHandler_CreateTransaction(t *testing.T) {
	tests := []struct {
		name         string
		body         interface{}
		setupMocks   func(*mockService)
		setupContext func(*testutil.APITest)
		wantStatus   int
	}{
		{
			name: "success",
			body: CreateTransactionRequest{
				FinanceCategoryID: 1,
				Type:              "INCOME",
				Amount:            5000000,
				TransactionDate:   "2026-01-15",
			},
			setupMocks: func(svc *mockService) {
				svc.On("CreateTransaction", mock.Anything, mock.AnythingOfType("*finance.CreateTransactionRequest")).Return(nil)
			},
			setupContext: func(at *testutil.APITest) {
				at.WithAuthContext(&infrastructure.MyClaims{
					UserID:    1,
					CompanyID: 1,
				})
			},
			wantStatus: http.StatusCreated,
		},
		{
			name: "service error",
			body: CreateTransactionRequest{
				FinanceCategoryID: 1,
				Type:              "INCOME",
				Amount:            5000000,
				TransactionDate:   "2026-01-15",
			},
			setupMocks: func(svc *mockService) {
				svc.On("CreateTransaction", mock.Anything, mock.AnythingOfType("*finance.CreateTransactionRequest")).Return(errors.New("db error"))
			},
			setupContext: func(at *testutil.APITest) {
				at.WithAuthContext(&infrastructure.MyClaims{
					UserID:    1,
					CompanyID: 1,
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

			at := testutil.NewAPITest(t, http.MethodPost, "/api/finance/transaction", tt.body)
			if tt.setupContext != nil {
				tt.setupContext(at)
			}

			rec, err := at.Execute(handler.CreateTransaction)
			require.NoError(t, err)
			assert.Equal(t, tt.wantStatus, rec.Code)
		})
	}
}

func TestHandler_GetAllTransactions(t *testing.T) {
	tests := []struct {
		name        string
		queryParams string
		setupMocks  func(*mockService)
		wantStatus  int
	}{
		{
			name:        "success",
			queryParams: "?limit=10",
			setupMocks: func(svc *mockService) {
				svc.On("GetTransactions", mock.Anything, mock.AnythingOfType("finance.TransactionFilter")).
					Return([]TransactionListResponse{}, (*response.Meta)(nil), nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name:        "service error",
			queryParams: "?limit=10",
			setupMocks: func(svc *mockService) {
				svc.On("GetTransactions", mock.Anything, mock.AnythingOfType("finance.TransactionFilter")).
					Return([]TransactionListResponse(nil), (*response.Meta)(nil), errors.New("db error"))
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := new(mockService)
			tt.setupMocks(svc)
			handler := NewHandler(svc)

			at := testutil.NewAPITest(t, http.MethodGet, "/api/finance/transaction"+tt.queryParams, nil)
			at.WithAuthContext(&infrastructure.MyClaims{UserID: 1, CompanyID: 1})

			rec, err := at.Execute(handler.GetAllTransactions)
			require.NoError(t, err)
			assert.Equal(t, tt.wantStatus, rec.Code)
		})
	}
}

func TestHandler_GetTransactionDetail(t *testing.T) {
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
				svc.On("GetTransactionDetail", mock.Anything, uint(1)).Return(&TransactionDetailResponse{ID: 1}, nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name:       "not found",
			pathParams: map[string]string{"id": "99"},
			setupMocks: func(svc *mockService) {
				svc.On("GetTransactionDetail", mock.Anything, uint(99)).Return(nil, errors.New("not found"))
			},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := new(mockService)
			tt.setupMocks(svc)
			handler := NewHandler(svc)

			at := testutil.NewAPITest(t, http.MethodGet, "/api/finance/transaction/:id", nil)
			at.WithPathParams(tt.pathParams)

			rec, err := at.Execute(handler.GetTransactionDetail)
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
			name: "approve success",
			pathParams: map[string]string{"id": "1"},
			body: ActionRequest{
				Action: "APPROVE",
			},
			setupMocks: func(svc *mockService) {
				svc.On("ProcessAction", mock.Anything, mock.AnythingOfType("*finance.ActionRequest")).Return(nil)
			},
			setupContext: func(at *testutil.APITest) {
				at.WithAuthContext(&infrastructure.MyClaims{
					UserID:    2,
					CompanyID: 1,
				})
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "service error",
			pathParams: map[string]string{"id": "1"},
			body: ActionRequest{
				Action: "APPROVE",
			},
			setupMocks: func(svc *mockService) {
				svc.On("ProcessAction", mock.Anything, mock.AnythingOfType("*finance.ActionRequest")).Return(errors.New("not pending"))
			},
			setupContext: func(at *testutil.APITest) {
				at.WithAuthContext(&infrastructure.MyClaims{
					UserID:    2,
					CompanyID: 1,
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

			at := testutil.NewAPITest(t, http.MethodPost, "/api/finance/transaction/:id/action", tt.body)
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

func TestHandler_CreateCategory(t *testing.T) {
	tests := []struct {
		name       string
		body       interface{}
		setupMocks func(*mockService)
		wantStatus int
	}{
		{
			name: "success",
			body: CategoryRequest{
				Name: "Bonus",
				Type: "INCOME",
			},
			setupMocks: func(svc *mockService) {
				svc.On("CreateCategory", mock.Anything, mock.AnythingOfType("*finance.CategoryRequest")).Return(nil)
			},
			wantStatus: http.StatusCreated,
		},
		{
			name: "service error",
			body: CategoryRequest{
				Name: "Bonus",
				Type: "INCOME",
			},
			setupMocks: func(svc *mockService) {
				svc.On("CreateCategory", mock.Anything, mock.AnythingOfType("*finance.CategoryRequest")).Return(errors.New("db error"))
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := new(mockService)
			tt.setupMocks(svc)
			handler := NewHandler(svc)

			at := testutil.NewAPITest(t, http.MethodPost, "/api/finance/category", tt.body)
			at.WithAuthContext(&infrastructure.MyClaims{UserID: 1, CompanyID: 1})

			rec, err := at.Execute(handler.CreateCategory)
			require.NoError(t, err)
			assert.Equal(t, tt.wantStatus, rec.Code)
		})
	}
}

func TestHandler_GetCategories(t *testing.T) {
	tests := []struct {
		name        string
		queryParams string
		setupMocks  func(*mockService)
		wantStatus  int
	}{
		{
			name:        "success",
			queryParams: "?type=INCOME",
			setupMocks: func(svc *mockService) {
				svc.On("GetCategories", mock.Anything, "INCOME").Return([]CategoryResponse{
					{ID: 1, Name: "Salary", Type: constants.FinanceTypeIncome},
				}, nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name:        "service error",
			queryParams: "",
			setupMocks: func(svc *mockService) {
				svc.On("GetCategories", mock.Anything, "").Return([]CategoryResponse(nil), errors.New("db error"))
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := new(mockService)
			tt.setupMocks(svc)
			handler := NewHandler(svc)

			at := testutil.NewAPITest(t, http.MethodGet, "/api/finance/category"+tt.queryParams, nil)
			at.WithAuthContext(&infrastructure.MyClaims{UserID: 1, CompanyID: 1})

			rec, err := at.Execute(handler.GetCategories)
			require.NoError(t, err)
			assert.Equal(t, tt.wantStatus, rec.Code)
		})
	}
}

func TestHandler_UpdateCategory(t *testing.T) {
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
			body: CategoryRequest{
				Name: "Updated Salary",
				Type: "INCOME",
			},
			setupMocks: func(svc *mockService) {
				svc.On("UpdateCategory", mock.Anything, uint(1), mock.AnythingOfType("*finance.CategoryRequest")).Return(nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name:       "service error",
			pathParams: map[string]string{"id": "99"},
			body: CategoryRequest{
				Name: "Test",
				Type: "INCOME",
			},
			setupMocks: func(svc *mockService) {
				svc.On("UpdateCategory", mock.Anything, uint(99), mock.AnythingOfType("*finance.CategoryRequest")).Return(errors.New("not found"))
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := new(mockService)
			tt.setupMocks(svc)
			handler := NewHandler(svc)

			at := testutil.NewAPITest(t, http.MethodPut, "/api/finance/category/:id", tt.body)
			at.WithPathParams(tt.pathParams)
			at.WithAuthContext(&infrastructure.MyClaims{UserID: 1, CompanyID: 1})

			rec, err := at.Execute(handler.UpdateCategory)
			require.NoError(t, err)
			assert.Equal(t, tt.wantStatus, rec.Code)
		})
	}
}

func TestHandler_DeleteCategory(t *testing.T) {
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
				svc.On("DeleteCategory", mock.Anything, uint(1)).Return(nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name:       "error",
			pathParams: map[string]string{"id": "99"},
			setupMocks: func(svc *mockService) {
				svc.On("DeleteCategory", mock.Anything, uint(99)).Return(errors.New("not found"))
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := new(mockService)
			tt.setupMocks(svc)
			handler := NewHandler(svc)

			at := testutil.NewAPITest(t, http.MethodDelete, "/api/finance/category/:id", nil)
			at.WithPathParams(tt.pathParams)
			at.WithAuthContext(&infrastructure.MyClaims{UserID: 1, CompanyID: 1})

			rec, err := at.Execute(handler.DeleteCategory)
			require.NoError(t, err)
			assert.Equal(t, tt.wantStatus, rec.Code)
		})
	}
}

func TestHandler_GetDashboard(t *testing.T) {
	tests := []struct {
		name        string
		queryParams string
		setupMocks  func(*mockService)
		wantStatus  int
	}{
		{
			name:        "success",
			queryParams: "?start_date=2026-01-01&end_date=2026-12-31",
			setupMocks: func(svc *mockService) {
				svc.On("GetDashboard", mock.Anything, "2026-01-01", "2026-12-31").Return(&DashboardResponse{
					TotalIncome: 10000000,
				}, nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name:        "service error",
			queryParams: "",
			setupMocks: func(svc *mockService) {
				svc.On("GetDashboard", mock.Anything, "", "").Return(nil, errors.New("db error"))
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := new(mockService)
			tt.setupMocks(svc)
			handler := NewHandler(svc)

			at := testutil.NewAPITest(t, http.MethodGet, "/api/finance/dashboard"+tt.queryParams, nil)
			at.WithAuthContext(&infrastructure.MyClaims{UserID: 1, CompanyID: 1})

			rec, err := at.Execute(handler.GetDashboard)
			require.NoError(t, err)
			assert.Equal(t, tt.wantStatus, rec.Code)

			if tt.wantStatus < 400 {
				var resp map[string]interface{}
				require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
				assert.Nil(t, resp["error"])
			}
		})
	}
}
