package sysservicelogic

import (
	"context"
	"database/sql"

	"github.com/gcc798/lightning/application/sys-rpc/internal/svc"
	"github.com/gcc798/lightning/application/sys-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ConfigDataLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewConfigDataLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ConfigDataLogic {
	return &ConfigDataLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ConfigDataLogic) ConfigData(in *pb.ConfigCodeQueryReq) (*pb.ConfigDataResp, error) {
	var row struct {
		Data sql.NullString `db:"data"`
	}
	if err := l.svcCtx.DB.QueryRowCtx(l.ctx, &row, `select data from public.s_config where code = $1 order by id desc limit 1`, in.Code); err != nil {
		return nil, err
	}
	return &pb.ConfigDataResp{Code: in.Code, DataJson: nullString(row.Data)}, nil
}
