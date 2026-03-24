package rbac

type CreateRoleRequest struct {
	Name string `json:"name" validate:"required"`
}

type RolePermissionsResponse struct {
	RoleID      uint         `json:"role_id"`
	RoleName    string       `json:"role_name"`
	Permissions []Permission `json:"permissions"`
}

type AssignPermissionsRequest struct {
	PermissionIDs []uint `json:"permission_ids" validate:"required,min=1"`
}
