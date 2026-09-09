package storageenv

import (
	"context"

	"github.com/gcc798/lightning/application/sys-api/internal/svc"
	"github.com/gcc798/lightning/application/sys-api/internal/types"
	"github.com/gcc798/lightning/application/sys-rpc/client/sysservice"
	"github.com/zeromicro/go-zero/core/logx"
)

type StorageEnvTestLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewStorageEnvTestLogic(ctx context.Context, svcCtx *svc.ServiceContext) *StorageEnvTestLogic {
	return &StorageEnvTestLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}
func (l *StorageEnvTestLogic) StorageEnvTest(req *types.StorageEnvTestReq) (resp *types.CommonResp, err error) {
	data, err := l.svcCtx.SysRpcClient.StorageEnvTest(l.ctx, &sysservice.StorageEnvTestReq{Id: req.Id})
	if err != nil {
		return &types.CommonResp{Code: 500, Msg: err.Error()}, nil
	}
	return &types.CommonResp{Code: 200, Msg: "操作成功", Data: data}, nil
}
