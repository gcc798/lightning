package sysservicelogic

import (
	"context"
	"errors"

	"github.com/gcc798/microservice-kit/application/sys-rpc/internal/svc"
	"github.com/gcc798/microservice-kit/application/sys-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type RoleAssignLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRoleAssignLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RoleAssignLogic {
	return &RoleAssignLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RoleAssignLogic) RoleAssign(in *pb.AssignRoleReq) (*pb.Ack, error) {
	if _, err := getUserByID(l.ctx, l.svcCtx, in.UserId); err != nil {
		return nil, err
	}
	if _, err := getRoleByID(l.ctx, l.svcCtx, in.RoleId); err != nil {
		return nil, err
	}
	var totalCount int64
	if err := l.svcCtx.DB.QueryRowCtx(l.ctx, &totalCount, `select count(1) from public.m_user_role where user_id = $1 and role_id = $2`, in.UserId, in.RoleId); err != nil {
		return nil, err
	}
	if totalCount > 0 {
		return nil, errors.New("用户已拥有该角色")
	}
	if _, err := l.svcCtx.DB.ExecCtx(l.ctx, `insert into public.m_user_role (user_id, role_id, create_by, update_by, created_time, updated_time) values ($1, $2, null, null, now(), now())`, in.UserId, in.RoleId); err != nil {
		return nil, err
	}
	return &pb.Ack{Msg: "ok"}, nil
}
