package sysservicelogic

import (
	"context"
	"database/sql"
	"time"

	"github.com/gcc798/lightning/application/sys-rpc/internal/svc"
	"github.com/gcc798/lightning/application/sys-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginLogCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogCreateLogic {
	return &LoginLogCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LoginLogCreateLogic) LoginLogCreate(in *pb.LoginLogReq) (*pb.Ack, error) {
	if _, err := l.svcCtx.DB.ExecCtx(l.ctx, `insert into public.s_login_log (user_name, ipaddr, login_location, browser, os, status, msg, login_time, client_id) values ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		sql.NullString{String: in.UserName, Valid: in.UserName != ""},
		sql.NullString{String: in.Ipaddr, Valid: in.Ipaddr != ""},
		sql.NullString{String: in.LoginLocation, Valid: in.LoginLocation != ""},
		sql.NullString{String: in.Browser, Valid: in.Browser != ""},
		sql.NullString{String: in.Os, Valid: in.Os != ""},
		in.Status,
		sql.NullString{String: in.Msg, Valid: in.Msg != ""},
		time.Now(),
		sql.NullString{String: in.ClientId, Valid: in.ClientId != ""},
	); err != nil {
		return nil, err
	}
	return &pb.Ack{Msg: "ok"}, nil
}
