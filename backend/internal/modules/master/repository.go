package master

import (
	"context"

	"basekarya-backend/pkg/utils"

	"gorm.io/gorm"
)

type Repository interface {
	FindAllShifts(ctx context.Context) ([]Shift, error)
	FindAllLeaveTypes(ctx context.Context) ([]LeaveType, error)
	FindShiftByName(ctx context.Context, name string) (*Shift, error)
	SeedDefaults(ctx context.Context, companyID uint) error
}
type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) FindAllShifts(ctx context.Context) ([]Shift, error) {
	var shifts []Shift
	db := utils.TenantScope(ctx, utils.GetDBFromContext(ctx, r.db))
	if err := db.Model(&Shift{}).Find(&shifts).Error; err != nil {
		return nil, err
	}

	return shifts, nil
}

func (r *repository) FindAllLeaveTypes(ctx context.Context) ([]LeaveType, error) {
	var leaveTypes []LeaveType
	db := utils.TenantScope(ctx, utils.GetDBFromContext(ctx, r.db))
	if err := db.Model(&LeaveType{}).Find(&leaveTypes).Error; err != nil {
		return nil, err
	}

	return leaveTypes, nil
}

func (r *repository) FindShiftByName(ctx context.Context, name string) (*Shift, error) {
	var shift Shift
	db := utils.TenantScope(ctx, utils.GetDBFromContext(ctx, r.db))
	if err := db.Model(&Shift{}).Where("name = ?", name).First(&shift).Error; err != nil {
		return nil, err
	}

	return &shift, nil
}

func (r *repository) SeedDefaults(ctx context.Context, companyID uint) error {
	db := utils.GetDBFromContext(ctx, r.db)

	regularShift := Shift{Name: "Regular", StartTime: "09:00:00", EndTime: "18:00:00", CompanyID: companyID}
	if err := db.Where(Shift{Name: "Regular", CompanyID: companyID}).FirstOrCreate(&regularShift).Error; err != nil {
		return err
	}

	leaveTypes := []LeaveType{
		{Name: "Annual", DefaultQuota: 12, IsDeducted: true, CompanyID: companyID},
		{Name: "Sick", DefaultQuota: 15, IsDeducted: false, CompanyID: companyID},
		{Name: "Unpaid", DefaultQuota: 0, IsDeducted: false, CompanyID: companyID},
	}
	for _, lt := range leaveTypes {
		if err := db.Where(LeaveType{Name: lt.Name, CompanyID: companyID}).FirstOrCreate(&lt).Error; err != nil {
			return err
		}
	}
	return nil
}


