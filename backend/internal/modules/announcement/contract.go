package announcement

import "context"

type UserProvider interface {
	FindAllUserIDs(ctx context.Context) ([]uint, error)
}

type NotificationProvider interface {
	BlastNotification(userIDs []uint,
		Type string,
		Title string,
		Message string,
		relatedID uint) error
}
