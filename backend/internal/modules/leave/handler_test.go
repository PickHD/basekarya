package leave

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

func TestHandler_Apply(t *testing.T) {
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
			body: ApplyRequest{
				LeaveTypeID: 1,
				StartDate:   "2026-06-01",
				EndDate:     "2026-06-02",
				Reason:      "Family event",
			},
			setupMocks: func(svc *mockService) {
				svc.On("Apply", mock.Anything, mock.AnythingOfType("*leave.ApplyRequest")).Return(nil)
			},
			setupContext: func(at *testutil.APITest) {
				at.WithAuthContext(&infrastructure.MyClaims{
					UserID:     1,
					CompanyID:  1,
					EmployeeID: &employeeID,
					Permissions: []string{constants.VIEW_LEAVE},
				})
			},
			wantStatus: http.StatusCreated,
		},
		{
			name: "service error",
			body: ApplyRequest{
				LeaveTypeID: 1,
				StartDate:   "2026-06-01",
				EndDate:     "2026-06-02",
				Reason:      "Family event",
			},
			setupMocks: func(svc *mockService) {
				svc.On("Apply", mock.Anything, mock.AnythingOfType("*leave.ApplyRequest")).Return(errors.New("insufficient balance"))
			},
			setupContext: func(at *testutil.APITest) {
				at.WithAuthContext(&infrastructure.MyClaims{
					UserID:     1,
					CompanyID:  1,
					EmployeeID: &employeeID,
					Permissions: []string{constants.VIEW_LEAVE},
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

			at := testutil.NewAPITest(t, http.MethodPost, "/api/leave", tt.body)
			tt.setupContext(at)

			rec, err := at.Execute(handler.Apply)
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
				svc.On("GetList", mock.Anything, mock.AnythingOfType("*leave.LeaveFilter")).
					Return([]LeaveRequestListResponse{}, (*response.Meta)(nil), nil)
			},
			setupContext: func(at *testutil.APITest) {
				at.WithAuthContext(&infrastructure.MyClaims{
					UserID:      1,
					CompanyID:   1,
					Permissions: []string{constants.VIEW_LEAVE},
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

			at := testutil.NewAPITest(t, http.MethodGet, "/api/leave"+tt.queryParams, nil)
			tt.setupContext(at)

			rec, err := at.Execute(handler.GetAll)
			require.NoError(t, err)
			assert.Equal(t, tt.wantStatus, rec.Code)
		})
	}
}

func TestHandler_GetDetail(t *testing.T) {
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
				svc.On("GetDetail", mock.Anything, uint(1)).Return(&LeaveRequestDetailResponse{ID: 1}, nil)
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

			at := testutil.NewAPITest(t, http.MethodGet, "/api/leave/:id", nil)
			at.WithPathParams(tt.pathParams)

			rec, err := at.Execute(handler.GetDetail)
			require.NoError(t, err)
			assert.Equal(t, tt.wantStatus, rec.Code)
		})
	}
}

func TestHandler_RequestAction(t *testing.T) {
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
			body: LeaveActionRequest{
				Action:          "APPROVE",
				RejectionReason: "",
			},
			setupMocks: func(svc *mockService) {
				svc.On("RequestAction", mock.Anything, mock.AnythingOfType("*leave.LeaveActionRequest")).Return(nil)
			},
			setupContext: func(at *testutil.APITest) {
				at.WithAuthContext(&infrastructure.MyClaims{
					UserID:      10,
					CompanyID:   1,
					Permissions: []string{constants.VIEW_LEAVE},
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

			at := testutil.NewAPITest(t, http.MethodPost, "/api/leave/:id/action", tt.body)
			at.WithPathParams(tt.pathParams)
			if tt.setupContext != nil {
				tt.setupContext(at)
			}

			rec, err := at.Execute(handler.RequestAction)
			require.NoError(t, err)
			assert.Equal(t, tt.wantStatus, rec.Code)
		})
	}
}
