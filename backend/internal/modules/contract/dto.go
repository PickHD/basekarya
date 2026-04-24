package contract

import (
	"basekarya-backend/pkg/constants"
	"time"
)

type UpsertContractRequest struct {
	EmployeeID       uint                   `json:"employee_id" validate:"required"`
	ContractType     constants.ContractType `json:"contract_type" validate:"required"`
	ContractNumber   string                 `json:"contract_number"`
	StartDate        string                 `json:"start_date" validate:"required"`
	EndDate          string                 `json:"end_date"`
	Notes            string                 `json:"notes"`
	AttachmentBase64 string                 `json:"attachment_base64"`
}

type ContractFilter struct {
	ContractType       string
	ExpiringWithinDays int
	Page               int
	Limit              int
	Search             string
}

type ContractListResponse struct {
	ID             uint                   `json:"id"`
	EmployeeID     uint                   `json:"employee_id"`
	EmployeeName   string                 `json:"employee_name"`
	EmployeeNIK    string                 `json:"employee_nik"`
	ContractType   constants.ContractType `json:"contract_type"`
	ContractNumber string                 `json:"contract_number"`
	StartDate      time.Time              `json:"start_date"`
	EndDate        *time.Time             `json:"end_date"`
	CreatedAt      time.Time              `json:"created_at"`
}

type ContractDetailResponse struct {
	ID             uint                   `json:"id"`
	EmployeeID     uint                   `json:"employee_id"`
	EmployeeName   string                 `json:"employee_name"`
	EmployeeNIK    string                 `json:"employee_nik"`
	ContractType   constants.ContractType `json:"contract_type"`
	ContractNumber string                 `json:"contract_number"`
	StartDate      time.Time              `json:"start_date"`
	EndDate        *time.Time             `json:"end_date"`
	Notes          string                 `json:"notes"`
	AttachmentURL  string                 `json:"attachment_url"`
	CreatedAt      time.Time              `json:"created_at"`
}
