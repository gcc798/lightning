package modules

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"sync"

	"github.com/gcc798/microservice-kit/internal/platform/captcha"
	"github.com/gcc798/microservice-kit/internal/runtimeconfig"
)

type CaptchaModule struct {
	mu      sync.RWMutex
	cont    Container
	digest  [sha256.Size]byte
	config  runtimeconfig.CaptchaConfig
	manager *captcha.CaptchaManager
	sms     SMS
	email   Email
}

func NewCaptchaModule() *CaptchaModule { return &CaptchaModule{} }
func (*CaptchaModule) Name() string    { return CaptchaName }

func (m *CaptchaModule) Init(ctx context.Context, cont Container) error {
	m.cont = cont
	var err error
	m.sms, err = GetSMS(cont)
	if err != nil {
		return err
	}
	m.email, err = GetEmail(cont)
	if err != nil {
		return err
	}
	return m.reload(ctx)
}

func (*CaptchaModule) Start(context.Context) error { return nil }
func (*CaptchaModule) Stop(context.Context) error  { return nil }

func (m *CaptchaModule) Refresh(ctx context.Context, req ModuleRefreshRequest) error {
	if !containsCode(req.Codes, runtimeconfig.CodeCaptcha) {
		return nil
	}
	return m.reload(ctx)
}

func (m *CaptchaModule) Generate(ctx context.Context, captchaType captcha.CaptchaType, params any) (*captcha.CaptchaData, error) {
	manager, err := m.current(ctx)
	if err != nil {
		return nil, err
	}
	return manager.Generate(ctx, captchaType, params)
}

func (m *CaptchaModule) Verify(ctx context.Context, captchaType captcha.CaptchaType, params any) error {
	manager, err := m.current(ctx)
	if err != nil {
		return err
	}
	return manager.Verify(ctx, captchaType, params)
}

func (m *CaptchaModule) GetEnabledTypes(ctx context.Context) ([]captcha.CaptchaType, error) {
	manager, err := m.current(ctx)
	if err != nil {
		return nil, err
	}
	return manager.GetEnabledTypes(), nil
}

func (m *CaptchaModule) ImageEnabled(ctx context.Context) (bool, error) {
	if _, err := m.current(ctx); err != nil {
		return false, err
	}
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.config.Image.Enabled, nil
}

func (m *CaptchaModule) current(ctx context.Context) (*captcha.CaptchaManager, error) {
	if err := m.reloadIfChanged(ctx); err != nil {
		return nil, err
	}
	m.mu.RLock()
	defer m.mu.RUnlock()
	if m.manager == nil {
		return nil, errors.New("captcha module is not initialized")
	}
	return m.manager, nil
}

func (m *CaptchaModule) reloadIfChanged(ctx context.Context) error {
	raw, err := m.cont.GetRuntimeConfig().GetRaw(ctx, runtimeconfig.CodeCaptcha)
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

func (m *CaptchaModule) reload(ctx context.Context) error {
	raw, err := m.cont.GetRuntimeConfig().GetRaw(ctx, runtimeconfig.CodeCaptcha)
	if err != nil {
		return err
	}
	return m.apply(raw, sha256.Sum256(raw))
}

func (m *CaptchaModule) apply(raw []byte, digest [sha256.Size]byte) error {
	var cfg runtimeconfig.CaptchaConfig
	if err := json.Unmarshal(raw, &cfg); err != nil {
		return fmt.Errorf("decode captcha configuration: %w", err)
	}
	manager := captcha.NewCaptchaManager()
	if cfg.Image.Enabled {
		manager.RegisterProvider(captcha.NewImageCaptchaProvider(&captcha.ImageCaptchaConfig{
			Enabled: true, Length: cfg.Image.Length, Width: cfg.Image.Width,
			Height: cfg.Image.Height, Expire: cfg.Image.Expire,
		}, m.cont.GetRedis()))
	}
	if cfg.SMS.Enabled {
		manager.RegisterProvider(captcha.NewSMSCaptchaProvider(&captcha.SMSCaptchaConfig{
			Enabled: true, Length: cfg.SMS.Length, Expire: cfg.SMS.Expire,
			Template: cfg.SMS.Template, Provider: cfg.SMS.Provider,
		}, m.cont.GetRedis(), m.sms))
	}
	if cfg.Email.Enabled {
		manager.RegisterProvider(captcha.NewEmailCaptchaProvider(&captcha.EmailCaptchaConfig{
			Enabled: true, Length: cfg.Email.Length, Expire: cfg.Email.Expire,
			Template: cfg.Email.Template,
		}, m.cont.GetRedis(), m.email))
	}
	m.mu.Lock()
	m.digest, m.config, m.manager = digest, cfg, manager
	m.mu.Unlock()
	return nil
}
