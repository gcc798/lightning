package router

import (
	"fmt"

	iamv1 "github.com/gcc798/lightning/internal/api/iam/v1"
	sysv1 "github.com/gcc798/lightning/internal/api/sys/v1"
	"github.com/gcc798/lightning/internal/container"
	middleware "github.com/gcc798/lightning/internal/httpmiddleware"
	"github.com/gcc798/lightning/internal/httpx"
	"github.com/labstack/echo/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type RouterContext struct {
	Container         container.Container
	PermissionService middleware.PermissionChecker
	AuthMiddleware    echo.MiddlewareFunc
	SystemAPI         sysv1.API
}

func Setup(r *httpx.Router, c container.Container, security iamv1.API, systemAPI sysv1.API) error {
	ctx := &RouterContext{Container: c, PermissionService: security, AuthMiddleware: middleware.Auth(security, c.GetConfig()), SystemAPI: systemAPI}
	r.Use(middleware.PrometheusMiddleware())
	r.GET("/metrics", ctx.AuthMiddleware, echo.WrapHandler(promhttp.Handler()))
	registerCommonRoutes(r, ctx, true)
	if err := registerAuthRoutes(r, ctx); err != nil {
		return fmt.Errorf("auth routes: %w", err)
	}
	if err := registerCaptchaRoutes(r, ctx); err != nil {
		return fmt.Errorf("captcha routes: %w", err)
	}
	registerUserRoutes(r, ctx)
	registerRoleRoutes(r, ctx)
	registerApiPermissionRoutes(r, ctx)
	registerOrgRoutes(r, ctx)
	registerMenuRoutes(r, ctx)
	return nil
}
