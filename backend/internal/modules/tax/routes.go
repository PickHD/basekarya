package tax

import (
	"basekarya-backend/internal/middleware"
	"basekarya-backend/pkg/constants"

	"github.com/labstack/echo/v4"
)

func RegisterRoutes(e *echo.Group, handler *Handler, auth *middleware.AuthMiddleware) {
	g := e.Group("/admin/tax")
	g.GET("/ter-brackets", handler.ListTERBrackets, auth.GrantPermission(constants.VIEW_TAX_CONFIG))
	g.POST("/ter-brackets", handler.CreateTERBracket, auth.GrantPermission(constants.MANAGE_TAX_CONFIG))
	g.GET("/ter-brackets/:id", handler.GetTERBracketByID, auth.GrantPermission(constants.VIEW_TAX_CONFIG))
	g.PUT("/ter-brackets/:id", handler.UpdateTERBracket, auth.GrantPermission(constants.MANAGE_TAX_CONFIG))
	g.DELETE("/ter-brackets/:id", handler.DeleteTERBracket, auth.GrantPermission(constants.MANAGE_TAX_CONFIG))
	g.GET("/ptkp-configs", handler.ListPTKPConfigs, auth.GrantPermission(constants.VIEW_TAX_CONFIG))
	g.POST("/ptkp-configs", handler.CreatePTKPConfig, auth.GrantPermission(constants.MANAGE_TAX_CONFIG))
	g.PUT("/ptkp-configs/:id", handler.UpdatePTKPConfig, auth.GrantPermission(constants.MANAGE_TAX_CONFIG))
	g.DELETE("/ptkp-configs/:id", handler.DeletePTKPConfig, auth.GrantPermission(constants.MANAGE_TAX_CONFIG))
}
