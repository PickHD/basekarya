# Subscription Redis Caching Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add Redis caching to subscription feature checks to eliminate per-request DB queries, replace FlushDB with targeted invalidation, add a subscription expiry cron job, and add an admin cache refresh endpoint.

**Architecture:** A new `PlanCacheService` in the subscription module handles cache-aside (Redis hit→return, miss→DB→cache). The middleware delegates to it via a `ModuleAccessProvider` interface. Targeted `Del` calls replace `FlushDB` on plan approval and status change. A daily cron marks expired subscriptions and invalidates their cache.

**Tech Stack:** Go, Echo v4, GORM, Redis (go-redis/v9), robfig/cron v3, miniredis (for tests)

---

## File Map

| File | Action | Responsibility |
|------|--------|---------------|
| `internal/modules/subscription/plancache.go` | **Create** | `PlanCacheService` — `HasAccess`, `CheckEmployeeLimit`, cache-aside logic |
| `internal/modules/subscription/plancache_test.go` | **Create** | Unit tests with miniredis |
| `internal/modules/subscription/contract.go` | **Modify** | Add `Del` to `CacheProvider`, remove `FlushDB` |
| `internal/modules/subscription/service.go` | **Modify** | Targeted invalidation in `ReviewRequest`+`UpdateCompanyStatus`, remove `FlushDB` |
| `internal/modules/subscription/service_test.go` | **Modify** | Update mock expectations: `FlushDB`→`Del` |
| `internal/modules/subscription/mocks_test.go` | **Modify** | Replace `FlushDB` with `Del` on `mockCache` |
| `internal/modules/subscription/repository.go` | **Modify** | Add `FindExpiredCompanies` |
| `internal/modules/subscription/scheduler.go` | **Create** | Expiry cron job |
| `internal/modules/subscription/handler.go` | **Modify** | Add `RefreshCompanyCache` handler |
| `internal/modules/subscription/dto.go` | **Modify** | Add `RefreshCacheResponse` |
| `internal/middleware/subscription.go` | **Modify** | Inject `ModuleAccessProvider` interface, delegate to it |
| `internal/middleware/subscription_test.go` | **Modify** | Update tests for new interface |
| `pkg/constants/cache_key.go` | **Modify** | Add `SUBSCRIPTION_FEATURES_CACHE_KEY` |
| `internal/routes/subscription.go` | **Modify** | Wire admin refresh endpoint |
| `internal/bootstrap/container.go` | **Modify** | Create `PlanCacheService`, wire middleware + scheduler |
| `cmd/api/main.go` | **Modify** | Start subscription expiry scheduler |

---

### Task 1: Add cache key constant

**Files:**
- Modify: `backend/pkg/constants/cache_key.go`

- [ ] **Step 1: Add the cache key constant**

```go
package constants

const (
	PERMISSION_CACHE_KEY           = "permission:all"
	ROLE_CACHE_KEY                 = "role:all"
	ROLE_PERMISSION_CACHE_KEY      = "role:permission:%d"
	USER_CACHE_KEY                 = "user:%d"
	DEPARTMEN_CACHE_KEY            = "department:all"
	SHIFT_CACHE_KEY                = "shift:all"
	LEAVE_TYPE_CACHE_KEY           = "leave_type:all"
	COMPANY_PROFILE_CACHE_KEY      = "company:profile:%d"
	SUBSCRIPTION_FEATURES_CACHE_KEY = "subscription:features:%d"
)
```

- [ ] **Step 2: Verify compilation**

Run: `go build ./pkg/constants/...`
Expected: success

- [ ] **Step 3: Commit**

```bash
git add backend/pkg/constants/cache_key.go
git commit -m "feat: add SUBSCRIPTION_FEATURES_CACHE_KEY constant"
```

---

### Task 2: Create PlanCacheService

**Files:**
- Create: `backend/internal/modules/subscription/plancache.go`
- Create: `backend/internal/modules/subscription/plancache_test.go`

- [ ] **Step 1: Write the PlanCacheService**

```go
package subscription

import (
	"basekarya-backend/internal/infrastructure"
	"basekarya-backend/pkg/constants"
	"basekarya-backend/pkg/logger"
	"basekarya-backend/pkg/utils"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type PlanCacheService struct {
	db    *gorm.DB
	redis *infrastructure.RedisClientProvider
}

type planFeatures struct {
	Modules []string `json:"modules"`
}

func NewPlanCacheService(db *gorm.DB, redis *infrastructure.RedisClientProvider) *PlanCacheService {
	return &PlanCacheService{db: db, redis: redis}
}

func (s *PlanCacheService) HasAccess(ctx context.Context, companyID uint, module string) (bool, error) {
	key := fmt.Sprintf(constants.SUBSCRIPTION_FEATURES_CACHE_KEY, companyID)

	featuresJSON, err := s.redis.Get(ctx, key)
	if err == nil && featuresJSON != "" {
		var features planFeatures
		if err := json.Unmarshal([]byte(featuresJSON), &features); err != nil {
			return false, fmt.Errorf("failed to parse cached features: %w", err)
		}
		return moduleInFeatures(features.Modules, module), nil
	}

	var dbJSON string
	err = s.db.Table("subscription_plans").
		Select("subscription_plans.features").
		Joins("JOIN companies ON companies.subscription_plan_id = subscription_plans.id").
		Where("companies.id = ?", companyID).
		Scan(&dbJSON).Error
	if err != nil || dbJSON == "" {
		return false, fmt.Errorf("subscription plan not found")
	}

	if setErr := s.redis.Set(ctx, key, dbJSON, 24*time.Hour); setErr != nil {
		logger.Errorw("PlanCacheService: failed to cache features", "key", key, "err", setErr)
	}

	var features planFeatures
	if err := json.Unmarshal([]byte(dbJSON), &features); err != nil {
		return false, fmt.Errorf("failed to parse features: %w", err)
	}
	return moduleInFeatures(features.Modules, module), nil
}

func (s *PlanCacheService) CheckEmployeeLimit(ctx context.Context) (bool, error) {
	companyID := utils.GetCompanyIDFromCtx(ctx)
	if companyID == 0 {
		return true, nil
	}

	var maxEmployees int
	err := s.db.Table("subscription_plans").
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
	s.db.Table("users").
		Joins("JOIN roles ON roles.id = users.role_id").
		Where("users.company_id = ? AND roles.name = ? AND users.is_active = ?", companyID, "EMPLOYEE", true).
		Count(&count)

	return count < int64(maxEmployees), nil
}

func (s *PlanCacheService) Del(ctx context.Context, key string) error {
	return s.redis.Del(ctx, key)
}

func (s *PlanCacheService) Invalidate(ctx context.Context, companyID uint) {
	key := fmt.Sprintf(constants.SUBSCRIPTION_FEATURES_CACHE_KEY, companyID)
	if err := s.redis.Del(ctx, key); err != nil {
		logger.Errorw("PlanCacheService: failed to delete features cache", "key", key, "err", err)
	}
	profileKey := fmt.Sprintf(constants.COMPANY_PROFILE_CACHE_KEY, companyID)
	if err := s.redis.Del(ctx, profileKey); err != nil {
		logger.Errorw("PlanCacheService: failed to delete company profile cache", "key", profileKey, "err", err)
	}
}

func moduleInFeatures(modules []string, target string) bool {
	for _, m := range modules {
		if m == target {
			return true
		}
	}
	return false
}
```

- [ ] **Step 2: Write the PlanCacheService tests**

```go
package subscription

import (
	"basekarya-backend/internal/infrastructure"
	"basekarya-backend/pkg/constants"
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupPlanCacheTest(t *testing.T) (*PlanCacheService, *miniredis.Miniredis, *gorm.DB) {
	t.Helper()

	mr, err := miniredis.Run()
	require.NoError(t, err)

	rdb := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	redisProvider := &infrastructure.RedisClientProvider{Client: rdb}

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	err = db.Exec(`CREATE TABLE subscription_plans (id INTEGER PRIMARY KEY, name TEXT, slug TEXT, max_employees INTEGER, price_monthly REAL, features TEXT, is_active INTEGER, created_at DATETIME, updated_at DATETIME)`).Error
	require.NoError(t, err)

	err = db.Exec(`CREATE TABLE companies (id INTEGER PRIMARY KEY, name TEXT, subscription_plan_id INTEGER, subscription_status TEXT, subscription_expires_at DATETIME, created_at DATETIME, updated_at DATETIME)`).Error
	require.NoError(t, err)

	err = db.Exec(`CREATE TABLE users (id INTEGER PRIMARY KEY, company_id INTEGER, role_id INTEGER, is_active INTEGER, created_at DATETIME, updated_at DATETIME)`).Error
	require.NoError(t, err)

	err = db.Exec(`CREATE TABLE roles (id INTEGER PRIMARY KEY, name TEXT, company_id INTEGER, created_at DATETIME, updated_at DATETIME)`).Error
	require.NoError(t, err)

	svc := &PlanCacheService{db: db, redis: redisProvider}
	return svc, mr, db
}

func TestPlanCacheService_RedisHit(t *testing.T) {
	svc, mr, db := setupPlanCacheTest(t)
	defer mr.Close()

	db.Exec(`INSERT INTO subscription_plans (id, name, slug, max_employees, price_monthly, features, is_active) VALUES (1, 'Pro', 'pro', 10, 99.00, '{"modules":["payroll","attendance"]}', 1)`)
	db.Exec(`INSERT INTO companies (id, name, subscription_plan_id, subscription_status) VALUES (1, 'TestCo', 1, 'ACTIVE')`)

	key := constants.SUBSCRIPTION_FEATURES_CACHE_KEY
	cacheKey := "subscription:features:1"
	mr.Set(cacheKey, `{"modules":["payroll","attendance"]}`)

	hasAccess, err := svc.HasAccess(context.Background(), 1, "payroll")
	require.NoError(t, err)
	assert.True(t, hasAccess)
	_ = key
}

func TestPlanCacheService_RedisMiss_PopulatesCache(t *testing.T) {
	svc, mr, db := setupPlanCacheTest(t)
	defer mr.Close()

	db.Exec(`INSERT INTO subscription_plans (id, name, slug, max_employees, price_monthly, features, is_active) VALUES (1, 'Basic', 'basic', 5, 29.00, '{"modules":["attendance"]}', 1)`)
	db.Exec(`INSERT INTO companies (id, name, subscription_plan_id, subscription_status) VALUES (1, 'TestCo', 1, 'ACTIVE')`)

	hasAccess, err := svc.HasAccess(context.Background(), 1, "attendance")
	require.NoError(t, err)
	assert.True(t, hasAccess)

	cached, err := mr.Get("subscription:features:1")
	require.NoError(t, err)
	assert.Equal(t, `{"modules":["attendance"]}`, cached)

	ttl := mr.TTL("subscription:features:1")
	assert.InDelta(t, 24*time.Hour.Seconds(), ttl.Seconds(), 5)
}

func TestPlanCacheService_ModuleFound(t *testing.T) {
	svc, mr, db := setupPlanCacheTest(t)
	defer mr.Close()

	db.Exec(`INSERT INTO subscription_plans (id, name, slug, max_employees, price_monthly, features, is_active) VALUES (1, 'Pro', 'pro', 10, 99.00, '{"modules":["payroll","attendance","recruitment"]}', 1)`)
	db.Exec(`INSERT INTO companies (id, name, subscription_plan_id, subscription_status) VALUES (1, 'TestCo', 1, 'ACTIVE')`)

	hasAccess, err := svc.HasAccess(context.Background(), 1, "recruitment")
	require.NoError(t, err)
	assert.True(t, hasAccess)
}

func TestPlanCacheService_ModuleNotFound(t *testing.T) {
	svc, mr, db := setupPlanCacheTest(t)
	defer mr.Close()

	db.Exec(`INSERT INTO subscription_plans (id, name, slug, max_employees, price_monthly, features, is_active) VALUES (1, 'Basic', 'basic', 5, 29.00, '{"modules":["attendance"]}', 1)`)
	db.Exec(`INSERT INTO companies (id, name, subscription_plan_id, subscription_status) VALUES (1, 'TestCo', 1, 'ACTIVE')`)

	hasAccess, err := svc.HasAccess(context.Background(), 1, "payroll")
	require.NoError(t, err)
	assert.False(t, hasAccess)
}

func TestPlanCacheService_NoPlan(t *testing.T) {
	svc, mr, db := setupPlanCacheTest(t)
	defer mr.Close()

	db.Exec(`INSERT INTO companies (id, name, subscription_status) VALUES (1, 'TestCo', 'ACTIVE')`)

	hasAccess, err := svc.HasAccess(context.Background(), 1, "payroll")
	require.Error(t, err)
	assert.False(t, hasAccess)
	assert.Contains(t, err.Error(), "not found")
}

func TestPlanCacheService_Invalidate(t *testing.T) {
	svc, mr, _ := setupPlanCacheTest(t)
	defer mr.Close()

	mr.Set("subscription:features:1", `{"modules":["payroll"]}`)
	mr.Set("company:profile:1", `{"name":"TestCo"}`)

	svc.Invalidate(context.Background(), 1)

	_, err := mr.Get("subscription:features:1")
	assert.Equal(t, miniredis.ErrKeyNotFound, err)

	_, err = mr.Get("company:profile:1")
	assert.Equal(t, miniredis.ErrKeyNotFound, err)
}

func TestPlanCacheService_CheckEmployeeLimit(t *testing.T) {
	svc, mr, db := setupPlanCacheTest(t)
	defer mr.Close()

	db.Exec(`INSERT INTO subscription_plans (id, name, slug, max_employees, price_monthly, features, is_active) VALUES (1, 'Basic', 'basic', 3, 29.00, '{"modules":["attendance"]}', 1)`)
	db.Exec(`INSERT INTO companies (id, name, subscription_plan_id, subscription_status) VALUES (1, 'TestCo', 1, 'ACTIVE')`)
	db.Exec(`INSERT INTO roles (id, name, company_id) VALUES (1, 'EMPLOYEE', 1)`)
	db.Exec(`INSERT INTO users (id, company_id, role_id, is_active) VALUES (1, 1, 1, 1), (2, 1, 1, 1), (3, 1, 1, 1)`)

	ctx := context.WithValue(context.Background(), constants.CompanyIDContextKey, uint(1))

	allowed, err := svc.CheckEmployeeLimit(ctx)
	require.NoError(t, err)
	assert.False(t, allowed)
}
```

- [ ] **Step 3: Run the tests to verify they pass**

Run: `go test ./internal/modules/subscription/ -run TestPlanCacheService -v -count=1`
Expected: PASS

- [ ] **Step 4: Commit**

```bash
git add backend/internal/modules/subscription/plancache.go backend/internal/modules/subscription/plancache_test.go
git commit -m "feat: add PlanCacheService with Redis caching for subscription features"
```

---

### Task 3: Update CacheProvider interface and mock

**Files:**
- Modify: `backend/internal/modules/subscription/contract.go`
- Modify: `backend/internal/modules/subscription/mocks_test.go`

- [ ] **Step 1: Update contract.go — add Del, remove FlushDB**

```go
package subscription

import (
	"context"
)

type RoleProvider interface {
	FindPermissionIDsByGroupNames(ctx context.Context, groupNames []string) ([]uint, error)
	FindRoleIDsByCompanyID(ctx context.Context, companyID uint) ([]uint, error)
	AssignPermissions(ctx context.Context, roleID uint, permissionIDs []uint, companyID uint) error
}

type CacheProvider interface {
	Del(ctx context.Context, key string) error
}

type UserProvider interface {
	ForceResetPasswordByCompanyID(ctx context.Context, companyID uint) error
}
```

- [ ] **Step 2: Update mocks_test.go — replace FlushDB with Del on mockCache**

Replace:
```go
type mockCache struct{ mock.Mock }

func (m *mockCache) FlushDB(ctx context.Context) error {
	return m.Called(ctx).Error(0)
}
```

With:
```go
type mockCache struct{ mock.Mock }

func (m *mockCache) Del(ctx context.Context, key string) error {
	return m.Called(ctx, key).Error(0)
}
```

- [ ] **Step 3: Verify compilation**

Run: `go build ./internal/modules/subscription/...`
Expected: compilation errors in service.go and service_test.go (expected — those still reference FlushDB, fixed in next task)

- [ ] **Step 4: Commit**

```bash
git add backend/internal/modules/subscription/contract.go backend/internal/modules/subscription/mocks_test.go
git commit -m "feat: replace FlushDB with Del on CacheProvider interface"
```

---

### Task 4: Update subscription service — targeted invalidation

**Files:**
- Modify: `backend/internal/modules/subscription/service.go`
- Modify: `backend/internal/modules/subscription/service_test.go`

- [ ] **Step 1: Replace FlushDB call with planCache.Invalidate in ReviewRequest**

In `service.go`, change the `ReviewRequest` method:

Add `planCache` to the `service` struct — but wait, actually the service doesn't have a `PlanCacheService` reference. Instead, we should modify the `CacheProvider` call pattern. The `CacheProvider.Del` will be used to delete specific subscription feature cache keys.

First, change the `ReviewRequest` method's cache invalidation from:
```go
_ = s.cache.FlushDB(context.Background())
```
To:
```go
_ = s.cache.Del(context.Background(), fmt.Sprintf(constants.SUBSCRIPTION_FEATURES_CACHE_KEY, subReq.CompanyID))
_ = s.cache.Del(context.Background(), fmt.Sprintf(constants.COMPANY_PROFILE_CACHE_KEY, subReq.CompanyID))
```

Full `ReviewRequest` method (only the approval block changes):

```go
func (s *service) ReviewRequest(ctx context.Context, requestID uint, req *ReviewRequest) error {
	subReq, err := s.repo.FindRequestByID(ctx, requestID)
	if err != nil {
		return errors.New("request not found")
	}

	if subReq.Status != constants.SubReqStatusPending {
		return errors.New("request already reviewed")
	}

	reviewerID := utils.GetUserIDFromCtx(ctx)
	now := time.Now()

	subReq.Status = req.Status
	subReq.ReviewedBy = &reviewerID
	subReq.ReviewedAt = &now
	subReq.Notes = req.Notes

	if err := s.repo.UpdateRequest(ctx, subReq); err != nil {
		return fmt.Errorf("failed to update request: %w", err)
	}

	if req.Status == constants.SubReqStatusApproved {
		comp, err := s.company.FindByID(ctx, subReq.CompanyID)
		if err != nil {
			return errors.New("company not found")
		}

		comp.SubscriptionPlanID = &subReq.RequestedPlanID
		comp.SubscriptionStatus = constants.SubStatusActive
		expiresAt := time.Now().Add(30 * 24 * time.Hour)
		comp.SubscriptionExpiresAt = &expiresAt

		if err := s.company.Update(ctx, comp); err != nil {
			return fmt.Errorf("failed to update company plan: %w", err)
		}

		plan, err := s.repo.FindPlanByID(ctx, subReq.RequestedPlanID)
		if err == nil && plan != nil {
			allowedGroups := buildAllowedGroups(plan.Slug)
			permissionIDs, err := s.role.FindPermissionIDsByGroupNames(ctx, allowedGroups)
			if err == nil && len(permissionIDs) > 0 {
				roleIDs, _ := s.role.FindRoleIDsByCompanyID(ctx, subReq.CompanyID)
				for _, roleID := range roleIDs {
					_ = s.role.AssignPermissions(ctx, roleID, permissionIDs, subReq.CompanyID)
				}
			}
		}

		_ = s.cache.Del(ctx, fmt.Sprintf(constants.SUBSCRIPTION_FEATURES_CACHE_KEY, subReq.CompanyID))
		_ = s.cache.Del(ctx, fmt.Sprintf(constants.COMPANY_PROFILE_CACHE_KEY, subReq.CompanyID))

		_ = s.user.ForceResetPasswordByCompanyID(ctx, subReq.CompanyID)
	}

	return nil
}
```

Also add `"fmt"` to the imports if not already present.

- [ ] **Step 2: Add invalidation to UpdateCompanyStatus**

```go
func (s *service) UpdateCompanyStatus(ctx context.Context, companyID uint, req *UpdateCompanyStatusRequest) error {
	if err := s.repo.UpdateCompanyStatus(ctx, companyID, req.SubscriptionStatus); err != nil {
		return err
	}

	_ = s.cache.Del(ctx, fmt.Sprintf(constants.SUBSCRIPTION_FEATURES_CACHE_KEY, companyID))
	_ = s.cache.Del(ctx, fmt.Sprintf(constants.COMPANY_PROFILE_CACHE_KEY, companyID))

	return nil
}
```

- [ ] **Step 3: Update service_test.go — change FlushDB expectations to Del**

In the "approve success" test case (around line 300), replace:
```go
cache.On("FlushDB", mock.Anything).Return(nil)
```
With:
```go
cache.On("Del", mock.Anything, "subscription:features:1").Return(nil)
cache.On("Del", mock.Anything, "company:profile:1").Return(nil)
```

- [ ] **Step 4: Run service tests to verify they pass**

Run: `go test ./internal/modules/subscription/ -run TestService -v -count=1`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add backend/internal/modules/subscription/service.go backend/internal/modules/subscription/service_test.go
git commit -m "feat: replace FlushDB with targeted Del invalidation in subscription service"
```

---

### Task 5: Update middleware to use PlanCacheService

**Files:**
- Modify: `backend/internal/middleware/subscription.go`
- Modify: `backend/internal/middleware/subscription_test.go`

- [ ] **Step 1: Rewrite middleware to inject ModuleAccessProvider**

```go
package middleware

import (
	"basekarya-backend/pkg/response"
	"basekarya-backend/pkg/utils"
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
)

type ModuleAccessProvider interface {
	HasAccess(ctx context.Context, companyID uint, module string) (bool, error)
}

type SubscriptionProvider interface {
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
```

Important: The `ModuleAccessProvider` interface must include BOTH `HasAccess` and `CheckEmployeeLimit` since the middleware delegates both. Update the interface:

```go
type ModuleAccessProvider interface {
	HasAccess(ctx context.Context, companyID uint, module string) (bool, error)
	CheckEmployeeLimit(ctx context.Context) (bool, error)
}
```

And remove the `SubscriptionProvider` interface since `PlanCacheService` now implements both methods directly.

- [ ] **Step 2: Update middleware tests**

Rewrite `subscription_test.go` to use a mock PlanCacheService instead of a real DB:

```go
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
	hasAccess        func(ctx context.Context, companyID uint, module string) (bool, error)
	checkEmpLimit    func(ctx context.Context) (bool, error)
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
```

- [ ] **Step 3: Run middleware tests**

Run: `go test ./internal/middleware/ -run TestSubscriptionMiddleware -v -count=1`
Expected: PASS

- [ ] **Step 4: Commit**

```bash
git add backend/internal/middleware/subscription.go backend/internal/middleware/subscription_test.go
git commit -m "refactor: inject ModuleAccessProvider into SubscriptionMiddleware"
```

---

### Task 6: Add FindExpiredCompanies to repository

**Files:**
- Modify: `backend/internal/modules/subscription/repository.go`

- [ ] **Step 1: Add interface method and implementation**

Add to the `Repository` interface:
```go
FindExpiredCompanies(ctx context.Context) ([]uint, error)
```

Add the implementation at the end of `repository.go`:
```go
func (r *repository) FindExpiredCompanies(ctx context.Context) ([]uint, error) {
	var ids []uint
	err := r.db.Table("companies").
		Select("id").
		Where("subscription_status != ?", constants.SubStatusExpired).
		Where("subscription_expires_at IS NOT NULL AND subscription_expires_at < NOW()").
		Pluck("id", &ids).Error
	return ids, err
}
```

- [ ] **Step 2: Verify compilation**

Run: `go build ./internal/modules/subscription/...`
Expected: success (but the build may need mocks to be updated if `FindExpiredCompanies` is in the Repository interface)

- [ ] **Step 3: Commit**

```bash
git add backend/internal/modules/subscription/repository.go
git commit -m "feat: add FindExpiredCompanies to subscription repository"
```

---

### Task 7: Create expiry scheduler

**Files:**
- Create: `backend/internal/modules/subscription/scheduler.go`

- [ ] **Step 1: Write the scheduler**

```go
package subscription

import (
	"basekarya-backend/internal/infrastructure"
	"basekarya-backend/pkg/constants"
	"basekarya-backend/pkg/logger"
	"context"
	"fmt"
)

type Scheduler interface {
	Start()
	Stop()
}

type subscriptionScheduler struct {
	cronProvider *infrastructure.CronProvider
	repo         Repository
	cache        CacheProvider
}

func NewScheduler(cronProvider *infrastructure.CronProvider, repo Repository, cache CacheProvider) Scheduler {
	return &subscriptionScheduler{cronProvider, repo, cache}
}

func (sch *subscriptionScheduler) Start() {
	logger.Info("Subscription Expiry Scheduler Started...")

	_, err := sch.cronProvider.GetCron().AddFunc("0 0 * * *", func() {
		logger.Info("[SCHEDULER] Starting subscription expiry check...")

		ctx := context.Background()

		ids, err := sch.repo.FindExpiredCompanies(ctx)
		if err != nil {
			logger.Errorf("[SCHEDULER] Failed to find expired companies: %v", err)
			return
		}

		for _, companyID := range ids {
			if err := sch.repo.UpdateCompanyStatus(ctx, companyID, constants.SubStatusExpired); err != nil {
				logger.Errorf("[SCHEDULER] Failed to update status for company %d: %v", companyID, err)
				continue
			}

			_ = sch.cache.Del(ctx, fmt.Sprintf(constants.SUBSCRIPTION_FEATURES_CACHE_KEY, companyID))
			_ = sch.cache.Del(ctx, fmt.Sprintf(constants.COMPANY_PROFILE_CACHE_KEY, companyID))
		}

		if len(ids) > 0 {
			logger.Infof("[SCHEDULER] Expired %d companies", len(ids))
		} else {
			logger.Info("[SCHEDULER] No expired companies found")
		}
	})

	if err != nil {
		logger.Errorf("Failed to start subscription expiry scheduler: %v", err)
	}

	sch.cronProvider.GetCron().Start()
}

func (sch *subscriptionScheduler) Stop() {
	if sch.cronProvider != nil && sch.cronProvider.GetCron() != nil {
		sch.cronProvider.GetCron().Stop()
		logger.Info("Subscription Expiry Scheduler Stopped.")
	}
}
```

- [ ] **Step 2: Verify compilation**

Run: `go build ./internal/modules/subscription/...`
Expected: success

- [ ] **Step 3: Commit**

```bash
git add backend/internal/modules/subscription/scheduler.go
git commit -m "feat: add subscription expiry scheduler"
```

---

### Task 8: Add admin refresh cache endpoint

**Files:**
- Modify: `backend/internal/modules/subscription/handler.go`
- Modify: `backend/internal/modules/subscription/dto.go`
- Modify: `backend/internal/routes/subscription.go`

- [ ] **Step 1: Add handler method to Handler**

First, add a `CacheRefresher` interface and the handler:

```go
// Add to handler.go

type CacheRefresher interface {
	RefreshCompanyCache(ctx context.Context, companyID uint) error
}
```

Add `RefreshCompanyCache` method:

```go
func (h *Handler) RefreshCompanyCache(ctx echo.Context) error {
	id := ctx.Param("id")
	var companyID uint
	if _, err := fmt.Sscanf(id, "%d", &companyID); err != nil {
		return response.NewResponses[any](ctx, http.StatusBadRequest, "Invalid company ID", nil, err, nil)
	}

	// Check company exists — use the existing service method
	_, err := h.service.GetCompanyDetail(ctx.Request().Context(), companyID)
	if err != nil {
		return response.NewResponses[any](ctx, http.StatusNotFound, "Company not found", nil, nil, nil)
	}

	// h.service doesn't have cache invalidation yet. Add a Service method.
	return response.NewResponses[any](ctx, http.StatusOK, "Cache refreshed for company", nil, nil, nil)
}
```

Wait — the handler currently has `Service` interface. We should add a `RefreshCompanyCache` method to the Service interface. Or we can pass `*PlanCacheService` to the handler.

Better approach: Add `RefreshCompanyCache(ctx context.Context, companyID uint) error` to the `Service` interface and implement it in service:

In `service.go`, add to the `Service` interface:
```go
RefreshCompanyCache(ctx context.Context, companyID uint) error
```

And the implementation:
```go
func (s *service) RefreshCompanyCache(ctx context.Context, companyID uint) error {
	s.cache.Del(ctx, fmt.Sprintf(constants.SUBSCRIPTION_FEATURES_CACHE_KEY, companyID))
	s.cache.Del(ctx, fmt.Sprintf(constants.COMPANY_PROFILE_CACHE_KEY, companyID))
	return nil
}
```

Then in `handler.go` add:

```go
func (h *Handler) RefreshCompanyCache(ctx echo.Context) error {
	id := ctx.Param("id")
	var companyID uint
	if _, err := fmt.Sscanf(id, "%d", &companyID); err != nil {
		return response.NewResponses[any](ctx, http.StatusBadRequest, "Invalid company ID", nil, err, nil)
	}

	if err := h.service.RefreshCompanyCache(ctx.Request().Context(), companyID); err != nil {
		logger.Errorw("Failed to refresh company cache: ", err)
		return response.NewResponses[any](ctx, http.StatusInternalServerError, err.Error(), nil, err, nil)
	}

	return response.NewResponses[any](ctx, http.StatusOK, "Cache refreshed for company", nil, nil, nil)
}
```

- [ ] **Step 2: Wire the route**

In `routes/subscription.go`, add to `SetupSubscriptionAdminRoutes`:

```go
func (r *Router) SetupSubscriptionAdminRoutes(e *echo.Group) {
	g := e.Group("", middleware.RequirePlatformAdmin(r.container.AuthMiddleware))
	g.GET("/pending", r.container.SubscriptionHandler.ListPendingRequests)
	g.GET("/requests", r.container.SubscriptionHandler.ListAllRequests)
	g.PUT("/:id/review", r.container.SubscriptionHandler.ReviewRequest)
	g.GET("/companies", r.container.SubscriptionHandler.ListCompanies)
	g.GET("/companies/:id", r.container.SubscriptionHandler.GetCompanyDetail)
	g.PUT("/companies/:id/status", r.container.SubscriptionHandler.UpdateCompanyStatus)
	g.DELETE("/companies/:id/cache", r.container.SubscriptionHandler.RefreshCompanyCache)
	g.GET("/dashboard", r.container.SubscriptionHandler.GetDashboardStats)
}
```

- [ ] **Step 3: Verify compilation**

Run: `go build ./internal/modules/subscription/... ./internal/routes/...`
Expected: success

- [ ] **Step 4: Commit**

```bash
git add backend/internal/modules/subscription/handler.go backend/internal/modules/subscription/service.go backend/internal/routes/subscription.go
git commit -m "feat: add admin company cache refresh endpoint"
```

---

### Task 9: Wire container and main.go

**Files:**
- Modify: `backend/internal/bootstrap/container.go`
- Modify: `backend/cmd/api/main.go`

- [ ] **Step 1: Update container.go**

Add `SubscriptionScheduler` to the `Container` struct:
```go
SubscriptionScheduler subscription.Scheduler
```

Change line 111 from:
```go
subscriptionMW := middleware.NewSubscriptionMiddleware(db.GetDB())
```
To:
```go
planCache := subscription.NewPlanCacheService(db.GetDB(), redis)
subscriptionMW := middleware.NewSubscriptionMiddleware(planCache)
```

Change line 133 from:
```go
subscriptionSvc := subscription.NewService(subscriptionRepo, companyRepo, rbacRepo, userRepo, redis)
```
To:
```go
subscriptionSvc := subscription.NewService(subscriptionRepo, companyRepo, rbacRepo, userRepo, planCache)
```

Add scheduler creation after line 163:
```go
subscriptionScheduler := subscription.NewScheduler(cronScheduler, subscriptionRepo, planCache)
```

Add to the return `&Container{` block:
```go
SubscriptionScheduler: subscriptionScheduler,
```

Add to `Close()`:
```go
if c.SubscriptionScheduler != nil {
    c.SubscriptionScheduler.Stop()
}
```

Wait — let me reconsider. The service currently takes `CacheProvider` which is `*infrastructure.RedisClientProvider`. But in the plan, we changed `CacheProvider` to have `Del` instead of `FlushDB`. The `*infrastructure.RedisClientProvider` already has `Del(ctx, key) error`. So we can pass `planCache` as the `CacheProvider` — but we need the service to use `planCache.Del` for cache invalidation, not `redis.Del` directly.

Actually, let me look at the current wiring:
```go
subscriptionSvc := subscription.NewService(subscriptionRepo, companyRepo, rbacRepo, userRepo, redis)
```

The `redis` is `*infrastructure.RedisClientProvider` and it satisfies `CacheProvider` (which now has `Del`). But we want the service to call `planCache.Invalidate()` for proper company-scoped invalidation.

Better approach: Pass `planCache` as the CacheProvider. The `PlanCacheService` implements `Del`:

Add to `plancache.go`:
```go
func (s *PlanCacheService) Del(ctx context.Context, key string) error {
	return s.redis.Del(ctx, key)
}
```

Then in container.go:
```go
planCache := subscription.NewPlanCacheService(db.GetDB(), redis)
subscriptionSvc := subscription.NewService(subscriptionRepo, companyRepo, rbacRepo, userRepo, planCache)
```

And `PlanCacheService` satisfies `CacheProvider` since it has `Del`.

Let me update the container code properly:

```go
subscriptionSvc := subscription.NewService(subscriptionRepo, companyRepo, rbacRepo, userRepo, planCache)
```

So the service's `cache` field will be the `PlanCacheService`, and `s.cache.Del(...)` will call the `Del` method on `PlanCacheService` which delegates to Redis. Clean.

But wait — in the service, `UpdateCompanyStatus` and `ReviewRequest` call `s.cache.Del(...)`. If we pass `planCache` as the cache, and `planCache` implements `Del`, then `s.cache.Del(...)` will work and we don't need to change the service further. We already called it `s.cache.Del`. Good.

The container wiring changes for `main.go` are simpler — just start the scheduler.

- [ ] **Step 1 (revised): Update container.go**

In `container.go`, find line 111:
```go
subscriptionMW := middleware.NewSubscriptionMiddleware(db.GetDB())
```
Replace with:
```go
planCache := subscription.NewPlanCacheService(db.GetDB(), redis)
subscriptionMW := middleware.NewSubscriptionMiddleware(planCache)
```

Find line 133:
```go
subscriptionSvc := subscription.NewService(subscriptionRepo, companyRepo, rbacRepo, userRepo, redis)
```
Replace with:
```go
subscriptionSvc := subscription.NewService(subscriptionRepo, companyRepo, rbacRepo, userRepo, planCache)
```

After line 163 (`contractScheduler := contract.NewScheduler(cronScheduler, contractSvc)`), add:
```go
subscriptionScheduler := subscription.NewScheduler(cronScheduler, subscriptionRepo, planCache)
```

In the `Container` struct, add after `ContractScheduler`:
```go
SubscriptionScheduler subscription.Scheduler
```

In the `return &Container{` block, add after `ContractScheduler`:
```go
SubscriptionScheduler: subscriptionScheduler,
```

In `Close()`, add after the `ContractScheduler` stop block:
```go
if c.SubscriptionScheduler != nil {
    c.SubscriptionScheduler.Stop()
}
```

- [ ] **Step 2: Update main.go**

Add after line 75:
```go
appContainer.SubscriptionScheduler.Start()
```

- [ ] **Step 3: Verify compilation**

Run: `go build ./...`
Expected: success

- [ ] **Step 4: Run all tests**

Run: `go test ./... -count=1`
Expected: all tests pass

- [ ] **Step 5: Commit**

```bash
git add backend/internal/bootstrap/container.go backend/cmd/api/main.go
git commit -m "feat: wire PlanCacheService into container, start subscription scheduler"
```

---

### Task 10: Update mockRepo for FindExpiredCompanies

**Files:**
- Modify: `backend/internal/modules/subscription/mocks_test.go`

- [ ] **Step 1: Add mock method**

After the `GetDashboardStats` mock, add:
```go
func (m *mockRepo) FindExpiredCompanies(ctx context.Context) ([]uint, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]uint), args.Error(1)
}
```

- [ ] **Step 2: Run tests**

Run: `go test ./internal/modules/subscription/... -count=1`
Expected: PASS (all tests including PlanCache tests)

- [ ] **Step 3: Commit**

```bash
git add backend/internal/modules/subscription/mocks_test.go
git commit -m "test: add FindExpiredCompanies mock"
```

---

### Task 11: Final integration — run full test suite

- [ ] **Step 1: Run all backend tests**

Run: `go test ./... -count=1`
Expected: all tests pass

- [ ] **Step 2: Fix any compilation or test failures**

If any tests fail due to the new `FindExpiredCompanies` being in the `Repository` interface but not mocked in non-subscription tests, add the mock method where needed.

If the mockService in `mocks_test.go` also needs `RefreshCompanyCache`, add it:
```go
func (m *mockService) RefreshCompanyCache(ctx context.Context, companyID uint) error {
	return m.Called(ctx, companyID).Error(0)
}
```

- [ ] **Step 3: Commit final fixes if any**

```bash
git commit -m "fix: update mocks for new interface methods"
```
