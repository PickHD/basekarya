package subscription

import (
	"context"
)

type Repository interface {
	FindAllPlans(ctx context.Context) ([]SubscriptionPlan, error)
	FindPlanBySlug(ctx context.Context, slug string) (*SubscriptionPlan, error)
	FindPlanByID(ctx context.Context, id uint) (*SubscriptionPlan, error)
	CreateRequest(ctx context.Context, req *SubscriptionRequest) error
	FindPendingRequestByCompanyID(ctx context.Context, companyID uint) (*SubscriptionRequest, error)
	FindRequestByID(ctx context.Context, id uint) (*SubscriptionRequest, error)
	FindAllPendingRequests(ctx context.Context) ([]SubscriptionRequestResponse, error)
	FindAllRequests(ctx context.Context) ([]SubscriptionRequestResponse, error)
	UpdateRequest(ctx context.Context, req *SubscriptionRequest) error
	FindAllCompanies(ctx context.Context, search string) ([]CompanyListItem, error)
	FindCompanyDetailByID(ctx context.Context, id uint) (*CompanyDetail, error)
	UpdateCompanyStatus(ctx context.Context, companyID uint, status string) error
	GetDashboardStats(ctx context.Context) (*DashboardStatsResponse, error)
}

type RoleProvider interface {
	FindPermissionIDsByGroupNames(ctx context.Context, groupNames []string) ([]uint, error)
	FindRoleIDsByCompanyID(ctx context.Context, companyID uint) ([]uint, error)
	AssignPermissions(ctx context.Context, roleID uint, permissionIDs []uint, companyID uint) error
}

type CacheProvider interface {
	FlushDB(ctx context.Context) error
}

type UserProvider interface {
	ForceResetPasswordByCompanyID(ctx context.Context, companyID uint) error
}
