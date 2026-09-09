package data

import (
	"context"
	"fmt"
	"strings"
	"time"

	v1 "github.com/gcc798/lightning/kratos/api/system/v1"
	"github.com/gcc798/lightning/kratos/application/sys-api/internal/conf"
	"github.com/gcc798/lightning/kratos/pkg/configx"
	"github.com/gcc798/lightning/kratos/pkg/grpcx"
	"github.com/gcc798/lightning/kratos/pkg/registryx"
	grpc "google.golang.org/grpc"
)

type RPCClientSet struct {
	conn       *grpc.ClientConn
	cleanup    func()
	Health     v1.HealthServiceClient
	Auth       v1.AuthServiceClient
	Captcha    v1.CaptchaServiceClient
	Menu       v1.MenuServiceClient
	User       v1.UserServiceClient
	Role       v1.RoleServiceClient
	Org        v1.OrgServiceClient
	Config     v1.ConfigServiceClient
	Dict       v1.DictServiceClient
	LoginLog   v1.LoginLogServiceClient
	OperLog    v1.OperLogServiceClient
	StorageEnv v1.StorageEnvServiceClient
	Attachment v1.AttachmentServiceClient
}

func NewRPCClientSet(cfg *conf.RPC) (*RPCClientSet, error) {
	conn, cleanup, err := newRPCConn(cfg)
	if err != nil {
		return nil, err
	}
	return &RPCClientSet{
		conn:       conn,
		cleanup:    cleanup,
		Health:     v1.NewHealthServiceClient(conn),
		Auth:       v1.NewAuthServiceClient(conn),
		Captcha:    v1.NewCaptchaServiceClient(conn),
		Menu:       v1.NewMenuServiceClient(conn),
		User:       v1.NewUserServiceClient(conn),
		Role:       v1.NewRoleServiceClient(conn),
		Org:        v1.NewOrgServiceClient(conn),
		Config:     v1.NewConfigServiceClient(conn),
		Dict:       v1.NewDictServiceClient(conn),
		LoginLog:   v1.NewLoginLogServiceClient(conn),
		OperLog:    v1.NewOperLogServiceClient(conn),
		StorageEnv: v1.NewStorageEnvServiceClient(conn),
		Attachment: v1.NewAttachmentServiceClient(conn),
	}, nil
}

func (c *RPCClientSet) Close() error {
	if c == nil {
		return nil
	}
	if c.cleanup != nil {
		c.cleanup()
	}
	if c.conn == nil {
		return nil
	}
	return c.conn.Close()
}

func newRPCConn(cfg *conf.RPC) (*grpc.ClientConn, func(), error) {
	timeout := 2 * time.Second
	if strings.TrimSpace(cfg.GetTimeout()) != "" {
		parsed, err := time.ParseDuration(strings.TrimSpace(cfg.GetTimeout()))
		if err != nil {
			return nil, nil, fmt.Errorf("parse rpc timeout: %w", err)
		}
		timeout = parsed
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return grpcx.DialInsecure(ctx, grpcx.ClientConfig{
		Mode:     cfg.GetMode(),
		Endpoint: cfg.GetEndpoint(),
		Timeout:  timeout,
		Discovery: registryx.DiscoveryConfig{
			Driver:    cfg.GetDiscovery().GetDriver(),
			Service:   cfg.GetDiscovery().GetService(),
			Namespace: cfg.GetDiscovery().GetEtcd().GetNamespace(),
			Etcd: registryx.EtcdConfig{
				Endpoints: cfg.GetDiscovery().GetEtcd().GetEndpoints(),
				Username:  cfg.GetDiscovery().GetEtcd().GetUsername(),
				Password:  cfg.GetDiscovery().GetEtcd().GetPassword(),
				Timeout:   configx.ParseDurationOrDefault(cfg.GetDiscovery().GetEtcd().GetTimeout(), 5*time.Second),
			},
		},
	}, outgoingUserInterceptor)
}
