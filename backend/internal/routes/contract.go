package routes

import (
	"basekarya-backend/pkg/constants"

	"github.com/labstack/echo/v4"
)

func (r *Router) SetupContractRoutes(e *echo.Group) {
	e.GET("", r.container.ContractHandler.GetAll, r.container.AuthMiddleware.GrantPermission(constants.VIEW_CONTRACT))
	e.GET("/:id", r.container.ContractHandler.GetDetail, r.container.AuthMiddleware.GrantPermission(constants.VIEW_CONTRACT))
	e.GET("/employee/:employeeId", r.container.ContractHandler.GetByEmployee, r.container.AuthMiddleware.GrantPermission(constants.VIEW_CONTRACT))
	e.PUT("", r.container.ContractHandler.Upsert, r.container.AuthMiddleware.GrantPermission(constants.CREATE_CONTRACT))
	e.DELETE("/:id", r.container.ContractHandler.Delete, r.container.AuthMiddleware.GrantPermission(constants.UPDATE_CONTRACT))
	e.GET("/export", r.container.ContractHandler.Export, r.container.AuthMiddleware.GrantPermission(constants.EXPORT_CONTRACT))
}
