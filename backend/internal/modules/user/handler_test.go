package user

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

func TestHandler_GetProfile(t *testing.T) {
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
				svc.On("GetProfile", uint(1)).Return(&UserProfileResponse{
					ID:       1,
					Username: "john.doe",
					Role:     "EMPLOYEE",
					FullName: "John Doe",
				}, nil)
			},
			setupContext: func(at *testutil.APITest) {
				at.WithAuthContext(&infrastructure.MyClaims{
					UserID:     1,
					CompanyID:  1,
					EmployeeID: &employeeID,
				})
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "service error",
			setupMocks: func(svc *mockService) {
				svc.On("GetProfile", uint(1)).Return(nil, errors.New("not found"))
			},
			setupContext: func(at *testutil.APITest) {
				at.WithAuthContext(&infrastructure.MyClaims{
					UserID:     1,
					CompanyID:  1,
					EmployeeID: &employeeID,
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

			at := testutil.NewAPITest(t, http.MethodGet, "/api/user/profile", nil)
			tt.setupContext(at)

			rec, err := at.Execute(handler.GetProfile)
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

func TestHandler_UpdateProfile(t *testing.T) {
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
			body: UpdateProfileRequest{
				FullName:    "John Updated",
				PhoneNumber: "081234567890",
			},
			setupMocks: func(svc *mockService) {
				svc.On("UpdateProfile", mock.Anything, uint(1), mock.AnythingOfType("*user.UpdateProfileRequest"), mock.Anything).Return(nil)
			},
			setupContext: func(at *testutil.APITest) {
				at.WithAuthContext(&infrastructure.MyClaims{
					UserID:     1,
					CompanyID:  1,
					EmployeeID: &employeeID,
				})
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "service error",
			body: UpdateProfileRequest{
				FullName: "John Updated",
			},
			setupMocks: func(svc *mockService) {
				svc.On("UpdateProfile", mock.Anything, uint(1), mock.AnythingOfType("*user.UpdateProfileRequest"), mock.Anything).Return(errors.New("db error"))
			},
			setupContext: func(at *testutil.APITest) {
				at.WithAuthContext(&infrastructure.MyClaims{
					UserID:     1,
					CompanyID:  1,
					EmployeeID: &employeeID,
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

			at := testutil.NewAPITest(t, http.MethodPut, "/api/user/profile", tt.body)
			tt.setupContext(at)

			rec, err := at.Execute(handler.UpdateProfile)
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

func TestHandler_ChangePassword(t *testing.T) {
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
			body: ChangePasswordRequest{
				OldPassword:     "oldpass",
				NewPassword:     "newpass123",
				ConfirmPassword: "newpass123",
			},
			setupMocks: func(svc *mockService) {
				svc.On("ChangePassword", mock.Anything, uint(1), mock.AnythingOfType("*user.ChangePasswordRequest")).Return(nil)
			},
			setupContext: func(at *testutil.APITest) {
				at.WithAuthContext(&infrastructure.MyClaims{
					UserID:     1,
					CompanyID:  1,
					EmployeeID: &employeeID,
				})
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "service error",
			body: ChangePasswordRequest{
				OldPassword:     "oldpass",
				NewPassword:     "newpass123",
				ConfirmPassword: "newpass123",
			},
			setupMocks: func(svc *mockService) {
				svc.On("ChangePassword", mock.Anything, uint(1), mock.AnythingOfType("*user.ChangePasswordRequest")).Return(errors.New("invalid old password"))
			},
			setupContext: func(at *testutil.APITest) {
				at.WithAuthContext(&infrastructure.MyClaims{
					UserID:     1,
					CompanyID:  1,
					EmployeeID: &employeeID,
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

			at := testutil.NewAPITest(t, http.MethodPost, "/api/user/change-password", tt.body)
			tt.setupContext(at)

			rec, err := at.Execute(handler.ChangePassword)
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

func TestHandler_GetAllEmployees(t *testing.T) {
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
				svc.On("GetAllEmployees", mock.Anything, 1, 10, "").
					Return([]EmployeeListResponse{}, (*response.Meta)(nil), nil)
			},
			setupContext: func(at *testutil.APITest) {
				at.WithAuthContext(&infrastructure.MyClaims{
					UserID:    1,
					CompanyID: 1,
				})
			},
			wantStatus: http.StatusOK,
		},
		{
			name:        "service error",
			queryParams: "?page=1&limit=10",
			setupMocks: func(svc *mockService) {
				svc.On("GetAllEmployees", mock.Anything, 1, 10, "").
					Return([]EmployeeListResponse(nil), (*response.Meta)(nil), errors.New("db error"))
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

			at := testutil.NewAPITest(t, http.MethodGet, "/api/employees"+tt.queryParams, nil)
			tt.setupContext(at)

			rec, err := at.Execute(handler.GetAllEmployees)
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

func TestHandler_CreateEmployee(t *testing.T) {
	tests := []struct {
		name         string
		body         interface{}
		setupMocks   func(*mockService)
		setupContext func(*testutil.APITest)
		wantStatus   int
	}{
		{
			name: "success",
			body: CreateEmployeeRequest{
				FullName:     "Jane Doe",
				NIK:          "EMP002",
				DepartmentID: 1,
				ShiftID:      1,
				RoleID:       1,
				BaseSalary:   5000000,
				Email:        "jane@example.com",
				Position:     "Designer",
			},
			setupMocks: func(svc *mockService) {
				svc.On("CreateEmployee", mock.Anything, mock.AnythingOfType("*user.CreateEmployeeRequest")).
					Return(&CreateEmployeeResponse{Username: "jane.doe"}, nil)
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
			body: CreateEmployeeRequest{
				FullName:     "Jane Doe",
				NIK:          "EMP002",
				DepartmentID: 1,
				ShiftID:      1,
				RoleID:       1,
				BaseSalary:   5000000,
				Email:        "jane@example.com",
				Position:     "Designer",
			},
			setupMocks: func(svc *mockService) {
				svc.On("CreateEmployee", mock.Anything, mock.AnythingOfType("*user.CreateEmployeeRequest")).
					Return(nil, errors.New("db error"))
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

			at := testutil.NewAPITest(t, http.MethodPost, "/api/employees", tt.body)
			tt.setupContext(at)

			rec, err := at.Execute(handler.CreateEmployee)
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

func TestHandler_UpdateEmployee(t *testing.T) {
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
			body: UpdateEmployeeRequest{
				FullName: "John Updated",
				Position: "Senior Developer",
			},
			setupMocks: func(svc *mockService) {
				svc.On("UpdateEmployee", mock.Anything, uint(1), mock.AnythingOfType("*user.UpdateEmployeeRequest")).Return(nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name:       "service error",
			pathParams: map[string]string{"id": "1"},
			body: UpdateEmployeeRequest{
				FullName: "John Updated",
			},
			setupMocks: func(svc *mockService) {
				svc.On("UpdateEmployee", mock.Anything, uint(1), mock.AnythingOfType("*user.UpdateEmployeeRequest")).Return(errors.New("db error"))
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := new(mockService)
			tt.setupMocks(svc)
			handler := NewHandler(svc)

			at := testutil.NewAPITest(t, http.MethodPut, "/api/employees/:id", tt.body)
			at.WithPathParams(tt.pathParams)

			rec, err := at.Execute(handler.UpdateEmployee)
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

func TestHandler_DeleteEmployee(t *testing.T) {
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
				svc.On("DeleteEmployee", mock.Anything, uint(1)).Return(nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name:       "service error",
			pathParams: map[string]string{"id": "1"},
			setupMocks: func(svc *mockService) {
				svc.On("DeleteEmployee", mock.Anything, uint(1)).Return(errors.New("db error"))
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := new(mockService)
			tt.setupMocks(svc)
			handler := NewHandler(svc)

			at := testutil.NewAPITest(t, http.MethodDelete, "/api/employees/:id", nil)
			at.WithPathParams(tt.pathParams)

			rec, err := at.Execute(handler.DeleteEmployee)
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
