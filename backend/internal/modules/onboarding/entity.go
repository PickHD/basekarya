package onboarding

import (
	"time"

	"basekarya-backend/internal/modules/user"
)

// OnboardingWorkflow is a per-hire instance of the onboarding process.
type OnboardingWorkflow struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	ApplicantID *uint  `json:"applicant_id"`
	EmployeeID  *uint  `json:"employee_id"`
	CompanyID   uint   `gorm:"index;not null" json:"company_id"`
	NewHireName  string `gorm:"type:varchar(100);not null" json:"new_hire_name"`
	NewHireEmail string `gorm:"type:varchar(255);not null" json:"new_hire_email"`
	Position     string `gorm:"type:varchar(100)" json:"position"`
	Department   string `gorm:"type:varchar(100)" json:"department"`
	StartDate    *time.Time `gorm:"type:date" json:"start_date"`

	Status           string `gorm:"type:varchar(15);default:'IN_PROGRESS'" json:"status"`
	WelcomeEmailSent bool   `gorm:"default:false" json:"welcome_email_sent"`

	Tasks []OnboardingTask `gorm:"foreignKey:OnboardingWorkflowID" json:"tasks,omitempty"`
}

func (OnboardingWorkflow) TableName() string { return "onboarding_workflows" }

// OnboardingTask is an individual task assigned to a workflow.
type OnboardingTask struct {
	ID                   uint           `gorm:"primaryKey" json:"id"`
	OnboardingWorkflowID uint           `gorm:"not null" json:"onboarding_workflow_id"`
	CompanyID            uint           `gorm:"index;not null" json:"company_id"`

	TaskName    string     `gorm:"type:varchar(255);not null" json:"task_name"`
	Description string     `gorm:"type:text" json:"description"`
	IsCompleted bool       `gorm:"default:false" json:"is_completed"`
	CompletedBy *uint      `json:"completed_by"`
	CompletedAt *time.Time `json:"completed_at"`
	Notes       string     `gorm:"type:text" json:"notes"`
	SortOrder   int        `gorm:"default:0" json:"sort_order"`

	// Relations
	CompletedByUser *user.User `gorm:"foreignKey:CompletedBy;references:ID" json:"completed_by_user,omitempty"`
}

func (OnboardingTask) TableName() string { return "onboarding_tasks" }

// OnboardingWorkflowStatus constants.
const (
	WorkflowStatusInProgress = "IN_PROGRESS"
	WorkflowStatusCompleted  = "COMPLETED"
)

