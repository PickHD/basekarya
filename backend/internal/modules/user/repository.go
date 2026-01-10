package user

import (
	"hris-backend/pkg/constants"
	"hris-backend/pkg/logger"

	"gorm.io/gorm"
)

type Repository interface {
	FindByUsername(username string) (*User, error)
	FindByID(id uint) (*User, error)
	UpdateEmployee(emp *Employee) error
	UpdateUser(user *User) error
	FindAllEmployees(page, limit int, search string) ([]User, int64, error)
	CreateUser(tx *gorm.DB, user *User) error
	CreateEmployee(tx *gorm.DB, emp *Employee) error
	DeleteUser(id uint) error
	FindEmployeeByID(id uint) (*Employee, error)
	StartTX() *gorm.DB
	CountActiveEmployee() (int64, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) FindByUsername(username string) (*User, error) {
	var user User

	err := r.db.Preload("Employee").Where("username = ?", username).First(&user).Error
	if err != nil {
		logger.Errorw("UserRepository.FindByUsername ERROR: ", err)

		return nil, err
	}

	return &user, nil
}

func (r *repository) FindByID(id uint) (*User, error) {
	var user User

	err := r.db.Preload("Employee.Department").Preload("Employee.Shift").First(&user, id).Error
	if err != nil {
		logger.Errorw("UserRepository.FindByID ERROR: ", err)

		return nil, err
	}

	return &user, nil
}

func (r *repository) UpdateEmployee(emp *Employee) error {
	return r.db.Save(emp).Error
}

func (r *repository) UpdateUser(user *User) error {
	return r.db.Save(user).Error
}

func (r *repository) FindAllEmployees(page, limit int, search string) ([]User, int64, error) {
	var users []User
	var total int64

	query := r.db.Model(&User{}).
		Joins("JOIN employees ON employees.user_id = users.id").
		Preload("Employee").
		Preload("Employee.Department").
		Preload("Employee.Shift")

	// filter search by fullname or NIK/ID
	if search != "" {
		searchParam := "%" + search + "%"
		query = query.Where("LOWER(employees.full_name) LIKE LOWER(?) OR LOWER(employees.nik) LIKE LOWER(?)", searchParam, searchParam)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	err := query.Limit(limit).Offset(offset).Order("employees.full_name ASC").Find(&users).Error

	return users, total, err
}

func (r *repository) CreateUser(tx *gorm.DB, user *User) error {
	return tx.Create(user).Error
}

func (r *repository) CreateEmployee(tx *gorm.DB, emp *Employee) error {
	return tx.Create(emp).Error
}

func (r *repository) DeleteUser(id uint) error {
	return r.db.Delete(&User{}, id).Error
}

func (r *repository) FindEmployeeByID(id uint) (*Employee, error) {
	var emp Employee
	err := r.db.Preload("User").First(&emp, id).Error
	return &emp, err
}

func (r *repository) StartTX() *gorm.DB {
	return r.db.Begin()
}

func (r *repository) CountActiveEmployee() (int64, error) {
	var totalActive int64
	if err := r.db.Model(&User{}).
		Where("is_active = ? AND role = ?", true, string(constants.UserRoleEmployee)).
		Count(&totalActive).Error; err != nil {
		return 0, err
	}

	return totalActive, nil
}
