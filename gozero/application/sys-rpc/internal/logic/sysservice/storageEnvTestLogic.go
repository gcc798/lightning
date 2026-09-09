package sysservicelogic

import (
	"context"

	"github.com/gcc798/lightning/application/sys-rpc/internal/svc"
	"github.com/gcc798/lightning/application/sys-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type StorageEnvTestLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewStorageEnvTestLogic(ctx context.Context, svcCtx *svc.ServiceContext) *StorageEnvTestLogic {
	return &StorageEnvTestLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *StorageEnvTestLogic) StorageEnvTest(in *pb.StorageEnvTestReq) (*pb.StorageEnvTestResp, error) {
	if _, err := getStorageEnvByID(l.ctx, l.svcCtx, in.Id); err != nil {
		return nil, err
	}
	return &pb.StorageEnvTestResp{Connected: true, Msg: "ok"}, nil
}
