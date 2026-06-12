package asset

import (
	"basekarya-backend/internal/modules/user"
	"basekarya-backend/pkg/constants"
	"database/sql"
	"time"
)

type AssetCategory struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	CompanyID   uint   `gorm:"index;not null" json:"company_id"`
	Name        string `gorm:"not null" json:"name"`
	Description string `gorm:"type:text" json:"description"`
}

func (AssetCategory) TableName() string {
	return "asset_categories"
}

type Asset struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	CompanyID       uint                 `gorm:"index;not null" json:"company_id"`
	AssetCategoryID uint                 `gorm:"not null" json:"asset_category_id"`
	AssetCategory   AssetCategory        `gorm:"foreignKey:AssetCategoryID" json:"asset_category,omitempty"`
	Name            string               `gorm:"not null" json:"name"`
	Description     string               `gorm:"type:text" json:"description"`
	SerialNumber    string               `gorm:"type:varchar(100)" json:"serial_number"`
	Status          constants.AssetStatus `gorm:"type:enum('AVAILABLE','ASSIGNED','MAINTENANCE','RETIRED');default:'AVAILABLE'" json:"status"`
	Condition       constants.AssetCondition `gorm:"type:enum('GOOD','FAIR','DAMAGED','LOST');default:'GOOD';column:condition" json:"condition"`
}

func (Asset) TableName() string {
	return "assets"
}

type AssetAssignment struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	CompanyID uint `gorm:"index;not null" json:"company_id"`

	AssetID   uint  `gorm:"not null" json:"asset_id"`
	Asset     Asset `gorm:"foreignKey:AssetID" json:"asset,omitempty"`

	EmployeeID uint          `gorm:"not null" json:"employee_id"`
	Employee   user.Employee `gorm:"foreignKey:EmployeeID" json:"employee,omitempty"`

	UserID uint      `gorm:"not null" json:"user_id"`
	User   user.User `gorm:"foreignKey:UserID" json:"user,omitempty"`

	ApprovedBy *uint      `json:"approved_by"`
	Approver   *user.User `gorm:"foreignKey:ApprovedBy" json:"approver,omitempty"`

	Purpose            string  `gorm:"type:text" json:"purpose"`
	ExpectedReturnDate *string `gorm:"type:date" json:"expected_return_date"`
	ActualReturnDate   *string `gorm:"type:date" json:"actual_return_date"`
	Notes              string  `gorm:"type:text" json:"notes"`

	Status          constants.AssetAssignmentStatus `gorm:"type:enum('PENDING','ACTIVE','RETURNED','REJECTED');default:'PENDING'" json:"status"`
	RejectionReason sql.NullString                  `gorm:"type:text" json:"rejection_reason"`
}

func (AssetAssignment) TableName() string {
	return "asset_assignments"
}
