package sysservicelogic

import (
	"context"
	"database/sql"

	"github.com/gcc798/microservice-kit/application/sys-rpc/internal/svc"
	"github.com/gcc798/microservice-kit/application/sys-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type DictLabelLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDictLabelLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DictLabelLogic {
	return &DictLabelLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DictLabelLogic) DictLabel(in *pb.DictLabelQueryReq) (*pb.DictLabelResp, error) {
	var row struct {
		Label sql.NullString `db:"dict_label"`
	}
	if err := l.svcCtx.DB.QueryRowCtx(l.ctx, &row, `select dict_label from public.s_dict_data where dict_type = $1 and dict_value = $2 and status = 0 order by sort asc, id asc limit 1`, in.DictType, in.DictValue); err != nil {
		return nil, err
	}
	return &pb.DictLabelResp{Label: row.Label.String}, nil
}
