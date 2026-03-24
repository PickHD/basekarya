package routes

import (
	"basekarya-backend/pkg/constants"

	"github.com/labstack/echo/v4"
)

func (r *Router) SetupReimbursementRoutes(e *echo.Group) {
	e.GET("", r.container.ReimbursementHandler.GetAll, r.container.AuthMiddleware.GrantAnyPermission(constants.VIEW_REIMBURSEMENT, constants.VIEW_SELF_REIMBURSEMENT))
	e.GET("/export", r.container.ReimbursementHandler.Export, r.container.AuthMiddleware.GrantPermission(constants.EXPORT_REIMBURSEMENT))
	e.POST("", r.container.ReimbursementHandler.Create, r.container.AuthMiddleware.GrantPermission(constants.CREATE_REIMBURSEMENT))
	e.GET("/:id", r.container.ReimbursementHandler.GetDetail, r.container.AuthMiddleware.GrantAnyPermission(constants.VIEW_REIMBURSEMENT, constants.VIEW_SELF_REIMBURSEMENT))
	e.PUT("/:id/action", r.container.ReimbursementHandler.ProcessAction, r.container.AuthMiddleware.GrantPermission(constants.APPROVAL_REIMBURSEMENT))
}
