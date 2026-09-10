// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package apipermission

import (
	"context"

	"github.com/gcc798/microservice-kit/application/sys-api/internal/logic/commonutil"
	"github.com/gcc798/microservice-kit/application/sys-api/internal/svc"
	"github.com/gcc798/microservice-kit/application/sys-api/internal/types"
	"github.com/gcc798/microservice-kit/application/sys-rpc/client/sysservice"

	"github.com/zeromicro/go-zero/core/logx"
)

type RoleApiPermissionsAssignLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRoleApiPermissionsAssignLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RoleApiPermissionsAssignLogic {
	return &RoleApiPermissionsAssignLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RoleApiPermissionsAssignLogic) RoleApiPermissionsAssign(req *types.RoleApiPermissionsAssignReq) (resp *types.CommonResp, err error) {
	if _, err := l.svcCtx.SysRpcClient.RoleApiPermissionsAssign(l.ctx, &sysservice.RoleApiPermissionsReq{
		RoleId:        req.RoleId,
		PermissionIds: req.PermissionIds,
		OperatorId:    commonutil.UserIDFromContext(l.ctx),
	}); err != nil {
		return failure(err), nil
	}
	return success("ok"), nil
}
