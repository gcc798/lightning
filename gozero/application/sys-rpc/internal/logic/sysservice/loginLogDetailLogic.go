package sysservicelogic

import (
	"context"

	"github.com/gcc798/lightning/application/sys-rpc/internal/svc"
	"github.com/gcc798/lightning/application/sys-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginLogDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogDetailLogic {
	return &LoginLogDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LoginLogDetailLogic) LoginLogDetail(in *pb.IdReq) (*pb.LoginLog, error) {
	row, err := getLoginLogByID(l.ctx, l.svcCtx, in.Id)
	if err != nil {
		return nil, err
	}
	return toLoginLogPB(*row), nil
}
