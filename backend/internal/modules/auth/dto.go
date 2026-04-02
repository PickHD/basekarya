package auth

type LoginRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50,alphanum"`
	Password string `json:"password" validate:"required"`
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
	Code     string `json:"code" validate:"required,len=6"`
	Password string `json:"password" validate:"required"`
}
