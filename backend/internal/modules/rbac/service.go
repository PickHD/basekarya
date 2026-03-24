package rbac

import (
	"basekarya-backend/internal/infrastructure"
	"context"
	"errors"
)

type Service interface {
	CreateRole(ctx context.Context, req *CreateRoleRequest) error
	GetRolePermissions(ctx context.Context, roleID uint) (*RolePermissionsResponse, error)
	AssignPermissions(ctx context.Context, roleID uint, req *AssignPermissionsRequest) error
}

type service struct {
	repo               Repository
	transactionManager infrastructure.TransactionManager
}

func NewService(repo Repository, transactionManager infrastructure.TransactionManager) Service {
	return &service{
		repo:               repo,
		transactionManager: transactionManager,
	}
}

func (s *service) CreateRole(ctx context.Context, req *CreateRoleRequest) error {
	return s.transactionManager.RunInTransaction(ctx, func(ctx context.Context) error {
		_, err := s.repo.FindRoleByName(ctx, req.Name)
		if err == nil {
			return errors.New("role already exists")
		}

		return s.repo.Create(ctx, Role{
			Name: req.Name,
		})
	})
}

func (s *service) GetRolePermissions(ctx context.Context, roleID uint) (*RolePermissionsResponse, error) {
	role, err := s.repo.FindRoleByID(ctx, roleID)
	if err != nil {
		return nil, errors.New("role not found")
	}

	permissions := role.Permissions
	if permissions == nil {
		permissions = []Permission{}
	}

	return &RolePermissionsResponse{
		RoleID:      role.ID,
		RoleName:    role.Name,
		Permissions: permissions,
	}, nil
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

		return s.repo.ReplacingRolePermissions(ctx, roleID, req.PermissionIDs)
	})
}
