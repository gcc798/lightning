package sysservicelogic

import (
	"context"

	"github.com/gcc798/lightning/application/sys-rpc/internal/svc"
	"github.com/gcc798/lightning/application/sys-rpc/pb"
	"github.com/zeromicro/go-zero/core/logx"
)

type MenuListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewMenuListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MenuListLogic {
	return &MenuListLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}
func (l *MenuListLogic) MenuList(in *pb.Empty) (*pb.MenuListResp, error) {
	rows, err := getAllMenus(l.ctx, l.svcCtx)
	if err != nil {
		return nil, err
	}
	return &pb.MenuListResp{Records: toMenuList(rows)}, nil
}
