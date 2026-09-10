package router

import (
	"github.com/gcc798/microservice-kit/application/iam/internal/controller"
	iam "github.com/gcc798/microservice-kit/application/iam/internal/domain"
	"github.com/gcc798/microservice-kit/internal/constants"
	"github.com/gcc798/microservice-kit/internal/httpmiddleware"
	"github.com/gcc798/microservice-kit/internal/httpx"
)

// registerMenuRoutes 注册菜单管理路由
func registerMenuRoutes(r *httpx.Router, ctx *RouterContext) {
	// 创建 MenuService
	menuService := iam.NewMenuService(ctx.Container.GetDB())

	// 创建 MenuController
	menuController := controller.NewMenuController(menuService)

	// 菜单管理路由组（需要认证和权限）
	menus := r.Group("/api/v1/menu")
	menus.Use(ctx.AuthMiddleware)
	{
		// 获取当前用户的菜单树（用于前端路由生成）- 所有登录用户都可以访问
		menus.GET("/user/tree", menuController.GetUserMenuTree)

		// 菜单查询 - 需要 menu.read 权限
		menus.GET("/tree", middleware.Permission(ctx.PermissionService, constants.ResourceMenuRead, "read"), menuController.GetMenuTree)
		menus.GET("", middleware.Permission(ctx.PermissionService, constants.ResourceMenuRead, "read"), menuController.GetMenuList)
		menus.GET("/:id", middleware.Permission(ctx.PermissionService, constants.ResourceMenuRead, "read"), menuController.GetMenuById)

		// 菜单创建 - 需要 menu.create 权限
		menus.POST("", middleware.Permission(ctx.PermissionService, constants.ResourceMenuCreate, "write"), menuController.CreateMenu)

		// 菜单更新 - 需要 menu.update 权限
		menus.PUT("/:id", middleware.Permission(ctx.PermissionService, constants.ResourceMenuUpdate, "write"), menuController.UpdateMenu)

		// 菜单删除 - 需要 menu.delete 权限
		menus.DELETE("/:id", middleware.Permission(ctx.PermissionService, constants.ResourceMenuDelete, "write"), menuController.DeleteMenu)
	}
}
