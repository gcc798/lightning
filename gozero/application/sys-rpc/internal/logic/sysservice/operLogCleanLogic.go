package sysservicelogic

import (
	"context"
	"fmt"

	"github.com/gcc798/lightning/application/sys-rpc/internal/svc"
	"github.com/gcc798/lightning/application/sys-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type OperLogCleanLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewOperLogCleanLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OperLogCleanLogic {
	return &OperLogCleanLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *OperLogCleanLogic) OperLogClean(in *pb.LogCleanReq) (*pb.Ack, error) {
	if in.Days <= 0 {
		return nil, fmt.Errorf("天数必须大于0")
	}
	if _, err := l.svcCtx.DB.ExecCtx(l.ctx, `delete from public.s_oper_log where oper_time < now() - ($1 || ' day')::interval`, in.Days); err != nil {
		return nil, err
	}
	return &pb.Ack{Msg: "ok"}, nil
}
