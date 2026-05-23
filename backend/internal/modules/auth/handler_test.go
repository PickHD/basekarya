package auth

import (
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"basekarya-backend/internal/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestHandler_Login(t *testing.T) {
	tests := []struct {
		name       string
		body       interface{}
		setupMocks func(*mockService)
		wantStatus int
	}{
		{
			name: "success",
			body: LoginRequest{
				Username: "admin",
				Password: "pass123",
			},
			setupMocks: func(svc *mockService) {
				svc.On("Login", mock.Anything, "admin", "pass123").Return(&LoginResponse{
					Token:              "jwt-token",
					MustChangePassword: false,
				}, nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "invalid credentials",
			body: LoginRequest{
				Username: "admin",
				Password: "wrong",
			},
			setupMocks: func(svc *mockService) {
				svc.On("Login", mock.Anything, "admin", "wrong").Return(nil, errors.New("invalid credentials"))
			},
			wantStatus: http.StatusInternalServerError,
		},
		{
			name:       "missing fields",
			body:       LoginRequest{},
			setupMocks: func(svc *mockService) {},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := new(mockService)
			tt.setupMocks(svc)
			handler := NewHandler(svc)

			at := testutil.NewAPITest(t, http.MethodPost, "/api/auth/login", tt.body)
			rec, err := at.Execute(handler.Login)

			require.NoError(t, err)
			assert.Equal(t, tt.wantStatus, rec.Code)

			var resp map[string]interface{}
			require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
			if tt.wantStatus < 400 {
				assert.Nil(t, resp["error"])
			}
		})
	}
}

func TestHandler_ForgotPassword(t *testing.T) {
	tests := []struct {
		name       string
		body       interface{}
		setupMocks func(*mockService)
		wantStatus int
	}{
		{
			name: "success",
			body: SendOrResendOTPRequest{
				Email: "user@example.com",
			},
			setupMocks: func(svc *mockService) {
				svc.On("SendOrResendOTP", mock.Anything, mock.AnythingOfType("*auth.SendOrResendOTPRequest")).Return(nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "employee not found",
			body: SendOrResendOTPRequest{
				Email: "unknown@example.com",
			},
			setupMocks: func(svc *mockService) {
				svc.On("SendOrResendOTP", mock.Anything, mock.AnythingOfType("*auth.SendOrResendOTPRequest")).Return(errors.New("employee not found"))
			},
			wantStatus: http.StatusInternalServerError,
		},
		{
			name:       "invalid email",
			body:       SendOrResendOTPRequest{},
			setupMocks: func(svc *mockService) {},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := new(mockService)
			tt.setupMocks(svc)
			handler := NewHandler(svc)

			at := testutil.NewAPITest(t, http.MethodPost, "/api/auth/forgot-password", tt.body)
			rec, err := at.Execute(handler.ForgotPassword)

			require.NoError(t, err)
			assert.Equal(t, tt.wantStatus, rec.Code)
		})
	}
}

func TestHandler_VerifyOTP(t *testing.T) {
	tests := []struct {
		name       string
		body       interface{}
		setupMocks func(*mockService)
		wantStatus int
	}{
		{
			name: "valid otp",
			body: VerifyOTPRequest{
				Code: "123456",
			},
			setupMocks: func(svc *mockService) {
				svc.On("VerifyOTP", mock.Anything, mock.AnythingOfType("*auth.VerifyOTPRequest")).Return(&VerifyOTPResponse{IsValid: true}, nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "invalid otp",
			body: VerifyOTPRequest{
				Code: "000000",
			},
			setupMocks: func(svc *mockService) {
				svc.On("VerifyOTP", mock.Anything, mock.AnythingOfType("*auth.VerifyOTPRequest")).Return(&VerifyOTPResponse{IsValid: false}, nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name:       "missing code",
			body:       VerifyOTPRequest{},
			setupMocks: func(svc *mockService) {},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := new(mockService)
			tt.setupMocks(svc)
			handler := NewHandler(svc)

			at := testutil.NewAPITest(t, http.MethodPost, "/api/auth/verify-otp", tt.body)
			rec, err := at.Execute(handler.VerifyOTP)

			require.NoError(t, err)
			assert.Equal(t, tt.wantStatus, rec.Code)
		})
	}
}

func TestHandler_ResetPassword(t *testing.T) {
	tests := []struct {
		name       string
		body       interface{}
		setupMocks func(*mockService)
		wantStatus int
	}{
		{
			name: "success",
			body: ResetPasswordRequest{
				Email:    "test@email.com",
				Code:     "123456",
				Password: "newpass123",
			},
			setupMocks: func(svc *mockService) {
				svc.On("ResetPassword", mock.Anything, mock.AnythingOfType("*auth.ResetPasswordRequest")).Return(nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "invalid otp",
			body: ResetPasswordRequest{
				Email:    "test@email.com",
				Code:     "000000",
				Password: "newpass123",
			},
			setupMocks: func(svc *mockService) {
				svc.On("ResetPassword", mock.Anything, mock.AnythingOfType("*auth.ResetPasswordRequest")).Return(errors.New("invalid OTP"))
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "missing fields",
			body:       ResetPasswordRequest{},
			setupMocks: func(svc *mockService) {},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := new(mockService)
			tt.setupMocks(svc)
			handler := NewHandler(svc)

			at := testutil.NewAPITest(t, http.MethodPost, "/api/auth/reset-password", tt.body)
			rec, err := at.Execute(handler.ResetPassword)

			require.NoError(t, err)
			assert.Equal(t, tt.wantStatus, rec.Code)
		})
	}
}

func TestHandler_RegisterCompany(t *testing.T) {
	tests := []struct {
		name       string
		body       interface{}
		setupMocks func(*mockService)
		wantStatus int
	}{
		{
			name: "success",
			body: RegisterCompanyRequest{
				CompanyName: "Test Co",
				AdminName:   "John",
				AdminEmail:  "john@test.com",
				Password:    "pass123456",
				PhoneNumber: "081234567890",
				PlanSlug:    "free",
			},
			setupMocks: func(svc *mockService) {
				svc.On("RegisterCompany", mock.Anything, mock.AnythingOfType("*auth.RegisterCompanyRequest")).Return(&RegisterCompanyResponse{
					Username: "john_abc",
				}, nil)
			},
			wantStatus: http.StatusCreated,
		},
		{
			name: "email already registered",
			body: RegisterCompanyRequest{
				CompanyName: "Test Co",
				AdminName:   "John",
				AdminEmail:  "john@test.com",
				Password:    "pass123456",
				PhoneNumber: "081234567890",
				PlanSlug:    "free",
			},
			setupMocks: func(svc *mockService) {
				svc.On("RegisterCompany", mock.Anything, mock.AnythingOfType("*auth.RegisterCompanyRequest")).Return(nil, errors.New("email already registered"))
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "missing required fields",
			body: RegisterCompanyRequest{
				CompanyName: "Test Co",
			},
			setupMocks: func(svc *mockService) {},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := new(mockService)
			tt.setupMocks(svc)
			handler := NewHandler(svc)

			at := testutil.NewAPITest(t, http.MethodPost, "/api/auth/register", tt.body)
			rec, err := at.Execute(handler.RegisterCompany)

			require.NoError(t, err)
			assert.Equal(t, tt.wantStatus, rec.Code)

			var resp map[string]interface{}
			require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
			if tt.wantStatus < 400 {
				assert.Nil(t, resp["error"])
			}
		})
	}
}
