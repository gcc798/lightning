package sysservicelogic

import (
	"context"
	"fmt"

	"github.com/gcc798/microservice-kit/application/sys-rpc/internal/svc"
	"github.com/gcc798/microservice-kit/application/sys-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ApiPermissionDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewApiPermissionDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApiPermissionDeleteLogic {
	return &ApiPermissionDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ApiPermissionDeleteLogic) ApiPermissionDelete(in *pb.IdReq) (*pb.Ack, error) {
	var childCount int64
	if err := l.svcCtx.DB.QueryRowCtx(l.ctx, &childCount, `select count(1) from public.s_api_permission where parent_id = $1`, in.Id); err != nil {
		return nil, fmt.Errorf("检查子权限失败: %w", err)
	}
	if childCount > 0 {
		return nil, fmt.Errorf("存在子权限，无法删除")
	}
	roleIDs, userIDs, err := findAffectedPermissionSubjects(l.ctx, l.svcCtx, in.Id)
	if err != nil {
		l.Errorf("查询受影响主体失败: %v", err)
		roleIDs, userIDs = nil, nil
	}
	if err := l.svcCtx.DB.TransactCtx(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		if _, err := session.ExecCtx(ctx, `delete from public.m_role_api_permission where permission_id = $1`, in.Id); err != nil {
			return fmt.Errorf("删除角色 API 权限关联失败: %w", err)
		}
		if _, err := session.ExecCtx(ctx, `delete from public.m_user_api_permission where permission_id = $1`, in.Id); err != nil {
			return fmt.Errorf("删除用户 API 权限关联失败: %w", err)
		}
		result, err := session.ExecCtx(ctx, `delete from public.s_api_permission where id = $1`, in.Id)
		if err != nil {
			return fmt.Errorf("删除 API 权限失败: %w", err)
		}
		if rows, err := result.RowsAffected(); err == nil && rows == 0 {
			return fmt.Errorf("API 权限不存在")
		}
		return nil
	}); err != nil {
		return nil, err
	}
	for _, roleID := range roleIDs {
		if err := syncRolePermissionRedis(l.ctx, l.svcCtx, roleID); err != nil {
			l.Errorf("同步角色权限缓存失败 roleID=%d: %v", roleID, err)
		}
	}
	for _, userID := range userIDs {
		if err := syncUserPermissionRedis(l.ctx, l.svcCtx, userID); err != nil {
			l.Errorf("同步用户权限缓存失败 userID=%d: %v", userID, err)
		}
	}
	return &pb.Ack{Msg: "ok"}, nil
}
