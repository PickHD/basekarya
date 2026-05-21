package routes

import (
	"basekarya-backend/pkg/constants"

	"github.com/labstack/echo/v4"
)

func (r *Router) SetupFinanceRoutes(e *echo.Group) {
	e.GET("/dashboard", r.container.FinanceHandler.GetDashboard, r.container.AuthMiddleware.GrantPermission(constants.VIEW_FINANCE_DASHBOARD))

	e.GET("/categories", r.container.FinanceHandler.GetCategories, r.container.AuthMiddleware.GrantPermission(constants.MANAGE_FINANCE_CATEGORY))
	e.POST("/categories", r.container.FinanceHandler.CreateCategory, r.container.AuthMiddleware.GrantPermission(constants.MANAGE_FINANCE_CATEGORY))
	e.PUT("/categories/:id", r.container.FinanceHandler.UpdateCategory, r.container.AuthMiddleware.GrantPermission(constants.MANAGE_FINANCE_CATEGORY))
	e.DELETE("/categories/:id", r.container.FinanceHandler.DeleteCategory, r.container.AuthMiddleware.GrantPermission(constants.MANAGE_FINANCE_CATEGORY))

	e.GET("/transactions", r.container.FinanceHandler.GetAllTransactions, r.container.AuthMiddleware.GrantPermission(constants.VIEW_FINANCE))
	e.GET("/transactions/export", r.container.FinanceHandler.ExportTransactions, r.container.AuthMiddleware.GrantPermission(constants.EXPORT_FINANCE))
	e.GET("/transactions/:id", r.container.FinanceHandler.GetTransactionDetail, r.container.AuthMiddleware.GrantPermission(constants.VIEW_FINANCE))
	e.POST("/transactions", r.container.FinanceHandler.CreateTransaction, r.container.AuthMiddleware.GrantPermission(constants.CREATE_FINANCE))
	e.PUT("/transactions/:id/action", r.container.FinanceHandler.ProcessAction, r.container.AuthMiddleware.GrantPermission(constants.APPROVAL_FINANCE))
}
