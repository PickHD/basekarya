package middleware

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"basekarya-backend/internal/infrastructure"
	"basekarya-backend/internal/testutil"
	"basekarya-backend/pkg/constants"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockPlanCache struct {
	hasAccess     func(ctx context.Context, companyID uint, module string) (bool, error)
	checkEmpLimit func(ctx context.Context) (bool, error)
}

func (m *mockPlanCache) HasAccess(ctx context.Context, companyID uint, module string) (bool, error) {
	if m.hasAccess != nil {
		return m.hasAccess(ctx, companyID, module)
	}
	return false, errors.New("not implemented")
}

func (m *mockPlanCache) CheckEmployeeLimit(ctx context.Context) (bool, error) {
	if m.checkEmpLimit != nil {
		return m.checkEmpLimit(ctx)
	}
	return true, nil
}

func TestSubscriptionMiddleware_RequireModule_PlatformAdminPasses(t *testing.T) {
	mock := &mockPlanCache{}
	mw := NewSubscriptionMiddleware(mock)

	claims := &infrastructure.MyClaims{
		UserID:          1,
		CompanyID:       1,
		IsPlatformAdmin: true,
	}

	at := testutil.NewAPITest(t, http.MethodGet, "/test", nil)
	at.WithAuthContext(claims)

	handler := mw.RequireModule("payroll")(func(ctx echo.Context) error {
		return ctx.String(http.StatusOK, "ok")
	})

	rec, err := at.Execute(handler)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestSubscriptionMiddleware_RequireModule_CompanyHasModule(t *testing.T) {
	mock := &mockPlanCache{
		hasAccess: func(ctx context.Context, companyID uint, module string) (bool, error) {
			return true, nil
		},
	}
	mw := NewSubscriptionMiddleware(mock)

	at := testutil.NewAPITest(t, http.MethodGet, "/test", nil)

	ctx := context.Background()
	ctx = context.WithValue(ctx, constants.CompanyIDContextKey, uint(1))
	ctx = context.WithValue(ctx, constants.IsPlatformAdminContextKey, false)
	ctx = context.WithValue(ctx, constants.UserIDContextKey, uint(1))
	at.Context.SetRequest(at.Context.Request().WithContext(ctx))

	handler := mw.RequireModule("payroll")(func(ctx echo.Context) error {
		return ctx.String(http.StatusOK, "ok")
	})

	rec, err := at.Execute(handler)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestSubscriptionMiddleware_RequireModule_CompanyMissingModule(t *testing.T) {
	mock := &mockPlanCache{
		hasAccess: func(ctx context.Context, companyID uint, module string) (bool, error) {
			return false, nil
		},
	}
	mw := NewSubscriptionMiddleware(mock)

	at := testutil.NewAPITest(t, http.MethodGet, "/test", nil)

	ctx := context.Background()
	ctx = context.WithValue(ctx, constants.CompanyIDContextKey, uint(1))
	ctx = context.WithValue(ctx, constants.IsPlatformAdminContextKey, false)
	ctx = context.WithValue(ctx, constants.UserIDContextKey, uint(1))
	at.Context.SetRequest(at.Context.Request().WithContext(ctx))

	handler := mw.RequireModule("payroll")(func(ctx echo.Context) error {
		return ctx.String(http.StatusOK, "ok")
	})

	rec, err := at.Execute(handler)
	require.NoError(t, err)
	assert.Equal(t, http.StatusForbidden, rec.Code)
}

func TestSubscriptionMiddleware_RequireModule_CompanyNoPlan(t *testing.T) {
	mock := &mockPlanCache{
		hasAccess: func(ctx context.Context, companyID uint, module string) (bool, error) {
			return false, errors.New("subscription plan not found")
		},
	}
	mw := NewSubscriptionMiddleware(mock)

	at := testutil.NewAPITest(t, http.MethodGet, "/test", nil)

	ctx := context.Background()
	ctx = context.WithValue(ctx, constants.CompanyIDContextKey, uint(1))
	ctx = context.WithValue(ctx, constants.IsPlatformAdminContextKey, false)
	ctx = context.WithValue(ctx, constants.UserIDContextKey, uint(1))
	at.Context.SetRequest(at.Context.Request().WithContext(ctx))

	handler := mw.RequireModule("payroll")(func(ctx echo.Context) error {
		return ctx.String(http.StatusOK, "ok")
	})

	rec, err := at.Execute(handler)
	require.NoError(t, err)
	assert.Equal(t, http.StatusForbidden, rec.Code)
}
