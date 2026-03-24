package routes

import "github.com/labstack/echo/v4"

func (r *Router) SetupUserRoutes(e *echo.Group) {
	e.GET("/me", r.container.UserHandler.GetProfile)
	e.PUT("/profile", r.container.UserHandler.UpdateProfile)
	e.PUT("/change-password", r.container.UserHandler.ChangePassword)
}
