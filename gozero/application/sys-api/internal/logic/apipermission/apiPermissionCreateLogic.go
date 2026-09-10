// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package apipermission

import (
	"context"

	"github.com/gcc798/microservice-kit/application/sys-api/internal/svc"
	"github.com/gcc798/microservice-kit/application/sys-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ApiPermissionCreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewApiPermissionCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApiPermissionCreateLogic {
	return &ApiPermissionCreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ApiPermissionCreateLogic) ApiPermissionCreate(req *types.ApiPermissionSaveReq) (resp *types.CommonResp, err error) {
	data, err := l.svcCtx.SysRpcClient.ApiPermissionCreate(l.ctx, saveReq(l.ctx, req))
	if err != nil {
		return failure(err), nil
	}
	return success(data), nil
}
