package routes

import (
	"hris-backend/internal/bootstrap"
	"hris-backend/pkg/constants"
	"hris-backend/pkg/utils"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Router struct {
	container *bootstrap.Container
	app       *echo.Echo
}

func newRouter(container *bootstrap.Container) *Router {
	app := echo.New()

	return &Router{
		container: container,
		app:       app,
	}
}

func (r *Router) setupMiddleware() {
	r.app.Use(middleware.Recover())
	r.app.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:5173", "http://127.0.0.1:5173", "http://localhost:8080", "http://127.0.0.1:8080"},
		AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	}))
	r.app.Use(middleware.RequestID())
	r.app.Validator = utils.NewValidator()
}

func (r *Router) setupRoutes() {
	// public
	r.app.GET("/health", r.container.HealthCheckHandler.HealthCheck, r.container.RateLimiterMiddleware.Init())

	api := r.app.Group("/api/v1")
	api.POST("/auth/login", r.container.AuthHandler.Login, r.container.RateLimiterMiddleware.Init())

	// protected global
	protected := api.Group("", r.container.AuthMiddleware.VerifyToken)

	// admin or employee can access
	userOnly := protected.Group("",
		r.container.AuthMiddleware.GrantRole(string(constants.UserRoleEmployee),
			string(constants.UserRoleSuperadmin)))
	{
		userOnly.GET("/users/me", r.container.UserHandler.GetProfile)
		userOnly.PUT("/users/profile", r.container.UserHandler.UpdateProfile)
		userOnly.PUT("/users/change-password", r.container.UserHandler.ChangePassword)

		userOnly.POST("/attendances/clock", r.container.AttendanceHandler.Clock)
		userOnly.GET("/attendances/today", r.container.AttendanceHandler.GetTodayStatus)
		userOnly.GET("/attendances/history", r.container.AttendanceHandler.GetHistory)

		userOnly.GET("/reimbursements", r.container.ReimbursementHandler.GetAll)
		userOnly.POST("/reimbursements", r.container.ReimbursementHandler.Create)

		userOnly.GET("/reimbursements/:id", r.container.ReimbursementHandler.GetDetail)
		userOnly.PUT("/reimbursements/:id/action", r.container.ReimbursementHandler.ProcessAction)
	}

	// only admin can access
	adminOnly := protected.Group("/admin",
		r.container.AuthMiddleware.GrantRole(string(constants.UserRoleSuperadmin)))
	{
		adminOnly.GET("/employees", r.container.UserHandler.GetAllEmployees)
		adminOnly.POST("/employees", r.container.UserHandler.CreateEmployee)
		adminOnly.PUT("/employees/:id", r.container.UserHandler.UpdateEmployee)
		adminOnly.DELETE("/employees/:id", r.container.UserHandler.DeleteEmployee)

		adminOnly.GET("/attendances/recap", r.container.AttendanceHandler.GetAllAttendanceRecap)
		adminOnly.GET("/attendances/export", r.container.AttendanceHandler.ExportAttendance)

		adminOnly.GET("/departments", r.container.MasterHandler.GetDepartments)
		adminOnly.GET("/shifts", r.container.MasterHandler.GetShifts)

		adminOnly.GET("/dashboard/stats", r.container.AttendanceHandler.GetDashboardStats)

		adminOnly.GET("/payrolls", r.container.PayrollHandler.GetList)
		adminOnly.POST("/payrolls/generate", r.container.PayrollHandler.Generate)
		adminOnly.GET("/payrolls/:id", r.container.PayrollHandler.GetDetail)
		adminOnly.GET("/payrolls/:id/download", r.container.PayrollHandler.DownloadPayslipPDF)
		adminOnly.PUT("/payrolls/:id/status", r.container.PayrollHandler.MarkAsPaid)
	}
}

func ServeHTTP(container *bootstrap.Container) *echo.Echo {
	router := newRouter(container)
	router.setupMiddleware()
	router.setupRoutes()

	return router.app
}
