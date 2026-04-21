package routes

import (
	"basekarya-backend/pkg/constants"

	"github.com/labstack/echo/v4"
)

func (r *Router) SetupOnboardingRoutes(e *echo.Group) {
	// Templates (require MANAGE_ONBOARDING_TEMPLATE)
	e.POST("/templates", r.container.OnboardingHandler.CreateTemplate, r.container.AuthMiddleware.GrantPermission(constants.MANAGE_ONBOARDING_TEMPLATE))
	e.GET("/templates", r.container.OnboardingHandler.GetTemplates, r.container.AuthMiddleware.GrantPermission(constants.MANAGE_ONBOARDING_TEMPLATE))
	e.PUT("/templates/:id", r.container.OnboardingHandler.UpdateTemplate, r.container.AuthMiddleware.GrantPermission(constants.MANAGE_ONBOARDING_TEMPLATE))
	e.DELETE("/templates/:id", r.container.OnboardingHandler.DeleteTemplate, r.container.AuthMiddleware.GrantPermission(constants.MANAGE_ONBOARDING_TEMPLATE))

	// Workflows (require VIEW_ONBOARDING / MANAGE_ONBOARDING_TEMPLATE for creation)
	e.POST("/workflows", r.container.OnboardingHandler.CreateWorkflow, r.container.AuthMiddleware.GrantPermission(constants.MANAGE_ONBOARDING_TEMPLATE))
	e.GET("/workflows", r.container.OnboardingHandler.GetWorkflows, r.container.AuthMiddleware.GrantPermission(constants.VIEW_ONBOARDING))
	e.GET("/workflows/:id", r.container.OnboardingHandler.GetWorkflowDetail, r.container.AuthMiddleware.GrantPermission(constants.VIEW_ONBOARDING))

	// Tasks
	e.PUT("/tasks/:id/complete", r.container.OnboardingHandler.CompleteTask, r.container.AuthMiddleware.GrantPermission(constants.UPDATE_ONBOARDING_TASK))
}
