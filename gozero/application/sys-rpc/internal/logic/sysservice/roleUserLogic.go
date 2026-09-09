package sysservicelogic

import (
	"context"

	"github.com/gcc798/lightning/application/sys-rpc/internal/svc"
	"github.com/gcc798/lightning/application/sys-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type RoleUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRoleUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RoleUserLogic {
	return &RoleUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RoleUserLogic) RoleUser(in *pb.UserRoleQueryReq) (*pb.RoleListResp, error) {
	var rows []roleRow
	if err := l.svcCtx.DB.QueryRowsCtx(l.ctx, &rows, `
		select r.id, r.role_key, r.role_name, r.sort, r.status, r.data_scope, r.is_system, r.remark, r.create_by, r.created_time
		from public.s_role r
		inner join public.m_user_role ur on ur.role_id = r.id
		where ur.user_id = $1
		order by r.sort asc, r.id asc
	`, in.UserId); err != nil {
		return nil, err
	}
	return &pb.RoleListResp{Records: toRoleList(rows)}, nil
}
