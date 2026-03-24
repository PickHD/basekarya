package routes

import (
	"basekarya-backend/pkg/constants"

	"github.com/labstack/echo/v4"
)

func (r *Router) SetupAttendanceRoutes(e *echo.Group) {
	e.POST("/clock", r.container.AttendanceHandler.Clock, r.container.AuthMiddleware.GrantPermission(constants.CREATE_ATTENDANCE))
	e.GET("/today", r.container.AttendanceHandler.GetTodayStatus, r.container.AuthMiddleware.GrantPermission(constants.VIEW_SELF_ATTENDANCE))
	e.GET("/history", r.container.AttendanceHandler.GetHistory, r.container.AuthMiddleware.GrantPermission(constants.VIEW_SELF_ATTENDANCE))
	e.GET("/recap", r.container.AttendanceHandler.GetAllAttendanceRecap, r.container.AuthMiddleware.GrantPermission(constants.VIEW_ATTENDANCE))
	e.GET("/export", r.container.AttendanceHandler.ExportAttendance, r.container.AuthMiddleware.GrantPermission(constants.EXPORT_ATTENDANCE))
	e.GET("/dashboard/stats", r.container.AttendanceHandler.GetDashboardStats, r.container.AuthMiddleware.GrantPermission(constants.VIEW_ATTENDANCE))
}
