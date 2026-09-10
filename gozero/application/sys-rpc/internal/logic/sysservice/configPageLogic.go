package sysservicelogic

import (
	"context"
	"fmt"
	"strings"

	"github.com/gcc798/microservice-kit/application/sys-rpc/internal/svc"
	"github.com/gcc798/microservice-kit/application/sys-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ConfigPageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewConfigPageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ConfigPageLogic {
	return &ConfigPageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ConfigPageLogic) ConfigPage(in *pb.ConfigPageReq) (*pb.ConfigPageResp, error) {
	pageNum, pageSize := normalizePage(in.PageNum, in.PageSize)
	where := []string{"1=1"}
	args := make([]interface{}, 0)
	if in.Name != "" {
		args = append(args, "%"+in.Name+"%")
		where = append(where, fmt.Sprintf("name like $%d", len(args)))
	}
	if in.Code != "" {
		args = append(args, "%"+in.Code+"%")
		where = append(where, fmt.Sprintf("code like $%d", len(args)))
	}
	whereSQL := strings.Join(where, " and ")
	var total int64
	if err := l.svcCtx.DB.QueryRowCtx(l.ctx, &total, "select count(1) from public.s_config where "+whereSQL, args...); err != nil {
		return nil, err
	}
	queryArgs := append(append([]interface{}{}, args...), pageSize, (pageNum-1)*pageSize)
	var rows []configRow
	query := `select id, name, code, data, remark, create_by, created_time, update_by, updated_time from public.s_config where ` + whereSQL + ` order by code asc, id asc limit $` + fmt.Sprint(len(args)+1) + ` offset $` + fmt.Sprint(len(args)+2)
	if err := l.svcCtx.DB.QueryRowsCtx(l.ctx, &rows, query, queryArgs...); err != nil {
		return nil, err
	}
	return &pb.ConfigPageResp{Records: toConfigList(rows), Page: toPageInfo(total, pageNum, pageSize)}, nil
}
