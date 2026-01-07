package health

import (
	"gorm.io/gorm"
)

type Repository interface {
	Check() error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) Check() error {
	sqlDB, err := r.db.DB()
	if err != nil {
		return err
	}

	err = sqlDB.Ping()
	if err != nil {
		return err
	}

	return nil
}
