package sysservicelogic

import (
	"context"

	"github.com/gcc798/microservice-kit/application/sys-rpc/internal/svc"
	"github.com/gcc798/microservice-kit/application/sys-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type CaptchaSmsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCaptchaSmsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CaptchaSmsLogic {
	return &CaptchaSmsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CaptchaSmsLogic) CaptchaSms(in *pb.CaptchaPhoneReq) (*pb.CaptchaDataResp, error) {
	return generateSmsCaptcha(l.ctx, l.svcCtx, in.Phonenumber)
}
