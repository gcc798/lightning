package sysservicelogic

import (
	"context"

	"github.com/gcc798/lightning/application/sys-rpc/internal/svc"
	"github.com/gcc798/lightning/application/sys-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type StorageEnvDefaultLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewStorageEnvDefaultLogic(ctx context.Context, svcCtx *svc.ServiceContext) *StorageEnvDefaultLogic {
	return &StorageEnvDefaultLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *StorageEnvDefaultLogic) StorageEnvDefault(in *pb.Empty) (*pb.StorageEnv, error) {
	var row storageEnvRow
	if err := l.svcCtx.DB.QueryRowCtx(l.ctx, &row, `select id, name, code, storage_type, is_default, status, config, remark, create_by, created_time, update_by, updated_time from public.s_storage_env where is_default = true and status = 0 order by id desc limit 1`); err != nil {
		return nil, err
	}
	return toStorageEnvPB(row), nil
}
