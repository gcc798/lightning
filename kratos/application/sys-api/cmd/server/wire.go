//go:build wireinject

package main

import (
	"github.com/google/wire"

	"github.com/gcc798/microservice-kit/kratos/application/sys-api/internal/biz"
	"github.com/gcc798/microservice-kit/kratos/application/sys-api/internal/conf"
	"github.com/gcc798/microservice-kit/kratos/application/sys-api/internal/data"
	"github.com/gcc798/microservice-kit/kratos/application/sys-api/internal/server"
	"github.com/gcc798/microservice-kit/kratos/application/sys-api/internal/service"
	kratos "github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	khttp "github.com/go-kratos/kratos/v2/transport/http"
)

func newApp(logger log.Logger, httpSrv *khttp.Server) *kratos.App {
	return kratos.New(kratos.Name("sys-api"), kratos.Logger(logger), kratos.Server(httpSrv))
}

func wireApp(cfg *conf.Bootstrap, logger log.Logger) (*kratos.App, func(), error) {
	panic(wire.Build(server.ProviderSet, data.ProviderSet, biz.ProviderSet, service.ProviderSet, newApp))
}
