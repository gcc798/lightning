package router

import (
	"github.com/gcc798/microservice-kit/application/iam/internal/controller"
	"github.com/gcc798/microservice-kit/internal/constants"
	"github.com/gcc798/microservice-kit/internal/httpmiddleware"
	"github.com/gcc798/microservice-kit/internal/httpx"
)

// registerApiPermissionRoutes 注册 API 权限管理路由。
func registerApiPermissionRoutes(r *httpx.Router, ctx *RouterContext) {
	apiPermissionController := controller.NewApiPermissionController(ctx.Container)

	apiPermissions := r.Group("/api/v1/api-permission")
	apiPermissions.Use(ctx.AuthMiddleware)
	{
		apiPermissions.GET("/tree", middleware.Permission(ctx.PermissionService, constants.ResourceApiPermissionRead, "read"), apiPermissionController.Tree)
		apiPermissions.GET("", middleware.Permission(ctx.PermissionService, constants.ResourceApiPermissionRead, "read"), apiPermissionController.List)
		apiPermissions.POST("", middleware.Permission(ctx.PermissionService, constants.ResourceApiPermissionCreate, "write"), apiPermissionController.Create)
		apiPermissions.PUT("/:id", middleware.Permission(ctx.PermissionService, constants.ResourceApiPermissionUpdate, "write"), apiPermissionController.Update)
		apiPermissions.DELETE("/:id", middleware.Permission(ctx.PermissionService, constants.ResourceApiPermissionDelete, "write"), apiPermissionController.Delete)
	}

	roles := r.Group("/api/v1/role")
	roles.Use(ctx.AuthMiddleware)
	{
		roles.GET("/:roleId/api-permissions", middleware.Permission(ctx.PermissionService, constants.ResourceApiPermissionAssign, "write"), apiPermissionController.GetRolePermissions)
		roles.POST("/:roleId/api-permissions", middleware.Permission(ctx.PermissionService, constants.ResourceApiPermissionAssign, "write"), apiPermissionController.AssignRolePermissions)
	}

	users := r.Group("/api/v1/user")
	users.Use(ctx.AuthMiddleware)
	{
		users.GET("/:id/api-permissions", middleware.Permission(ctx.PermissionService, constants.ResourceApiPermissionAssign, "write"), apiPermissionController.GetUserPermissions)
		users.POST("/:id/api-permissions", middleware.Permission(ctx.PermissionService, constants.ResourceApiPermissionAssign, "write"), apiPermissionController.AssignUserPermissions)
	}
}
