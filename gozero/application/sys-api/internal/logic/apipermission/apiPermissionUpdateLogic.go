// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package apipermission

import (
	"context"

	"github.com/gcc798/lightning/application/sys-api/internal/svc"
	"github.com/gcc798/lightning/application/sys-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ApiPermissionUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewApiPermissionUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApiPermissionUpdateLogic {
	return &ApiPermissionUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ApiPermissionUpdateLogic) ApiPermissionUpdate(req *types.ApiPermissionUpdateReq) (resp *types.CommonResp, err error) {
	if _, err := l.svcCtx.SysRpcClient.ApiPermissionUpdate(l.ctx, updateReq(l.ctx, req)); err != nil {
		return failure(err), nil
	}
	return success("ok"), nil
}
