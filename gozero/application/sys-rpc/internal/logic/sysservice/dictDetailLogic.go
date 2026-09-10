package sysservicelogic

import (
	"context"

	"github.com/gcc798/microservice-kit/application/sys-rpc/internal/svc"
	"github.com/gcc798/microservice-kit/application/sys-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type DictDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDictDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DictDetailLogic {
	return &DictDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DictDetailLogic) DictDetail(in *pb.IdReq) (*pb.Dict, error) {
	row, err := getDictByID(l.ctx, l.svcCtx, in.Id)
	if err != nil {
		return nil, err
	}
	return toDictPB(*row), nil
}
