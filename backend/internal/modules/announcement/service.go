package announcement

import (
	"basekarya-backend/pkg/constants"
	"context"
)

type Service interface {
	Publish(ctx context.Context, req *CreateAnnouncementRequest) error
}

type service struct {
	userProvider  UserProvider
	notifProvider NotificationProvider
}

func NewService(userProvider UserProvider, notifProvider NotificationProvider) Service {
	return &service{
		userProvider:  userProvider,
		notifProvider: notifProvider,
	}
}

func (s *service) Publish(ctx context.Context, req *CreateAnnouncementRequest) error {
	userIDs, err := s.userProvider.FindAllUserIDs(ctx)
	if err != nil {
		return err
	}

	err = s.notifProvider.BlastNotification(userIDs, string(constants.NotificationTypeAnnouncement), req.Title, req.Body, 0)
	if err != nil {
		return err
	}
	return nil
}
