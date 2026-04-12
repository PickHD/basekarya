package recruitment

import (
	"basekarya-backend/pkg/utils"
	"context"

	"gorm.io/gorm"
)

type Repository interface {
	CreateRequisition(ctx context.Context, req *JobRequisition) error
	FindRequisitionByID(ctx context.Context, id uint) (*JobRequisition, error)
	FindAllRequisitions(ctx context.Context, filter *RequisitionFilter) ([]JobRequisition, int64, error)
	UpdateRequisitionStatus(ctx context.Context, id uint, status string, approvedBy *uint, rejectionReason string) error
	SoftDeleteRequisition(ctx context.Context, id uint) error

	CreateApplicant(ctx context.Context, applicant *Applicant) error
	FindApplicantByID(ctx context.Context, id uint) (*Applicant, error)
	FindApplicantsByRequisitionID(ctx context.Context, requisitionID uint) ([]Applicant, error)
	UpdateApplicantStage(ctx context.Context, id uint, stage string, stageOrder int, notes, rejectionReason string) error
	CreateStageHistory(ctx context.Context, history *ApplicantStageHistory) error
	CountApplicantsByRequisitionAndStage(ctx context.Context, requisitionID uint, stage string) (int64, error)
	SoftDeleteApplicant(ctx context.Context, id uint) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) CreateRequisition(ctx context.Context, req *JobRequisition) error {
	db := utils.GetDBFromContext(ctx, r.db)
	return db.Create(req).Error
}

func (r *repository) FindRequisitionByID(ctx context.Context, id uint) (*JobRequisition, error) {
	db := utils.GetDBFromContext(ctx, r.db)
	var req JobRequisition
	err := db.
		Preload("Requester.Employee").
		Preload("Approver.Employee").
		Preload("Department").
		First(&req, id).Error
	if err != nil {
		return nil, err
	}
	return &req, nil
}

func (r *repository) FindAllRequisitions(ctx context.Context, filter *RequisitionFilter) ([]JobRequisition, int64, error) {
	var items []JobRequisition
	var total int64

	db := utils.GetDBFromContext(ctx, r.db)
	q := db.Model(&JobRequisition{})

	if filter.Status != "" {
		q = q.Where("status = ?", filter.Status)
	}
	if filter.DepartmentID > 0 {
		q = q.Where("department_id = ?", filter.DepartmentID)
	}
	if filter.Priority != "" {
		q = q.Where("priority = ?", filter.Priority)
	}
	if filter.Search != "" {
		q = q.Where("title LIKE ?", "%"+filter.Search+"%")
	}

	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (filter.Page - 1) * filter.Limit
	err := q.
		Preload("Requester.Employee").
		Preload("Department").
		Joins("LEFT JOIN users u ON u.id = job_requisitions.requester_id").
		Joins("LEFT JOIN ref_departments d ON d.id = job_requisitions.department_id").
		Order("created_at DESC").
		Limit(filter.Limit).
		Offset(offset).
		Find(&items).Error

	return items, total, err
}

func (r *repository) UpdateRequisitionStatus(ctx context.Context, id uint, status string, approvedBy *uint, rejectionReason string) error {
	db := utils.GetDBFromContext(ctx, r.db)
	updates := map[string]interface{}{
		"status":           status,
		"approved_by":      approvedBy,
		"rejection_reason": rejectionReason,
	}
	return db.Model(&JobRequisition{}).Where("id = ?", id).Updates(updates).Error
}

func (r *repository) SoftDeleteRequisition(ctx context.Context, id uint) error {
	db := utils.GetDBFromContext(ctx, r.db)
	return db.Delete(&JobRequisition{}, id).Error
}

func (r *repository) CreateApplicant(ctx context.Context, applicant *Applicant) error {
	db := utils.GetDBFromContext(ctx, r.db)
	return db.Create(applicant).Error
}

func (r *repository) FindApplicantByID(ctx context.Context, id uint) (*Applicant, error) {
	db := utils.GetDBFromContext(ctx, r.db)
	var applicant Applicant
	err := db.
		Preload("JobRequisition").
		Preload("StageHistories").
		Preload("StageHistories.ChangedByUser").
		First(&applicant, id).Error
	if err != nil {
		return nil, err
	}
	return &applicant, nil
}

func (r *repository) FindApplicantsByRequisitionID(ctx context.Context, requisitionID uint) ([]Applicant, error) {
	db := utils.GetDBFromContext(ctx, r.db)
	var applicants []Applicant
	err := db.
		Where("job_requisition_id = ?", requisitionID).
		Order("stage_order ASC").
		Find(&applicants).Error
	return applicants, err
}

func (r *repository) UpdateApplicantStage(ctx context.Context, id uint, stage string, stageOrder int, notes, rejectionReason string) error {
	db := utils.GetDBFromContext(ctx, r.db)
	updates := map[string]interface{}{
		"stage":            stage,
		"stage_order":      stageOrder,
		"notes":            notes,
		"rejection_reason": rejectionReason,
	}
	return db.Model(&Applicant{}).Where("id = ?", id).Updates(updates).Error
}

func (r *repository) CreateStageHistory(ctx context.Context, history *ApplicantStageHistory) error {
	db := utils.GetDBFromContext(ctx, r.db)
	return db.Create(history).Error
}

func (r *repository) CountApplicantsByRequisitionAndStage(ctx context.Context, requisitionID uint, stage string) (int64, error) {
	db := utils.GetDBFromContext(ctx, r.db)
	var count int64
	err := db.Model(&Applicant{}).
		Where("job_requisition_id = ? AND stage = ?", requisitionID, stage).
		Count(&count).Error
	return count, err
}

func (r *repository) SoftDeleteApplicant(ctx context.Context, id uint) error {
	db := utils.GetDBFromContext(ctx, r.db)
	return db.Delete(&Applicant{}, id).Error
}
