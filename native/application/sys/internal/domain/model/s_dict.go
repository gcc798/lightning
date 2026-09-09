package model

import (
	"github.com/gcc798/lightning/internal/utils"

	"gorm.io/gorm"
)

// DictData 字典数据
type DictData struct {
	ID          int64           `gorm:"column:id;type:bigint;primaryKey;autoIncrement:false" autogen:"int64" json:"id"` // 字典编码（使用分布式ID）
	ParentId    int64           `gorm:"column:parent_id;type:bigint;default:0;index" json:"parentId"`                   // 父字典ID（0表示根节点）
	Sort        int64           `gorm:"column:sort;type:bigint;default:0" json:"sort"`                                  // 字典排序
	DictLabel   string          `gorm:"column:dict_label;type:varchar(128)" json:"dictLabel"`                           // 字典标签
	DictValue   string          `gorm:"column:dict_value;type:varchar(128)" json:"dictValue"`                           // 字典键值
	DictType    string          `gorm:"column:dict_type;type:varchar(128);index" json:"dictType"`                       // 字典类型
	IsDefault   bool            `gorm:"column:is_default;type:boolean;default:false" json:"isDefault"`                  // 是否默认
	Status      int32           `gorm:"column:status;type:smallint;default:0" json:"status"`                            // 状态：0正常 1停用
	Remark      string          `gorm:"column:remark;type:varchar(500)" json:"remark"`                                  // 备注
	CreateBy    int64           `gorm:"column:create_by;type:bigint" json:"createBy"`                                   // 创建者
	CreatedTime utils.LocalTime `gorm:"column:created_time;type:timestamptz;autoCreateTime" json:"createdTime"`         // 创建时间
	UpdateBy    int64           `gorm:"column:update_by;type:bigint" json:"updateBy"`                                   // 更新者
	UpdatedTime utils.LocalTime `gorm:"column:updated_time;type:timestamptz;autoUpdateTime" json:"updatedTime"`         // 更新时间
}

// TableName 返回数据库表名。
func (*DictData) TableName() string {
	return "s_dict_data"
}

// FindByID 根据ID查询字典
func (*DictData) FindByID(db *gorm.DB, id int64) (*DictData, error) {
	var dict DictData
	err := db.Where("id = ?", id).First(&dict).Error
	if err != nil {
		return nil, err
	}
	return &dict, nil
}

// FindByType 根据字典类型查询字典数据列表
func (*DictData) FindByType(db *gorm.DB, dictType string) ([]DictData, error) {
	var dicts []DictData
	err := db.Where("dict_type = ? AND status = 0", dictType).
		Order("sort ASC, id ASC").
		Find(&dicts).Error
	return dicts, err
}

// FindByTypeAndParent 根据字典类型和父ID查询子字典列表
func (*DictData) FindByTypeAndParent(db *gorm.DB, dictType string, parentId int64) ([]DictData, error) {
	var dicts []DictData
	err := db.Where("dict_type = ? AND parent_id = ? AND status = 0", dictType, parentId).
		Order("sort ASC, id ASC").
		Find(&dicts).Error
	return dicts, err
}

// FindChildren 查询子字典列表
func (*DictData) FindChildren(db *gorm.DB, parentId int64) ([]DictData, error) {
	var dicts []DictData
	err := db.Where("parent_id = ?", parentId).
		Order("sort ASC, id ASC").
		Find(&dicts).Error
	return dicts, err
}

// CountChildren 统计子字典数量
func (*DictData) CountChildren(db *gorm.DB, parentId int64) (int64, error) {
	var count int64
	err := db.Model(&DictData{}).Where("parent_id = ?", parentId).Count(&count).Error
	return count, err
}

// CheckDictValueExists 检查字典值是否存在（同类型下）
func (*DictData) CheckDictValueExists(db *gorm.DB, dictType, dictValue string) (bool, error) {
	var count int64
	err := db.Model(&DictData{}).
		Where("dict_type = ? AND dict_value = ?", dictType, dictValue).
		Count(&count).Error
	return count > 0, err
}

// CheckDictValueExistsExcludingSelf 检查字典值是否被其他字典占用
func (*DictData) CheckDictValueExistsExcludingSelf(db *gorm.DB, id int64, dictType, dictValue string) (bool, error) {
	var count int64
	err := db.Model(&DictData{}).
		Where("dict_type = ? AND dict_value = ? AND id != ?", dictType, dictValue, id).
		Count(&count).Error
	return count > 0, err
}

// GetDictLabel 根据字典类型和键值获取标签
func (*DictData) GetDictLabel(db *gorm.DB, dictType, dictValue string) (string, error) {
	var dict DictData
	err := db.Where("dict_type = ? AND dict_value = ? AND status = 0", dictType, dictValue).
		First(&dict).Error
	if err != nil {
		return "", err
	}
	return dict.DictLabel, nil
}

// GetDictValue 根据字典类型和标签获取键值
func (*DictData) GetDictValue(db *gorm.DB, dictType, dictLabel string) (string, error) {
	var dict DictData
	err := db.Where("dict_type = ? AND dict_label = ? AND status = 0", dictType, dictLabel).
		First(&dict).Error
	if err != nil {
		return "", err
	}
	return dict.DictValue, nil
}

// GetDefaultDict 获取指定类型的默认字典项
func (*DictData) GetDefaultDict(db *gorm.DB, dictType string) (*DictData, error) {
	var dict DictData
	err := db.Where("dict_type = ? AND is_default = ? AND status = 0", dictType, true).
		First(&dict).Error
	if err != nil {
		return nil, err
	}
	return &dict, nil
}

// Create 创建字典数据
func (d *DictData) Create(db *gorm.DB) error {
	return db.Create(d).Error
}

// Update 更新字典数据
func (d *DictData) Update(db *gorm.DB, id int64, updates map[string]interface{}) error {
	return db.Model(&DictData{}).Where("id = ?", id).Updates(updates).Error
}

// Delete 删除字典数据
func (*DictData) Delete(db *gorm.DB, id int64) error {
	return db.Where("id = ?", id).Delete(&DictData{}).Error
}

// BatchDelete 批量删除字典数据
func (*DictData) BatchDelete(db *gorm.DB, ids []int64) (int64, error) {
	result := db.Where("id IN ?", ids).Delete(&DictData{})
	return result.RowsAffected, result.Error
}

// List 分页查询字典列表
func (*DictData) List(db *gorm.DB, offset, limit int, dictType, dictLabel string, status int32) ([]DictData, int64, error) {
	var dicts []DictData
	var total int64

	query := db.Model(&DictData{})

	// 条件过滤
	if dictType != "" {
		query = query.Where("dict_type = ?", dictType)
	}
	if dictLabel != "" {
		query = query.Where("dict_label LIKE ?", "%"+dictLabel+"%")
	}
	if status >= 0 {
		query = query.Where("status = ?", status)
	}

	// 统计总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	if err := query.Offset(offset).Limit(limit).
		Order("dict_type ASC, sort ASC, id ASC").
		Find(&dicts).Error; err != nil {
		return nil, 0, err
	}

	return dicts, total, nil
}

// IsActive 判断字典是否激活
func (d *DictData) IsActive() bool {
	return d.Status == 0
}

// IsDefaultDict 判断是否为默认字典
func (d *DictData) IsDefaultDict() bool {
	return d.IsDefault
}

// HasChildren 判断是否有子字典
func (d *DictData) HasChildren(db *gorm.DB) (bool, error) {
	count, err := (&DictData{}).CountChildren(db, d.ID)
	return count > 0, err
}
