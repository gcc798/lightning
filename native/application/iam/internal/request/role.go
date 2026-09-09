package request

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/gcc798/lightning/internal/utils/pagination"
)

// Int64ID 兼容前端雪花 ID 字符串与数字两种 JSON 传参形式。
type Int64ID int64

// UnmarshalJSON 执行业务逻辑。
func (id *Int64ID) UnmarshalJSON(data []byte) error {
	raw := strings.TrimSpace(string(data))
	if raw == "" || raw == "null" {
		*id = 0
		return nil
	}

	var text string
	if err := json.Unmarshal(data, &text); err == nil {
		value, err := strconv.ParseInt(strings.TrimSpace(text), 10, 64)
		if err != nil {
			return fmt.Errorf("无效ID: %s", text)
		}
		*id = Int64ID(value)
		return nil
	}

	value, err := strconv.ParseInt(raw, 10, 64)
	if err != nil {
		return fmt.Errorf("无效ID: %s", raw)
	}
	*id = Int64ID(value)
	return nil
}

// Int64 返回 int64 类型 ID。
func (id Int64ID) Int64() int64 {
	return int64(id)
}

// PageRoleRequest 查询角色列表请求
type PageRoleRequest struct {
	pagination.PageQuery        // 嵌入分页参数
	RoleName             string `json:"roleName" example:"管理员"` // 角色名称（模糊查询）
	Status               int32  `json:"status" example:"0"`     // 状态：0正常 1停用，-1表示不过滤
}

// CreateRoleRequest 创建角色请求
type CreateRoleRequest struct {
	RoleKey   string `json:"roleKey" binding:"required" msg:"角色标识不能为空" example:"user_manager"` // 角色标识（唯一）
	RoleName  string `json:"roleName" binding:"required" msg:"角色名称不能为空" example:"用户管理员"`       // 角色名称
	Sort      int64  `json:"sort" example:"1"`                                                 // 显示顺序
	Status    int32  `json:"status" example:"0"`                                               // 状态：0正常 1停用
	DataScope int32  `json:"dataScope" example:"2"`                                            // 数据范围：1全部 2自定义 3本组织 4本组织及以下 5仅本人
	Remark    string `json:"remark" example:"负责用户管理"`                                          // 备注
}

// UpdateRoleRequest 更新角色请求
type UpdateRoleRequest struct {
	RoleId    int64  `json:"roleId" binding:"required" msg:"角色ID不能为空" example:"1"`       // 角色ID
	RoleName  string `json:"roleName" binding:"required" msg:"角色名称不能为空" example:"用户管理员"` // 角色名称
	Sort      int64  `json:"sort" example:"1"`                                           // 显示顺序
	Status    int32  `json:"status" example:"0"`                                         // 状态：0正常 1停用
	DataScope int32  `json:"dataScope" example:"2"`                                      // 数据范围
	Remark    string `json:"remark" example:"负责用户管理"`                                    // 备注
}

// AssignRoleToUserRequest 为用户分配角色请求
type AssignRoleToUserRequest struct {
	UserId Int64ID `json:"userId" binding:"required" msg:"用户ID不能为空" example:"1001"` // 用户ID
	RoleId Int64ID `json:"roleId" binding:"required" msg:"角色ID不能为空" example:"1"`    // 角色ID
}

// BatchRoleUsersRequest 批量调整角色用户请求
type BatchRoleUsersRequest struct {
	UserIds []Int64ID `json:"userIds" binding:"required" msg:"用户ID不能为空" example:"1001,1002"` // 用户ID列表
}

// AssignRoleMenusRequest 分配角色菜单请求
type AssignRoleMenusRequest struct {
	MenuIds []Int64ID `json:"menuIds" msg:"菜单ID列表" example:"1001,1002"` // 菜单ID列表
}
