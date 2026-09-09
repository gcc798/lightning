package modules

import (
	"context"

	"github.com/gcc798/lightning/internal/platform/websocket"
)

// WebSocketModule 拥有进程内 WebSocket Hub 的生命周期。
// 关闭时禁用该模块，Hub 为 nil，处理器会直接返回不可用。
type WebSocketModule struct {
	enabled bool
	hub     *websocket.Hub
}

func NewWebSocketModule(enabled bool) *WebSocketModule {
	return &WebSocketModule{enabled: enabled}
}

func (*WebSocketModule) Name() string { return WebSocketName }

func (m *WebSocketModule) Init(_ context.Context, cont Container) error {
	if m.enabled {
		m.hub = websocket.NewHub(cont.GetLogger())
	}
	return nil
}

func (m *WebSocketModule) Start(context.Context) error {
	if m.hub != nil {
		m.hub.Start()
	}
	return nil
}

func (m *WebSocketModule) Stop(context.Context) error {
	if m.hub != nil {
		m.hub.Close()
	}
	return nil
}

func (*WebSocketModule) Refresh(context.Context, ModuleRefreshRequest) error { return nil }

// Hub 返回连接管理中心，模块被禁用时返回 nil。
func (m *WebSocketModule) Hub() *websocket.Hub { return m.hub }
