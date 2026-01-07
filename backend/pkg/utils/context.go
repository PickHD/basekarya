package utils

import (
	"errors"
	"hris-backend/internal/infrastructure"

	"github.com/labstack/echo/v4"
)

func GetUserContext(ctx echo.Context) (*infrastructure.MyClaims, error) {
	userContext := ctx.Get("user")
	if claims, ok := userContext.(*infrastructure.MyClaims); ok {
		return claims, nil
	}
	return nil, errors.New("failed to get user from context")
}
