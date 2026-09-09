package iam

import "github.com/gcc798/lightning/internal/utils/pagination"

// CreateUserRequest 创建用户请求
type CreateUserRequest struct {
	UserName string `json:"userName" binding:"required,min=3,max=20" msg:"用户名长度必须是3-20个字符"`
	NickName string `json:"nickName" binding:"required" msg:"请输入昵称"`
	Password string `json:"password" binding:"required,min=6" msg:"密码长度不能少于6位"`
	UserType int32  `json:"userType"` // 用户类型：0系统用户 1微信用户 2APP用户
	//Email       string `json:"email" binding:"omitempty,email" msg:"邮箱格式不正确，请输入有效的邮箱地址"`
	Email       string `json:"email" binding:"omitempty,email"`
	Phonenumber string `json:"phonenumber" binding:"omitempty,len=11" msg:"手机号必须是11位数字"`
	Sex         int32  `json:"sex" binding:"omitempty,oneof=0 1 2"` // 性别：0男 1女 2未知
	Avatar      string `json:"avatar"`
	Status      int32  `json:"status" binding:"omitempty,oneof=0 1"` // 状态：0正常 1停用
	Remark      string `json:"remark"`
	CreateBy    int64  `json:"-"` // 从上下文获取，不从 JSON 解析
	UpdateBy    int64  `json:"-"` // 从上下文获取，不从 JSON 解析
}

// UpdateUserRequest 更新用户请求
type UpdateUserRequest struct {
	UserId      int64  `json:"-"` // 从路径参数获取，不从 JSON 解析
	UserName    string `json:"userName" binding:"omitempty,min=3,max=20"`
	NickName    string `json:"nickName"`
	UserType    int32  `json:"userType"` // 用户类型：0系统用户 1微信用户 2APP用户
	Email       string `json:"email" binding:"omitempty,email"`
	Phonenumber string `json:"phonenumber" binding:"omitempty,len=11"`
	Sex         int32  `json:"sex" binding:"omitempty,oneof=0 1 2"` // 性别：0男 1女 2未知
	Avatar      string `json:"avatar"`
	Status      int32  `json:"status" binding:"omitempty,oneof=0 1"` // 状态：0正常 1停用
	Remark      string `json:"remark"`
	UpdateBy    int64  `json:"-"` // 从上下文获取，不从 JSON 解析
}

// BatchDeleteUsersRequest 批量删除用户请求
type BatchDeleteUsersRequest struct {
	IDs []int64 `json:"ids" binding:"required,min=1"`
}

// PageUsersRequest 查询用户列表请求
type PageUsersRequest struct {
	pagination.PageQuery        // 嵌入分页参数
	UserName             string `json:"username"`
	Phonenumber          string `json:"phonenumber"`
	Status               int32  `json:"status" binding:"omitempty,oneof=0 1"` // 状态：0正常 1停用
}

// BatchImportUsersRequest 批量导入用户请求
type BatchImportUsersRequest struct {
	Users []CreateUserRequest `json:"users" binding:"required,min=1,dive"`
}

// ResetPasswordRequest 重置密码请求
type ResetPasswordRequest struct {
	NewPassword string `json:"newPassword" binding:"required,min=6"`
}

// ChangePasswordRequest 用户修改密码请求
type ChangePasswordRequest struct {
	OldPassword string `json:"oldPassword" binding:"required,min=6"`
	NewPassword string `json:"newPassword" binding:"required,min=6"`
}
