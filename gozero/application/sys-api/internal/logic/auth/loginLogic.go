package auth

import (
	"context"

	"github.com/gcc798/lightning/application/sys-api/internal/svc"
	"github.com/gcc798/lightning/application/sys-api/internal/types"
	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}
func (l *LoginLogic) Login(req *types.LoginReq) (resp *types.CommonResp, err error) {
	user, client, err := loginWithRPC(l.ctx, l.svcCtx, req)
	if err != nil {
		return &types.CommonResp{Code: 401, Msg: err.Error()}, nil
	}
	return buildLoginResponse(l.ctx, l.svcCtx, user, client)
}
