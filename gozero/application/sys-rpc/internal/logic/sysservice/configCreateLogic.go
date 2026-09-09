package sysservicelogic

import (
	"context"
	"database/sql"
	"errors"

	"github.com/gcc798/lightning/application/sys-rpc/internal/svc"
	"github.com/gcc798/lightning/application/sys-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ConfigCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewConfigCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ConfigCreateLogic {
	return &ConfigCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ConfigCreateLogic) ConfigCreate(in *pb.ConfigCreateReq) (*pb.Ack, error) {
	if in.Name == "" || in.Code == "" {
		return nil, errors.New("名称和编码不能为空")
	}
	exists, err := configNameExists(l.ctx, l.svcCtx, in.Name, 0)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("配置名称已存在")
	}
	if _, err := l.svcCtx.DB.ExecCtx(l.ctx, `insert into public.s_config (name, code, data, remark, create_by, update_by, created_time, updated_time) values ($1, $2, $3, $4, nullif($5, 0), nullif($6, 0), now(), now())`,
		in.Name, in.Code, sql.NullString{String: in.DataJson, Valid: in.DataJson != ""}, sql.NullString{String: in.Remark, Valid: in.Remark != ""}, in.CreateBy, in.UpdateBy); err != nil {
		return nil, err
	}
	return &pb.Ack{Msg: "ok"}, nil
}
