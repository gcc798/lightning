package model

import (
	"github.com/gcc798/lightning/internal/utils"
)

// ApiPermission API 权限资源树。
type ApiPermission struct {
	ID          int64           `gorm:"column:id;type:bigint;primaryKey;autoIncrement:false" autogen:"int64" json:"id"`
	ParentId    int64           `gorm:"column:parent_id;type:bigint;default:0;index" json:"parentId"`
	Module      string          `gorm:"column:module;type:varchar(64);not null;index" json:"module"`
	Code        string          `gorm:"column:code;type:varchar(128);uniqueIndex;not null" json:"code"`
	Name        string          `gorm:"column:name;type:varchar(64);not null" json:"name"`
	NodeType    int32           `gorm:"column:node_type;type:smallint;not null;default:2" json:"nodeType"`
	Action      string          `gorm:"column:action;type:varchar(32);not null;default:'*'" json:"action"`
	Method      string          `gorm:"column:method;type:varchar(16)" json:"method"`
	Path        string          `gorm:"column:path;type:varchar(255)" json:"path"`
	Sort        int64           `gorm:"column:sort;type:bigint;default:0" json:"sort"`
	Status      int32           `gorm:"column:status;type:smallint;default:0" json:"status"`
	Remark      string          `gorm:"column:remark;type:varchar(500)" json:"remark"`
	CreateBy    int64           `gorm:"column:create_by;type:bigint" json:"createBy"`
	UpdateBy    int64           `gorm:"column:update_by;type:bigint" json:"updateBy"`
	CreatedTime utils.LocalTime `gorm:"column:created_time;type:timestamptz;autoCreateTime" json:"createdTime"`
	UpdatedTime utils.LocalTime `gorm:"column:updated_time;type:timestamptz;autoUpdateTime" json:"updatedTime"`
}

func (*ApiPermission) TableName() string { return "s_api_permission" }
