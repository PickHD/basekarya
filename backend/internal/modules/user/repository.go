package user

import (
	"context"
	"hris-backend/pkg/constants"
	"hris-backend/pkg/logger"
	"hris-backend/pkg/utils"

	"gorm.io/gorm"
)

type Repository interface {
	FindByUsername(ctx context.Context, username string) (*User, error)
	FindByID(ctx context.Context, id uint) (*User, error)
	UpdateEmployee(ctx context.Context, emp *Employee) error
	UpdateUser(ctx context.Context, user *User) error
	FindAllEmployees(ctx context.Context, page, limit int, search string) ([]User, int64, error)
	CreateUser(ctx context.Context, user *User) error
	CreateEmployee(ctx context.Context, emp *Employee) error
	DeleteUser(ctx context.Context, id uint) error
	FindEmployeeByID(ctx context.Context, id uint) (*Employee, error)
	CountActiveEmployee(ctx context.Context) (int64, error)
	FindAllEmployeeActive(ctx context.Context) ([]Employee, error)
	FindAdminID(ctx context.Context) (uint, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) FindByUsername(ctx context.Context, username string) (*User, error) {
	db := utils.GetDBFromContext(ctx, r.db)
	var user User

	err := db.Preload("Employee").Where("username = ?", username).First(&user).Error
	if err != nil {
		logger.Errorw("UserRepository.FindByUsername ERROR: ", err)

		return nil, err
	}

	return &user, nil
}

func (r *repository) FindByID(ctx context.Context, id uint) (*User, error) {
	db := utils.GetDBFromContext(ctx, r.db)
	var user User

	err := db.Preload("Employee.Department").Preload("Employee.Shift").First(&user, id).Error
	if err != nil {
		logger.Errorw("UserRepository.FindByID ERROR: ", err)

		return nil, err
	}

	return &user, nil
}

func (r *repository) UpdateEmployee(ctx context.Context, emp *Employee) error {
	db := utils.GetDBFromContext(ctx, r.db)
	return db.Save(emp).Error
}

func (r *repository) UpdateUser(ctx context.Context, user *User) error {
	db := utils.GetDBFromContext(ctx, r.db)
	return db.Save(user).Error
}

func (r *repository) FindAllEmployees(ctx context.Context, page, limit int, search string) ([]User, int64, error) {
	db := utils.GetDBFromContext(ctx, r.db)
	var users []User
	var total int64

	query := db.Model(&User{}).
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

func (r *repository) CreateUser(ctx context.Context, user *User) error {
	db := utils.GetDBFromContext(ctx, r.db)
	return db.Create(user).Error
}

func (r *repository) CreateEmployee(ctx context.Context, emp *Employee) error {
	db := utils.GetDBFromContext(ctx, r.db)
	return db.Create(emp).Error
}

func (r *repository) DeleteUser(ctx context.Context, id uint) error {
	db := utils.GetDBFromContext(ctx, r.db)
	return db.Delete(&User{}, id).Error
}

func (r *repository) FindEmployeeByID(ctx context.Context, id uint) (*Employee, error) {
	db := utils.GetDBFromContext(ctx, r.db)
	var emp Employee
	err := db.Preload("User").First(&emp, id).Error
	return &emp, err
}

func (r *repository) CountActiveEmployee(ctx context.Context) (int64, error) {
	db := utils.GetDBFromContext(ctx, r.db)
	var totalActive int64
	if err := db.Model(&User{}).
		Where("is_active = ? AND role = ?", true, string(constants.UserRoleEmployee)).
		Count(&totalActive).Error; err != nil {
		return 0, err
	}

	return totalActive, nil
}

func (r *repository) FindAllEmployeeActive(ctx context.Context) ([]Employee, error) {
	db := utils.GetDBFromContext(ctx, r.db)
	var employees []Employee

	if err := db.Model(&Employee{}).
		Joins("User").
		Where("User.is_active = ? AND User.role = ?", true, string(constants.UserRoleEmployee)).
		Preload("User").
		Preload("Department").
		Preload("Shift").
		Find(&employees).Error; err != nil {
		return nil, err
	}

	return employees, nil
}

func (r *repository) FindAdminID(ctx context.Context) (uint, error) {
	db := utils.GetDBFromContext(ctx, r.db)
	var id uint
	err := db.Model(&User{}).
		Select("id").
		Where("role = ?", string(constants.UserRoleSuperadmin)).
		Scan(&id).Error

	if err != nil {
		logger.Errorw("UserRepository.FindAdminID ERROR: ", err)

		return 0, err
	}

	return id, nil
}
