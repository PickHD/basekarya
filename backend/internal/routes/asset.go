package routes

import (
	"basekarya-backend/internal/middleware"
	"basekarya-backend/pkg/constants"

	"github.com/labstack/echo/v4"
)

func (r *Router) SetupAssetRoutes(e *echo.Group, sub *middleware.SubscriptionMiddleware) {
	g := e.Group("", sub.RequireModule("asset"))
	g.GET("/categories", r.container.AssetHandler.GetAllCategories, r.container.AuthMiddleware.GrantPermission(constants.VIEW_ASSET))
	g.GET("/categories/:id", r.container.AssetHandler.GetCategoryDetail, r.container.AuthMiddleware.GrantPermission(constants.VIEW_ASSET))
	g.POST("/categories", r.container.AssetHandler.CreateCategory, r.container.AuthMiddleware.GrantPermission(constants.MANAGE_ASSET))
	g.PUT("/categories/:id", r.container.AssetHandler.UpdateCategory, r.container.AuthMiddleware.GrantPermission(constants.MANAGE_ASSET))
	g.DELETE("/categories/:id", r.container.AssetHandler.DeleteCategory, r.container.AuthMiddleware.GrantPermission(constants.MANAGE_ASSET))

	g.GET("", r.container.AssetHandler.GetAllAssets, r.container.AuthMiddleware.GrantPermission(constants.VIEW_ASSET))
	g.GET("/:id", r.container.AssetHandler.GetAssetDetail, r.container.AuthMiddleware.GrantPermission(constants.VIEW_ASSET))
	g.POST("", r.container.AssetHandler.CreateAsset, r.container.AuthMiddleware.GrantPermission(constants.MANAGE_ASSET))
	g.PUT("/:id", r.container.AssetHandler.UpdateAsset, r.container.AuthMiddleware.GrantPermission(constants.MANAGE_ASSET))
	g.DELETE("/:id", r.container.AssetHandler.DeleteAsset, r.container.AuthMiddleware.GrantPermission(constants.MANAGE_ASSET))

	g.GET("/assignments", r.container.AssetHandler.GetAllAssignments, r.container.AuthMiddleware.GrantAnyPermission(constants.VIEW_ASSET, constants.VIEW_SELF_ASSET))
	g.GET("/assignments/:id", r.container.AssetHandler.GetAssignmentDetail, r.container.AuthMiddleware.GrantAnyPermission(constants.VIEW_ASSET, constants.VIEW_SELF_ASSET))
	g.POST("/assignments", r.container.AssetHandler.CreateAssignment, r.container.AuthMiddleware.GrantPermission(constants.CREATE_ASSET))
	g.PUT("/assignments/:id/action", r.container.AssetHandler.ProcessAction, r.container.AuthMiddleware.GrantPermission(constants.APPROVAL_ASSET))
	g.PUT("/assignments/:id/return", r.container.AssetHandler.ProcessReturn, r.container.AuthMiddleware.GrantAnyPermission(constants.CREATE_ASSET, constants.APPROVAL_ASSET))

	g.GET("/export", r.container.AssetHandler.ExportAssets, r.container.AuthMiddleware.GrantPermission(constants.EXPORT_ASSET))
}
