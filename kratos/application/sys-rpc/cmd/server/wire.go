//go:build wireinject

package main

import (
	"strings"

	"github.com/google/wire"

	"github.com/gcc798/lightning/kratos/application/sys-rpc/internal/biz"
	"github.com/gcc798/lightning/kratos/application/sys-rpc/internal/conf"
	"github.com/gcc798/lightning/kratos/application/sys-rpc/internal/data"
	"github.com/gcc798/lightning/kratos/application/sys-rpc/internal/server"
	"github.com/gcc798/lightning/kratos/application/sys-rpc/internal/service"
	kratos "github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/registry"
	kgrpc "github.com/go-kratos/kratos/v2/transport/grpc"
)

func newApp(cfg *conf.Bootstrap, logger log.Logger, grpcSrv *kgrpc.Server, registrar registry.Registrar) *kratos.App {
	serviceName := strings.TrimSpace(cfg.GetRegistry().GetService())
	if serviceName == "" {
		serviceName = "sys-rpc"
	}
	options := []kratos.Option{kratos.Name(serviceName), kratos.Logger(logger), kratos.Server(grpcSrv)}
	if registrar != nil {
		options = append(options, kratos.Registrar(registrar))
	}
	return kratos.New(options...)
}

func wireApp(cfg *conf.Bootstrap, logger log.Logger) (*kratos.App, func(), error) {
	panic(wire.Build(server.ProviderSet, data.ProviderSet, biz.ProviderSet, service.ProviderSet, newApp))
}
