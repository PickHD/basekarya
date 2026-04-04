package rbac

import "time"

type Role struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"unique;not null" json:"name"`
	CreatedAt time.Time `json:"created_at"`

	Permissions []Permission `gorm:"many2many:role_permissions;" json:"permissions,omitempty"`
}

type Permission struct {
	ID                uint      `gorm:"primaryKey" json:"id"`
	Name              string    `gorm:"unique;not null" json:"name"`
	DisplayName       string    `gorm:"not null" json:"display_name"`
	Description       string    `gorm:"null" json:"description"`
	PermissionGroupID uint      `gorm:"null" json:"permission_group_id"`
	CreatedAt         time.Time `json:"created_at"`

	PermissionGroup PermissionGroup `gorm:"foreignKey:PermissionGroupID;references:ID" json:"permission_group,omitempty"`
}

type PermissionGroup struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"unique;not null" json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type RolePermission struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	RoleID       uint      `gorm:"uniqueIndex:idx_role_permission" json:"role_id"`
	PermissionID uint      `gorm:"uniqueIndex:idx_role_permission" json:"permission_id"`
	CreatedAt    time.Time `json:"created_at"`
}

func (Role) TableName() string {
	return "roles"
}

func (Permission) TableName() string {
	return "permissions"
}

func (RolePermission) TableName() string {
	return "role_permissions"
}
