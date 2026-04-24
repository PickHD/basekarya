package onboarding

import "time"

// ── Template DTOs ─────────────────────────────────────────────────────────────

type TemplateItemRequest struct {
	TaskName    string `json:"task_name" validate:"required"`
	Description string `json:"description"`
	SortOrder   int    `json:"sort_order"`
}

type CreateTemplateRequest struct {
	Name       string                `json:"name" validate:"required"`
	Department string                `json:"department" validate:"required"`
	Items      []TemplateItemRequest `json:"items"`
}

type UpdateTemplateRequest struct {
	Name       string                `json:"name" validate:"required"`
	Department string                `json:"department" validate:"required"`
	Items      []TemplateItemRequest `json:"items"`
}

type TemplateItemResponse struct {
	ID          uint   `json:"id"`
	TaskName    string `json:"task_name"`
	Description string `json:"description"`
	SortOrder   int    `json:"sort_order"`
}

type TemplateResponse struct {
	ID         uint                   `json:"id"`
	Name       string                 `json:"name"`
	Department string                 `json:"department"`
	Items      []TemplateItemResponse `json:"items"`
	CreatedAt  time.Time              `json:"created_at"`
}

// ── Workflow DTOs ─────────────────────────────────────────────────────────────

type CreateWorkflowRequest struct {
	ApplicantID  *uint  `json:"applicant_id"`
	EmployeeID   *uint  `json:"employee_id"`
	NewHireName  string `json:"new_hire_name" validate:"required"`
	NewHireEmail string `json:"new_hire_email" validate:"required,email"`
	Position     string `json:"position"`
	Department   string `json:"department"`
	StartDate    string `json:"start_date"` // YYYY-MM-DD
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
	Department  string     `json:"department"`
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
	ITTasks          []TaskResponse `json:"it_tasks"`
	HRTasks          []TaskResponse `json:"hr_tasks"`
	OtherTasks       []TaskResponse `json:"other_tasks"`
}

// ── Task Complete ─────────────────────────────────────────────────────────────

type CompleteTaskRequest struct {
	Notes string `json:"notes"`
}
