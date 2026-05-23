package auth

import (
	"basekarya-backend/internal/modules/rbac"
	"basekarya-backend/internal/modules/user"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLogin_Success(t *testing.T) {
	svc, userProv, hasher, tokenProv, _, _, _, _, _ := newTestAuthService()
	ctx := context.Background()

	userProv.On("FindByUsername", ctx, "admin").Return(&user.User{
		ID: 1, Username: "admin", PasswordHash: "hashed", Role: &rbac.Role{Name: "SUPERADMIN"}, IsPlatformAdmin: true,
	}, nil)
	hasher.On("CheckPasswordHash", "pass123", "hashed").Return(true)
	tokenProv.On("GenerateToken", uint(1), uint(0), true, "SUPERADMIN", (*uint)(nil), []string(nil)).Return("jwt-token", nil)

	resp, err := svc.Login(ctx, "admin", "pass123")

	require.NoError(t, err)
	assert.Equal(t, "jwt-token", resp.Token)
	assert.False(t, resp.MustChangePassword)
}

func TestLogin_UserNotFound(t *testing.T) {
	svc, userProv, _, _, _, _, _, _, _ := newTestAuthService()
	ctx := context.Background()

	userProv.On("FindByUsername", ctx, "unknown").Return(nil, errors.New("not found"))

	resp, err := svc.Login(ctx, "unknown", "any")

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, "invalid credentials", err.Error())
}

func TestLogin_WrongPassword(t *testing.T) {
	svc, userProv, hasher, _, _, _, _, _, _ := newTestAuthService()
	ctx := context.Background()

	userProv.On("FindByUsername", ctx, "admin").Return(&user.User{
		ID: 1, Username: "admin", PasswordHash: "hashed", Role: &rbac.Role{Name: "SUPERADMIN"},
	}, nil)
	hasher.On("CheckPasswordHash", "wrong", "hashed").Return(false)

	resp, err := svc.Login(ctx, "admin", "wrong")

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, "invalid credentials", err.Error())
}

func TestVerifyOTP_Valid(t *testing.T) {
	svc, _, _, _, cache, _, _, _, _ := newTestAuthService()
	ctx := context.Background()

	cache.On("Get", ctx, "123456").Return("user@email.com", nil)

	resp, err := svc.VerifyOTP(ctx, &VerifyOTPRequest{Code: "123456"})

	require.NoError(t, err)
	assert.True(t, resp.IsValid)
}

func TestVerifyOTP_Invalid(t *testing.T) {
	svc, _, _, _, cache, _, _, _, _ := newTestAuthService()
	ctx := context.Background()

	cache.On("Get", ctx, "000000").Return("", errors.New("not found"))

	resp, err := svc.VerifyOTP(ctx, &VerifyOTPRequest{Code: "000000"})

	require.NoError(t, err)
	assert.False(t, resp.IsValid)
}

func TestResetPassword_Success(t *testing.T) {
	svc, userProv, hasher, _, cache, _, _, _, _ := newTestAuthService()
	ctx := context.Background()

	cache.On("Get", ctx, "123456").Return("user@email.com", nil)
	hasher.On("HashPassword", "newpass123").Return("newhash", nil)
	userProv.On("UpdatePasswordByEmail", ctx, "user@email.com", "newhash").Return(nil)

	err := svc.ResetPassword(ctx, &ResetPasswordRequest{Code: "123456", Password: "newpass123"})

	require.NoError(t, err)
}

func TestResetPassword_InvalidOTP(t *testing.T) {
	svc, _, _, _, cache, _, _, _, _ := newTestAuthService()
	ctx := context.Background()

	cache.On("Get", ctx, "000000").Return("", errors.New("not found"))

	err := svc.ResetPassword(ctx, &ResetPasswordRequest{Code: "000000", Password: "newpass123"})

	assert.Error(t, err)
	assert.Equal(t, "invalid OTP", err.Error())
}
