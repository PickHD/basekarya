package reimburstment

import "gorm.io/gorm"

type Repository interface {
	Create(reimburstment *Reimbursement) error
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

func (r *repository) Create(reimburstment *Reimbursement) error {
	return r.db.Create(reimburstment).Error
}

func (r *repository) FindByID(id uint) (*Reimbursement, error) {
	var reimburstment Reimbursement

	err := r.db.Preload("User").First(&reimburstment).Error
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

	err := query.Preload("User").
		Limit(filter.Limit).
		Offset(filter.Offset).
		Order("created_at DESC").
		Find(&reimbursements).Error

	if err != nil {
		return nil, 0, err
	}

	return reimbursements, total, nil
}

func (r *repository) Update(reimbursement *Reimbursement) error {
	return r.db.Save(reimbursement).Error
}
