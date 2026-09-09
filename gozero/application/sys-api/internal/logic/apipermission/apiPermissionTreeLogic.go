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

type ApiPermissionTreeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewApiPermissionTreeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApiPermissionTreeLogic {
	return &ApiPermissionTreeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ApiPermissionTreeLogic) ApiPermissionTree() (resp *types.CommonResp, err error) {
	data, err := l.svcCtx.SysRpcClient.ApiPermissionTree(l.ctx, &sysservice.Empty{})
	if err != nil {
		return failure(err), nil
	}
	return success(data.Records), nil
}
