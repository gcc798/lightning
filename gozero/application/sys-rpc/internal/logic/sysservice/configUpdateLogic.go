package sysservicelogic

import (
	"context"
	"database/sql"
	"errors"

	"github.com/gcc798/lightning/application/sys-rpc/internal/svc"
	"github.com/gcc798/lightning/application/sys-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ConfigUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewConfigUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ConfigUpdateLogic {
	return &ConfigUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ConfigUpdateLogic) ConfigUpdate(in *pb.ConfigUpdateReq) (*pb.Ack, error) {
	oldRow, err := getConfigByID(l.ctx, l.svcCtx, in.Id)
	if err != nil {
		return nil, err
	}
	name := in.Name
	if name == "" {
		name = oldRow.Name
	}
	exists, err := configNameExists(l.ctx, l.svcCtx, name, in.Id)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("配置名称已存在")
	}
	code := in.Code
	if code == "" {
		code = oldRow.Code
	}
	if _, err := l.svcCtx.DB.ExecCtx(l.ctx, `update public.s_config set name = $2, code = $3, data = $4, remark = $5, update_by = nullif($6, 0), updated_time = now() where id = $1`,
		in.Id, name, code, sql.NullString{String: in.DataJson, Valid: in.DataJson != ""}, sql.NullString{String: in.Remark, Valid: in.Remark != ""}, in.UpdateBy); err != nil {
		return nil, err
	}
	return &pb.Ack{Msg: "ok"}, nil
}
