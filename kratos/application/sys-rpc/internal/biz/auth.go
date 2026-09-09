package biz

import (
	"context"
	"strings"
	"time"

	v1 "github.com/gcc798/lightning/kratos/api/system/v1"
	"github.com/gcc798/lightning/kratos/application/sys-rpc/internal/data"
	"golang.org/x/crypto/bcrypt"
)

type AuthUsecase struct {
	res *data.Resources
}

func NewAuthUsecase(res *data.Resources) *AuthUsecase {
	return &AuthUsecase{res: res}
}

func (uc *AuthUsecase) Login(ctx context.Context, req *v1.LoginRequest) (*v1.LoginReply, error) {
	if err := uc.validateLoginCaptcha(ctx, req); err != nil {
		return nil, err
	}
	client, err := uc.res.AuthenticateClient(ctx, req.GetClientKey(), req.GetClientSecret(), req.GetGrantType())
	if err != nil {
		_ = uc.res.CreateLoginLogEntry(ctx, fallbackLoginAccount(req), currentClientIP(ctx), "", currentUserAgent(ctx), 1, err.Error())
		return nil, unauthorized(err.Error())
	}
	username := req.GetUsername()
	if username == "" {
		username = fallbackLoginAccount(req)
	}
	if username == "" || uc.res == nil {
		return nil, badRequest("登录参数错误")
	}

	var (
		user           *v1.UserItem
		hashedPassword string
	)
	if strings.TrimSpace(req.GetGrantType()) == "xcx" {
		user, err = uc.res.ResolveXcxUser(ctx, req.GetPhonenumber(), req.GetWxCode())
		if err != nil {
			_ = uc.res.CreateLoginLogEntry(ctx, username, currentClientIP(ctx), client.ClientID, currentUserAgent(ctx), 1, err.Error())
			return nil, unauthorized(err.Error())
		}
	} else {
		user, hashedPassword, err = uc.res.FindUserByAccount(ctx, username)
		if err != nil {
			return nil, err
		}
		if user == nil {
			_ = uc.res.CreateLoginLogEntry(ctx, username, currentClientIP(ctx), client.ClientID, currentUserAgent(ctx), 1, "用户不存在")
			return nil, unauthorized("用户不存在")
		}
	}
	if user.GetUserType() < 0 {
		return nil, unauthorized("用户状态异常")
	}
	if user.GetStatus() != 0 {
		_ = uc.res.CreateLoginLogEntry(ctx, username, currentClientIP(ctx), client.ClientID, currentUserAgent(ctx), 1, "用户已停用")
		return nil, unauthorized("用户已停用")
	}
	if req.GetPassword() != "" && bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(req.GetPassword())) != nil {
		_ = uc.res.CreateLoginLogEntry(ctx, username, currentClientIP(ctx), client.ClientID, currentUserAgent(ctx), 1, "用户名或密码错误")
		return nil, unauthorized("用户名或密码错误")
	}
	accessToken, refreshToken, err := uc.res.IssueSession(ctx, user.GetUserId(), client)
	if err != nil {
		return nil, err
	}
	_ = uc.res.UpdateUserLoginState(ctx, user.GetUserId(), currentClientIP(ctx))
	_ = uc.res.CreateLoginLogEntry(ctx, user.GetUserName(), currentClientIP(ctx), client.ClientID, currentUserAgent(ctx), 0, "登录成功")
	return &v1.LoginReply{
		AccessToken:      accessToken,
		RefreshToken:     refreshToken,
		ExpiresIn:        authClientAccessTTL(client),
		RefreshExpiresIn: authClientRefreshTTL(client),
		UserInfo: &v1.UserInfo{
			UserId:      user.GetUserId(),
			Username:    user.GetUserName(),
			Nickname:    user.GetNickName(),
			Phonenumber: user.GetPhonenumber(),
			Email:       user.GetEmail(),
			Avatar:      user.GetAvatar(),
			UserType:    user.GetUserType(),
		},
	}, nil
}

func (uc *AuthUsecase) Logout(ctx context.Context) (*v1.MessageReply, error) {
	if err := uc.res.RevokeSession(ctx, currentAccessToken(ctx), ""); err != nil {
		return nil, err
	}
	return okReply(), nil
}

func (uc *AuthUsecase) RefreshToken(ctx context.Context, req *v1.RefreshTokenRequest) (*v1.RefreshTokenReply, error) {
	client, err := uc.res.AuthenticateClient(ctx, req.GetClientKey(), req.GetClientSecret(), "refresh")
	if err != nil {
		return nil, unauthorized(err.Error())
	}
	_, accessToken, refreshToken, err := uc.res.RefreshSession(ctx, req.GetRefreshToken(), client)
	if err != nil {
		return nil, unauthorized(err.Error())
	}
	return &v1.RefreshTokenReply{
		AccessToken:      accessToken,
		RefreshToken:     refreshToken,
		ExpiresIn:        authClientAccessTTL(client),
		RefreshExpiresIn: authClientRefreshTTL(client),
	}, nil
}

func (uc *AuthUsecase) ValidateAccessToken(ctx context.Context, req *v1.ValidateAccessTokenRequest) (*v1.ValidateAccessTokenReply, error) {
	if strings.TrimSpace(req.GetAccessToken()) == "" {
		return &v1.ValidateAccessTokenReply{Valid: false}, nil
	}
	return uc.res.ValidateAccessToken(ctx, req.GetAccessToken())
}

func (uc *AuthUsecase) CheckPermission(ctx context.Context, req *v1.CheckPermissionRequest) (*v1.CheckPermissionReply, error) {
	allowed, err := uc.res.CheckPermission(ctx, req.GetUserId(), req.GetResource(), req.GetAction())
	if err != nil {
		return nil, err
	}
	return &v1.CheckPermissionReply{Allowed: allowed}, nil
}

func fallbackLoginAccount(req *v1.LoginRequest) string {
	for _, item := range []string{req.GetPhonenumber(), req.GetEmail(), req.GetWxCode()} {
		if strings.TrimSpace(item) != "" {
			return item
		}
	}
	return ""
}

func authClientAccessTTL(client *data.AuthClientInfo) int64 {
	if client == nil || client.ActiveTimeout <= 0 {
		return int64((30 * time.Minute).Seconds())
	}
	return client.ActiveTimeout
}

func authClientRefreshTTL(client *data.AuthClientInfo) int64 {
	if client == nil || client.Timeout <= 0 {
		return int64((7 * 24 * time.Hour).Seconds())
	}
	return client.Timeout
}

func (uc *AuthUsecase) validateLoginCaptcha(ctx context.Context, req *v1.LoginRequest) error {
	if uc.res == nil {
		return badRequest("登录参数错误")
	}
	switch strings.TrimSpace(req.GetGrantType()) {
	case "password":
		if strings.TrimSpace(req.GetUuid()) == "" || strings.TrimSpace(req.GetCode()) == "" {
			return badRequest("验证码不能为空")
		}
		return uc.res.VerifyImageCaptcha(ctx, req.GetUuid(), req.GetCode())
	case "email":
		if strings.TrimSpace(req.GetEmail()) == "" || strings.TrimSpace(req.GetCode()) == "" {
			return badRequest("邮箱和验证码不能为空")
		}
		return uc.res.VerifyEmailCaptcha(ctx, req.GetEmail(), req.GetCode())
	case "xcx":
		if strings.TrimSpace(req.GetPhonenumber()) == "" || strings.TrimSpace(req.GetCode()) == "" {
			return badRequest("手机号和验证码不能为空")
		}
		return uc.res.VerifySMSCaptcha(ctx, req.GetPhonenumber(), req.GetCode())
	default:
		return nil
	}
}
