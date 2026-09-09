package sysservicelogic

import (
	"context"

	"github.com/gcc798/lightning/application/sys-rpc/internal/svc"
	"github.com/gcc798/lightning/application/sys-rpc/pb"
	"github.com/zeromicro/go-zero/core/logx"
)

type UserProfileLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserProfileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserProfileLogic {
	return &UserProfileLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}
func (l *UserProfileLogic) UserProfile(in *pb.IdReq) (*pb.User, error) {
	return getUserByID(l.ctx, l.svcCtx, in.Id)
}
