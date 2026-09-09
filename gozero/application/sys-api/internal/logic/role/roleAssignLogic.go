package role

import (
	"context"

	"github.com/gcc798/lightning/application/sys-api/internal/svc"
	"github.com/gcc798/lightning/application/sys-api/internal/types"
	"github.com/gcc798/lightning/application/sys-rpc/client/sysservice"
	"github.com/zeromicro/go-zero/core/logx"
)

type RoleAssignLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRoleAssignLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RoleAssignLogic {
	return &RoleAssignLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}
func (l *RoleAssignLogic) RoleAssign(req *types.AssignRoleReq) (resp *types.CommonResp, err error) {
	if _, err := l.svcCtx.SysRpcClient.RoleAssign(l.ctx, &sysservice.AssignRoleReq{UserId: req.UserId, RoleId: req.RoleId}); err != nil {
		return &types.CommonResp{Code: 500, Msg: err.Error()}, nil
	}
	return &types.CommonResp{Code: 200, Msg: "操作成功", Data: "ok"}, nil
}
