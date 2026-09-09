package sysservicelogic

import (
	"context"
	"fmt"
	"strings"

	"github.com/gcc798/lightning/application/sys-rpc/internal/svc"
	"github.com/gcc798/lightning/application/sys-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type OrgPageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewOrgPageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrgPageLogic {
	return &OrgPageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *OrgPageLogic) OrgPage(in *pb.OrgPageReq) (*pb.OrgPageResp, error) {
	pageNum, pageSize := normalizePage(in.PageNum, in.PageSize)
	where := []string{"1=1"}
	args := make([]interface{}, 0)
	if in.OrgName != "" {
		args = append(args, "%"+in.OrgName+"%")
		where = append(where, fmt.Sprintf("org_name like $%d", len(args)))
	}
	if in.OrgCode != "" {
		args = append(args, in.OrgCode)
		where = append(where, fmt.Sprintf("org_code = $%d", len(args)))
	}
	if in.Status == 0 || in.Status == 1 {
		args = append(args, in.Status)
		where = append(where, fmt.Sprintf("status = $%d", len(args)))
	}
	if in.ParentId > 0 {
		args = append(args, in.ParentId)
		where = append(where, fmt.Sprintf("parent_id = $%d", len(args)))
	}
	whereSQL := strings.Join(where, " and ")
	var total int64
	if err := l.svcCtx.DB.QueryRowCtx(l.ctx, &total, "select count(1) from public.s_org where "+whereSQL, args...); err != nil {
		return nil, err
	}
	queryArgs := append(append([]interface{}{}, args...), pageSize, (pageNum-1)*pageSize)
	var rows []orgRow
	query := `select id, parent_id, ancestors, org_name, org_code, org_type, leader, phone, email, status, sort, remark, created_time, updated_time
		from public.s_org where ` + whereSQL + ` order by sort asc, id asc limit $` + fmt.Sprint(len(args)+1) + ` offset $` + fmt.Sprint(len(args)+2)
	if err := l.svcCtx.DB.QueryRowsCtx(l.ctx, &rows, query, queryArgs...); err != nil {
		return nil, err
	}
	return &pb.OrgPageResp{Records: toOrgList(rows), Page: toPageInfo(total, pageNum, pageSize)}, nil
}
