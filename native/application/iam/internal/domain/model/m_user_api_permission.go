package model

import (
	"github.com/gcc798/lightning/internal/utils"
)

// MUserApiPermission 用户 API 权限关联表。
type MUserApiPermission struct {
	ID           int64           `gorm:"column:id;type:bigint;primaryKey;autoIncrement:false" autogen:"int64" json:"id"`
	UserId       int64           `gorm:"column:user_id;type:bigint;not null;uniqueIndex:idx_user_api_permission" json:"userId"`
	PermissionId int64           `gorm:"column:permission_id;type:bigint;not null;uniqueIndex:idx_user_api_permission" json:"permissionId"`
	Source       int32           `gorm:"column:source;type:smallint;not null;default:0;uniqueIndex:idx_user_api_permission" json:"source"`
	CreateBy     int64           `gorm:"column:create_by;type:bigint" json:"createBy"`
	UpdateBy     int64           `gorm:"column:update_by;type:bigint" json:"updateBy"`
	CreatedTime  utils.LocalTime `gorm:"column:created_time;type:timestamptz;autoCreateTime" json:"createdTime"`
	UpdatedTime  utils.LocalTime `gorm:"column:updated_time;type:timestamptz;autoUpdateTime" json:"updatedTime"`
}

func (*MUserApiPermission) TableName() string { return "m_user_api_permission" }
