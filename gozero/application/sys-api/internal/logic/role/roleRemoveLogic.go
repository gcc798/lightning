package role

import (
	"context"

	"github.com/gcc798/lightning/application/sys-api/internal/svc"
	"github.com/gcc798/lightning/application/sys-api/internal/types"
	"github.com/gcc798/lightning/application/sys-rpc/client/sysservice"
	"github.com/zeromicro/go-zero/core/logx"
)

type RoleRemoveLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRoleRemoveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RoleRemoveLogic {
	return &RoleRemoveLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}
func (l *RoleRemoveLogic) RoleRemove(req *types.RemoveRoleReq) (resp *types.CommonResp, err error) {
	if _, err := l.svcCtx.SysRpcClient.RoleRemove(l.ctx, &sysservice.RemoveRoleReq{UserId: req.UserId, RoleId: req.RoleId}); err != nil {
		return &types.CommonResp{Code: 500, Msg: err.Error()}, nil
	}
	return &types.CommonResp{Code: 200, Msg: "操作成功", Data: "ok"}, nil
}
