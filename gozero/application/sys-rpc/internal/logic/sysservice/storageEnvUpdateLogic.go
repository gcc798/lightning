package sysservicelogic

import (
	"context"
	"database/sql"
	"errors"

	"github.com/gcc798/lightning/application/sys-rpc/internal/svc"
	"github.com/gcc798/lightning/application/sys-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type StorageEnvUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewStorageEnvUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *StorageEnvUpdateLogic {
	return &StorageEnvUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *StorageEnvUpdateLogic) StorageEnvUpdate(in *pb.StorageEnvReq) (*pb.Ack, error) {
	if _, err := getStorageEnvByID(l.ctx, l.svcCtx, in.Id); err != nil {
		return nil, err
	}
	exists, err := storageEnvCodeExists(l.ctx, l.svcCtx, in.Code, in.Id)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("环境编码已存在")
	}
	if in.IsDefault {
		_, _ = l.svcCtx.DB.ExecCtx(l.ctx, `update public.s_storage_env set is_default = false where is_default = true and id <> $1`, in.Id)
	}
	if _, err := l.svcCtx.DB.ExecCtx(l.ctx, `update public.s_storage_env set name = $2, code = $3, storage_type = $4, is_default = $5, status = $6, config = $7, remark = $8, updated_time = now() where id = $1`,
		in.Id, in.Name, in.Code, in.StorageType, in.IsDefault, in.Status, in.ConfigJson, sql.NullString{String: in.Remark, Valid: in.Remark != ""}); err != nil {
		return nil, err
	}
	return &pb.Ack{Msg: "ok"}, nil
}
