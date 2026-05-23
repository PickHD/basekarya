package announcement

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type mockUserProvider struct{ mock.Mock }

func (m *mockUserProvider) FindAllUserIDs(ctx context.Context) ([]uint, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]uint), args.Error(1)
}

type mockNotificationProvider struct{ mock.Mock }

func (m *mockNotificationProvider) BlastNotification(ctx context.Context, userIDs []uint, notifType string, title string, message string, relatedID uint) error {
	args := m.Called(ctx, userIDs, notifType, title, message, relatedID)
	return args.Error(0)
}

type mockService struct{ mock.Mock }

func (m *mockService) Publish(ctx context.Context, req *CreateAnnouncementRequest) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}
