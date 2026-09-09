package model

import (
	"github.com/gcc798/lightning/internal/utils"
	"gorm.io/gorm"
)

// MUserRole 用户角色关联表（映射表）
type MUserRole struct {
	Id          int64           `gorm:"column:id;type:bigint;primaryKey;autoIncrement:false" autogen:"int64" json:"id"` // 使用分布式ID
	UserId      int64           `gorm:"column:user_id;type:bigint;not null;uniqueIndex:idx_user_role" json:"userId"`    // 用户ID
	RoleId      int64           `gorm:"column:role_id;type:bigint;not null;uniqueIndex:idx_user_role" json:"roleId"`    // 角色ID
	CreateBy    int64           `gorm:"column:create_by;type:bigint" json:"createBy"`                                   // 创建人
	UpdateBy    int64           `gorm:"column:update_by;type:bigint" json:"updateBy"`                                   // 更新人
	CreatedTime utils.LocalTime `gorm:"column:created_time;type:timestamptz;autoCreateTime" json:"createdTime"`
	UpdatedTime utils.LocalTime `gorm:"column:updated_time;type:timestamptz;autoUpdateTime" json:"updatedTime"`
}

// TableName 返回数据库表名。
func (*MUserRole) TableName() string { return "m_user_role" }

// FindByUserId 根据用户ID查询角色关联
func (m *MUserRole) FindByUserId(db *gorm.DB, userId int64) ([]MUserRole, error) {
	var userRoles []MUserRole
	err := db.Where("user_id = ?", userId).Find(&userRoles).Error
	return userRoles, err
}

// FindByRoleId 根据角色ID查询用户关联
func (m *MUserRole) FindByRoleId(db *gorm.DB, roleId int64) ([]MUserRole, error) {
	var userRoles []MUserRole
	err := db.Where("role_id = ?", roleId).Find(&userRoles).Error
	return userRoles, err
}

// Exists 检查用户角色关联是否存在
func (m *MUserRole) Exists(db *gorm.DB, userId, roleId int64) (bool, error) {
	var count int64
	err := db.Model(&MUserRole{}).
		Where("user_id = ? AND role_id = ?", userId, roleId).
		Count(&count).Error
	return count > 0, err
}

// Create 创建用户角色关联
func (m *MUserRole) Create(db *gorm.DB, userRole *MUserRole) error {
	return db.Create(userRole).Error
}

// Delete 删除用户角色关联
func (m *MUserRole) Delete(db *gorm.DB, userId, roleId int64) error {
	return db.Where("user_id = ? AND role_id = ?", userId, roleId).Delete(&MUserRole{}).Error
}

// DeleteByUserId 删除用户的所有角色关联
func (m *MUserRole) DeleteByUserId(db *gorm.DB, userId int64) error {
	return db.Where("user_id = ?", userId).Delete(&MUserRole{}).Error
}

// DeleteByRoleId 删除角色的所有用户关联
func (m *MUserRole) DeleteByRoleId(db *gorm.DB, roleId int64) error {
	return db.Where("role_id = ?", roleId).Delete(&MUserRole{}).Error
}

// CountByRoleId 统计使用该角色的用户数量
func (m *MUserRole) CountByRoleId(db *gorm.DB, roleId int64) (int64, error) {
	var count int64
	err := db.Model(&MUserRole{}).Where("role_id = ?", roleId).Count(&count).Error
	return count, err
}
