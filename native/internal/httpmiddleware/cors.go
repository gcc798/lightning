package middleware

import (
	"net/http"

	"github.com/labstack/echo/v5"
)

// CORS 跨域中间件
func CORS() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			method := c.Request().Method
			origin := c.Request().Header.Get("Origin")

			// 允许所有来源（开发环境）
			// 生产环境建议配置具体的域名
			if origin != "" {
				c.Response().Header().Set("Access-Control-Allow-Origin", origin)
			} else {
				c.Response().Header().Set("Access-Control-Allow-Origin", "*")
			}

			// 允许的请求方法
			c.Response().Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH, OPTIONS")

			// 允许的请求头。
			// 浏览器在预检请求里可能会把自定义头名规范成小写，因此同时放行 clientId/clientid。
			c.Response().Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With, clientId, clientid")

			// 允许浏览器访问的响应头
			c.Response().Header().Set("Access-Control-Expose-Headers", "Content-Length, Content-Type, Authorization")

			// 允许携带凭证（cookies）
			c.Response().Header().Set("Access-Control-Allow-Credentials", "true")

			// 预检请求缓存时间（秒）
			c.Response().Header().Set("Access-Control-Max-Age", "86400")

			// 处理 OPTIONS 预检请求
			if method == "OPTIONS" {
				return c.NoContent(http.StatusNoContent)
			}

			return next(c)
		}
	}
}
