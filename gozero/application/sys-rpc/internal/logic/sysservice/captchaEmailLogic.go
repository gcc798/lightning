package sysservicelogic

import (
	"context"

	"github.com/gcc798/lightning/application/sys-rpc/internal/svc"
	"github.com/gcc798/lightning/application/sys-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type CaptchaEmailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCaptchaEmailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CaptchaEmailLogic {
	return &CaptchaEmailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CaptchaEmailLogic) CaptchaEmail(in *pb.CaptchaEmailReq) (*pb.CaptchaDataResp, error) {
	return generateEmailCaptcha(l.ctx, l.svcCtx, in.Email)
}
