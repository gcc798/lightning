package iam

import (
	"errors"
	"fmt"
	"sort"

	"github.com/gcc798/lightning/application/iam/internal/domain/model"
	"gorm.io/gorm"
)

// MenuService 定义业务数据结构。
type MenuService struct {
	db *gorm.DB
}

// NewMenuService 创建组件实例。
func NewMenuService(db *gorm.DB) *MenuService {
	return &MenuService{
		db: db,
	}
}

// MenuTree 菜单树节点
type MenuTree struct {
	model.Menu
	Children []*MenuTree `json:"children,omitempty"`
}

// GetUserMenuTree 获取用户的菜单树（用于前端路由生成）
func (s *MenuService) GetUserMenuTree(userId int64) ([]*MenuTree, error) {
	// 1. 获取用户的角色列表
	roleModel := &model.Role{}
	roles, err := roleModel.FindByUserId(s.db, userId)
	if err != nil {
		return nil, fmt.Errorf("failed to get user roles: %w", err)
	}

	if len(roles) == 0 {
		return []*MenuTree{}, nil
	}

	// 2. 检查是否是超级管理员
	isSuperAdmin := false
	roleIds := make([]int64, len(roles))
	for i, role := range roles {
		roleIds[i] = role.ID
		// 超级管理员的 role_key 是 super_admin
		if role.RoleKey == "super_admin" {
			isSuperAdmin = true
		}
	}

	var menus []model.Menu
	menuModel := &model.Menu{}

	// 3. 如果是超级管理员，获取所有菜单
	if isSuperAdmin {
		menus, err = menuModel.FindAll(s.db)
		if err != nil {
			return nil, fmt.Errorf("failed to get all menus: %w", err)
		}
	} else {
		// 普通用户，根据角色获取菜单
		menus, err = menuModel.FindByRoleIds(s.db, roleIds)
		if err != nil {
			return nil, fmt.Errorf("failed to get menus by roles: %w", err)
		}
	}

	// 4. 过滤停用/隐藏的菜单，但保留按钮类型（前端需要按钮的权限标识）
	var filteredMenus []model.Menu
	for _, menu := range menus {
		// 保留所有状态正常的菜单（包括按钮类型）
		// 按钮类型的菜单虽然不生成路由，但其 perms 字段用于权限控制
		if menu.Status == 0 {
			filteredMenus = append(filteredMenus, menu)
		}
	}

	// 5. 构建菜单树
	filteredMenus, err = s.withAncestorMenus(filteredMenus)
	if err != nil {
		return nil, err
	}
	return s.buildMenuTree(filteredMenus, 0), nil
}

// GetAllMenuTree 获取所有菜单树（用于菜单管理页面）
func (s *MenuService) GetAllMenuTree() ([]*MenuTree, error) {
	menuModel := &model.Menu{}
	menus, err := menuModel.FindAll(s.db)
	if err != nil {
		return nil, err
	}
	if err := s.loadAPIPermissionIDs(menus); err != nil {
		return nil, err
	}
	return s.buildMenuTree(menus, 0), nil
}

// GetMenuList 获取菜单列表
func (s *MenuService) GetMenuList() ([]model.Menu, error) {
	menuModel := &model.Menu{}
	menus, err := menuModel.FindAll(s.db)
	if err != nil {
		return nil, err
	}
	return menus, s.loadAPIPermissionIDs(menus)
}

// GetMenuById 根据ID获取菜单
func (s *MenuService) GetMenuById(menuId int64) (*model.Menu, error) {
	menu, err := (&model.Menu{}).FindByID(s.db, menuId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("菜单不存在")
		}
		return nil, fmt.Errorf("查询菜单失败: %w", err)
	}
	var ids []int64
	if err := s.db.Model(&model.MMenuApiPermission{}).Where("menu_id = ?", menuId).Pluck("permission_id", &ids).Error; err != nil {
		return nil, fmt.Errorf("查询菜单 API 权限失败: %w", err)
	}
	menu.ApiPermissionIds = ids
	return menu, nil
}

// CreateMenu 创建菜单
func (s *MenuService) CreateMenu(menu *model.Menu) error {
	// 场景规则：创建时 ID 必须为空
	if menu.ID != 0 {
		return errors.New("创建时不能指定菜单ID")
	}

	// 场景规则：验证父菜单是否存在
	if menu.ParentId != 0 {
		parent, err := (&model.Menu{}).FindByID(s.db, menu.ParentId)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("父菜单不存在")
			}
			return fmt.Errorf("查询父菜单失败: %w", err)
		}

		// 使用领域规则检查是否可以创建子菜单
		if !parent.CanHaveChild(menu.MenuType) {
			if parent.IsDirectory() && menu.IsButton() {
				return errors.New("目录下不能直接创建按钮")
			}
			if parent.IsMenu() && !menu.IsButton() {
				return errors.New("菜单下只能创建按钮")
			}
			return errors.New("不允许创建此类型的子菜单")
		}
	}

	// 场景规则：检查同级菜单名称唯一性
	exists, err := (&model.Menu{}).CheckMenuNameExists(s.db, menu.MenuName, menu.ParentId)
	if err != nil {
		return fmt.Errorf("检查菜单名称失败: %w", err)
	}
	if exists {
		return errors.New("同级菜单名称已存在")
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		if err := menu.Create(tx, menu); err != nil {
			return fmt.Errorf("创建菜单失败: %w", err)
		}
		return s.replaceMenuAPIPermissions(tx, menu)
	})
}

// UpdateMenu 更新菜单
func (s *MenuService) UpdateMenu(menu *model.Menu) error {
	// 场景规则：更新时 ID 必须不为空
	if menu.ID == 0 {
		return errors.New("更新时必须指定菜单ID")
	}

	// 场景规则：检查菜单是否存在
	existing, err := (&model.Menu{}).FindByID(s.db, menu.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("菜单不存在")
		}
		return fmt.Errorf("查询菜单失败: %w", err)
	}

	// 场景规则：不能将自己设置为父菜单
	if menu.ParentId == menu.ID {
		return errors.New("不能将自己设置为父菜单")
	}

	// 场景规则：验证父菜单是否存在
	if menu.ParentId != 0 {
		_, err := (&model.Menu{}).FindByID(s.db, menu.ParentId)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("父菜单不存在")
			}
			return fmt.Errorf("查询父菜单失败: %w", err)
		}
	}

	// 场景规则：检查同级菜单名称唯一性（排除自己）
	exists, err := (&model.Menu{}).CheckMenuNameExistsExcludingSelf(s.db, menu.ID, menu.MenuName, menu.ParentId)
	if err != nil {
		return fmt.Errorf("检查菜单名称失败: %w", err)
	}
	if exists {
		return errors.New("同级菜单名称已存在")
	}

	// 保留原有的创建信息
	menu.CreateBy = existing.CreateBy
	menu.CreatedTime = existing.CreatedTime

	return s.db.Transaction(func(tx *gorm.DB) error {
		if err := menu.Update(tx, menu); err != nil {
			return fmt.Errorf("更新菜单失败: %w", err)
		}
		return s.replaceMenuAPIPermissions(tx, menu)
	})
}

func (s *MenuService) replaceMenuAPIPermissions(tx *gorm.DB, menu *model.Menu) error {
	ids := uniqueInt64(menu.ApiPermissionIds)
	if len(ids) > 0 {
		var count int64
		if err := tx.Model(&model.ApiPermission{}).Where("id IN ?", ids).Count(&count).Error; err != nil {
			return fmt.Errorf("校验菜单 API 权限失败: %w", err)
		}
		if count != int64(len(ids)) {
			return errors.New("菜单包含不存在的 API 权限")
		}
	}
	if err := tx.Where("menu_id = ?", menu.ID).Delete(&model.MMenuApiPermission{}).Error; err != nil {
		return fmt.Errorf("删除旧菜单 API 权限失败: %w", err)
	}
	rows := make([]model.MMenuApiPermission, 0, len(ids))
	for _, id := range ids {
		rows = append(rows, model.MMenuApiPermission{MenuId: menu.ID, PermissionId: id, CreateBy: menu.UpdateBy, UpdateBy: menu.UpdateBy})
	}
	if len(rows) > 0 {
		if err := tx.Create(&rows).Error; err != nil {
			return fmt.Errorf("保存菜单 API 权限失败: %w", err)
		}
	}
	return nil
}

func (s *MenuService) loadAPIPermissionIDs(menus []model.Menu) error {
	if len(menus) == 0 {
		return nil
	}
	menuByID := make(map[int64]*model.Menu, len(menus))
	ids := make([]int64, 0, len(menus))
	for i := range menus {
		menuByID[menus[i].ID] = &menus[i]
		ids = append(ids, menus[i].ID)
	}
	var rows []model.MMenuApiPermission
	if err := s.db.Where("menu_id IN ?", ids).Find(&rows).Error; err != nil {
		return fmt.Errorf("查询菜单 API 权限失败: %w", err)
	}
	for _, row := range rows {
		menuByID[row.MenuId].ApiPermissionIds = append(menuByID[row.MenuId].ApiPermissionIds, row.PermissionId)
	}
	return nil
}

// DeleteMenu 删除菜单
func (s *MenuService) DeleteMenu(menuId int64) error {
	// 场景规则：ID 必须不为空
	if menuId == 0 {
		return errors.New("菜单ID不能为空")
	}

	// 查询菜单
	menu, err := (&model.Menu{}).FindByID(s.db, menuId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("菜单不存在")
		}
		return fmt.Errorf("查询菜单失败: %w", err)
	}

	// 检查是否有子菜单
	hasChildren, err := menu.HasChildren(s.db)
	if err != nil {
		return fmt.Errorf("检查子菜单失败: %w", err)
	}
	if hasChildren {
		return errors.New("存在子菜单，无法删除")
	}

	// 调用模型层的删除方法
	return s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("menu_id = ?", menuId).Delete(&model.MMenuApiPermission{}).Error; err != nil {
			return fmt.Errorf("删除菜单 API 权限失败: %w", err)
		}
		if err := menu.Delete(tx, menuId); err != nil {
			return fmt.Errorf("删除菜单失败: %w", err)
		}
		return nil
	})
}

// buildMenuTree 构建菜单树
func (s *MenuService) buildMenuTree(menus []model.Menu, parentId int64) []*MenuTree {
	var tree []*MenuTree

	for _, menu := range menus {
		if menu.ParentId == parentId {
			node := &MenuTree{
				Menu:     menu,
				Children: s.buildMenuTree(menus, menu.ID),
			}
			tree = append(tree, node)
		}
	}

	return tree
}

func (s *MenuService) withAncestorMenus(menus []model.Menu) ([]model.Menu, error) {
	if len(menus) == 0 {
		return menus, nil
	}

	menuByID := make(map[int64]model.Menu, len(menus))
	missingParentIDs := make([]int64, 0)
	for _, menu := range menus {
		menuByID[menu.ID] = menu
	}
	for _, menu := range menus {
		if menu.ParentId != 0 {
			if _, ok := menuByID[menu.ParentId]; !ok {
				missingParentIDs = append(missingParentIDs, menu.ParentId)
			}
		}
	}

	for len(missingParentIDs) > 0 {
		parentIDs := uniqueMenuIDs(missingParentIDs)
		missingParentIDs = nil

		var parents []model.Menu
		if err := s.db.
			Where("id IN ? AND status = 0", parentIDs).
			Find(&parents).Error; err != nil {
			return nil, fmt.Errorf("failed to get ancestor menus: %w", err)
		}

		for _, parent := range parents {
			if _, exists := menuByID[parent.ID]; exists {
				continue
			}
			menuByID[parent.ID] = parent
			if parent.ParentId != 0 {
				if _, ok := menuByID[parent.ParentId]; !ok {
					missingParentIDs = append(missingParentIDs, parent.ParentId)
				}
			}
		}
	}

	result := make([]model.Menu, 0, len(menuByID))
	for _, menu := range menuByID {
		result = append(result, menu)
	}
	sort.Slice(result, func(i, j int) bool {
		if result[i].Sort == result[j].Sort {
			return result[i].ID < result[j].ID
		}
		return result[i].Sort < result[j].Sort
	})
	return result, nil
}

func uniqueMenuIDs(ids []int64) []int64 {
	if len(ids) == 0 {
		return nil
	}

	seen := make(map[int64]struct{}, len(ids))
	result := make([]int64, 0, len(ids))
	for _, id := range ids {
		if id == 0 {
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		result = append(result, id)
	}
	return result
}
