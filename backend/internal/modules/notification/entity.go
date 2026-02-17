package notification

import (
	"hris-backend/internal/modules/user"
	"time"
)

type Notification struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`

	UserID uint `json:"user_id"`

	Type      string `gorm:"type:varchar(50);not null" json:"type"`
	Title     string `gorm:"type:varchar(255);not null" json:"title"`
	Message   string `gorm:"type:text" json:"message"`
	IsRead    bool   `gorm:"default:false;index" json:"is_read"`
	RelatedID uint   `json:"related_id"`

	User user.User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

func (Notification) TableName() string {
	return "notifications"
}
