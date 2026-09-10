// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package role

import (
	"context"

	"github.com/gcc798/microservice-kit/application/sys-api/internal/logic/commonutil"
	"github.com/gcc798/microservice-kit/application/sys-api/internal/svc"
	"github.com/gcc798/microservice-kit/application/sys-api/internal/types"
	"github.com/gcc798/microservice-kit/application/sys-rpc/client/sysservice"

	"github.com/zeromicro/go-zero/core/logx"
)

type RoleAssignUsersLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRoleAssignUsersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RoleAssignUsersLogic {
	return &RoleAssignUsersLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RoleAssignUsersLogic) RoleAssignUsers(req *types.RoleUsersReq) (resp *types.CommonResp, err error) {
	if _, err := l.svcCtx.SysRpcClient.RoleAssignUsers(l.ctx, &sysservice.RoleUsersReq{
		RoleId:     req.RoleId,
		UserIds:    req.UserIds,
		OperatorId: commonutil.UserIDFromContext(l.ctx),
	}); err != nil {
		return &types.CommonResp{Code: 500, Msg: err.Error()}, nil
	}
	return &types.CommonResp{Code: 200, Msg: "操作成功", Data: "ok"}, nil
}
