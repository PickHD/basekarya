package routes

import (
	"basekarya-backend/pkg/constants"

	"github.com/labstack/echo/v4"
)

func (r *Router) SetupCompanyRoutes(e *echo.Group) {
	e.GET("/profile", r.container.CompanyHandler.GetProfile, r.container.AuthMiddleware.GrantPermission(constants.VIEW_COMPANY))
	e.PUT("/profile", r.container.CompanyHandler.UpdateProfile, r.container.AuthMiddleware.GrantPermission(constants.UPDATE_COMPANY))
}
