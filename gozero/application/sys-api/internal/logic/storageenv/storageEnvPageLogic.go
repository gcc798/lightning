package storageenv

import (
	"context"

	"github.com/gcc798/microservice-kit/application/sys-api/internal/logic/commonutil"
	"github.com/gcc798/microservice-kit/application/sys-api/internal/svc"
	"github.com/gcc798/microservice-kit/application/sys-api/internal/types"
	"github.com/gcc798/microservice-kit/application/sys-rpc/client/sysservice"
	"github.com/zeromicro/go-zero/core/logx"
)

type StorageEnvPageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewStorageEnvPageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *StorageEnvPageLogic {
	return &StorageEnvPageLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}
func (l *StorageEnvPageLogic) StorageEnvPage(req *types.StorageEnvPageReq) (resp *types.CommonResp, err error) {
	data, err := l.svcCtx.SysRpcClient.StorageEnvPage(l.ctx, &sysservice.StorageEnvPageReq{
		PageNum:     req.PageNum,
		PageSize:    req.PageSize,
		Name:        req.Name,
		StorageType: req.StorageType,
	})
	if err != nil {
		return &types.CommonResp{Code: 500, Msg: err.Error()}, nil
	}
	return &types.CommonResp{Code: 200, Msg: "操作成功", Data: commonutil.PageData(data.Records, data.Page.Total, data.Page.Current, data.Page.Size)}, nil
}
