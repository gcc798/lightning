package sysservicelogic

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/gcc798/lightning/application/sys-rpc/internal/svc"
	"github.com/gcc798/lightning/application/sys-rpc/pb"
	"github.com/lib/pq"
	gzsqlx "github.com/zeromicro/go-zero/core/stores/sqlx"
)

type apiPermissionRow struct {
	Id          int64          `db:"id"`
	ParentId    int64          `db:"parent_id"`
	Module      string         `db:"module"`
	Code        string         `db:"code"`
	Name        string         `db:"name"`
	NodeType    int64          `db:"node_type"`
	Action      string         `db:"action"`
	Method      sql.NullString `db:"method"`
	Path        sql.NullString `db:"path"`
	Sort        int64          `db:"sort"`
	Status      int64          `db:"status"`
	Remark      sql.NullString `db:"remark"`
	CreateBy    sql.NullInt64  `db:"create_by"`
	UpdateBy    sql.NullInt64  `db:"update_by"`
	CreatedTime sql.NullTime   `db:"created_time"`
	UpdatedTime sql.NullTime   `db:"updated_time"`
}

func listApiPermissions(ctx context.Context, svcCtx *svc.ServiceContext) ([]apiPermissionRow, error) {
	var rows []apiPermissionRow
	err := svcCtx.DB.QueryRowsCtx(ctx, &rows, `
		select id, parent_id, module, code, name, node_type, action, method, path, sort, status, remark, create_by, update_by, created_time, updated_time
		from public.s_api_permission
		order by sort asc, created_time asc
	`)
	if err != nil {
		return nil, fmt.Errorf("查询 API 权限失败: %w", err)
	}
	return rows, nil
}

func getApiPermissionByID(ctx context.Context, svcCtx *svc.ServiceContext, id int64) (*apiPermissionRow, error) {
	var row apiPermissionRow
	err := svcCtx.DB.QueryRowCtx(ctx, &row, `
		select id, parent_id, module, code, name, node_type, action, method, path, sort, status, remark, create_by, update_by, created_time, updated_time
		from public.s_api_permission
		where id = $1
		limit 1
	`, id)
	if err != nil {
		if errors.Is(err, gzsqlx.ErrNotFound) {
			return nil, fmt.Errorf("API 权限不存在")
		}
		return nil, err
	}
	return &row, nil
}

func toApiPermissionPB(row apiPermissionRow) *pb.ApiPermission {
	return &pb.ApiPermission{
		Id:          row.Id,
		ParentId:    row.ParentId,
		Module:      row.Module,
		Code:        row.Code,
		Name:        row.Name,
		NodeType:    int32(row.NodeType),
		Action:      row.Action,
		Method:      nullString(row.Method),
		Path:        nullString(row.Path),
		Sort:        row.Sort,
		Status:      int32(row.Status),
		Remark:      nullString(row.Remark),
		CreateBy:    nullInt64(row.CreateBy),
		UpdateBy:    nullInt64(row.UpdateBy),
		CreatedTime: nullTime(row.CreatedTime),
		UpdatedTime: nullTime(row.UpdatedTime),
	}
}

func toApiPermissionList(rows []apiPermissionRow) []*pb.ApiPermission {
	list := make([]*pb.ApiPermission, 0, len(rows))
	for _, row := range rows {
		list = append(list, toApiPermissionPB(row))
	}
	return list
}

func buildApiPermissionTree(rows []apiPermissionRow, parentID int64) []*pb.ApiPermission {
	tree := make([]*pb.ApiPermission, 0)
	for _, row := range rows {
		if row.ParentId != parentID {
			continue
		}
		node := toApiPermissionPB(row)
		node.Children = buildApiPermissionTree(rows, row.Id)
		tree = append(tree, node)
	}
	return tree
}

func validateApiPermission(ctx context.Context, svcCtx *svc.ServiceContext, req *pb.ApiPermissionSaveReq) error {
	if strings.TrimSpace(req.Module) == "" || strings.TrimSpace(req.Code) == "" || strings.TrimSpace(req.Name) == "" {
		return fmt.Errorf("模块、权限标识和权限名称不能为空")
	}
	if strings.TrimSpace(req.Action) == "" {
		return fmt.Errorf("操作类型不能为空")
	}
	if req.Id != 0 && req.ParentId == req.Id {
		return fmt.Errorf("不能将自己设置为父级权限")
	}
	if err := validateApiPermissionParent(ctx, svcCtx, req.ParentId, req.Id); err != nil {
		return err
	}
	var count int64
	if err := svcCtx.DB.QueryRowCtx(ctx, &count, `
		select count(1) from public.s_api_permission
		where code = $1 and id <> $2
	`, req.Code, req.Id); err != nil {
		return fmt.Errorf("检查权限标识失败: %w", err)
	}
	if count > 0 {
		return fmt.Errorf("权限标识已存在")
	}
	return nil
}

func validateApiPermissionParent(ctx context.Context, svcCtx *svc.ServiceContext, parentID, selfID int64) error {
	visited := make(map[int64]struct{})
	for parentID != 0 {
		if _, ok := visited[parentID]; ok {
			return fmt.Errorf("父级权限存在循环引用")
		}
		visited[parentID] = struct{}{}
		if selfID != 0 && parentID == selfID {
			return fmt.Errorf("不能将自己的下级设置为父级权限")
		}
		parent, err := getApiPermissionByID(ctx, svcCtx, parentID)
		if err != nil {
			return fmt.Errorf("父级权限不存在")
		}
		parentID = parent.ParentId
	}
	return nil
}

func apiPermissionIDsByOwner(ctx context.Context, svcCtx *svc.ServiceContext, table, ownerColumn string, ownerID int64) ([]int64, error) {
	query := fmt.Sprintf(`select permission_id from public.%s where %s = $1 order by permission_id asc`, table, ownerColumn)
	var ids []int64
	if err := svcCtx.DB.QueryRowsCtx(ctx, &ids, query, ownerID); err != nil {
		return nil, err
	}
	return ids, nil
}

func resolveAssignableApiPermissions(ctx context.Context, svcCtx *svc.ServiceContext, ids []int64) ([]apiPermissionRow, []int64, error) {
	uniqueIDs := uniqueInt64Values(ids)
	if len(uniqueIDs) == 0 {
		return nil, nil, nil
	}
	var rows []apiPermissionRow
	if err := svcCtx.DB.QueryRowsCtx(ctx, &rows, `
		select id, parent_id, module, code, name, node_type, action, method, path, sort, status, remark, create_by, update_by, created_time, updated_time
		from public.s_api_permission
		where id = any($1) and status = 0
	`, pq.Array(uniqueIDs)); err != nil {
		return nil, nil, fmt.Errorf("查询 API 权限失败: %w", err)
	}
	if len(rows) != len(uniqueIDs) {
		return nil, nil, fmt.Errorf("存在无效或停用的 API 权限")
	}
	normalized := normalizeCoveredApiPermissions(rows)
	normalizedIDs := make([]int64, 0, len(normalized))
	for _, row := range normalized {
		normalizedIDs = append(normalizedIDs, row.Id)
	}
	sort.Slice(normalizedIDs, func(i, j int) bool { return normalizedIDs[i] < normalizedIDs[j] })
	return normalized, normalizedIDs, nil
}

func replaceRoleApiPermissions(ctx context.Context, svcCtx *svc.ServiceContext, roleID int64, roleKey string, permissionIDs []int64, operatorID int64) error {
	permissions, normalizedIDs, err := resolveAssignableApiPermissions(ctx, svcCtx, permissionIDs)
	if err != nil {
		return err
	}
	if err := svcCtx.DB.TransactCtx(ctx, func(ctx context.Context, session gzsqlx.Session) error {
		if _, err := session.ExecCtx(ctx, `delete from public.m_role_api_permission where role_id = $1`, roleID); err != nil {
			return fmt.Errorf("删除旧角色 API 权限失败: %w", err)
		}
		for _, permissionID := range normalizedIDs {
			if _, err := session.ExecCtx(ctx, `
				insert into public.m_role_api_permission (role_id, permission_id, create_by, update_by, created_time, updated_time)
				values ($1, $2, nullif($3, 0), nullif($3, 0), now(), now())
			`, roleID, permissionID, operatorID); err != nil {
				return fmt.Errorf("保存角色 API 权限失败: %w", err)
			}
		}
		return nil
	}); err != nil {
		return err
	}
	return replacePermissionRedis(ctx, svcCtx, "casbin:role:"+roleKey+":permissions", permissions)
}

func replaceUserApiPermissions(ctx context.Context, svcCtx *svc.ServiceContext, userID int64, permissionIDs []int64, operatorID int64) error {
	permissions, normalizedIDs, err := resolveAssignableApiPermissions(ctx, svcCtx, permissionIDs)
	if err != nil {
		return err
	}
	if err := svcCtx.DB.TransactCtx(ctx, func(ctx context.Context, session gzsqlx.Session) error {
		if _, err := session.ExecCtx(ctx, `delete from public.m_user_api_permission where user_id = $1`, userID); err != nil {
			return fmt.Errorf("删除旧用户 API 权限失败: %w", err)
		}
		for _, permissionID := range normalizedIDs {
			if _, err := session.ExecCtx(ctx, `
				insert into public.m_user_api_permission (user_id, permission_id, create_by, update_by, created_time, updated_time)
				values ($1, $2, nullif($3, 0), nullif($3, 0), now(), now())
			`, userID, permissionID, operatorID); err != nil {
				return fmt.Errorf("保存用户 API 权限失败: %w", err)
			}
		}
		return nil
	}); err != nil {
		return err
	}
	return replacePermissionRedis(ctx, svcCtx, fmt.Sprintf("casbin:user:%d:permissions", userID), permissions)
}

func replacePermissionRedis(ctx context.Context, svcCtx *svc.ServiceContext, key string, permissions []apiPermissionRow) error {
	if err := svcCtx.Redis.Del(ctx, key).Err(); err != nil {
		return err
	}
	if len(permissions) == 0 {
		return nil
	}
	values := make([]interface{}, 0, len(permissions))
	for _, permission := range permissions {
		values = append(values, permission.Code+"::"+normalizeApiPermissionAction(permission.Code, permission.Action))
	}
	return svcCtx.Redis.SAdd(ctx, key, values...).Err()
}

func normalizeCoveredApiPermissions(permissions []apiPermissionRow) []apiPermissionRow {
	sort.Slice(permissions, func(i, j int) bool {
		return len(permissions[i].Code) < len(permissions[j].Code)
	})
	selected := make([]apiPermissionRow, 0, len(permissions))
	wildcards := make([]string, 0)
	for _, permission := range permissions {
		covered := false
		for _, wildcard := range wildcards {
			prefix := strings.TrimSuffix(wildcard, "*")
			if permission.Code != wildcard && strings.HasPrefix(permission.Code, prefix) {
				covered = true
				break
			}
		}
		if covered {
			continue
		}
		selected = append(selected, permission)
		if strings.HasSuffix(permission.Code, ".*") || permission.Code == "*" {
			wildcards = append(wildcards, permission.Code)
		}
	}
	return selected
}

func normalizeApiPermissionAction(code, action string) string {
	if code == "*" || strings.HasSuffix(code, ".*") {
		return "*"
	}
	if action == "" {
		if strings.HasSuffix(code, ".read") {
			return "read"
		}
		return "write"
	}
	return action
}

func uniqueInt64Values(values []int64) []int64 {
	seen := make(map[int64]struct{}, len(values))
	result := make([]int64, 0, len(values))
	for _, value := range values {
		if value == 0 {
			continue
		}
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		result = append(result, value)
	}
	sort.Slice(result, func(i, j int) bool { return result[i] < result[j] })
	return result
}

func findAffectedPermissionSubjects(ctx context.Context, svcCtx *svc.ServiceContext, permissionID int64) ([]int64, []int64, error) {
	var roleIDs []int64
	if err := svcCtx.DB.QueryRowsCtx(ctx, &roleIDs, `select distinct role_id from public.m_role_api_permission where permission_id = $1`, permissionID); err != nil {
		return nil, nil, fmt.Errorf("查询受影响角色失败: %w", err)
	}
	var userIDs []int64
	if err := svcCtx.DB.QueryRowsCtx(ctx, &userIDs, `select distinct user_id from public.m_user_api_permission where permission_id = $1`, permissionID); err != nil {
		return nil, nil, fmt.Errorf("查询受影响用户失败: %w", err)
	}
	return roleIDs, userIDs, nil
}

func syncRolePermissionRedis(ctx context.Context, svcCtx *svc.ServiceContext, roleID int64) error {
	var roleKey string
	if err := svcCtx.DB.QueryRowCtx(ctx, &roleKey, `select role_key from public.s_role where id = $1`, roleID); err != nil {
		return nil
	}
	var rows []apiPermissionRow
	if err := svcCtx.DB.QueryRowsCtx(ctx, &rows, `
		select p.id, p.parent_id, p.module, p.code, p.name, p.node_type, p.action, p.method, p.path, p.sort, p.status, p.remark, p.create_by, p.update_by, p.created_time, p.updated_time
		from public.s_api_permission p
		join public.m_role_api_permission rp on rp.permission_id = p.id
		where rp.role_id = $1 and p.status = 0
	`, roleID); err != nil {
		return fmt.Errorf("查询角色API权限失败: %w", err)
	}
	normalized := normalizeCoveredApiPermissions(rows)
	return replacePermissionRedis(ctx, svcCtx, "casbin:role:"+roleKey+":permissions", normalized)
}

func syncUserPermissionRedis(ctx context.Context, svcCtx *svc.ServiceContext, userID int64) error {
	var rows []apiPermissionRow
	if err := svcCtx.DB.QueryRowsCtx(ctx, &rows, `
		select p.id, p.parent_id, p.module, p.code, p.name, p.node_type, p.action, p.method, p.path, p.sort, p.status, p.remark, p.create_by, p.update_by, p.created_time, p.updated_time
		from public.s_api_permission p
		join public.m_user_api_permission up on up.permission_id = p.id
		where up.user_id = $1 and p.status = 0
	`, userID); err != nil {
		return fmt.Errorf("查询用户API权限失败: %w", err)
	}
	normalized := normalizeCoveredApiPermissions(rows)
	return replacePermissionRedis(ctx, svcCtx, fmt.Sprintf("casbin:user:%d:permissions", userID), normalized)
}
