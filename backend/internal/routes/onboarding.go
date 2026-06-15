package routes

import (
	"basekarya-backend/internal/middleware"
	"basekarya-backend/pkg/constants"

	"github.com/labstack/echo/v4"
)

func (r *Router) SetupOnboardingRoutes(e *echo.Group, sub *middleware.SubscriptionMiddleware) {
	g := e.Group("", sub.RequireModule("onboarding"))

	// Workflows
	g.POST("/workflows", r.container.OnboardingHandler.CreateWorkflow, r.container.AuthMiddleware.GrantPermission(constants.VIEW_ONBOARDING))
	g.GET("/workflows", r.container.OnboardingHandler.GetWorkflows, r.container.AuthMiddleware.GrantPermission(constants.VIEW_ONBOARDING))
	g.GET("/workflows/:id", r.container.OnboardingHandler.GetWorkflowDetail, r.container.AuthMiddleware.GrantPermission(constants.VIEW_ONBOARDING))

	// Tasks
	g.PUT("/tasks/:id/complete", r.container.OnboardingHandler.CompleteTask, r.container.AuthMiddleware.GrantPermission(constants.UPDATE_ONBOARDING_TASK))
	g.PUT("/workflows/:id/tasks", r.container.OnboardingHandler.UpdateWorkflowTasks, r.container.AuthMiddleware.GrantPermission(constants.VIEW_ONBOARDING))
}
