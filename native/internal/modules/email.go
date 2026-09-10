package modules

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"sync"

	thirdpartyemail "github.com/gcc798/microservice-kit/internal/platform/thirdparty/email"
	"github.com/gcc798/microservice-kit/internal/runtimeconfig"
)

type EmailModule struct {
	mu      sync.RWMutex
	cont    Container
	digest  [sha256.Size]byte
	enabled bool
	manager *thirdpartyemail.Manager
}

func NewEmailModule() *EmailModule { return &EmailModule{} }
func (*EmailModule) Name() string  { return EmailName }

func (m *EmailModule) Init(ctx context.Context, cont Container) error {
	m.cont = cont
	return m.reload(ctx)
}

func (*EmailModule) Start(context.Context) error { return nil }
func (*EmailModule) Stop(context.Context) error  { return nil }

func (m *EmailModule) Refresh(ctx context.Context, req ModuleRefreshRequest) error {
	if !containsCode(req.Codes, runtimeconfig.CodeEmail) {
		return nil
	}
	return m.reload(ctx)
}

func (m *EmailModule) SendEmail(ctx context.Context, address, code, template string) error {
	manager, err := m.current(ctx)
	if err != nil {
		return err
	}
	return manager.SendWithTemplate(address, code, template)
}

func (m *EmailModule) current(ctx context.Context) (*thirdpartyemail.Manager, error) {
	if err := m.reloadIfChanged(ctx); err != nil {
		return nil, err
	}
	m.mu.RLock()
	defer m.mu.RUnlock()
	if !m.enabled || m.manager == nil {
		return nil, fmt.Errorf("email: %w", ErrDisabled)
	}
	return m.manager, nil
}

func (m *EmailModule) reloadIfChanged(ctx context.Context) error {
	raw, err := m.cont.GetRuntimeConfig().GetRaw(ctx, runtimeconfig.CodeEmail)
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

func (m *EmailModule) reload(ctx context.Context) error {
	raw, err := m.cont.GetRuntimeConfig().GetRaw(ctx, runtimeconfig.CodeEmail)
	if err != nil {
		return err
	}
	return m.apply(raw, sha256.Sum256(raw))
}

func (m *EmailModule) apply(raw []byte, digest [sha256.Size]byte) error {
	var cfg runtimeconfig.EmailConfig
	if err := json.Unmarshal(raw, &cfg); err != nil {
		return fmt.Errorf("decode email configuration: %w", err)
	}
	var manager *thirdpartyemail.Manager
	if cfg.Enabled {
		created, err := thirdpartyemail.NewManager(thirdpartyemail.Config{
			Host: cfg.Host, Port: cfg.Port, Username: cfg.Username,
			Password: cfg.Password, From: cfg.From,
		}, m.cont.GetRedis(), m.cont.GetLogger())
		if err != nil {
			return err
		}
		manager = created
	}
	m.mu.Lock()
	m.digest, m.enabled, m.manager = digest, cfg.Enabled, manager
	m.mu.Unlock()
	return nil
}
