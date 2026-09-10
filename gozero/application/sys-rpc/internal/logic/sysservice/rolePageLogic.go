package sysservicelogic

import (
	"context"
	"fmt"
	"strings"

	"github.com/gcc798/microservice-kit/application/sys-rpc/internal/svc"
	"github.com/gcc798/microservice-kit/application/sys-rpc/pb"
	"github.com/zeromicro/go-zero/core/logx"
)

type RolePageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRolePageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RolePageLogic {
	return &RolePageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RolePageLogic) RolePage(in *pb.RolePageReq) (*pb.RolePageResp, error) {
	pageNum, pageSize := normalizePage(in.PageNum, in.PageSize)
	where := []string{"1=1"}
	args := make([]interface{}, 0)
	if in.RoleName != "" {
		args = append(args, "%"+in.RoleName+"%")
		where = append(where, fmt.Sprintf("role_name like $%d", len(args)))
	}
	if in.Status == 0 || in.Status == 1 {
		args = append(args, in.Status)
		where = append(where, fmt.Sprintf("status = $%d", len(args)))
	}
	whereSQL := strings.Join(where, " and ")
	var total int64
	if err := l.svcCtx.DB.QueryRowCtx(l.ctx, &total, "select count(1) from public.s_role where "+whereSQL, args...); err != nil {
		return nil, err
	}
	queryArgs := append(append([]interface{}{}, args...), pageSize, (pageNum-1)*pageSize)
	var rows []roleRow
	query := `select id, role_key, role_name, sort, status, data_scope, is_system, remark, create_by, created_time
		from public.s_role where ` + whereSQL + ` order by id desc limit $` + fmt.Sprint(len(args)+1) + ` offset $` + fmt.Sprint(len(args)+2)
	if err := l.svcCtx.DB.QueryRowsCtx(l.ctx, &rows, query, queryArgs...); err != nil {
		return nil, err
	}
	return &pb.RolePageResp{Records: toRoleList(rows), Page: toPageInfo(total, pageNum, pageSize)}, nil
}
