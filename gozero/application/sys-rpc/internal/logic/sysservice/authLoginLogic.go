package sysservicelogic

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/gcc798/lightning/application/sys-rpc/internal/svc"
	"github.com/gcc798/lightning/application/sys-rpc/pb"
	"github.com/zeromicro/go-zero/core/logx"
)

type AuthLoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAuthLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AuthLoginLogic {
	return &AuthLoginLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *AuthLoginLogic) AuthLogin(in *pb.AuthLoginReq) (*pb.AuthLoginResp, error) {
	client, err := authenticateClient(l.ctx, l.svcCtx, in.ClientKey, in.GrantType)
	if err != nil {
		return nil, err
	}
	var user *userAuthRow
	switch in.GrantType {
	case "password":
		if err := verifyImageCaptcha(l.ctx, l.svcCtx, in.Uuid, in.Code); err != nil {
			return nil, err
		}
		user, err = authenticatePassword(l.ctx, l.svcCtx, in.Username, in.Password)
	case "email":
		user, err = authenticateEmail(l.ctx, l.svcCtx, in.Email, in.Uuid, in.Code)
	case "sms":
		user, err = authenticateSms(l.ctx, l.svcCtx, in.Phonenumber, in.Uuid, in.Code)
	case "xcx":
		user, err = authenticateXcx(l.ctx, l.svcCtx, in.Phonenumber, in.Code, in.WxCode)
	case "wechat":
		user, err = authenticateWechat(l.ctx, l.svcCtx, in.WxCode)
	default:
		return nil, errUnsupportedGrantType()
	}
	if err != nil {
		return nil, err
	}
	// Record login log
	l.svcCtx.DB.ExecCtx(l.ctx, `insert into public.s_login_log (user_name, status, msg, login_time, client_id) values ($1, $2, $3, $4, $5)`,
		sql.NullString{String: user.UserName, Valid: true},
		int64(0),
		sql.NullString{String: "登录成功", Valid: true},
		time.Now(),
		sql.NullString{String: client.ClientId, Valid: true},
	)
	return &pb.AuthLoginResp{
		ClientId:      client.ClientId,
		ClientKey:     client.ClientKey,
		DeviceType:    nullString(client.DeviceType),
		Timeout:       client.Timeout,
		ActiveTimeout: client.ActiveTimeout,
		UserInfo: &pb.UserInfo{
			UserId:      user.Id,
			Username:    user.UserName,
			Nickname:    nullString(user.NickName),
			Phonenumber: nullString(user.Phonenumber),
			Email:       nullString(user.Email),
			Avatar:      nullString(user.Avatar),
			UserType:    int32(user.UserType),
			OrgId:       nullInt64(user.OrgId),
		},
	}, nil
}

func errUnsupportedGrantType() error {
	return errors.New("不支持的授权类型，支持的类型: password, email, sms, xcx, wechat")
}
