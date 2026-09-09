package sysservicelogic

import (
	"context"

	"github.com/gcc798/lightning/application/sys-rpc/internal/svc"
	"github.com/gcc798/lightning/application/sys-rpc/pb"
	"github.com/zeromicro/go-zero/core/logx"
)

type UserDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserDeleteLogic {
	return &UserDeleteLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}
func (l *UserDeleteLogic) UserDelete(in *pb.IdReq) (*pb.Ack, error) {
	if _, err := l.svcCtx.DB.ExecCtx(l.ctx, `delete from public.s_user where id = $1`, in.Id); err != nil {
		return nil, err
	}
	return &pb.Ack{Msg: "ok"}, nil
}
