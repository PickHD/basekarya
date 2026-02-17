package constants

type NotificationType string

const (
	NotificationTypeApproved NotificationType = "APPROVED"
	NotificationTypeRejected NotificationType = "REJECTED"

	NotificationTypeLeaveApprovalReq     NotificationType = "LEAVE_APPROVAL_REQ"
	NotificationTypeReimburseApprovalReq NotificationType = "REIMBURSE_APPROVAL_REQ"
)
