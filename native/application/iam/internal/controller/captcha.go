package controller

import (
	"context"

	iam "github.com/gcc798/lightning/application/iam/internal/domain"
	"github.com/gcc798/lightning/application/iam/internal/response"
	"github.com/gcc798/lightning/internal/platform/captcha"

	"github.com/labstack/echo/v5"
)

// CaptchaController 验证码控制器
type CaptchaController struct {
	captchaService iam.CaptchaService
	smsSender      interface {
		SendVerificationCode(ctx context.Context, phonenumber string) (string, error)
	}
}

// NewCaptchaController 创建验证码控制器
func NewCaptchaController(captchaService iam.CaptchaService, smsSender interface {
	SendVerificationCode(ctx context.Context, phonenumber string) (string, error)
}) *CaptchaController {
	return &CaptchaController{
		captchaService: captchaService,
		smsSender:      smsSender,
	}
}

// GenerateImageCaptcha 生成图形验证码
// @Summary 生成图形验证码
// @Tags 验证码
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=captcha.CaptchaData}
// @Router /captcha/image [get]
func (c *CaptchaController) GenerateImageCaptcha(ctx *echo.Context) {
	data, err := c.captchaService.Generate(ctx.Request().Context(), captcha.CaptchaTypeImage, "")
	if err != nil {
		response.Fail(ctx, err.Error())
		return
	}
	response.Success(ctx, data)
}

// SendSMSCaptchaRequest 发送短信验证码请求
type SendSMSCaptchaRequest struct {
	Phone string `json:"phone" binding:"required"`
}

// SendSMSCaptcha 发送短信验证码
// @Summary 发送短信验证码
// @Tags 验证码
// @Accept json
// @Produce json
// @Param request body SendSMSCaptchaRequest true "手机号"
// @Success 200 {object} response.Response{data=captcha.CaptchaData}
// @Router /captcha/sms [post]
func (c *CaptchaController) SendSMSCaptcha(ctx *echo.Context) {
	var req SendSMSCaptchaRequest
	if err := ctx.Bind(&req); err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}

	data, err := c.captchaService.Generate(ctx.Request().Context(), captcha.CaptchaTypeSMS, req.Phone)
	if err != nil {
		response.Fail(ctx, err.Error())
		return
	}
	response.Success(ctx, data)
}

// ResourceSMSCode 执行业务逻辑。
func (c *CaptchaController) ResourceSMSCode(ctx *echo.Context) {
	phonenumber := ctx.QueryParam("phonenumber")
	if phonenumber == "" {
		phonenumber = ctx.QueryParam("phone")
	}
	if phonenumber == "" {
		response.BadRequest(ctx, "手机号不能为空")
		return
	}
	if c.smsSender == nil {
		response.Fail(ctx, "短信服务未配置")
		return
	}
	if _, err := c.smsSender.SendVerificationCode(ctx.Request().Context(), phonenumber); err != nil {
		response.Fail(ctx, err.Error())
		return
	}
	response.Success(ctx, nil)
}

// SendEmailCaptchaRequest 发送邮箱验证码请求
type SendEmailCaptchaRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// SendEmailCaptcha 发送邮箱验证码
// @Summary 发送邮箱验证码
// @Tags 验证码
// @Accept json
// @Produce json
// @Param request body SendEmailCaptchaRequest true "邮箱"
// @Success 200 {object} response.Response{data=captcha.CaptchaData}
// @Router /captcha/email [post]
func (c *CaptchaController) SendEmailCaptcha(ctx *echo.Context) {
	var req SendEmailCaptchaRequest
	if err := ctx.Bind(&req); err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}

	data, err := c.captchaService.Generate(ctx.Request().Context(), captcha.CaptchaTypeEmail, req.Email)
	if err != nil {
		response.Fail(ctx, err.Error())
		return
	}
	response.Success(ctx, data)
}

// GetEnabledTypes 获取已启用的验证码类型
// @Summary 获取已启用的验证码类型
// @Tags 验证码
// @Produce json
// @Success 200 {object} response.Response{data=[]string}
// @Router /captcha/types [get]
func (c *CaptchaController) GetEnabledTypes(ctx *echo.Context) {
	types, err := c.captchaService.GetEnabledTypes(ctx.Request().Context())
	if err != nil {
		response.Fail(ctx, err.Error())
		return
	}
	response.Success(ctx, types)
}
