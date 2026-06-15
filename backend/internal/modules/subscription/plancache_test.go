package subscription

import (
	"basekarya-backend/internal/infrastructure"
	"basekarya-backend/pkg/constants"
	"context"
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

	mr.Set("subscription:features:1", `{"modules":["payroll","attendance"]}`)

	hasAccess, err := svc.HasAccess(context.Background(), 1, "payroll")
	require.NoError(t, err)
	assert.True(t, hasAccess)
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
