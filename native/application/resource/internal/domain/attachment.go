// Package resource 提供附件与对象存储领域能力。
package resource

import (
	"encoding/json"

	"github.com/gcc798/microservice-kit/internal/utils"
)

// 附件状态。
const (
	StatusNormal  int32 = 0
	StatusDeleted int32 = 1
)

// Attachment 附件。字段的 json 标签是对外 HTTP 契约，改动会影响前端。
type Attachment struct {
	ID            int64            `gorm:"column:id;type:bigint;primaryKey;autoIncrement:false" autogen:"int64" json:"id"` // 使用分布式ID
	FileName      string           `gorm:"column:file_name;type:varchar(255);not null" json:"fileName"`                    // 原始文件名
	FileKey       string           `gorm:"column:file_key;type:varchar(512);not null;index" json:"fileKey"`                // 存储路径/Key
	FileSize      int64            `gorm:"column:file_size;type:bigint;not null" json:"fileSize"`                          // 文件大小（字节）
	FileType      string           `gorm:"column:file_type;type:varchar(128)" json:"fileType"`                             // 文件类型（MIME Type）
	FileExt       string           `gorm:"column:file_ext;type:varchar(32)" json:"fileExt"`                                // 文件扩展名
	BusinessType  string           `gorm:"column:business_type;type:varchar(64);index" json:"businessType"`                // 业务类型
	BusinessId    string           `gorm:"column:business_id;type:varchar(64);index" json:"businessId"`                    // 业务ID
	BusinessField string           `gorm:"column:business_field;type:varchar(64)" json:"businessField"`                    // 业务字段
	IsPublic      bool             `gorm:"column:is_public;type:boolean;default:false" json:"isPublic"`                    // 是否公开访问
	AccessUrl     string           `gorm:"column:access_url;type:varchar(1024)" json:"accessUrl"`                          // 访问URL
	Metadata      *json.RawMessage `gorm:"column:metadata;type:jsonb" json:"metadata"`                                     // JSON元数据
	Status        int32            `gorm:"column:status;type:smallint;default:0" json:"status"`                            // 状态：0正常 1已删除
	ExpireTime    utils.LocalTime  `gorm:"column:expire_time;type:timestamptz" json:"expireTime"`                          // 过期时间
	CreateBy      int64            `gorm:"column:create_by;type:bigint" json:"createBy"`                                   // 创建人
	CreateTime    utils.LocalTime  `gorm:"column:create_time;type:timestamptz;autoCreateTime" json:"createTime"`           // 创建时间
	UpdateTime    utils.LocalTime  `gorm:"column:update_time;type:timestamptz;autoUpdateTime" json:"updateTime"`           // 更新时间
}

// TableName 返回数据库表名。
func (*Attachment) TableName() string { return "biz_attachment" }
