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

type UserApiPermissionsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserApiPermissionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserApiPermissionsLogic {
	return &UserApiPermissionsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserApiPermissionsLogic) UserApiPermissions(req *types.UserApiPermissionsPathReq) (resp *types.CommonResp, err error) {
	data, err := l.svcCtx.SysRpcClient.UserApiPermissions(l.ctx, &sysservice.IdReq{Id: req.Id})
	if err != nil {
		return failure(err), nil
	}
	return success(data.PermissionIds), nil
}
