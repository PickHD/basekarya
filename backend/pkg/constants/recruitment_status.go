package constants

// Requisition statuses
const (
	RequisitionStatusDraft    = "DRAFT"
	RequisitionStatusPending  = "PENDING"
	RequisitionStatusApproved = "APPROVED"
	RequisitionStatusRejected = "REJECTED"
	RequisitionStatusClosed   = "CLOSED"
)

// Requisition priorities
const (
	RequisitionPriorityLow    = "LOW"
	RequisitionPriorityMedium = "MEDIUM"
	RequisitionPriorityHigh   = "HIGH"
	RequisitionPriorityUrgent = "URGENT"
)

// Applicant stages
const (
	ApplicantStageScreening  = "SCREENING"
	ApplicantStageInterview  = "INTERVIEW"
	ApplicantStageOffering   = "OFFERING"
	ApplicantStageHired      = "HIRED"
	ApplicantStageRejected   = "REJECTED"
)
