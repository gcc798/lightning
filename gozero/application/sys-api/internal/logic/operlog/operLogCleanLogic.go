package operlog

import (
	"context"

	"github.com/gcc798/microservice-kit/application/sys-api/internal/svc"
	"github.com/gcc798/microservice-kit/application/sys-api/internal/types"
	"github.com/gcc798/microservice-kit/application/sys-rpc/client/sysservice"
	"github.com/zeromicro/go-zero/core/logx"
)

type OperLogCleanLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOperLogCleanLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OperLogCleanLogic {
	return &OperLogCleanLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}
func (l *OperLogCleanLogic) OperLogClean(req *types.LogCleanReq) (resp *types.CommonResp, err error) {
	if _, err := l.svcCtx.SysRpcClient.OperLogClean(l.ctx, &sysservice.LogCleanReq{Days: req.Days}); err != nil {
		return &types.CommonResp{Code: 500, Msg: err.Error()}, nil
	}
	return &types.CommonResp{Code: 200, Msg: "操作成功", Data: nil}, nil
}
