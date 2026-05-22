package utils

import (
	"basekarya-backend/internal/infrastructure"
	"basekarya-backend/pkg/constants"
	"context"
	"errors"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func GetUserContext(ctx echo.Context) (*infrastructure.MyClaims, error) {
	userContext := ctx.Get("user")
	if claims, ok := userContext.(*infrastructure.MyClaims); ok {
		return claims, nil
	}
	return nil, errors.New("failed to get user from context")
}

func GetDBFromContext(ctx context.Context, defaultDB *gorm.DB) *gorm.DB {
	if tx, ok := ctx.Value(constants.TxContextKey).(*gorm.DB); ok {
		return tx
	}
	return defaultDB
}

func GetCompanyIDFromCtx(ctx context.Context) uint {
	if v, ok := ctx.Value(constants.CompanyIDContextKey).(uint); ok {
		return v
	}
	return 0
}

func GetUserIDFromCtx(ctx context.Context) uint {
	if v, ok := ctx.Value(constants.UserIDContextKey).(uint); ok {
		return v
	}
	return 0
}

func IsPlatformAdminFromCtx(ctx context.Context) bool {
	if v, ok := ctx.Value(constants.IsPlatformAdminContextKey).(bool); ok {
		return v
	}
	return false
}

func TenantScope(ctx context.Context, db *gorm.DB) *gorm.DB {
	if IsPlatformAdminFromCtx(ctx) {
		return db
	}
	companyID := GetCompanyIDFromCtx(ctx)
	if companyID == 0 {
		return db
	}

	tableName := db.Statement.Table
	if tableName == "" && db.Statement.Model != nil {
		if err := db.Statement.Parse(db.Statement.Model); err == nil {
			tableName = db.Statement.Table
		}
	}

	if tableName != "" {
		return db.Where(tableName+".company_id = ?", companyID)
	}
	return db.Where("company_id = ?", companyID)
}

// DetachContext returns a new context with tenant values preserved but
// without any transaction reference. Use this when spawning goroutines
// that need to perform DB operations independently of the parent transaction.
func DetachContext(ctx context.Context) context.Context {
	detached := context.Background()
	detached = context.WithValue(detached, constants.CompanyIDContextKey, GetCompanyIDFromCtx(ctx))
	detached = context.WithValue(detached, constants.IsPlatformAdminContextKey, IsPlatformAdminFromCtx(ctx))
	detached = context.WithValue(detached, constants.UserIDContextKey, GetUserIDFromCtx(ctx))
	return detached
}
