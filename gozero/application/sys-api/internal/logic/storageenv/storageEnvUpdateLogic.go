package storageenv

import (
	"context"

	"github.com/gcc798/lightning/application/sys-api/internal/logic/commonutil"
	"github.com/gcc798/lightning/application/sys-api/internal/svc"
	"github.com/gcc798/lightning/application/sys-api/internal/types"
	"github.com/gcc798/lightning/application/sys-rpc/client/sysservice"
	"github.com/zeromicro/go-zero/core/logx"
)

type StorageEnvUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewStorageEnvUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *StorageEnvUpdateLogic {
	return &StorageEnvUpdateLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}
func (l *StorageEnvUpdateLogic) StorageEnvUpdate(req *types.StorageEnvUpdateReq) (resp *types.CommonResp, err error) {
	config, err := commonutil.InterfaceToJSONString(req.Config)
	if err != nil {
		return &types.CommonResp{Code: 400, Msg: "config 格式错误"}, nil
	}
	if _, err := l.svcCtx.SysRpcClient.StorageEnvUpdate(l.ctx, &sysservice.StorageEnvReq{
		Id:          req.Id,
		Name:        req.Name,
		Code:        req.Code,
		StorageType: req.StorageType,
		IsDefault:   req.IsDefault,
		Status:      int32(req.Status),
		ConfigJson:  config,
		Remark:      req.Remark,
	}); err != nil {
		return &types.CommonResp{Code: 500, Msg: err.Error()}, nil
	}
	return &types.CommonResp{Code: 200, Msg: "操作成功", Data: "ok"}, nil
}
