package routes

import "github.com/labstack/echo/v4"

func (r *Router) SetupAuthRoutes(e *echo.Group) {
	e.POST("/login", r.container.AuthHandler.Login, r.container.RateLimiterMiddleware.Init())
}
