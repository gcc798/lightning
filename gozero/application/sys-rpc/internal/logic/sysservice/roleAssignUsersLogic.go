package sysservicelogic

import (
	"context"
	"fmt"

	"github.com/gcc798/microservice-kit/application/sys-rpc/internal/svc"
	"github.com/gcc798/microservice-kit/application/sys-rpc/pb"
	"github.com/lib/pq"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type RoleAssignUsersLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRoleAssignUsersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RoleAssignUsersLogic {
	return &RoleAssignUsersLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RoleAssignUsersLogic) RoleAssignUsers(in *pb.RoleUsersReq) (*pb.Ack, error) {
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
	var userCount int64
	if err := l.svcCtx.DB.QueryRowCtx(l.ctx, &userCount, `select count(1) from public.s_user where id = any($1)`, pq.Array(userIDs)); err != nil {
		return nil, fmt.Errorf("查询用户失败: %w", err)
	}
	if userCount != int64(len(userIDs)) {
		return nil, fmt.Errorf("部分用户不存在")
	}
	if err := l.svcCtx.DB.TransactCtx(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		for _, userID := range userIDs {
			var totalCount int64
			if err := session.QueryRowCtx(ctx, &totalCount, `select count(1) from public.m_user_role where user_id = $1 and role_id = $2`, userID, in.RoleId); err != nil {
				return err
			}
			if totalCount > 0 {
				continue
			}
			if _, err := session.ExecCtx(ctx, `insert into public.m_user_role (user_id, role_id, create_by, update_by, created_time, updated_time) values ($1, $2, nullif($3, 0), nullif($3, 0), now(), now())`, userID, in.RoleId, in.OperatorId); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return &pb.Ack{Msg: "ok"}, nil
}
