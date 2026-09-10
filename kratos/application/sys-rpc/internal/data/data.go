package data

import (
	"github.com/google/wire"

	"github.com/gcc798/microservice-kit/kratos/application/sys-rpc/internal/conf"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(
	ProvideDataConfig,
	ProvideAuthConfig,
	ProvideJWTConfig,
	ProvideObservabilityConfig,
	NewData,
)

func ProvideDataConfig(cfg *conf.Bootstrap) *conf.Data {
	if cfg == nil {
		return nil
	}
	return cfg.GetData()
}

func ProvideAuthConfig(cfg *conf.Bootstrap) *conf.Auth {
	if cfg == nil {
		return nil
	}
	return cfg.GetAuth()
}

func ProvideJWTConfig(cfg *conf.Bootstrap) *conf.JWT {
	if cfg == nil {
		return nil
	}
	return cfg.GetJwt()
}

func ProvideObservabilityConfig(cfg *conf.Bootstrap) *conf.Observability {
	if cfg == nil {
		return nil
	}
	return cfg.GetObservability()
}

func NewData(dataCfg *conf.Data, authCfg *conf.Auth, jwtCfg *conf.JWT, obsCfg *conf.Observability) (*Resources, func(), error) {
	resources, err := NewResources(dataCfg, authCfg, jwtCfg, obsCfg)
	if err != nil {
		return nil, nil, err
	}
	return resources, func() {
		_ = resources.Close()
	}, nil
}
