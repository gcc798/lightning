// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package captcha

import (
	"context"

	"github.com/gcc798/lightning/application/sys-api/internal/logic/commonutil"
	"github.com/gcc798/lightning/application/sys-api/internal/svc"
	"github.com/gcc798/lightning/application/sys-api/internal/types"
	"github.com/gcc798/lightning/application/sys-rpc/client/sysservice"

	"github.com/zeromicro/go-zero/core/logx"
)

type CaptchaSmsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCaptchaSmsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CaptchaSmsLogic {
	return &CaptchaSmsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CaptchaSmsLogic) CaptchaSms(req *types.CaptchaPhoneReq) (resp *types.CommonResp, err error) {
	data, err := l.svcCtx.SysRpcClient.CaptchaSms(l.ctx, &sysservice.CaptchaPhoneReq{Phonenumber: req.Phone})
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
