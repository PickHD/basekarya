package routes

import (
	"basekarya-backend/internal/middleware"
	"basekarya-backend/pkg/constants"

	"github.com/labstack/echo/v4"
)

func (r *Router) SetupSubscriptionRoutes(e *echo.Group) {
	e.POST("/upgrade", r.container.SubscriptionHandler.RequestUpgrade, r.container.AuthMiddleware.GrantPermission(constants.VIEW_COMPANY))
}

func (r *Router) SetupSubscriptionAdminRoutes(e *echo.Group) {
	g := e.Group("", middleware.RequirePlatformAdmin(r.container.AuthMiddleware))
	g.GET("/pending", r.container.SubscriptionHandler.ListPendingRequests)
	g.GET("/requests", r.container.SubscriptionHandler.ListAllRequests)
	g.PUT("/:id/review", r.container.SubscriptionHandler.ReviewRequest)
	g.GET("/companies", r.container.SubscriptionHandler.ListCompanies)
	g.GET("/companies/:id", r.container.SubscriptionHandler.GetCompanyDetail)
	g.PUT("/companies/:id/status", r.container.SubscriptionHandler.UpdateCompanyStatus)
	g.GET("/dashboard", r.container.SubscriptionHandler.GetDashboardStats)
}
