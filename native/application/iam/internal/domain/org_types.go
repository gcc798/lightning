package iam

import "github.com/gcc798/microservice-kit/internal/utils/pagination"

// CreateOrgRequest 创建组织请求
type CreateOrgRequest struct {
	ParentId int64  `json:"parentId"`                                                                                             // 父组织ID，0表示根组织
	OrgName  string `json:"orgName" binding:"required,min=2,max=50" msg:"组织名称必须是2-50个字符"`                                         // 组织名称
	OrgCode  string `json:"orgCode" binding:"required,min=2,max=30" msg:"组织编码必须是2-30个字符"`                                         // 组织编码
	OrgType  string `json:"orgType" binding:"omitempty,oneof=company department group" msg:"组织类型必须是 company/department/group 之一"` // 组织类型
	Leader   string `json:"leader" binding:"omitempty,max=50" msg:"负责人不能超过50个字符"`                                                 // 负责人
	Phone    string `json:"phone" binding:"omitempty,min=11,max=11" msg:"联系电话必须是11位数字"`                                           // 联系电话
	Email    string `json:"email" binding:"omitempty,email" msg:"邮箱格式不正确"`                                                        // 邮箱
	Status   int32  `json:"status" binding:"omitempty,oneof=0 1" msg:"状态必须是0（正常）或1（停用）"`                                          // 状态：0正常 1停用
	Sort     int64  `json:"sort"`                                                                                                 // 显示顺序
	Remark   string `json:"remark" binding:"omitempty,max=500" msg:"备注不能超过500个字符"`                                                // 备注
	CreateBy int64  `json:"-"`                                                                                                    // 从上下文获取
	UpdateBy int64  `json:"-"`                                                                                                    // 从上下文获取
}

// UpdateOrgRequest 更新组织请求
type UpdateOrgRequest struct {
	OrgId    int64  `json:"-"`                                                                                                    // 从路径参数获取
	ParentId int64  `json:"parentId"`                                                                                             // 父组织ID
	OrgName  string `json:"orgName" binding:"omitempty,min=2,max=50" msg:"组织名称必须是2-50个字符"`                                        // 组织名称
	OrgCode  string `json:"orgCode" binding:"omitempty,min=2,max=30" msg:"组织编码必须是2-30个字符"`                                        // 组织编码
	OrgType  string `json:"orgType" binding:"omitempty,oneof=company department group" msg:"组织类型必须是 company/department/group 之一"` // 组织类型
	Leader   string `json:"leader" binding:"omitempty,max=50" msg:"负责人不能超过50个字符"`                                                 // 负责人
	Phone    string `json:"phone" binding:"omitempty,min=11,max=11" msg:"联系电话必须是11位数字"`                                           // 联系电话
	Email    string `json:"email" binding:"omitempty,email" msg:"邮箱格式不正确"`                                                        // 邮箱
	Status   int32  `json:"status" binding:"omitempty,oneof=0 1" msg:"状态必须是0（正常）或1（停用）"`                                          // 状态：0正常 1停用
	Sort     int64  `json:"sort"`                                                                                                 // 显示顺序
	Remark   string `json:"remark" binding:"omitempty,max=500" msg:"备注不能超过500个字符"`                                                // 备注
	UpdateBy int64  `json:"-"`                                                                                                    // 从上下文获取
}

// BatchDeleteOrgsRequest 批量删除组织请求
type BatchDeleteOrgsRequest struct {
	IDs []int64 `json:"ids" binding:"required,min=1" msg:"请至少选择一个组织"` // 组织ID列表
}

// PageOrgsRequest 查询组织列表请求
type PageOrgsRequest struct {
	pagination.PageQuery        // 嵌入分页参数
	OrgName              string `json:"orgName"`                                                     // 组织名称（模糊查询）
	OrgCode              string `json:"orgCode"`                                                     // 组织编码
	Status               int32  `json:"status" binding:"omitempty,oneof=0 1" msg:"状态必须是0（正常）或1（停用）"` // 状态：0正常 1停用
	ParentId             *int64 `json:"parentId"`                                                    // 父组织ID（可选）
}
