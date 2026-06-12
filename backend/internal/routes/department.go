package routes

import (
	"basekarya-backend/pkg/constants"

	"github.com/labstack/echo/v4"
)

func (r *Router) SetupDepartmentRoutes(e *echo.Group) {
	e.GET("", r.container.DepartmentHandler.GetAll, r.container.AuthMiddleware.GrantPermission(constants.VIEW_MASTER))
	e.GET("/:id", r.container.DepartmentHandler.GetByID, r.container.AuthMiddleware.GrantPermission(constants.VIEW_MASTER))
	e.POST("", r.container.DepartmentHandler.Create, r.container.AuthMiddleware.GrantPermission(constants.MANAGE_MASTER))
	e.PUT("/:id", r.container.DepartmentHandler.Update, r.container.AuthMiddleware.GrantPermission(constants.MANAGE_MASTER))
	e.DELETE("/:id", r.container.DepartmentHandler.Delete, r.container.AuthMiddleware.GrantPermission(constants.MANAGE_MASTER))
}
