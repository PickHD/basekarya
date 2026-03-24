package routes

import (
	"basekarya-backend/pkg/constants"

	"github.com/labstack/echo/v4"
)

func (r *Router) SetupOvertimeRoutes(e *echo.Group) {
	e.GET("", r.container.OvertimeHandler.GetAll, r.container.AuthMiddleware.GrantAnyPermission(constants.VIEW_OVERTIME, constants.VIEW_SELF_OVERTIME))
	e.GET("/export", r.container.OvertimeHandler.Export, r.container.AuthMiddleware.GrantPermission(constants.EXPORT_OVERTIME))
	e.POST("", r.container.OvertimeHandler.Create, r.container.AuthMiddleware.GrantPermission(constants.CREATE_OVERTIME))
	e.GET("/:id", r.container.OvertimeHandler.GetDetail, r.container.AuthMiddleware.GrantAnyPermission(constants.VIEW_OVERTIME, constants.VIEW_SELF_OVERTIME))
	e.PUT("/:id/action", r.container.OvertimeHandler.ProcessAction, r.container.AuthMiddleware.GrantPermission(constants.APPROVAL_OVERTIME))
}
