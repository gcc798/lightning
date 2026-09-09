package iam

import (
	"context"

	"github.com/gcc798/lightning/internal/platform/captcha"
)

// CaptchaService 验证码服务接口
type CaptchaService interface {
	// Generate 生成验证码
	Generate(ctx context.Context, captchaType captcha.CaptchaType, params interface{}) (*captcha.CaptchaData, error)

	// Verify 验证验证码
	Verify(ctx context.Context, captchaType captcha.CaptchaType, params interface{}) error

	// GetEnabledTypes 获取已启用的验证码类型
	GetEnabledTypes(ctx context.Context) ([]captcha.CaptchaType, error)

	ImageEnabled(ctx context.Context) (bool, error)
}

type CaptchaProvider interface {
	Generate(ctx context.Context, captchaType captcha.CaptchaType, params any) (*captcha.CaptchaData, error)
	Verify(ctx context.Context, captchaType captcha.CaptchaType, params any) error
	GetEnabledTypes(ctx context.Context) ([]captcha.CaptchaType, error)
	ImageEnabled(ctx context.Context) (bool, error)
}

// captchaService 验证码服务实现
type captchaService struct {
	provider CaptchaProvider
}

// NewCaptchaService 创建验证码服务
func NewCaptchaService(provider CaptchaProvider) CaptchaService {
	return &captchaService{
		provider: provider,
	}
}

// Generate 执行业务逻辑。
func (s *captchaService) Generate(ctx context.Context, captchaType captcha.CaptchaType, params interface{}) (*captcha.CaptchaData, error) {
	return s.provider.Generate(ctx, captchaType, params)
}

// Verify 执行业务逻辑。
func (s *captchaService) Verify(ctx context.Context, captchaType captcha.CaptchaType, params interface{}) error {
	return s.provider.Verify(ctx, captchaType, params)
}

// GetEnabledTypes 获取业务数据。
func (s *captchaService) GetEnabledTypes(ctx context.Context) ([]captcha.CaptchaType, error) {
	return s.provider.GetEnabledTypes(ctx)
}

func (s *captchaService) ImageEnabled(ctx context.Context) (bool, error) {
	return s.provider.ImageEnabled(ctx)
}
