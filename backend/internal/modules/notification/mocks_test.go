package notification

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type mockRepo struct{ mock.Mock }

func (m *mockRepo) Create(ctx context.Context, notification *Notification) error {
	return m.Called(ctx, notification).Error(0)
}

func (m *mockRepo) CreateBatch(ctx context.Context, notifications []*Notification) error {
	return m.Called(ctx, notifications).Error(0)
}

func (m *mockRepo) FindByID(ctx context.Context, id uint) (*Notification, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Notification), args.Error(1)
}

func (m *mockRepo) FindAllByUserID(ctx context.Context, userID uint) ([]Notification, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]Notification), args.Error(1)
}

func (m *mockRepo) MarkAsRead(ctx context.Context, id uint) error {
	return m.Called(ctx, id).Error(0)
}

func (m *mockRepo) DeleteReadOlderThan(ctx context.Context, days int) error {
	return m.Called(ctx, days).Error(0)
}

type mockService struct{ mock.Mock }

func (m *mockService) SendNotification(ctx context.Context, userID uint, notifType string, title string, message string, relatedID uint) error {
	return m.Called(ctx, userID, notifType, title, message, relatedID).Error(0)
}

func (m *mockService) GetList(ctx context.Context, userID uint) ([]NotificationListResponse, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]NotificationListResponse), args.Error(1)
}

func (m *mockService) MarkAsRead(ctx context.Context, id uint) error {
	return m.Called(ctx, id).Error(0)
}

func (m *mockService) DeleteReadOlderThan(days int) error {
	return m.Called(days).Error(0)
}

func (m *mockService) BlastNotification(ctx context.Context, userIDs []uint, notifType string, title string, message string, relatedID uint) error {
	return m.Called(ctx, userIDs, notifType, title, message, relatedID).Error(0)
}
