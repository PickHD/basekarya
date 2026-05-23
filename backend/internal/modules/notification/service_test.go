package notification

import (
	"errors"
	"testing"

	"basekarya-backend/internal/testutil"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestService_GetList(t *testing.T) {
	ctx := testutil.CtxWithTenant(1, 1, false)

	tests := []struct {
		name       string
		setupMocks func(*mockRepo)
		wantErr    bool
	}{
		{
			name: "success with data",
			setupMocks: func(repo *mockRepo) {
				repo.On("FindAllByUserID", mock.Anything, uint(1)).Return([]Notification{
					{ID: 1, UserID: 1, Type: "LEAVE", Title: "Test", Message: "Msg", IsRead: false},
				}, nil)
			},
			wantErr: false,
		},
		{
			name: "success empty",
			setupMocks: func(repo *mockRepo) {
				repo.On("FindAllByUserID", mock.Anything, uint(1)).Return([]Notification{}, nil)
			},
			wantErr: false,
		},
		{
			name: "repo error",
			setupMocks: func(repo *mockRepo) {
				repo.On("FindAllByUserID", mock.Anything, uint(1)).Return([]Notification(nil), errors.New("db error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(mockRepo)

			tt.setupMocks(repo)

			svc := NewService(nil, repo)
			_, err := svc.GetList(ctx, 1)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestService_MarkAsRead(t *testing.T) {
	ctx := testutil.CtxWithTenant(1, 1, false)

	tests := []struct {
		name       string
		setupMocks func(*mockRepo)
		wantErr    bool
	}{
		{
			name: "success",
			setupMocks: func(repo *mockRepo) {
				repo.On("FindByID", mock.Anything, uint(1)).Return(&Notification{ID: 1, IsRead: false}, nil)
				repo.On("MarkAsRead", mock.Anything, uint(1)).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "not found",
			setupMocks: func(repo *mockRepo) {
				repo.On("FindByID", mock.Anything, uint(99)).Return(nil, errors.New("not found"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(mockRepo)
			tt.setupMocks(repo)

			svc := NewService(nil, repo)

			id := uint(1)
			if tt.name == "not found" {
				id = 99
			}
			err := svc.MarkAsRead(ctx, id)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestService_DeleteReadOlderThan(t *testing.T) {
	repo := new(mockRepo)
	repo.On("DeleteReadOlderThan", mock.Anything, 30).Return(nil)

	svc := NewService(nil, repo)
	err := svc.DeleteReadOlderThan(30)

	require.NoError(t, err)
}
