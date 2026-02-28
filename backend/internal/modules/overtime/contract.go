package overtime

import "context"

type NotificationProvider interface {
	SendNotification(userID uint,
		Type string,
		Title string,
		Message string, relatedID uint) error
}

type UserProvider interface {
	FindAdminID(ctx context.Context) (uint, error)
}
