package controller

import (
	"strings"

	iam "github.com/gcc798/lightning/application/iam/internal/domain"
	"github.com/gcc798/lightning/application/iam/internal/request"
	"github.com/gcc798/lightning/application/iam/internal/response"
	sysv1 "github.com/gcc798/lightning/internal/api/sys/v1"
	"github.com/gcc798/lightning/internal/config"
	"github.com/gcc798/lightning/internal/container"
	logging "github.com/gcc798/lightning/internal/logger"
	"github.com/gcc798/lightning/internal/modules"
	"github.com/labstack/echo/v5"
	"go.uber.org/zap"
)

type AuthController interface {
	Login(c *echo.Context)
	Logout(c *echo.Context)
	RefreshToken(c *echo.Context)
}

type authController struct {
	config  *config.Config
	logger  logging.Logger
	service iam.AuthService
}

func NewAuthController(c container.Container, logs sysv1.API) (AuthController, error) {
	captchaModule, err := modules.GetCaptcha(c)
	if err != nil {
		return nil, err
	}
	wechatModule, err := modules.GetWeChat(c)
	if err != nil {
		return nil, err
	}
	clientService := iam.NewClientService(c.GetDB(), c.GetRedis(), c.GetLogger())
	tokenManager := iam.NewTokenManager(c.GetJWT(), c.GetRedis(), c.GetLogger())
	captchaService := iam.NewCaptchaService(captchaModule)
	authService := iam.NewAuthService(
		c.GetDB(),
		c.GetRedis(),
		c.GetConfig(),
		c.GetLogger(),
		clientService,
		tokenManager,
		captchaService,
		wechatModule,
		logs,
	)
	return &authController{config: c.GetConfig(), logger: c.GetLogger(), service: authService}, nil
}

// Login godoc
//
//	@Summary		用户登录
//	@Description	支持密码、邮箱验证码、短信验证码和纯微信登录
//	@Tags			认证
//	@Accept			json
//	@Produce		json
//	@Param			request	body		request.LoginRequest	true	"登录请求"
//	@Success		200		{object}	response.Response{data=response.LoginResponse}
//	@Failure		400		{object}	response.Response
//	@Failure		401		{object}	response.Response
//	@Router			/login [post]
func (h *authController) Login(c *echo.Context) {
	var req request.LoginRequest
	if err := c.Bind(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}
	req.LoginIP = c.RealIP()
	req.UserAgent = c.Request().UserAgent()
	result, err := h.service.Login(c.Request().Context(), &req)
	if err != nil {
		h.logger.Warn("login failed", zap.String("grantType", req.GrantType), zap.Error(err))
		response.FailCode(c, response.CodeUnauthorized, err.Error())
		return
	}
	response.Success(c, result)
}

// Logout godoc
//
//	@Summary	用户登出
//	@Tags		认证
//	@Security	Bearer
//	@Success	200	{object}	response.Response{data=string}
//	@Router		/logout [post]
func (h *authController) Logout(c *echo.Context) {
	token := strings.TrimPrefix(c.Request().Header.Get(h.config.Auth.TokenHeader), "Bearer ")
	if err := h.service.Logout(c.Request().Context(), token); err != nil {
		h.logger.Warn("logout failed", zap.Error(err))
		response.InternalServerError(c, "登出失败")
		return
	}
	response.Success(c, "ok")
}

// RefreshToken godoc
//
//	@Summary		刷新访问令牌
//	@Description	Refresh Token 单次使用，刷新后立即轮换
//	@Tags			认证
//	@Accept			json
//	@Produce		json
//	@Param			request	body		response.RefreshTokenRequest	true	"刷新令牌请求"
//	@Success		200		{object}	response.Response{data=response.RefreshTokenResponse}
//	@Failure		400		{object}	response.Response
//	@Failure		401		{object}	response.Response
//	@Router			/auth/refresh [post]
func (h *authController) RefreshToken(c *echo.Context) {
	var req response.RefreshTokenRequest
	if err := c.Bind(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}
	result, err := h.service.Refresh(c.Request().Context(), &req)
	if err != nil {
		h.logger.Warn("refresh token failed", zap.String("clientId", req.ClientID), zap.Error(err))
		response.FailCode(c, response.CodeUnauthorized, err.Error())
		return
	}
	response.Success(c, result)
}
