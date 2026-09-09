package sysservicelogic

import (
	"context"
	"database/sql"

	"github.com/gcc798/lightning/application/sys-rpc/internal/svc"
	"github.com/gcc798/lightning/application/sys-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type RoleUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRoleUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RoleUpdateLogic {
	return &RoleUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RoleUpdateLogic) RoleUpdate(in *pb.RoleUpdateReq) (*pb.Ack, error) {
	row, err := getRoleByID(l.ctx, l.svcCtx, in.RoleId)
	if err != nil {
		return nil, err
	}
	roleName := row.RoleName
	if in.RoleName != "" {
		roleName = in.RoleName
	}
	if _, err := l.svcCtx.DB.ExecCtx(l.ctx, `update public.s_role set role_name = $2, sort = $3, status = $4, data_scope = $5, remark = $6, updated_time = now() where id = $1`,
		in.RoleId, roleName, in.Sort, in.Status, in.DataScope, sql.NullString{String: in.Remark, Valid: in.Remark != ""}); err != nil {
		return nil, err
	}
	return &pb.Ack{Msg: "ok"}, nil
}
