package utils

import (
	"context"
	"errors"
	"hris-backend/internal/infrastructure"
	"hris-backend/pkg/constants"

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

// Helper function for Repositories to get the correct DB (TX or Standard)
func GetDBFromContext(ctx context.Context, defaultDB *gorm.DB) *gorm.DB {
	if tx, ok := ctx.Value(constants.TxContextKey).(*gorm.DB); ok {
		return tx
	}
	return defaultDB
}
