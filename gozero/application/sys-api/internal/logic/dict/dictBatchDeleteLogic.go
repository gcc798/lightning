package dict

import (
	"context"

	"github.com/gcc798/microservice-kit/application/sys-api/internal/svc"
	"github.com/gcc798/microservice-kit/application/sys-api/internal/types"
	"github.com/gcc798/microservice-kit/application/sys-rpc/client/sysservice"
	"github.com/zeromicro/go-zero/core/logx"
)

type DictBatchDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDictBatchDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DictBatchDeleteLogic {
	return &DictBatchDeleteLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}
func (l *DictBatchDeleteLogic) DictBatchDelete(req *types.BatchIdsReq) (resp *types.CommonResp, err error) {
	if len(req.Ids) == 0 {
		return &types.CommonResp{Code: 400, Msg: "ids 不能为空"}, nil
	}
	if _, err := l.svcCtx.SysRpcClient.DictBatchDelete(l.ctx, &sysservice.BatchIdsReq{Ids: req.Ids}); err != nil {
		return &types.CommonResp{Code: 500, Msg: err.Error()}, nil
	}
	return &types.CommonResp{Code: 200, Msg: "操作成功", Data: "ok"}, nil
}
