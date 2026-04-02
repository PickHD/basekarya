package auth

import (
	"basekarya-backend/pkg/constants"
	"basekarya-backend/pkg/utils"
	"context"
	"errors"
	"fmt"
	"time"
)

type Service interface {
	Login(ctx context.Context, username, password string) (*LoginResponse, error)
	SendOrResendOTP(ctx context.Context, req *SendOrResendOTPRequest) error
	VerifyOTP(ctx context.Context, req *VerifyOTPRequest) (*VerifyOTPResponse, error)
	ResetPassword(ctx context.Context, req *ResetPasswordRequest) error
}

type service struct {
	user   UserProvider
	hasher Hasher
	token  TokenProvider
	cache  CacheProvider
	email  EmailProvider
}

func NewService(user UserProvider, hasher Hasher, token TokenProvider, cache CacheProvider, email EmailProvider) Service {
	return &service{
		user:   user,
		hasher: hasher,
		token:  token,
		cache:  cache,
		email:  email,
	}
}

func (s *service) Login(ctx context.Context, username, password string) (*LoginResponse, error) {
	foundUser, err := s.user.FindByUsername(ctx, username)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if !s.hasher.CheckPasswordHash(password, foundUser.PasswordHash) {
		return nil, errors.New("invalid credentials")
	}

	var employeeID *uint
	if foundUser.Employee != nil && foundUser.Role != nil && foundUser.Role.Name != string(constants.UserRoleSuperadmin) {
		employeeID = &foundUser.Employee.ID
	}

	var permissions []string
	if foundUser.Role != nil {
		for _, permission := range foundUser.Role.Permissions {
			permissions = append(permissions, permission.Name)
		}
	}

	tokenString, err := s.token.GenerateToken(foundUser.ID, foundUser.Role.Name, employeeID, permissions)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		Token:              tokenString,
		MustChangePassword: foundUser.MustChangePassword,
	}, nil
}

func (s *service) SendOrResendOTP(ctx context.Context, req *SendOrResendOTPRequest) error {
	_, err := s.user.FindEmployeeByEmail(ctx, req.Email)
	if err != nil {
		return errors.New("employee not found")
	}

	code := utils.GenerateRandomNumber(6)
	subject := fmt.Sprintf("Basekarya - Kode OTP: %s", code)
	htmlBody := fmt.Sprintf(`
		<h1>Basekarya - Kode OTP</h1>
		<p>Kode OTP Anda adalah: <strong>%s</strong></p>
		<p>kode akan kadaluarsa dalam 5 menit</p>
	`, code)

	err = s.email.Send(req.Email, subject, htmlBody)
	if err != nil {
		return err
	}

	err = s.cache.Set(ctx, code, req.Email, 5*time.Minute)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) VerifyOTP(ctx context.Context, req *VerifyOTPRequest) (*VerifyOTPResponse, error) {
	_, err := s.cache.Get(ctx, req.Code)
	if err != nil {
		return &VerifyOTPResponse{
			IsValid: false,
		}, nil
	}

	return &VerifyOTPResponse{
		IsValid: true,
	}, nil
}

func (s *service) ResetPassword(ctx context.Context, req *ResetPasswordRequest) error {
	email, err := s.cache.Get(ctx, req.Code)
	if err != nil {
		return errors.New("invalid OTP")
	}

	passwordHash, err := s.hasher.HashPassword(req.Password)
	if err != nil {
		return err
	}

	err = s.user.UpdatePasswordByEmail(ctx, email, passwordHash)
	if err != nil {
		return err
	}

	err = s.cache.Del(ctx, req.Code)
	if err != nil {
		return err
	}

	return nil
}
