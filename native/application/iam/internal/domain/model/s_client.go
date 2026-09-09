package model

import (
	"strings"

	"github.com/gcc798/lightning/internal/utils"

	"gorm.io/gorm"
)

// AuthClient 客户端配置表
type AuthClient struct {
	ClientId      string          `gorm:"column:client_id;type:varchar(64);primaryKey;comment:客户端ID" json:"clientId"`
	GrantType     string          `gorm:"column:grant_type;type:varchar(255);comment:授权类型(逗号分隔)" json:"grantType"`
	DeviceType    string          `gorm:"column:device_type;type:varchar(32);comment:设备类型" json:"deviceType"`
	Status        int             `gorm:"column:status;type:smallint;default:0;comment:状态(0正常 1停用)" json:"status"`
	Timeout       int64           `gorm:"column:timeout;type:bigint;default:604800;comment:固定超时时间(秒),默认7天" json:"timeout"`
	ActiveTimeout int64           `gorm:"column:active_timeout;type:bigint;default:1800;comment:活动超时时间(秒),默认30分钟" json:"activeTimeout"`
	Remark        string          `gorm:"column:remark;type:varchar(500);comment:备注" json:"remark"`
	CreateBy      int64           `gorm:"column:create_by;type:bigint;comment:创建者" json:"createBy"`
	CreatedTime   utils.LocalTime `gorm:"column:created_time;type:timestamptz;autoCreateTime;comment:创建时间" json:"createdTime"`
	UpdateBy      int64           `gorm:"column:update_by;type:bigint;comment:更新者" json:"updateBy"`
	UpdatedTime   utils.LocalTime `gorm:"column:updated_time;type:timestamptz;autoUpdateTime;comment:更新时间" json:"updatedTime"`
}

// TableName 指定表名
func (*AuthClient) TableName() string {
	return "s_auth_client"
}

// BeforeCreate 设置客户端超时默认值。客户端 ID 必须由调用方显式提供。
func (c *AuthClient) BeforeCreate(_ *gorm.DB) error {
	if c.Timeout == 0 {
		c.Timeout = 604800 // 7天
	}
	if c.ActiveTimeout == 0 {
		c.ActiveTimeout = 1800 // 30分钟
	}
	return nil
}

// IsGrantTypeSupported 检查是否支持指定的授权类型
func (c *AuthClient) IsGrantTypeSupported(grantType string) bool {
	if c.GrantType == "" {
		return false
	}
	types := strings.Split(c.GrantType, ",")
	for _, t := range types {
		if strings.TrimSpace(t) == grantType {
			return true
		}
	}
	return false
}

// IsActive 检查客户端是否启用
func (c *AuthClient) IsActive() bool {
	return c.Status == 0
}

// FindByClientId 根据clientId查询
func (*AuthClient) FindByClientId(db *gorm.DB, clientId string) (*AuthClient, error) {
	var client AuthClient
	err := db.Where("client_id = ?", clientId).First(&client).Error
	if err != nil {
		return nil, err
	}
	return &client, nil
}

// Create 创建客户端
func (c *AuthClient) Create(db *gorm.DB) error {
	return db.Create(c).Error
}

// Update 更新客户端
func (c *AuthClient) Update(db *gorm.DB) error {
	return db.Model(&AuthClient{}).Where("client_id = ?", c.ClientId).Updates(c).Error
}

// Delete 删除客户端
func (*AuthClient) Delete(db *gorm.DB, clientId string) error {
	return db.Where("client_id = ?", clientId).Delete(&AuthClient{}).Error
}

// List 分页查询客户端列表
func (*AuthClient) List(db *gorm.DB, pageNum, pageSize int, status *int, clientID string) ([]AuthClient, int64, error) {
	var clients []AuthClient
	var total int64

	query := db.Model(&AuthClient{})

	// 条件查询
	if status != nil {
		query = query.Where("status = ?", *status)
	}
	if clientID != "" {
		query = query.Where("client_id LIKE ?", "%"+clientID+"%")
	}

	// 统计总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (pageNum - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Order("created_time DESC").Find(&clients).Error

	return clients, total, err
}
