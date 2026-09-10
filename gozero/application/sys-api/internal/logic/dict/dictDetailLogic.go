package dict

import (
	"context"

	"github.com/gcc798/microservice-kit/application/sys-api/internal/svc"
	"github.com/gcc798/microservice-kit/application/sys-api/internal/types"
	"github.com/gcc798/microservice-kit/application/sys-rpc/client/sysservice"
	"github.com/zeromicro/go-zero/core/logx"
)

type DictDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDictDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DictDetailLogic {
	return &DictDetailLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}
func (l *DictDetailLogic) DictDetail(req *types.IdPathReq) (resp *types.CommonResp, err error) {
	row, err := l.svcCtx.SysRpcClient.DictDetail(l.ctx, &sysservice.IdReq{Id: req.Id})
	if err != nil {
		return &types.CommonResp{Code: 500, Msg: err.Error()}, nil
	}
	return &types.CommonResp{Code: 200, Msg: "操作成功", Data: row}, nil
}
