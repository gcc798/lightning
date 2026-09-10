// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package role

import (
	"context"

	"github.com/gcc798/microservice-kit/application/sys-api/internal/svc"
	"github.com/gcc798/microservice-kit/application/sys-api/internal/types"
	"github.com/gcc798/microservice-kit/application/sys-rpc/client/sysservice"

	"github.com/zeromicro/go-zero/core/logx"
)

type RoleMenusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRoleMenusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RoleMenusLogic {
	return &RoleMenusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RoleMenusLogic) RoleMenus(req *types.RoleMenusPathReq) (resp *types.CommonResp, err error) {
	data, err := l.svcCtx.SysRpcClient.RoleMenus(l.ctx, &sysservice.RoleMenusReq{RoleId: req.RoleId})
	if err != nil {
		return &types.CommonResp{Code: 500, Msg: err.Error()}, nil
	}
	return &types.CommonResp{Code: 200, Msg: "操作成功", Data: data.MenuIds}, nil
}
