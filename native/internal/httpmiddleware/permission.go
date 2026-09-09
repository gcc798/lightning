package middleware

import (
	"context"

	"github.com/gcc798/lightning/internal/httpresponse"
	"github.com/labstack/echo/v5"
)

type PermissionChecker interface {
	CheckPermission(context.Context, int64, string, string) (bool, error)
}

// Permission checks one API permission after authentication.
func Permission(permissionService PermissionChecker, resource, action string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			userID, ok := contextUserID(c)
			if !ok {
				response.Forbidden(c, "用户信息不存在")
				return nil
			}
			allowed, err := permissionService.CheckPermission(c.Request().Context(), userID, resource, action)
			if err != nil {
				response.InternalServerError(c, "权限检查失败: "+err.Error())
				return nil
			}
			if !allowed {
				response.Forbidden(c, "无权限访问")
				return nil
			}
			return next(c)
		}
	}
}

func contextUserID(c *echo.Context) (int64, bool) {
	value := c.Get("userId")
	userID, ok := value.(int64)
	return userID, ok
}
