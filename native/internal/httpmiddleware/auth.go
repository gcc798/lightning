package middleware

import (
	"context"
	"strings"

	"github.com/gcc798/lightning/internal/config"
	"github.com/gcc798/lightning/internal/httpresponse"
	"github.com/gcc798/lightning/internal/platform/jwt"
	"github.com/labstack/echo/v5"
)

type AccessTokenValidator interface {
	ValidateAccessToken(context.Context, string) (*jwt.Claims, error)
}

// Auth 认证中间件
// 1. 从配置的请求头读取 AccessToken，WebSocket 路径兼容 query Token
// 2. 验证 AccessToken
// 3. 设置用户信息到 context
func Auth(tokenManager AccessTokenValidator, cfg *config.Config) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			// 从配置的请求头读取 Token，WebSocket 路径兼容小程序 query 握手。
			tokenHeader := cfg.Auth.TokenHeader
			token := c.Request().Header.Get(tokenHeader)
			if token == "" && isWebSocketPath(c.Path()) {
				token = c.QueryParam(tokenHeader)
			}
			if token == "" {
				response.Unauthorized(c, "未登录")
				return nil
			}

			// 去除 "Bearer " 前缀
			token = strings.TrimPrefix(token, "Bearer ")

			// 验证 AccessToken
			claims, err := tokenManager.ValidateAccessToken(c.Request().Context(), token)
			if err != nil {
				response.Unauthorized(c, err.Error())
				return nil
			}

			// 可选：验证请求头中的 clientId 是否与 Token 中的一致
			headerClientId := c.Request().Header.Get("clientid")
			if headerClientId == "" {
				headerClientId = c.QueryParam("clientid")
			}
			if headerClientId != "" && claims.ClientId != headerClientId {
				response.Unauthorized(c, "客户端ID与Token不匹配")
				return nil
			}

			// 设置用户信息到 context
			c.Set("userId", claims.UserId)
			c.Set("userName", claims.UserName)
			c.Set("clientId", claims.ClientId)
			c.Set("deviceType", claims.DeviceType)
			return next(c)
		}
	}
}

func isWebSocketPath(path string) bool {
	return path == "/ws" || path == "/resource/websocket"
}
