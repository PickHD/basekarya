package routes

import (
	"basekarya-backend/pkg/constants"

	"github.com/labstack/echo/v4"
)

func (r *Router) SetupEmployeeRoutes(e *echo.Group) {
	e.GET("", r.container.UserHandler.GetAllEmployees, r.container.AuthMiddleware.GrantPermission(constants.VIEW_EMPLOYEE))
	e.POST("", r.container.UserHandler.CreateEmployee, r.container.AuthMiddleware.GrantPermission(constants.CREATE_EMPLOYEE))
	e.PUT("/:id", r.container.UserHandler.UpdateEmployee, r.container.AuthMiddleware.GrantPermission(constants.UPDATE_EMPLOYEE))
	e.DELETE("/:id", r.container.UserHandler.DeleteEmployee, r.container.AuthMiddleware.GrantPermission(constants.DELETE_EMPLOYEE))
}