package company

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

func TestHandler_GetProfile(t *testing.T) {
	tests := []struct {
		name         string
		setupMocks   func(*mockService)
		setupContext func(*testutil.APITest)
		wantStatus   int
	}{
		{
			name: "success",
			setupMocks: func(svc *mockService) {
				svc.On("GetProfile", mock.Anything).Return(&CompanyProfileResponse{
					ID:   1,
					Name: "Test Co",
				}, nil)
			},
			setupContext: func(at *testutil.APITest) {
				at.WithAuthContext(&infrastructure.MyClaims{
					UserID:      1,
					CompanyID:   1,
					Permissions: []string{constants.VIEW_EMPLOYEE},
				})
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "service error",
			setupMocks: func(svc *mockService) {
				svc.On("GetProfile", mock.Anything).Return(nil, errors.New("not found"))
			},
			setupContext: func(at *testutil.APITest) {
				at.WithAuthContext(&infrastructure.MyClaims{
					UserID:      1,
					CompanyID:   1,
					Permissions: []string{constants.VIEW_EMPLOYEE},
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

			at := testutil.NewAPITest(t, http.MethodGet, "/api/company/profile", nil)
			tt.setupContext(at)

			rec, err := at.Execute(handler.GetProfile)

			require.NoError(t, err)
			assert.Equal(t, tt.wantStatus, rec.Code)

			var resp map[string]interface{}
			require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
		})
	}
}
