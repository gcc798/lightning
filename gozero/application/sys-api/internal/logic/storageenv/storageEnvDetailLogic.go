package storageenv

import (
	"context"

	"github.com/gcc798/lightning/application/sys-api/internal/svc"
	"github.com/gcc798/lightning/application/sys-api/internal/types"
	"github.com/gcc798/lightning/application/sys-rpc/client/sysservice"
	"github.com/zeromicro/go-zero/core/logx"
)

type StorageEnvDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewStorageEnvDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *StorageEnvDetailLogic {
	return &StorageEnvDetailLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}
func (l *StorageEnvDetailLogic) StorageEnvDetail(req *types.IdPathReq) (resp *types.CommonResp, err error) {
	row, err := l.svcCtx.SysRpcClient.StorageEnvDetail(l.ctx, &sysservice.IdReq{Id: req.Id})
	if err != nil {
		return &types.CommonResp{Code: 404, Msg: err.Error()}, nil
	}
	return &types.CommonResp{Code: 200, Msg: "操作成功", Data: row}, nil
}
