package server

import (
	"context"
	"strings"

	"github.com/gcc798/lightning/kratos/application/sys-api/internal/data"
	kerrors "github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
)

func authMiddleware(deps *GatewayDeps) middleware.Middleware {
	return func(next middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req any) (any, error) {
			operation := transportOperation(ctx)
			if isPublicOperation(operation) {
				return next(ctx, req)
			}
			if deps == nil || deps.Auth == nil {
				return nil, kerrors.InternalServer("AUTH_GATEWAY_MISSING", "认证服务未初始化")
			}
			token := authorizationToken(ctx)
			if token == "" {
				return nil, kerrors.Unauthorized("UNAUTHORIZED", "未登录")
			}
			validated, err := deps.Auth.ValidateAccessToken(ctx, token)
			if err != nil {
				return nil, err
			}
			if validated == nil || !validated.GetValid() || validated.GetUserId() <= 0 {
				return nil, kerrors.Unauthorized("UNAUTHORIZED", "登录已失效")
			}
			ctx = data.WithCurrentAuth(ctx, validated.GetUserId(), validated.GetUserName(), token)
			return next(ctx, req)
		}
	}
}

func permissionMiddleware(deps *GatewayDeps) middleware.Middleware {
	return func(next middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req any) (any, error) {
			operation := transportOperation(ctx)
			permission, ok := permissionForOperation(operation)
			if !ok {
				return next(ctx, req)
			}
			if deps == nil || deps.Auth == nil {
				return nil, kerrors.InternalServer("AUTH_GATEWAY_MISSING", "权限服务未初始化")
			}
			userID := data.CurrentUserID(ctx)
			if userID <= 0 {
				return nil, kerrors.Unauthorized("UNAUTHORIZED", "未登录")
			}
			allowed, err := deps.Auth.CheckPermission(ctx, userID, permission, permissionAction(permission))
			if err != nil {
				return nil, err
			}
			if allowed == nil || !allowed.GetAllowed() {
				return nil, kerrors.Forbidden("FORBIDDEN", "无权限访问")
			}
			return next(ctx, req)
		}
	}
}

func transportOperation(ctx context.Context) string {
	if tr, ok := transport.FromServerContext(ctx); ok {
		return strings.TrimSpace(tr.Operation())
	}
	return ""
}

func authorizationToken(ctx context.Context) string {
	if tr, ok := transport.FromServerContext(ctx); ok {
		value := strings.TrimSpace(tr.RequestHeader().Get("Authorization"))
		if value == "" {
			return ""
		}
		parts := strings.SplitN(value, " ", 2)
		if len(parts) == 2 {
			value = strings.TrimSpace(parts[1])
		}
		return value
	}
	return ""
}

func permissionAction(resource string) string {
	if strings.HasSuffix(strings.TrimSpace(resource), ".read") {
		return "read"
	}
	return "write"
}
