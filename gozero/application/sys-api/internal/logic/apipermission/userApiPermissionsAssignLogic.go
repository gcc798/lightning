// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package apipermission

import (
	"context"

	"github.com/gcc798/lightning/application/sys-api/internal/logic/commonutil"
	"github.com/gcc798/lightning/application/sys-api/internal/svc"
	"github.com/gcc798/lightning/application/sys-api/internal/types"
	"github.com/gcc798/lightning/application/sys-rpc/client/sysservice"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserApiPermissionsAssignLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserApiPermissionsAssignLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserApiPermissionsAssignLogic {
	return &UserApiPermissionsAssignLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserApiPermissionsAssignLogic) UserApiPermissionsAssign(req *types.UserApiPermissionsAssignReq) (resp *types.CommonResp, err error) {
	if _, err := l.svcCtx.SysRpcClient.UserApiPermissionsAssign(l.ctx, &sysservice.UserApiPermissionsReq{
		UserId:        req.Id,
		PermissionIds: req.PermissionIds,
		OperatorId:    commonutil.UserIDFromContext(l.ctx),
	}); err != nil {
		return failure(err), nil
	}
	return success("ok"), nil
}
