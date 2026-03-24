package routes

import "github.com/labstack/echo/v4"

func (r *Router) SetupNotificationRoutes(e *echo.Group) {
	e.GET("", r.container.NotificationHandler.GetAll)
	e.PUT("/:id/read", r.container.NotificationHandler.MarkAsRead)
}
