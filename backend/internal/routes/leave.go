package routes

import (
	"basekarya-backend/pkg/constants"

	"github.com/labstack/echo/v4"
)

func (r *Router) SetupLeaveRoutes(e *echo.Group) {
	e.GET("", r.container.LeaveHandler.GetAll, r.container.AuthMiddleware.GrantAnyPermission(constants.VIEW_LEAVE, constants.VIEW_SELF_LEAVE))
	e.GET("/:id", r.container.LeaveHandler.GetDetail, r.container.AuthMiddleware.GrantAnyPermission(constants.VIEW_LEAVE, constants.VIEW_SELF_LEAVE))
	e.POST("/apply", r.container.LeaveHandler.Apply, r.container.AuthMiddleware.GrantPermission(constants.CREATE_LEAVE))
	e.PUT("/:id/action", r.container.LeaveHandler.RequestAction, r.container.AuthMiddleware.GrantPermission(constants.APPROVAL_LEAVE))
	e.GET("/export", r.container.LeaveHandler.Export, r.container.AuthMiddleware.GrantPermission(constants.EXPORT_LEAVE))
}
