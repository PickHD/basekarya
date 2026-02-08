package constants

type LeaveStatus string

const (
	LeaveStatusPending  LeaveStatus = "PENDING"
	LeaveStatusApproved LeaveStatus = "APPROVED"
	LeaveStatusRejected LeaveStatus = "REJECTED"
)
