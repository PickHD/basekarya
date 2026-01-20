package constants

type ReimbursementStatus string

const (
	ReimbursementStatusPending  ReimbursementStatus = "PENDING"
	ReimbursementStatusApproved ReimbursementStatus = "APPROVED"
	ReimbursementStatusRejected ReimbursementStatus = "REJECTED"
)
