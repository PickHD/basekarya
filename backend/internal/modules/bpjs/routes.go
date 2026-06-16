package bpjs

import (
	"basekarya-backend/internal/middleware"
	"basekarya-backend/pkg/constants"

	"github.com/labstack/echo/v4"
)

func RegisterRoutes(e *echo.Group, handler *Handler, auth *middleware.AuthMiddleware) {
	g := e.Group("/admin/bpjs")
	g.GET("/configs", handler.List, auth.GrantPermission(constants.VIEW_BPJS_CONFIG))
	g.POST("/configs", handler.Create, auth.GrantPermission(constants.MANAGE_BPJS_CONFIG))
	g.GET("/configs/:id", handler.GetByID, auth.GrantPermission(constants.VIEW_BPJS_CONFIG))
	g.PUT("/configs/:id", handler.Update, auth.GrantPermission(constants.MANAGE_BPJS_CONFIG))
	g.DELETE("/configs/:id", handler.Delete, auth.GrantPermission(constants.MANAGE_BPJS_CONFIG))
}
