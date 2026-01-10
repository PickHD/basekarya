package routes

import (
	"hris-backend/internal/bootstrap"
	"hris-backend/pkg/constants"

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
}

func (r *Router) setupRoutes() {
	// public
	r.app.GET("/health", r.container.HealthCheckHandler.HealthCheck)

	api := r.app.Group("/api/v1")
	api.POST("/auth/login", r.container.AuthHandler.Login)

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
	}
}

func ServeHTTP(container *bootstrap.Container) *echo.Echo {
	router := newRouter(container)
	router.setupMiddleware()
	router.setupRoutes()

	return router.app
}
