package router

import (
	"github.com/gcc798/microservice-kit/application/sys/internal/controller"
	"github.com/gcc798/microservice-kit/internal/constants"
	"github.com/gcc798/microservice-kit/internal/httpmiddleware"
	"github.com/gcc798/microservice-kit/internal/httpx"
)

// registerConfigRoutes 注册配置管理路由
func registerConfigRoutes(r *httpx.Router, ctx *RouterContext) {
	configController := controller.NewConfigController(ctx.Container)

	// API v1 路由组
	v1 := r.Group("/api/v1")
	{
		config := v1.Group("/config")
		config.Use(ctx.AuthMiddleware) // 添加认证中间件
		{
			// 创建配置 - 需要 config.create 权限
			config.POST("", middleware.Permission(ctx.PermissionService, constants.ResourceConfigCreate, "write"), configController.CreateConfig)

			// 分页查询配置列表 - 需要 config.read 权限
			config.POST("/page", middleware.Permission(ctx.PermissionService, constants.ResourceConfigRead, "read"), configController.PageConfig)

			// 批量删除配置 - 需要 config.delete 权限
			config.DELETE("/batch", middleware.Permission(ctx.PermissionService, constants.ResourceConfigDelete, "write"), configController.BatchDeleteConfig)

			// 根据编码获取配置列表 - 需要 config.read 权限
			config.GET("/code", middleware.Permission(ctx.PermissionService, constants.ResourceConfigRead, "read"), configController.GetConfigByCode)

			// 根据编码获取配置数据（仅返回data字段）- 需要 config.read 权限
			config.GET("/data", middleware.Permission(ctx.PermissionService, constants.ResourceConfigRead, "read"), configController.GetConfigDataByCode)

			// 更新、查询和删除配置 - 需要 config.update/read/delete 权限（带参数的路由放在最后，避免路径冲突）
			config.PUT("/:id", middleware.Permission(ctx.PermissionService, constants.ResourceConfigUpdate, "write"), configController.UpdateConfig)
			config.GET("/:id", middleware.Permission(ctx.PermissionService, constants.ResourceConfigRead, "read"), configController.GetConfigById)
			config.DELETE("/:id", middleware.Permission(ctx.PermissionService, constants.ResourceConfigDelete, "write"), configController.DeleteConfig)
		}
	}
}
