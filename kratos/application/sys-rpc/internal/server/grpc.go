package server

import (
	"time"

	v1 "github.com/gcc798/microservice-kit/kratos/api/system/v1"
	"github.com/gcc798/microservice-kit/kratos/application/sys-rpc/internal/conf"
	"github.com/gcc798/microservice-kit/kratos/application/sys-rpc/internal/service"
	kgrpc "github.com/go-kratos/kratos/v2/transport/grpc"
)

func NewGRPCServer(c *conf.GRPC,
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
) *kgrpc.Server {
	timeout := time.Second
	if c.GetTimeout() != "" {
		if parsed, err := time.ParseDuration(c.GetTimeout()); err == nil {
			timeout = parsed
		}
	}
	srv := kgrpc.NewServer(kgrpc.Network(c.GetNetwork()), kgrpc.Address(c.GetAddr()), kgrpc.Timeout(timeout))
	v1.RegisterHealthServiceServer(srv, healthSvc)
	v1.RegisterAuthServiceServer(srv, authSvc)
	v1.RegisterCaptchaServiceServer(srv, captchaSvc)
	v1.RegisterMenuServiceServer(srv, menuSvc)
	v1.RegisterUserServiceServer(srv, userSvc)
	v1.RegisterRoleServiceServer(srv, roleSvc)
	v1.RegisterOrgServiceServer(srv, orgSvc)
	v1.RegisterConfigServiceServer(srv, configSvc)
	v1.RegisterDictServiceServer(srv, dictSvc)
	v1.RegisterLoginLogServiceServer(srv, loginLogSvc)
	v1.RegisterOperLogServiceServer(srv, operLogSvc)
	v1.RegisterStorageEnvServiceServer(srv, storageSvc)
	v1.RegisterAttachmentServiceServer(srv, attachmentSvc)
	return srv
}
