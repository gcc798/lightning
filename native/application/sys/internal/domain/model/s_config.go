package model

import (
	"encoding/json"

	"github.com/gcc798/lightning/internal/utils"
	"gorm.io/gorm"
)

// Config 配置表
type Config struct {
	ID          int64           `gorm:"column:id;type:bigint;primaryKey;autoIncrement:false" autogen:"int64" json:"id"`   // 配置ID（使用分布式ID）
	Name        string          `gorm:"column:name;type:varchar(128);not null" json:"name"`                               // 配置名称
	Code        string          `gorm:"column:code;type:varchar(128);not null;uniqueIndex:idx_s_config_code" json:"code"` // 配置编码
	Data        json.RawMessage `gorm:"column:data;type:jsonb;not null" json:"data"`                                      // 配置数据（JSON格式）
	Remark      string          `gorm:"column:remark;type:varchar(500)" json:"remark"`                                    // 备注
	CreateBy    int64           `gorm:"column:create_by;type:bigint" json:"createBy"`                                     // 创建者
	CreatedTime utils.LocalTime `gorm:"column:created_time;type:timestamptz;autoCreateTime" json:"createdTime"`           // 创建时间
	UpdateBy    int64           `gorm:"column:update_by;type:bigint" json:"updateBy"`                                     // 更新者
	UpdatedTime utils.LocalTime `gorm:"column:updated_time;type:timestamptz;autoUpdateTime" json:"updatedTime"`           // 更新时间
}

// TableName 返回数据库表名。
func (*Config) TableName() string {
	return "s_config"
}

// FindByID 根据ID查询配置
func (*Config) FindByID(db *gorm.DB, id int64) (*Config, error) {
	var config Config
	err := db.Where("id = ?", id).First(&config).Error
	if err != nil {
		return nil, err
	}
	return &config, nil
}

// FindByCode 根据唯一配置编码查询配置。
func (*Config) FindByCode(db *gorm.DB, configCode string) (*Config, error) {
	var config Config
	if err := db.Where("code = ?", configCode).First(&config).Error; err != nil {
		return nil, err
	}
	return &config, nil
}

// CheckNameExists 检查配置名称是否存在
func (*Config) CheckNameExists(db *gorm.DB, name string) (bool, error) {
	var count int64
	err := db.Model(&Config{}).Where("name = ?", name).Count(&count).Error
	return count > 0, err
}

// CheckNameExistsExcludingSelf 检查配置名称是否被其他配置占用
func (*Config) CheckNameExistsExcludingSelf(db *gorm.DB, id int64, name string) (bool, error) {
	var count int64
	err := db.Model(&Config{}).
		Where("name = ? AND id != ?", name, id).
		Count(&count).Error
	return count > 0, err
}

// Create 创建配置
func (c *Config) Create(db *gorm.DB) error {
	return db.Create(c).Error
}

// Update 更新配置
func (c *Config) Update(db *gorm.DB, id int64, updates map[string]interface{}) error {
	return db.Model(&Config{}).Where("id = ?", id).Updates(updates).Error
}

// Delete 删除配置
func (*Config) Delete(db *gorm.DB, id int64) error {
	return db.Where("id = ?", id).Delete(&Config{}).Error
}

// GetDataByCode 根据唯一编码获取配置数据。
func (*Config) GetDataByCode(db *gorm.DB, configCode string) (json.RawMessage, error) {
	var config Config
	err := db.Where("code = ?", configCode).First(&config).Error
	if err != nil {
		return nil, err
	}
	return config.Data, nil
}
