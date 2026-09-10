package server

import (
	"context"
	"time"

	"github.com/gcc798/microservice-kit/kratos/application/sys-api/internal/conf"
	"github.com/gcc798/microservice-kit/kratos/application/sys-api/internal/service"
	selector "github.com/go-kratos/kratos/v2/middleware/selector"
	khttp "github.com/go-kratos/kratos/v2/transport/http"
)

func NewHTTPServer(c *conf.HTTP,
	deps *GatewayDeps,
	healthSvc *service.HealthServiceService,
	authSvc *service.AuthServiceService,
	captchaSvc *service.CaptchaServiceService,
	menuSvc *service.MenuServiceService,
	userSvc *service.UserServiceService,
	roleSvc *service.RoleServiceService,
	orgSvc *service.OrgServiceService,
	configSvc *service.ConfigServiceService,
	dictSvc *service.DictServiceService,
	loginLogSvc *service.LoginLogServiceService,
	operLogSvc *service.OperLogServiceService,
	storageSvc *service.StorageEnvServiceService,
	attachmentSvc *service.AttachmentServiceService,
) *khttp.Server {
	timeout := time.Second
	if c.GetTimeout() != "" {
		if parsed, err := time.ParseDuration(c.GetTimeout()); err == nil {
			timeout = parsed
		}
	}
	srv := khttp.NewServer(
		khttp.Network(c.GetNetwork()),
		khttp.Address(c.GetAddr()),
		khttp.Timeout(timeout),
		khttp.ResponseEncoder(responseEncoder),
		khttp.ErrorEncoder(errorEncoder),
		khttp.Middleware(
			metricsMiddleware(),
			selector.Server(authMiddleware(deps)).Match(func(ctx context.Context, operation string) bool {
				return !isPublicOperation(operation)
			}).Build(),
			selector.Server(permissionMiddleware(deps)).Match(func(ctx context.Context, operation string) bool {
				_, ok := permissionForOperation(operation)
				return ok
			}).Build(),
			operLogMiddleware(deps),
		),
	)
	registerProtoHTTPRoutes(srv, healthSvc, authSvc, captchaSvc, menuSvc, userSvc, roleSvc, orgSvc, configSvc, dictSvc, loginLogSvc, operLogSvc, storageSvc)
	registerManualHTTPRoutes(srv, deps, attachmentSvc)
	return srv
}
