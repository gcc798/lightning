package model

import (
	"time"

	"github.com/gcc798/lightning/internal/utils"
	"gorm.io/gorm"
)

// OperLog 操作日志
type OperLog struct {
	ID            int64           `gorm:"column:id;type:bigint;primaryKey;autoIncrement:false" autogen:"int64" json:"id"` // 日志ID（使用分布式ID）
	Title         string          `gorm:"column:title;type:varchar(128)" json:"title"`                                    // 模块标题
	BusinessType  string          `gorm:"column:business_type;type:varchar(32)" json:"businessType"`                      // 业务类型
	Method        string          `gorm:"column:method;type:varchar(255)" json:"method"`                                  // 调用方法
	RequestMethod string          `gorm:"column:request_method;type:varchar(16)" json:"requestMethod"`                    // 请求方式：GET/POST
	DeviceType    string          `gorm:"column:device_type;type:varchar(32)" json:"deviceType"`                          // 终端类型：web/ios/android/wechat
	OperName      string          `gorm:"column:oper_name;type:varchar(64)" json:"operName"`                              // 操作者
	OperUrl       string          `gorm:"column:oper_url;type:varchar(1024)" json:"operUrl"`                              // 请求URL
	OperIp        string          `gorm:"column:oper_ip;type:varchar(64)" json:"operIp"`                                  // 操作IP
	OperLocation  string          `gorm:"column:oper_location;type:varchar(128)" json:"operLocation"`                     // 操作地点
	OperParam     string          `gorm:"column:oper_param;type:text" json:"operParam"`                                   // 请求参数
	JsonResult    string          `gorm:"column:json_result;type:text" json:"jsonResult"`                                 // 返回结果
	Status        string          `gorm:"column:status;type:char(1)" json:"status"`                                       // 操作状态：0成功 1失败
	ErrorMsg      string          `gorm:"column:error_msg;type:text" json:"errorMsg"`                                     // 错误信息
	OperTime      utils.LocalTime `gorm:"column:oper_time;type:timestamptz;index" json:"operTime"`                        // 操作时间
	CostTime      int64           `gorm:"column:cost_time;type:bigint" json:"costTime"`                                   // 耗时（毫秒）
	UserAgent     string          `gorm:"column:user_agent;type:varchar(512)" json:"userAgent"`                           // UA
}

// TableName 返回数据库表名。
func (*OperLog) TableName() string {
	return "s_oper_log"
}

// FindByID 根据ID查询操作日志
func (*OperLog) FindByID(db *gorm.DB, id int64) (*OperLog, error) {
	var log OperLog
	err := db.Where("id = ?", id).First(&log).Error
	if err != nil {
		return nil, err
	}
	return &log, nil
}

// Create 创建操作日志
func (l *OperLog) Create(db *gorm.DB) error {
	return db.Create(l).Error
}

// Update 更新操作日志
func (l *OperLog) Update(db *gorm.DB, id int64, updates map[string]interface{}) error {
	return db.Model(&OperLog{}).Where("id = ?", id).Updates(updates).Error
}

// Delete 删除操作日志（物理删除）
func (*OperLog) Delete(db *gorm.DB, id int64) error {
	return db.Where("id = ?", id).Delete(&OperLog{}).Error
}

// BatchDelete 批量删除操作日志
func (*OperLog) BatchDelete(db *gorm.DB, ids []int64) (int64, error) {
	result := db.Where("id IN ?", ids).Delete(&OperLog{})
	return result.RowsAffected, result.Error
}

// IsSuccess 判断操作是否成功
func (l *OperLog) IsSuccess() bool {
	return l.Status == "0"
}

// CleanOldLogs 清理指定天数之前的日志
func (*OperLog) CleanOldLogs(db *gorm.DB, days int) (int64, error) {
	result := db.Where("oper_time < ?", time.Now().AddDate(0, 0, -days)).Delete(&OperLog{})
	return result.RowsAffected, result.Error
}
