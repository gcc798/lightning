package model

import (
	"github.com/gcc798/lightning/internal/utils"

	"gorm.io/gorm"
)

// User 系统用户
type User struct {
	ID          int64           `gorm:"column:id;type:bigint;primaryKey;autoIncrement:false" autogen:"int64" json:"id"`                                    // 用户ID（使用分布式ID）
	UserName    string          `gorm:"column:user_name;type:varchar(64);uniqueIndex;not null" json:"userName"`                                            // 用户名（登录账号）
	NickName    string          `gorm:"column:nick_name;type:varchar(64)" json:"nickName"`                                                                 // 昵称（显示名称）
	UserType    int32           `gorm:"column:user_type;type:smallint;default:0" json:"userType"`                                                          // 用户类型：0系统用户 1微信用户 2APP用户
	OrgID       int64           `gorm:"column:org_id;type:bigint;default:0;index" json:"orgId"`                                                            // 组织ID
	Email       string          `gorm:"column:email;type:varchar(128);uniqueIndex:idx_s_user_email,where:email <> ''" json:"email"`                        // 邮箱
	Phonenumber string          `gorm:"column:phonenumber;type:varchar(32);uniqueIndex:idx_s_user_phonenumber,where:phonenumber <> ''" json:"phonenumber"` // 手机号
	Sex         int32           `gorm:"column:sex;type:smallint;default:2" json:"sex"`                                                                     // 性别：0男 1女 2未知
	Avatar      string          `gorm:"column:avatar;type:varchar(512)" json:"avatar"`                                                                     // 头像URL
	Password    string          `gorm:"column:password;type:varchar(255)" json:"-"`                                                                        // 密码（加密）
	Status      int32           `gorm:"column:status;type:smallint;default:0" json:"status"`                                                               // 状态：0正常 1停用
	Sort        int64           `gorm:"column:sort;type:bigint;default:0" json:"sort"`                                                                     // 排序字段
	LoginIp     string          `gorm:"column:login_ip;type:varchar(64)" json:"loginIp"`                                                                   // 最后登录IP
	LoginDate   int64           `gorm:"column:login_date;type:bigint" json:"loginDate"`                                                                    // 最后登录时间（时间戳）
	OpenId      string          `gorm:"column:open_id;type:varchar(128);uniqueIndex:idx_s_user_open_id,where:open_id <> ''" json:"openId"`                 // 微信OpenID
	UnionId     string          `gorm:"column:union_id;type:varchar(128);uniqueIndex:idx_s_user_union_id,where:union_id <> ''" json:"unionId"`             // 微信UnionID
	Remark      string          `gorm:"column:remark;type:varchar(500)" json:"remark"`                                                                     // 备注
	CreateBy    int64           `gorm:"column:create_by;type:bigint" json:"createBy"`                                                                      // 创建人
	UpdateBy    int64           `gorm:"column:update_by;type:bigint" json:"updateBy"`                                                                      // 更新人
	CreatedTime utils.LocalTime `gorm:"column:created_time;type:timestamptz;autoCreateTime" json:"createdTime"`
	UpdatedTime utils.LocalTime `gorm:"column:updated_time;type:timestamptz;autoUpdateTime" json:"updatedTime"`
}

// TableName 返回数据库表名。
func (*User) TableName() string { return "s_user" }

// FindByUsername 执行业务逻辑。
func (u *User) FindByUsername(db *gorm.DB, username string) (*User, error) {
	var out User
	tx := db.Where("user_name = ?", username).Limit(1).Find(&out)
	if tx.Error != nil {
		return nil, tx.Error
	}
	if tx.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &out, nil
}

// FindByOpenId 执行业务逻辑。
func (u *User) FindByOpenId(db *gorm.DB, openId string) (*User, error) {
	var out User
	tx := db.Where("open_id = ?", openId).Limit(1).Find(&out)
	if tx.Error != nil {
		return nil, tx.Error
	}
	if tx.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &out, nil
}

// FindByPhonenumber 执行业务逻辑。
func (u *User) FindByPhonenumber(db *gorm.DB, phonenumber string) (*User, error) {
	var out User
	tx := db.Where("phonenumber = ?", phonenumber).Limit(1).Find(&out)
	if tx.Error != nil {
		return nil, tx.Error
	}
	if tx.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &out, nil
}

// FindByEmail 根据邮箱查询用户
func (u *User) FindByEmail(db *gorm.DB, email string) (*User, error) {
	var out User
	tx := db.Where("email = ?", email).Limit(1).Find(&out)
	if tx.Error != nil {
		return nil, tx.Error
	}
	if tx.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &out, nil
}

// Create 创建业务数据。
func (u *User) Create(db *gorm.DB, nu *User) error {
	return db.Create(nu).Error
}

// UpdateLoginInfo 更新登录信息（IP与时间戳）
func (u *User) UpdateLoginInfo(db *gorm.DB, userId int64, ip string, ts int64) error {
	return db.Model(&User{}).Where("id = ?", userId).Updates(map[string]any{
		"login_ip":   ip,
		"login_date": ts,
	}).Error
}

// FindConflicts 查找用户名、手机号、邮箱冲突（用于创建时的唯一性校验）
// 一次查询返回所有冲突的用户，避免多次数据库查询
func (u *User) FindConflicts(db *gorm.DB, username, phone, email string) ([]User, error) {
	var users []User
	query := db.Model(&User{})

	conditions := make([]string, 0, 3)
	args := make([]interface{}, 0, 3)

	if username != "" {
		conditions = append(conditions, "user_name = ?")
		args = append(args, username)
	}
	if phone != "" {
		conditions = append(conditions, "phonenumber = ?")
		args = append(args, phone)
	}
	if email != "" {
		conditions = append(conditions, "email = ?")
		args = append(args, email)
	}

	if len(conditions) == 0 {
		return nil, nil
	}

	// 使用 OR 条件查询所有可能的冲突
	whereClause := conditions[0]
	for i := 1; i < len(conditions); i++ {
		whereClause += " OR " + conditions[i]
	}

	err := query.Where(whereClause, args...).Find(&users).Error
	return users, err
}

// FindConflictsExcludingSelf 查找冲突但排除自己（用于更新时的唯一性校验）
func (u *User) FindConflictsExcludingSelf(db *gorm.DB, userId int64, username, phone, email string) ([]User, error) {
	var users []User
	query := db.Model(&User{}).Where("id != ?", userId)

	conditions := make([]string, 0, 3)
	args := make([]interface{}, 0, 3)

	if username != "" {
		conditions = append(conditions, "user_name = ?")
		args = append(args, username)
	}
	if phone != "" {
		conditions = append(conditions, "phonenumber = ?")
		args = append(args, phone)
	}
	if email != "" {
		conditions = append(conditions, "email = ?")
		args = append(args, email)
	}

	if len(conditions) == 0 {
		return nil, nil
	}

	whereClause := conditions[0]
	for i := 1; i < len(conditions); i++ {
		whereClause += " OR " + conditions[i]
	}

	err := query.Where(whereClause, args...).Find(&users).Error
	return users, err
}

// FindByID 根据ID查询用户
func (u *User) FindByID(db *gorm.DB, userId int64) (*User, error) {
	var user User
	err := db.Where("id = ?", userId).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Update 更新用户信息
func (u *User) Update(db *gorm.DB, userId int64, updates map[string]interface{}) error {
	return db.Model(&User{}).Where("id = ?", userId).Updates(updates).Error
}

// Delete 删除用户
func (u *User) Delete(db *gorm.DB, userId int64) error {
	return db.Where("id = ?", userId).Delete(&User{}).Error
}

// BatchDelete 批量删除用户
func (u *User) BatchDelete(db *gorm.DB, userIds []int64) (int64, error) {
	result := db.Where("id IN ?", userIds).Delete(&User{})
	return result.RowsAffected, result.Error
}

// List 分页查询用户列表
func (u *User) List(db *gorm.DB, offset, limit int, username, phonenumber string, status int32) ([]User, int64, error) {
	var users []User
	var total int64

	query := db.Model(&User{})

	// 条件过滤
	if username != "" {
		query = query.Where("user_name LIKE ?", "%"+username+"%")
	}
	if phonenumber != "" {
		query = query.Where("phonenumber LIKE ?", "%"+phonenumber+"%")
	}
	if status >= 0 {
		query = query.Where("status = ?", status)
	}

	// 统计总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	if err := query.Offset(offset).Limit(limit).Order("sort ASC, created_time DESC").Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// UpdatePassword 更新密码
func (u *User) UpdatePassword(db *gorm.DB, userId int64, hashedPassword string) error {
	return db.Model(&User{}).Where("id = ?", userId).Update("password", hashedPassword).Error
}

// ClearPassword 清空密码字段（用于返回给前端）
func (u *User) ClearPassword() {
	u.Password = ""
}

// IsActive 判断用户是否激活
func (u *User) IsActive() bool {
	return u.Status == 0
}

// CanLogin 判断用户是否可以登录
func (u *User) CanLogin() bool {
	return u.IsActive()
}
