package sysservicelogic

import (
	"context"

	"github.com/gcc798/lightning/application/sys-rpc/internal/svc"
	"github.com/gcc798/lightning/application/sys-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ConfigDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewConfigDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ConfigDetailLogic {
	return &ConfigDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ConfigDetailLogic) ConfigDetail(in *pb.IdReq) (*pb.Config, error) {
	row, err := getConfigByID(l.ctx, l.svcCtx, in.Id)
	if err != nil {
		return nil, err
	}
	return toConfigPB(*row), nil
}
