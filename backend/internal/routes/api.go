package routes

import (
	"basekarya-backend/internal/bootstrap"
	customMiddleware "basekarya-backend/internal/middleware"
	"basekarya-backend/pkg/utils"
	"os"
	"strings"

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
	r.app.Use(customMiddleware.SecurityHeaders())
	r.app.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: strings.Split(getEnv("CORS_ALLOWED_ORIGINS", "http://localhost:5173"), ","),
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
	api.GET("/subscriptions/plans", r.container.SubscriptionHandler.ListPlans)

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
	r.SetupPayrollRoutes(protected.Group("/payrolls"), r.container.SubscriptionMiddleware)
	r.SetupReimbursementRoutes(protected.Group("/reimbursements"))
	r.SetupRoleRoutes(protected.Group("/roles"))
	r.SetupPermissionRoutes(protected.Group("/permissions"))
	r.SetupUserRoutes(protected.Group("/users"))
	r.SetupAnnouncementRoutes(protected.Group("/announcements"))
	r.SetupContractRoutes(protected.Group("/contracts"), r.container.SubscriptionMiddleware)
	r.SetupRecruitmentRoutes(protected.Group("/recruitments"), r.container.SubscriptionMiddleware)
	r.SetupOnboardingRoutes(protected.Group("/onboarding"), r.container.SubscriptionMiddleware)
	r.SetupFinanceRoutes(protected.Group("/finances"), r.container.SubscriptionMiddleware)
	r.SetupAssetRoutes(protected.Group("/assets"), r.container.SubscriptionMiddleware)
	r.SetupSubscriptionRoutes(protected.Group("/subscriptions"))
	r.SetupSubscriptionAdminRoutes(protected.Group("/admin/subscriptions"))
}

func ServeHTTP(container *bootstrap.Container) *echo.Echo {
	router := newRouter(container)
	router.setupMiddleware()
	router.setupRoutes()

	return router.app
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
