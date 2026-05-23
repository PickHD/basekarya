package rbac

import (
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"basekarya-backend/internal/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestHandler_CreateRole(t *testing.T) {
	tests := []struct {
		name       string
		body       interface{}
		setupMocks func(*mockService)
		wantStatus int
	}{
		{
			name: "success",
			body: CreateRoleRequest{Name: "MANAGER"},
			setupMocks: func(svc *mockService) {
				svc.On("CreateRole", mock.Anything, mock.AnythingOfType("*rbac.CreateRoleRequest")).Return(nil)
			},
			wantStatus: http.StatusCreated,
		},
		{
			name: "service error",
			body: CreateRoleRequest{Name: "MANAGER"},
			setupMocks: func(svc *mockService) {
				svc.On("CreateRole", mock.Anything, mock.AnythingOfType("*rbac.CreateRoleRequest")).Return(errors.New("role already exists"))
			},
			wantStatus: http.StatusInternalServerError,
		},
		{
			name:       "missing name",
			body:       CreateRoleRequest{},
			setupMocks: func(svc *mockService) {},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := new(mockService)
			tt.setupMocks(svc)
			handler := NewHandler(svc)

			at := testutil.NewAPITest(t, http.MethodPost, "/api/roles", tt.body)
			rec, err := at.Execute(handler.CreateRole)

			require.NoError(t, err)
			assert.Equal(t, tt.wantStatus, rec.Code)
		})
	}
}

func TestHandler_GetRolePermissions(t *testing.T) {
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
				svc.On("GetRolePermissions", mock.Anything, uint(1)).Return(&RolePermissionsResponse{
					RoleID:      1,
					RoleName:    "SUPERADMIN",
					Permissions: []Permission{},
				}, nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name:       "not found",
			pathParams: map[string]string{"id": "99"},
			setupMocks: func(svc *mockService) {
				svc.On("GetRolePermissions", mock.Anything, uint(99)).Return(nil, errors.New("not found"))
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

			at := testutil.NewAPITest(t, http.MethodGet, "/api/roles/:id/permissions", nil)
			at.WithPathParams(tt.pathParams)

			rec, err := at.Execute(handler.GetRolePermissions)

			require.NoError(t, err)
			assert.Equal(t, tt.wantStatus, rec.Code)
		})
	}
}

func TestHandler_AssignPermissions(t *testing.T) {
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
			body:       AssignPermissionsRequest{PermissionIDs: []uint{1, 2}},
			setupMocks: func(svc *mockService) {
				svc.On("AssignPermissions", mock.Anything, uint(1), mock.AnythingOfType("*rbac.AssignPermissionsRequest")).Return(nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name:       "service error",
			pathParams: map[string]string{"id": "1"},
			body:       AssignPermissionsRequest{PermissionIDs: []uint{1}},
			setupMocks: func(svc *mockService) {
				svc.On("AssignPermissions", mock.Anything, uint(1), mock.AnythingOfType("*rbac.AssignPermissionsRequest")).Return(errors.New("role not found"))
			},
			wantStatus: http.StatusInternalServerError,
		},
		{
			name:       "invalid id",
			pathParams: map[string]string{"id": "abc"},
			body:       AssignPermissionsRequest{PermissionIDs: []uint{1}},
			setupMocks: func(svc *mockService) {},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "empty permission ids",
			pathParams: map[string]string{"id": "1"},
			body:       AssignPermissionsRequest{},
			setupMocks: func(svc *mockService) {},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := new(mockService)
			tt.setupMocks(svc)
			handler := NewHandler(svc)

			at := testutil.NewAPITest(t, http.MethodPost, "/api/roles/:id/permissions", tt.body)
			at.WithPathParams(tt.pathParams)

			rec, err := at.Execute(handler.AssignPermissions)

			require.NoError(t, err)
			assert.Equal(t, tt.wantStatus, rec.Code)
		})
	}
}

func TestHandler_GetAllPermissions(t *testing.T) {
	tests := []struct {
		name       string
		setupMocks func(*mockService)
		wantStatus int
	}{
		{
			name: "success",
			setupMocks: func(svc *mockService) {
				svc.On("GetAllPermissions", mock.Anything).Return([]PermissionResponse{}, nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "service error",
			setupMocks: func(svc *mockService) {
				svc.On("GetAllPermissions", mock.Anything).Return(nil, errors.New("db error"))
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := new(mockService)
			tt.setupMocks(svc)
			handler := NewHandler(svc)

			at := testutil.NewAPITest(t, http.MethodGet, "/api/permissions", nil)
			rec, err := at.Execute(handler.GetAllPermissions)

			require.NoError(t, err)
			assert.Equal(t, tt.wantStatus, rec.Code)

			var resp map[string]interface{}
			require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
		})
	}
}

func TestHandler_GetAllRoles(t *testing.T) {
	tests := []struct {
		name       string
		setupMocks func(*mockService)
		wantStatus int
	}{
		{
			name: "success",
			setupMocks: func(svc *mockService) {
				svc.On("GetAllRoles", mock.Anything).Return([]RoleResponse{
					{ID: 1, Name: "SUPERADMIN"},
				}, nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "service error",
			setupMocks: func(svc *mockService) {
				svc.On("GetAllRoles", mock.Anything).Return(nil, errors.New("db error"))
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := new(mockService)
			tt.setupMocks(svc)
			handler := NewHandler(svc)

			at := testutil.NewAPITest(t, http.MethodGet, "/api/roles", nil)
			rec, err := at.Execute(handler.GetAllRoles)

			require.NoError(t, err)
			assert.Equal(t, tt.wantStatus, rec.Code)
		})
	}
}
