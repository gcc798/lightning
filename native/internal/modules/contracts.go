// Package modules contains runtime contracts and capabilities composed by application entrypoints.
package modules

import (
	"context"
	"errors"
	"fmt"

	"github.com/gcc798/lightning/internal/platform/captcha"
	"github.com/gcc798/lightning/internal/platform/thirdparty/wechat"
	"github.com/gcc798/lightning/internal/platform/websocket"
)

const (
	WeChatName    = "wechat"
	SMSName       = "sms"
	EmailName     = "email"
	CaptchaName   = "captcha"
	SchedulerName = "scheduler"
	WebSocketName = "websocket"
)

var ErrDisabled = errors.New("module is disabled")

type WeChat interface {
	Module
	Code2Session(ctx context.Context, wxCode string) (*wechat.Code2SessionResponse, error)
}

type SMS interface {
	Module
	SendVerificationCode(ctx context.Context, phonenumber string) (string, error)
	SendSMS(ctx context.Context, phone, code, template string) error
}

type Email interface {
	Module
	SendEmail(ctx context.Context, address, code, template string) error
}

type Captcha interface {
	Module
	Generate(ctx context.Context, captchaType captcha.CaptchaType, params any) (*captcha.CaptchaData, error)
	Verify(ctx context.Context, captchaType captcha.CaptchaType, params any) error
	GetEnabledTypes(ctx context.Context) ([]captcha.CaptchaType, error)
	ImageEnabled(ctx context.Context) (bool, error)
}

func GetWeChat(cont Container) (WeChat, error) {
	value, ok := cont.GetModule(WeChatName).(WeChat)
	if !ok {
		return nil, fmt.Errorf("module %q is not registered as a WeChat capability", WeChatName)
	}
	return value, nil
}

func GetSMS(cont Container) (SMS, error) {
	value, ok := cont.GetModule(SMSName).(SMS)
	if !ok {
		return nil, fmt.Errorf("module %q is not registered as an SMS capability", SMSName)
	}
	return value, nil
}

func GetEmail(cont Container) (Email, error) {
	value, ok := cont.GetModule(EmailName).(Email)
	if !ok {
		return nil, fmt.Errorf("module %q is not registered as an email capability", EmailName)
	}
	return value, nil
}

func GetCaptcha(cont Container) (Captcha, error) {
	value, ok := cont.GetModule(CaptchaName).(Captcha)
	if !ok {
		return nil, fmt.Errorf("module %q is not registered as a captcha capability", CaptchaName)
	}
	return value, nil
}

// WebSocketHub 返回已注册的连接管理中心，模块未注册或被禁用时返回 nil。
func WebSocketHub(cont Container) *websocket.Hub {
	value, ok := cont.GetModule(WebSocketName).(*WebSocketModule)
	if !ok {
		return nil
	}
	return value.Hub()
}
