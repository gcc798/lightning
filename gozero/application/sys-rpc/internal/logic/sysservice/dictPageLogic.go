package sysservicelogic

import (
	"context"
	"fmt"
	"strings"

	"github.com/gcc798/microservice-kit/application/sys-rpc/internal/svc"
	"github.com/gcc798/microservice-kit/application/sys-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type DictPageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDictPageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DictPageLogic {
	return &DictPageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DictPageLogic) DictPage(in *pb.DictPageReq) (*pb.DictPageResp, error) {
	pageNum, pageSize := normalizePage(in.PageNum, in.PageSize)
	where := []string{"(parent_id = 0 or parent_id is null)"}
	args := make([]interface{}, 0)
	if in.DictType != "" {
		args = append(args, "%"+in.DictType+"%")
		where = append(where, fmt.Sprintf("dict_type like $%d", len(args)))
	}
	if in.DictLabel != "" {
		args = append(args, "%"+in.DictLabel+"%")
		where = append(where, fmt.Sprintf("dict_label like $%d", len(args)))
	}
	if in.Status == 0 || in.Status == 1 {
		args = append(args, in.Status)
		where = append(where, fmt.Sprintf("status = $%d", len(args)))
	}
	whereSQL := strings.Join(where, " and ")
	var total int64
	if err := l.svcCtx.DB.QueryRowCtx(l.ctx, &total, "select count(1) from public.s_dict_data where "+whereSQL, args...); err != nil {
		return nil, err
	}
	queryArgs := append(append([]interface{}{}, args...), pageSize, (pageNum-1)*pageSize)
	var rows []dictRow
	query := `
		select id, parent_id, dict_type, dict_label, dict_value, sort, is_default, status, remark, create_by, update_by, created_time, updated_time
		from public.s_dict_data
		where ` + whereSQL + `
		order by sort asc, id asc
		limit $` + fmt.Sprint(len(args)+1) + ` offset $` + fmt.Sprint(len(args)+2)
	if err := l.svcCtx.DB.QueryRowsCtx(l.ctx, &rows, query, queryArgs...); err != nil {
		return nil, err
	}
	return &pb.DictPageResp{Records: toDictList(rows), Page: toPageInfo(total, pageNum, pageSize)}, nil
}
