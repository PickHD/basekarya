package rbac

import (
	"basekarya-backend/pkg/constants"
	"basekarya-backend/pkg/utils"
	"context"

	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, role Role) error
	FindRoleByID(ctx context.Context, id uint) (*Role, error)
	FindRoleByName(ctx context.Context, name string) (*Role, error)
	ReplacingRolePermissions(ctx context.Context, roleID uint, permissionIDs []uint) error
	FindPermissionsByIDs(ctx context.Context, ids []uint) ([]Permission, error)
	FindAllPermissions(ctx context.Context) ([]Permission, error)
	FindAllRoles(ctx context.Context) ([]Role, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) Create(ctx context.Context, role Role) error {
	db := utils.GetDBFromContext(ctx, r.db)

	return db.Create(&role).Error
}

func (r *repository) FindRoleByID(ctx context.Context, id uint) (*Role, error) {
	db := utils.GetDBFromContext(ctx, r.db)
	var role Role

	err := db.Preload("Permissions").First(&role, id).Error
	if err != nil {
		return nil, err
	}

	return &role, nil
}

func (r *repository) FindPermissionsByIDs(ctx context.Context, ids []uint) ([]Permission, error) {
	db := utils.GetDBFromContext(ctx, r.db)
	var permissions []Permission

	err := db.Where("id IN ?", ids).Find(&permissions).Error
	if err != nil {
		return nil, err
	}

	return permissions, nil
}

func (r *repository) FindRoleByName(ctx context.Context, name string) (*Role, error) {
	db := utils.GetDBFromContext(ctx, r.db)
	var role Role

	err := db.Where("name = ?", name).First(&role).Error
	if err != nil {
		return nil, err
	}

	return &role, nil
}

func (r *repository) ReplacingRolePermissions(ctx context.Context, roleID uint, permissionIDs []uint) error {
	db := utils.GetDBFromContext(ctx, r.db)

	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("role_id = ?", roleID).Delete(&RolePermission{}).Error; err != nil {
			return err
		}

		if len(permissionIDs) > 0 {
			var rolePermissions []RolePermission
			for _, pid := range permissionIDs {
				rolePermissions = append(rolePermissions, RolePermission{
					RoleID:       roleID,
					PermissionID: pid,
				})
			}

			if err := tx.Create(&rolePermissions).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *repository) FindAllPermissions(ctx context.Context) ([]Permission, error) {
	db := utils.GetDBFromContext(ctx, r.db)

	var permissions []Permission

	err := db.Where("name != ?", constants.VIEW_PERMISSION).Find(&permissions).Error
	if err != nil {
		return nil, err
	}

	return permissions, nil
}

func (r *repository) FindAllRoles(ctx context.Context) ([]Role, error) {
	db := utils.GetDBFromContext(ctx, r.db)

	var roles []Role

	err := db.Find(&roles).Error
	if err != nil {
		return nil, err
	}

	return roles, nil
}
