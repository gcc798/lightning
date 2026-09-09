// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package apipermission

import (
	"context"

	"github.com/gcc798/lightning/application/sys-api/internal/svc"
	"github.com/gcc798/lightning/application/sys-api/internal/types"
	"github.com/gcc798/lightning/application/sys-rpc/client/sysservice"

	"github.com/zeromicro/go-zero/core/logx"
)

type RoleApiPermissionsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRoleApiPermissionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RoleApiPermissionsLogic {
	return &RoleApiPermissionsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RoleApiPermissionsLogic) RoleApiPermissions(req *types.RoleApiPermissionsPathReq) (resp *types.CommonResp, err error) {
	data, err := l.svcCtx.SysRpcClient.RoleApiPermissions(l.ctx, &sysservice.IdReq{Id: req.RoleId})
	if err != nil {
		return failure(err), nil
	}
	return success(data.PermissionIds), nil
}
