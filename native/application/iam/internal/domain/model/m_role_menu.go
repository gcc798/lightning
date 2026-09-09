package model

import (
	"github.com/gcc798/lightning/internal/utils"
	"gorm.io/gorm"
)

// MRoleMenu 角色菜单权限关联表（映射表）
type MRoleMenu struct {
	Id          int64           `gorm:"column:id;type:bigint;primaryKey;autoIncrement:false" autogen:"int64" json:"id"` // 使用分布式ID
	RoleId      int64           `gorm:"column:role_id;type:bigint;not null;index:idx_role_menu" json:"roleId"`          // 角色ID
	MenuId      int64           `gorm:"column:menu_id;type:bigint;not null;index:idx_role_menu" json:"menuId"`          // 菜单ID
	CreateBy    int64           `gorm:"column:create_by;type:bigint" json:"createBy"`                                   // 创建人
	UpdateBy    int64           `gorm:"column:update_by;type:bigint" json:"updateBy"`                                   // 更新人
	CreatedTime utils.LocalTime `gorm:"column:created_time;type:timestamptz;autoCreateTime" json:"createdTime"`
	UpdatedTime utils.LocalTime `gorm:"column:updated_time;type:timestamptz;autoUpdateTime" json:"updatedTime"`
}

// TableName 返回数据库表名。
func (*MRoleMenu) TableName() string { return "m_role_menu" }

// FindByRoleId 根据角色ID查询菜单关联
func (m *MRoleMenu) FindByRoleId(db *gorm.DB, roleId int64) ([]MRoleMenu, error) {
	var roleMenus []MRoleMenu
	err := db.Where("role_id = ?", roleId).Find(&roleMenus).Error
	return roleMenus, err
}

// FindByMenuId 根据菜单ID查询角色关联
func (m *MRoleMenu) FindByMenuId(db *gorm.DB, menuId int64) ([]MRoleMenu, error) {
	var roleMenus []MRoleMenu
	err := db.Where("menu_id = ?", menuId).Find(&roleMenus).Error
	return roleMenus, err
}

// Exists 检查角色菜单关联是否存在
func (m *MRoleMenu) Exists(db *gorm.DB, roleId, menuId int64) (bool, error) {
	var count int64
	err := db.Model(&MRoleMenu{}).
		Where("role_id = ? AND menu_id = ?", roleId, menuId).
		Count(&count).Error
	return count > 0, err
}

// Create 创建角色菜单关联
func (m *MRoleMenu) Create(db *gorm.DB, roleMenu *MRoleMenu) error {
	return db.Create(roleMenu).Error
}

// BatchCreate 批量创建角色菜单关联
func (m *MRoleMenu) BatchCreate(db *gorm.DB, roleMenus []MRoleMenu) error {
	if len(roleMenus) == 0 {
		return nil
	}
	return db.Create(&roleMenus).Error
}

// Delete 删除角色菜单关联
func (m *MRoleMenu) Delete(db *gorm.DB, roleId, menuId int64) error {
	return db.Where("role_id = ? AND menu_id = ?", roleId, menuId).Delete(&MRoleMenu{}).Error
}

// DeleteByRoleId 删除角色的所有菜单关联
func (m *MRoleMenu) DeleteByRoleId(db *gorm.DB, roleId int64) error {
	return db.Where("role_id = ?", roleId).Delete(&MRoleMenu{}).Error
}

// DeleteByMenuId 删除菜单的所有角色关联
func (m *MRoleMenu) DeleteByMenuId(db *gorm.DB, menuId int64) error {
	return db.Where("menu_id = ?", menuId).Delete(&MRoleMenu{}).Error
}

// CountByMenuId 统计使用该菜单的角色数量
func (m *MRoleMenu) CountByMenuId(db *gorm.DB, menuId int64) (int64, error) {
	var count int64
	err := db.Model(&MRoleMenu{}).Where("menu_id = ?", menuId).Count(&count).Error
	return count, err
}
