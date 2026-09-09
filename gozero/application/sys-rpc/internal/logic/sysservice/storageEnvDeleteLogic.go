package sysservicelogic

import (
	"context"
	"errors"

	"github.com/gcc798/lightning/application/sys-rpc/internal/svc"
	"github.com/gcc798/lightning/application/sys-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type StorageEnvDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewStorageEnvDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *StorageEnvDeleteLogic {
	return &StorageEnvDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *StorageEnvDeleteLogic) StorageEnvDelete(in *pb.IdReq) (*pb.Ack, error) {
	row, err := getStorageEnvByID(l.ctx, l.svcCtx, in.Id)
	if err != nil {
		return nil, err
	}
	if row.IsDefault {
		return nil, errors.New("默认存储环境不能删除")
	}
	var count int64
	if err := l.svcCtx.DB.QueryRowCtx(l.ctx, &count, `select count(1) from public.biz_attachment where env_id = $1 and status = '0'`, in.Id); err == nil && count > 0 {
		return nil, errors.New("该环境下仍有关联附件，无法删除")
	}
	if _, err := l.svcCtx.DB.ExecCtx(l.ctx, `delete from public.s_storage_env where id = $1`, in.Id); err != nil {
		return nil, err
	}
	return &pb.Ack{Msg: "ok"}, nil
}
