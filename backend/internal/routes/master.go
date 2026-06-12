package routes

import (
	"basekarya-backend/pkg/constants"

	"github.com/labstack/echo/v4"
)

func (r *Router) SetupMasterRoutes(e *echo.Group) {
	e.GET("/shifts", r.container.MasterHandler.GetShifts, r.container.AuthMiddleware.GrantPermission(constants.VIEW_MASTER))
	e.GET("/leaves/types", r.container.MasterHandler.GetLeaveTypes, r.container.AuthMiddleware.GrantPermission(constants.VIEW_MASTER))
}
