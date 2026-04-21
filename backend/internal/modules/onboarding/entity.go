package onboarding

import (
	"time"

	"basekarya-backend/internal/modules/user"

	"gorm.io/gorm"
)

// OnboardingTemplate is a reusable checklist template (e.g. "IT Setup", "HR Document Collection").
type OnboardingTemplate struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Name       string `gorm:"type:varchar(100);not null" json:"name"`
	Department string `gorm:"type:varchar(50);not null" json:"department"`

	Items []OnboardingTemplateItem `gorm:"foreignKey:TemplateID" json:"items,omitempty"`
}

func (OnboardingTemplate) TableName() string { return "onboarding_templates" }

// OnboardingTemplateItem is a single task inside a template.
type OnboardingTemplateItem struct {
	ID         uint   `gorm:"primaryKey" json:"id"`
	TemplateID uint   `gorm:"not null" json:"template_id"`
	TaskName   string `gorm:"type:varchar(255);not null" json:"task_name"`
	Description string `gorm:"type:text" json:"description"`
	SortOrder  int    `gorm:"default:0" json:"sort_order"`
}

func (OnboardingTemplateItem) TableName() string { return "onboarding_template_items" }

// OnboardingWorkflow is a per-hire instance of the onboarding process.
type OnboardingWorkflow struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	ApplicantID *uint  `json:"applicant_id"`
	EmployeeID  *uint  `json:"employee_id"`
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
	TemplateItemID       *uint          `json:"template_item_id"`

	TaskName    string     `gorm:"type:varchar(255);not null" json:"task_name"`
	Description string     `gorm:"type:text" json:"description"`
	Department  string     `gorm:"type:varchar(50);not null" json:"department"`
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

// DeletedAt is not on workflows/tasks; we use hard deletes only via cascade.
// Soft-delete is only on the template level.
type OnboardingTemplateWithDeleted struct {
	OnboardingTemplate
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
