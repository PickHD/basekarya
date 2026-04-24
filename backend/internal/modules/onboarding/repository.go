package onboarding

import (
	"context"
	"time"

	"basekarya-backend/pkg/utils"

	"gorm.io/gorm"
)

type Repository interface {
	// Templates
	CreateTemplate(ctx context.Context, t *OnboardingTemplate) error
	FindAllTemplates(ctx context.Context) ([]OnboardingTemplate, error)
	FindTemplateByID(ctx context.Context, id uint) (*OnboardingTemplate, error)
	UpdateTemplate(ctx context.Context, t *OnboardingTemplate) error
	DeleteTemplate(ctx context.Context, id uint) error

	// Workflows
	CreateWorkflow(ctx context.Context, w *OnboardingWorkflow) error
	CreateTasks(ctx context.Context, tasks []OnboardingTask) error
	FindAllWorkflows(ctx context.Context, filter *WorkflowFilter) ([]OnboardingWorkflow, int64, error)
	FindWorkflowByID(ctx context.Context, id uint) (*OnboardingWorkflow, error)
	MarkWorkflowEmailSent(ctx context.Context, id uint) error

	// Tasks
	FindTaskByID(ctx context.Context, id uint) (*OnboardingTask, error)
	CompleteTask(ctx context.Context, id uint, completedBy uint, notes string) error
	CountPendingTasks(ctx context.Context, workflowID uint) (int64, error)
	CountTotalTasks(ctx context.Context, workflowID uint) (int64, error)
	MarkWorkflowCompleted(ctx context.Context, id uint) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) getDB(ctx context.Context) *gorm.DB {
	return utils.GetDBFromContext(ctx, r.db)
}

// ── Templates ─────────────────────────────────────────────────────────────────

func (r *repository) CreateTemplate(ctx context.Context, t *OnboardingTemplate) error {
	return r.getDB(ctx).Create(t).Error
}

func (r *repository) FindAllTemplates(ctx context.Context) ([]OnboardingTemplate, error) {
	var templates []OnboardingTemplate
	err := r.getDB(ctx).Preload("Items", func(db *gorm.DB) *gorm.DB {
		return db.Order("sort_order ASC")
	}).Order("department ASC, name ASC").Find(&templates).Error
	return templates, err
}

func (r *repository) FindTemplateByID(ctx context.Context, id uint) (*OnboardingTemplate, error) {
	var t OnboardingTemplate
	err := r.getDB(ctx).Preload("Items", func(db *gorm.DB) *gorm.DB {
		return db.Order("sort_order ASC")
	}).First(&t, id).Error
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *repository) UpdateTemplate(ctx context.Context, t *OnboardingTemplate) error {
	db := r.getDB(ctx)
	// Replace items: delete old, create new
	if err := db.Where("template_id = ?", t.ID).Delete(&OnboardingTemplateItem{}).Error; err != nil {
		return err
	}
	return db.Session(&gorm.Session{FullSaveAssociations: true}).Save(t).Error
}

func (r *repository) DeleteTemplate(ctx context.Context, id uint) error {
	return r.getDB(ctx).Delete(&OnboardingTemplate{}, id).Error
}

// ── Workflows ─────────────────────────────────────────────────────────────────

func (r *repository) CreateWorkflow(ctx context.Context, w *OnboardingWorkflow) error {
	return r.getDB(ctx).Create(w).Error
}

func (r *repository) CreateTasks(ctx context.Context, tasks []OnboardingTask) error {
	if len(tasks) == 0 {
		return nil
	}
	return r.getDB(ctx).Create(&tasks).Error
}

func (r *repository) FindAllWorkflows(ctx context.Context, filter *WorkflowFilter) ([]OnboardingWorkflow, int64, error) {
	var workflows []OnboardingWorkflow
	var total int64

	q := r.getDB(ctx).Model(&OnboardingWorkflow{})

	if filter.Status != "" {
		q = q.Where("status = ?", filter.Status)
	}
	if filter.Search != "" {
		like := "%" + filter.Search + "%"
		q = q.Where("new_hire_name LIKE ? OR new_hire_email LIKE ?", like, like)
	}

	q.Count(&total)

	offset := (filter.Page - 1) * filter.Limit
	err := q.Preload("Tasks").Order("created_at DESC").Offset(offset).Limit(filter.Limit).Find(&workflows).Error

	return workflows, total, err
}

func (r *repository) FindWorkflowByID(ctx context.Context, id uint) (*OnboardingWorkflow, error) {
	var w OnboardingWorkflow
	err := r.getDB(ctx).Preload("Tasks", func(db *gorm.DB) *gorm.DB {
		return db.Order("department ASC, sort_order ASC")
	}).Preload("Tasks.CompletedByUser.Employee").First(&w, id).Error
	if err != nil {
		return nil, err
	}
	return &w, nil
}

func (r *repository) MarkWorkflowEmailSent(ctx context.Context, id uint) error {
	return r.getDB(ctx).Model(&OnboardingWorkflow{}).Where("id = ?", id).
		Update("welcome_email_sent", true).Error
}

// ── Tasks ─────────────────────────────────────────────────────────────────────

func (r *repository) FindTaskByID(ctx context.Context, id uint) (*OnboardingTask, error) {
	var t OnboardingTask
	err := r.getDB(ctx).First(&t, id).Error
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *repository) CompleteTask(ctx context.Context, id uint, completedBy uint, notes string) error {
	now := time.Now()
	return r.getDB(ctx).Model(&OnboardingTask{}).Where("id = ?", id).Updates(map[string]any{
		"is_completed": true,
		"completed_by": completedBy,
		"completed_at": now,
		"notes":        notes,
	}).Error
}

func (r *repository) CountPendingTasks(ctx context.Context, workflowID uint) (int64, error) {
	var count int64
	err := r.getDB(ctx).Model(&OnboardingTask{}).
		Where("onboarding_workflow_id = ? AND is_completed = ?", workflowID, false).
		Count(&count).Error
	return count, err
}

func (r *repository) CountTotalTasks(ctx context.Context, workflowID uint) (int64, error) {
	var count int64
	err := r.getDB(ctx).Model(&OnboardingTask{}).
		Where("onboarding_workflow_id = ?", workflowID).
		Count(&count).Error
	return count, err
}

func (r *repository) MarkWorkflowCompleted(ctx context.Context, id uint) error {
	return r.getDB(ctx).Model(&OnboardingWorkflow{}).Where("id = ?", id).
		Update("status", WorkflowStatusCompleted).Error
}
