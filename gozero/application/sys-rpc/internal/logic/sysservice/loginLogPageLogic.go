package sysservicelogic

import (
	"context"
	"fmt"
	"strings"

	"github.com/gcc798/lightning/application/sys-rpc/internal/svc"
	"github.com/gcc798/lightning/application/sys-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogPageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginLogPageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogPageLogic {
	return &LoginLogPageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LoginLogPageLogic) LoginLogPage(in *pb.LoginLogPageReq) (*pb.LoginLogPageResp, error) {
	pageNum, pageSize := normalizePage(in.PageNum, in.PageSize)
	where := []string{"1=1"}
	args := make([]interface{}, 0)
	if in.UserName != "" {
		args = append(args, "%"+in.UserName+"%")
		where = append(where, fmt.Sprintf("user_name like $%d", len(args)))
	}
	if in.Ipaddr != "" {
		args = append(args, "%"+in.Ipaddr+"%")
		where = append(where, fmt.Sprintf("ipaddr like $%d", len(args)))
	}
	if in.Status == 0 || in.Status == 1 {
		args = append(args, in.Status)
		where = append(where, fmt.Sprintf("status = $%d", len(args)))
	}
	if in.StartTime != "" {
		args = append(args, in.StartTime)
		where = append(where, fmt.Sprintf("login_time >= $%d", len(args)))
	}
	if in.EndTime != "" {
		args = append(args, in.EndTime)
		where = append(where, fmt.Sprintf("login_time <= $%d", len(args)))
	}
	whereSQL := strings.Join(where, " and ")
	var total int64
	if err := l.svcCtx.DB.QueryRowCtx(l.ctx, &total, "select count(1) from public.s_login_log where "+whereSQL, args...); err != nil {
		return nil, err
	}
	queryArgs := append(append([]interface{}{}, args...), pageSize, (pageNum-1)*pageSize)
	var rows []loginLogRow
	query := `select id, user_name, ipaddr, login_location, browser, os, status, msg, login_time, client_id from public.s_login_log where ` + whereSQL + ` order by login_time desc nulls last, id desc limit $` + fmt.Sprint(len(args)+1) + ` offset $` + fmt.Sprint(len(args)+2)
	if err := l.svcCtx.DB.QueryRowsCtx(l.ctx, &rows, query, queryArgs...); err != nil {
		return nil, err
	}
	return &pb.LoginLogPageResp{Records: toLoginLogList(rows), Page: toPageInfo(total, pageNum, pageSize)}, nil
}
