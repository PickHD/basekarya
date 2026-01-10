package middleware

import (
	"hris-backend/internal/infrastructure"
	"hris-backend/pkg/response"
	"hris-backend/pkg/utils"
	"net/http"
	"slices"
	"strings"

	"github.com/labstack/echo/v4"
)

type AuthMiddleware struct {
	jwtProvider *infrastructure.JwtProvider
}

func NewAuthMiddleware(jwtProvider *infrastructure.JwtProvider) *AuthMiddleware {
	return &AuthMiddleware{
		jwtProvider: jwtProvider,
	}
}

func (m *AuthMiddleware) VerifyToken(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		authHeader := ctx.Request().Header.Get("Authorization")
		if authHeader == "" {
			return response.NewResponses[any](ctx, http.StatusUnauthorized, "missing authorization header", nil, nil, nil)
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return response.NewResponses[any](ctx, http.StatusUnauthorized, "invalid authorization format", nil, nil, nil)
		}

		tokenString := parts[1]

		claims, err := m.jwtProvider.ValidateToken(tokenString)
		if err != nil {
			return response.NewResponses[any](ctx, http.StatusUnauthorized, "invalid or expired token", nil, nil, nil)
		}

		ctx.Set("user", claims)

		return next(ctx)
	}
}

func (m *AuthMiddleware) GrantRole(roles ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			userContext, err := utils.GetUserContext(ctx)
			if err != nil {
				return response.NewResponses[any](ctx, http.StatusInternalServerError, err.Error(), nil, err, nil)
			}

			if !slices.Contains(roles, userContext.Role) {
				return response.NewResponses[any](ctx, http.StatusForbidden, "You dont have access to this resource", nil, nil, nil)
			}

			return next(ctx)
		}
	}
}
