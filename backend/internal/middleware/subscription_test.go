package middleware

import (
	"context"
	"net/http"
	"testing"

	"basekarya-backend/internal/infrastructure"
	"basekarya-backend/internal/modules/company"
	"basekarya-backend/internal/modules/subscription"
	"basekarya-backend/internal/testutil"
	"basekarya-backend/pkg/constants"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newSubscriptionTestDB(t *testing.T) *testutil.TestDB {
	t.Helper()
	tdb := testutil.NewTestDB(&subscription.SubscriptionPlan{}, &company.Company{})
	return tdb
}

func TestSubscriptionMiddleware_RequireModule_PlatformAdminPasses(t *testing.T) {
	tdb := newSubscriptionTestDB(t)
	defer tdb.Close()

	mw := NewSubscriptionMiddleware(tdb.DB)

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
	tdb := newSubscriptionTestDB(t)
	defer tdb.Close()

	tdb.DB.Exec(`INSERT INTO subscription_plans (id, name, slug, max_employees, price_monthly, features, is_active, created_at, updated_at) VALUES (1, 'Pro', 'pro', 10, 99.00, '{"modules":["payroll","attendance"]}', 1, datetime('now'), datetime('now'))`)
	tdb.DB.Exec(`INSERT INTO companies (id, name, subscription_plan_id, subscription_status, created_at, updated_at) VALUES (1, 'TestCo', 1, 'ACTIVE', datetime('now'), datetime('now'))`)

	mw := NewSubscriptionMiddleware(tdb.DB)

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
	tdb := newSubscriptionTestDB(t)
	defer tdb.Close()

	tdb.DB.Exec(`INSERT INTO subscription_plans (id, name, slug, max_employees, price_monthly, features, is_active, created_at, updated_at) VALUES (1, 'Basic', 'basic', 10, 29.00, '{"modules":["attendance"]}', 1, datetime('now'), datetime('now'))`)
	tdb.DB.Exec(`INSERT INTO companies (id, name, subscription_plan_id, subscription_status, created_at, updated_at) VALUES (1, 'TestCo', 1, 'ACTIVE', datetime('now'), datetime('now'))`)

	mw := NewSubscriptionMiddleware(tdb.DB)

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
	tdb := newSubscriptionTestDB(t)
	defer tdb.Close()

	tdb.DB.Exec(`INSERT INTO companies (id, name, subscription_status, created_at, updated_at) VALUES (1, 'TestCo', 'ACTIVE', datetime('now'), datetime('now'))`)

	mw := NewSubscriptionMiddleware(tdb.DB)

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
