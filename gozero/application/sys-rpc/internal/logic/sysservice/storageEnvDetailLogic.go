package sysservicelogic

import (
	"context"

	"github.com/gcc798/microservice-kit/application/sys-rpc/internal/svc"
	"github.com/gcc798/microservice-kit/application/sys-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type StorageEnvDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewStorageEnvDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *StorageEnvDetailLogic {
	return &StorageEnvDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *StorageEnvDetailLogic) StorageEnvDetail(in *pb.IdReq) (*pb.StorageEnv, error) {
	row, err := getStorageEnvByID(l.ctx, l.svcCtx, in.Id)
	if err != nil {
		return nil, err
	}
	return toStorageEnvPB(*row), nil
}
