package sysservicelogic

import (
	"context"

	"github.com/gcc798/lightning/application/sys-rpc/internal/svc"
	"github.com/gcc798/lightning/application/sys-rpc/pb"
	"github.com/zeromicro/go-zero/core/logx"
)

type MenuUserTreeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewMenuUserTreeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MenuUserTreeLogic {
	return &MenuUserTreeLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}
func (l *MenuUserTreeLogic) MenuUserTree(in *pb.IdReq) (*pb.MenuListResp, error) {
	rows, err := getUserMenus(l.ctx, l.svcCtx, in.Id)
	if err != nil {
		return nil, err
	}
	return &pb.MenuListResp{Records: buildMenuTree(rows, 0)}, nil
}
