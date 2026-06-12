package master

import (
	"errors"
	"net/http"
	"testing"

	"basekarya-backend/internal/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestHandler_GetShifts(t *testing.T) {
	tests := []struct {
		name       string
		setupMocks func(*mockService)
		wantStatus int
	}{
		{
			name: "success",
			setupMocks: func(svc *mockService) {
				svc.On("GetAllShifts", mock.Anything).Return([]LookupResponse{
					{ID: 1, Name: "Day"},
				}, nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "service error",
			setupMocks: func(svc *mockService) {
				svc.On("GetAllShifts", mock.Anything).Return(nil, errors.New("db error"))
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := new(mockService)
			tt.setupMocks(svc)
			handler := NewHandler(svc)

			at := testutil.NewAPITest(t, http.MethodGet, "/api/master/shifts", nil)
			rec, err := at.Execute(handler.GetShifts)

			require.NoError(t, err)
			assert.Equal(t, tt.wantStatus, rec.Code)
		})
	}
}

func TestHandler_GetLeaveTypes(t *testing.T) {
	tests := []struct {
		name       string
		setupMocks func(*mockService)
		wantStatus int
	}{
		{
			name: "success",
			setupMocks: func(svc *mockService) {
				svc.On("GetAllLeaveTypes", mock.Anything).Return([]LookupLeaveTypeResponse{
					{ID: 1, Name: "Annual", DefaultQuota: 12, IsDeducted: true},
				}, nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "service error",
			setupMocks: func(svc *mockService) {
				svc.On("GetAllLeaveTypes", mock.Anything).Return(nil, errors.New("db error"))
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := new(mockService)
			tt.setupMocks(svc)
			handler := NewHandler(svc)

			at := testutil.NewAPITest(t, http.MethodGet, "/api/master/leave-types", nil)
			rec, err := at.Execute(handler.GetLeaveTypes)

			require.NoError(t, err)
			assert.Equal(t, tt.wantStatus, rec.Code)
		})
	}
}
