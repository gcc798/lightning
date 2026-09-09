package dict

import (
	"context"

	"github.com/gcc798/lightning/application/sys-api/internal/svc"
	"github.com/gcc798/lightning/application/sys-api/internal/types"
	"github.com/gcc798/lightning/application/sys-rpc/client/sysservice"
	"github.com/zeromicro/go-zero/core/logx"
)

type DictDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDictDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DictDeleteLogic {
	return &DictDeleteLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}
func (l *DictDeleteLogic) DictDelete(req *types.IdPathReq) (resp *types.CommonResp, err error) {
	if _, err := l.svcCtx.SysRpcClient.DictDelete(l.ctx, &sysservice.IdReq{Id: req.Id}); err != nil {
		return &types.CommonResp{Code: 500, Msg: err.Error()}, nil
	}
	return &types.CommonResp{Code: 200, Msg: "操作成功", Data: "ok"}, nil
}
