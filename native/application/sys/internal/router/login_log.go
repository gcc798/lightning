package router

import (
	"github.com/gcc798/lightning/application/sys/internal/controller"
	"github.com/gcc798/lightning/internal/constants"
	"github.com/gcc798/lightning/internal/httpmiddleware"
	"github.com/gcc798/lightning/internal/httpx"
)

// registerLoginLogRoutes 注册登录日志路由
func registerLoginLogRoutes(r *httpx.Router, ctx *RouterContext) {
	loginLogController := controller.NewLoginLogController(ctx.Container)

	// API v1 路由组
	v1 := r.Group("/api/v1")
	{
		loginLog := v1.Group("/loginLog")
		loginLog.Use(ctx.AuthMiddleware) // 添加认证中间件
		{
			// 创建登录日志 - 需要 login_log.create 权限
			loginLog.POST("", middleware.Permission(ctx.PermissionService, constants.ResourceLoginLogCreate, "write"), loginLogController.CreateLoginLog)

			// 分页查询登录日志列表 - 需要 login_log.read 权限
			loginLog.POST("/page", middleware.Permission(ctx.PermissionService, constants.ResourceLoginLogRead, "read"), loginLogController.PageLoginLog)

			// 批量删除登录日志 - 需要 login_log.delete 权限
			loginLog.DELETE("/batch", middleware.Permission(ctx.PermissionService, constants.ResourceLoginLogDelete, "write"), loginLogController.BatchDeleteLoginLog)

			// 清理登录日志 - 需要 login_log.delete 权限
			loginLog.POST("/clean", middleware.Permission(ctx.PermissionService, constants.ResourceLoginLogDelete, "write"), loginLogController.CleanLoginLog)

			// 更新、查询和删除登录日志 - 需要 login_log.update/read/delete 权限（带参数的路由放在最后）
			loginLog.PUT("/:id", middleware.Permission(ctx.PermissionService, constants.ResourceLoginLogUpdate, "write"), loginLogController.UpdateLoginLog)
			loginLog.GET("/:id", middleware.Permission(ctx.PermissionService, constants.ResourceLoginLogRead, "read"), loginLogController.GetLoginLogById)
			loginLog.DELETE("/:id", middleware.Permission(ctx.PermissionService, constants.ResourceLoginLogDelete, "write"), loginLogController.DeleteLoginLog)
		}
	}
}
