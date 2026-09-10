// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package captcha

import (
	"context"

	"github.com/gcc798/microservice-kit/application/sys-api/internal/logic/commonutil"
	"github.com/gcc798/microservice-kit/application/sys-api/internal/svc"
	"github.com/gcc798/microservice-kit/application/sys-api/internal/types"
	"github.com/gcc798/microservice-kit/application/sys-rpc/client/sysservice"

	"github.com/zeromicro/go-zero/core/logx"
)

type CaptchaEmailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCaptchaEmailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CaptchaEmailLogic {
	return &CaptchaEmailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CaptchaEmailLogic) CaptchaEmail(req *types.CaptchaEmailReq) (resp *types.CommonResp, err error) {
	data, err := l.svcCtx.SysRpcClient.CaptchaEmail(l.ctx, &sysservice.CaptchaEmailReq{Email: req.Email})
	if err != nil {
		return &types.CommonResp{Code: 500, Msg: err.Error()}, nil
	}
	return &types.CommonResp{Code: 200, Msg: "操作成功", Data: map[string]interface{}{
		"id":       data.Id,
		"type":     data.Type,
		"data":     commonutil.JSONStringToValue(data.DataJson),
		"expireAt": data.ExpireAt,
	}}, nil
}
