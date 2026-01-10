package user

import (
	"hris-backend/internal/modules/master"
	"time"
)

type User struct {
	ID                 uint      `gorm:"primaryKey" json:"id"`
	Username           string    `gorm:"unique;not null" json:"username"`
	PasswordHash       string    `json:"-"`
	Role               string    `gorm:"type:enum('SUPERADMIN','EMPLOYEE');default:'EMPLOYEE'" json:"role"`
	MustChangePassword bool      `json:"must_change_password"`
	IsActive           bool      `gorm:"default:true" json:"is_active"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`

	Employee *Employee `gorm:"foreignKey:UserID;references:ID" json:"employee,omitempty"`
}

type Employee struct {
	ID                uint   `gorm:"primaryKey" json:"id"`
	UserID            uint   `gorm:"unique;not null" json:"user_id"`
	DepartmentID      uint   `json:"department_id"`
	ShiftID           uint   `json:"shift_id"`
	NIK               string `gorm:"unique;not null" json:"nik"`
	FullName          string `json:"full_name"`
	PhoneNumber       string `json:"phone_number"`
	ProfilePictureUrl string `json:"profile_picture_url"`

	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`

	Department *master.Department `gorm:"foreignKey:DepartmentID" json:"department,omitempty"`
	Shift      *master.Shift      `gorm:"foreignKey:ShiftID" json:"shift,omitempty"`
}
