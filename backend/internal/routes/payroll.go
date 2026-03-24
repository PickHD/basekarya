package routes

import (
	"basekarya-backend/pkg/constants"

	"github.com/labstack/echo/v4"
)

func (r *Router) SetupPayrollRoutes(e *echo.Group) {
	e.GET("", r.container.PayrollHandler.GetList, r.container.AuthMiddleware.GrantPermission(constants.VIEW_PAYROLL))
	e.POST("/generate", r.container.PayrollHandler.Generate, r.container.AuthMiddleware.GrantPermission(constants.GENERATE_PAYROLL))
	e.GET("/:id", r.container.PayrollHandler.GetDetail, r.container.AuthMiddleware.GrantPermission(constants.VIEW_PAYROLL))
	e.GET("/:id/download", r.container.PayrollHandler.DownloadPayslipPDF, r.container.AuthMiddleware.GrantPermission(constants.DOWNLOAD_PAYSLIP))
	e.PUT("/:id/status", r.container.PayrollHandler.MarkAsPaid, r.container.AuthMiddleware.GrantPermission(constants.MARK_AS_PAID))
	e.POST("/:id/send-email", r.container.PayrollHandler.BlastPayslipEmail, r.container.AuthMiddleware.GrantPermission(constants.SEND_PAYSLIP))
}