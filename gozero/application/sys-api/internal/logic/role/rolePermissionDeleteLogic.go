package role

import (
	"context"

	"github.com/gcc798/lightning/application/sys-api/internal/svc"
	"github.com/gcc798/lightning/application/sys-api/internal/types"
	"github.com/gcc798/lightning/application/sys-rpc/client/sysservice"
	"github.com/zeromicro/go-zero/core/logx"
)

type RolePermissionDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRolePermissionDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RolePermissionDeleteLogic {
	return &RolePermissionDeleteLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}
func (l *RolePermissionDeleteLogic) RolePermissionDelete(req *types.RolePermissionDeleteReq) (resp *types.CommonResp, err error) {
	if _, err := l.svcCtx.SysRpcClient.RolePermissionDelete(l.ctx, &sysservice.RolePermissionReq{RoleKey: req.RoleKey, Resource: req.Resource, Action: req.Action}); err != nil {
		return &types.CommonResp{Code: 500, Msg: err.Error()}, nil
	}
	return &types.CommonResp{Code: 200, Msg: "操作成功", Data: "ok"}, nil
}
