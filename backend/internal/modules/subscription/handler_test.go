package subscription

import (
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"basekarya-backend/internal/infrastructure"
	"basekarya-backend/internal/testutil"
	"basekarya-backend/pkg/constants"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestHandler_ListPlans(t *testing.T) {
	tests := []struct {
		name       string
		setupMocks func(*mockService)
		wantStatus int
	}{
		{
			name: "success",
			setupMocks: func(svc *mockService) {
				svc.On("ListPlans", mock.Anything).Return([]PlanResponse{
					{ID: 1, Name: "Basic", Slug: "basic"},
				}, nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "service error",
			setupMocks: func(svc *mockService) {
				svc.On("ListPlans", mock.Anything).Return(nil, errors.New("db error"))
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := new(mockService)
			tt.setupMocks(svc)
			handler := NewHandler(svc)

			at := testutil.NewAPITest(t, http.MethodGet, "/api/subscription/plans", nil)

			rec, err := at.Execute(handler.ListPlans)
			require.NoError(t, err)
			assert.Equal(t, tt.wantStatus, rec.Code)
		})
	}
}

func TestHandler_RequestUpgrade(t *testing.T) {
	tests := []struct {
		name         string
		body         interface{}
		setupMocks   func(*mockService)
		setupContext func(*testutil.APITest)
		wantStatus   int
	}{
		{
			name: "success",
			body: UpgradeRequest{PlanSlug: "pro"},
			setupMocks: func(svc *mockService) {
				svc.On("RequestUpgrade", mock.Anything, mock.AnythingOfType("*subscription.UpgradeRequest")).Return(&UpgradeResponse{ID: 1, RequestedPlanID: 2, Status: "PENDING"}, nil)
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
			name: "invalid plan slug",
			body: UpgradeRequest{PlanSlug: ""},
			setupMocks: func(svc *mockService) {},
			setupContext: func(at *testutil.APITest) {
				at.WithAuthContext(&infrastructure.MyClaims{
					UserID:    1,
					CompanyID: 1,
				})
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "service error",
			body: UpgradeRequest{PlanSlug: "pro"},
			setupMocks: func(svc *mockService) {
				svc.On("RequestUpgrade", mock.Anything, mock.AnythingOfType("*subscription.UpgradeRequest")).Return(nil, errors.New("already on this plan"))
			},
			setupContext: func(at *testutil.APITest) {
				at.WithAuthContext(&infrastructure.MyClaims{
					UserID:    1,
					CompanyID: 1,
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

			at := testutil.NewAPITest(t, http.MethodPost, "/api/subscription/upgrade", tt.body)
			tt.setupContext(at)

			rec, err := at.Execute(handler.RequestUpgrade)
			require.NoError(t, err)
			assert.Equal(t, tt.wantStatus, rec.Code)
		})
	}
}

func TestHandler_ListPendingRequests(t *testing.T) {
	tests := []struct {
		name         string
		setupMocks   func(*mockService)
		setupContext func(*testutil.APITest)
		wantStatus   int
	}{
		{
			name: "success",
			setupMocks: func(svc *mockService) {
				svc.On("ListPendingRequests", mock.Anything).Return([]SubscriptionRequestResponse{
					{ID: 1, CompanyName: "Acme", Status: "PENDING"},
				}, nil)
			},
			setupContext: func(at *testutil.APITest) {
				at.WithAuthContext(&infrastructure.MyClaims{
					UserID:          1,
					CompanyID:       1,
					IsPlatformAdmin: true,
				})
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "service error",
			setupMocks: func(svc *mockService) {
				svc.On("ListPendingRequests", mock.Anything).Return(nil, errors.New("db error"))
			},
			setupContext: func(at *testutil.APITest) {
				at.WithAuthContext(&infrastructure.MyClaims{
					UserID:          1,
					CompanyID:       1,
					IsPlatformAdmin: true,
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

			at := testutil.NewAPITest(t, http.MethodGet, "/api/subscription/requests/pending", nil)
			tt.setupContext(at)

			rec, err := at.Execute(handler.ListPendingRequests)
			require.NoError(t, err)
			assert.Equal(t, tt.wantStatus, rec.Code)
		})
	}
}

func TestHandler_ReviewRequest(t *testing.T) {
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
			body:       ReviewRequest{Status: constants.SubReqStatusApproved, Notes: "ok"},
			setupMocks: func(svc *mockService) {
				svc.On("ReviewRequest", mock.Anything, uint(1), mock.AnythingOfType("*subscription.ReviewRequest")).Return(nil)
			},
			setupContext: func(at *testutil.APITest) {
				at.WithAuthContext(&infrastructure.MyClaims{
					UserID:          10,
					CompanyID:       1,
					IsPlatformAdmin: true,
				})
			},
			wantStatus: http.StatusOK,
		},
		{
			name:       "invalid request id",
			pathParams: map[string]string{"id": "abc"},
			body:       ReviewRequest{Status: constants.SubReqStatusApproved},
			setupMocks: func(svc *mockService) {},
			setupContext: func(at *testutil.APITest) {
				at.WithAuthContext(&infrastructure.MyClaims{
					UserID:          10,
					CompanyID:       1,
					IsPlatformAdmin: true,
				})
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "invalid status",
			pathParams: map[string]string{"id": "1"},
			body:       ReviewRequest{Status: "INVALID"},
			setupMocks: func(svc *mockService) {},
			setupContext: func(at *testutil.APITest) {
				at.WithAuthContext(&infrastructure.MyClaims{
					UserID:          10,
					CompanyID:       1,
					IsPlatformAdmin: true,
				})
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "service error",
			pathParams: map[string]string{"id": "1"},
			body:       ReviewRequest{Status: constants.SubReqStatusApproved},
			setupMocks: func(svc *mockService) {
				svc.On("ReviewRequest", mock.Anything, uint(1), mock.AnythingOfType("*subscription.ReviewRequest")).Return(errors.New("request already reviewed"))
			},
			setupContext: func(at *testutil.APITest) {
				at.WithAuthContext(&infrastructure.MyClaims{
					UserID:          10,
					CompanyID:       1,
					IsPlatformAdmin: true,
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

			at := testutil.NewAPITest(t, http.MethodPost, "/api/subscription/requests/:id/review", tt.body)
			at.WithPathParams(tt.pathParams)
			tt.setupContext(at)

			rec, err := at.Execute(handler.ReviewRequest)
			require.NoError(t, err)
			assert.Equal(t, tt.wantStatus, rec.Code)
		})
	}
}

func TestHandler_ListCompanies(t *testing.T) {
	tests := []struct {
		name         string
		setupMocks   func(*mockService)
		setupContext func(*testutil.APITest)
		wantStatus   int
	}{
		{
			name: "success",
			setupMocks: func(svc *mockService) {
				svc.On("ListCompanies", mock.Anything, "").Return([]CompanyListItem{
					{ID: 1, Name: "Acme"},
				}, nil)
			},
			setupContext: func(at *testutil.APITest) {
				at.WithAuthContext(&infrastructure.MyClaims{
					UserID:          1,
					CompanyID:       1,
					IsPlatformAdmin: true,
				})
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "service error",
			setupMocks: func(svc *mockService) {
				svc.On("ListCompanies", mock.Anything, "").Return(nil, errors.New("db error"))
			},
			setupContext: func(at *testutil.APITest) {
				at.WithAuthContext(&infrastructure.MyClaims{
					UserID:          1,
					CompanyID:       1,
					IsPlatformAdmin: true,
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

			at := testutil.NewAPITest(t, http.MethodGet, "/api/subscription/companies", nil)
			tt.setupContext(at)

			rec, err := at.Execute(handler.ListCompanies)
			require.NoError(t, err)
			assert.Equal(t, tt.wantStatus, rec.Code)
		})
	}
}

func TestHandler_GetCompanyDetail(t *testing.T) {
	tests := []struct {
		name         string
		pathParams   map[string]string
		setupMocks   func(*mockService)
		wantStatus   int
	}{
		{
			name:       "success",
			pathParams: map[string]string{"id": "1"},
			setupMocks: func(svc *mockService) {
				svc.On("GetCompanyDetail", mock.Anything, uint(1)).Return(&CompanyDetail{ID: 1, Name: "Acme"}, nil)
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
			pathParams: map[string]string{"id": "99"},
			setupMocks: func(svc *mockService) {
				svc.On("GetCompanyDetail", mock.Anything, uint(99)).Return(nil, errors.New("not found"))
			},
			wantStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := new(mockService)
			tt.setupMocks(svc)
			handler := NewHandler(svc)

			at := testutil.NewAPITest(t, http.MethodGet, "/api/subscription/companies/:id", nil)
			at.WithPathParams(tt.pathParams)

			rec, err := at.Execute(handler.GetCompanyDetail)
			require.NoError(t, err)
			assert.Equal(t, tt.wantStatus, rec.Code)
		})
	}
}

func TestHandler_UpdateCompanyStatus(t *testing.T) {
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
			body:       UpdateCompanyStatusRequest{SubscriptionStatus: constants.SubStatusActive},
			setupMocks: func(svc *mockService) {
				svc.On("UpdateCompanyStatus", mock.Anything, uint(1), mock.AnythingOfType("*subscription.UpdateCompanyStatusRequest")).Return(nil)
			},
			setupContext: func(at *testutil.APITest) {
				at.WithAuthContext(&infrastructure.MyClaims{
					UserID:          1,
					CompanyID:       1,
					IsPlatformAdmin: true,
				})
			},
			wantStatus: http.StatusOK,
		},
		{
			name:       "invalid id",
			pathParams: map[string]string{"id": "abc"},
			body:       UpdateCompanyStatusRequest{SubscriptionStatus: constants.SubStatusActive},
			setupMocks: func(svc *mockService) {},
			setupContext: func(at *testutil.APITest) {
				at.WithAuthContext(&infrastructure.MyClaims{
					UserID:          1,
					CompanyID:       1,
					IsPlatformAdmin: true,
				})
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "invalid status",
			pathParams: map[string]string{"id": "1"},
			body:       UpdateCompanyStatusRequest{SubscriptionStatus: "INVALID"},
			setupMocks: func(svc *mockService) {},
			setupContext: func(at *testutil.APITest) {
				at.WithAuthContext(&infrastructure.MyClaims{
					UserID:          1,
					CompanyID:       1,
					IsPlatformAdmin: true,
				})
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "service error",
			pathParams: map[string]string{"id": "1"},
			body:       UpdateCompanyStatusRequest{SubscriptionStatus: constants.SubStatusExpired},
			setupMocks: func(svc *mockService) {
				svc.On("UpdateCompanyStatus", mock.Anything, uint(1), mock.AnythingOfType("*subscription.UpdateCompanyStatusRequest")).Return(errors.New("db error"))
			},
			setupContext: func(at *testutil.APITest) {
				at.WithAuthContext(&infrastructure.MyClaims{
					UserID:          1,
					CompanyID:       1,
					IsPlatformAdmin: true,
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

			at := testutil.NewAPITest(t, http.MethodPut, "/api/subscription/companies/:id/status", tt.body)
			at.WithPathParams(tt.pathParams)
			tt.setupContext(at)

			rec, err := at.Execute(handler.UpdateCompanyStatus)
			require.NoError(t, err)
			assert.Equal(t, tt.wantStatus, rec.Code)
		})
	}
}

func TestHandler_GetDashboardStats(t *testing.T) {
	tests := []struct {
		name         string
		setupMocks   func(*mockService)
		setupContext func(*testutil.APITest)
		wantStatus   int
	}{
		{
			name: "success",
			setupMocks: func(svc *mockService) {
				svc.On("GetDashboardStats", mock.Anything).Return(&DashboardStatsResponse{
					TotalCompanies:      10,
					ActiveSubscriptions: 8,
					TotalRevenue:        800,
				}, nil)
			},
			setupContext: func(at *testutil.APITest) {
				at.WithAuthContext(&infrastructure.MyClaims{
					UserID:          1,
					CompanyID:       1,
					IsPlatformAdmin: true,
				})
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "service error",
			setupMocks: func(svc *mockService) {
				svc.On("GetDashboardStats", mock.Anything).Return(nil, errors.New("db error"))
			},
			setupContext: func(at *testutil.APITest) {
				at.WithAuthContext(&infrastructure.MyClaims{
					UserID:          1,
					CompanyID:       1,
					IsPlatformAdmin: true,
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

			at := testutil.NewAPITest(t, http.MethodGet, "/api/subscription/dashboard", nil)
			tt.setupContext(at)

			rec, err := at.Execute(handler.GetDashboardStats)
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
