package router

import (
	"fmt"

	"github.com/gcc798/microservice-kit/application/iam/internal/controller"
	"github.com/gcc798/microservice-kit/internal/httpx"
)

// 注册认证相关路由。
func registerAuthRoutes(r *httpx.Router, ctx *RouterContext) error {
	// 初始化controller
	authController, err := controller.NewAuthController(ctx.Container, ctx.SystemAPI)
	if err != nil {
		return fmt.Errorf("create auth controller: %w", err)
	}

	// 公开路由（无需认证）
	r.POST("/login", authController.Login)               // 统一登录接口
	r.POST("/logout", authController.Logout)             // 登出
	r.POST("/auth/refresh", authController.RefreshToken) // 刷新Token
	return nil
}
