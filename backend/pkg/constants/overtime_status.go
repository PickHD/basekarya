package constants

type OvertimeStatus string

const (
	OvertimeStatusPending  OvertimeStatus = "PENDING"
	OvertimeStatusApproved OvertimeStatus = "APPROVED"
	OvertimeStatusRejected OvertimeStatus = "REJECTED"
	OvertimeStatusPaid     OvertimeStatus = "PAID"
)
