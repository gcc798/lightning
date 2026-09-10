package router

import (
	"fmt"

	"github.com/gcc798/microservice-kit/application/iam/internal/controller"
	iam "github.com/gcc798/microservice-kit/application/iam/internal/domain"
	"github.com/gcc798/microservice-kit/internal/httpx"
	"github.com/gcc798/microservice-kit/internal/modules"
)

// 注册验证码相关路由。
func registerCaptchaRoutes(r *httpx.Router, ctx *RouterContext) error {
	captchaModule, err := modules.GetCaptcha(ctx.Container)
	if err != nil {
		return fmt.Errorf("get captcha module: %w", err)
	}
	smsModule, err := modules.GetSMS(ctx.Container)
	if err != nil {
		return fmt.Errorf("get SMS module: %w", err)
	}
	captchaService := iam.NewCaptchaService(captchaModule)
	captchaController := controller.NewCaptchaController(captchaService, smsModule)

	// 公开路由（无需认证）
	r.GET("/resource/sms/code", captchaController.ResourceSMSCode)

	captcha := r.Group("/captcha")
	{
		captcha.GET("/image", captchaController.GenerateImageCaptcha)    // 生成图形验证码
		captcha.POST("/sms", captchaController.SendSMSCaptcha)           // 发送短信验证码
		captcha.POST("/email", captchaController.SendEmailCaptcha)       // 发送邮箱验证码
		captcha.GET("/enabled-types", captchaController.GetEnabledTypes) // 获取启用的验证码类型
	}
	return nil
}
