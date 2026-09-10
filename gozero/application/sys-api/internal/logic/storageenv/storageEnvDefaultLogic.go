package storageenv

import (
	"context"

	"github.com/gcc798/microservice-kit/application/sys-api/internal/svc"
	"github.com/gcc798/microservice-kit/application/sys-api/internal/types"
	"github.com/gcc798/microservice-kit/application/sys-rpc/client/sysservice"
	"github.com/zeromicro/go-zero/core/logx"
)

type StorageEnvDefaultLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewStorageEnvDefaultLogic(ctx context.Context, svcCtx *svc.ServiceContext) *StorageEnvDefaultLogic {
	return &StorageEnvDefaultLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}
func (l *StorageEnvDefaultLogic) StorageEnvDefault() (resp *types.CommonResp, err error) {
	row, err := l.svcCtx.SysRpcClient.StorageEnvDefault(l.ctx, &sysservice.Empty{})
	if err != nil {
		return &types.CommonResp{Code: 500, Msg: err.Error()}, nil
	}
	return &types.CommonResp{Code: 200, Msg: "操作成功", Data: row}, nil
}
