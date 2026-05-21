package routes

import (
	"basekarya-backend/internal/middleware"
	"basekarya-backend/pkg/constants"

	"github.com/labstack/echo/v4"
)

func (r *Router) SetupPayrollRoutes(e *echo.Group, sub *middleware.SubscriptionMiddleware) {
	g := e.Group("", sub.RequireModule("payroll"))
	g.GET("", r.container.PayrollHandler.GetList, r.container.AuthMiddleware.GrantPermission(constants.VIEW_PAYROLL))
	g.POST("/generate", r.container.PayrollHandler.Generate, r.container.AuthMiddleware.GrantPermission(constants.GENERATE_PAYROLL))
	g.GET("/:id", r.container.PayrollHandler.GetDetail, r.container.AuthMiddleware.GrantPermission(constants.VIEW_PAYROLL))
	g.GET("/:id/download", r.container.PayrollHandler.DownloadPayslipPDF, r.container.AuthMiddleware.GrantPermission(constants.DOWNLOAD_PAYSLIP))
	g.PUT("/:id/status", r.container.PayrollHandler.MarkAsPaid, r.container.AuthMiddleware.GrantPermission(constants.MARK_AS_PAID))
	g.POST("/:id/send-email", r.container.PayrollHandler.BlastPayslipEmail, r.container.AuthMiddleware.GrantPermission(constants.SEND_PAYSLIP))
}