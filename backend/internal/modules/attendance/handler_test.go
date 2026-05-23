package attendance

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

func TestHandler_Clock(t *testing.T) {
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
			body: ClockRequest{
				Latitude:    -6.2,
				Longitude:   106.8,
				ImageBase64: "aGVsbG8=",
			},
			setupMocks: func(svc *mockService) {
				svc.On("Clock", mock.Anything, uint(1), mock.AnythingOfType("*attendance.ClockRequest")).Return(&AttendanceResponse{
					Type:    string(constants.AttendanceTypeCheckIn),
					Status:  string(constants.AttendanceStatusPresent),
					Message: "Check-in succesful",
				}, nil)
			},
			setupContext: func(at *testutil.APITest) {
				at.WithAuthContext(&infrastructure.MyClaims{
					UserID:      1,
					CompanyID:   1,
					EmployeeID:  &employeeID,
					Permissions: []string{constants.CREATE_ATTENDANCE},
				})
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "service error",
			body: ClockRequest{
				Latitude:    -6.2,
				Longitude:   106.8,
				ImageBase64: "aGVsbG8=",
			},
			setupMocks: func(svc *mockService) {
				svc.On("Clock", mock.Anything, uint(1), mock.AnythingOfType("*attendance.ClockRequest")).Return(nil, errors.New("employee data not found"))
			},
			setupContext: func(at *testutil.APITest) {
				at.WithAuthContext(&infrastructure.MyClaims{
					UserID:      1,
					CompanyID:   1,
					EmployeeID:  &employeeID,
					Permissions: []string{constants.CREATE_ATTENDANCE},
				})
			},
			wantStatus: http.StatusInternalServerError,
		},
		{
			name: "invalid request missing fields",
			body: map[string]interface{}{
				"latitude": -6.2,
			},
			setupMocks: func(svc *mockService) {},
			setupContext: func(at *testutil.APITest) {
				at.WithAuthContext(&infrastructure.MyClaims{
					UserID:      1,
					CompanyID:   1,
					EmployeeID:  &employeeID,
					Permissions: []string{constants.CREATE_ATTENDANCE},
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

			at := testutil.NewAPITest(t, http.MethodPost, "/api/attendance/clock", tt.body)
			tt.setupContext(at)

			rec, err := at.Execute(handler.Clock)
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

func TestHandler_GetTodayStatus(t *testing.T) {
	employeeID := uint(1)

	tests := []struct {
		name         string
		setupMocks   func(*mockService)
		setupContext func(*testutil.APITest)
		wantStatus   int
	}{
		{
			name: "success",
			setupMocks: func(svc *mockService) {
				svc.On("GetTodayStatus", mock.Anything, uint(1)).Return(&TodayStatusResponse{
					Status: string(constants.AttendanceStatusPresent),
					Type:   string(constants.AttendanceTypeCheckIn),
				}, nil)
			},
			setupContext: func(at *testutil.APITest) {
				at.WithAuthContext(&infrastructure.MyClaims{
					UserID:      1,
					CompanyID:   1,
					EmployeeID:  &employeeID,
					Permissions: []string{constants.VIEW_SELF_ATTENDANCE},
				})
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "service error",
			setupMocks: func(svc *mockService) {
				svc.On("GetTodayStatus", mock.Anything, uint(1)).Return(nil, errors.New("employee not found"))
			},
			setupContext: func(at *testutil.APITest) {
				at.WithAuthContext(&infrastructure.MyClaims{
					UserID:      1,
					CompanyID:   1,
					EmployeeID:  &employeeID,
					Permissions: []string{constants.VIEW_SELF_ATTENDANCE},
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

			at := testutil.NewAPITest(t, http.MethodGet, "/api/attendance/today", nil)
			tt.setupContext(at)

			rec, err := at.Execute(handler.GetTodayStatus)
			require.NoError(t, err)
			assert.Equal(t, tt.wantStatus, rec.Code)
		})
	}
}

func TestHandler_GetHistory(t *testing.T) {
	employeeID := uint(1)

	tests := []struct {
		name         string
		queryParams  string
		setupMocks   func(*mockService)
		setupContext func(*testutil.APITest)
		wantStatus   int
	}{
		{
			name:        "success",
			queryParams: "?month=5&year=2026",
			setupMocks: func(svc *mockService) {
				svc.On("GetMyHistory", mock.Anything, uint(1), 5, 2026, 10, "").Return([]Attendance{}, (*response.Meta)(nil), nil)
			},
			setupContext: func(at *testutil.APITest) {
				at.WithAuthContext(&infrastructure.MyClaims{
					UserID:      1,
					CompanyID:   1,
					EmployeeID:  &employeeID,
					Permissions: []string{constants.VIEW_SELF_ATTENDANCE},
				})
			},
			wantStatus: http.StatusOK,
		},
		{
			name:        "service error",
			queryParams: "?month=5&year=2026",
			setupMocks: func(svc *mockService) {
				svc.On("GetMyHistory", mock.Anything, uint(1), 5, 2026, 10, "").Return(nil, nil, errors.New("employee not found"))
			},
			setupContext: func(at *testutil.APITest) {
				at.WithAuthContext(&infrastructure.MyClaims{
					UserID:      1,
					CompanyID:   1,
					EmployeeID:  &employeeID,
					Permissions: []string{constants.VIEW_SELF_ATTENDANCE},
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

			at := testutil.NewAPITest(t, http.MethodGet, "/api/attendance/history"+tt.queryParams, nil)
			tt.setupContext(at)

			rec, err := at.Execute(handler.GetHistory)
			require.NoError(t, err)
			assert.Equal(t, tt.wantStatus, rec.Code)
		})
	}
}

func TestHandler_GetAllAttendanceRecap(t *testing.T) {
	tests := []struct {
		name         string
		queryParams  string
		setupMocks   func(*mockService)
		setupContext func(*testutil.APITest)
		wantStatus   int
	}{
		{
			name:        "success",
			queryParams: "?limit=10",
			setupMocks: func(svc *mockService) {
				svc.On("GetAllRecap", mock.Anything, mock.AnythingOfType("*attendance.FilterParams")).Return([]RecapResponse{}, (*response.Meta)(nil), nil)
			},
			setupContext: func(at *testutil.APITest) {
				at.WithAuthContext(&infrastructure.MyClaims{
					UserID:      1,
					CompanyID:   1,
					Permissions: []string{constants.VIEW_ATTENDANCE},
				})
			},
			wantStatus: http.StatusOK,
		},
		{
			name:        "service error",
			queryParams: "?limit=10",
			setupMocks: func(svc *mockService) {
				svc.On("GetAllRecap", mock.Anything, mock.AnythingOfType("*attendance.FilterParams")).Return(nil, nil, errors.New("db error"))
			},
			setupContext: func(at *testutil.APITest) {
				at.WithAuthContext(&infrastructure.MyClaims{
					UserID:      1,
					CompanyID:   1,
					Permissions: []string{constants.VIEW_ATTENDANCE},
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

			at := testutil.NewAPITest(t, http.MethodGet, "/api/attendance/recap"+tt.queryParams, nil)
			tt.setupContext(at)

			rec, err := at.Execute(handler.GetAllAttendanceRecap)
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
				svc.On("GetDashboardStats", mock.Anything).Return(&DashboardStatResponse{
					TotalEmployees: 100,
					PresentToday:   80,
					LateToday:      10,
					AbsentToday:    20,
				}, nil)
			},
			setupContext: func(at *testutil.APITest) {
				at.WithAuthContext(&infrastructure.MyClaims{
					UserID:      1,
					CompanyID:   1,
					Permissions: []string{constants.VIEW_ATTENDANCE},
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
					UserID:      1,
					CompanyID:   1,
					Permissions: []string{constants.VIEW_ATTENDANCE},
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

			at := testutil.NewAPITest(t, http.MethodGet, "/api/attendance/dashboard", nil)
			tt.setupContext(at)

			rec, err := at.Execute(handler.GetDashboardStats)
			require.NoError(t, err)
			assert.Equal(t, tt.wantStatus, rec.Code)
		})
	}
}

func TestHandler_ExportAttendance(t *testing.T) {
	tests := []struct {
		name         string
		setupMocks   func(*mockService)
		setupContext func(*testutil.APITest)
		wantStatus   int
	}{
		{
			name: "success",
			setupMocks: func(svc *mockService) {
				svc.On("GenerateExcel", mock.Anything, mock.AnythingOfType("*attendance.FilterParams")).Return([]byte("fake-excel-data"), nil)
			},
			setupContext: func(at *testutil.APITest) {
				at.WithAuthContext(&infrastructure.MyClaims{
					UserID:      1,
					CompanyID:   1,
					Permissions: []string{constants.EXPORT_ATTENDANCE},
				})
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "service error",
			setupMocks: func(svc *mockService) {
				svc.On("GenerateExcel", mock.Anything, mock.AnythingOfType("*attendance.FilterParams")).Return(nil, errors.New("db error"))
			},
			setupContext: func(at *testutil.APITest) {
				at.WithAuthContext(&infrastructure.MyClaims{
					UserID:      1,
					CompanyID:   1,
					Permissions: []string{constants.EXPORT_ATTENDANCE},
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

			at := testutil.NewAPITest(t, http.MethodGet, "/api/attendance/export", nil)
			tt.setupContext(at)

			rec, err := at.Execute(handler.ExportAttendance)
			require.NoError(t, err)
			assert.Equal(t, tt.wantStatus, rec.Code)
		})
	}
}
