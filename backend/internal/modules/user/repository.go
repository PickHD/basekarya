package user

import (
	"hris-backend/pkg/logger"

	"gorm.io/gorm"
)

type Repository interface {
	FindByUsername(username string) (*User, error)
	FindByID(id uint) (*User, error)
	UpdateEmployee(emp *Employee) error
	UpdateUser(user *User) error
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
