package iam

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/gcc798/microservice-kit/application/iam/internal/domain/model"
	"gorm.io/gorm"
)

// ApiPermissionTree API 权限树节点。
type ApiPermissionTree struct {
	model.ApiPermission
	Children []*ApiPermissionTree `json:"children,omitempty"`
}

// RolePermissionGrants separates editable grants from menu-derived grants.
type RolePermissionGrants struct {
	ManualIds  []int64 `json:"manualIds"`
	DerivedIds []int64 `json:"derivedIds"`
}

// ApiPermissionService 管理 API 权限资源和授权映射。
type ApiPermissionService interface {
	Tree(ctx context.Context) ([]*ApiPermissionTree, error)
	List(ctx context.Context) ([]model.ApiPermission, error)
	Create(ctx context.Context, req *ApiPermissionSaveRequest, userId int64) (*model.ApiPermission, error)
	Update(ctx context.Context, id int64, req *ApiPermissionSaveRequest, userId int64) error
	Delete(ctx context.Context, id int64) error
	GetRolePermissions(ctx context.Context, roleId int64) (*RolePermissionGrants, error)
	AssignRolePermissions(ctx context.Context, roleId int64, permissionIds []int64, userId int64) error
	GetUserPermissionIds(ctx context.Context, userId int64) ([]int64, error)
	AssignUserPermissions(ctx context.Context, targetUserId int64, permissionIds []int64, operatorId int64) error
}

type apiPermissionService struct {
	db *gorm.DB
}

func NewApiPermissionService(db *gorm.DB) ApiPermissionService {
	return &apiPermissionService{db: db}
}

func (s *apiPermissionService) Tree(ctx context.Context) ([]*ApiPermissionTree, error) {
	permissions, err := s.List(ctx)
	if err != nil {
		return nil, err
	}
	return buildApiPermissionTree(permissions, 0), nil
}

func (s *apiPermissionService) List(ctx context.Context) ([]model.ApiPermission, error) {
	var permissions []model.ApiPermission
	err := s.db.WithContext(ctx).
		Order("sort ASC, created_time ASC").
		Find(&permissions).Error
	if err != nil {
		return nil, fmt.Errorf("查询 API 权限失败: %w", err)
	}
	return permissions, nil
}

func (s *apiPermissionService) Create(ctx context.Context, req *ApiPermissionSaveRequest, userId int64) (*model.ApiPermission, error) {
	permission := &model.ApiPermission{
		ParentId: req.ParentId,
		Module:   req.Module,
		Code:     req.Code,
		Name:     req.Name,
		NodeType: req.NodeType,
		Action:   normalizeAction(req.Code, req.Action),
		Method:   strings.ToUpper(req.Method),
		Path:     req.Path,
		Sort:     req.Sort,
		Status:   req.Status,
		Remark:   req.Remark,
		CreateBy: userId,
		UpdateBy: userId,
	}
	if err := s.validatePermission(ctx, permission, 0); err != nil {
		return nil, err
	}
	if err := s.db.WithContext(ctx).Create(permission).Error; err != nil {
		return nil, fmt.Errorf("创建 API 权限失败: %w", err)
	}
	return permission, nil
}

func (s *apiPermissionService) Update(ctx context.Context, id int64, req *ApiPermissionSaveRequest, userId int64) error {
	var existing model.ApiPermission
	if err := s.db.WithContext(ctx).Where("id = ?", id).First(&existing).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("API 权限不存在")
		}
		return fmt.Errorf("查询 API 权限失败: %w", err)
	}
	existing.ParentId = req.ParentId
	existing.Module = req.Module
	existing.Code = req.Code
	existing.Name = req.Name
	existing.NodeType = req.NodeType
	existing.Action = normalizeAction(req.Code, req.Action)
	existing.Method = strings.ToUpper(req.Method)
	existing.Path = req.Path
	existing.Sort = req.Sort
	existing.Status = req.Status
	existing.Remark = req.Remark
	existing.UpdateBy = userId
	if err := s.validatePermission(ctx, &existing, id); err != nil {
		return err
	}
	if err := s.db.WithContext(ctx).Save(&existing).Error; err != nil {
		return fmt.Errorf("更新 API 权限失败: %w", err)
	}
	return nil
}

func (s *apiPermissionService) Delete(ctx context.Context, id int64) error {
	var childCount int64
	if err := s.db.WithContext(ctx).Model(&model.ApiPermission{}).Where("parent_id = ?", id).Count(&childCount).Error; err != nil {
		return fmt.Errorf("检查子权限失败: %w", err)
	}
	if childCount > 0 {
		return fmt.Errorf("存在子权限，无法删除")
	}
	if err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if _, err := gorm.G[model.MMenuApiPermission](tx).Where("permission_id = ?", id).Delete(ctx); err != nil {
			return fmt.Errorf("删除菜单 API 权限关联失败: %w", err)
		}
		if _, err := gorm.G[model.MRoleApiPermission](tx).Where("permission_id = ?", id).Delete(ctx); err != nil {
			return fmt.Errorf("删除角色 API 权限关联失败: %w", err)
		}
		if _, err := gorm.G[model.MUserApiPermission](tx).Where("permission_id = ?", id).Delete(ctx); err != nil {
			return fmt.Errorf("删除用户 API 权限关联失败: %w", err)
		}
		if _, err := gorm.G[model.ApiPermission](tx).Where("id = ?", id).Delete(ctx); err != nil {
			return fmt.Errorf("删除 API 权限失败: %w", err)
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

func (s *apiPermissionService) GetRolePermissions(ctx context.Context, roleId int64) (*RolePermissionGrants, error) {
	var rows []model.MRoleApiPermission
	if err := s.db.WithContext(ctx).Where("role_id = ?", roleId).Find(&rows).Error; err != nil {
		return nil, fmt.Errorf("查询角色 API 权限失败: %w", err)
	}
	grants := &RolePermissionGrants{ManualIds: []int64{}, DerivedIds: []int64{}}
	for _, row := range rows {
		if row.Source == model.PermissionSourceDerived {
			grants.DerivedIds = append(grants.DerivedIds, row.PermissionId)
		} else {
			grants.ManualIds = append(grants.ManualIds, row.PermissionId)
		}
	}
	return grants, nil
}

func (s *apiPermissionService) AssignRolePermissions(ctx context.Context, roleId int64, permissionIds []int64, userId int64) error {
	var role model.Role
	if err := s.db.WithContext(ctx).Where("id = ?", roleId).First(&role).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("角色不存在")
		}
		return fmt.Errorf("查询角色失败: %w", err)
	}
	_, normalizedIds, err := s.resolveAssignablePermissions(ctx, permissionIds)
	if err != nil {
		return err
	}
	if err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if _, err := gorm.G[model.MRoleApiPermission](tx).
			Where("role_id = ? AND source = ?", roleId, model.PermissionSourceManual).Delete(ctx); err != nil {
			return fmt.Errorf("删除旧角色 API 权限失败: %w", err)
		}
		rows := make([]model.MRoleApiPermission, 0, len(normalizedIds))
		for _, permissionId := range normalizedIds {
			rows = append(rows, model.MRoleApiPermission{RoleId: roleId, PermissionId: permissionId, Source: model.PermissionSourceManual, CreateBy: userId, UpdateBy: userId})
		}
		if len(rows) > 0 {
			if err := tx.Create(&rows).Error; err != nil {
				return fmt.Errorf("保存角色 API 权限失败: %w", err)
			}
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

func (s *apiPermissionService) GetUserPermissionIds(ctx context.Context, userId int64) ([]int64, error) {
	var rows []model.MUserApiPermission
	if err := s.db.WithContext(ctx).Where("user_id = ?", userId).Find(&rows).Error; err != nil {
		return nil, fmt.Errorf("查询用户 API 权限失败: %w", err)
	}
	ids := make([]int64, 0, len(rows))
	for _, row := range rows {
		ids = append(ids, row.PermissionId)
	}
	return ids, nil
}

func (s *apiPermissionService) AssignUserPermissions(ctx context.Context, targetUserId int64, permissionIds []int64, operatorId int64) error {
	if err := s.db.WithContext(ctx).Where("id = ?", targetUserId).First(&model.User{}).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("用户不存在")
		}
		return fmt.Errorf("查询用户失败: %w", err)
	}
	_, normalizedIds, err := s.resolveAssignablePermissions(ctx, permissionIds)
	if err != nil {
		return err
	}
	if err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if _, err := gorm.G[model.MUserApiPermission](tx).
			Where("user_id = ? AND source = ?", targetUserId, model.PermissionSourceManual).Delete(ctx); err != nil {
			return fmt.Errorf("删除旧用户 API 权限失败: %w", err)
		}
		rows := make([]model.MUserApiPermission, 0, len(normalizedIds))
		for _, permissionId := range normalizedIds {
			rows = append(rows, model.MUserApiPermission{UserId: targetUserId, PermissionId: permissionId, Source: model.PermissionSourceManual, CreateBy: operatorId, UpdateBy: operatorId})
		}
		if len(rows) > 0 {
			if err := tx.Create(&rows).Error; err != nil {
				return fmt.Errorf("保存用户 API 权限失败: %w", err)
			}
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

func (s *apiPermissionService) validatePermission(ctx context.Context, permission *model.ApiPermission, selfId int64) error {
	if permission.ID != 0 && permission.ParentId == permission.ID {
		return fmt.Errorf("不能将自己设置为父级权限")
	}
	if err := s.validateParentChain(ctx, permission.ParentId, selfId); err != nil {
		return err
	}
	var count int64
	query := s.db.WithContext(ctx).Model(&model.ApiPermission{}).Where("code = ?", permission.Code)
	if selfId != 0 {
		query = query.Where("id != ?", selfId)
	}
	if err := query.Count(&count).Error; err != nil {
		return fmt.Errorf("检查权限标识失败: %w", err)
	}
	if count > 0 {
		return fmt.Errorf("权限标识已存在")
	}
	return nil
}

func (s *apiPermissionService) validateParentChain(ctx context.Context, parentId int64, selfId int64) error {
	visited := make(map[int64]struct{})
	for parentId != 0 {
		if _, ok := visited[parentId]; ok {
			return fmt.Errorf("父级权限存在循环引用")
		}
		visited[parentId] = struct{}{}
		if selfId != 0 && parentId == selfId {
			return fmt.Errorf("不能将自己的下级设置为父级权限")
		}
		var parent model.ApiPermission
		if err := s.db.WithContext(ctx).Where("id = ?", parentId).First(&parent).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("父级权限不存在")
			}
			return fmt.Errorf("查询父级权限失败: %w", err)
		}
		parentId = parent.ParentId
	}
	return nil
}

func (s *apiPermissionService) resolveAssignablePermissions(ctx context.Context, permissionIds []int64) ([]model.ApiPermission, []int64, error) {
	if len(permissionIds) == 0 {
		return nil, nil, nil
	}
	uniqueIds := make([]int64, 0, len(permissionIds))
	seen := make(map[int64]struct{}, len(permissionIds))
	for _, id := range permissionIds {
		if id == 0 {
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		uniqueIds = append(uniqueIds, id)
	}
	var permissions []model.ApiPermission
	if err := s.db.WithContext(ctx).Where("id IN ? AND status = 0", uniqueIds).Find(&permissions).Error; err != nil {
		return nil, nil, fmt.Errorf("查询 API 权限失败: %w", err)
	}
	if len(permissions) != len(uniqueIds) {
		return nil, nil, fmt.Errorf("存在无效或停用的 API 权限")
	}
	normalized := normalizeCoveredPermissions(permissions)
	normalizedIds := make([]int64, 0, len(normalized))
	for _, item := range normalized {
		normalizedIds = append(normalizedIds, item.ID)
	}
	sort.Slice(normalizedIds, func(i, j int) bool { return normalizedIds[i] < normalizedIds[j] })
	return normalized, normalizedIds, nil
}

func uniqueInt64(values []int64) []int64 {
	if len(values) == 0 {
		return nil
	}
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

func normalizeCoveredPermissions(permissions []model.ApiPermission) []model.ApiPermission {
	sort.Slice(permissions, func(i, j int) bool {
		return len(permissions[i].Code) < len(permissions[j].Code)
	})
	selected := make([]model.ApiPermission, 0, len(permissions))
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

func normalizeAction(code, action string) string {
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

func buildApiPermissionTree(permissions []model.ApiPermission, parentId int64) []*ApiPermissionTree {
	tree := make([]*ApiPermissionTree, 0)
	for _, permission := range permissions {
		if permission.ParentId != parentId {
			continue
		}
		node := &ApiPermissionTree{
			ApiPermission: permission,
			Children:      buildApiPermissionTree(permissions, permission.ID),
		}
		tree = append(tree, node)
	}
	return tree
}
