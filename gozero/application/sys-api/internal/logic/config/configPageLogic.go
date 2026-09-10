package config

import (
	"context"

	"github.com/gcc798/microservice-kit/application/sys-api/internal/logic/commonutil"
	"github.com/gcc798/microservice-kit/application/sys-api/internal/svc"
	"github.com/gcc798/microservice-kit/application/sys-api/internal/types"
	"github.com/gcc798/microservice-kit/application/sys-rpc/client/sysservice"
	"github.com/zeromicro/go-zero/core/logx"
)

type ConfigPageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewConfigPageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ConfigPageLogic {
	return &ConfigPageLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}
func (l *ConfigPageLogic) ConfigPage(req *types.ConfigPageReq) (resp *types.CommonResp, err error) {
	data, err := l.svcCtx.SysRpcClient.ConfigPage(l.ctx, &sysservice.ConfigPageReq{
		PageNum:  req.PageNum,
		PageSize: req.PageSize,
		Name:     req.Name,
		Code:     req.Code,
	})
	if err != nil {
		return &types.CommonResp{Code: 500, Msg: err.Error()}, nil
	}
	return &types.CommonResp{Code: 200, Msg: "操作成功", Data: commonutil.PageData(data.Records, data.Page.Total, data.Page.Current, data.Page.Size)}, nil
}
