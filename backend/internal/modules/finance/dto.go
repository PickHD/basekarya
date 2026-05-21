package finance

import (
	"basekarya-backend/pkg/constants"
	"time"
)

type TransactionFilter struct {
	CreatedBy uint
	Type      string
	Status    string
	StartDate string
	EndDate   string
	Page      int
	Limit     int
}

type CreateTransactionRequest struct {
	CreatedBy         uint    `json:"-"`
	FinanceCategoryID uint    `json:"finance_category_id" validate:"required"`
	Type              string  `json:"type" validate:"required"`
	Amount            float64 `json:"amount" validate:"required"`
	Description       string  `json:"description" validate:"omitempty"`
	TransactionDate   string  `json:"transaction_date" validate:"required"`
	ReferenceNumber   string  `json:"reference_number" validate:"omitempty"`
}

type ActionRequest struct {
	ID              uint   `json:"-"`
	SuperAdminID    uint   `json:"-"`
	Action          string `json:"action" validate:"required"`
	RejectionReason string `json:"rejection_reason" validate:"omitempty"`
}

type CategoryRequest struct {
	Name        string `json:"name" validate:"required"`
	Type        string `json:"type" validate:"required"`
	Description string `json:"description" validate:"omitempty"`
}

type TransactionListResponse struct {
	ID                uint                    `json:"id"`
	CreatorName       string                  `json:"creator_name"`
	CategoryName      string                  `json:"category_name"`
	Type              constants.FinanceType   `json:"type"`
	Amount            float64                 `json:"amount"`
	TransactionDate   time.Time               `json:"transaction_date"`
	ReferenceNumber   string                  `json:"reference_number"`
	Status            constants.FinanceStatus `json:"status"`
	CreatedAt         time.Time               `json:"created_at"`
}

type TransactionDetailResponse struct {
	ID                uint                    `json:"id"`
	CreatorName       string                  `json:"creator_name"`
	CategoryName      string                  `json:"category_name"`
	CategoryType      constants.FinanceType   `json:"category_type"`
	Type              constants.FinanceType   `json:"type"`
	Amount            float64                 `json:"amount"`
	Description       string                  `json:"description"`
	TransactionDate   time.Time               `json:"transaction_date"`
	ReferenceNumber   string                  `json:"reference_number"`
	Status            constants.FinanceStatus `json:"status"`
	RejectionReason   string                  `json:"rejection_reason"`
	ApprovedBy        *uint                   `json:"approved_by"`
	ApproverName      string                  `json:"approver_name"`
	CreatedAt         time.Time               `json:"created_at"`
}

type CategoryResponse struct {
	ID          uint                  `json:"id"`
	Name        string                `json:"name"`
	Type        constants.FinanceType `json:"type"`
	Description string                `json:"description"`
	CreatedAt   time.Time             `json:"created_at"`
	UpdatedAt   time.Time             `json:"updated_at"`
}

type DashboardResponse struct {
	TotalIncome        float64                    `json:"total_income"`
	TotalExpense       float64                    `json:"total_expense"`
	NetBalance         float64                    `json:"net_balance"`
	TransactionCount   int64                      `json:"transaction_count"`
	MonthlySummary     []MonthlySummaryItem       `json:"monthly_summary"`
	CategoryBreakdown  []CategoryBreakdownItem    `json:"category_breakdown"`
	RecentTransactions []TransactionListResponse  `json:"recent_transactions"`
}

type MonthlySummaryItem struct {
	Month   string  `json:"month"`
	Income  float64 `json:"income"`
	Expense float64 `json:"expense"`
}

type CategoryBreakdownItem struct {
	CategoryName string  `json:"category_name"`
	Type         string  `json:"type"`
	Total        float64 `json:"total"`
}
