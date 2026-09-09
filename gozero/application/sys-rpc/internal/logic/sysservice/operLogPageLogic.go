package sysservicelogic

import (
	"context"
	"fmt"
	"strings"

	"github.com/gcc798/lightning/application/sys-rpc/internal/svc"
	"github.com/gcc798/lightning/application/sys-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type OperLogPageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewOperLogPageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OperLogPageLogic {
	return &OperLogPageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *OperLogPageLogic) OperLogPage(in *pb.OperLogPageReq) (*pb.OperLogPageResp, error) {
	pageNum, pageSize := normalizePage(in.PageNum, in.PageSize)
	where := []string{"1=1"}
	args := make([]interface{}, 0)
	if in.Title != "" {
		args = append(args, "%"+in.Title+"%")
		where = append(where, fmt.Sprintf("title like $%d", len(args)))
	}
	if in.OperName != "" {
		args = append(args, "%"+in.OperName+"%")
		where = append(where, fmt.Sprintf("oper_name like $%d", len(args)))
	}
	if in.BusinessType != "" {
		args = append(args, in.BusinessType)
		where = append(where, fmt.Sprintf("business_type = $%d", len(args)))
	}
	if in.Status != "" {
		args = append(args, in.Status)
		where = append(where, fmt.Sprintf("status = $%d", len(args)))
	}
	if in.StartTime != "" {
		args = append(args, in.StartTime)
		where = append(where, fmt.Sprintf("oper_time >= $%d", len(args)))
	}
	if in.EndTime != "" {
		args = append(args, in.EndTime)
		where = append(where, fmt.Sprintf("oper_time <= $%d", len(args)))
	}
	whereSQL := strings.Join(where, " and ")
	var total int64
	if err := l.svcCtx.DB.QueryRowCtx(l.ctx, &total, "select count(1) from public.s_oper_log where "+whereSQL, args...); err != nil {
		return nil, err
	}
	queryArgs := append(append([]interface{}{}, args...), pageSize, (pageNum-1)*pageSize)
	var rows []operLogRow
	query := `select id, title, business_type, method, request_method, device_type, oper_name, oper_url, oper_ip, oper_location, oper_param, json_result, status, error_msg, oper_time, cost_time, user_agent from public.s_oper_log where ` + whereSQL + ` order by oper_time desc nulls last, id desc limit $` + fmt.Sprint(len(args)+1) + ` offset $` + fmt.Sprint(len(args)+2)
	if err := l.svcCtx.DB.QueryRowsCtx(l.ctx, &rows, query, queryArgs...); err != nil {
		return nil, err
	}
	return &pb.OperLogPageResp{Records: toOperLogList(rows), Page: toPageInfo(total, pageNum, pageSize)}, nil
}
