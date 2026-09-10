package router

import (
	"github.com/gcc798/microservice-kit/application/sys/internal/controller"
	"github.com/gcc798/microservice-kit/internal/constants"
	"github.com/gcc798/microservice-kit/internal/httpmiddleware"
	"github.com/gcc798/microservice-kit/internal/httpx"
)

// registerOperLogRoutes 注册操作日志路由
func registerOperLogRoutes(r *httpx.Router, ctx *RouterContext) {
	operLogController := controller.NewOperLogController(ctx.Container)

	// API v1 路由组
	v1 := r.Group("/api/v1")
	{
		operLog := v1.Group("/operLog")
		operLog.Use(ctx.AuthMiddleware) // 添加认证中间件
		{
			// 创建操作日志 - 需要 oper_log.create 权限
			operLog.POST("", middleware.Permission(ctx.PermissionService, constants.ResourceOperLogCreate, "write"), operLogController.CreateOperLog)

			// 分页查询操作日志列表 - 需要 oper_log.read 权限
			operLog.POST("/page", middleware.Permission(ctx.PermissionService, constants.ResourceOperLogRead, "read"), operLogController.PageOperLog)

			// 批量删除操作日志 - 需要 oper_log.delete 权限
			operLog.DELETE("/batch", middleware.Permission(ctx.PermissionService, constants.ResourceOperLogDelete, "write"), operLogController.BatchDeleteOperLog)

			// 清理操作日志 - 需要 oper_log.delete 权限
			operLog.POST("/clean", middleware.Permission(ctx.PermissionService, constants.ResourceOperLogDelete, "write"), operLogController.CleanOperLog)

			// 更新、查询和删除操作日志 - 需要 oper_log.update/read/delete 权限（带参数的路由放在最后）
			operLog.PUT("/:id", middleware.Permission(ctx.PermissionService, constants.ResourceOperLogUpdate, "write"), operLogController.UpdateOperLog)
			operLog.GET("/:id", middleware.Permission(ctx.PermissionService, constants.ResourceOperLogRead, "read"), operLogController.GetOperLogById)
			operLog.DELETE("/:id", middleware.Permission(ctx.PermissionService, constants.ResourceOperLogDelete, "write"), operLogController.DeleteOperLog)
		}
	}
}
