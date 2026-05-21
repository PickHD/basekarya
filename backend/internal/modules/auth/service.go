package auth

import (
	"basekarya-backend/internal/modules/company"
	"basekarya-backend/internal/modules/rbac"
	"basekarya-backend/internal/modules/user"
	"basekarya-backend/pkg/constants"
	"basekarya-backend/pkg/utils"
	"context"
	"errors"
	"fmt"
	"time"
)

type Service interface {
	Login(ctx context.Context, username, password string) (*LoginResponse, error)
	RegisterCompany(ctx context.Context, req *RegisterCompanyRequest) (*RegisterCompanyResponse, error)
	SendOrResendOTP(ctx context.Context, req *SendOrResendOTPRequest) error
	VerifyOTP(ctx context.Context, req *VerifyOTPRequest) (*VerifyOTPResponse, error)
	ResetPassword(ctx context.Context, req *ResetPasswordRequest) error
}

type service struct {
	user    UserProvider
	hasher  Hasher
	token   TokenProvider
	cache   CacheProvider
	email   EmailProvider
	company CompanyProvider
	role    RoleProvider
	master  MasterProvider
}

func NewService(user UserProvider, hasher Hasher, token TokenProvider, cache CacheProvider, email EmailProvider, company CompanyProvider, role RoleProvider, master MasterProvider) Service {
	return &service{
		user:    user,
		hasher:  hasher,
		token:   token,
		cache:   cache,
		email:   email,
		company: company,
		role:    role,
		master:  master,
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
	if foundUser.Employee != nil && foundUser.Role != nil && foundUser.Role.Name != "PLATFORM_ADMIN" {
		employeeID = &foundUser.Employee.ID
	}

	var permissions []string
	if foundUser.Role != nil {
		for _, permission := range foundUser.Role.Permissions {
			permissions = append(permissions, permission.Name)
		}
	}

	tokenString, err := s.token.GenerateToken(foundUser.ID, foundUser.CompanyID, foundUser.IsPlatformAdmin, foundUser.Role.Name, employeeID, permissions)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		Token:              tokenString,
		MustChangePassword: foundUser.MustChangePassword,
	}, nil
}

func (s *service) RegisterCompany(ctx context.Context, req *RegisterCompanyRequest) (*RegisterCompanyResponse, error) {
	existing, _ := s.user.FindByUsername(ctx, req.AdminEmail)
	if existing != nil && existing.ID != 0 {
		return nil, errors.New("email already registered")
	}

	planID, err := s.company.FindPlanIDBySlug(ctx, req.PlanSlug)
	if err != nil || planID == 0 {
		return nil, errors.New("invalid plan selected")
	}

	subscriptionStatus := constants.SubStatusActive
	if req.PlanSlug != "free" {
		subscriptionStatus = constants.SubStatusPendingPayment
	}

	newCompany := &company.Company{
		Name:               req.CompanyName,
		PhoneNumber:        req.PhoneNumber,
		Email:              req.AdminEmail,
		SubscriptionPlanID: &planID,
		SubscriptionStatus: subscriptionStatus,
	}
	if err := s.company.CreateCompany(ctx, newCompany); err != nil {
		return nil, errors.New("failed to create company")
	}

	if err := s.master.SeedDefaults(ctx, newCompany.ID); err != nil {
		return nil, errors.New("failed to seed master data")
	}

	superadminRole := &rbac.Role{
		Name:      "SUPERADMIN",
		CompanyID: newCompany.ID,
	}
	if err := s.role.Create(ctx, superadminRole); err != nil {
		return nil, errors.New("failed to create superadmin role")
	}

	employeeRole := &rbac.Role{
		Name:      "EMPLOYEE",
		CompanyID: newCompany.ID,
	}
	if err := s.role.Create(ctx, employeeRole); err != nil {
		return nil, errors.New("failed to create employee role")
	}

	allowedGroups := buildAllowedGroups(req.PlanSlug)
	permissionIDs, err := s.role.FindPermissionIDsByGroupNames(ctx, allowedGroups)
	if err != nil {
		return nil, errors.New("failed to get permissions")
	}
	if err := s.role.AssignPermissions(ctx, superadminRole.ID, permissionIDs, newCompany.ID); err != nil {
		return nil, errors.New("failed to assign permissions")
	}

	hashPass, err := s.hasher.HashPassword(req.Password)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	newUser := &user.User{
		Username:           utils.GenerateUsername(req.AdminName),
		PasswordHash:       hashPass,
		RoleID:             superadminRole.ID,
		CompanyID:          newCompany.ID,
		MustChangePassword: false,
		IsActive:           true,
	}
	if err := s.user.CreateUser(ctx, newUser); err != nil {
		return nil, errors.New("failed to create admin user")
	}

	_ = s.cache.FlushDB(context.Background())

	return &RegisterCompanyResponse{
		Username: newUser.Username,
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

	return nil
}

func buildAllowedGroups(planSlug string) []string {
	groups := make([]string, len(constants.AlwaysAvailableGroups))
	copy(groups, constants.AlwaysAvailableGroups)

	planModules, ok := constants.PlanModules[planSlug]
	if !ok {
		return groups
	}

	for _, mod := range planModules {
		if modGroups, exists := constants.ModulePermissionGroups[mod]; exists {
			groups = append(groups, modGroups...)
		}
	}

	return groups
}
