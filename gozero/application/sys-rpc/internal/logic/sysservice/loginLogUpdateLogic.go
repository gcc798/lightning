package sysservicelogic

import (
	"context"
	"database/sql"

	"github.com/gcc798/lightning/application/sys-rpc/internal/svc"
	"github.com/gcc798/lightning/application/sys-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginLogUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogUpdateLogic {
	return &LoginLogUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LoginLogUpdateLogic) LoginLogUpdate(in *pb.LoginLogUpdateReq) (*pb.Ack, error) {
	if _, err := l.svcCtx.DB.ExecCtx(l.ctx, `update public.s_login_log set user_name = $2, ipaddr = $3, login_location = $4, browser = $5, os = $6, status = $7, msg = $8, client_id = $9 where id = $1`,
		in.Id,
		sql.NullString{String: in.UserName, Valid: in.UserName != ""},
		sql.NullString{String: in.Ipaddr, Valid: in.Ipaddr != ""},
		sql.NullString{String: in.LoginLocation, Valid: in.LoginLocation != ""},
		sql.NullString{String: in.Browser, Valid: in.Browser != ""},
		sql.NullString{String: in.Os, Valid: in.Os != ""},
		in.Status,
		sql.NullString{String: in.Msg, Valid: in.Msg != ""},
		sql.NullString{String: in.ClientId, Valid: in.ClientId != ""},
	); err != nil {
		return nil, err
	}
	return &pb.Ack{Msg: "ok"}, nil
}
