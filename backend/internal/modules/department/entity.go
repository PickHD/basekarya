package department

import "time"

type Department struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"not null" json:"name"`
	CompanyID uint      `gorm:"index;not null" json:"company_id"`
	CreatedAt time.Time `json:"created_at"`
}

func (Department) TableName() string {
	return "ref_departments"
}
