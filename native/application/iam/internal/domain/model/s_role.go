package model

import (
	"github.com/gcc798/lightning/internal/utils"
	"gorm.io/gorm"
)

// Role 系统角色表
type Role struct {
	ID          int64           `gorm:"column:id;type:bigint;primaryKey;autoIncrement:false" autogen:"int64" json:"id"` // 角色ID（使用分布式ID）
	RoleKey     string          `gorm:"column:role_key;type:varchar(64);uniqueIndex;not null" json:"roleKey"`           // 角色标识（唯一，用于权限匹配）
	RoleName    string          `gorm:"column:role_name;type:varchar(64);not null" json:"roleName"`                     // 角色名称
	Sort        int64           `gorm:"column:sort;type:bigint;default:0" json:"sort"`                                  // 显示顺序
	Status      int32           `gorm:"column:status;type:smallint;default:0" json:"status"`                            // 状态：0正常 1停用
	DataScope   int32           `gorm:"column:data_scope;type:smallint;default:1" json:"dataScope"`                     // 数据范围：1全部 2自定义 3本组织 4本组织及以下 5仅本人
	IsSystem    bool            `gorm:"column:is_system;type:boolean;default:false" json:"isSystem"`                    // 是否系统内置角色（内置角色不可删除）
	Remark      string          `gorm:"column:remark;type:varchar(500)" json:"remark"`                                  // 备注
	CreateBy    int64           `gorm:"column:create_by;type:bigint" json:"createBy"`                                   // 创建人
	UpdateBy    int64           `gorm:"column:update_by;type:bigint" json:"updateBy"`                                   // 更新人
	CreatedTime utils.LocalTime `gorm:"column:created_time;type:timestamptz;autoCreateTime" json:"createdTime"`
	UpdatedTime utils.LocalTime `gorm:"column:updated_time;type:timestamptz;autoUpdateTime" json:"updatedTime"`
}

// TableName 返回数据库表名。
func (*Role) TableName() string { return "s_role" }

// FindByRoleKey 根据角色标识查询角色
func (r *Role) FindByRoleKey(db *gorm.DB, roleKey string) (*Role, error) {
	var role Role
	err := db.Where("role_key = ?", roleKey).First(&role).Error
	return &role, err
}

// FindByRoleId 根据角色ID查询角色
func (r *Role) FindByRoleId(db *gorm.DB, roleId int64) (*Role, error) {
	var role Role
	err := db.Where("id = ?", roleId).First(&role).Error
	return &role, err
}

// FindByID 根据ID查询角色
func (r *Role) FindByID(db *gorm.DB, roleId int64) (*Role, error) {
	var role Role
	err := db.Where("id = ?", roleId).First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

// CheckRoleKeyExists 检查角色标识是否存在
func (r *Role) CheckRoleKeyExists(db *gorm.DB, roleKey string) (bool, error) {
	var count int64
	err := db.Model(&Role{}).Where("role_key = ?", roleKey).Count(&count).Error
	return count > 0, err
}

// CountUserRoles 统计使用该角色的用户数量
func (r *Role) CountUserRoles(db *gorm.DB, roleId int64) (int64, error) {
	var count int64
	err := db.Model(&MUserRole{}).Where("role_id = ?", roleId).Count(&count).Error
	return count, err
}

// Create 创建角色
func (r *Role) Create(db *gorm.DB, role *Role) error {
	return db.Create(role).Error
}

// Update 更新角色
func (r *Role) Update(db *gorm.DB, roleId int64, updates map[string]any) error {
	return db.Model(&Role{}).Where("id = ?", roleId).Updates(updates).Error
}

// Delete 删除角色
func (r *Role) Delete(db *gorm.DB, roleId int64) error {
	return db.Where("id = ?", roleId).Delete(&Role{}).Error
}

// IsSystemRole 判断是否为系统内置角色
func (r *Role) IsSystemRole() bool {
	return r.IsSystem
}

// IsActiveRole 判断角色是否激活
func (r *Role) IsActiveRole() bool {
	return r.Status == 0
}

// FindByUserId 根据用户ID查询角色列表
func (r *Role) FindByUserId(db *gorm.DB, userId int64) ([]Role, error) {
	var roles []Role
	err := db.Table("s_role r").
		Joins("INNER JOIN m_user_role ur ON r.id = ur.role_id").
		Where("ur.user_id = ? AND r.status = 0", userId).
		Order("r.sort ASC").
		Find(&roles).Error
	return roles, err
}

// List 分页查询角色列表
func (r *Role) List(db *gorm.DB, offset, limit int, roleName string, status int32) ([]Role, int64, error) {
	var roles []Role
	var total int64

	query := db.Model(&Role{})

	// 条件过滤
	if roleName != "" {
		query = query.Where("role_name LIKE ?", "%"+roleName+"%")
	}
	if status >= 0 {
		query = query.Where("status = ?", status)
	}

	// 统计总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	if err := query.Offset(offset).Limit(limit).Order("sort ASC, created_time DESC").Find(&roles).Error; err != nil {
		return nil, 0, err
	}

	return roles, total, nil
}

// FindAll 查询所有角色
func (r *Role) FindAll(db *gorm.DB) ([]Role, error) {
	var roles []Role
	err := db.Where("status = 0").Order("sort ASC").Find(&roles).Error
	return roles, err
}

// BatchDelete 批量删除角色
func (r *Role) BatchDelete(db *gorm.DB, roleIds []int64) (int64, error) {
	result := db.Where("id IN ?", roleIds).Delete(&Role{})
	return result.RowsAffected, result.Error
}
