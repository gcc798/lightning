package config

import (
	"context"

	"github.com/gcc798/microservice-kit/application/sys-api/internal/logic/commonutil"
	"github.com/gcc798/microservice-kit/application/sys-api/internal/svc"
	"github.com/gcc798/microservice-kit/application/sys-api/internal/types"
	"github.com/gcc798/microservice-kit/application/sys-rpc/client/sysservice"
	"github.com/zeromicro/go-zero/core/logx"
)

type ConfigUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewConfigUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ConfigUpdateLogic {
	return &ConfigUpdateLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}
func (l *ConfigUpdateLogic) ConfigUpdate(req *types.ConfigUpdateReq) (resp *types.CommonResp, err error) {
	data, err := commonutil.InterfaceToJSONString(req.Data)
	if err != nil {
		return &types.CommonResp{Code: 400, Msg: "data 格式错误"}, nil
	}
	if _, err := l.svcCtx.SysRpcClient.ConfigUpdate(l.ctx, &sysservice.ConfigUpdateReq{
		Id:       req.Id,
		Name:     req.Name,
		Code:     req.Code,
		DataJson: data,
		Remark:   req.Remark,
		UpdateBy: req.UpdateBy,
	}); err != nil {
		return &types.CommonResp{Code: 500, Msg: err.Error()}, nil
	}
	return &types.CommonResp{Code: 200, Msg: "操作成功", Data: "ok"}, nil
}
