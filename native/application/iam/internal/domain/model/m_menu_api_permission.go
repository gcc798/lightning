package model

import "github.com/gcc798/microservice-kit/internal/utils"

// MMenuApiPermission declares which API permissions a menu requires.
type MMenuApiPermission struct {
	ID           int64           `gorm:"column:id;type:bigint;primaryKey;autoIncrement:false" autogen:"int64" json:"id"`
	MenuId       int64           `gorm:"column:menu_id;type:bigint;not null;uniqueIndex:idx_menu_api_permission" json:"menuId"`
	PermissionId int64           `gorm:"column:permission_id;type:bigint;not null;uniqueIndex:idx_menu_api_permission" json:"permissionId"`
	CreateBy     int64           `gorm:"column:create_by;type:bigint" json:"createBy"`
	UpdateBy     int64           `gorm:"column:update_by;type:bigint" json:"updateBy"`
	CreatedTime  utils.LocalTime `gorm:"column:created_time;type:timestamptz;autoCreateTime" json:"createdTime"`
	UpdatedTime  utils.LocalTime `gorm:"column:updated_time;type:timestamptz;autoUpdateTime" json:"updatedTime"`
}

func (*MMenuApiPermission) TableName() string { return "m_menu_api_permission" }
