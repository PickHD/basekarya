package routes

import (
	"basekarya-backend/pkg/constants"

	"github.com/labstack/echo/v4"
)

func (r *Router) SetupLoanRoutes(e *echo.Group) {
	e.GET("", r.container.LoanHandler.GetAll, r.container.AuthMiddleware.GrantAnyPermission(constants.VIEW_LOAN, constants.VIEW_SELF_LOAN))
	e.GET("/:id", r.container.LoanHandler.GetDetail, r.container.AuthMiddleware.GrantAnyPermission(constants.VIEW_LOAN, constants.VIEW_SELF_LOAN))
	e.POST("", r.container.LoanHandler.Create, r.container.AuthMiddleware.GrantPermission(constants.CREATE_LOAN))
	e.PUT("/:id/action", r.container.LoanHandler.ProcessAction, r.container.AuthMiddleware.GrantPermission(constants.APPROVAL_LOAN))
	e.GET("/export", r.container.LoanHandler.Export, r.container.AuthMiddleware.GrantPermission(constants.EXPORT_LOAN))
}
