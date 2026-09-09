package router

import (
	"github.com/gcc798/lightning/application/iam/internal/controller"
	"github.com/gcc798/lightning/internal/constants"
	"github.com/gcc798/lightning/internal/httpmiddleware"
	"github.com/gcc798/lightning/internal/httpx"
)

// registerRoleRoutes 注册角色管理路由
func registerRoleRoutes(r *httpx.Router, ctx *RouterContext) {
	// 初始化 controller
	roleController := controller.NewRoleController(ctx.Container)

	// 角色管理路由组（需要认证和权限）
	roles := r.Group("/api/v1/role")
	roles.Use(ctx.AuthMiddleware)
	{
		// 角色创建 - 需要 role.create 权限
		roles.POST("", middleware.Permission(ctx.PermissionService, constants.ResourceRoleCreate, "write"), roleController.CreateRole)

		// 角色查询 - 需要 role.read 权限
		roles.POST("/page", middleware.Permission(ctx.PermissionService, constants.ResourceRoleRead, "read"), roleController.PageRole)

		// 用户角色管理 - 需要 role.assign 权限
		roles.POST("/assign", middleware.Permission(ctx.PermissionService, constants.ResourceRoleAssign, "write"), roleController.AssignRoleToUser)
		roles.DELETE("/remove", middleware.Permission(ctx.PermissionService, constants.ResourceRoleAssign, "write"), roleController.RemoveRoleFromUser)
		roles.GET("/user", middleware.Permission(ctx.PermissionService, constants.ResourceRoleRead, "read"), roleController.GetUserRoles)

		roles.GET("/:roleId/users", middleware.Permission(ctx.PermissionService, constants.ResourceRoleRead, "read"), roleController.GetRoleUsers)
		roles.POST("/:roleId/users", middleware.Permission(ctx.PermissionService, constants.ResourceRoleAssign, "write"), roleController.AssignUsersToRole)
		roles.DELETE("/:roleId/users", middleware.Permission(ctx.PermissionService, constants.ResourceRoleAssign, "write"), roleController.RemoveUsersFromRole)
		roles.GET("/:roleId/menus", middleware.Permission(ctx.PermissionService, constants.ResourceRoleRead, "read"), roleController.GetRoleMenus)
		roles.POST("/:roleId/menus", middleware.Permission(ctx.PermissionService, constants.ResourceRoleUpdate, "write"), roleController.AssignRoleMenus)

		// 角色更新、查询和删除 - 需要 role.update/read/delete 权限（带参数的路由放在最后）
		roles.PUT("/:roleId", middleware.Permission(ctx.PermissionService, constants.ResourceRoleUpdate, "write"), roleController.UpdateRole)
		roles.GET("/:roleId", middleware.Permission(ctx.PermissionService, constants.ResourceRoleRead, "read"), roleController.GetRole)
		roles.DELETE("/:roleId", middleware.Permission(ctx.PermissionService, constants.ResourceRoleDelete, "write"), roleController.DeleteRole)
	}
}
