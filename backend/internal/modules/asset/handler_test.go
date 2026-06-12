package asset

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

func TestHandler_CreateCategory(t *testing.T) {
	tests := []struct {
		name       string
		body       interface{}
		setupMocks func(*mockService)
		wantStatus int
	}{
		{
			name: "success",
			body: CreateAssetCategoryRequest{Name: "Laptop", Description: "Company laptops"},
			setupMocks: func(svc *mockService) {
				svc.On("CreateCategory", mock.Anything, mock.AnythingOfType("*asset.CreateAssetCategoryRequest")).Return(nil)
			},
			wantStatus: http.StatusCreated,
		},
		{
			name: "service error",
			body: CreateAssetCategoryRequest{Name: "Laptop"},
			setupMocks: func(svc *mockService) {
				svc.On("CreateCategory", mock.Anything, mock.AnythingOfType("*asset.CreateAssetCategoryRequest")).Return(errors.New("db error"))
			},
			wantStatus: http.StatusInternalServerError,
		},
		{
			name: "validation error missing name",
			body: CreateAssetCategoryRequest{Description: "desc"},
			setupMocks: func(svc *mockService) {},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := new(mockService)
			tt.setupMocks(svc)
			handler := NewHandler(svc)

			at := testutil.NewAPITest(t, http.MethodPost, "/api/assets/categories", tt.body)
			rec, err := at.Execute(handler.CreateCategory)
			require.NoError(t, err)
			assert.Equal(t, tt.wantStatus, rec.Code)
		})
	}
}

func TestHandler_GetAllCategories(t *testing.T) {
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
				svc.On("GetCategories", mock.Anything, mock.AnythingOfType("asset.AssetCategoryFilter")).
					Return([]AssetCategoryResponse{}, (*response.Meta)(nil), nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name:        "success with default pagination",
			queryParams: "",
			setupMocks: func(svc *mockService) {
				svc.On("GetCategories", mock.Anything, mock.MatchedBy(func(f AssetCategoryFilter) bool {
					return f.Page == 1 && f.Limit == 10
				})).Return([]AssetCategoryResponse{}, (*response.Meta)(nil), nil)
			},
			wantStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := new(mockService)
			tt.setupMocks(svc)
			handler := NewHandler(svc)

			at := testutil.NewAPITest(t, http.MethodGet, "/api/assets/categories"+tt.queryParams, nil)
			rec, err := at.Execute(handler.GetAllCategories)
			require.NoError(t, err)
			assert.Equal(t, tt.wantStatus, rec.Code)
		})
	}
}

func TestHandler_CreateAsset(t *testing.T) {
	tests := []struct {
		name       string
		body       interface{}
		setupMocks func(*mockService)
		wantStatus int
	}{
		{
			name: "success",
			body: CreateAssetRequest{AssetCategoryID: 1, Name: "MacBook Pro", Condition: constants.AssetConditionGood},
			setupMocks: func(svc *mockService) {
				svc.On("CreateAsset", mock.Anything, mock.AnythingOfType("*asset.CreateAssetRequest")).Return(nil)
			},
			wantStatus: http.StatusCreated,
		},
		{
			name: "validation error missing category",
			body: CreateAssetRequest{Name: "MacBook Pro"},
			setupMocks: func(svc *mockService) {},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := new(mockService)
			tt.setupMocks(svc)
			handler := NewHandler(svc)

			at := testutil.NewAPITest(t, http.MethodPost, "/api/assets", tt.body)
			rec, err := at.Execute(handler.CreateAsset)
			require.NoError(t, err)
			assert.Equal(t, tt.wantStatus, rec.Code)
		})
	}
}

func TestHandler_GetAllAssets(t *testing.T) {
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
				svc.On("GetAssets", mock.Anything, mock.AnythingOfType("asset.AssetFilter")).
					Return([]AssetListResponse{}, (*response.Meta)(nil), nil)
			},
			wantStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := new(mockService)
			tt.setupMocks(svc)
			handler := NewHandler(svc)

			at := testutil.NewAPITest(t, http.MethodGet, "/api/assets"+tt.queryParams, nil)
			rec, err := at.Execute(handler.GetAllAssets)
			require.NoError(t, err)
			assert.Equal(t, tt.wantStatus, rec.Code)
		})
	}
}

func TestHandler_GetAssetDetail(t *testing.T) {
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
				svc.On("GetAssetDetail", mock.Anything, uint(1)).Return(&AssetDetailResponse{ID: 1, Name: "MacBook Pro"}, nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name:       "not found",
			pathParams: map[string]string{"id": "99"},
			setupMocks: func(svc *mockService) {
				svc.On("GetAssetDetail", mock.Anything, uint(99)).Return(nil, errors.New("not found"))
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

			at := testutil.NewAPITest(t, http.MethodGet, "/api/assets/:id", nil)
			at.WithPathParams(tt.pathParams)
			rec, err := at.Execute(handler.GetAssetDetail)
			require.NoError(t, err)
			assert.Equal(t, tt.wantStatus, rec.Code)
		})
	}
}

func TestHandler_UpdateAsset(t *testing.T) {
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
			body:       UpdateAssetRequest{Condition: constants.AssetConditionDamaged},
			setupMocks: func(svc *mockService) {
				svc.On("UpdateAsset", mock.Anything, mock.AnythingOfType("*asset.UpdateAssetRequest")).Return(nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name:       "service error",
			pathParams: map[string]string{"id": "1"},
			body:       UpdateAssetRequest{Name: "Updated"},
			setupMocks: func(svc *mockService) {
				svc.On("UpdateAsset", mock.Anything, mock.AnythingOfType("*asset.UpdateAssetRequest")).Return(errors.New("db error"))
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := new(mockService)
			tt.setupMocks(svc)
			handler := NewHandler(svc)

			at := testutil.NewAPITest(t, http.MethodPut, "/api/assets/:id", tt.body)
			at.WithPathParams(tt.pathParams)
			rec, err := at.Execute(handler.UpdateAsset)
			require.NoError(t, err)
			assert.Equal(t, tt.wantStatus, rec.Code)
		})
	}
}

func TestHandler_CreateAssignment(t *testing.T) {
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
			body: CreateAssetAssignmentRequest{
				AssetID:            1,
				Purpose:            "Need for presentation",
				ExpectedReturnDate: "2025-12-31",
			},
			setupMocks: func(svc *mockService) {
				svc.On("CreateAssignment", mock.Anything, mock.AnythingOfType("*asset.CreateAssetAssignmentRequest")).Return(nil)
			},
			setupContext: func(at *testutil.APITest) {
				at.WithAuthContext(&infrastructure.MyClaims{
					UserID:      1,
					CompanyID:   1,
					EmployeeID:  &employeeID,
					Permissions: []string{constants.CREATE_ASSET},
				})
			},
			wantStatus: http.StatusCreated,
		},
		{
			name: "service error",
			body: CreateAssetAssignmentRequest{AssetID: 1, Purpose: "Need"},
			setupMocks: func(svc *mockService) {
				svc.On("CreateAssignment", mock.Anything, mock.AnythingOfType("*asset.CreateAssetAssignmentRequest")).Return(errors.New("asset not available"))
			},
			setupContext: func(at *testutil.APITest) {
				at.WithAuthContext(&infrastructure.MyClaims{
					UserID:      1,
					CompanyID:   1,
					EmployeeID:  &employeeID,
					Permissions: []string{constants.CREATE_ASSET},
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

			at := testutil.NewAPITest(t, http.MethodPost, "/api/assets/assignments", tt.body)
			tt.setupContext(at)

			rec, err := at.Execute(handler.CreateAssignment)
			require.NoError(t, err)
			assert.Equal(t, tt.wantStatus, rec.Code)
		})
	}
}

func TestHandler_GetAllAssignments(t *testing.T) {
	tests := []struct {
		name         string
		queryParams  string
		setupMocks   func(*mockService)
		setupContext func(*testutil.APITest)
		wantStatus   int
	}{
		{
			name:        "success with view asset permission",
			queryParams: "?page=1&limit=10",
			setupMocks: func(svc *mockService) {
				svc.On("GetAssignments", mock.Anything, mock.AnythingOfType("asset.AssetAssignmentFilter")).
					Return([]AssetAssignmentListResponse{}, (*response.Meta)(nil), nil)
			},
			setupContext: func(at *testutil.APITest) {
				at.WithAuthContext(&infrastructure.MyClaims{
					UserID:      1,
					CompanyID:   1,
					Permissions: []string{constants.VIEW_ASSET},
				})
			},
			wantStatus: http.StatusOK,
		},
		{
			name:        "success with self asset permission only",
			queryParams: "?page=1&limit=10",
			setupMocks: func(svc *mockService) {
				svc.On("GetAssignments", mock.Anything, mock.MatchedBy(func(f AssetAssignmentFilter) bool {
					return f.UserID == 1
				})).Return([]AssetAssignmentListResponse{}, (*response.Meta)(nil), nil)
			},
			setupContext: func(at *testutil.APITest) {
				at.WithAuthContext(&infrastructure.MyClaims{
					UserID:      1,
					CompanyID:   1,
					Permissions: []string{constants.VIEW_SELF_ASSET},
				})
			},
			wantStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := new(mockService)
			tt.setupMocks(svc)
			handler := NewHandler(svc)

			at := testutil.NewAPITest(t, http.MethodGet, "/api/assets/assignments"+tt.queryParams, nil)
			tt.setupContext(at)

			rec, err := at.Execute(handler.GetAllAssignments)
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
			body:       ActionRequest{Action: "APPROVE"},
			setupMocks: func(svc *mockService) {
				svc.On("ProcessAction", mock.Anything, mock.AnythingOfType("*asset.ActionRequest")).Return(nil)
			},
			setupContext: func(at *testutil.APITest) {
				at.WithAuthContext(&infrastructure.MyClaims{
					UserID:      10,
					CompanyID:   1,
					Permissions: []string{constants.APPROVAL_ASSET},
				})
			},
			wantStatus: http.StatusOK,
		},
		{
			name:       "reject success",
			pathParams: map[string]string{"id": "2"},
			body:       ActionRequest{Action: "REJECT", RejectionReason: "Not available"},
			setupMocks: func(svc *mockService) {
				svc.On("ProcessAction", mock.Anything, mock.AnythingOfType("*asset.ActionRequest")).Return(nil)
			},
			setupContext: func(at *testutil.APITest) {
				at.WithAuthContext(&infrastructure.MyClaims{
					UserID:      10,
					CompanyID:   1,
					Permissions: []string{constants.APPROVAL_ASSET},
				})
			},
			wantStatus: http.StatusOK,
		},
		{
			name:       "service error",
			pathParams: map[string]string{"id": "1"},
			body:       ActionRequest{Action: "APPROVE"},
			setupMocks: func(svc *mockService) {
				svc.On("ProcessAction", mock.Anything, mock.AnythingOfType("*asset.ActionRequest")).Return(errors.New("cannot process"))
			},
			setupContext: func(at *testutil.APITest) {
				at.WithAuthContext(&infrastructure.MyClaims{
					UserID:      10,
					CompanyID:   1,
					Permissions: []string{constants.APPROVAL_ASSET},
				})
			},
			wantStatus: http.StatusInternalServerError,
		},
		{
			name:       "invalid id",
			pathParams: map[string]string{"id": "abc"},
			body:       ActionRequest{Action: "APPROVE"},
			setupMocks: func(svc *mockService) {},
			setupContext: func(at *testutil.APITest) {
				at.WithAuthContext(&infrastructure.MyClaims{
					UserID:      10,
					CompanyID:   1,
					Permissions: []string{constants.APPROVAL_ASSET},
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

			at := testutil.NewAPITest(t, http.MethodPut, "/api/assets/assignments/:id/action", tt.body)
			at.WithPathParams(tt.pathParams)
			tt.setupContext(at)

			rec, err := at.Execute(handler.ProcessAction)
			require.NoError(t, err)
			assert.Equal(t, tt.wantStatus, rec.Code)
		})
	}
}

func TestHandler_ProcessReturn(t *testing.T) {
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
				svc.On("ProcessReturn", mock.Anything, mock.AnythingOfType("*asset.ReturnRequest")).Return(nil)
			},
			setupContext: func(at *testutil.APITest) {
				at.WithAuthContext(&infrastructure.MyClaims{
					UserID:      10,
					CompanyID:   1,
					Permissions: []string{constants.APPROVAL_ASSET},
				})
			},
			wantStatus: http.StatusOK,
		},
		{
			name:       "service error",
			pathParams: map[string]string{"id": "1"},
			setupMocks: func(svc *mockService) {
				svc.On("ProcessReturn", mock.Anything, mock.AnythingOfType("*asset.ReturnRequest")).Return(errors.New("cannot return"))
			},
			setupContext: func(at *testutil.APITest) {
				at.WithAuthContext(&infrastructure.MyClaims{
					UserID:      10,
					CompanyID:   1,
					Permissions: []string{constants.APPROVAL_ASSET},
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

			at := testutil.NewAPITest(t, http.MethodPut, "/api/assets/assignments/:id/return", nil)
			at.WithPathParams(tt.pathParams)
			tt.setupContext(at)

			rec, err := at.Execute(handler.ProcessReturn)
			require.NoError(t, err)
			assert.Equal(t, tt.wantStatus, rec.Code)
		})
	}
}

func TestHandler_ExportAssets(t *testing.T) {
	tests := []struct {
		name         string
		queryParams  string
		setupMocks   func(*mockService)
		setupContext func(*testutil.APITest)
		wantStatus   int
	}{
		{
			name:        "success",
			queryParams: "?status=AVAILABLE",
			setupMocks: func(svc *mockService) {
				svc.On("Export", mock.Anything, mock.AnythingOfType("asset.AssetFilter")).Return([]byte("fake-excel"), nil)
			},
			setupContext: func(at *testutil.APITest) {
				at.WithAuthContext(&infrastructure.MyClaims{
					UserID:      1,
					CompanyID:   1,
					Permissions: []string{constants.EXPORT_ASSET},
				})
			},
			wantStatus: http.StatusOK,
		},
		{
			name:        "service error",
			queryParams: "",
			setupMocks: func(svc *mockService) {
				svc.On("Export", mock.Anything, mock.AnythingOfType("asset.AssetFilter")).Return(nil, errors.New("db error"))
			},
			setupContext: func(at *testutil.APITest) {
				at.WithAuthContext(&infrastructure.MyClaims{
					UserID:      1,
					CompanyID:   1,
					Permissions: []string{constants.EXPORT_ASSET},
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

			at := testutil.NewAPITest(t, http.MethodGet, "/api/assets/export"+tt.queryParams, nil)
			tt.setupContext(at)

			rec, err := at.Execute(handler.ExportAssets)
			require.NoError(t, err)
			assert.Equal(t, tt.wantStatus, rec.Code)

			if tt.wantStatus == http.StatusOK {
				assert.Equal(t, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", rec.Header().Get("Content-Type"))
				assert.Contains(t, rec.Header().Get("Content-Disposition"), "attachment; filename=assets.xlsx")
			}
		})
	}
}

func TestHandler_JSONResponse(t *testing.T) {
	svc := new(mockService)
	svc.On("GetAssetDetail", mock.Anything, uint(1)).Return(&AssetDetailResponse{
		ID:           1,
		Name:         "MacBook Pro",
		Condition:    constants.AssetConditionGood,
		Status:       constants.AssetStatusAvailable,
		CategoryName: "Laptop",
	}, nil)
	handler := NewHandler(svc)

	at := testutil.NewAPITest(t, http.MethodGet, "/api/assets/:id", nil)
	at.WithPathParams(map[string]string{"id": "1"})

	rec, err := at.Execute(handler.GetAssetDetail)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	assert.NotNil(t, resp["data"])
	data := resp["data"].(map[string]interface{})
	assert.Equal(t, "MacBook Pro", data["name"])
	assert.Equal(t, "GOOD", data["condition"])
	assert.Equal(t, "Laptop", data["category_name"])
}
