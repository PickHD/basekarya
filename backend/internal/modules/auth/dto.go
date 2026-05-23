package auth

type LoginRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Password string `json:"password" validate:"required,min=1,max=72"`
}

type LoginResponse struct {
	Token              string `json:"token"`
	MustChangePassword bool   `json:"must_change_password"`
}

type SendOrResendOTPRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type VerifyOTPRequest struct {
	Code string `json:"code" validate:"required,len=6"`
}

type VerifyOTPResponse struct {
	IsValid bool `json:"is_valid"`
}

type ResetPasswordRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Code     string `json:"code" validate:"required,len=6"`
	Password string `json:"password" validate:"required,min=6,max=72"`
}

type RegisterCompanyRequest struct {
	CompanyName string `json:"company_name" validate:"required,min=2,max=255"`
	AdminName   string `json:"admin_name" validate:"required,min=2,max=100"`
	AdminEmail  string `json:"admin_email" validate:"required,email"`
	Password    string `json:"password" validate:"required,min=6"`
	PhoneNumber string `json:"phone_number" validate:"required"`
	PlanSlug    string `json:"plan_slug" validate:"required,oneof=free basic pro"`
}

type RegisterCompanyResponse struct {
	Username string `json:"username"`
}
