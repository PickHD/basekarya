package user

type UserProfileResponse struct {
	ID                 uint    `json:"id"`
	Username           string  `json:"username"`
	Role               string  `json:"role"`
	FullName           string  `json:"full_name"`
	NIK                string  `json:"nik"`
	DepartmentName     string  `json:"department_name"`
	ShiftName          string  `json:"shift_name"`
	ShiftStartTime     string  `json:"shift_start_time"`
	ShiftEndTime       string  `json:"shift_end_time"`
	PhoneNumber        string  `json:"phone_number"`
	ProfilePictureUrl  string  `json:"profile_picture_url"`
	MustChangePassword bool    `json:"must_change_password"`
	BankName           string  `json:"bank_name"`
	BankAccountNumber  string  `json:"bank_account_number"`
	BankAccountHolder  string  `json:"bank_account_holder"`
	NPWP               string  `json:"npwp"`
	BaseSalary         float64 `json:"base_salary"`
}

type UpdateProfileRequest struct {
	FullName          string `form:"full_name" json:"full_name" validate:"omitempty,min=3,max=100"`
	PhoneNumber       string `form:"phone_number" json:"phone_number" validate:"omitempty,numeric,min=10,max=15"`
	BankName          string `form:"bank_name" json:"bank_name" validate:"omitempty,min=1"`
	BankAccountNumber string `form:"bank_account_number" json:"bank_account_number" validate:"omitempty,min=3"`
	BankAccountHolder string `form:"bank_account_holder" json:"bank_account_holder" validate:"omitempty,min=3"`
	NPWP              string `form:"npwp" json:"npwp" validate:"omitempty,min=12"`
}

type ChangePasswordRequest struct {
	OldPassword     string `json:"old_password" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required,min=6,max=72,nefield=OldPassword"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=NewPassword"`
}

type EmployeeListResponse struct {
	ID             uint    `json:"id"`
	FullName       string  `json:"full_name"`
	NIK            string  `json:"nik"`
	Username       string  `json:"username"`
	DepartmentName string  `json:"department_name"`
	ShiftName      string  `json:"shift_name"`
	BaseSalary     float64 `json:"base_salary"`
}

type CreateEmployeeRequest struct {
	Username     string  `json:"username" validate:"required"`
	FullName     string  `json:"full_name" validate:"required"`
	NIK          string  `json:"nik" validate:"required"`
	DepartmentID uint    `json:"department_id" validate:"required"`
	ShiftID      uint    `json:"shift_id" validate:"required"`
	BaseSalary   float64 `json:"base_salary" validate:"required"`
}

type UpdateEmployeeRequest struct {
	FullName     string  `json:"full_name"`
	NIK          string  `json:"nik"`
	DepartmentID uint    `json:"department_id"`
	ShiftID      uint    `json:"shift_id"`
	BaseSalary   float64 `json:"base_salary"`
}
