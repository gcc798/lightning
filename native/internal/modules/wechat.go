package modules

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/gcc798/microservice-kit/internal/platform/thirdparty/wechat"
	"github.com/gcc798/microservice-kit/internal/runtimeconfig"
)

type WeChatModule struct {
	mu      sync.RWMutex
	cont    Container
	digest  [sha256.Size]byte
	enabled bool
	manager *wechat.Manager
}

func NewWeChatModule() *WeChatModule { return &WeChatModule{} }
func (*WeChatModule) Name() string   { return WeChatName }

func (m *WeChatModule) Init(ctx context.Context, cont Container) error {
	m.cont = cont
	return m.reload(ctx)
}

func (*WeChatModule) Start(context.Context) error { return nil }
func (*WeChatModule) Stop(context.Context) error  { return nil }

func (m *WeChatModule) Refresh(ctx context.Context, req ModuleRefreshRequest) error {
	if !containsCode(req.Codes, runtimeconfig.CodeWeChat) {
		return nil
	}
	return m.reload(ctx)
}

func (m *WeChatModule) Code2Session(ctx context.Context, wxCode string) (*wechat.Code2SessionResponse, error) {
	if err := m.reloadIfChanged(ctx); err != nil {
		return nil, err
	}
	m.mu.RLock()
	enabled, manager := m.enabled, m.manager
	m.mu.RUnlock()
	if !enabled || manager == nil {
		return nil, fmt.Errorf("wechat: %w", ErrDisabled)
	}
	return manager.Code2Session(ctx, wxCode)
}

func (m *WeChatModule) reloadIfChanged(ctx context.Context) error {
	raw, err := m.cont.GetRuntimeConfig().GetRaw(ctx, runtimeconfig.CodeWeChat)
	if err != nil {
		return err
	}
	digest := sha256.Sum256(raw)
	m.mu.RLock()
	unchanged := digest == m.digest
	m.mu.RUnlock()
	if unchanged {
		return nil
	}
	return m.apply(raw, digest)
}

func (m *WeChatModule) reload(ctx context.Context) error {
	raw, err := m.cont.GetRuntimeConfig().GetRaw(ctx, runtimeconfig.CodeWeChat)
	if err != nil {
		return err
	}
	return m.apply(raw, sha256.Sum256(raw))
}

func (m *WeChatModule) apply(raw []byte, digest [sha256.Size]byte) error {
	var cfg runtimeconfig.WeChatConfig
	if err := json.Unmarshal(raw, &cfg); err != nil {
		return fmt.Errorf("decode WeChat configuration: %w", err)
	}
	var manager *wechat.Manager
	if cfg.Enabled {
		manager = wechat.NewManager(wechat.Config{Enabled: true, AppID: cfg.AppID, Secret: cfg.Secret}, m.cont.GetLogger(), m.cont.GetRedis())
	}
	m.mu.Lock()
	m.digest, m.enabled, m.manager = digest, cfg.Enabled, manager
	m.mu.Unlock()
	return nil
}
