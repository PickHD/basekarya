package loan

import "context"

type NotificationProvider interface {
	SendNotification(userID uint,
		Type string,
		Title string,
		Message string,
		relatedID uint) error
	BlastNotification(userIDs []uint,
		Type string,
		Title string,
		Message string,
		relatedID uint) error
}

type UserProvider interface {
	FindApprovalUsers(ctx context.Context, permissionApprovalName string) ([]uint, error)
}
