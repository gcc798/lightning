package sysservicelogic

import (
	"context"

	"github.com/gcc798/lightning/application/sys-rpc/internal/svc"
	"github.com/gcc798/lightning/application/sys-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type CaptchaImageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCaptchaImageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CaptchaImageLogic {
	return &CaptchaImageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CaptchaImageLogic) CaptchaImage(in *pb.CaptchaReq) (*pb.CaptchaDataResp, error) {
	return generateImageCaptcha(l.ctx, l.svcCtx)
}
