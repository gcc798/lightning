package router

import (
	"github.com/gcc798/lightning/application/sys/internal/controller"
	"github.com/gcc798/lightning/internal/constants"
	"github.com/gcc798/lightning/internal/httpmiddleware"
	"github.com/gcc798/lightning/internal/httpx"
)

// registerDictRoutes 注册字典管理路由
func registerDictRoutes(r *httpx.Router, ctx *RouterContext) {
	dictController := controller.NewDictController(ctx.Container)

	// API v1 路由组
	v1 := r.Group("/api/v1")
	{
		dict := v1.Group("/dict")
		dict.Use(ctx.AuthMiddleware) // 添加认证中间件
		{
			// 字典管理（需要认证和权限）
			dict.POST("", middleware.Permission(ctx.PermissionService, constants.ResourceDictCreate, "write"), dictController.CreateDict)              // 创建字典
			dict.POST("/page", middleware.Permission(ctx.PermissionService, constants.ResourceDictRead, "read"), dictController.PageDict)              // 分页查询字典列表
			dict.DELETE("/batch", middleware.Permission(ctx.PermissionService, constants.ResourceDictDelete, "write"), dictController.BatchDeleteDict) // 批量删除字典

			// 字典数据获取（需要认证，但权限要求较低）
			dict.GET("/type", middleware.Permission(ctx.PermissionService, constants.ResourceDictRead, "read"), dictController.GetDictByType) // 根据类型获取字典列表
			dict.GET("/label", middleware.Permission(ctx.PermissionService, constants.ResourceDictRead, "read"), dictController.GetDictLabel) // 根据类型和键值获取标签

			// 字典更新和删除（带参数的路由放在最后）
			dict.PUT("/:id", middleware.Permission(ctx.PermissionService, constants.ResourceDictUpdate, "write"), dictController.UpdateDict)    // 更新字典
			dict.GET("/:id", middleware.Permission(ctx.PermissionService, constants.ResourceDictRead, "read"), dictController.GetDictById)      // 根据ID查询字典
			dict.DELETE("/:id", middleware.Permission(ctx.PermissionService, constants.ResourceDictDelete, "write"), dictController.DeleteDict) // 删除字典
		}
	}
}
