package role

import (
	"context"

	"github.com/gcc798/microservice-kit/application/sys-api/internal/svc"
	"github.com/gcc798/microservice-kit/application/sys-api/internal/types"
	"github.com/gcc798/microservice-kit/application/sys-rpc/client/sysservice"
	"github.com/zeromicro/go-zero/core/logx"
)

type RolePermissionAddLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRolePermissionAddLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RolePermissionAddLogic {
	return &RolePermissionAddLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}
func (l *RolePermissionAddLogic) RolePermissionAdd(req *types.RolePermissionReq) (resp *types.CommonResp, err error) {
	if _, err := l.svcCtx.SysRpcClient.RolePermissionAdd(l.ctx, &sysservice.RolePermissionReq{RoleKey: req.RoleKey, Resource: req.Resource, Action: req.Action}); err != nil {
		return &types.CommonResp{Code: 500, Msg: err.Error()}, nil
	}
	return &types.CommonResp{Code: 200, Msg: "操作成功", Data: "ok"}, nil
}
