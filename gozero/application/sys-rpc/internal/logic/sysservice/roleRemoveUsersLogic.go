package sysservicelogic

import (
	"context"
	"fmt"

	"github.com/gcc798/microservice-kit/application/sys-rpc/internal/svc"
	"github.com/gcc798/microservice-kit/application/sys-rpc/pb"
	"github.com/lib/pq"

	"github.com/zeromicro/go-zero/core/logx"
)

type RoleRemoveUsersLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRoleRemoveUsersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RoleRemoveUsersLogic {
	return &RoleRemoveUsersLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RoleRemoveUsersLogic) RoleRemoveUsers(in *pb.RoleUsersReq) (*pb.Ack, error) {
	if in.RoleId <= 0 {
		return nil, fmt.Errorf("角色ID不能为空")
	}
	if _, err := getRoleByID(l.ctx, l.svcCtx, in.RoleId); err != nil {
		return nil, err
	}
	userIDs := uniqueInt64Values(in.UserIds)
	if len(userIDs) == 0 {
		return &pb.Ack{Msg: "ok"}, nil
	}
	if _, err := l.svcCtx.DB.ExecCtx(l.ctx, `
		delete from public.m_user_role
		where role_id = $1 and user_id = any($2)
	`, in.RoleId, pq.Array(userIDs)); err != nil {
		return nil, fmt.Errorf("批量移除角色用户失败: %w", err)
	}
	return &pb.Ack{Msg: "ok"}, nil
}
