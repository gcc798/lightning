package iam

import (
	"context"
	"errors"
	"fmt"

	"github.com/gcc798/lightning/application/iam/internal/domain/model"
	"github.com/gcc798/lightning/internal/logger"
	"github.com/gcc798/lightning/internal/utils/pagination"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// RoleService 角色管理服务接口
type RoleService interface {
	// Create 创建角色
	Create(ctx context.Context, role *model.Role) error

	// Update 更新角色
	Update(ctx context.Context, role *model.Role) error

	// Delete 删除角色
	Delete(ctx context.Context, roleId int64) error

	// GetById 根据ID查询角色
	GetById(ctx context.Context, roleId int64) (*model.Role, error)

	// GetByRoleKey 根据角色标识查询角色
	GetByRoleKey(ctx context.Context, roleKey string) (*model.Role, error)

	// Page 分页查询角色列表
	Page(ctx context.Context, pageNum, pageSize int, roleName string, status int32) (*pagination.Page[model.Role], error)

	// AssignRoleToUser 为用户分配角色
	AssignRoleToUser(ctx context.Context, userId, roleId int64) error

	// RemoveRoleFromUser 移除用户的角色
	RemoveRoleFromUser(ctx context.Context, userId, roleId int64) error

	// GetUserRoles 获取用户的所有角色
	GetUserRoles(ctx context.Context, userId int64) ([]model.Role, error)

	// GetRoleUsers 获取拥有该角色的用户
	GetRoleUsers(ctx context.Context, roleId int64) ([]model.User, error)

	// AssignUsersToRole 批量为角色添加用户
	AssignUsersToRole(ctx context.Context, roleId int64, userIds []int64, operatorId int64) error

	// RemoveUsersFromRole 批量移除角色下的用户
	RemoveUsersFromRole(ctx context.Context, roleId int64, userIds []int64) error

	// AssignMenusToRole 为角色分配菜单权限
	AssignMenusToRole(ctx context.Context, roleId int64, menuIds []int64) error

	// GetRoleMenus 获取角色的所有菜单
	GetRoleMenus(ctx context.Context, roleId int64) ([]model.Menu, error)
}

type roleService struct {
	db     *gorm.DB
	logger logger.Logger
}

// NewRoleService 创建角色服务实例
func NewRoleService(db *gorm.DB, logger logger.Logger) RoleService {
	return &roleService{
		db:     db,
		logger: logger,
	}
}

// Create 创建角色
func (s *roleService) Create(ctx context.Context, role *model.Role) error {
	// 检查角色标识是否已存在
	count, err := gorm.G[model.Role](s.db).Where("role_key = ?", role.RoleKey).Count(ctx, "id")
	if err != nil {
		s.logger.Error("检查角色标识失败", zap.Error(err))
		return fmt.Errorf("检查角色标识失败: %w", err)
	}
	if count > 0 {
		return fmt.Errorf("角色标识已存在: %s", role.RoleKey)
	}

	// 创建角色
	if err := gorm.G[model.Role](s.db).Create(ctx, role); err != nil {
		s.logger.Error("创建角色失败", zap.Error(err))
		return fmt.Errorf("创建角色失败: %w", err)
	}

	s.logger.Info("创建角色成功", zap.Int64("roleId", role.ID), zap.String("roleKey", role.RoleKey))
	return nil
}

// Update 更新角色
func (s *roleService) Update(ctx context.Context, role *model.Role) error {
	// 检查角色是否存在
	existingRole, err := gorm.G[model.Role](s.db).Where("id = ?", role.ID).First(ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("角色不存在")
		}
		s.logger.Error("查询角色失败", zap.Error(err))
		return fmt.Errorf("查询角色失败: %w", err)
	}

	// 系统内置角色不允许修改角色标识
	if existingRole.IsSystem && existingRole.RoleKey != role.RoleKey {
		return fmt.Errorf("系统内置角色不允许修改角色标识")
	}

	// 更新角色
	updates := map[string]any{
		"role_name":  role.RoleName,
		"sort":       role.Sort,
		"status":     role.Status,
		"data_scope": role.DataScope,
		"remark":     role.Remark,
		"update_by":  role.UpdateBy,
	}

	if err := s.db.Model(&model.Role{}).Where("id = ?", role.ID).Updates(updates).Error; err != nil {
		s.logger.Error("更新角色失败", zap.Error(err))
		return fmt.Errorf("更新角色失败: %w", err)
	}

	s.logger.Info("更新角色成功", zap.Int64("roleId", role.ID))
	return nil
}

// Delete 删除角色
func (s *roleService) Delete(ctx context.Context, roleId int64) error {
	// 检查角色是否存在
	role, err := gorm.G[model.Role](s.db).Where("id = ?", roleId).First(ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("角色不存在")
		}
		s.logger.Error("查询角色失败", zap.Error(err))
		return fmt.Errorf("查询角色失败: %w", err)
	}

	// 系统内置角色不允许删除
	if role.IsSystem {
		return fmt.Errorf("系统内置角色不允许删除")
	}

	// 检查是否有用户使用该角色
	userRoleCount, err := gorm.G[model.MUserRole](s.db).Where("role_id = ?", roleId).Count(ctx, "id")
	if err != nil {
		return fmt.Errorf("检查角色使用情况失败: %w", err)
	}
	if userRoleCount > 0 {
		return fmt.Errorf("该角色正在被使用，无法删除")
	}

	// 开启事务删除角色及相关数据
	return s.db.Transaction(func(tx *gorm.DB) error {
		// 删除角色菜单关联
		if _, err := gorm.G[model.MRoleMenu](tx).Where("role_id = ?", roleId).Delete(ctx); err != nil {
			return fmt.Errorf("删除角色菜单关联失败: %w", err)
		}
		if _, err := gorm.G[model.MRoleApiPermission](tx).Where("role_id = ?", roleId).Delete(ctx); err != nil {
			return fmt.Errorf("删除角色 API 权限关联失败: %w", err)
		}

		// 删除角色
		if _, err := gorm.G[model.Role](tx).Where("id = ?", roleId).Delete(ctx); err != nil {
			return fmt.Errorf("删除角色失败: %w", err)
		}

		s.logger.Info("删除角色成功", zap.Int64("roleId", roleId))
		return nil
	})
}

// GetById 根据ID查询角色
func (s *roleService) GetById(ctx context.Context, roleId int64) (*model.Role, error) {
	role, err := gorm.G[model.Role](s.db).Where("id = ?", roleId).First(ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("角色不存在")
		}
		s.logger.Error("查询角色失败", zap.Error(err))
		return nil, fmt.Errorf("查询角色失败: %w", err)
	}

	return &role, nil
}

// GetByRoleKey 根据角色标识查询角色
func (s *roleService) GetByRoleKey(ctx context.Context, roleKey string) (*model.Role, error) {
	role, err := gorm.G[model.Role](s.db).Where("role_key = ?", roleKey).First(ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("角色不存在")
		}
		s.logger.Error("查询角色失败", zap.Error(err))
		return nil, fmt.Errorf("查询角色失败: %w", err)
	}

	return &role, nil
}

// Page 分页查询角色列表
func (s *roleService) Page(ctx context.Context, pageNum, pageSize int, roleName string, status int32) (*pagination.Page[model.Role], error) {
	query := s.db.Model(&model.Role{})

	// 条件查询
	if roleName != "" {
		query = query.Where("role_name LIKE ?", "%"+roleName+"%")
	}
	if status >= 0 {
		query = query.Where("status = ?", status)
	}

	// 构建 PageQuery
	pageQuery := &pagination.PageQuery{
		PageNum:  pageNum,
		PageSize: pageSize,
	}

	// 使用 Paginator 进行分页
	page, err := pagination.New[model.Role](query, pageQuery).Find()
	if err != nil {
		s.logger.Error("分页查询角色列表失败", zap.Error(err))
		return nil, fmt.Errorf("分页查询角色列表失败: %w", err)
	}

	return page, nil
}

// AssignRoleToUser 为用户分配角色
func (s *roleService) AssignRoleToUser(ctx context.Context, userId, roleId int64) error {
	// 使用事务确保数据一致性
	return s.db.Transaction(func(tx *gorm.DB) error {
		// 1. 检查用户是否存在
		if _, err := gorm.G[model.User](tx).Where("id = ?", userId).First(ctx); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("用户不存在")
			}
			return fmt.Errorf("查询用户失败: %w", err)
		}

		// 2. 检查角色是否存在
		role, err := gorm.G[model.Role](tx).Where("id = ?", roleId).First(ctx)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("角色不存在")
			}
			return fmt.Errorf("查询角色失败: %w", err)
		}

		// 3. 检查是否已分配
		count, err := gorm.G[model.MUserRole](tx).
			Where("user_id = ? AND role_id = ?", userId, roleId).
			Count(ctx, "id")
		if err != nil {
			return fmt.Errorf("检查用户角色关系失败: %w", err)
		}
		if count > 0 {
			return fmt.Errorf("用户已拥有该角色")
		}

		// 4. 创建用户角色关联
		userRole := &model.MUserRole{
			UserId: userId,
			RoleId: roleId,
		}
		if err := gorm.G[model.MUserRole](tx).Create(ctx, userRole); err != nil {
			s.logger.Error("分配用户角色失败", zap.Error(err))
			return fmt.Errorf("分配用户角色失败: %w", err)
		}

		s.logger.Info("为用户分配角色成功",
			zap.Int64("userId", userId),
			zap.Int64("roleId", roleId),
			zap.String("roleKey", role.RoleKey))

		return nil
	})
}

// RemoveRoleFromUser 移除用户的角色
func (s *roleService) RemoveRoleFromUser(ctx context.Context, userId, roleId int64) error {
	// 使用事务确保数据一致性
	return s.db.Transaction(func(tx *gorm.DB) error {
		// 1. 检查角色是否存在
		role, err := gorm.G[model.Role](tx).Where("id = ?", roleId).First(ctx)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("角色不存在")
			}
			return fmt.Errorf("查询角色失败: %w", err)
		}

		// 2. 从数据库移除用户角色关联
		rowsAffected, err := gorm.G[model.MUserRole](tx).Where("user_id = ? AND role_id = ?", userId, roleId).Delete(ctx)

		if err != nil {
			s.logger.Error("移除用户角色失败", zap.Error(err))
			return fmt.Errorf("移除用户角色失败: %w", err)
		}

		if rowsAffected == 0 {
			return fmt.Errorf("用户角色关系不存在")
		}

		s.logger.Info("移除用户角色成功",
			zap.Int64("userId", userId),
			zap.Int64("roleId", roleId),
			zap.String("roleKey", role.RoleKey))

		return nil
	})
}

// GetUserRoles 获取用户的所有角色
func (s *roleService) GetUserRoles(ctx context.Context, userId int64) ([]model.Role, error) {
	var roles []model.Role

	err := s.db.Table("s_role r").
		Joins("INNER JOIN m_user_role ur ON r.id = ur.role_id").
		Where("ur.user_id = ? AND r.status = 0", userId).
		Order("r.sort ASC").
		Find(&roles).Error

	if err != nil {
		s.logger.Error("查询用户角色失败", zap.Error(err))
		return nil, fmt.Errorf("查询用户角色失败: %w", err)
	}

	return roles, nil
}

// GetRoleUsers 获取拥有该角色的用户
func (s *roleService) GetRoleUsers(ctx context.Context, roleId int64) ([]model.User, error) {
	if _, err := gorm.G[model.Role](s.db).Where("id = ?", roleId).First(ctx); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("角色不存在")
		}
		return nil, fmt.Errorf("查询角色失败: %w", err)
	}

	var users []model.User
	err := s.db.Table("s_user u").
		Joins("INNER JOIN m_user_role ur ON u.id = ur.user_id").
		Where("ur.role_id = ?", roleId).
		Order("u.created_time DESC").
		Find(&users).Error
	if err != nil {
		s.logger.Error("查询角色用户失败", zap.Error(err))
		return nil, fmt.Errorf("查询角色用户失败: %w", err)
	}

	return users, nil
}

// AssignUsersToRole 批量为角色添加用户
func (s *roleService) AssignUsersToRole(ctx context.Context, roleId int64, userIds []int64, operatorId int64) error {
	userIds = uniqueInt64(userIds)
	if len(userIds) == 0 {
		return nil
	}

	_, err := gorm.G[model.Role](s.db).Where("id = ?", roleId).First(ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("角色不存在")
		}
		return fmt.Errorf("查询角色失败: %w", err)
	}

	var existingUsers []model.User
	if err := s.db.WithContext(ctx).
		Where("id IN ?", userIds).
		Find(&existingUsers).Error; err != nil {
		return fmt.Errorf("查询用户失败: %w", err)
	}
	if len(existingUsers) != len(userIds) {
		return fmt.Errorf("部分用户不存在")
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		var rows []model.MUserRole
		if err := tx.WithContext(ctx).
			Where("role_id = ? AND user_id IN ?", roleId, userIds).
			Find(&rows).Error; err != nil {
			return fmt.Errorf("查询用户角色关系失败: %w", err)
		}

		existing := make(map[int64]struct{}, len(rows))
		for _, row := range rows {
			existing[row.UserId] = struct{}{}
		}

		toCreate := make([]model.MUserRole, 0, len(userIds)-len(existing))
		for _, userId := range userIds {
			if _, ok := existing[userId]; ok {
				continue
			}
			toCreate = append(toCreate, model.MUserRole{
				UserId:   userId,
				RoleId:   roleId,
				CreateBy: operatorId,
				UpdateBy: operatorId,
			})
		}

		if len(toCreate) == 0 {
			return nil
		}

		if err := tx.Create(&toCreate).Error; err != nil {
			return fmt.Errorf("批量分配角色用户失败: %w", err)
		}

		s.logger.Info("批量为角色添加用户成功",
			zap.Int64("roleId", roleId),
			zap.Int("userCount", len(toCreate)))

		return nil
	})
}

// RemoveUsersFromRole 批量移除角色下的用户
func (s *roleService) RemoveUsersFromRole(ctx context.Context, roleId int64, userIds []int64) error {
	userIds = uniqueInt64(userIds)
	if len(userIds) == 0 {
		return nil
	}

	_, err := gorm.G[model.Role](s.db).Where("id = ?", roleId).First(ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("角色不存在")
		}
		return fmt.Errorf("查询角色失败: %w", err)
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.WithContext(ctx).
			Where("role_id = ? AND user_id IN ?", roleId, userIds).
			Delete(&model.MUserRole{}).Error; err != nil {
			return fmt.Errorf("批量移除角色用户失败: %w", err)
		}

		s.logger.Info("批量移除角色用户成功",
			zap.Int64("roleId", roleId),
			zap.Int("userCount", len(userIds)))

		return nil
	})
}

// AssignMenusToRole 为角色分配菜单权限
func (s *roleService) AssignMenusToRole(ctx context.Context, roleId int64, menuIds []int64) error {
	// 检查角色是否存在
	_, err := gorm.G[model.Role](s.db).Where("id = ?", roleId).First(ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("角色不存在")
		}
		return fmt.Errorf("查询角色失败: %w", err)
	}
	menuIds = uniqueInt64(menuIds)
	permissionIDs, err := s.menuAPIPermissionIDs(ctx, menuIds)
	if err != nil {
		return err
	}

	// 开启事务
	if err := s.db.Transaction(func(tx *gorm.DB) error {
		// 删除旧的菜单权限
		if err := tx.WithContext(ctx).Where("role_id = ?", roleId).Delete(&model.MRoleMenu{}).Error; err != nil {
			return fmt.Errorf("删除旧菜单权限失败: %w", err)
		}

		// 添加新的菜单权限
		if len(menuIds) > 0 {
			roleMenus := make([]model.MRoleMenu, 0, len(menuIds))
			for _, menuId := range menuIds {
				roleMenus = append(roleMenus, model.MRoleMenu{
					RoleId: roleId,
					MenuId: menuId,
				})
			}
			if err := tx.Create(&roleMenus).Error; err != nil {
				return fmt.Errorf("添加新菜单权限失败: %w", err)
			}
		}

		if err := s.replaceDerivedAPIPermissions(ctx, tx, roleId, permissionIDs); err != nil {
			return err
		}

		s.logger.Info("为角色分配菜单权限成功",
			zap.Int64("roleId", roleId),
			zap.Int("menuCount", len(menuIds)))

		return nil
	}); err != nil {
		return err
	}

	return nil
}

func (s *roleService) menuAPIPermissionIDs(ctx context.Context, menuIds []int64) ([]int64, error) {
	if len(menuIds) == 0 {
		return nil, nil
	}
	var ids []int64
	if err := s.db.WithContext(ctx).
		Table("m_menu_api_permission map").
		Distinct("map.permission_id").
		Joins("JOIN s_api_permission p ON p.id = map.permission_id AND p.status = 0").
		Where("map.menu_id IN ?", menuIds).
		Pluck("map.permission_id", &ids).Error; err != nil {
		return nil, fmt.Errorf("查询菜单关联 API 权限失败: %w", err)
	}
	return ids, nil
}

func (s *roleService) replaceDerivedAPIPermissions(ctx context.Context, tx *gorm.DB, roleId int64, permissionIDs []int64) error {
	if err := tx.WithContext(ctx).
		Where("role_id = ? AND source = ?", roleId, model.PermissionSourceDerived).
		Delete(&model.MRoleApiPermission{}).Error; err != nil {
		return fmt.Errorf("删除旧菜单派生 API 权限失败: %w", err)
	}
	rows := make([]model.MRoleApiPermission, 0, len(permissionIDs))
	for _, permissionID := range permissionIDs {
		rows = append(rows, model.MRoleApiPermission{RoleId: roleId, PermissionId: permissionID, Source: model.PermissionSourceDerived})
	}
	if len(rows) > 0 {
		if err := tx.WithContext(ctx).Create(&rows).Error; err != nil {
			return fmt.Errorf("保存菜单派生 API 权限失败: %w", err)
		}
	}
	return nil
}

// GetRoleMenus 获取角色的所有菜单
func (s *roleService) GetRoleMenus(ctx context.Context, roleId int64) ([]model.Menu, error) {
	var menus []model.Menu

	err := s.db.Table("s_menu m").
		Joins("INNER JOIN m_role_menu rm ON m.id = rm.menu_id").
		Where("rm.role_id = ? AND m.status = 0", roleId).
		Group("m.id").
		Order("m.sort ASC").
		Find(&menus).Error

	if err != nil {
		s.logger.Error("查询角色菜单失败", zap.Error(err))
		return nil, fmt.Errorf("查询角色菜单失败: %w", err)
	}

	return menus, nil
}
