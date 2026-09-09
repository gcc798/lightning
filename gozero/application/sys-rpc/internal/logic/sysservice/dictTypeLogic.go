package sysservicelogic

import (
	"context"
	"fmt"

	"github.com/gcc798/lightning/application/sys-rpc/internal/svc"
	"github.com/gcc798/lightning/application/sys-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type DictTypeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDictTypeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DictTypeLogic {
	return &DictTypeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DictTypeLogic) DictType(in *pb.DictTypeQueryReq) (*pb.DictListResp, error) {
	args := []interface{}{in.DictType}
	query := `
		select id, parent_id, dict_type, dict_label, dict_value, sort, is_default, status, remark, create_by, update_by, created_time, updated_time
		from public.s_dict_data
		where dict_type = $1 and status = 0`
	if in.ParentId > 0 {
		args = append(args, in.ParentId)
		query += fmt.Sprintf(" and parent_id = $%d", len(args))
	}
	query += ` order by sort asc, id asc`
	var rows []dictRow
	if err := l.svcCtx.DB.QueryRowsCtx(l.ctx, &rows, query, args...); err != nil {
		return nil, err
	}
	return &pb.DictListResp{Records: toDictList(rows)}, nil
}
