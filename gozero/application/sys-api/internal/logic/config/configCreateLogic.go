package config

import (
	"context"

	"github.com/gcc798/lightning/application/sys-api/internal/logic/commonutil"
	"github.com/gcc798/lightning/application/sys-api/internal/svc"
	"github.com/gcc798/lightning/application/sys-api/internal/types"
	"github.com/gcc798/lightning/application/sys-rpc/client/sysservice"
	"github.com/zeromicro/go-zero/core/logx"
)

type ConfigCreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewConfigCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ConfigCreateLogic {
	return &ConfigCreateLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}
func (l *ConfigCreateLogic) ConfigCreate(req *types.ConfigCreateReq) (resp *types.CommonResp, err error) {
	if req.Name == "" || req.Code == "" {
		return &types.CommonResp{Code: 400, Msg: "名称和编码不能为空"}, nil
	}
	data, err := commonutil.InterfaceToJSONString(req.Data)
	if err != nil {
		return &types.CommonResp{Code: 400, Msg: "data 格式错误"}, nil
	}
	if _, err := l.svcCtx.SysRpcClient.ConfigCreate(l.ctx, &sysservice.ConfigCreateReq{
		Name:     req.Name,
		Code:     req.Code,
		DataJson: data,
		Remark:   req.Remark,
		CreateBy: req.CreateBy,
		UpdateBy: req.UpdateBy,
	}); err != nil {
		return &types.CommonResp{Code: 500, Msg: err.Error()}, nil
	}
	return &types.CommonResp{Code: 200, Msg: "操作成功", Data: "ok"}, nil
}
