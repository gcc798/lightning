package sysservicelogic

import (
	"context"
	"database/sql"
	"errors"

	"github.com/gcc798/microservice-kit/application/sys-rpc/internal/svc"
	"github.com/gcc798/microservice-kit/application/sys-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type StorageEnvCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewStorageEnvCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *StorageEnvCreateLogic {
	return &StorageEnvCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *StorageEnvCreateLogic) StorageEnvCreate(in *pb.StorageEnvReq) (*pb.Ack, error) {
	if in.Name == "" || in.Code == "" || in.StorageType == "" {
		return nil, errors.New("名称、编码和存储类型不能为空")
	}
	exists, err := storageEnvCodeExists(l.ctx, l.svcCtx, in.Code, 0)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("环境编码已存在")
	}
	if in.IsDefault {
		_, _ = l.svcCtx.DB.ExecCtx(l.ctx, `update public.s_storage_env set is_default = false where is_default = true`)
	}
	if _, err := l.svcCtx.DB.ExecCtx(l.ctx, `insert into public.s_storage_env (name, code, storage_type, is_default, status, config, remark, created_time, updated_time) values ($1, $2, $3, $4, $5, $6, $7, now(), now())`,
		in.Name, in.Code, in.StorageType, in.IsDefault, in.Status, in.ConfigJson, sql.NullString{String: in.Remark, Valid: in.Remark != ""}); err != nil {
		return nil, err
	}
	return &pb.Ack{Msg: "ok"}, nil
}
