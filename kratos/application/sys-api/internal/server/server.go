package server

import (
	"github.com/google/wire"

	"github.com/gcc798/microservice-kit/kratos/application/sys-api/internal/biz"
	"github.com/gcc798/microservice-kit/kratos/application/sys-api/internal/conf"
	"github.com/gcc798/microservice-kit/kratos/application/sys-api/internal/data"
)

// ProviderSet is server providers.
var ProviderSet = wire.NewSet(
	ProvideHTTPConfig,
	NewWebSocketHubProvider,
	ProvideAuthGateway,
	ProvideAttachmentGateway,
	ProvideOperLogGateway,
	NewGatewayDeps,
	NewHTTPServer,
)

func ProvideHTTPConfig(cfg *conf.Bootstrap) *conf.HTTP {
	if cfg == nil {
		return nil
	}
	return cfg.GetServer().GetHttp()
}

func NewWebSocketHubProvider() (*WebSocketHub, func(), error) {
	hub := NewWebSocketHub()
	return hub, func() {
		hub.Stop()
	}, nil
}

func ProvideAuthGateway(repo *data.AuthRepo) AuthGateway {
	return repo
}

func ProvideAttachmentGateway(uc *biz.AttachmentUsecase) AttachmentGateway {
	return uc
}

func ProvideOperLogGateway(uc *biz.OperLogUsecase) OperLogGateway {
	return uc
}

func NewGatewayDeps(auth AuthGateway, attachment AttachmentGateway, operLog OperLogGateway, wsHub *WebSocketHub) *GatewayDeps {
	return &GatewayDeps{
		Auth:       auth,
		Attachment: attachment,
		OperLog:    operLog,
		WebSocket:  wsHub,
	}
}
