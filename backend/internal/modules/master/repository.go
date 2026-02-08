package master

import "gorm.io/gorm"

type Repository interface {
	FindAllDepartments() ([]Department, error)
	FindAllShifts() ([]Shift, error)
	FindAllLeaveTypes() ([]LeaveType, error)
}
type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) FindAllDepartments() ([]Department, error) {
	var deps []Department
	if err := r.db.Model(&Department{}).Find(&deps).Error; err != nil {
		return nil, err
	}

	return deps, nil
}

func (r *repository) FindAllShifts() ([]Shift, error) {
	var shifts []Shift
	if err := r.db.Model(&Shift{}).Find(&shifts).Error; err != nil {
		return nil, err
	}

	return shifts, nil
}

func (r *repository) FindAllLeaveTypes() ([]LeaveType, error) {
	var leaveTypes []LeaveType
	if err := r.db.Model(&LeaveType{}).Find(&leaveTypes).Error; err != nil {
		return nil, err
	}

	return leaveTypes, nil
}
