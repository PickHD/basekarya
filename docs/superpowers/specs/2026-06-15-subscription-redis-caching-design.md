# Subscription Module Redis Caching

**Date**: 2026-06-15
**Status**: Draft

## Problem

`RequireModule` middleware in `backend/internal/middleware/subscription.go` hits the database on every request to 6 module groups (payrolls, contracts, recruitments, onboarding, finances, assets). As the number of companies grows, this becomes increasingly expensive. The middleware does not cache results, and the existing `FlushDB` call on plan approval nukes all Redis keys — too aggressive.

Additionally, `RequireModule` does not check `subscription_status` or `subscription_expires_at`, meaning expired companies can still access gated modules.

## Solution

Add Redis caching for subscription feature checks using the existing cache-aside pattern already established across the codebase. Replace `FlushDB` with targeted per-company cache invalidation. Add a subscription expiry cron job and a manual admin cache refresh endpoint.

---

## Architecture

```
Request → RequireModule (middleware) → ModuleAccessProvider (interface)
                                            ↓
                              PlanCacheService (subscription module)
                                    ↓
                            Redis (hit) → parse → return
                            Redis (miss) → DB query → cache → return
```

### Components

1. **`ModuleAccessProvider` interface** — defined in `middleware/subscription.go`:
   ```go
   type ModuleAccessProvider interface {
       HasAccess(ctx context.Context, companyID uint, module string) (bool, error)
   }
   ```

2. **`PlanCacheService`** in `internal/modules/subscription/plancache.go`:
   - Dependencies: `*gorm.DB`, `*infrastructure.RedisClientProvider`
   - Implements cache-aside: Redis hit → parse → return; Redis miss → DB query → cache set → return
   - Received by `SubscriptionMiddleware` via constructor injection

3. **`SubscriptionMiddleware`** changes:
   - Constructor accepts `ModuleAccessProvider` instead of `*gorm.DB`
   - `RequireModule` delegates the access check to `HasAccess`
   - `CheckEmployeeLimit` stays as-is (calls `planCache.CheckEmployeeLimit` instead of raw DB)

4. **`CacheProvider` interface** in `subscription/contract.go` gets a `Del(ctx, key string) error` method. No more `FlushDB`.

5. **`ExpiryScheduler`** in `subscription/expiry_scheduler.go`:
   - Cron job (daily) queries expired companies, updates status, deletes cache keys

6. **Admin endpoint**: `DELETE /api/v1/admin/companies/:id/cache` — manual cache refresh

### Bootstrap wiring (`container.go`)
```go
planCache := subscription.NewPlanCacheService(db.GetDB(), redis)
subscriptionMW := middleware.NewSubscriptionMiddleware(planCache)
```

---

## Data Flow

### Cache Key & Payload

- **Key**: `subscription:features:<companyID>`
- **Payload**: Raw `features` JSON string from `subscription_plans` (e.g., `{"modules":["payroll","attendance"]}`)
- **TTL**: 24 hours

### Read Path (every gated request)

1. `GET subscription:features:<companyID>` from Redis
2. On hit: unmarshal JSON, check if `moduleName` is in `modules[]`, return result
3. On miss: `SELECT features FROM subscription_plans JOIN companies WHERE companies.id = ?`, `SET subscription:features:<companyID>`, unmarshal, check, return

### Error Handling

- Redis unavailable: fall through to DB query (fail-open), log warning
- DB query fails: return 403 Forbidden (same as today)
- Cache set fails: log warning, continue (pessimistic — next request retries)

---

## Cache Invalidation

### Triggers

| Trigger | Location | Keys Deleted |
|---------|----------|-------------|
| Plan upgrade/downgrade approved | `subscription/service.go:ReviewRequest` | `subscription:features:<companyID>`, `company:profile:<companyID>` |
| Company subscription status changed | `subscription/service.go:UpdateCompanyStatus` | `subscription:features:<companyID>`, `company:profile:<companyID>` |
| Subscription expiry (cron) | `subscription/expiry_scheduler.go` | `subscription:features:<companyID>`, `company:profile:<companyID>` |
| Manual admin refresh | New endpoint | `subscription:features:<companyID>`, `company:profile:<companyID>` |

### Removed

- `s.cache.FlushDB(context.Background())` call in `ReviewRequest` line 167 — replaced with targeted `Del` calls.
- The `FlushDB` method on `CacheProvider` interface.

---

## Cron Job: Expiry Check

- **Schedule**: Once daily at midnight
- **Logic**:
  1. Query `companies` where `subscription_expires_at < NOW()` AND `subscription_status != 'EXPIRED'`
  2. For each: update `subscription_status` to `EXPIRED`
  3. For each: `DEL subscription:features:<companyID>`, `DEL company:profile:<companyID>`
- **Location**: `subscription/expiry_scheduler.go`, started in `cmd/api/main.go`

---

## Admin Refresh Endpoint

- **Route**: `DELETE /api/v1/admin/companies/:id/cache`
- **Auth**: Platform admin only (`RequirePlatformAdmin` middleware + `GrantPermission`)
- **Behavior**:
  - Validate company exists (404 if not)
  - `DEL subscription:features:<companyID>`, `DEL company:profile:<companyID>`
  - Return 204 No Content
- **Location**: `subscription/handler.go` (new handler method), wired in `routes/subscription.go`

---

## Cache Key Constants

Add to `pkg/constants/cache_key.go`:
```go
SUBSCRIPTION_FEATURES_CACHE_KEY = "subscription:features:%d"
```

---

## Files Changed / Created

| File | Action |
|------|--------|
| `backend/internal/middleware/subscription.go` | Refactor: inject `ModuleAccessProvider` interface, thin out `RequireModule`, move `CheckEmployeeLimit` to `PlanCacheService` |
| `backend/internal/middleware/subscription_test.go` | Update tests for new interface |
| `backend/internal/modules/subscription/plancache.go` | **New**: `PlanCacheService` struct with `HasAccess`, `CheckEmployeeLimit` |
| `backend/internal/modules/subscription/plancache_test.go` | **New**: unit tests with miniredis |
| `backend/internal/modules/subscription/contract.go` | Add `Del` to `CacheProvider`, remove `FlushDB` |
| `backend/internal/modules/subscription/service.go` | `ReviewRequest`, `UpdateCompanyStatus`: add targeted invalidation, remove `FlushDB` |
| `backend/internal/modules/subscription/expiry_scheduler.go` | **New**: cron job for expired subscriptions |
| `backend/internal/modules/subscription/handler.go` | Add admin cache refresh handler |
| `backend/internal/modules/subscription/dto.go` | Add response types if needed |
| `backend/internal/routes/subscription.go` | Wire admin cache refresh endpoint |
| `backend/internal/bootstrap/container.go` | Create `PlanCacheService`, wire into middleware |
| `backend/pkg/constants/cache_key.go` | Add `SUBSCRIPTION_FEATURES_CACHE_KEY` |
| `backend/cmd/api/main.go` | Start expiry scheduler |
| `backend/internal/infrastructure/redis.go` | No changes needed (already supports `Get`, `Set`, `Del`) |

---

## Test Plan

| Test | Scope |
|------|-------|
| `TestPlanCacheService_RedisHit` | Cached features returned, no DB call |
| `TestPlanCacheService_RedisMiss` | DB queried + cache populated on miss |
| `TestPlanCacheService_ModuleFound` | Module in plan features → HasAccess true |
| `TestPlanCacheService_ModuleNotFound` | Module not in features → HasAccess false |
| `TestPlanCacheService_RedisDown_FallsThroughDB` | Redis unavailable → DB fallback succeeds |
| `TestRequireModule_CachedAccess` | Middleware delegates to PlanCacheService |
| `TestInvalidation_PlanApproved` | Cache + company profile deleted on approval |
| `TestInvalidation_StatusChange` | Cache + company profile deleted on UpdateCompanyStatus |
| `TestCronExpiry_UpdatesExpiredCompanies` | Expired companies get status=EXPIRED + cache deleted |
| `TestCronExpiry_SkipsAlreadyExpired` | Already EXPIRED companies are skipped |
| `TestAdminRefreshCache_Success` | 204, cache keys deleted, requires platform admin |
| `TestAdminRefreshCache_NotFound` | 404 if company doesn't exist |
