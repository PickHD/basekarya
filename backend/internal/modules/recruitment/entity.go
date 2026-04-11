package recruitment

import (
	"time"

	"basekarya-backend/internal/modules/master"
	"basekarya-backend/internal/modules/user"

	"gorm.io/gorm"
)

// JobRequisition represents a job opening request.
type JobRequisition struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	RequesterID    uint   `gorm:"not null" json:"requester_id"`
	DepartmentID   uint   `gorm:"not null" json:"department_id"`
	Title          string `gorm:"type:varchar(255);not null" json:"title"`
	Description    string `gorm:"type:text" json:"description"`
	Quantity       int    `gorm:"default:1" json:"quantity"`
	EmploymentType string `gorm:"type:varchar(10);not null" json:"employment_type"`
	Priority       string `gorm:"type:varchar(10);default:'MEDIUM'" json:"priority"`
	Status         string `gorm:"type:varchar(10);default:'DRAFT'" json:"status"`

	ApprovedBy      *uint  `json:"approved_by"`
	RejectionReason string `gorm:"type:text" json:"rejection_reason"`
	TargetDate      *time.Time `gorm:"type:date" json:"target_date"`

	// Relations
	Requester  *user.User        `gorm:"foreignKey:RequesterID" json:"requester,omitempty"`
	Approver   *user.User        `gorm:"foreignKey:ApprovedBy" json:"approver,omitempty"`
	Department *master.Department `gorm:"foreignKey:DepartmentID" json:"department,omitempty"`
	Applicants []Applicant       `gorm:"foreignKey:JobRequisitionID" json:"applicants,omitempty"`
}

func (JobRequisition) TableName() string {
	return "job_requisitions"
}

// Applicant represents a candidate for a job requisition.
type Applicant struct {
	ID               uint           `gorm:"primaryKey" json:"id"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"`

	JobRequisitionID uint   `gorm:"not null" json:"job_requisition_id"`
	FullName         string `gorm:"type:varchar(100);not null" json:"full_name"`
	Email            string `gorm:"type:varchar(255);not null" json:"email"`
	PhoneNumber      string `gorm:"type:varchar(20)" json:"phone_number"`
	ResumeURL        string `gorm:"type:varchar(255)" json:"resume_url"`

	Stage      string `gorm:"type:varchar(15);default:'SCREENING'" json:"stage"`
	StageOrder int    `gorm:"default:0" json:"stage_order"`

	Notes           string `gorm:"type:text" json:"notes"`
	RejectionReason string `gorm:"type:text" json:"rejection_reason"`

	// Relations
	JobRequisition *JobRequisition      `gorm:"foreignKey:JobRequisitionID" json:"job_requisition,omitempty"`
	StageHistories []ApplicantStageHistory `gorm:"foreignKey:ApplicantID" json:"stage_histories,omitempty"`
}

func (Applicant) TableName() string {
	return "applicants"
}

// ApplicantStageHistory records stage transitions for audit trail.
type ApplicantStageHistory struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	CreatedAt   time.Time `json:"created_at"`

	ApplicantID uint   `gorm:"not null" json:"applicant_id"`
	FromStage   string `gorm:"type:varchar(15)" json:"from_stage"`
	ToStage     string `gorm:"type:varchar(15);not null" json:"to_stage"`
	ChangedBy   uint   `gorm:"not null" json:"changed_by"`
	Notes       string `gorm:"type:text" json:"notes"`

	// Relations
	ChangedByUser *user.User `gorm:"foreignKey:ChangedBy" json:"changed_by_user,omitempty"`
}

func (ApplicantStageHistory) TableName() string {
	return "applicant_stage_histories"
}
