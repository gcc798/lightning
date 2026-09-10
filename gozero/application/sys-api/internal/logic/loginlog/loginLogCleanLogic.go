package loginlog

import (
	"context"

	"github.com/gcc798/microservice-kit/application/sys-api/internal/svc"
	"github.com/gcc798/microservice-kit/application/sys-api/internal/types"
	"github.com/gcc798/microservice-kit/application/sys-rpc/client/sysservice"
	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogCleanLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogCleanLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogCleanLogic {
	return &LoginLogCleanLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}
func (l *LoginLogCleanLogic) LoginLogClean(req *types.LogCleanReq) (resp *types.CommonResp, err error) {
	if _, err := l.svcCtx.SysRpcClient.LoginLogClean(l.ctx, &sysservice.LogCleanReq{Days: req.Days}); err != nil {
		return &types.CommonResp{Code: 500, Msg: err.Error()}, nil
	}
	return &types.CommonResp{Code: 200, Msg: "操作成功", Data: nil}, nil
}
