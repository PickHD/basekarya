package attendance

import (
	"hris-backend/internal/modules/master"
	"hris-backend/internal/modules/user"
	"time"
)

type Attendance struct {
	ID uint `gorm:"primaryKey" json:"id"`

	EmployeeID uint `gorm:"not null;uniqueIndex:idx_emp_date,priority:1" json:"employee_id"`
	ShiftID    uint `gorm:"not null" json:"shift_id"`

	Date time.Time `gorm:"type:date;not null;uniqueIndex:idx_emp_date,priority:2;index:idx_date_status,priority:1" json:"date"`

	CheckInTime     time.Time `gorm:"not null" json:"check_in_time"`
	CheckInLat      float64   `gorm:"type:decimal(10,8);not null" json:"check_in_lat"`
	CheckInLong     float64   `gorm:"type:decimal(11,8);not null" json:"check_in_long"`
	CheckInImageURL string    `gorm:"size:255;not null" json:"check_in_image_url"`

	CheckInAddress string `gorm:"type:varchar(500);not null" json:"check_in_address"`

	CheckOutTime     *time.Time `json:"check_out_time"`
	CheckOutLat      *float64   `gorm:"type:decimal(10,8)" json:"check_out_lat"`
	CheckOutLong     *float64   `gorm:"type:decimal(11,8)" json:"check_out_long"`
	CheckOutImageURL *string    `gorm:"size:255" json:"check_out_image_url"`
	CheckOutAddress  *string    `gorm:"type:varchar(500)" json:"check_out_address"`

	Status string `gorm:"type:enum('PRESENT', 'LATE', 'EXCUSED', 'ABSENT');default:'ABSENT';index:idx_date_status,priority:2" json:"status"`

	IsSuspicious bool   `gorm:"default:false;index" json:"is_suspicious"`
	Notes        string `gorm:"type:varchar(500)" json:"notes"`

	LateDurationMinute int `gorm:"default:0" json:"late_duration_minute"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Employee *user.Employee `gorm:"foreignKey:EmployeeID" json:"employee,omitempty"`
	Shift    *master.Shift  `gorm:"foreignKey:ShiftID" json:"shift,omitempty"`
}

func (Attendance) TableName() string {
	return "attendances"
}
