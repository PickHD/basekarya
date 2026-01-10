package user

type UserProfileResponse struct {
	ID                uint   `json:"id"`
	Username          string `json:"username"`
	Role              string `json:"role"`
	FullName          string `json:"full_name"`
	NIK               string `json:"nik"`
	DepartmentName    string `json:"department_name"`
	ShiftName         string `json:"shift_name"`
	ShiftStartTime    string `json:"shift_start_time"`
	ShiftEndTime      string `json:"shift_end_time"`
	PhoneNumber       string `json:"phone_number"`
	ProfilePictureUrl string `json:"profile_picture_url"`
}

type UpdateProfileRequest struct {
	PhoneNumber string `form:"phone_number"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=6"`
}

type EmployeeListResponse struct {
	ID             uint   `json:"id"`
	FullName       string `json:"full_name"`
	NIK            string `json:"nik"`
	Username       string `json:"username"`
	DepartmentName string `json:"department_name"`
	ShiftName      string `json:"shift_name"`
}

type CreateEmployeeRequest struct {
	Username     string `json:"username" validate:"required"`
	Password     string `json:"password" validate:"required,min=6"`
	FullName     string `json:"full_name" validate:"required"`
	NIK          string `json:"nik" validate:"required"`
	DepartmentID uint   `json:"department_id" validate:"required"`
	ShiftID      uint   `json:"shift_id" validate:"required"`
}

type UpdateEmployeeRequest struct {
	FullName     string `json:"full_name"`
	NIK          string `json:"nik"`
	DepartmentID uint   `json:"department_id"`
	ShiftID      uint   `json:"shift_id"`
}
