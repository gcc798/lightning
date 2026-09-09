package server

import (
	v1 "github.com/gcc798/lightning/kratos/api/system/v1"
	"github.com/gcc798/lightning/kratos/application/sys-api/internal/service"
	khttp "github.com/go-kratos/kratos/v2/transport/http"
)

// registerProtoHTTPRoutes keeps the default Kratos path: routes that can be
// described cleanly by proto annotations stay in generated Register*HTTPServer.
func registerProtoHTTPRoutes(
	srv *khttp.Server,
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
) {
	v1.RegisterHealthServiceHTTPServer(srv, healthSvc)
	v1.RegisterAuthServiceHTTPServer(srv, authSvc)
	v1.RegisterCaptchaServiceHTTPServer(srv, captchaSvc)
	v1.RegisterMenuServiceHTTPServer(srv, menuSvc)
	v1.RegisterUserServiceHTTPServer(srv, userSvc)
	v1.RegisterRoleServiceHTTPServer(srv, roleSvc)
	v1.RegisterOrgServiceHTTPServer(srv, orgSvc)
	v1.RegisterConfigServiceHTTPServer(srv, configSvc)
	v1.RegisterDictServiceHTTPServer(srv, dictSvc)
	v1.RegisterLoginLogServiceHTTPServer(srv, loginLogSvc)
	v1.RegisterOperLogServiceHTTPServer(srv, operLogSvc)
	v1.RegisterStorageEnvServiceHTTPServer(srv, storageSvc)
}

// registerManualHTTPRoutes keeps the exceptional cases explicit: multipart
// upload, binary download, probes, metrics, swagger and websocket all need
// transport-specific behavior that proto-generated HTTP handlers do not model
// well.
func registerManualHTTPRoutes(srv *khttp.Server, deps *GatewayDeps, attachmentSvc *service.AttachmentServiceService) {
	registerAttachmentHTTPRoutes(srv, deps, attachmentSvc)
	registerMetricsEndpoint(srv)
	registerHealthEndpoints(srv)
	registerSwaggerEndpoints(srv)
	registerWebSocketEndpoint(srv, deps)
}
