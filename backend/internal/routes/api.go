package routes

import (
	"basekarya-backend/internal/bootstrap"
	"basekarya-backend/pkg/utils"

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
	r.SetupAuthRoutes(api.Group("/auth"))

	// protected global
	protected := api.Group("", r.container.AuthMiddleware.VerifyToken)

	// setup websocket
	protected.GET("/ws", r.container.NotificationHandler.HandleWebSocket)

	// setup module routes
	r.SetupAttendanceRoutes(protected.Group("/attendances"))
	r.SetupCompanyRoutes(protected.Group("/companies"))
	r.SetupEmployeeRoutes(protected.Group("/employees"))
	r.SetupLeaveRoutes(protected.Group("/leaves"))
	r.SetupLoanRoutes(protected.Group("/loans"))
	r.SetupMasterRoutes(protected.Group("/masters"))
	r.SetupNotificationRoutes(protected.Group("/notifications"))
	r.SetupOvertimeRoutes(protected.Group("/overtimes"))
	r.SetupPayrollRoutes(protected.Group("/payrolls"))
	r.SetupReimbursementRoutes(protected.Group("/reimbursements"))
	r.SetupRoleRoutes(protected.Group("/roles"))
	r.SetupPermissionRoutes(protected.Group("/permissions"))
	r.SetupUserRoutes(protected.Group("/users"))
	r.SetupAnnouncementRoutes(protected.Group("/announcements"))
}

func ServeHTTP(container *bootstrap.Container) *echo.Echo {
	router := newRouter(container)
	router.setupMiddleware()
	router.setupRoutes()

	return router.app
}
