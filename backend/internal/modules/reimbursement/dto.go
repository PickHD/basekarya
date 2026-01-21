package reimbursement

import (
	"mime/multipart"
	"time"
)

type ReimbursementFilter struct {
	UserID uint
	Status string
	Page   int
	Limit  int
}

type ReimbursementRequest struct {
	UserID      uint                  `form:"-"`
	Title       string                `form:"title" validate:"required,max=255"`
	Description string                `form:"description" validate:"omitempty"`
	Amount      float64               `form:"amount" validate:"required,min=1000"`
	Date        string                `form:"date" validate:"required"`
	File        *multipart.FileHeader `form:"file" validate:"required"`
}

type ActionRequest struct {
	ID              uint   `json:"-"`
	SuperAdminID    uint   `json:"-"`
	Action          string `json:"action" validate:"required"`
	RejectionReason string `json:"rejection_reason" validate:"omitempty"`
}

type ReimbursementDetailResponse struct {
	ID              uint      `json:"id"`
	Title           string    `json:"title"`
	Description     string    `json:"description"`
	Amount          float64   `json:"amount"`
	DateOfExpense   time.Time `json:"date_of_expense"`
	ProofFileURL    string    `json:"proof_file_url"`
	Status          string    `json:"status"`
	RejectionReason *string   `json:"rejection_reason"`

	RequesterName string `json:"requester_name"`
}

type ReimbursementListResponse struct {
	ID            uint      `json:"id"`
	Title         string    `json:"title"`
	Amount        float64   `json:"amount"`
	DateOfExpense time.Time `json:"date_of_expense"`
	ProofFileURL  string    `json:"proof_file_url"`
	Status        string    `json:"status"`
}
