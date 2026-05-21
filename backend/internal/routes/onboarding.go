package routes

import (
	"basekarya-backend/internal/middleware"
	"basekarya-backend/pkg/constants"

	"github.com/labstack/echo/v4"
)

func (r *Router) SetupOnboardingRoutes(e *echo.Group, sub *middleware.SubscriptionMiddleware) {
	g := e.Group("", sub.RequireModule("onboarding"))
	// Templates (require MANAGE_ONBOARDING_TEMPLATE)
	g.POST("/templates", r.container.OnboardingHandler.CreateTemplate, r.container.AuthMiddleware.GrantPermission(constants.MANAGE_ONBOARDING_TEMPLATE))
	g.GET("/templates", r.container.OnboardingHandler.GetTemplates, r.container.AuthMiddleware.GrantPermission(constants.MANAGE_ONBOARDING_TEMPLATE))
	g.PUT("/templates/:id", r.container.OnboardingHandler.UpdateTemplate, r.container.AuthMiddleware.GrantPermission(constants.MANAGE_ONBOARDING_TEMPLATE))
	g.DELETE("/templates/:id", r.container.OnboardingHandler.DeleteTemplate, r.container.AuthMiddleware.GrantPermission(constants.MANAGE_ONBOARDING_TEMPLATE))

	// Workflows (require VIEW_ONBOARDING / MANAGE_ONBOARDING_TEMPLATE for creation)
	g.POST("/workflows", r.container.OnboardingHandler.CreateWorkflow, r.container.AuthMiddleware.GrantPermission(constants.MANAGE_ONBOARDING_TEMPLATE))
	g.GET("/workflows", r.container.OnboardingHandler.GetWorkflows, r.container.AuthMiddleware.GrantPermission(constants.VIEW_ONBOARDING))
	g.GET("/workflows/:id", r.container.OnboardingHandler.GetWorkflowDetail, r.container.AuthMiddleware.GrantPermission(constants.VIEW_ONBOARDING))

	// Tasks
	g.PUT("/tasks/:id/complete", r.container.OnboardingHandler.CompleteTask, r.container.AuthMiddleware.GrantPermission(constants.UPDATE_ONBOARDING_TASK))
}
