package routes

import (
	"basekarya-backend/internal/middleware"
	"basekarya-backend/pkg/constants"

	"github.com/labstack/echo/v4"
)

func (r *Router) SetupContractRoutes(e *echo.Group, sub *middleware.SubscriptionMiddleware) {
	g := e.Group("", sub.RequireModule("contract"))
	g.GET("", r.container.ContractHandler.GetAll, r.container.AuthMiddleware.GrantPermission(constants.VIEW_CONTRACT))
	g.GET("/:id", r.container.ContractHandler.GetDetail, r.container.AuthMiddleware.GrantPermission(constants.VIEW_CONTRACT))
	g.GET("/employee/:employeeId", r.container.ContractHandler.GetByEmployee, r.container.AuthMiddleware.GrantPermission(constants.VIEW_CONTRACT))
	g.PUT("", r.container.ContractHandler.Upsert, r.container.AuthMiddleware.GrantPermission(constants.CREATE_CONTRACT))
	g.DELETE("/:id", r.container.ContractHandler.Delete, r.container.AuthMiddleware.GrantPermission(constants.UPDATE_CONTRACT))
	g.GET("/export", r.container.ContractHandler.Export, r.container.AuthMiddleware.GrantPermission(constants.EXPORT_CONTRACT))
}
