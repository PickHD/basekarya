package auth

import (
	"basekarya-backend/internal/modules/company"
	"basekarya-backend/internal/modules/rbac"
	"basekarya-backend/internal/modules/user"
	"context"
	"time"
)

type Hasher interface {
	CheckPasswordHash(password, hash string) bool
	HashPassword(password string) (string, error)
}

type TokenProvider interface {
	GenerateToken(userID uint, companyID uint, isPlatformAdmin bool, role string, employeeID *uint, permissions []string) (string, error)
}

type UserProvider interface {
	FindByUsername(ctx context.Context, username string) (*user.User, error)
	FindEmployeeByEmail(ctx context.Context, email string) (*user.Employee, error)
	UpdatePasswordByEmail(ctx context.Context, email string, password string) error
	CreateUser(ctx context.Context, user *user.User) error
}

type CacheProvider interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value any, expiration time.Duration) error
	Del(ctx context.Context, key string) error
	FlushDB(ctx context.Context) error
}

type EmailProvider interface {
	Send(to, subject, htmlBody string) error
}

type CompanyProvider interface {
	CreateCompany(ctx context.Context, c *company.Company) error
	FindPlanIDBySlug(ctx context.Context, slug string) (uint, error)
}

type RoleProvider interface {
	Create(ctx context.Context, role *rbac.Role) error
	FindRoleByName(ctx context.Context, name string) (*rbac.Role, error)
	FindAllPermissionIDs(ctx context.Context) ([]uint, error)
	FindPermissionIDsByGroupNames(ctx context.Context, groupNames []string) ([]uint, error)
	AssignPermissions(ctx context.Context, roleID uint, permissionIDs []uint, companyID uint) error
}

type MasterProvider interface {
	SeedDefaults(ctx context.Context, companyID uint) error
}

type MasterSeeder interface {
	SeedDefaults(ctx context.Context, companyID uint) error
}
