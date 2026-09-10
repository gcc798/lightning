package dict

import (
	"context"

	"github.com/gcc798/microservice-kit/application/sys-api/internal/svc"
	"github.com/gcc798/microservice-kit/application/sys-api/internal/types"
	"github.com/gcc798/microservice-kit/application/sys-rpc/client/sysservice"
	"github.com/zeromicro/go-zero/core/logx"
)

type DictLabelLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDictLabelLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DictLabelLogic {
	return &DictLabelLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}
func (l *DictLabelLogic) DictLabel(req *types.DictLabelQueryReq) (resp *types.CommonResp, err error) {
	data, err := l.svcCtx.SysRpcClient.DictLabel(l.ctx, &sysservice.DictLabelQueryReq{DictType: req.DictType, DictValue: req.DictValue})
	if err != nil {
		return &types.CommonResp{Code: 500, Msg: err.Error()}, nil
	}
	return &types.CommonResp{Code: 200, Msg: "操作成功", Data: data.Label}, nil
}
