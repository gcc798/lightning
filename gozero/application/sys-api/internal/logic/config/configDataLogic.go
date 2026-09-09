package config

import (
	"context"

	"github.com/gcc798/lightning/application/sys-api/internal/logic/commonutil"
	"github.com/gcc798/lightning/application/sys-api/internal/svc"
	"github.com/gcc798/lightning/application/sys-api/internal/types"
	"github.com/gcc798/lightning/application/sys-rpc/client/sysservice"
	"github.com/zeromicro/go-zero/core/logx"
)

type ConfigDataLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewConfigDataLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ConfigDataLogic {
	return &ConfigDataLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}
func (l *ConfigDataLogic) ConfigData(req *types.ConfigCodeQueryReq) (resp *types.CommonResp, err error) {
	data, err := l.svcCtx.SysRpcClient.ConfigData(l.ctx, &sysservice.ConfigCodeQueryReq{Code: req.Code})
	if err != nil {
		return &types.CommonResp{Code: 500, Msg: err.Error()}, nil
	}
	return &types.CommonResp{Code: 200, Msg: "操作成功", Data: map[string]interface{}{
		"code": req.Code,
		"data": commonutil.JSONStringToValue(data.DataJson),
	}}, nil
}
