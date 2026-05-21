package routes

import (
	"basekarya-backend/internal/middleware"
	"basekarya-backend/pkg/constants"

	"github.com/labstack/echo/v4"
)

func (r *Router) SetupFinanceRoutes(e *echo.Group, sub *middleware.SubscriptionMiddleware) {
	g := e.Group("", sub.RequireModule("finance"))
	g.GET("/dashboard", r.container.FinanceHandler.GetDashboard, r.container.AuthMiddleware.GrantPermission(constants.VIEW_FINANCE_DASHBOARD))

	g.GET("/categories", r.container.FinanceHandler.GetCategories, r.container.AuthMiddleware.GrantPermission(constants.MANAGE_FINANCE_CATEGORY))
	g.POST("/categories", r.container.FinanceHandler.CreateCategory, r.container.AuthMiddleware.GrantPermission(constants.MANAGE_FINANCE_CATEGORY))
	g.PUT("/categories/:id", r.container.FinanceHandler.UpdateCategory, r.container.AuthMiddleware.GrantPermission(constants.MANAGE_FINANCE_CATEGORY))
	g.DELETE("/categories/:id", r.container.FinanceHandler.DeleteCategory, r.container.AuthMiddleware.GrantPermission(constants.MANAGE_FINANCE_CATEGORY))

	g.GET("/transactions", r.container.FinanceHandler.GetAllTransactions, r.container.AuthMiddleware.GrantPermission(constants.VIEW_FINANCE))
	g.GET("/transactions/export", r.container.FinanceHandler.ExportTransactions, r.container.AuthMiddleware.GrantPermission(constants.EXPORT_FINANCE))
	g.GET("/transactions/:id", r.container.FinanceHandler.GetTransactionDetail, r.container.AuthMiddleware.GrantPermission(constants.VIEW_FINANCE))
	g.POST("/transactions", r.container.FinanceHandler.CreateTransaction, r.container.AuthMiddleware.GrantPermission(constants.CREATE_FINANCE))
	g.PUT("/transactions/:id/action", r.container.FinanceHandler.ProcessAction, r.container.AuthMiddleware.GrantPermission(constants.APPROVAL_FINANCE))
}
