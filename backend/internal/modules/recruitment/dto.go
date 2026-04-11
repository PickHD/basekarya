package recruitment

import "time"

// ── Requisition DTOs ──────────────────────────────────────────────────────────

type CreateRequisitionRequest struct {
	DepartmentID   uint   `json:"department_id" validate:"required"`
	Title          string `json:"title" validate:"required,max=255"`
	Description    string `json:"description"`
	Quantity       int    `json:"quantity" validate:"min=1"`
	EmploymentType string `json:"employment_type" validate:"required,oneof=PKWT PKWTT"`
	Priority       string `json:"priority" validate:"required,oneof=LOW MEDIUM HIGH URGENT"`
	TargetDate     string `json:"target_date"` // format: YYYY-MM-DD
}

type RequisitionActionRequest struct {
	Action          string `json:"action" validate:"required,oneof=APPROVE REJECT"`
	RejectionReason string `json:"rejection_reason"`
}

type RequisitionFilter struct {
	Status       string
	DepartmentID uint
	Priority     string
	Search       string
	Page         int
	Limit        int
}

type RequisitionListResponse struct {
	ID             uint       `json:"id"`
	Title          string     `json:"title"`
	DepartmentID   uint       `json:"department_id"`
	DepartmentName string     `json:"department_name"`
	EmploymentType string     `json:"employment_type"`
	Quantity       int        `json:"quantity"`
	Priority       string     `json:"priority"`
	Status         string     `json:"status"`
	RequesterID    uint       `json:"requester_id"`
	RequesterName  string     `json:"requester_name"`
	TargetDate     *time.Time `json:"target_date"`
	CreatedAt      time.Time  `json:"created_at"`
}

type RequisitionDetailResponse struct {
	ID              uint       `json:"id"`
	Title           string     `json:"title"`
	Description     string     `json:"description"`
	DepartmentID    uint       `json:"department_id"`
	DepartmentName  string     `json:"department_name"`
	EmploymentType  string     `json:"employment_type"`
	Quantity        int        `json:"quantity"`
	Priority        string     `json:"priority"`
	Status          string     `json:"status"`
	RequesterID     uint       `json:"requester_id"`
	RequesterName   string     `json:"requester_name"`
	ApprovedBy      *uint      `json:"approved_by"`
	ApproverName    string     `json:"approver_name"`
	RejectionReason string     `json:"rejection_reason"`
	TargetDate      *time.Time `json:"target_date"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

// ── Applicant DTOs ────────────────────────────────────────────────────────────

type CreateApplicantRequest struct {
	FullName        string `json:"full_name" validate:"required,max=100"`
	Email           string `json:"email" validate:"required,email"`
	PhoneNumber     string `json:"phone_number"`
	ResumeBase64    string `json:"resume_base64"`
}

type UpdateApplicantStageRequest struct {
	Stage           string `json:"stage" validate:"required,oneof=SCREENING INTERVIEW OFFERING HIRED REJECTED"`
	Notes           string `json:"notes"`
	RejectionReason string `json:"rejection_reason"`
}

type ApplicantFilter struct {
	JobRequisitionID uint
	Stage            string
	Search           string
	Page             int
	Limit            int
}

type ApplicantListResponse struct {
	ID               uint      `json:"id"`
	JobRequisitionID uint      `json:"job_requisition_id"`
	FullName         string    `json:"full_name"`
	Email            string    `json:"email"`
	PhoneNumber      string    `json:"phone_number"`
	ResumeURL        string    `json:"resume_url"`
	Stage            string    `json:"stage"`
	StageOrder       int       `json:"stage_order"`
	CreatedAt        time.Time `json:"created_at"`
}

type ApplicantDetailResponse struct {
	ID               uint                      `json:"id"`
	JobRequisitionID uint                      `json:"job_requisition_id"`
	FullName         string                    `json:"full_name"`
	Email            string                    `json:"email"`
	PhoneNumber      string                    `json:"phone_number"`
	ResumeURL        string                    `json:"resume_url"`
	Stage            string                    `json:"stage"`
	Notes            string                    `json:"notes"`
	RejectionReason  string                    `json:"rejection_reason"`
	CreatedAt        time.Time                 `json:"created_at"`
	StageHistories   []StageHistoryResponse    `json:"stage_histories"`
}

type StageHistoryResponse struct {
	ID            uint      `json:"id"`
	FromStage     string    `json:"from_stage"`
	ToStage       string    `json:"to_stage"`
	ChangedByName string    `json:"changed_by_name"`
	Notes         string    `json:"notes"`
	CreatedAt     time.Time `json:"created_at"`
}

// KanbanBoardResponse groups applicants by stage for the Kanban view.
type KanbanBoardResponse struct {
	Screening []ApplicantListResponse `json:"SCREENING"`
	Interview []ApplicantListResponse `json:"INTERVIEW"`
	Offering  []ApplicantListResponse `json:"OFFERING"`
	Hired     []ApplicantListResponse `json:"HIRED"`
	Rejected  []ApplicantListResponse `json:"REJECTED"`
}
