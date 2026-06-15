package middleware

import (
	"context"
	"net/http"

	"basekarya-backend/pkg/response"
	"basekarya-backend/pkg/utils"

	"github.com/labstack/echo/v4"
)

type ModuleAccessProvider interface {
	HasAccess(ctx context.Context, companyID uint, module string) (bool, error)
	CheckEmployeeLimit(ctx context.Context) (bool, error)
}

type SubscriptionMiddleware struct {
	planCache ModuleAccessProvider
}

func NewSubscriptionMiddleware(planCache ModuleAccessProvider) *SubscriptionMiddleware {
	return &SubscriptionMiddleware{planCache: planCache}
}

func (m *SubscriptionMiddleware) RequireModule(moduleName string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			if utils.IsPlatformAdminFromCtx(ctx.Request().Context()) {
				return next(ctx)
			}

			companyID := utils.GetCompanyIDFromCtx(ctx.Request().Context())
			if companyID == 0 {
				return next(ctx)
			}

			hasAccess, err := m.planCache.HasAccess(ctx.Request().Context(), companyID, moduleName)
			if err != nil {
				return response.NewResponses[any](ctx, http.StatusForbidden, "subscription plan not found", nil, nil, nil)
			}

			if !hasAccess {
				return response.NewResponses[any](ctx, http.StatusForbidden, "Module not available in your subscription plan", nil, nil, nil)
			}

			return next(ctx)
		}
	}
}

func (m *SubscriptionMiddleware) CheckEmployeeLimit(ctx context.Context) (bool, error) {
	return m.planCache.CheckEmployeeLimit(ctx)
}
