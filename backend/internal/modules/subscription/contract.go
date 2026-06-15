package subscription

import (
	"context"
)



type RoleProvider interface {
	FindPermissionIDsByGroupNames(ctx context.Context, groupNames []string) ([]uint, error)
	FindRoleIDsByCompanyID(ctx context.Context, companyID uint) ([]uint, error)
	AssignPermissions(ctx context.Context, roleID uint, permissionIDs []uint, companyID uint) error
}

type CacheProvider interface {
	Del(ctx context.Context, key string) error
}

type UserProvider interface {
	ForceResetPasswordByCompanyID(ctx context.Context, companyID uint) error
}
