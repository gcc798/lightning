// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package captcha

import (
	"context"

	"github.com/gcc798/lightning/application/sys-api/internal/svc"
	"github.com/gcc798/lightning/application/sys-api/internal/types"
	"github.com/gcc798/lightning/application/sys-rpc/client/sysservice"

	"github.com/zeromicro/go-zero/core/logx"
)

type CaptchaEnabledTypesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCaptchaEnabledTypesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CaptchaEnabledTypesLogic {
	return &CaptchaEnabledTypesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CaptchaEnabledTypesLogic) CaptchaEnabledTypes() (resp *types.CommonResp, err error) {
	data, err := l.svcCtx.SysRpcClient.CaptchaEnabledTypes(l.ctx, &sysservice.CaptchaReq{})
	if err != nil {
		return &types.CommonResp{Code: 500, Msg: err.Error()}, nil
	}
	return &types.CommonResp{Code: 200, Msg: "操作成功", Data: data.Types}, nil
}
