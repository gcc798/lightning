package sysservicelogic

import (
	"context"

	"github.com/gcc798/lightning/application/sys-rpc/internal/svc"
	"github.com/gcc798/lightning/application/sys-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type OperLogDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewOperLogDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OperLogDeleteLogic {
	return &OperLogDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *OperLogDeleteLogic) OperLogDelete(in *pb.IdReq) (*pb.Ack, error) {
	if _, err := l.svcCtx.DB.ExecCtx(l.ctx, `delete from public.s_oper_log where id = $1`, in.Id); err != nil {
		return nil, err
	}
	return &pb.Ack{Msg: "ok"}, nil
}
