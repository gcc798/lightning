package router

import (
	"github.com/gcc798/microservice-kit/application/iam/internal/controller"
	"github.com/gcc798/microservice-kit/internal/constants"
	"github.com/gcc798/microservice-kit/internal/httpmiddleware"
	"github.com/gcc798/microservice-kit/internal/httpx"
)

// registerUserRoutes 注册用户管理路由
func registerUserRoutes(r *httpx.Router, ctx *RouterContext) {
	// 初始化 controller
	userController := controller.NewUserController(ctx.Container)

	// 用户管理路由组（需要认证和权限）
	users := r.Group("/api/v1/user")
	users.Use(ctx.AuthMiddleware)
	{
		// 用户创建 - 需要 user.create 权限
		users.POST("", middleware.Permission(ctx.PermissionService, constants.ResourceUserCreate, "write"), userController.Create)
		users.POST("/import", middleware.Permission(ctx.PermissionService, constants.ResourceUserCreate, "write"), userController.BatchImport)

		// 用户查询 - 需要 user.read 权限
		users.POST("/page", middleware.Permission(ctx.PermissionService, constants.ResourceUserRead, "read"), userController.PageUser)

		// 批量删除 - 需要 user.delete 权限
		users.DELETE("/batch", middleware.Permission(ctx.PermissionService, constants.ResourceUserDelete, "write"), userController.BatchDelete)

		// 用户修改密码 - 不需要特殊权限（用户修改自己的密码）
		users.POST("/password/change", userController.ChangePassword)

		// 用户更新 - 需要 user.update 权限（带参数的路由放在后面）
		users.PUT("/:id", middleware.Permission(ctx.PermissionService, constants.ResourceUserUpdate, "write"), userController.Update)
		users.PUT("/:id/password", middleware.Permission(ctx.PermissionService, constants.ResourceUserUpdate, "write"), userController.ResetPassword)

		// 用户查询 - 需要 user.read 权限（带参数的路由放在最后）
		users.GET("/:id", middleware.Permission(ctx.PermissionService, constants.ResourceUserRead, "read"), userController.GetById)

		// 用户删除 - 需要 user.delete 权限（带参数的路由放在最后）
		users.DELETE("/:id", middleware.Permission(ctx.PermissionService, constants.ResourceUserDelete, "write"), userController.Delete)
	}

	systemUsers := r.Group("/system/user", ctx.AuthMiddleware)
	{
		systemUsers.POST("/xcxGetInfo", userController.XcxGetInfo)
	}
}
