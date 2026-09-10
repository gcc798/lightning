package sysservicelogic

import (
	"context"

	"github.com/gcc798/microservice-kit/application/sys-rpc/internal/svc"
	"github.com/gcc798/microservice-kit/application/sys-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ConfigCodeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewConfigCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ConfigCodeLogic {
	return &ConfigCodeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ConfigCodeLogic) ConfigCode(in *pb.ConfigCodeQueryReq) (*pb.ConfigListResp, error) {
	var rows []configRow
	if err := l.svcCtx.DB.QueryRowsCtx(l.ctx, &rows, `select id, name, code, data, remark, create_by, created_time, update_by, updated_time from public.s_config where code = $1 order by id asc`, in.Code); err != nil {
		return nil, err
	}
	return &pb.ConfigListResp{Records: toConfigList(rows)}, nil
}
