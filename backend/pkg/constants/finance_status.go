package constants

type FinanceStatus string

const (
	FinanceStatusPending  FinanceStatus = "PENDING"
	FinanceStatusApproved FinanceStatus = "APPROVED"
	FinanceStatusRejected FinanceStatus = "REJECTED"
)
