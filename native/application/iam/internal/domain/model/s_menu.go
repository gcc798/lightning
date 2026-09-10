package model

import (
	"github.com/gcc798/microservice-kit/internal/utils"

	"gorm.io/gorm"
)

// Menu 系统菜单权限表
type Menu struct {
	ID               int64           `gorm:"column:id;type:bigint;primaryKey;autoIncrement:false" autogen:"int64" json:"id"` // 菜单ID（使用分布式ID）
	MenuName         string          `gorm:"column:menu_name;type:varchar(64);not null" json:"menuName"`                     // 菜单名称
	ParentId         int64           `gorm:"column:parent_id;type:bigint;default:0;index" json:"parentId"`                   // 父菜单ID（0表示根菜单）
	Sort             int64           `gorm:"column:sort;type:bigint;default:0" json:"sort"`                                  // 显示顺序
	Path             string          `gorm:"column:path;type:varchar(255)" json:"path"`                                      // 路由地址
	Component        string          `gorm:"column:component;type:varchar(255)" json:"component"`                            // 组件路径
	Query            string          `gorm:"column:query;type:varchar(255)" json:"query"`                                    // 路由参数
	IsFrame          int32           `gorm:"column:is_frame;type:smallint;default:0" json:"isFrame"`                         // 是否外链：0否 1是
	IsCache          int32           `gorm:"column:is_cache;type:smallint;default:0" json:"isCache"`                         // 是否缓存：0否 1是
	MenuType         int32           `gorm:"column:menu_type;type:smallint;not null" json:"menuType"`                        // 菜单类型：0目录 1菜单 2按钮
	Visible          int32           `gorm:"column:visible;type:smallint;default:0" json:"visible"`                          // 显示状态：0显示 1隐藏
	Status           int32           `gorm:"column:status;type:smallint;default:0" json:"status"`                            // 状态：0正常 1停用
	Perms            string          `gorm:"column:perms;type:varchar(255)" json:"perms"`                                    // 前端按钮显隐标识，不参与后端 API 权限派生
	ApiPermissionIds []int64         `gorm:"-" json:"apiPermissionIds"`                                                      // 菜单派生的后端 API 权限
	Icon             string          `gorm:"column:icon;type:varchar(64)" json:"icon"`                                       // 菜单图标
	Remark           string          `gorm:"column:remark;type:varchar(500)" json:"remark"`                                  // 备注
	CreateBy         int64           `gorm:"column:create_by;type:bigint" json:"createBy"`                                   // 创建人
	UpdateBy         int64           `gorm:"column:update_by;type:bigint" json:"updateBy"`                                   // 更新人
	CreatedTime      utils.LocalTime `gorm:"column:created_time;type:timestamptz;autoCreateTime" json:"createdTime"`
	UpdatedTime      utils.LocalTime `gorm:"column:updated_time;type:timestamptz;autoUpdateTime" json:"updatedTime"`
}

// TableName 返回数据库表名。
func (*Menu) TableName() string { return "s_menu" }

// FindByMenuId 根据菜单ID查询菜单
func (m *Menu) FindByMenuId(db *gorm.DB, menuId int64) (*Menu, error) {
	var menu Menu
	err := db.Where("id = ?", menuId).First(&menu).Error
	return &menu, err
}

// FindByID 根据ID查询菜单
func (m *Menu) FindByID(db *gorm.DB, menuId int64) (*Menu, error) {
	var menu Menu
	err := db.Where("id = ?", menuId).First(&menu).Error
	if err != nil {
		return nil, err
	}
	return &menu, nil
}

// FindByParentId 查询子菜单列表
func (m *Menu) FindByParentId(db *gorm.DB, parentId int64) ([]Menu, error) {
	var menus []Menu
	err := db.Where("parent_id = ?", parentId).Order("sort ASC").Find(&menus).Error
	return menus, err
}

// FindByRoleId 根据角色ID查询菜单列表
func (m *Menu) FindByRoleId(db *gorm.DB, roleId int64) ([]Menu, error) {
	var menus []Menu
	err := db.Table("s_menu m").
		Joins("INNER JOIN m_role_menu rm ON m.id = rm.menu_id").
		Where("rm.role_id = ? AND m.status = 0", roleId).
		Group("m.id").
		Order("m.sort ASC").
		Find(&menus).Error
	return menus, err
}

// FindAll 查询所有菜单
func (m *Menu) FindAll(db *gorm.DB) ([]Menu, error) {
	var menus []Menu
	err := db.Where("status = 0").Order("sort ASC").Find(&menus).Error
	return menus, err
}

// FindByRoleIds 根据多个角色ID查询菜单列表
func (m *Menu) FindByRoleIds(db *gorm.DB, roleIds []int64) ([]Menu, error) {
	var menus []Menu
	err := db.Table("s_menu m").
		Joins("INNER JOIN m_role_menu rm ON m.id = rm.menu_id").
		Where("rm.role_id IN ? AND m.status = 0", roleIds).
		Group("m.id").
		Order("m.sort ASC").
		Find(&menus).Error
	return menus, err
}

// CheckMenuNameExists 检查同级菜单名称是否存在
func (m *Menu) CheckMenuNameExists(db *gorm.DB, menuName string, parentId int64) (bool, error) {
	var count int64
	err := db.Model(&Menu{}).
		Where("menu_name = ? AND parent_id = ? AND status = 0", menuName, parentId).
		Count(&count).Error
	return count > 0, err
}

// CheckMenuNameExistsExcludingSelf 检查同级菜单名称是否被其他菜单占用
func (m *Menu) CheckMenuNameExistsExcludingSelf(db *gorm.DB, menuId int64, menuName string, parentId int64) (bool, error) {
	var count int64
	err := db.Model(&Menu{}).
		Where("menu_name = ? AND parent_id = ? AND id != ? AND status = 0", menuName, parentId, menuId).
		Count(&count).Error
	return count > 0, err
}

// CountChildren 统计子菜单数量
func (m *Menu) CountChildren(db *gorm.DB, menuId int64) (int64, error) {
	var count int64
	err := db.Model(&Menu{}).Where("parent_id = ?", menuId).Count(&count).Error
	return count, err
}

// Create 创建菜单
func (m *Menu) Create(db *gorm.DB, menu *Menu) error {
	return db.Create(menu).Error
}

// Update 更新菜单
func (m *Menu) Update(db *gorm.DB, menu *Menu) error {
	return db.Save(menu).Error
}

// Delete 删除菜单
func (m *Menu) Delete(db *gorm.DB, menuId int64) error {
	return db.Where("id = ?", menuId).Delete(&Menu{}).Error
}

// IsActive 判断菜单是否激活
func (m *Menu) IsActive() bool {
	return m.Status == 0
}

// IsVisible 判断菜单是否可见
func (m *Menu) IsVisible() bool {
	return m.Visible == 0
}

// IsDirectory 判断是否为目录
func (m *Menu) IsDirectory() bool {
	return m.MenuType == 0
}

// IsMenu 判断是否为菜单
func (m *Menu) IsMenu() bool {
	return m.MenuType == 1
}

// IsButton 判断是否为按钮
func (m *Menu) IsButton() bool {
	return m.MenuType == 2
}

// CanHaveChild 判断是否可以有子菜单
func (m *Menu) CanHaveChild(childType int32) bool {
	// 目录(0)下可以创建目录(0)或菜单(1)
	if m.MenuType == 0 && (childType == 0 || childType == 1) {
		return true
	}
	// 菜单(1)下只能创建按钮(2)
	if m.MenuType == 1 && childType == 2 {
		return true
	}
	return false
}

// HasChildren 判断是否有子菜单
func (m *Menu) HasChildren(db *gorm.DB) (bool, error) {
	count, err := m.CountChildren(db, m.ID)
	return count > 0, err
}
