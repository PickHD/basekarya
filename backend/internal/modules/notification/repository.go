package notification

import (
	"context"
	"hris-backend/pkg/utils"
	"time"

	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, notification *Notification) error
	FindByID(id uint) (*Notification, error)
	FindAllByUserID(userID uint) ([]Notification, error)
	MarkAsRead(id uint) error
	DeleteReadOlderThan(days int) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) Create(ctx context.Context, notification *Notification) error {
	db := utils.GetDBFromContext(ctx, r.db)
	return db.Create(notification).Error
}

func (r *repository) FindByID(id uint) (*Notification, error) {
	var notification Notification
	err := r.db.
		First(&notification, id).Error

	return &notification, err
}

func (r *repository) FindAllByUserID(userID uint) ([]Notification, error) {
	var logs []Notification

	query := r.db.Model(&Notification{}).
		Select("notifications.*").
		Where("notifications.user_id = ?", userID).
		Order("notifications.created_at DESC")

	if err := query.Find(&logs).Error; err != nil {
		return nil, err
	}

	return logs, nil
}

func (r *repository) MarkAsRead(id uint) error {
	err := r.db.Model(&Notification{}).Where("id = ?", id).Update("is_read", true).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) DeleteReadOlderThan(days int) error {
	cutoffDate := time.Now().AddDate(0, 0, -days)
	err := r.db.Unscoped().
		Where("is_read = ? AND created_at < ?", true, cutoffDate).
		Delete(&Notification{}).Error

	if err != nil {
		return err
	}

	return nil
}
