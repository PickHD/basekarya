package rbac

import (
	"basekarya-backend/internal/infrastructure"
	"basekarya-backend/pkg/constants"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type Service interface {
	CreateRole(ctx context.Context, req *CreateRoleRequest) error
	GetRolePermissions(ctx context.Context, roleID uint) (*RolePermissionsResponse, error)
	AssignPermissions(ctx context.Context, roleID uint, req *AssignPermissionsRequest) error
	GetAllPermissions(ctx context.Context) ([]PermissionResponse, error)
	GetAllRoles(ctx context.Context) ([]RoleResponse, error)
}

type service struct {
	repo               Repository
	cache              CacheProvider
	transactionManager infrastructure.TransactionManager
}

func NewService(repo Repository, cache CacheProvider, transactionManager infrastructure.TransactionManager) Service {
	return &service{
		repo:               repo,
		cache:              cache,
		transactionManager: transactionManager,
	}
}

func (s *service) CreateRole(ctx context.Context, req *CreateRoleRequest) error {
	return s.transactionManager.RunInTransaction(ctx, func(ctx context.Context) error {
		_, err := s.repo.FindRoleByName(ctx, req.Name)
		if err == nil {
			return errors.New("role already exists")
		}

		err = s.repo.Create(ctx, Role{
			Name: req.Name,
		})
		if err != nil {
			return errors.New("failed to create role")
		}

		err = s.cache.Del(ctx, constants.ROLE_CACHE_KEY)
		if err != nil {
			return errors.New("failed to delete role cache")
		}

		return nil
	})
}

func (s *service) GetRolePermissions(ctx context.Context, roleID uint) (*RolePermissionsResponse, error) {
	cacheKey := fmt.Sprintf(constants.ROLE_PERMISSION_CACHE_KEY, roleID)

	cacheData, err := s.cache.Get(ctx, cacheKey)
	if err == redis.Nil {
		role, err := s.repo.FindRoleByID(ctx, roleID)
		if err != nil {
			return nil, errors.New("role not found")
		}

		permissions := role.Permissions
		if permissions == nil {
			permissions = []Permission{}
		}

		parsedData, err := json.Marshal(&RolePermissionsResponse{
			RoleID:      role.ID,
			RoleName:    role.Name,
			Permissions: permissions,
		})
		if err != nil {
			return nil, err
		}

		err = s.cache.Set(ctx, cacheKey, parsedData, 24*time.Hour)
		if err != nil {
			return nil, err
		}

		return &RolePermissionsResponse{
			RoleID:      role.ID,
			RoleName:    role.Name,
			Permissions: permissions,
		}, nil
	} else if err != nil {
		return nil, err
	}

	var resp RolePermissionsResponse
	err = json.Unmarshal([]byte(cacheData), &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func (s *service) AssignPermissions(ctx context.Context, roleID uint, req *AssignPermissionsRequest) error {
	return s.transactionManager.RunInTransaction(ctx, func(ctx context.Context) error {
		_, err := s.repo.FindRoleByID(ctx, roleID)
		if err != nil {
			return errors.New("role not found")
		}

		permissions, err := s.repo.FindPermissionsByIDs(ctx, req.PermissionIDs)
		if err != nil {
			return err
		}

		if len(permissions) != len(req.PermissionIDs) {
			return errors.New("one or more permissions are invalid")
		}

		err = s.repo.ReplacingRolePermissions(ctx, roleID, req.PermissionIDs)
		if err != nil {
			return err
		}

		err = s.cache.Del(ctx, fmt.Sprintf(constants.ROLE_PERMISSION_CACHE_KEY, roleID))
		if err != nil {
			return err
		}

		return nil
	})
}

func (s *service) GetAllPermissions(ctx context.Context) ([]PermissionResponse, error) {
	cacheData, err := s.cache.Get(ctx, constants.PERMISSION_CACHE_KEY)
	if err == redis.Nil {
		data, err := s.repo.FindAllPermissions(ctx)
		if err != nil {
			return nil, err
		}

		if len(data) == 0 {
			return []PermissionResponse{}, nil
		}

		permissions := make([]PermissionResponse, len(data))

		for i, perm := range data {
			permissions[i] = PermissionResponse{
				ID:   perm.ID,
				Name: perm.Name,
			}
		}

		parsedData, err := json.Marshal(permissions)
		if err != nil {
			return nil, err
		}

		err = s.cache.Set(ctx, constants.PERMISSION_CACHE_KEY, parsedData, 24*time.Hour)
		if err != nil {
			return nil, err
		}

		return permissions, nil
	} else if err != nil {
		return nil, err
	}

	var permissions []PermissionResponse
	err = json.Unmarshal([]byte(cacheData), &permissions)
	if err != nil {
		return nil, err
	}

	return permissions, nil
}

func (s *service) GetAllRoles(ctx context.Context) ([]RoleResponse, error) {
	cacheData, err := s.cache.Get(ctx, constants.ROLE_CACHE_KEY)
	if err == redis.Nil {
		data, err := s.repo.FindAllRoles(ctx)
		if err != nil {
			return nil, err
		}

		if len(data) == 0 {
			return []RoleResponse{}, nil
		}

		roles := make([]RoleResponse, len(data))

		for i, role := range data {
			roles[i] = RoleResponse{
				ID:   role.ID,
				Name: role.Name,
			}
		}

		parsedData, err := json.Marshal(roles)
		if err != nil {
			return nil, err
		}

		err = s.cache.Set(ctx, constants.ROLE_CACHE_KEY, parsedData, 24*time.Hour)
		if err != nil {
			return nil, err
		}

		return roles, nil
	} else if err != nil {
		return nil, err
	}

	var roles []RoleResponse
	err = json.Unmarshal([]byte(cacheData), &roles)
	if err != nil {
		return nil, err
	}

	return roles, nil
}
