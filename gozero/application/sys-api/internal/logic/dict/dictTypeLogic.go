package dict

import (
	"context"

	"github.com/gcc798/microservice-kit/application/sys-api/internal/svc"
	"github.com/gcc798/microservice-kit/application/sys-api/internal/types"
	"github.com/gcc798/microservice-kit/application/sys-rpc/client/sysservice"
	"github.com/zeromicro/go-zero/core/logx"
)

type DictTypeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDictTypeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DictTypeLogic {
	return &DictTypeLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}
func (l *DictTypeLogic) DictType(req *types.DictTypeQueryReq) (resp *types.CommonResp, err error) {
	data, err := l.svcCtx.SysRpcClient.DictType(l.ctx, &sysservice.DictTypeQueryReq{DictType: req.DictType, ParentId: req.ParentId})
	if err != nil {
		return &types.CommonResp{Code: 500, Msg: err.Error()}, nil
	}
	return &types.CommonResp{Code: 200, Msg: "操作成功", Data: data.Records}, nil
}
