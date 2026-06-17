package routes

import (
	"basekarya-backend/pkg/constants"

	"github.com/labstack/echo/v4"
)

func (r *Router) SetupTaxRoutes(e *echo.Group) {
	e.GET("/ter-brackets", r.container.TaxHandler.ListTERBrackets, r.container.AuthMiddleware.GrantPermission(constants.VIEW_TAX_CONFIG))
	e.POST("/ter-brackets", r.container.TaxHandler.CreateTERBracket, r.container.AuthMiddleware.GrantPermission(constants.MANAGE_TAX_CONFIG))
	e.GET("/ter-brackets/:id", r.container.TaxHandler.GetTERBracketByID, r.container.AuthMiddleware.GrantPermission(constants.VIEW_TAX_CONFIG))
	e.PUT("/ter-brackets/:id", r.container.TaxHandler.UpdateTERBracket, r.container.AuthMiddleware.GrantPermission(constants.MANAGE_TAX_CONFIG))
	e.DELETE("/ter-brackets/:id", r.container.TaxHandler.DeleteTERBracket, r.container.AuthMiddleware.GrantPermission(constants.MANAGE_TAX_CONFIG))
	e.GET("/ptkp-configs", r.container.TaxHandler.ListPTKPConfigs, r.container.AuthMiddleware.GrantPermission(constants.VIEW_TAX_CONFIG))
	e.POST("/ptkp-configs", r.container.TaxHandler.CreatePTKPConfig, r.container.AuthMiddleware.GrantPermission(constants.MANAGE_TAX_CONFIG))
	e.PUT("/ptkp-configs/:id", r.container.TaxHandler.UpdatePTKPConfig, r.container.AuthMiddleware.GrantPermission(constants.MANAGE_TAX_CONFIG))
	e.DELETE("/ptkp-configs/:id", r.container.TaxHandler.DeletePTKPConfig, r.container.AuthMiddleware.GrantPermission(constants.MANAGE_TAX_CONFIG))
}
