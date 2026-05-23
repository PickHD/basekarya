package announcement

import (
	"errors"
	"testing"

	"basekarya-backend/internal/testutil"
	"basekarya-backend/pkg/constants"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func newTestAnnouncementService() (Service, *mockUserProvider, *mockNotificationProvider) {
	userProv := new(mockUserProvider)
	notifProv := new(mockNotificationProvider)
	svc := NewService(userProv, notifProv)
	return svc, userProv, notifProv
}

func TestService_Publish(t *testing.T) {
	ctx := testutil.CtxWithTenant(1, 1, false)

	tests := []struct {
		name       string
		req        *CreateAnnouncementRequest
		setupMocks func(*mockUserProvider, *mockNotificationProvider)
		wantErr    bool
		errMsg     string
	}{
		{
			name: "success",
			req: &CreateAnnouncementRequest{
				Title: "Important Announcement",
				Body:  "Please read this carefully.",
			},
			setupMocks: func(userProv *mockUserProvider, notifProv *mockNotificationProvider) {
				userProv.On("FindAllUserIDs", mock.Anything).Return([]uint{1, 2, 3}, nil)
				notifProv.On("BlastNotification", mock.Anything, []uint{1, 2, 3}, string(constants.NotificationTypeAnnouncement), "Important Announcement", "Please read this carefully.", uint(0)).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "error find all user ids fails",
			req: &CreateAnnouncementRequest{
				Title: "Test",
				Body:  "Body",
			},
			setupMocks: func(userProv *mockUserProvider, notifProv *mockNotificationProvider) {
				userProv.On("FindAllUserIDs", mock.Anything).Return(nil, errors.New("db error"))
			},
			wantErr: true,
			errMsg:  "db error",
		},
		{
			name: "error blast notification fails",
			req: &CreateAnnouncementRequest{
				Title: "Test",
				Body:  "Body",
			},
			setupMocks: func(userProv *mockUserProvider, notifProv *mockNotificationProvider) {
				userProv.On("FindAllUserIDs", mock.Anything).Return([]uint{1}, nil)
				notifProv.On("BlastNotification", mock.Anything, []uint{1}, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(errors.New("notif error"))
			},
			wantErr: true,
			errMsg:  "notif error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, userProv, notifProv := newTestAnnouncementService()
			tt.setupMocks(userProv, notifProv)

			err := svc.Publish(ctx, tt.req)

			if tt.wantErr {
				require.Error(t, err)
				assert.Equal(t, tt.errMsg, err.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}
