package asset

import (
	"basekarya-backend/pkg/constants"
	"time"
)

type AssetCategoryFilter struct {
	Page  int
	Limit int
}

type AssetFilter struct {
	Status     string
	Condition  string
	CategoryID uint
	Page       int
	Limit      int
}

type AssetAssignmentFilter struct {
	UserID uint
	Status string
	Page   int
	Limit  int
}

type CreateAssetCategoryRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
}

type UpdateAssetCategoryRequest struct {
	ID          uint   `json:"-"`
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
}

type CreateAssetRequest struct {
	AssetCategoryID uint                   `json:"asset_category_id" validate:"required"`
	Name            string                 `json:"name" validate:"required"`
	Description     string                 `json:"description"`
	SerialNumber    string                 `json:"serial_number"`
	Condition       constants.AssetCondition `json:"condition"`
}

type UpdateAssetRequest struct {
	ID              uint                    `json:"-"`
	AssetCategoryID uint                    `json:"asset_category_id"`
	Name            string                  `json:"name"`
	Description     string                  `json:"description"`
	SerialNumber    string                  `json:"serial_number"`
	Status          constants.AssetStatus    `json:"status"`
	Condition       constants.AssetCondition `json:"condition"`
}

type CreateAssetAssignmentRequest struct {
	UserID             uint   `json:"-"`
	EmployeeID         uint   `json:"-"`
	AssetID            uint   `json:"asset_id" validate:"required"`
	Purpose            string `json:"purpose" validate:"required"`
	ExpectedReturnDate string `json:"expected_return_date"`
}

type ActionRequest struct {
	ID              uint   `json:"-"`
	SuperAdminID    uint   `json:"-"`
	Action          string `json:"action" validate:"required"`
	RejectionReason string `json:"rejection_reason" validate:"omitempty"`
}

type ReturnRequest struct {
	ID   uint `json:"-"`
	UserID uint `json:"-"`
}

type AssetCategoryResponse struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type AssetListResponse struct {
	ID               uint                     `json:"id"`
	Name             string                   `json:"name"`
	Description      string                   `json:"description"`
	SerialNumber     string                   `json:"serial_number"`
	AssetCategoryID  uint                     `json:"asset_category_id"`
	CategoryName     string                   `json:"category_name"`
	Status           constants.AssetStatus    `json:"status"`
	Condition        constants.AssetCondition `json:"condition"`
	CurrentEmployee  string                   `json:"current_employee"`
	CreatedAt        time.Time                `json:"created_at"`
	UpdatedAt        time.Time                `json:"updated_at"`
}

type AssetDetailResponse struct {
	ID              uint                     `json:"id"`
	Name            string                   `json:"name"`
	Description     string                   `json:"description"`
	SerialNumber    string                   `json:"serial_number"`
	AssetCategoryID uint                     `json:"asset_category_id"`
	CategoryName    string                   `json:"category_name"`
	Status          constants.AssetStatus    `json:"status"`
	Condition       constants.AssetCondition `json:"condition"`
	CurrentEmployee string                   `json:"current_employee"`
	CreatedAt       time.Time                `json:"created_at"`
	UpdatedAt       time.Time                `json:"updated_at"`
}

type AssetAssignmentListResponse struct {
	ID                 uint                           `json:"id"`
	AssetID            uint                           `json:"asset_id"`
	AssetName          string                         `json:"asset_name"`
	EmployeeID         uint                           `json:"employee_id"`
	EmployeeName       string                         `json:"employee_name"`
	EmployeeNIK        string                         `json:"employee_nik"`
	Purpose            string                         `json:"purpose"`
	ExpectedReturnDate *string                         `json:"expected_return_date"`
	ActualReturnDate   *string                         `json:"actual_return_date"`
	Status             constants.AssetAssignmentStatus `json:"status"`
	CreatedAt          time.Time                      `json:"created_at"`
}

type AssetAssignmentDetailResponse struct {
	ID                 uint                           `json:"id"`
	AssetID            uint                           `json:"asset_id"`
	AssetName          string                         `json:"asset_name"`
	EmployeeID         uint                           `json:"employee_id"`
	EmployeeName       string                         `json:"employee_name"`
	EmployeeNIK        string                         `json:"employee_nik"`
	Purpose            string                         `json:"purpose"`
	ExpectedReturnDate *string                         `json:"expected_return_date"`
	ActualReturnDate   *string                         `json:"actual_return_date"`
	Notes              string                         `json:"notes"`
	Status             constants.AssetAssignmentStatus `json:"status"`
	RejectionReason    string                         `json:"rejection_reason"`
	CreatedAt          time.Time                      `json:"created_at"`
}
