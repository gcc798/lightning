package constants

// 用户相关枚举
const (
	// 性别
	SexMale    int32 = 0 // 男
	SexFemale  int32 = 1 // 女
	SexUnknown int32 = 2 // 未知

	// 用户类型
	UserTypeSystem int32 = 0 // 系统用户
	UserTypeWechat int32 = 1 // 微信用户
	UserTypeApp    int32 = 2 // APP用户

	// 状态（通用）
	StatusNormal   int32 = 0 // 正常
	StatusDisabled int32 = 1 // 停用
)

// 角色相关枚举
const (
	// 数据范围
	DataScopeAll       int32 = 1 // 全部数据权限
	DataScopeCustom    int32 = 2 // 自定义数据权限
	DataScopeOrg       int32 = 3 // 本组织数据权限
	DataScopeOrgAndSub int32 = 4 // 本组织及以下数据权限
	DataScopeSelf      int32 = 5 // 仅本人数据权限
)

// 菜单相关枚举
const (
	// 菜单类型
	MenuTypeDir    int32 = 0 // 目录
	MenuTypeMenu   int32 = 1 // 菜单
	MenuTypeButton int32 = 2 // 按钮

	// 是否外链
	IsFrameNo  int32 = 0 // 否
	IsFrameYes int32 = 1 // 是

	// 是否缓存
	IsCacheNo  int32 = 0 // 否
	IsCacheYes int32 = 1 // 是

	// 是否可见
	VisibleYes int32 = 0 // 显示
	VisibleNo  int32 = 1 // 隐藏
)

// 登录日志相关枚举
const (
	// 登录状态
	LoginStatusSuccess int32 = 0 // 成功
	LoginStatusFailed  int32 = 1 // 失败
)

// 字典相关枚举
const (
	// 是否默认
	IsDefaultNo  bool = false // 否
	IsDefaultYes bool = true  // 是
)

// 操作日志相关枚举
const (
	// 业务类型
	BusinessTypeOther  = "OTHER"  // 其他
	BusinessTypeQuery  = "QUERY"  // 查询
	BusinessTypeCreate = "CREATE" // 新增
	BusinessTypeUpdate = "UPDATE" // 修改
	BusinessTypeDelete = "DELETE" // 删除
	BusinessTypeExport = "EXPORT" // 导出
	BusinessTypeImport = "IMPORT" // 导入
	BusinessTypeGrant  = "GRANT"  // 授权
	BusinessTypeClean  = "CLEAN"  // 清空
)
