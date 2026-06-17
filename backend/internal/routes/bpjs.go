package routes

import (
	"basekarya-backend/pkg/constants"

	"github.com/labstack/echo/v4"
)

func (r *Router) SetupBpjsRoutes(e *echo.Group) {
	e.GET("/configs", r.container.BpjsHandler.List, r.container.AuthMiddleware.GrantPermission(constants.VIEW_BPJS_CONFIG))
	e.POST("/configs", r.container.BpjsHandler.Create, r.container.AuthMiddleware.GrantPermission(constants.MANAGE_BPJS_CONFIG))
	e.GET("/configs/:id", r.container.BpjsHandler.GetByID, r.container.AuthMiddleware.GrantPermission(constants.VIEW_BPJS_CONFIG))
	e.PUT("/configs/:id", r.container.BpjsHandler.Update, r.container.AuthMiddleware.GrantPermission(constants.MANAGE_BPJS_CONFIG))
	e.DELETE("/configs/:id", r.container.BpjsHandler.Delete, r.container.AuthMiddleware.GrantPermission(constants.MANAGE_BPJS_CONFIG))
}
