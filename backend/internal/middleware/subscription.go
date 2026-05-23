package middleware

import (
	"basekarya-backend/pkg/response"
	"basekarya-backend/pkg/utils"
	"context"
	"encoding/json"
	"net/http"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type SubscriptionMiddleware struct {
	db *gorm.DB
}

func NewSubscriptionMiddleware(db *gorm.DB) *SubscriptionMiddleware {
	return &SubscriptionMiddleware{db: db}
}

type planFeatures struct {
	Modules []string `json:"modules"`
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

			var featuresJSON string
			err := m.db.Table("subscription_plans").
				Select("subscription_plans.features").
				Joins("JOIN companies ON companies.subscription_plan_id = subscription_plans.id").
				Where("companies.id = ?", companyID).
				Scan(&featuresJSON).Error

		if err != nil || featuresJSON == "" {
			return response.NewResponses[any](ctx, http.StatusForbidden, "subscription plan not found", nil, nil, nil)
		}

		var features planFeatures
		if err := json.Unmarshal([]byte(featuresJSON), &features); err != nil {
			return response.NewResponses[any](ctx, http.StatusForbidden, "failed to parse subscription features", nil, nil, nil)
		}

			for _, mod := range features.Modules {
				if mod == moduleName {
					return next(ctx)
				}
			}

			return response.NewResponses[any](ctx, http.StatusForbidden, "Module not available in your subscription plan", nil, nil, nil)
		}
	}
}

func (m *SubscriptionMiddleware) CheckEmployeeLimit(ctx context.Context) (bool, error) {
	companyID := utils.GetCompanyIDFromCtx(ctx)
	if companyID == 0 {
		return true, nil
	}

	var maxEmployees int
	err := m.db.Table("subscription_plans").
		Select("subscription_plans.max_employees").
		Joins("JOIN companies ON companies.subscription_plan_id = subscription_plans.id").
		Where("companies.id = ?", companyID).
		Scan(&maxEmployees).Error
	if err != nil {
		return true, err
	}

	if maxEmployees == 0 {
		return true, nil
	}

	var count int64
	m.db.Table("users").
		Joins("JOIN roles ON roles.id = users.role_id").
		Where("users.company_id = ? AND roles.name = ? AND users.is_active = ?", companyID, "EMPLOYEE", true).
		Count(&count)

	return count < int64(maxEmployees), nil
}
