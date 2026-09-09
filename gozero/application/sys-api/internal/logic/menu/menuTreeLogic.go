package menu

import (
	"context"

	"github.com/gcc798/lightning/application/sys-api/internal/svc"
	"github.com/gcc798/lightning/application/sys-api/internal/types"
	"github.com/gcc798/lightning/application/sys-rpc/client/sysservice"
	"github.com/zeromicro/go-zero/core/logx"
)

type MenuTreeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMenuTreeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MenuTreeLogic {
	return &MenuTreeLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}
func (l *MenuTreeLogic) MenuTree() (resp *types.CommonResp, err error) {
	data, err := l.svcCtx.SysRpcClient.MenuTree(l.ctx, &sysservice.Empty{})
	if err != nil {
		return &types.CommonResp{Code: 500, Msg: err.Error()}, nil
	}
	return &types.CommonResp{Code: 200, Msg: "操作成功", Data: toNativeMenuTreeList(data.Records)}, nil
}
