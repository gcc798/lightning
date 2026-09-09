package sysservicelogic

import (
	"context"

	"github.com/gcc798/lightning/application/sys-rpc/internal/svc"
	"github.com/gcc798/lightning/application/sys-rpc/pb"
	"github.com/zeromicro/go-zero/core/logx"
)

type MenuDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewMenuDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MenuDetailLogic {
	return &MenuDetailLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}
func (l *MenuDetailLogic) MenuDetail(in *pb.IdReq) (*pb.Menu, error) {
	return getMenuByID(l.ctx, l.svcCtx, in.Id)
}
