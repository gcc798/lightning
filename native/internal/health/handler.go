package health

import (
	"context"
	"fmt"
	"time"

	"github.com/gcc798/microservice-kit/internal/container"
	"github.com/labstack/echo/v5"
)

var (
	startTime = time.Now()
	version   = "1.0.0"
)

// Handler 定义业务数据结构。
type Handler interface {
	Health(c *echo.Context)  // 基础健康检查
	Ready(c *echo.Context)   // Kubernetes就绪检查
	Live(c *echo.Context)    // Kubernetes存活检查
	Startup(c *echo.Context) // Kubernetes启动检查
}

type handler struct {
	container container.Container
}

// NewHandler 创建组件实例。
func NewHandler(c container.Container) Handler {
	return &handler{
		container: c,
	}
}

// HealthResponse 定义业务数据结构。
type HealthResponse struct {
	Status    string            `json:"status"`
	Timestamp time.Time         `json:"timestamp"`
	Version   string            `json:"version"`
	Uptime    string            `json:"uptime"`
	Services  map[string]string `json:"services"`
}

// Health 基础健康检查
//
//	@Summary		健康检查
//	@Description	检查服务是否正常运行
//	@Tags			系统监控
//	@Produce		json
//	@Success		200	{object}	HealthResponse	"服务正常"
//	@Failure		503	{object}	HealthResponse	"服务异常"
//	@Router			/health [get]
func (h *handler) Health(c *echo.Context) {
	db := h.container.GetDB()
	redis := h.container.GetRedis()

	services := make(map[string]string)
	overallStatus := "healthy"

	if sqlDB, err := db.DB(); err == nil {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		if err := sqlDB.PingContext(ctx); err != nil {
			services["database"] = "down"
			overallStatus = "unhealthy"
		} else {
			services["database"] = "up"
		}
	} else {
		services["database"] = "down"
		overallStatus = "unhealthy"
	}

	if redis != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		if err := redis.Ping(ctx).Err(); err != nil {
			services["redis"] = "down"
			overallStatus = "unhealthy"
		} else {
			services["redis"] = "up"
		}
	}

	uptime := time.Since(startTime)
	uptimeStr := formatDuration(uptime)

	response := HealthResponse{
		Status:    overallStatus,
		Timestamp: time.Now(),
		Version:   version,
		Uptime:    uptimeStr,
		Services:  services,
	}

	if overallStatus == "healthy" {
		c.JSON(200, response)
	} else {
		c.JSON(503, response)
	}
}

// Ready Kubernetes 就绪检查
//
//	@Summary		就绪检查
//	@Description	检查服务是否准备好接收流量（Kubernetes Readiness Probe）
//	@Tags			系统监控
//	@Produce		json
//	@Success		200	{object}	map[string]string	"服务就绪"
//	@Failure		503	{object}	map[string]string	"服务未就绪"
//	@Router			/health/ready [get]
func (h *handler) Ready(c *echo.Context) {
	db := h.container.GetDB()
	redis := h.container.GetRedis()

	if sqlDB, err := db.DB(); err == nil {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		if err := sqlDB.PingContext(ctx); err != nil {
			c.JSON(503, map[string]any{
				"status":  "not ready",
				"reason":  "database connection failed",
				"message": err.Error(),
			})
			return
		}
	} else {
		c.JSON(503, map[string]any{
			"status":  "not ready",
			"reason":  "database not initialized",
			"message": err.Error(),
		})
		return
	}

	if redis != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		if err := redis.Ping(ctx).Err(); err != nil {
			c.JSON(503, map[string]any{
				"status":  "not ready",
				"reason":  "redis connection failed",
				"message": err.Error(),
			})
			return
		}
	}

	c.JSON(200, map[string]any{
		"status":  "ready",
		"message": "service is ready to accept traffic",
	})
}

// Live Kubernetes 存活检查
//
//	@Summary		存活检查
//	@Description	检查服务是否存活（Kubernetes Liveness Probe）
//	@Tags			系统监控
//	@Produce		json
//	@Success		200	{object}	map[string]string	"服务存活"
//	@Router			/health/live [get]
func (h *handler) Live(c *echo.Context) {
	c.JSON(200, map[string]any{
		"status":  "alive",
		"message": "service is alive",
		"uptime":  formatDuration(time.Since(startTime)),
	})
}

// Startup Kubernetes 启动检查
//
//	@Summary		启动检查
//	@Description	检查服务是否已完成启动（Kubernetes Startup Probe）
//	@Tags			系统监控
//	@Produce		json
//	@Success		200	{object}	map[string]string	"服务已启动"
//	@Failure		503	{object}	map[string]string	"服务启动中"
//	@Router			/health/startup [get]
func (h *handler) Startup(c *echo.Context) {
	if time.Since(startTime) < 5*time.Second {
		c.JSON(503, map[string]any{
			"status":  "starting",
			"message": "service is still starting up",
		})
		return
	}

	db := h.container.GetDB()
	if sqlDB, err := db.DB(); err == nil {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		if err := sqlDB.PingContext(ctx); err != nil {
			c.JSON(503, map[string]any{
				"status":  "starting",
				"message": "database not ready",
			})
			return
		}
	} else {
		c.JSON(503, map[string]any{
			"status":  "starting",
			"message": "database not initialized",
		})
		return
	}

	c.JSON(200, map[string]any{
		"status":  "started",
		"message": "service has started successfully",
	})
}

func formatDuration(d time.Duration) string {
	days := int(d.Hours() / 24)
	hours := int(d.Hours()) % 24
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60

	if days > 0 {
		return fmt.Sprintf("%dd %dh %dm %ds", days, hours, minutes, seconds)
	} else if hours > 0 {
		return fmt.Sprintf("%dh %dm %ds", hours, minutes, seconds)
	} else if minutes > 0 {
		return fmt.Sprintf("%dm %ds", minutes, seconds)
	}
	return fmt.Sprintf("%ds", seconds)
}
