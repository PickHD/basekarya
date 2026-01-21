package reimbursement

import "gorm.io/gorm"

type Repository interface {
	Create(reimbursement *Reimbursement) error
	FindByID(id uint) (*Reimbursement, error)
	FindAll(filter ReimbursementFilter) ([]Reimbursement, int64, error)
	Update(reimbursement *Reimbursement) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) Create(reimbursement *Reimbursement) error {
	return r.db.Create(reimbursement).Error
}

func (r *repository) FindByID(id uint) (*Reimbursement, error) {
	var reimburstment Reimbursement

	err := r.db.Preload("User").First(&reimburstment, id).Error
	if err != nil {
		return nil, err
	}

	return &reimburstment, nil
}

func (r *repository) FindAll(filter ReimbursementFilter) ([]Reimbursement, int64, error) {
	var reimbursements []Reimbursement
	var total int64

	query := r.db.Model(&Reimbursement{})

	if filter.UserID > 0 {
		query = query.Where("user_id = ?", filter.UserID)
	}

	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (filter.Page - 1) * filter.Limit
	err := query.
		Limit(filter.Limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&reimbursements).Error

	return reimbursements, total, err
}

func (r *repository) Update(reimbursement *Reimbursement) error {
	return r.db.Save(reimbursement).Error
}
