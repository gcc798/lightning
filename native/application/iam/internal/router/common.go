package router

import (
	"time"

	"github.com/gcc798/lightning/internal/health"
	"github.com/gcc798/lightning/internal/httpx"
	"github.com/gcc798/lightning/internal/modules"
	"github.com/gcc798/lightning/internal/platform/websocket"
)

// 注册公共路由（健康检查等）。
func registerCommonRoutes(r *httpx.Router, ctx *RouterContext, websocketEnabled bool) {
	c := ctx.Container
	logger := c.GetLogger()

	// 初始化健康检查控制器
	healthController := health.NewHandler(c)

	// 健康检查接口（公开接口，无需认证）
	r.GET("/health", healthController.Health)          // 基础健康检查
	r.GET("/health/ready", healthController.Ready)     // 就绪探针
	r.GET("/health/live", healthController.Live)       // 存活探针
	r.GET("/health/startup", healthController.Startup) // 启动探针

	if !websocketEnabled {
		return
	}

	// 实时连接路由（需要认证）
	wsHandler := websocket.NewHandler(modules.WebSocketHub(c), logger, c.GetConfig().CORS.Enabled)
	wsCfg := c.GetConfig().WebSocket
	if wsCfg.TimeoutEnabled {
		wsHandler.ConfigureTimeouts(
			time.Duration(wsCfg.ReadTimeoutSeconds)*time.Second,
			time.Duration(wsCfg.WriteTimeoutSeconds)*time.Second,
		)
	}
	if wsCfg.HeartbeatEnabled {
		wsHandler.RegisterHeartbeatMessageBuilder(wsCfg.MaxReadTimeouts, defaultHeartbeat)
	}
	r.GET("/resource/websocket", ctx.AuthMiddleware, wsHandler.ServeWs)
}

// defaultHeartbeat 读超时后服务端主动下发的探活帧，业务可按需替换。
func defaultHeartbeat(*websocket.Client) interface{} {
	return map[string]any{"type": "ping"}
}
