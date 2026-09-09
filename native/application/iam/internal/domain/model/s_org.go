package model

import (
	"fmt"

	"github.com/gcc798/lightning/internal/utils"

	"gorm.io/gorm"
)

// Org 系统组织表（多租户）
type Org struct {
	ID          int64           `gorm:"column:id;type:bigint;primaryKey;autoIncrement:false" autogen:"int64" json:"id"` // 组织ID（使用分布式ID）
	ParentId    int64           `gorm:"column:parent_id;type:bigint;default:0;index" json:"parentId"`                   // 父组织ID（0表示根组织）
	Ancestors   string          `gorm:"column:ancestors;type:varchar(512)" json:"ancestors"`                            // 祖级列表（逗号分隔，例如: "0,1,2"）
	OrgName     string          `gorm:"column:org_name;type:varchar(128);not null" json:"orgName"`                      // 组织名称
	OrgCode     string          `gorm:"column:org_code;type:varchar(64);uniqueIndex" json:"orgCode"`                    // 组织编码（唯一）
	OrgType     string          `gorm:"column:org_type;type:varchar(32);default:'company'" json:"orgType"`              // 组织类型：company公司 department部门 group集团
	Leader      string          `gorm:"column:leader;type:varchar(64)" json:"leader"`                                   // 负责人
	Phone       string          `gorm:"column:phone;type:varchar(32)" json:"phone"`                                     // 联系电话
	Email       string          `gorm:"column:email;type:varchar(128)" json:"email"`                                    // 邮箱
	Status      int32           `gorm:"column:status;type:smallint;default:0" json:"status"`                            // 状态：0正常 1停用
	Sort        int64           `gorm:"column:sort;type:bigint;default:0" json:"sort"`                                  // 显示顺序
	Remark      string          `gorm:"column:remark;type:varchar(500)" json:"remark"`                                  // 备注
	CreateBy    int64           `gorm:"column:create_by;type:bigint" json:"createBy"`                                   // 创建人
	UpdateBy    int64           `gorm:"column:update_by;type:bigint" json:"updateBy"`                                   // 更新人
	CreatedTime utils.LocalTime `gorm:"column:created_time;type:timestamptz;autoCreateTime" json:"createdTime"`
	UpdatedTime utils.LocalTime `gorm:"column:updated_time;type:timestamptz;autoUpdateTime" json:"updatedTime"`
}

// TableName 返回数据库表名。
func (*Org) TableName() string { return "s_org" }

// FindByOrgId 根据组织ID查询组织
func (o *Org) FindByOrgId(db *gorm.DB, orgId int64) (*Org, error) {
	var org Org
	err := db.Where("id = ?", orgId).First(&org).Error
	return &org, err
}

// FindByID 根据ID查询组织
func (o *Org) FindByID(db *gorm.DB, orgId int64) (*Org, error) {
	var org Org
	err := db.Where("id = ?", orgId).First(&org).Error
	if err != nil {
		return nil, err
	}
	return &org, nil
}

// FindByOrgCode 根据组织编码查询组织
func (o *Org) FindByOrgCode(db *gorm.DB, orgCode string) (*Org, error) {
	var org Org
	err := db.Where("org_code = ?", orgCode).First(&org).Error
	return &org, err
}

// FindChildren 查询子组织列表
func (o *Org) FindChildren(db *gorm.DB, parentId int64) ([]Org, error) {
	var orgs []Org
	err := db.Where("parent_id = ?", parentId).Order("sort ASC").Find(&orgs).Error
	return orgs, err
}

// CheckOrgCodeExists 检查组织编码是否存在
func (o *Org) CheckOrgCodeExists(db *gorm.DB, orgCode string) (bool, error) {
	var count int64
	err := db.Model(&Org{}).Where("org_code = ?", orgCode).Count(&count).Error
	return count > 0, err
}

// CheckOrgCodeExistsExcludingSelf 检查组织编码是否被其他组织占用
func (o *Org) CheckOrgCodeExistsExcludingSelf(db *gorm.DB, orgId int64, orgCode string) (bool, error) {
	var count int64
	err := db.Model(&Org{}).Where("org_code = ? AND id != ?", orgCode, orgId).Count(&count).Error
	return count > 0, err
}

// CountChildren 统计子组织数量
func (o *Org) CountChildren(db *gorm.DB, orgId int64) (int64, error) {
	var count int64
	err := db.Model(&Org{}).Where("parent_id = ?", orgId).Count(&count).Error
	return count, err
}

// CountUsers 统计组织下的用户数量
func (o *Org) CountUsers(db *gorm.DB, orgId int64) (int64, error) {
	var count int64
	err := db.Model(&User{}).Where("org_id = ?", orgId).Count(&count).Error
	return count, err
}

// Create 创建组织
func (o *Org) Create(db *gorm.DB, org *Org) error {
	return db.Create(org).Error
}

// Update 更新组织
func (o *Org) Update(db *gorm.DB, org *Org) error {
	return db.Save(org).Error
}

// Delete 删除组织
func (o *Org) Delete(db *gorm.DB, orgId int64) error {
	return db.Where("id = ?", orgId).Delete(&Org{}).Error
}

// List 分页查询组织列表
func (o *Org) List(db *gorm.DB, offset, limit int, orgName, orgCode string, status int32, parentId *int64) ([]Org, int64, error) {
	var orgs []Org
	var total int64

	query := db.Model(&Org{})

	// 条件过滤
	if orgName != "" {
		query = query.Where("org_name LIKE ?", "%"+orgName+"%")
	}
	if orgCode != "" {
		query = query.Where("org_code LIKE ?", "%"+orgCode+"%")
	}
	if status != 0 {
		query = query.Where("status = ?", status)
	}
	if parentId != nil {
		query = query.Where("parent_id = ?", *parentId)
	}

	// 统计总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	if err := query.Offset(offset).Limit(limit).Order("sort ASC, created_time DESC").Find(&orgs).Error; err != nil {
		return nil, 0, err
	}

	return orgs, total, nil
}

// FindAll 查询所有组织（用于构建树）
func (o *Org) FindAll(db *gorm.DB) ([]Org, error) {
	var orgs []Org
	err := db.Order("sort ASC, created_time DESC").Find(&orgs).Error
	return orgs, err
}

// BuildAncestors 构建祖级列表
func (o *Org) BuildAncestors(db *gorm.DB, parentId int64) (string, error) {
	if parentId == 0 {
		return "0", nil
	}

	parent, err := o.FindByID(db, parentId)
	if err != nil {
		return "", err
	}

	return parent.Ancestors + "," + fmt.Sprint(parentId), nil
}

// IsActive 判断组织是否激活
func (o *Org) IsActive() bool {
	return o.Status == 0
}

// HasChildren 判断是否有子组织
func (o *Org) HasChildren(db *gorm.DB) (bool, error) {
	count, err := o.CountChildren(db, o.ID)
	return count > 0, err
}

// HasUsers 判断是否有关联用户
func (o *Org) HasUsers(db *gorm.DB) (bool, error) {
	count, err := o.CountUsers(db, o.ID)
	return count > 0, err
}
