package modules

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"sync"

	thirdpartysms "github.com/gcc798/microservice-kit/internal/platform/thirdparty/sms"
	"github.com/gcc798/microservice-kit/internal/runtimeconfig"
)

type SMSModule struct {
	mu      sync.RWMutex
	cont    Container
	digest  [sha256.Size]byte
	enabled bool
	manager *thirdpartysms.Manager
}

func NewSMSModule() *SMSModule  { return &SMSModule{} }
func (*SMSModule) Name() string { return SMSName }

func (m *SMSModule) Init(ctx context.Context, cont Container) error {
	m.cont = cont
	return m.reload(ctx)
}

func (*SMSModule) Start(context.Context) error { return nil }
func (*SMSModule) Stop(context.Context) error  { return nil }

func (m *SMSModule) Refresh(ctx context.Context, req ModuleRefreshRequest) error {
	if !containsCode(req.Codes, runtimeconfig.CodeSMS) {
		return nil
	}
	return m.reload(ctx)
}

func (m *SMSModule) SendVerificationCode(ctx context.Context, phonenumber string) (string, error) {
	manager, err := m.current(ctx)
	if err != nil {
		return "", err
	}
	return manager.SendVerificationCode(ctx, phonenumber)
}

func (m *SMSModule) SendSMS(ctx context.Context, phone, code, template string) error {
	manager, err := m.current(ctx)
	if err != nil {
		return err
	}
	return manager.Send(phone, template, map[string]string{"code": code})
}

func (m *SMSModule) current(ctx context.Context) (*thirdpartysms.Manager, error) {
	if err := m.reloadIfChanged(ctx); err != nil {
		return nil, err
	}
	m.mu.RLock()
	defer m.mu.RUnlock()
	if !m.enabled || m.manager == nil {
		return nil, fmt.Errorf("sms: %w", ErrDisabled)
	}
	return m.manager, nil
}

func (m *SMSModule) reloadIfChanged(ctx context.Context) error {
	raw, err := m.cont.GetRuntimeConfig().GetRaw(ctx, runtimeconfig.CodeSMS)
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

func (m *SMSModule) reload(ctx context.Context) error {
	raw, err := m.cont.GetRuntimeConfig().GetRaw(ctx, runtimeconfig.CodeSMS)
	if err != nil {
		return err
	}
	return m.apply(raw, sha256.Sum256(raw))
}

func (m *SMSModule) apply(raw []byte, digest [sha256.Size]byte) error {
	var cfg runtimeconfig.SMSConfig
	if err := json.Unmarshal(raw, &cfg); err != nil {
		return fmt.Errorf("decode SMS configuration: %w", err)
	}
	var manager *thirdpartysms.Manager
	if cfg.Enabled {
		created, err := thirdpartysms.NewManager(thirdpartysms.Config{
			AccessKeyId: cfg.AccessKeyID, AccessKeySecret: cfg.AccessKeySecret,
			SignName: cfg.SignName, TemplateCode: cfg.TemplateCode,
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
