package routes

import (
	"basekarya-backend/pkg/constants"

	"github.com/labstack/echo/v4"
)

func (r *Router) SetupRoleRoutes(e *echo.Group) {
	e.GET("", r.container.RbacHandler.GetAllRoles, r.container.AuthMiddleware.GrantAnyPermission(constants.VIEW_ROLE, constants.VIEW_MASTER))
	e.POST("", r.container.RbacHandler.CreateRole, r.container.AuthMiddleware.GrantPermission(constants.CREATE_ROLE))
	e.GET("/:id/permissions", r.container.RbacHandler.GetRolePermissions, r.container.AuthMiddleware.GrantPermission(constants.VIEW_ROLE))
	e.PUT("/:id/permissions", r.container.RbacHandler.AssignPermissions, r.container.AuthMiddleware.GrantPermission(constants.ASSIGN_ROLE))
}

func (r *Router) SetupPermissionRoutes(e *echo.Group) {
	e.GET("", r.container.RbacHandler.GetAllPermissions, r.container.AuthMiddleware.GrantPermission(constants.VIEW_PERMISSION))
}
