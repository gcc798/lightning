package model

import (
	"github.com/gcc798/microservice-kit/internal/utils"
)

const (
	PermissionSourceManual  int32 = 0
	PermissionSourceDerived int32 = 1
)

// MRoleApiPermission 角色 API 权限关联表。
type MRoleApiPermission struct {
	ID           int64           `gorm:"column:id;type:bigint;primaryKey;autoIncrement:false" autogen:"int64" json:"id"`
	RoleId       int64           `gorm:"column:role_id;type:bigint;not null;uniqueIndex:idx_role_api_permission" json:"roleId"`
	PermissionId int64           `gorm:"column:permission_id;type:bigint;not null;uniqueIndex:idx_role_api_permission" json:"permissionId"`
	Source       int32           `gorm:"column:source;type:smallint;not null;default:0;uniqueIndex:idx_role_api_permission" json:"source"`
	CreateBy     int64           `gorm:"column:create_by;type:bigint" json:"createBy"`
	UpdateBy     int64           `gorm:"column:update_by;type:bigint" json:"updateBy"`
	CreatedTime  utils.LocalTime `gorm:"column:created_time;type:timestamptz;autoCreateTime" json:"createdTime"`
	UpdatedTime  utils.LocalTime `gorm:"column:updated_time;type:timestamptz;autoUpdateTime" json:"updatedTime"`
}

func (*MRoleApiPermission) TableName() string { return "m_role_api_permission" }
