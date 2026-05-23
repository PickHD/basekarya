package middleware

import (
	"github.com/labstack/echo/v4"
)

func SecurityHeaders() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			ctx.Response().Header().Set("X-Content-Type-Options", "nosniff")
			ctx.Response().Header().Set("X-Frame-Options", "DENY")
			ctx.Response().Header().Set("X-XSS-Protection", "0")
			ctx.Response().Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
			ctx.Response().Header().Set("Permissions-Policy", "camera=(), microphone=(), geolocation=()")
			ctx.Response().Header().Set("Cache-Control", "no-store")
			return next(ctx)
		}
	}
}
