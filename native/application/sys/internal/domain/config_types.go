package sys

import (
	"encoding/json"

	"github.com/gcc798/microservice-kit/internal/utils/pagination"
)

// CreateConfigRequest 创建配置请求
type CreateConfigRequest struct {
	Name     string          `json:"name" binding:"required" msg:"配置名称不能为空"`    // 配置名称
	Code     string          `json:"code" binding:"required" msg:"配置编码不能为空"`    // 配置编码
	Data     json.RawMessage `json:"data"`                                      // 配置数据（JSON格式）
	Remark   string          `json:"remark"`                                    // 备注
	CreateBy int64           `json:"createBy" binding:"required" msg:"创建者不能为空"` // 创建者
	UpdateBy int64           `json:"updateBy" binding:"required" msg:"更新者不能为空"` // 更新者
}

// UpdateConfigRequest 更新配置请求
type UpdateConfigRequest struct {
	ID       int64           `json:"id"`                                        // 配置ID（由路径参数注入）
	Name     string          `json:"name" binding:"required" msg:"配置名称不能为空"`    // 配置名称
	Code     string          `json:"code" binding:"required" msg:"配置编码不能为空"`    // 配置编码
	Data     json.RawMessage `json:"data"`                                      // 配置数据（JSON格式）
	Remark   string          `json:"remark"`                                    // 备注
	UpdateBy int64           `json:"updateBy" binding:"required" msg:"更新者不能为空"` // 更新者
}

// BatchDeleteConfigRequest 批量删除配置请求
type BatchDeleteConfigRequest struct {
	IDs []int64 `json:"ids" binding:"required,min=1" msg:"请至少选择一个配置"` // 配置ID列表
}

// PageConfigRequest 配置列表查询请求
type PageConfigRequest struct {
	pagination.PageQuery        // 嵌入分页参数
	Code                 string `json:"code"` // 配置编码(可选)
	Name                 string `json:"name"` // 配置名称(可选,模糊查询)
}

// GetConfigByCodeRequest 根据编码获取配置请求
type GetConfigByCodeRequest struct {
	Code string `form:"code" query:"code" binding:"required" msg:"配置编码不能为空"` // 配置编码
}
