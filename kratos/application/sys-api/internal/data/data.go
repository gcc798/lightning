package data

import (
	"github.com/google/wire"

	"github.com/gcc798/microservice-kit/kratos/application/sys-api/internal/conf"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(
	ProvideRPCConfig,
	NewData,
	NewHealthRepo,
	NewAuthRepo,
	NewCaptchaRepo,
	NewMenuRepo,
	NewUserRepo,
	NewRoleRepo,
	NewOrgRepo,
	NewConfigRepo,
	NewDictRepo,
	NewLoginLogRepo,
	NewOperLogRepo,
	NewStorageEnvRepo,
	NewAttachmentRepo,
)

func ProvideRPCConfig(cfg *conf.Bootstrap) *conf.RPC {
	if cfg == nil {
		return nil
	}
	return cfg.GetServer().GetRpc()
}

func NewData(cfg *conf.RPC) (*RPCClientSet, func(), error) {
	clients, err := NewRPCClientSet(cfg)
	if err != nil {
		return nil, nil, err
	}
	return clients, func() {
		_ = clients.Close()
	}, nil
}
