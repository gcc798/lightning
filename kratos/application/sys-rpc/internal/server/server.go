package server

import (
	"strings"
	"time"

	"github.com/google/wire"

	"github.com/gcc798/microservice-kit/kratos/application/sys-rpc/internal/conf"
	"github.com/gcc798/microservice-kit/kratos/pkg/configx"
	"github.com/gcc798/microservice-kit/kratos/pkg/registryx"
	"github.com/go-kratos/kratos/v2/registry"
)

// ProviderSet is server providers.
var ProviderSet = wire.NewSet(
	ProvideGRPCConfig,
	NewGRPCServer,
	NewRegistrar,
)

func ProvideGRPCConfig(cfg *conf.Bootstrap) *conf.GRPC {
	if cfg == nil {
		return nil
	}
	return cfg.GetServer().GetGrpc()
}

func NewRegistrar(cfg *conf.Bootstrap) (registry.Registrar, func(), error) {
	mode := strings.ToLower(strings.TrimSpace(cfg.GetRegistry().GetMode()))
	if mode == "" || mode == "direct" || mode == "disabled" || mode == "off" {
		return nil, func() {}, nil
	}
	driver := strings.ToLower(strings.TrimSpace(cfg.GetRegistry().GetDriver()))
	if driver == "" {
		return nil, func() {}, nil
	}
	return registryx.NewRegistrar(registryx.RegistrarConfig{
		Driver:      driver,
		Namespace:   cfg.GetRegistry().GetNamespace(),
		RegisterTTL: configx.ParseDurationOrDefault(cfg.GetRegistry().GetRegisterTTL(), 15*time.Second),
		Etcd: registryx.EtcdConfig{
			Endpoints: cfg.GetRegistry().GetEtcd().GetEndpoints(),
			Username:  cfg.GetRegistry().GetEtcd().GetUsername(),
			Password:  cfg.GetRegistry().GetEtcd().GetPassword(),
			Timeout:   configx.ParseDurationOrDefault(cfg.GetRegistry().GetEtcd().GetTimeout(), 5*time.Second),
		},
	})
}
