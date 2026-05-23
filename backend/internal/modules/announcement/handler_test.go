package announcement

import (
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"basekarya-backend/internal/infrastructure"
	"basekarya-backend/internal/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestHandler_PublishAnnouncement(t *testing.T) {
	tests := []struct {
		name       string
		body       interface{}
		setupMocks func(*mockService)
		wantStatus int
	}{
		{
			name: "success",
			body: CreateAnnouncementRequest{
				Title: "Important",
				Body:  "Please read",
			},
			setupMocks: func(svc *mockService) {
				svc.On("Publish", mock.Anything, mock.AnythingOfType("*announcement.CreateAnnouncementRequest")).Return(nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "service error",
			body: CreateAnnouncementRequest{
				Title: "Important",
				Body:  "Please read",
			},
			setupMocks: func(svc *mockService) {
				svc.On("Publish", mock.Anything, mock.AnythingOfType("*announcement.CreateAnnouncementRequest")).Return(errors.New("blast failed"))
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := new(mockService)
			tt.setupMocks(svc)
			handler := NewHandler(svc)

			at := testutil.NewAPITest(t, http.MethodPost, "/api/announcement", tt.body)
			at.WithAuthContext(&infrastructure.MyClaims{
				UserID:    1,
				CompanyID: 1,
			})

			rec, err := at.Execute(handler.PublishAnnouncement)
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
