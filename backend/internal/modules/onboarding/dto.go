package onboarding

import "time"

// ── Workflow DTOs ─────────────────────────────────────────────────────────────

type CreateWorkflowRequest struct {
	ApplicantID  *uint                 `json:"applicant_id"`
	EmployeeID   *uint                 `json:"employee_id"`
	NewHireName  string                `json:"new_hire_name" validate:"required"`
	NewHireEmail string                `json:"new_hire_email" validate:"required,email"`
	Position     string                `json:"position"`
	Department   string                `json:"department"`
	StartDate    string                `json:"start_date"` // YYYY-MM-DD
	Tasks        []WorkflowTaskRequest `json:"tasks"`
}

type WorkflowTaskRequest struct {
	TaskName    string `json:"task_name" validate:"required"`
	Description string `json:"description"`
	SortOrder   int    `json:"sort_order"`
}

type WorkflowFilter struct {
	Status string
	Search string
	Page   int
	Limit  int
}

type TaskResponse struct {
	ID          uint       `json:"id"`
	TaskName    string     `json:"task_name"`
	Description string     `json:"description"`
	IsCompleted bool       `json:"is_completed"`
	CompletedBy string     `json:"completed_by,omitempty"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
	Notes       string     `json:"notes"`
	SortOrder   int        `json:"sort_order"`
}

type WorkflowListResponse struct {
	ID           uint      `json:"id"`
	NewHireName  string    `json:"new_hire_name"`
	NewHireEmail string    `json:"new_hire_email"`
	Position     string    `json:"position"`
	Department   string    `json:"department"`
	StartDate    *time.Time `json:"start_date"`
	Status       string    `json:"status"`
	Progress     int       `json:"progress"` // percentage
	CreatedAt    time.Time `json:"created_at"`
}

type WorkflowDetailResponse struct {
	ID               uint           `json:"id"`
	NewHireName      string         `json:"new_hire_name"`
	NewHireEmail     string         `json:"new_hire_email"`
	Position         string         `json:"position"`
	Department       string         `json:"department"`
	StartDate        *time.Time     `json:"start_date"`
	Status           string         `json:"status"`
	Progress         int            `json:"progress"`
	WelcomeEmailSent bool           `json:"welcome_email_sent"`
	CreatedAt        time.Time      `json:"created_at"`
	Tasks            []TaskResponse `json:"tasks"`
}

// ── Task Update ───────────────────────────────────────────────────────────────

type UpdateWorkflowTasksRequest struct {
	Tasks []WorkflowTaskRequest `json:"tasks" validate:"required,min=1"`
}

// ── Task Complete ─────────────────────────────────────────────────────────────

type CompleteTaskRequest struct {
	Notes string `json:"notes"`
}
