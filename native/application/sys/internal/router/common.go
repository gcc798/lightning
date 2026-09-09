package router

import (
	"github.com/gcc798/lightning/internal/health"
	"github.com/gcc798/lightning/internal/httpx"
)

// 注册公共路由（健康检查等）。
func registerCommonRoutes(r *httpx.Router, ctx *RouterContext, _ bool) {
	c := ctx.Container
	healthController := health.NewHandler(c)

	// 健康检查接口（公开接口，无需认证）
	r.GET("/health", healthController.Health)          // 基础健康检查
	r.GET("/health/ready", healthController.Ready)     // 就绪探针
	r.GET("/health/live", healthController.Live)       // 存活探针
	r.GET("/health/startup", healthController.Startup) // 启动探针

}
