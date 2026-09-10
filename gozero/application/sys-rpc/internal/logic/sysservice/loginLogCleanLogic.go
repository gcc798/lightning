package sysservicelogic

import (
	"context"
	"fmt"

	"github.com/gcc798/microservice-kit/application/sys-rpc/internal/svc"
	"github.com/gcc798/microservice-kit/application/sys-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogCleanLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginLogCleanLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogCleanLogic {
	return &LoginLogCleanLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LoginLogCleanLogic) LoginLogClean(in *pb.LogCleanReq) (*pb.Ack, error) {
	if in.Days <= 0 {
		return nil, fmt.Errorf("天数必须大于0")
	}
	if _, err := l.svcCtx.DB.ExecCtx(l.ctx, `delete from public.s_login_log where login_time < now() - ($1 || ' day')::interval`, in.Days); err != nil {
		return nil, err
	}
	return &pb.Ack{Msg: "ok"}, nil
}
