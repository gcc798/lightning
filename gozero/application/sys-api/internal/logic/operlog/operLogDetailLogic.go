package operlog

import (
	"context"

	"github.com/gcc798/lightning/application/sys-api/internal/svc"
	"github.com/gcc798/lightning/application/sys-api/internal/types"
	"github.com/gcc798/lightning/application/sys-rpc/client/sysservice"
	"github.com/zeromicro/go-zero/core/logx"
)

type OperLogDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOperLogDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OperLogDetailLogic {
	return &OperLogDetailLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}
func (l *OperLogDetailLogic) OperLogDetail(req *types.IdPathReq) (resp *types.CommonResp, err error) {
	row, err := l.svcCtx.SysRpcClient.OperLogDetail(l.ctx, &sysservice.IdReq{Id: req.Id})
	if err != nil {
		return &types.CommonResp{Code: 404, Msg: err.Error()}, nil
	}
	return &types.CommonResp{Code: 200, Msg: "操作成功", Data: row}, nil
}
