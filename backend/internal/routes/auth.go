package routes

import "github.com/labstack/echo/v4"

func (r *Router) SetupAuthRoutes(e *echo.Group) {
	e.POST("/login", r.container.AuthHandler.Login, r.container.RateLimiterMiddleware.Init())
	e.POST("/forgot-password", r.container.AuthHandler.ForgotPassword, r.container.RateLimiterMiddleware.Init())
	e.POST("/resend-otp", r.container.AuthHandler.ResendOTP, r.container.RateLimiterMiddleware.Init())
	e.POST("/verify-otp", r.container.AuthHandler.VerifyOTP, r.container.RateLimiterMiddleware.Init())
	e.POST("/reset-password", r.container.AuthHandler.ResetPassword, r.container.RateLimiterMiddleware.Init())
}
