package utils

import (
	"context"
	"net/http/httptest"
	"testing"

	"basekarya-backend/internal/infrastructure"
	"basekarya-backend/pkg/constants"

	"github.com/labstack/echo/v4"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func TestGetUserContext_Success(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	claims := &infrastructure.MyClaims{
		UserID:    1,
		CompanyID: 10,
	}
	ctx.Set("user", claims)

	result, err := GetUserContext(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.UserID != 1 {
		t.Errorf("expected UserID 1, got %d", result.UserID)
	}
	if result.CompanyID != 10 {
		t.Errorf("expected CompanyID 10, got %d", result.CompanyID)
	}
}

func TestGetUserContext_NoUser(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	_, err := GetUserContext(ctx)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestGetDBFromContext_WithTx(t *testing.T) {
	db, _ := gorm.Open(nil, &gorm.Config{})
	ctx := context.WithValue(context.Background(), constants.TxContextKey, db)

	result := GetDBFromContext(ctx, nil)
	if result != db {
		t.Error("expected tx from context")
	}
}

func TestGetDBFromContext_DefaultDB(t *testing.T) {
	defaultDB, _ := gorm.Open(nil, &gorm.Config{})
	ctx := context.Background()

	result := GetDBFromContext(ctx, defaultDB)
	if result != defaultDB {
		t.Error("expected default DB")
	}
}

func TestGetCompanyIDFromCtx_Set(t *testing.T) {
	ctx := context.WithValue(context.Background(), constants.CompanyIDContextKey, uint(5))

	if got := GetCompanyIDFromCtx(ctx); got != 5 {
		t.Errorf("expected 5, got %d", got)
	}
}

func TestGetCompanyIDFromCtx_NotSet(t *testing.T) {
	if got := GetCompanyIDFromCtx(context.Background()); got != 0 {
		t.Errorf("expected 0, got %d", got)
	}
}

func TestGetUserIDFromCtx_Set(t *testing.T) {
	ctx := context.WithValue(context.Background(), constants.UserIDContextKey, uint(42))

	if got := GetUserIDFromCtx(ctx); got != 42 {
		t.Errorf("expected 42, got %d", got)
	}
}

func TestGetUserIDFromCtx_NotSet(t *testing.T) {
	if got := GetUserIDFromCtx(context.Background()); got != 0 {
		t.Errorf("expected 0, got %d", got)
	}
}

func TestIsPlatformAdminFromCtx_True(t *testing.T) {
	ctx := context.WithValue(context.Background(), constants.IsPlatformAdminContextKey, true)

	if got := IsPlatformAdminFromCtx(ctx); !got {
		t.Errorf("expected true, got %v", got)
	}
}

func TestIsPlatformAdminFromCtx_False(t *testing.T) {
	if got := IsPlatformAdminFromCtx(context.Background()); got {
		t.Errorf("expected false, got %v", got)
	}
}

func TestDetachContext(t *testing.T) {
	origCtx := context.Background()
	origCtx = context.WithValue(origCtx, constants.CompanyIDContextKey, uint(5))
	origCtx = context.WithValue(origCtx, constants.UserIDContextKey, uint(10))
	origCtx = context.WithValue(origCtx, constants.IsPlatformAdminContextKey, true)

	detached := DetachContext(origCtx)

	if got := GetCompanyIDFromCtx(detached); got != 5 {
		t.Errorf("expected CompanyID 5, got %d", got)
	}
	if got := GetUserIDFromCtx(detached); got != 10 {
		t.Errorf("expected UserID 10, got %d", got)
	}
	if got := IsPlatformAdminFromCtx(detached); !got {
		t.Errorf("expected IsPlatformAdmin true, got %v", got)
	}

	if detached.Value(constants.TxContextKey) != nil {
		t.Error("expected tx reference to be removed")
	}
}

type tenantTestModel struct {
	gorm.Model
	Name      string
	CompanyID uint
}

func (tenantTestModel) TableName() string { return "tenant_tests" }

func newTestDB(models ...interface{}) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		panic("failed to connect test database: " + err.Error())
	}
	if err := db.AutoMigrate(models...); err != nil {
		panic("failed to migrate: " + err.Error())
	}
	return db
}

func TestTenantScope_PlatformAdmin(t *testing.T) {
	db := newTestDB(&tenantTestModel{})

	ctx := context.Background()
	ctx = context.WithValue(ctx, constants.CompanyIDContextKey, uint(1))
	ctx = context.WithValue(ctx, constants.UserIDContextKey, uint(1))
	ctx = context.WithValue(ctx, constants.IsPlatformAdminContextKey, true)

	result := TenantScope(ctx, db.Session(&gorm.Session{DryRun: true}).Model(&tenantTestModel{}))
	if _, ok := result.Statement.Clauses["WHERE"]; ok {
		sql := result.Statement.SQL.String()
		t.Errorf("platform admin should skip company_id filter, got: %s", sql)
	}
}

func TestTenantScope_RegularUser(t *testing.T) {
	db := newTestDB(&tenantTestModel{})

	ctx := context.Background()
	ctx = context.WithValue(ctx, constants.CompanyIDContextKey, uint(5))
	ctx = context.WithValue(ctx, constants.UserIDContextKey, uint(10))
	ctx = context.WithValue(ctx, constants.IsPlatformAdminContextKey, false)

	result := TenantScope(ctx, db.Session(&gorm.Session{DryRun: true}).Model(&tenantTestModel{}))
	if _, ok := result.Statement.Clauses["WHERE"]; !ok {
		t.Error("expected WHERE clause for regular user")
	}
}

func containsSubstring(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
