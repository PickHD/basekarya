package routes

import (
	"basekarya-backend/pkg/constants"

	"github.com/labstack/echo/v4"
)

func (r *Router) SetupAnnouncementRoutes(e *echo.Group) {
	e.POST("/publish", r.container.AnnouncementHandler.PublishAnnouncement, r.container.AuthMiddleware.GrantPermission(constants.CREATE_ANNOUNCEMENT))
}
