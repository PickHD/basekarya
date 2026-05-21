package notification

import (
	"context"
	"basekarya-backend/pkg/utils"
	"time"

	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, notification *Notification) error
	CreateBatch(ctx context.Context, notifications []*Notification) error
	FindByID(ctx context.Context, id uint) (*Notification, error)
	FindAllByUserID(ctx context.Context, userID uint) ([]Notification, error)
	MarkAsRead(ctx context.Context, id uint) error
	DeleteReadOlderThan(ctx context.Context, days int) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) Create(ctx context.Context, notification *Notification) error {
	db := utils.TenantScope(ctx, utils.GetDBFromContext(ctx, r.db))
	return db.Create(notification).Error
}

func (r *repository) CreateBatch(ctx context.Context, notifications []*Notification) error {
	db := utils.TenantScope(ctx, utils.GetDBFromContext(ctx, r.db))
	return db.Create(&notifications).Error
}

func (r *repository) FindByID(ctx context.Context, id uint) (*Notification, error) {
	db := utils.TenantScope(ctx, r.db)
	var notification Notification
	err := db.
		First(&notification, id).Error

	return &notification, err
}

func (r *repository) FindAllByUserID(ctx context.Context, userID uint) ([]Notification, error) {
	db := utils.TenantScope(ctx, r.db)
	var logs []Notification

	query := db.Model(&Notification{}).
		Select("notifications.*").
		Where("notifications.user_id = ?", userID).
		Order("notifications.created_at DESC")

	if err := query.Find(&logs).Error; err != nil {
		return nil, err
	}

	return logs, nil
}

func (r *repository) MarkAsRead(ctx context.Context, id uint) error {
	db := utils.TenantScope(ctx, r.db)
	err := db.Model(&Notification{}).Where("id = ?", id).Update("is_read", true).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) DeleteReadOlderThan(ctx context.Context, days int) error {
	db := utils.TenantScope(ctx, r.db)
	cutoffDate := time.Now().AddDate(0, 0, -days)
	err := db.Unscoped().
		Where("is_read = ? AND created_at < ?", true, cutoffDate).
		Delete(&Notification{}).Error

	if err != nil {
		return err
	}

	return nil
}
