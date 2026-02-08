package master

import "time"

type Department struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"not null" json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type Shift struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"not null" json:"name"`
	StartTime string    `gorm:"not null" json:"start_time"`
	EndTime   string    `gorm:"not null" json:"end_time"`
	CreatedAt time.Time `json:"created_at"`
}

type LeaveType struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Name         string    `gorm:"unique;not null" json:"name"`
	DefaultQuota int       `json:"default_quota"`
	IsDeducted   bool      `gorm:"default:true" json:"is_deducted"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (Department) TableName() string {
	return "ref_departments"
}

func (Shift) TableName() string {
	return "ref_shifts"
}

func (LeaveType) TableName() string {
	return "ref_leave_types"
}
