package storageenv

import (
	"context"

	"github.com/gcc798/lightning/application/sys-api/internal/logic/commonutil"
	"github.com/gcc798/lightning/application/sys-api/internal/svc"
	"github.com/gcc798/lightning/application/sys-api/internal/types"
	"github.com/gcc798/lightning/application/sys-rpc/client/sysservice"
	"github.com/zeromicro/go-zero/core/logx"
)

type StorageEnvCreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewStorageEnvCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *StorageEnvCreateLogic {
	return &StorageEnvCreateLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}
func (l *StorageEnvCreateLogic) StorageEnvCreate(req *types.StorageEnvCreateReq) (resp *types.CommonResp, err error) {
	if req.Name == "" || req.Code == "" || req.StorageType == "" {
		return &types.CommonResp{Code: 400, Msg: "名称、编码和存储类型不能为空"}, nil
	}
	config, err := commonutil.InterfaceToJSONString(req.Config)
	if err != nil {
		return &types.CommonResp{Code: 400, Msg: "config 格式错误"}, nil
	}
	if _, err := l.svcCtx.SysRpcClient.StorageEnvCreate(l.ctx, &sysservice.StorageEnvReq{
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
