package sysservicelogic

import (
	"context"

	"github.com/gcc798/lightning/application/sys-rpc/internal/svc"
	"github.com/gcc798/lightning/application/sys-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type OperLogDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewOperLogDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OperLogDetailLogic {
	return &OperLogDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *OperLogDetailLogic) OperLogDetail(in *pb.IdReq) (*pb.OperLog, error) {
	row, err := getOperLogByID(l.ctx, l.svcCtx, in.Id)
	if err != nil {
		return nil, err
	}
	return toOperLogPB(*row), nil
}
