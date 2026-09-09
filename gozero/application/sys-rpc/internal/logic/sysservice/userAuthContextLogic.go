package sysservicelogic

import (
	"context"
	"database/sql"
	"fmt"
	"sort"
	"strings"

	"github.com/gcc798/lightning/application/sys-rpc/internal/svc"
	"github.com/gcc798/lightning/application/sys-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserAuthContextLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserAuthContextLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserAuthContextLogic {
	return &UserAuthContextLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UserAuthContextLogic) UserAuthContext(in *pb.UserAuthContextReq) (*pb.UserAuthContextResp, error) {
	if in.UserId == 0 {
		return nil, fmt.Errorf("用户ID不能为空")
	}
	orgID, err := getUserOrgID(l.ctx, l.svcCtx, in.UserId)
	if err != nil {
		return nil, err
	}
	roles, err := getUserRoleKeys(l.ctx, l.svcCtx, in.UserId)
	if err != nil {
		return nil, err
	}
	permissions, err := getUserPermissionCodes(l.ctx, l.svcCtx, in.UserId, roles)
	if err != nil {
		return nil, err
	}

	return &pb.UserAuthContextResp{
		OrgId:       orgID,
		Roles:       roles,
		Permissions: permissions,
	}, nil
}

func getUserOrgID(ctx context.Context, svcCtx *svc.ServiceContext, userID int64) (int64, error) {
	var row struct {
		OrgID sql.NullInt64 `db:"org_id"`
	}
	if err := svcCtx.DB.QueryRowCtx(ctx, &row, `
		select org_id
		from public.s_user
		where id = $1
		limit 1
	`, userID); err != nil {
		return 0, fmt.Errorf("查询用户组织失败: %w", err)
	}
	return nullInt64(row.OrgID), nil
}

func getUserRoleKeys(ctx context.Context, svcCtx *svc.ServiceContext, userID int64) ([]string, error) {
	var rows []string
	if err := svcCtx.DB.QueryRowsCtx(ctx, &rows, `
		select distinct r.role_key
		from public.s_role r
		inner join public.m_user_role ur on ur.role_id = r.id
		where ur.user_id = $1
		  and r.status = 0
		  and r.role_key <> ''
		order by r.role_key asc
	`, userID); err != nil {
		return nil, fmt.Errorf("查询用户角色失败: %w", err)
	}
	return uniqueStringValues(rows), nil
}

func getUserPermissionCodes(ctx context.Context, svcCtx *svc.ServiceContext, userID int64, roles []string) ([]string, error) {
	permissions := make(map[string]struct{})
	var rows []string
	if err := svcCtx.DB.QueryRowsCtx(ctx, &rows, `
		select distinct p.code
		from public.s_api_permission p
		inner join public.m_role_api_permission rap on rap.permission_id = p.id
		inner join public.m_user_role ur on ur.role_id = rap.role_id
		inner join public.s_role r on r.id = ur.role_id
		where ur.user_id = $1
		  and r.status = 0
		  and p.status = 0
		  and p.code <> ''
		union
		select distinct p.code
		from public.s_api_permission p
		inner join public.m_user_api_permission uap on uap.permission_id = p.id
		where uap.user_id = $1
		  and p.status = 0
		  and p.code <> ''
		union
		select distinct m.perms
		from public.s_menu m
		inner join public.m_role_menu rm on rm.menu_id = m.id
		inner join public.m_user_role ur on ur.role_id = rm.role_id
		inner join public.s_role r on r.id = ur.role_id
		where ur.user_id = $1
		  and r.status = 0
		  and m.status = 0
		  and m.perms <> ''
		order by 1 asc
	`, userID); err != nil {
		return nil, fmt.Errorf("查询用户权限失败: %w", err)
	}
	for _, row := range rows {
		addPermission(permissions, row)
	}
	for _, role := range roles {
		if err := mergeRedisPermissions(ctx, svcCtx, "casbin:role:"+role+":permissions", permissions); err != nil {
			return nil, err
		}
	}
	if err := mergeRedisPermissions(ctx, svcCtx, fmt.Sprintf("casbin:user:%d:permissions", userID), permissions); err != nil {
		return nil, err
	}
	return sortedStringKeys(permissions), nil
}

func mergeRedisPermissions(ctx context.Context, svcCtx *svc.ServiceContext, key string, permissions map[string]struct{}) error {
	rows, err := svcCtx.Redis.SMembers(ctx, key).Result()
	if err != nil {
		return fmt.Errorf("查询缓存权限失败: %w", err)
	}
	for _, row := range rows {
		resource, _, _ := strings.Cut(row, "::")
		addPermission(permissions, resource)
	}
	return nil
}

func addPermission(permissions map[string]struct{}, value string) {
	value = strings.TrimSpace(value)
	if value == "" {
		return
	}
	permissions[value] = struct{}{}
}

func uniqueStringValues(values []string) []string {
	seen := make(map[string]struct{}, len(values))
	result := make([]string, 0, len(values))
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value == "" {
			continue
		}
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		result = append(result, value)
	}
	sort.Strings(result)
	return result
}

func sortedStringKeys(values map[string]struct{}) []string {
	result := make([]string, 0, len(values))
	for value := range values {
		result = append(result, value)
	}
	sort.Strings(result)
	return result
}
