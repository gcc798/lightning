package config

import (
	"context"

	"github.com/gcc798/lightning/application/sys-api/internal/svc"
	"github.com/gcc798/lightning/application/sys-api/internal/types"
	"github.com/gcc798/lightning/application/sys-rpc/client/sysservice"
	"github.com/zeromicro/go-zero/core/logx"
)

type ConfigCodeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewConfigCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ConfigCodeLogic {
	return &ConfigCodeLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}
func (l *ConfigCodeLogic) ConfigCode(req *types.ConfigCodeQueryReq) (resp *types.CommonResp, err error) {
	data, err := l.svcCtx.SysRpcClient.ConfigCode(l.ctx, &sysservice.ConfigCodeQueryReq{Code: req.Code})
	if err != nil {
		return &types.CommonResp{Code: 500, Msg: err.Error()}, nil
	}
	return &types.CommonResp{Code: 200, Msg: "操作成功", Data: data.Records}, nil
}
